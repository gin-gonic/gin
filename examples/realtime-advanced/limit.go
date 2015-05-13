package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/manucorporat/stats"
)

var ips = stats.New()

func ratelimit(c *gin.Context) {
	ip := c.ClientIP()
	value := uint64(ips.Add(ip, 1))
	if value >= 1000 {
		if value%1000 == 0 {
			log.Printf("BlockedIP:%s Requests:%d\n", ip, value)
		}
		c.AbortWithStatus(401)
	}
}
