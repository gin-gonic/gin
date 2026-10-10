package gin_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// readdirOnlyFS exposes exactly http.FileSystem's documented interfaces.
type readdirOnlyFS struct{ http.FileSystem }

func (f readdirOnlyFS) Open(name string) (http.File, error) {
	file, err := f.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}
	return struct{ http.File }{file}, nil
}

func TestLoadHTMLFSHTTPFileSystemGlob(t *testing.T) {
	old := gin.Mode()
	t.Cleanup(func() { gin.SetMode(old) })
	for _, mode := range []string{gin.DebugMode, gin.ReleaseMode} {
		for _, backend := range []string{"http-only", "http.FS", "http.Dir"} {
			for _, pattern := range []string{"templates/page.html", "templates/*.html", "*/page.html"} {
				t.Run(mode+"/"+backend+"/"+pattern, func(t *testing.T) {
					gin.SetMode(mode)
					engine := gin.New()
					source := fstest.MapFS{"templates/page.html": &fstest.MapFile{Data: []byte("Hello {{.name}}")}}
					var filesystem http.FileSystem = readdirOnlyFS{http.FS(source)}
					switch backend {
					case "http.FS":
						filesystem = http.FS(source)
					case "http.Dir":
						root := t.TempDir()
						require.NoError(t, os.Mkdir(filepath.Join(root, "templates"), 0o700))
						require.NoError(t, os.WriteFile(filepath.Join(root, "templates", "page.html"), source["templates/page.html"].Data, 0o600))
						filesystem = http.Dir(root)
					}
					require.NotPanics(t, func() { engine.LoadHTMLFS(filesystem, pattern) })
					engine.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "page.html", gin.H{"name": "Gin"}) })
					response := httptest.NewRecorder()
					require.NotPanics(t, func() { engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil)) })
					assert.Equal(t, http.StatusOK, response.Code)
					assert.Equal(t, "Hello Gin", response.Body.String())
				})
			}
		}
	}
}
