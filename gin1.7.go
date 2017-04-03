// +build go1.7

package gin

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

// RunAutoTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It obtains and refreshes certificates automatically,
// as well as providing them to a TLS server via tls.Config.
// only from Go version 1.7 onward
func (engine *Engine) RunAutoTLS(addr string, cache string, domain ...string) (err error) {
	debugPrint("Listening and serving HTTPS on %s and host name is %s\n", addr, domain)
	defer func() { debugPrintError(err) }()
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,
	}

	//your domain here
	if len(domain) != 0 {
		m.HostPolicy = autocert.HostWhitelist(domain...)
	}

	// folder for storing certificates
	if cache != "" {
		m.Cache = autocert.DirCache(cache)
	}

	s := &http.Server{
		Addr:      addr,
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   engine,
	}
	err = s.ListenAndServeTLS("", "")
	return
}
