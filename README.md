#Gin Web Framework

[![GoDoc](https://godoc.org/github.com/gin-gonic/gin?status.svg)](https://godoc.org/github.com/gin-gonic/gin)
[![Build Status](https://travis-ci.org/gin-gonic/gin.svg)](https://travis-ci.org/gin-gonic/gin)

Gin is a web framework written in Golang. It features a martini-like API with much better performance, up to 40 times faster. If you need performance and good productivity, you will love Gin.  
![Gin console logger](http://gin-gonic.github.io/gin/other/console.png)

##Gin is new, will it be supported?

Yes, Gin is an internal project of [my](https://github.com/manucorporat) upcoming startup. We developed it and we are going to continue using and improve it.


##Roadmap for v1.0
- [x] Performance improments, reduce allocation and garbage collection overhead
- [x] Fix bugs
- [ ] Stable API
- [ ] Ask our designer for a cool logo
- [ ] Add tons of unit tests
- [ ] Add internal benchmarks suite
- [x] Improve logging system
- [x] Improve JSON/XML validation using bindings
- [x] Improve XML support
- [x] Flexible rendering system
- [ ] More powerful validation API
- [ ] Improve documentation
- [ ] Add more cool middlewares, for example redis caching (this also helps developers to understand the framework).
- [x] Continuous integration



## Start using it
Obviously, you need to have Git and Go! already installed to run Gin.  
Run this in your terminal

```
go get github.com/gin-gonic/gin
```
Then import it in your Go! code:

```
import "github.com/gin-gonic/gin"
```


##Community
If you'd like to help out with the project, there's a mailing list and IRC channel where Gin discussions normally happen.

* IRC
 * [irc.freenode.net #getgin](irc://irc.freenode.net:6667/getgin)
 * [Webchat](http://webchat.freenode.net?randomnick=1&channels=%23getgin)
* Mailing List
 * Subscribe: [getgin@librelist.org](mailto:getgin@librelist.org)
 * [Archives](http://librelist.com/browser/getgin/)


##API Examples

#### Create most basic PING/PONG HTTP endpoint
```go 
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Using GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
func main() {
	// Creates a gin router + logger and recovery (crash-free) middlewares
	r := gin.Default()

	r.GET("/someGet", getting)
	r.POST("/somePost", posting)
	r.PUT("/somePut", putting)
	r.DELETE("/someDelete", deleting)
	r.PATCH("/somePatch", patching)
	r.HEAD("/someHead", head)
	r.OPTIONS("/someOptions", options)

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Parameters in path

```go
func main() {
	r := gin.Default()
	
	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Params.ByName("name")
		message := "Hello "+name
		c.String(200, message)
	})

	// However, this one will match /user/john and also /user/john/send
	r.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Params.ByName("name")
		action := c.Params.ByName("action")
		message := name + " is " + action
		c.String(200, message)
	})
	
	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```


#### Grouping routes
```go
func main() {
	r := gin.Default()

	// Simple group: v1
	v1 := r.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := r.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```


#### Blank Gin without middlewares by default

Use

```go
r := gin.New()
```
instead of

```go
r := gin.Default()
```


#### Using middlewares
```go
func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Per route middlewares, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same than:
	authorized := r.Group("/")
	// per group middlewares! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Model binding and validation

To bind a request body into a type, use model binding. We currently support binding of JSON, XML and standard form values (foo=bar&boo=baz).

Note that you need to set the corresponding binding tag on all fields you want to bind. For example, when binding from JSON, set `json:"fieldname"`.

When using the Bind-method, Gin tries to infer the binder depending on the Content-Type header. If you are sure what you are binding, you can use BindWith. 

You can also specify that specific fields are required. If a field is decorated with `binding:"required"` and has a empty value when binding, the current request will fail with an error.

```go
// Binding from JSON
type LoginJSON struct {
	User     string `json:"user" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Binding from form values
type LoginForm struct {
    User     string `form:"user" binding:"required"`
    Password string `form:"password" binding:"required"`   
}

func main() {
	r := gin.Default()

    // Example for binding JSON ({"user": "manu", "password": "123"})
	r.POST("/login", func(c *gin.Context) {
		var json LoginJSON

        c.Bind(&json) // This will infer what binder to use depending on the content-type header.
        if json.User == "manu" && json.Password == "123" {
            c.JSON(200, gin.H{"status": "you are logged in"})
        } else {
            c.JSON(401, gin.H{"status": "unauthorized"})
        }
	})

    // Example for binding a HTLM form (user=manu&password=123)
    r.POST("/login", func(c *gin.Context) {
        var form LoginForm

        c.BindWith(&form, binding.Form) // You can also specify which binder to use. We support binding.Form, binding.JSON and binding.XML.
        if form.User == "manu" && form.Password == "123" {
            c.JSON(200, gin.H{"status": "you are logged in"})
        } else {
            c.JSON(401, gin.H{"status": "unauthorized"})
        }
    })

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### XML and JSON rendering

```go
func main() {
	r := gin.Default()

	// gin.H is a shortcup for map[string]interface{}
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hey", "status": 200})
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		// You also can use a struct
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// Note that msg.Name becomes "user" in the JSON
		// Will output  :   {"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(200, msg)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(200, gin.H{"message": "hey", "status": 200})
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```


####HTML rendering

Using LoadHTMLTemplates()

```go
func main() {
	r := gin.Default()
	r.LoadHTMLTemplates("templates/*")
	r.GET("/index", func(c *gin.Context) {
		obj := gin.H{"title": "Main website"}
		c.HTML(200, "index.tmpl", obj)
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

You can also use your own html template render

```go
import "html/template"

func main() {
	r := gin.Default()
	html := template.Must(template.ParseFiles("file1", "file2"))
	r.HTMLTemplates = html

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Redirects

Issuing a HTTP redirect is easy:

```go
r.GET("/test", func(c *gin.Context) {
	c.Redirect(301, "http://www.google.com/")
})
```
Both internal and external locations are supported.


#### Custom Middlewares

```go
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Set example variable
		c.Set("example", "12345")

		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// it would print: "12345"
		log.Println(example)
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### Using BasicAuth() middleware
```go
// similate some private data
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	r := gin.Default()

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets endpoint
	// hit "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// get user, it was setted by the BasicAuth middleware
		user := c.Get(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(200, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(200, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
```


#### Goroutines inside a middleware
When starting inside a middleware or handler, you **SHOULD NOT** use the original context inside it, you have to use a read-only copy.

```go
func main() {
	r := gin.Default()

	r.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		c_cp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note than you are using the copied context "c_cp", IMPORTANT
			log.Println("Done! in path " + c_cp.Request.URL.Path)
		}()
	})


	r.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// since we are NOT using a goroutine, we do not have to copy the context
		log.Println("Done! in path " + c.Request.URL.Path)
	})

    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
```

#### Custom HTTP configuration

Use `http.ListenAndServe()` directly, like this:

```go
func main() {
	router := gin.Default()
	http.ListenAndServe(":8080", router)
}
```
or

```go
func main() {
	router := gin.Default()

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
```
