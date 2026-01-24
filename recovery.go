// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin/internal/bytesconv"
)

const (
	dunno     = "???"
	stackSkip = 3
)

// RecoveryFunc defines the function passable to CustomRecovery.
type RecoveryFunc func(c *Context, err any)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
func Recovery() HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter)
}

// CustomRecovery returns a middleware that recovers from any panics and calls the provided handle func to handle it.
func CustomRecovery(handle RecoveryFunc) HandlerFunc {
	return RecoveryWithWriter(DefaultErrorWriter, handle)
}

// RecoveryWithWriter returns a middleware for a given writer that recovers from any panics and writes a 500 if there was one.
func RecoveryWithWriter(out io.Writer, recovery ...RecoveryFunc) HandlerFunc {
	if len(recovery) > 0 {
		return CustomRecoveryWithWriter(out, recovery[0])
	}
	return CustomRecoveryWithWriter(out, defaultHandleRecovery)
}

// CustomRecoveryWithWriter returns a middleware for a given writer that recovers from any panics and calls the provided handle func to handle it.
func CustomRecoveryWithWriter(out io.Writer, handle RecoveryFunc) HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(c *Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var isBrokenPipe bool
				err, ok := rec.(error)
				if ok {
					isBrokenPipe = errors.Is(err, syscall.EPIPE) ||
						errors.Is(err, syscall.ECONNRESET) ||
						errors.Is(err, http.ErrAbortHandler)
				}
				if logger != nil {
					if isBrokenPipe {
						logger.Printf("%s\n%s%s", rec, secureRequestDump(c.Request), reset)
					} else if IsDebugging() {
						logger.Printf("[Recovery] %s panic recovered:\n%s\n%s\n%s%s",
							timeFormat(time.Now()), secureRequestDump(c.Request), rec, stack(stackSkip), reset)
					} else {
						logger.Printf("[Recovery] %s panic recovered:\n%s\n%s%s",
							timeFormat(time.Now()), rec, stack(stackSkip), reset)
					}
				}
				if isBrokenPipe {
					// If the connection is dead, we can't write a status to it.
					c.Error(err) //nolint: errcheck
					c.Abort()
				} else {
					handle(c, rec)
				}
			}
		}()
		c.Next()
	}
}

// secureRequestDump returns a sanitized HTTP request dump where the Authorization header,
// if present, is replaced with a masked value ("Authorization: *") to avoid leaking sensitive credentials.
//
// Currently, only the Authorization header is sanitized. All other headers and request data remain unchanged.
func secureRequestDump(r *http.Request) string {
	httpRequest, _ := httputil.DumpRequest(r, false)
	lines := strings.Split(bytesconv.BytesToString(httpRequest), "\r\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "Authorization:") {
			lines[i] = "Authorization: *"
		}
	}
	return strings.Join(lines, "\r\n")
}

func defaultHandleRecovery(c *Context, _ any) {
	c.AbortWithStatus(http.StatusInternalServerError)
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var (
		nLine    string
		lastFile string
		err      error
	)
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			nLine, err = readNthLine(file, line-1)
			if err != nil {
				continue
			}
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), cmp.Or(nLine, dunno))
	}
	return buf.Bytes()
}

// readNthLine reads the nth line from the file.
// It returns the trimmed content of the line if found,
// or an empty string if the line doesn't exist.
// If there's an error opening the file, it returns the error.
func readNthLine(file string, n int) (string, error) {
	if n < 0 {
		return "", nil
	}

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return "", nil
		}
	}

	if scanner.Scan() {
		return strings.TrimSpace(scanner.Text()), nil
	}

	return "", nil
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := fn.Name()
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := strings.LastIndexByte(name, '/'); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := strings.IndexByte(name, '.'); period >= 0 {
		name = name[period+1:]
	}
	name = strings.ReplaceAll(name, "·", ".")
	return name
}

// timeFormat returns a customized time string for logger.
func timeFormat(t time.Time) string {
	return t.Format("2006/01/02 - 15:04:05")
}
