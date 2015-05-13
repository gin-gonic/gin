package main

import (
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/manucorporat/stats"
)

var messages = stats.New()

func main() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(ratelimit, gin.Recovery(), gin.Logger())

	router.LoadHTMLGlob("resources/*.templ.html")
	router.Static("/static", "resources/static")
	router.GET("/", index)
	router.GET("/room/:roomid", roomGET)
	router.POST("/room/:roomid", roomPOST)
	//router.DELETE("/room/:roomid", roomDELETE)
	router.GET("/stream/:roomid", streamRoom)

	router.Run("127.0.0.1:8080")
}

func index(c *gin.Context) {
	c.Redirect(301, "/room/hn")
}

func roomGET(c *gin.Context) {
	roomid := c.ParamValue("roomid")
	userid := c.FormValue("nick")
	if len(userid) > 13 {
		userid = userid[0:12] + "..."
	}
	c.HTML(200, "room_login.templ.html", gin.H{
		"roomid":    roomid,
		"nick":      userid,
		"timestamp": time.Now().Unix(),
	})

}

func roomPOST(c *gin.Context) {
	roomid := c.ParamValue("roomid")
	nick := c.FormValue("nick")
	message := c.PostFormValue("message")

	if len(message) > 200 || len(nick) > 13 {
		c.JSON(400, gin.H{
			"status": "failed",
			"error":  "the message or nickname is too long",
		})
		return
	}

	post := gin.H{
		"nick":    nick,
		"message": message,
	}
	messages.Add("inbound", 1)
	room(roomid).Submit(post)
	c.JSON(200, post)
}

func roomDELETE(c *gin.Context) {
	roomid := c.ParamValue("roomid")
	deleteBroadcast(roomid)
}

func streamRoom(c *gin.Context) {
	roomid := c.ParamValue("roomid")
	listener := openListener(roomid)
	ticker := time.NewTicker(1 * time.Second)
	defer closeListener(roomid, listener)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-listener:
			messages.Add("outbound", 1)
			c.SSEvent("message", msg)
		case <-ticker.C:
			c.SSEvent("stats", Stats())
		}
		return true
	})
}
