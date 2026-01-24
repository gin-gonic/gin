package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	ginjson "github.com/gin-gonic/gin/codec/json"
	"github.com/stretchr/testify/assert"
)

func TestJsonIterator(t *testing.T) {
	// Restore default json api after test
	originalAPI := ginjson.API
	defer func() {
		ginjson.API = originalAPI
	}()

	// Use custom json api
	ginjson.API = customJsonApi{}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"hello": "world",
			"foo":   "bar",
		})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify JSON response
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "world", response["hello"])
	assert.Equal(t, "bar", response["foo"])
}
