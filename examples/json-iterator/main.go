package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/codec/json"
	jsoniter "github.com/json-iterator/go"
)

var customConfig = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

// customJsonApi implement api.JsonApi
type customJsonApi struct {
}

func (j customJsonApi) Marshal(v any) ([]byte, error) {
	return customConfig.Marshal(v)
}

func (j customJsonApi) Unmarshal(data []byte, v any) error {
	return customConfig.Unmarshal(data, v)
}

func (j customJsonApi) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	return customConfig.MarshalIndent(v, prefix, indent)
}

func (j customJsonApi) NewEncoder(writer io.Writer) json.Encoder {
	return customConfig.NewEncoder(writer)
}

func (j customJsonApi) NewDecoder(reader io.Reader) json.Decoder {
	return customConfig.NewDecoder(reader)
}

func main() {
	// Replace the default json api with json-iterator
	json.API = customJsonApi{}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run(":8080")
}
