// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"errors"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		{
			name: "Cookie header is masked",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("Cookie", "session=super-secret; theme=dark")
				return r
			}(),
			wantContains:   "Cookie: *",
			wantNotContain: "super-secret",
		},
		{
			name: "Proxy-Authorization header is masked",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("Proxy-Authorization", "Basic dXNlcjpwYXNz")
				return r
			}(),
			wantContains:   "Proxy-Authorization: *",
			wantNotContain: "dXNlcjpwYXNz",
		},
		{
			name: "multiple sensitive headers are masked",
			req: func() *http.Request {
				r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
				r.Header.Set("Authorization", "Bearer tok")
				r.Header.Set("Cookie", "sid=abc")
				r.Header.Set("Proxy-Authorization", "Basic xyz")
				r.Header.Set("X-Request-Id", "req-1")
				return r
			}(),
			wantContains:   "X-Request-Id: req-1",
			wantNotContain: "Bearer tok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := secureRequestDump(tt.req)
			if tt.wantContains != "" && !strings.Contains(result, tt.wantContains) {
				t.Errorf("secureRequestDump() = %q, want contains %q", result, tt.wantContains)
			}
			if tt.wantNotContain != "" && strings.Contains(result, tt.wantNotContain) {
				t.Errorf("secureRequestDump() = %q, want NOT contain %q", result, tt.wantNotContain)
			}
		})
	}

	// Combined check: secrets from every sensitive header must be gone.
	t.Run("no secret values remain when all sensitive headers set", func(t *testing.T) {
		r, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
		r.Header.Set("Authorization", "Bearer tok")
		r.Header.Set("Cookie", "sid=abc")
		r.Header.Set("Proxy-Authorization", "Basic xyz")
		result := secureRequestDump(r)
		for _, secret := range []string{"Bearer tok", "sid=abc", "Basic xyz"} {
			if strings.Contains(result, secret) {
				t.Errorf("secureRequestDump leaked %q in %q", secret, result)
			}
		}
		for _, masked := range []string{"Authorization: *", "Cookie: *", "Proxy-Authorization: *"} {
			if !strings.Contains(result, masked) {
				t.Errorf("secureRequestDump missing %q in %q", masked, result)
			}
		}
	})
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
