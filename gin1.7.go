// +build go1.7

package gin

import (
	"crypto/tls"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

var AutoTLSManager = autocert.Manager{
	Prompt: autocert.AcceptTOS,
}

// RunAutoTLS attaches the router to a http.Server and starts listening and serving HTTPS (secure) requests.
// It obtains and refreshes certificates automatically,
// as well as providing them to a TLS server via tls.Config.
// only from Go version 1.7 onward
func (engine *Engine) RunAutoTLS(domain ...string) (err error) {
	debugPrint("Listening and serving HTTPS on host name is %s\n", domain)
	defer func() { debugPrintError(err) }()

	//your domain here
	if len(domain) != 0 {
		AutoTLSManager.HostPolicy = autocert.HostWhitelist(domain...)
	}

	s := &http.Server{
		Addr:      ":443",
		TLSConfig: &tls.Config{GetCertificate: AutoTLSManager.GetCertificate},
		Handler:   engine,
	}
	err = s.ListenAndServeTLS("", "")
	return
}
