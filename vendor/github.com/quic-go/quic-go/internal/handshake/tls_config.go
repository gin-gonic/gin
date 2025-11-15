package handshake

import (
	"crypto/tls"
	"net"
)

func setupConfigForServer(conf *tls.Config, localAddr, remoteAddr net.Addr) *tls.Config {
	// Workaround for https://github.com/golang/go/issues/60506.
	// This initializes the session tickets _before_ cloning the config.
	_, _ = conf.DecryptTicket(nil, tls.ConnectionState{})

	conf = conf.Clone()
	conf.MinVersion = tls.VersionTLS13

	// The tls.Config contains two callbacks that pass in a tls.ClientHelloInfo.
	// Since crypto/tls doesn't do it, we need to make sure to set the Conn field with a fake net.Conn
	// that allows the caller to get the local and the remote address.
	if conf.GetConfigForClient != nil {
		gcfc := conf.GetConfigForClient
		conf.GetConfigForClient = func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			info.Conn = &conn{localAddr: localAddr, remoteAddr: remoteAddr}
			c, err := gcfc(info)
			if c != nil {
				// we're returning a tls.Config here, so we need to apply this recursively
				c = setupConfigForServer(c, localAddr, remoteAddr)
			}
			return c, err
		}
	}
	if conf.GetCertificate != nil {
		gc := conf.GetCertificate
		conf.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			info.Conn = &conn{localAddr: localAddr, remoteAddr: remoteAddr}
			return gc(info)
		}
	}
	return conf
}
