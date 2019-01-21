// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build tools

// This package exists to cause `go mod` and `go get` to believe these tools
// are dependencies, even though they are not runtime dependencies of any
// gin package. This means they will appear in `go.mod` file, but will not
// be a part of the build.

package tools

import (
	_ "github.com/campoy/embedmd"
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/dustin/go-broadcast"
	_ "github.com/gin-gonic/autotls"
	_ "github.com/jessevdk/go-assets"
	_ "github.com/manucorporat/stats"
	_ "github.com/thinkerou/favicon"
	_ "golang.org/x/crypto/acme/autocert"
	_ "golang.org/x/lint/golint"
	_ "google.golang.org/grpc"
)
