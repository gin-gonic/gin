package ini

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

type Body struct {
	// json tag to de-serialize json body
	Name string `json:"name"`
	//For example, you can use struct tags to validate a custom product code format. The validator package offers many helpful string validator helpers.
	ProductCode string    `json:"productCode" binding:"required,startswith=PC,len=10"`
	StartDate   string    `json:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate     string    `json:"end_date" binding:"required" time_format:"2006-01-02"`
	EndDate1    time.Time `form:"end_date_1" binding:"required" time_format:"2006-01-02"`
}

type Person struct {
	ID string `uri:"id" binding:"required"`
	//ID   string `uri:"id" binding:"required,uuid"`
	Name string `uri:"name" binding:"required"`
}

func UrlInit(router *gin.Engine) {

	//普通url测试
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "You Can Try Another",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 测试AsciiJSON数据返回
	router.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}
		// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	// 正常的json数据返回
	router.GET("/someJSON2", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}
		// will output : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.JSON(http.StatusOK, data)
	})

	//Gin bindings are used to serialize JSON, XML, path parameters, form data, etc.
	//to structs and maps.
	//It also has a baked-in validation framework with complex validations.
	router.POST("/bingJson", func(c *gin.Context) {
		// one: de-serialize json body
		body := Body{}
		// using BindJson method to serialize body with struct
		// BindJSON reads the body buffer to de-serialize it to a struct.
		// BindJSON cannot be called on the same context twice because it flushes the body buffer.
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println(body)
		c.JSON(http.StatusAccepted, &body)
	})

	router.POST("/bingJson1", func(c *gin.Context) {
		// one: de-serialize json body
		body := Body{}
		// using BindJson method to serialize body with struct
		// BindJSON reads the body buffer to de-serialize it to a struct.
		// BindJSON cannot be called on the same context twice because it flushes the body buffer.
		if err := c.BindJSON(&body); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		body2 := Body{}
		if err := c.BindJSON(&body2); err != nil {
			//在Gin框架中，c.BindJSON()第二次调用会报错的原因是因为：
			//BindJSON()方法会读取并消耗HTTP请求的Body数据流。HTTP请求的Body是一个只能读取一次的io.ReadCloser接口实现。
			//当第一次调用c.BindJSON(&body)时，它会完整读取请求Body中的数据并解析到第一个结构体中，同时会将Body流关闭。
			//当第二次尝试调用c.BindJSON(&body2)时，Body流已经被关闭且数据已被消耗，所以会返回错误。
			//解决方案：
			//如果需要多次绑定同一个请求体，应该使用ShouldBindBodyWith()方法（如代码中/bingJson2路由所示），这个方法会将请求体内容缓存起来，允许后续多次绑定。
			//或者，可以在第一次绑定后将数据手动复制一份，供后续使用。
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println(body)
		c.JSON(http.StatusAccepted, &body)
		c.JSON(http.StatusAccepted, &body2)
	})

	router.POST("/bingJson2", func(c *gin.Context) {
		// one: de-serialize json body
		body := Body{}
		// using BindJson method to serialize body with struct
		// BindJSON reads the body buffer to de-serialize it to a struct.
		// BindJSON cannot be called on the same context twice because it flushes the body buffer.
		if err := c.ShouldBind(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				gin.H{
					"error":   "VALIDATEERR-1",
					"message": err.Error()})
			return
		}
		body2 := Body{}
		if err := c.ShouldBindBodyWith(&body2, binding.JSON); err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		fmt.Println(body)
		c.JSON(http.StatusAccepted, &body)
		c.JSON(http.StatusAccepted, &body2)
	})

	router.GET("/:name/:id", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
		c.JSON(200, gin.H{"name": person.Name, "uuid": person.ID})
	})

	router.GET("/cookie", func(c *gin.Context) {
		cookie, err := c.Cookie("gin_cookie")
		if err != nil {
			cookie = "NotSet"
			c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
			//Delete cookie by set max age to -1.
			//c.SetCookie("gin_cookie", "test", -1, "/", "localhost", false, true)
		}
		fmt.Printf("Cookie value: %s \n", cookie)
		c.JSON(http.StatusOK, gin.H{
			"cookie": cookie,
		})
	})
}
