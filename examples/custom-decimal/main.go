package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type QueryParams struct {
	Amount binding.CustomDecimal `form:"amount"`
}

func main() {
	r := gin.Default()

	r.GET("/amount", func(c *gin.Context) {
		var params QueryParams
		if err := c.BindQuery(&params); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"amount": params.Amount.String(),
		})
	})

	r.Run(":8080")
}
