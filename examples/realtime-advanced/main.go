package main

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("resources/*.templ.html")
	router.Static("/static", "resources/static")
	router.GET("/", index)
	router.GET("/room/:roomid", roomGET)
	router.POST("/room/:roomid", roomPOST)
	//router.DELETE("/room/:roomid", roomDELETE)
	router.GET("/stream/:roomid", streamRoom)

	router.Run(":8080")
}

func index(c *gin.Context) {
	c.Redirect(301, "/room/hn")
}

func roomGET(c *gin.Context) {
	roomid := c.ParamValue("roomid")
	userid := c.FormValue("nick")
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

	if len(message) > 200 || len(nick) > 20 {
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
			c.SSEvent("message", msg)
		case <-ticker.C:
			c.SSEvent("stats", Stats())
		}
		return true
	})
}
