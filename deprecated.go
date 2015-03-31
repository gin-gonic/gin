// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin/binding"
)

// DEPRECATED, use Bind() instead.
// Like ParseBody() but this method also writes a 400 error if the json is not valid.
func (c *Context) EnsureBody(item interface{}) bool {
	return c.Bind(item)
}

// DEPRECATED use bindings directly
// Parses the body content as a JSON input. It decodes the json payload into the struct specified as a pointer.
func (c *Context) ParseBody(item interface{}) error {
	return binding.JSON.Bind(c.Request, item)
}

// DEPRECATED use gin.Static() instead
// ServeFiles serves files from the given file system root.
// The path must end with "/*filepath", files are then served from the local
// path /defined/root/dir/*filepath.
// For example if root is "/etc" and *filepath is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use http.Dir:
//     router.ServeFiles("/src/*filepath", http.Dir("/var/www"))
func (engine *Engine) ServeFiles(path string, root http.FileSystem) {
	engine.router.ServeFiles(path, root)
}

// DEPRECATED use gin.LoadHTMLGlob() or gin.LoadHTMLFiles() instead
func (engine *Engine) LoadHTMLTemplates(pattern string) {
	engine.LoadHTMLGlob(pattern)
}

// DEPRECATED. Use NoRoute() instead
func (engine *Engine) NotFound404(handlers ...HandlerFunc) {
	engine.NoRoute(handlers...)
}

// the ForwardedFor middleware unwraps the X-Forwarded-For headers, be careful to only use this
// middleware if you've got servers in front of this server. The list with (known) proxies and
// local ips are being filtered out of the forwarded for list, giving the last not local ip being
// the real client ip.
func ForwardedFor(proxies ...interface{}) HandlerFunc {
	if len(proxies) == 0 {
		// default to local ips
		var reservedLocalIps = []string{"10.0.0.0/8", "127.0.0.1/32", "172.16.0.0/12", "192.168.0.0/16"}

		proxies = make([]interface{}, len(reservedLocalIps))

		for i, v := range reservedLocalIps {
			proxies[i] = v
		}
	}

	return func(c *Context) {
		// the X-Forwarded-For header contains an array with left most the client ip, then
		// comma separated, all proxies the request passed. The last proxy appears
		// as the remote address of the request. Returning the client
		// ip to comply with default RemoteAddr response.

		// check if remoteaddr is local ip or in list of defined proxies
		remoteIp := net.ParseIP(strings.Split(c.Request.RemoteAddr, ":")[0])

		if !ipInMasks(remoteIp, proxies) {
			return
		}

		if forwardedFor := c.Request.Header.Get("X-Forwarded-For"); forwardedFor != "" {
			parts := strings.Split(forwardedFor, ",")

			for i := len(parts) - 1; i >= 0; i-- {
				part := parts[i]

				ip := net.ParseIP(strings.TrimSpace(part))

				if ipInMasks(ip, proxies) {
					continue
				}

				// returning remote addr conform the original remote addr format
				c.Request.RemoteAddr = ip.String() + ":0"

				// remove forwarded for address
				c.Request.Header.Set("X-Forwarded-For", "")
				return
			}
		}
	}
}

func ipInMasks(ip net.IP, masks []interface{}) bool {
	for _, proxy := range masks {
		var mask *net.IPNet
		var err error

		switch t := proxy.(type) {
		case string:
			if _, mask, err = net.ParseCIDR(t); err != nil {
				log.Panic(err)
			}
		case net.IP:
			mask = &net.IPNet{IP: t, Mask: net.CIDRMask(len(t)*8, len(t)*8)}
		case net.IPNet:
			mask = &t
		}

		if mask.Contains(ip) {
			return true
		}
	}

	return false
}
