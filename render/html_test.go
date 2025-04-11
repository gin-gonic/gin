// render/html_test.go
package render

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHTMLProductionInstanceSuccess tests a successful template instance in production
func TestHTMLProductionInstanceSuccess(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse("Hello {{.}}"))
	renderer := HTMLProduction{Template: tmpl}

	instance := renderer.Instance("test", "World")
	html, ok := instance.(HTML)
	assert.True(t, ok, "Instance should be of type HTML")
	assert.Equal(t, "test", html.Name)
	assert.Equal(t, "World", html.Data)
}

// TestHTMLProductionInstanceNilTemplate tests when template is not initialized
func TestHTMLProductionInstanceNilTemplate(t *testing.T) {
	renderer := HTMLProduction{Template: nil}

	instance := renderer.Instance("test", "World")
	html, ok := instance.(HTML)
	assert.True(t, ok, "Instance should be of type HTML")
	assert.Nil(t, html.Template, "Template should be nil when not initialized")
}

// TestHTMLDebugInstanceSuccess tests a successful template instance in debug mode
func TestHTMLDebugInstanceSuccess(t *testing.T) {
	renderer := HTMLDebug{
		Files:   nil, // No files provided
		Delims:  Delims{Left: "{{", Right: "}}"},
		FuncMap: template.FuncMap{},
	}

	// Since no files or glob are provided, Template will be nil
	instance := renderer.Instance("inline", "World")
	html, ok := instance.(HTML)
	assert.True(t, ok, "Instance should be of type HTML")
	assert.Equal(t, "inline", html.Name)
	assert.Equal(t, "World", html.Data)
	assert.Nil(t, html.Template, "Template should be nil since no files or glob provided")
}

// TestHTMLDebugInstanceNoFilesOrGlob tests behavior when no files or glob are provided
func TestHTMLDebugInstanceNoFilesOrGlob(t *testing.T) {
	renderer := HTMLDebug{
		Files:   nil,
		Glob:    "",
		Delims:  Delims{Left: "{{", Right: "}}"},
		FuncMap: template.FuncMap{},
	}

	instance := renderer.Instance("test", "World")
	html, ok := instance.(HTML)
	assert.True(t, ok, "Instance should be of type HTML")
	assert.Equal(t, "test", html.Name)
	assert.Equal(t, "World", html.Data)
	assert.Nil(t, html.Template, "Template should be nil when no files or glob provided")
}
