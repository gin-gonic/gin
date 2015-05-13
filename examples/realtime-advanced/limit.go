package main

import "github.com/gin-gonic/gin"

import "github.com/manucorporat/stats"

var ips = stats.New()

func ratelimit(c *gin.Context) {
	ip := c.ClientIP()
	value := ips.Add(ip, 1)
	if value > 400 {
		c.AbortWithStatus(401)
	}
}
