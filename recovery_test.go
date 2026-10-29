// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecoveryDisableConsoleColor(t *testing.T) {
	previousMode := consoleColorMode
	previousWriter := DefaultErrorWriter
	t.Cleanup(func() {
		consoleColorMode = previousMode
		DefaultErrorWriter = previousWriter
	})

	DisableConsoleColor()
	buffer := new(strings.Builder)
	DefaultErrorWriter = buffer
	router := New()
	router.Use(Recovery())
	router.GET("/recovery", func(_ *Context) {
		panic("recovery color regression")
	})

	w := PerformRequest(router, http.MethodGet, "/recovery")
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "recovery color regression")
	assert.NotContains(t, buffer.String(), "\x1b[")
}

func TestRecoveryColorModes(t *testing.T) {
	previousMode := consoleColorMode
	previousWriter := DefaultErrorWriter
	previousGinMode := Mode()
	t.Cleanup(func() {
		consoleColorMode = previousMode
		DefaultErrorWriter = previousWriter
		SetMode(previousGinMode)
	})
	SetMode(DebugMode)

	handle := func(c *Context, _ any) {
		c.AbortWithStatus(http.StatusTeapot)
	}
	constructors := []struct {
		name   string
		new    func(io.Writer) HandlerFunc
		status int
	}{
		{"Recovery", func(out io.Writer) HandlerFunc {
			DefaultErrorWriter = out
			return Recovery()
		}, http.StatusInternalServerError},
		{"RecoveryWithWriter", func(out io.Writer) HandlerFunc {
			return RecoveryWithWriter(out)
		}, http.StatusInternalServerError},
		{"RecoveryWithWriterCustom", func(out io.Writer) HandlerFunc {
			return RecoveryWithWriter(out, handle)
		}, http.StatusTeapot},
		{"CustomRecovery", func(out io.Writer) HandlerFunc {
			DefaultErrorWriter = out
			return CustomRecovery(handle)
		}, http.StatusTeapot},
		{"CustomRecoveryWithWriter", func(out io.Writer) HandlerFunc {
			return CustomRecoveryWithWriter(out, handle)
		}, http.StatusTeapot},
	}
	modes := []struct {
		name      string
		set       func()
		term      string
		wantColor bool
	}{
		{"disabled", DisableConsoleColor, "xterm", false},
		{"forced", ForceConsoleColor, "xterm", true},
		{"automatic buffer", func() { consoleColorMode = autoColor }, "xterm", false},
		{"automatic dumb", func() { consoleColorMode = autoColor }, "dumb", false},
		{"forced dumb", ForceConsoleColor, "dumb", true},
	}
	panics := []struct {
		name  string
		value any
	}{
		{"normal", "color panic"},
		{"broken pipe", &net.OpError{Err: &os.SyscallError{Err: syscall.EPIPE}}},
	}

	for _, constructor := range constructors {
		for _, mode := range modes {
			for _, panicCase := range panics {
				t.Run(constructor.name+"/"+mode.name+"/"+panicCase.name, func(t *testing.T) {
					t.Setenv("TERM", mode.term)
					mode.set()
					buffer := new(strings.Builder)
					password := "Bearer secret-token"
					router := New()
					router.Use(constructor.new(buffer))
					router.GET("/recovery", func(_ *Context) {
						panic(panicCase.value)
					})

					w := PerformRequest(router, http.MethodGet, "/recovery", header{Key: "Authorization", Value: password})
					out := buffer.String()
					if panicCase.name == "broken pipe" {
						assert.Equal(t, http.StatusOK, w.Code)
						assert.Contains(t, strings.ToLower(out), "broken pipe")
						assert.NotContains(t, out, "panic recovered")
					} else {
						assert.Equal(t, constructor.status, w.Code)
						assert.Contains(t, out, "panic recovered")
						assert.Contains(t, out, "color panic")
					}
					assert.Contains(t, out, "Authorization: *")
					assert.NotContains(t, out, password)
					if mode.wantColor {
						assert.Contains(t, out, "\x1b[31m")
						assert.Contains(t, out, reset)
					} else {
						assert.NotContains(t, out, "\x1b[")
					}
				})
			}
		}
	}
}

func TestRecoveryColorModeChanges(t *testing.T) {
	previousMode := consoleColorMode
	t.Cleanup(func() { consoleColorMode = previousMode })

	buffer := new(strings.Builder)
	DisableConsoleColor()
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic("changing color mode")
	})

	modes := []struct {
		name      string
		set       func()
		wantColor bool
	}{
		{"forced after disabled", ForceConsoleColor, true},
		{"disabled after forced", DisableConsoleColor, false},
		{"automatic after disabled", func() { consoleColorMode = autoColor }, false},
		{"forced after automatic", ForceConsoleColor, true},
	}
	for _, mode := range modes {
		t.Run(mode.name, func(t *testing.T) {
			mode.set()
			buffer.Reset()
			w := PerformRequest(router, http.MethodGet, "/recovery")
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Contains(t, buffer.String(), "changing color mode")
			if mode.wantColor {
				assert.Contains(t, buffer.String(), "\x1b[31m")
				assert.Contains(t, buffer.String(), reset)
			} else {
				assert.NotContains(t, buffer.String(), "\x1b[")
			}
		})
	}
}

func TestRecoveryNilWriterColorModes(t *testing.T) {
	previousMode := consoleColorMode
	t.Cleanup(func() { consoleColorMode = previousMode })

	for _, mode := range []struct {
		name string
		set  func()
	}{
		{"disabled", DisableConsoleColor},
		{"forced", ForceConsoleColor},
		{"automatic", func() { consoleColorMode = autoColor }},
	} {
		t.Run(mode.name, func(t *testing.T) {
			mode.set()
			for _, custom := range []bool{false, true} {
				called := false
				handle := func(c *Context, err any) {
					called = true
					assert.Equal(t, "nil writer panic", err)
					c.AbortWithStatus(http.StatusTeapot)
				}
				router := New()
				if custom {
					router.Use(CustomRecoveryWithWriter(nil, handle))
				} else {
					router.Use(RecoveryWithWriter(nil))
				}
				router.GET("/recovery", func(_ *Context) { panic("nil writer panic") })
				w := PerformRequest(router, http.MethodGet, "/recovery")
				if custom {
					assert.True(t, called)
					assert.Equal(t, http.StatusTeapot, w.Code)
				} else {
					assert.Equal(t, http.StatusInternalServerError, w.Code)
				}
			}
		})
	}
}

func TestRecoveryConcurrentColorModes(t *testing.T) {
	previousMode := consoleColorMode
	t.Cleanup(func() { consoleColorMode = previousMode })

	for _, mode := range []struct {
		name      string
		set       func()
		wantColor bool
	}{
		{"disabled", DisableConsoleColor, false},
		{"forced", ForceConsoleColor, true},
		{"automatic", func() { consoleColorMode = autoColor }, false},
	} {
		t.Run(mode.name, func(t *testing.T) {
			mode.set()
			var buffer bytes.Buffer
			router := New()
			router.Use(RecoveryWithWriter(&buffer))
			router.GET("/recovery", func(_ *Context) { panic("concurrent panic") })

			const requests = 20
			statuses := make([]int, requests)
			var wg sync.WaitGroup
			for i := range requests {
				wg.Add(1)
				go func() {
					defer wg.Done()
					statuses[i] = PerformRequest(router, http.MethodGet, "/recovery").Code
				}()
			}
			wg.Wait()

			for _, status := range statuses {
				assert.Equal(t, http.StatusInternalServerError, status)
			}
			out := buffer.String()
			assert.Equal(t, requests, strings.Count(out, "panic recovered"))
			if mode.wantColor {
				assert.Equal(t, requests, strings.Count(out, "\x1b[31m"))
				assert.Equal(t, requests, strings.Count(out, reset))
			} else {
				assert.NotContains(t, out, "\x1b[")
			}
			entries := strings.Split(out, "\n\n")
			assert.Len(t, entries, requests+1)
			for _, entry := range entries[1:] {
				assert.Contains(t, entry, "panic recovered")
				if mode.wantColor {
					assert.True(t, strings.HasPrefix(strings.TrimLeft(entry, "\n"), "\x1b[31m"))
					assert.True(t, strings.HasSuffix(strings.TrimSpace(entry), reset))
				}
			}
		})
	}
}

func recoveryStackFrameHandler(_ *Context) {
	panic("stack frame panic")
}

func TestRecoveryStackFrames(t *testing.T) {
	previousMode := Mode()
	t.Cleanup(func() { SetMode(previousMode) })
	for _, mode := range []string{TestMode, DebugMode} {
		t.Run(mode, func(t *testing.T) {
			SetMode(mode)
			var buffer strings.Builder
			router := New()
			router.Use(RecoveryWithWriter(&buffer))
			router.GET("/recovery", recoveryStackFrameHandler)

			w := PerformRequest(router, http.MethodGet, "/recovery")
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Regexp(t, `stack frame panic\n[^\n]*recovery_test\.go:[0-9]+ \(0x[0-9a-f]+\)\n\trecoveryStackFrameHandler:`, buffer.String())
		})
	}
}

func TestPanicClean(t *testing.T) {
	buffer := new(strings.Builder)
	router := New()
	password := "my-super-secret-password"
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(c *Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic("Oops, Houston, we have a problem")
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery",
		header{
			Key:   "Host",
			Value: "www.google.com",
		},
		header{
			Key:   "Authorization",
			Value: "Bearer " + password,
		},
		header{
			Key:   "Content-Type",
			Value: "application/json",
		},
	)
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check the buffer does not have the secret key
	assert.NotContains(t, buffer.String(), password)
}

// TestPanicInHandler assert that panic has been recovered.
func TestPanicInHandler(t *testing.T) {
	buffer := new(strings.Builder)
	router := New()
	router.Use(RecoveryWithWriter(buffer))
	router.GET("/recovery", func(_ *Context) {
		panic("Oops, Houston, we have a problem")
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oops, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	SetMode(TestMode)
}

// TestPanicWithAbort assert that panic has been recovered even if context.Abort was used.
func TestPanicWithAbort(t *testing.T) {
	router := New()
	router.Use(RecoveryWithWriter(nil))
	router.GET("/recovery", func(c *Context) {
		c.AbortWithStatus(http.StatusBadRequest)
		panic("Oops, Houston, we have a problem")
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFunction(t *testing.T) {
	bs := function(1)
	assert.Equal(t, dunno, bs)
}

// TestPanicWithBrokenPipe asserts that recovery specifically handles
// writing responses to broken pipes
func TestPanicWithBrokenPipe(t *testing.T) {
	const expectCode = 204

	expectErrnos := []syscall.Errno{
		syscall.EPIPE,
		syscall.ECONNRESET,
	}

	for _, errno := range expectErrnos {
		t.Run("Recovery from "+errno.Error(), func(t *testing.T) {
			var buf strings.Builder

			router := New()
			router.Use(RecoveryWithWriter(&buf))
			router.GET("/recovery", func(c *Context) {
				// Start writing response
				c.Header("X-Test", "Value")
				c.Status(expectCode)

				// Oops. Client connection closed
				e := &net.OpError{Err: &os.SyscallError{Err: errno}}
				panic(e)
			})
			// RUN
			w := PerformRequest(router, http.MethodGet, "/recovery")
			// TEST
			assert.Equal(t, expectCode, w.Code)
			assert.Contains(t, strings.ToLower(buf.String()), errno.Error())
			assert.NotContains(t, strings.ToLower(buf.String()), "[Recovery]")
		})
	}
}

// TestPanicWithAbortHandler asserts that recovery handles http.ErrAbortHandler as broken pipe
func TestPanicWithAbortHandler(t *testing.T) {
	const expectCode = 204

	var buf strings.Builder
	router := New()
	router.Use(RecoveryWithWriter(&buf))
	router.GET("/recovery", func(c *Context) {
		// Start writing response
		c.Header("X-Test", "Value")
		c.Status(expectCode)

		// Panic with ErrAbortHandler which should be treated as broken pipe
		panic(http.ErrAbortHandler)
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, expectCode, w.Code)
	out := buf.String()
	assert.Contains(t, out, "net/http: abort Handler")
	assert.NotContains(t, out, "panic recovered")
}

func TestPanicInHandlerRecordsError(t *testing.T) {
	tests := []struct {
		name         string
		recoveredErr any
		expectedErr  string
	}{
		{
			name:         "string panic",
			recoveredErr: "Oops, Houston, we have a problem",
			expectedErr:  "Oops, Houston, we have a problem",
		},
		{
			name:         "error panic",
			recoveredErr: errors.New("recovered error"),
			expectedErr:  "recovered error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := New()

			var recoveredErrors errorMsgs
			router.Use(func(c *Context) {
				c.Next()
				recoveredErrors = c.Errors
			})
			router.Use(RecoveryWithWriter(nil))
			router.GET("/recovery", func(_ *Context) {
				panic(tt.recoveredErr)
			})

			w := PerformRequest(router, http.MethodGet, "/recovery")

			assert.Equal(t, http.StatusInternalServerError, w.Code)
			if assert.Len(t, recoveredErrors, 1) {
				assert.EqualError(t, recoveredErrors[0], tt.expectedErr)
				assert.Equal(t, ErrorTypePrivate, recoveredErrors[0].Type)
			}
		})
	}
}

func TestCustomRecoveryWithWriter(t *testing.T) {
	errBuffer := new(strings.Builder)
	buffer := new(strings.Builder)
	router := New()
	handleRecovery := func(c *Context, err any) {
		errBuffer.WriteString(err.(string))
		c.AbortWithStatus(http.StatusBadRequest)
	}
	router.Use(CustomRecoveryWithWriter(buffer, handleRecovery))
	router.GET("/recovery", func(_ *Context) {
		panic("Oops, Houston, we have a problem")
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oops, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	assert.Equal(t, strings.Repeat("Oops, Houston, we have a problem", 2), errBuffer.String())

	SetMode(TestMode)
}

func TestCustomRecovery(t *testing.T) {
	errBuffer := new(strings.Builder)
	buffer := new(strings.Builder)
	router := New()
	DefaultErrorWriter = buffer
	handleRecovery := func(c *Context, err any) {
		errBuffer.WriteString(err.(string))
		c.AbortWithStatus(http.StatusBadRequest)
	}
	router.Use(CustomRecovery(handleRecovery))
	router.GET("/recovery", func(_ *Context) {
		panic("Oops, Houston, we have a problem")
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oops, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	assert.Equal(t, strings.Repeat("Oops, Houston, we have a problem", 2), errBuffer.String())

	SetMode(TestMode)
}

func TestRecoveryWithWriterWithCustomRecovery(t *testing.T) {
	errBuffer := new(strings.Builder)
	buffer := new(strings.Builder)
	router := New()
	DefaultErrorWriter = buffer
	handleRecovery := func(c *Context, err any) {
		errBuffer.WriteString(err.(string))
		c.AbortWithStatus(http.StatusBadRequest)
	}
	router.Use(RecoveryWithWriter(DefaultErrorWriter, handleRecovery))
	router.GET("/recovery", func(_ *Context) {
		panic("Oops, Houston, we have a problem")
	})
	// RUN
	w := PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "panic recovered")
	assert.Contains(t, buffer.String(), "Oops, Houston, we have a problem")
	assert.Contains(t, buffer.String(), t.Name())
	assert.NotContains(t, buffer.String(), "GET /recovery")

	// Debug mode prints the request
	SetMode(DebugMode)
	// RUN
	w = PerformRequest(router, http.MethodGet, "/recovery")
	// TEST
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, buffer.String(), "GET /recovery")

	assert.Equal(t, strings.Repeat("Oops, Houston, we have a problem", 2), errBuffer.String())

	SetMode(TestMode)
}

func TestSecureRequestDump(t *testing.T) {
	tests := []struct {
		name           string
		req            *http.Request
		wantContains   string
		wantNotContain string
	}{
		{
			name: "Authorization header standard case",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("Authorization", "Bearer secret-token")
				return r
			}(),
			wantContains:   "Authorization: *",
			wantNotContain: "Bearer secret-token",
		},
		{
			name: "authorization header lowercase",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("authorization", "some-secret")
				return r
			}(),
			wantContains:   "Authorization: *",
			wantNotContain: "some-secret",
		},
		{
			name: "Authorization header mixed case",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("AuThOrIzAtIoN", "token123")
				return r
			}(),
			wantContains:   "Authorization: *",
			wantNotContain: "token123",
		},
		{
			name: "No Authorization header",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("Content-Type", "application/json")
				return r
			}(),
			wantContains:   "",
			wantNotContain: "Authorization: *",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := secureRequestDump(tt.req)
			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("maskHeaders() = %q, want contains %q", result, tt.wantContains)
			}
			if tt.wantNotContain != "" && strings.Contains(result, tt.wantNotContain) {
				t.Errorf("maskHeaders() = %q, want NOT contain %q", result, tt.wantNotContain)
			}
		})
	}
}

// TestReadNthLine tests the readNthLine function with various scenarios.
func TestReadNthLine(t *testing.T) {
	// Create a temporary test file
	testContent := "line 0 \n line 1  \nline 2 \nline 3  \nline 4"
	tempFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write test content to the temporary file
	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatal(err)
	}
	if err := tempFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test cases
	tests := []struct {
		name     string
		lineNum  int
		fileName string
		want     string
		wantErr  bool
	}{
		{name: "Read first line", lineNum: 0, fileName: tempFile.Name(), want: "line 0", wantErr: false},
		{name: "Read middle line", lineNum: 2, fileName: tempFile.Name(), want: "line 2", wantErr: false},
		{name: "Read last line", lineNum: 4, fileName: tempFile.Name(), want: "line 4", wantErr: false},
		{name: "Line number exceeds file length", lineNum: 10, fileName: tempFile.Name(), want: "", wantErr: false},
		{name: "Negative line number", lineNum: -1, fileName: tempFile.Name(), want: "", wantErr: false},
		{name: "Non-existent file", lineNum: 1, fileName: "/non/existent/file.txt", want: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readNthLine(tt.fileName, tt.lineNum)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkStack(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = stack(stackSkip)
	}
}
