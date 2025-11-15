package quic

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	extTypeSNI = 0
	extTypeECH = 0xfe0d
)

// findSNIAndECH parses the given byte slice as a ClientHello, and locates:
// - the position and length of the Server Name Indication (SNI) extension,
// - the position of the Encrypted Client Hello (ECH) extension.
// If no SNI extension is found, it returns -1 for the SNI position.
// If no ECH extension is found, it returns -1 for the ECH position.
func findSNIAndECH(data []byte) (sniPos, sniLen, echPos int, err error) {
	if len(data) < 4 {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	if data[0] != 1 {
		return 0, 0, 0, errors.New("not a ClientHello")
	}
	handshakeLen := int(data[1])<<16 | int(data[2])<<8 | int(data[3])
	if len(data) != 4+handshakeLen {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}

	parsePos := 4
	// Skip protocol version (2 bytes)
	if parsePos+2 > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	parsePos += 2
	// skip random (32 bytes)
	if parsePos+32 > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	parsePos += 32
	// session ID
	if parsePos+1 > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	sessionIDLen := int(data[parsePos])
	parsePos++
	if parsePos+sessionIDLen > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	parsePos += sessionIDLen
	// cipher suites
	if parsePos+2 > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	cipherSuitesLen := int(binary.BigEndian.Uint16(data[parsePos:]))
	parsePos += 2
	if parsePos+cipherSuitesLen > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	parsePos += cipherSuitesLen
	// compression methods
	if parsePos+1 > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	compressionMethodsLen := int(data[parsePos])
	parsePos++
	if parsePos+compressionMethodsLen > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	parsePos += compressionMethodsLen

	// extensions
	if parsePos+2 > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	extensionsLen := int(binary.BigEndian.Uint16(data[parsePos:]))
	parsePos += 2
	if parsePos+extensionsLen > len(data) {
		return 0, 0, 0, io.ErrUnexpectedEOF
	}
	extensionsStart := parsePos
	extensions := data[extensionsStart : extensionsStart+extensionsLen]

	// parse extensions
	var extPos int
	sniPos = -1
	echPos = -1
	for extPos+4 <= extensionsLen {
		extType := binary.BigEndian.Uint16(extensions[extPos:])
		extLen := int(binary.BigEndian.Uint16(extensions[extPos+2:]))
		if extPos+4+extLen > extensionsLen {
			return 0, 0, 0, io.ErrUnexpectedEOF
		}
		switch extType {
		case extTypeSNI:
			if sniPos != -1 {
				return 0, 0, 0, errors.New("multiple SNI extensions")
			}
			sniData := extensions[extPos+4 : extPos+4+extLen]
			if len(sniData) < 2 {
				return 0, 0, 0, io.ErrUnexpectedEOF
			}
			nameListLen := int(binary.BigEndian.Uint16(sniData))
			if len(sniData) != 2+nameListLen {
				return 0, 0, 0, io.ErrUnexpectedEOF
			}
			listPos := 2
			for listPos+3 <= nameListLen+2 {
				nameType := sniData[listPos]
				sniLen = int(binary.BigEndian.Uint16(sniData[listPos+1:]))
				if listPos+3+sniLen > len(sniData) {
					return 0, 0, 0, io.ErrUnexpectedEOF
				}
				if nameType == 0 { // host_name
					sniPos = extensionsStart + extPos + 4 + listPos + 3
					break // stop after first host_name
				}
				listPos += 3 + sniLen
			}
			if sniPos == 0 {
				return 0, 0, 0, errors.New("SNI host_name not found")
			}
		case extTypeECH:
			if echPos != -1 {
				return 0, 0, 0, errors.New("multiple ECH extensions")
			}
			echPos = extensionsStart + extPos
		}
		extPos += 4 + extLen
		if sniPos != -1 && echPos != -1 {
			break
		}
	}
	return sniPos, sniLen, echPos, nil
}
