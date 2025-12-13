package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Example middleware that reads the request body
	r.Use(func(c *gin.Context) {
		// Get the request body - this can be called multiple times
		body, err := c.GetRequestBody()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// Log the body length for demonstration
		fmt.Printf("Middleware: Request body length: %d bytes\n", len(body))

		// Store body in context for use by handlers
		c.Set("rawBody", body)
		c.Next()
	})

	// Handler that also reads the body
	r.POST("/echo", func(c *gin.Context) {
		// Get the body again - this will use the cached version
		body, err := c.GetRequestBody()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Also get the body that was stored by middleware
		storedBody, exists := c.Get("rawBody")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "body not found in context"})
			return
		}

		// Both should be identical
		if string(body) != string(storedBody.([]byte)) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "bodies don't match"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Body received and cached successfully",
			"body":    string(body),
			"length":  len(body),
		})
	})

	// Handler that uses binding after GetRequestBody
	r.POST("/bind-after-body", func(c *gin.Context) {
		// First bind to a struct - this caches the body
		var jsonData map[string]interface{}
		if err := c.ShouldBindBodyWithJSON(&jsonData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Then get the raw body - this will use the cached version
		rawBody, err := c.GetRequestBody()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"raw_body": string(rawBody),
			"parsed_data": jsonData,
		})
	})

	fmt.Println("Server starting on :9090")
	fmt.Println("Try: curl -X POST -d 'hello world' http://localhost:9090/echo")
	fmt.Println("Or: curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"test\"}' http://localhost:9090/bind-after-body")

	r.Run(":9090")
}
