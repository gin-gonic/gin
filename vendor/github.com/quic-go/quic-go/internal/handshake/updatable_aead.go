package handshake

import (
	"crypto"
	"crypto/cipher"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"sync/atomic"

	"github.com/quic-go/quic-go/internal/monotime"
	"github.com/quic-go/quic-go/internal/protocol"
	"github.com/quic-go/quic-go/internal/qerr"
	"github.com/quic-go/quic-go/internal/utils"
	"github.com/quic-go/quic-go/qlog"
	"github.com/quic-go/quic-go/qlogwriter"
)

var keyUpdateInterval atomic.Uint64

func init() {
	keyUpdateInterval.Store(protocol.KeyUpdateInterval)
}

func SetKeyUpdateInterval(v uint64) (reset func()) {
	old := keyUpdateInterval.Swap(v)
	return func() { keyUpdateInterval.Store(old) }
}

// FirstKeyUpdateInterval is the maximum number of packets we send or receive before initiating the first key update.
// It's a package-level variable to allow modifying it for testing purposes.
var FirstKeyUpdateInterval uint64 = 100

type updatableAEAD struct {
	suite *cipherSuite

	keyPhase           protocol.KeyPhase
	largestAcked       protocol.PacketNumber
	firstPacketNumber  protocol.PacketNumber
	handshakeConfirmed bool

	invalidPacketLimit uint64
	invalidPacketCount uint64

	// Time when the keys should be dropped. Keys are dropped on the next call to Open().
	prevRcvAEADExpiry monotime.Time
	prevRcvAEAD       cipher.AEAD

	firstRcvdWithCurrentKey protocol.PacketNumber
	firstSentWithCurrentKey protocol.PacketNumber
	highestRcvdPN           protocol.PacketNumber // highest packet number received (which could be successfully unprotected)
	numRcvdWithCurrentKey   uint64
	numSentWithCurrentKey   uint64
	rcvAEAD                 cipher.AEAD
	sendAEAD                cipher.AEAD
	// caches cipher.AEAD.Overhead(). This speeds up calls to Overhead().
	aeadOverhead int

	nextRcvAEAD           cipher.AEAD
	nextSendAEAD          cipher.AEAD
	nextRcvTrafficSecret  []byte
	nextSendTrafficSecret []byte

	headerDecrypter headerProtector
	headerEncrypter headerProtector

	rttStats *utils.RTTStats

	qlogger qlogwriter.Recorder
	logger  utils.Logger
	version protocol.Version

	// use a single slice to avoid allocations
	nonceBuf []byte
}

var (
	_ ShortHeaderOpener = &updatableAEAD{}
	_ ShortHeaderSealer = &updatableAEAD{}
)

func newUpdatableAEAD(rttStats *utils.RTTStats, qlogger qlogwriter.Recorder, logger utils.Logger, version protocol.Version) *updatableAEAD {
	return &updatableAEAD{
		firstPacketNumber:       protocol.InvalidPacketNumber,
		largestAcked:            protocol.InvalidPacketNumber,
		firstRcvdWithCurrentKey: protocol.InvalidPacketNumber,
		firstSentWithCurrentKey: protocol.InvalidPacketNumber,
		rttStats:                rttStats,
		qlogger:                 qlogger,
		logger:                  logger,
		version:                 version,
	}
}

func (a *updatableAEAD) rollKeys() {
	if a.prevRcvAEAD != nil {
		a.logger.Debugf("Dropping key phase %d ahead of scheduled time. Drop time was: %s", a.keyPhase-1, a.prevRcvAEADExpiry)
		if a.qlogger != nil {
			a.qlogger.RecordEvent(qlog.KeyDiscarded{
				KeyType:  qlog.KeyTypeClient1RTT,
				KeyPhase: a.keyPhase - 1,
			})
			a.qlogger.RecordEvent(qlog.KeyDiscarded{
				KeyType:  qlog.KeyTypeServer1RTT,
				KeyPhase: a.keyPhase - 1,
			})
		}
		a.prevRcvAEADExpiry = 0
	}

	a.keyPhase++
	a.firstRcvdWithCurrentKey = protocol.InvalidPacketNumber
	a.firstSentWithCurrentKey = protocol.InvalidPacketNumber
	a.numRcvdWithCurrentKey = 0
	a.numSentWithCurrentKey = 0
	a.prevRcvAEAD = a.rcvAEAD
	a.rcvAEAD = a.nextRcvAEAD
	a.sendAEAD = a.nextSendAEAD

	a.nextRcvTrafficSecret = a.getNextTrafficSecret(a.suite.Hash, a.nextRcvTrafficSecret)
	a.nextSendTrafficSecret = a.getNextTrafficSecret(a.suite.Hash, a.nextSendTrafficSecret)
	a.nextRcvAEAD = createAEAD(a.suite, a.nextRcvTrafficSecret, a.version)
	a.nextSendAEAD = createAEAD(a.suite, a.nextSendTrafficSecret, a.version)
}

func (a *updatableAEAD) startKeyDropTimer(now monotime.Time) {
	d := 3 * a.rttStats.PTO(true)
	a.logger.Debugf("Starting key drop timer to drop key phase %d (in %s)", a.keyPhase-1, d)
	a.prevRcvAEADExpiry = now.Add(d)
}

func (a *updatableAEAD) getNextTrafficSecret(hash crypto.Hash, ts []byte) []byte {
	return hkdfExpandLabel(hash, ts, []byte{}, "quic ku", hash.Size())
}

// SetReadKey sets the read key.
// For the client, this function is called before SetWriteKey.
// For the server, this function is called after SetWriteKey.
func (a *updatableAEAD) SetReadKey(suite *cipherSuite, trafficSecret []byte) {
	a.rcvAEAD = createAEAD(suite, trafficSecret, a.version)
	a.headerDecrypter = newHeaderProtector(suite, trafficSecret, false, a.version)
	if a.suite == nil {
		a.setAEADParameters(a.rcvAEAD, suite)
	}

	a.nextRcvTrafficSecret = a.getNextTrafficSecret(suite.Hash, trafficSecret)
	a.nextRcvAEAD = createAEAD(suite, a.nextRcvTrafficSecret, a.version)
}

// SetWriteKey sets the write key.
// For the client, this function is called after SetReadKey.
// For the server, this function is called before SetReadKey.
func (a *updatableAEAD) SetWriteKey(suite *cipherSuite, trafficSecret []byte) {
	a.sendAEAD = createAEAD(suite, trafficSecret, a.version)
	a.headerEncrypter = newHeaderProtector(suite, trafficSecret, false, a.version)
	if a.suite == nil {
		a.setAEADParameters(a.sendAEAD, suite)
	}

	a.nextSendTrafficSecret = a.getNextTrafficSecret(suite.Hash, trafficSecret)
	a.nextSendAEAD = createAEAD(suite, a.nextSendTrafficSecret, a.version)
}

func (a *updatableAEAD) setAEADParameters(aead cipher.AEAD, suite *cipherSuite) {
	a.nonceBuf = make([]byte, aead.NonceSize())
	a.aeadOverhead = aead.Overhead()
	a.suite = suite
	switch suite.ID {
	case tls.TLS_AES_128_GCM_SHA256, tls.TLS_AES_256_GCM_SHA384:
		a.invalidPacketLimit = protocol.InvalidPacketLimitAES
	case tls.TLS_CHACHA20_POLY1305_SHA256:
		a.invalidPacketLimit = protocol.InvalidPacketLimitChaCha
	default:
		panic(fmt.Sprintf("unknown cipher suite %d", suite.ID))
	}
}

func (a *updatableAEAD) DecodePacketNumber(wirePN protocol.PacketNumber, wirePNLen protocol.PacketNumberLen) protocol.PacketNumber {
	return protocol.DecodePacketNumber(wirePNLen, a.highestRcvdPN, wirePN)
}

func (a *updatableAEAD) Open(dst, src []byte, rcvTime monotime.Time, pn protocol.PacketNumber, kp protocol.KeyPhaseBit, ad []byte) ([]byte, error) {
	dec, err := a.open(dst, src, rcvTime, pn, kp, ad)
	if err == ErrDecryptionFailed {
		a.invalidPacketCount++
		if a.invalidPacketCount >= a.invalidPacketLimit {
			return nil, &qerr.TransportError{ErrorCode: qerr.AEADLimitReached}
		}
	}
	if err == nil {
		a.highestRcvdPN = max(a.highestRcvdPN, pn)
	}
	return dec, err
}

func (a *updatableAEAD) open(dst, src []byte, rcvTime monotime.Time, pn protocol.PacketNumber, kp protocol.KeyPhaseBit, ad []byte) ([]byte, error) {
	if a.prevRcvAEAD != nil && !a.prevRcvAEADExpiry.IsZero() && rcvTime.After(a.prevRcvAEADExpiry) {
		a.prevRcvAEAD = nil
		a.logger.Debugf("Dropping key phase %d", a.keyPhase-1)
		a.prevRcvAEADExpiry = 0
		if a.qlogger != nil {
			a.qlogger.RecordEvent(qlog.KeyDiscarded{
				KeyType:  qlog.KeyTypeClient1RTT,
				KeyPhase: a.keyPhase - 1,
			})
			a.qlogger.RecordEvent(qlog.KeyDiscarded{
				KeyType:  qlog.KeyTypeServer1RTT,
				KeyPhase: a.keyPhase - 1,
			})
		}
	}
	binary.BigEndian.PutUint64(a.nonceBuf[len(a.nonceBuf)-8:], uint64(pn))
	if kp != a.keyPhase.Bit() {
		if a.keyPhase > 0 && a.firstRcvdWithCurrentKey == protocol.InvalidPacketNumber || pn < a.firstRcvdWithCurrentKey {
			if a.prevRcvAEAD == nil {
				return nil, ErrKeysDropped
			}
			// we updated the key, but the peer hasn't updated yet
			dec, err := a.prevRcvAEAD.Open(dst, a.nonceBuf, src, ad)
			if err != nil {
				err = ErrDecryptionFailed
			}
			return dec, err
		}
		// try opening the packet with the next key phase
		dec, err := a.nextRcvAEAD.Open(dst, a.nonceBuf, src, ad)
		if err != nil {
			return nil, ErrDecryptionFailed
		}
		// Opening succeeded. Check if the peer was allowed to update.
		if a.keyPhase > 0 && a.firstSentWithCurrentKey == protocol.InvalidPacketNumber {
			return nil, &qerr.TransportError{
				ErrorCode:    qerr.KeyUpdateError,
				ErrorMessage: "keys updated too quickly",
			}
		}
		a.rollKeys()
		a.logger.Debugf("Peer updated keys to %d", a.keyPhase)
		// The peer initiated this key update. It's safe to drop the keys for the previous generation now.
		// Start a timer to drop the previous key generation.
		a.startKeyDropTimer(rcvTime)
		if a.qlogger != nil {
			a.qlogger.RecordEvent(qlog.KeyUpdated{
				Trigger:  qlog.KeyUpdateRemote,
				KeyType:  qlog.KeyTypeClient1RTT,
				KeyPhase: a.keyPhase,
			})
			a.qlogger.RecordEvent(qlog.KeyUpdated{
				Trigger:  qlog.KeyUpdateRemote,
				KeyType:  qlog.KeyTypeServer1RTT,
				KeyPhase: a.keyPhase,
			})
		}
		a.firstRcvdWithCurrentKey = pn
		return dec, err
	}
	// The AEAD we're using here will be the qtls.aeadAESGCM13.
	// It uses the nonce provided here and XOR it with the IV.
	dec, err := a.rcvAEAD.Open(dst, a.nonceBuf, src, ad)
	if err != nil {
		return dec, ErrDecryptionFailed
	}
	a.numRcvdWithCurrentKey++
	if a.firstRcvdWithCurrentKey == protocol.InvalidPacketNumber {
		// We initiated the key updated, and now we received the first packet protected with the new key phase.
		// Therefore, we are certain that the peer rolled its keys as well. Start a timer to drop the old keys.
		if a.keyPhase > 0 {
			a.logger.Debugf("Peer confirmed key update to phase %d", a.keyPhase)
			a.startKeyDropTimer(rcvTime)
		}
		a.firstRcvdWithCurrentKey = pn
	}
	return dec, err
}

func (a *updatableAEAD) Seal(dst, src []byte, pn protocol.PacketNumber, ad []byte) []byte {
	if a.firstSentWithCurrentKey == protocol.InvalidPacketNumber {
		a.firstSentWithCurrentKey = pn
	}
	if a.firstPacketNumber == protocol.InvalidPacketNumber {
		a.firstPacketNumber = pn
	}
	a.numSentWithCurrentKey++
	binary.BigEndian.PutUint64(a.nonceBuf[len(a.nonceBuf)-8:], uint64(pn))
	// The AEAD we're using here will be the qtls.aeadAESGCM13.
	// It uses the nonce provided here and XOR it with the IV.
	return a.sendAEAD.Seal(dst, a.nonceBuf, src, ad)
}

func (a *updatableAEAD) SetLargestAcked(pn protocol.PacketNumber) error {
	if a.firstSentWithCurrentKey != protocol.InvalidPacketNumber &&
		pn >= a.firstSentWithCurrentKey && a.numRcvdWithCurrentKey == 0 {
		return &qerr.TransportError{
			ErrorCode:    qerr.KeyUpdateError,
			ErrorMessage: fmt.Sprintf("received ACK for key phase %d, but peer didn't update keys", a.keyPhase),
		}
	}
	a.largestAcked = pn
	return nil
}

func (a *updatableAEAD) SetHandshakeConfirmed() {
	a.handshakeConfirmed = true
}

func (a *updatableAEAD) updateAllowed() bool {
	if !a.handshakeConfirmed {
		return false
	}
	// the first key update is allowed as soon as the handshake is confirmed
	return a.keyPhase == 0 ||
		// subsequent key updates as soon as a packet sent with that key phase has been acknowledged
		(a.firstSentWithCurrentKey != protocol.InvalidPacketNumber &&
			a.largestAcked != protocol.InvalidPacketNumber &&
			a.largestAcked >= a.firstSentWithCurrentKey)
}

func (a *updatableAEAD) shouldInitiateKeyUpdate() bool {
	if !a.updateAllowed() {
		return false
	}
	// Initiate the first key update shortly after the handshake, in order to exercise the key update mechanism.
	if a.keyPhase == 0 {
		if a.numRcvdWithCurrentKey >= FirstKeyUpdateInterval || a.numSentWithCurrentKey >= FirstKeyUpdateInterval {
			return true
		}
	}
	if a.numRcvdWithCurrentKey >= keyUpdateInterval.Load() {
		a.logger.Debugf("Received %d packets with current key phase. Initiating key update to the next key phase: %d", a.numRcvdWithCurrentKey, a.keyPhase+1)
		return true
	}
	if a.numSentWithCurrentKey >= keyUpdateInterval.Load() {
		a.logger.Debugf("Sent %d packets with current key phase. Initiating key update to the next key phase: %d", a.numSentWithCurrentKey, a.keyPhase+1)
		return true
	}
	return false
}

func (a *updatableAEAD) KeyPhase() protocol.KeyPhaseBit {
	if a.shouldInitiateKeyUpdate() {
		a.rollKeys()
		if a.qlogger != nil {
			a.qlogger.RecordEvent(qlog.KeyUpdated{
				Trigger:  qlog.KeyUpdateLocal,
				KeyType:  qlog.KeyTypeClient1RTT,
				KeyPhase: a.keyPhase,
			})
			a.qlogger.RecordEvent(qlog.KeyUpdated{
				Trigger:  qlog.KeyUpdateLocal,
				KeyType:  qlog.KeyTypeServer1RTT,
				KeyPhase: a.keyPhase,
			})
		}
	}
	return a.keyPhase.Bit()
}

func (a *updatableAEAD) Overhead() int {
	return a.aeadOverhead
}

func (a *updatableAEAD) EncryptHeader(sample []byte, firstByte *byte, hdrBytes []byte) {
	a.headerEncrypter.EncryptHeader(sample, firstByte, hdrBytes)
}

func (a *updatableAEAD) DecryptHeader(sample []byte, firstByte *byte, hdrBytes []byte) {
	a.headerDecrypter.DecryptHeader(sample, firstByte, hdrBytes)
}

func (a *updatableAEAD) FirstPacketNumber() protocol.PacketNumber {
	return a.firstPacketNumber
}
