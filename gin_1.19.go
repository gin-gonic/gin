// Copyright 2023 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build go1.19

package gin

import (
	"github.com/quic-go/quic-go/http3"
)

// RunQUIC attaches the router to a http.Server and starts listening and serving QUIC requests.
// It is a shortcut for http3.ListenAndServeQUIC(addr, certFile, keyFile, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunQUIC(addr, certFile, keyFile string) (err error) {
	debugPrint("Listening and serving QUIC on %s\n", addr)
	defer func() { debugPrintError(err) }()

	if engine.isUnsafeTrustedProxies() {
		debugPrint("[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
			"Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")
	}

	err = http3.ListenAndServeQUIC(addr, certFile, keyFile, engine.Handler())
	return
}
