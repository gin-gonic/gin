#Gin Web Framework
Gin is a web framework written in Golang. It features a martini-like API with much better performance, up to 40 times faster. If you need performance and good productivity, you will love Gin.  
[Check out the official web site](http://gin-gonic.github.io/gin/)

## Start using it
Run:

```
go get github.com/gin-gonic/gin
```
Then import it in your Golang code:

```
import "github.com/gin-gonic/gin"
```


##API Examples

#### Create most basic PING/PONG HTTP endpoint
```go 
import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("ping", func(c *gin.Context){
        c.String("pong")
    })
    
    // Listen and server on 0.0.0.0:8080
    r.Run(":80")
}
```

#### Using GET, POST, PUT, PATCH and DELETE

```go
func main() {
    // Creates a gin router + logger and recovery (crash-free) middlwares
    r := gin.Default()
    
    r.GET("/someGet", getting)
    r.POST("/somePost", posting)
    r.PUT("/somePut", putting)
    r.DELETE("/someDelete", deleting)
    r.PATCH("/somePATCH", patching)

    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
```

#### Parameters in path

```go
func main() {
    r := gin.Default()
    
    r.GET("/user/:name", func(c *gin.Context) {
        name := c.Params.ByName("name")
        message := "Hello "+name
        c.String(200, message)
    })
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
    
    // Simple group: v1
    v2 := r.Group("/v2")
    {
        v2.POST("/login", loginEndpoint)
        v2.POST("/submit"", submitEndpoint)
        v2.POST("/read"", readEndpoint)
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
    // Creates a router without any middlware by default
    r := gin.New()
    
    // Global middlwares
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    
    // Per route middlwares, you can add as many as you desire.
    r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

    // Authorization group
    // authorized := r.Group("/", AuthRequired())
    // exactly the same than:
    authorized := r.Group("/")
    // per group middlwares! in this case we use the custom created
    // AuthRequired() middlware just in the "authorized" group.
    authorized.Use(AuthRequired())
    {
        authorized.Use.POST("/login", loginEndpoint)
        authorized.Use.POST("/submit", submitEndpoint)
        authorized.Use.POST("/read", readEndpoint)
        
        // nested group
        testing := authorized.Group("testing")
        testing.GET("/analytics", analyticsEndpoint)
    }
   
    // Listen and server on 0.0.0.0:8080
    r.Run(":8080")
}
```


#### JSON parsing and validation

```go
type LoginJSON struct {
    User     string `json:"user" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func main() {
    r := gin.Default()
    
    r.POST("/login", func(c *gin.Context) {
        var json LoginJSON
        
        // If EnsureBody returns false, it will write automatically the error
        // in the HTTP stream and return a 400 error. If you want custom error 
        // handling you should use: c.ParseBody(interface{}) error
        if c.EnsureBody(&json) {
            if json.User=="manu" && json.Password=="123" {
                c.JSON(200, gin.H{"status": "you are logged in"})
            }else{
                c.JSON(401, gin.H{"status": "unauthorized"})
            }
        }
    })
}
```

#### XML, and JSON rendering

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
            Message string
            Status  int
        }
        msg.Message = "hey"
        msg.Status = 200
        c.JSON(200, msg.Status)
    })
    
    r.GET("/someXML", func(c *gin.Context) {
        c.XML(200, gin.H{"message": "hey", "status": 200})
    })
}
```


####HTML rendering

Using LoadHTMLTemplates()

```go
func main() {
    r := gin.Default()
    r.LoadHTMLTemplates("templates/*")
    r.GET("index", func(c *gin.Context) {
        obj := gin.h{"title": "Main website"}
        c.HTML(200, "templates/index.tmpl", obj)
    })
}
```

You can also use your own html template render

```go
import "html/template"
func main() {
    r := gin.Default()
    html := template.ParseFiles("file1", "file2")
    r.HTMLTemplates = html
}
```


#### Custom Middlewares

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t : time.Now()
        
        // Set example variable
        c.Set("example", "12345")
        
        // before request
        
        c.Next()
        
        // after request
        latency := time.Since(t)
        log.Print(latency)
    }
}

func main() {
    r := gin.New()
    r.Use(Logger())
    
    r.GET("test", func(c *gin.Context){
        example := r.Get("example").(string)
        
        // it would print: "12345"
        log.Println(example)
    })
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