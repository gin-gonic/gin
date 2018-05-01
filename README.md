# Gin Web Framework

<img align="right" width="159px" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">

[![Build Status](https://travis-ci.org/gin-gonic/gin.svg)](https://travis-ci.org/gin-gonic/gin)
 [![codecov](https://codecov.io/gh/gin-gonic/gin/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-gonic/gin)
 [![Go Report Card](https://goreportcard.com/badge/github.com/gin-gonic/gin)](https://goreportcard.com/report/github.com/gin-gonic/gin)
 [![GoDoc](https://godoc.org/github.com/gin-gonic/gin?status.svg)](https://godoc.org/github.com/gin-gonic/gin)
 [![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Open Source Helpers](https://www.codetriage.com/gin-gonic/gin/badges/users.svg)](https://www.codetriage.com/gin-gonic/gin)

Gin is a web framework written in Go (Golang). It features a martini-like API with much better performance, up to 40 times faster thanks to [httprouter](https://github.com/julienschmidt/httprouter). If you need performance and good productivity, you will love Gin.

![Gin console logger](https://gin-gonic.github.io/gin/other/console.png)

## Contents

- [Quick start](#quick-start)
- [Benchmarks](#benchmarks)
- [Gin v1.stable](#gin-v1-stable)
- [Start using it](#start-using-it)
- [Build with jsoniter](#build-with-jsoniter)
- [API Examples](#api-examples)
    - [Using GET,POST,PUT,PATCH,DELETE and OPTIONS](#using-get-post-put-patch-delete-and-options)
    - [Parameters in path](#parameters-in-path)
    - [Querystring parameters](#querystring-parameters)
    - [Multipart/Urlencoded Form](#multiparturlencoded-form)
    - [Another example: query + post form](#another-example-query--post-form)
    - [Upload files](#upload-files)
    - [Grouping routes](#grouping-routes)
    - [Blank Gin without middleware by default](#blank-gin-without-middleware-by-default)
    - [Using middleware](#using-middleware)
    - [How to write log file](#how-to-write-log-file)
    - [Model binding and validation](#model-binding-and-validation)
    - [Custom Validators](#custom-validators)
    - [Only Bind Query String](#only-bind-query-string)
    - [Bind Query String or Post Data](#bind-query-string-or-post-data)
    - [Bind HTML checkboxes](#bind-html-checkboxes)
    - [Multipart/Urlencoded binding](#multiparturlencoded-binding)
    - [XML, JSON and YAML rendering](#xml-json-and-yaml-rendering)
    - [JSONP rendering](#jsonp)
    - [Serving static files](#serving-static-files)
    - [HTML rendering](#html-rendering)
    - [Multitemplate](#multitemplate)
    - [Redirects](#redirects)
    - [Custom Middleware](#custom-middleware)
    - [Using BasicAuth() middleware](#using-basicauth-middleware)
    - [Goroutines inside a middleware](#goroutines-inside-a-middleware)
    - [Custom HTTP configuration](#custom-http-configuration)
    - [Support Let's Encrypt](#support-lets-encrypt)
    - [Run multiple service using Gin](#run-multiple-service-using-gin)
    - [Graceful restart or stop](#graceful-restart-or-stop)
    - [Build a single binary with templates](#build-a-single-binary-with-templates)
    - [Bind form-data request with custom struct](#bind-form-data-request-with-custom-struct)
- [Testing](#testing)
- [Users](#users--)

## Quick start
 
```sh
# assume the following codes in example.go file
$ cat example.go
```

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
```

```
# run example.go and visit 0.0.0.0:8080/ping on browser
$ go run example.go
```

## Benchmarks

Gin uses a custom version of [HttpRouter](https://github.com/julienschmidt/httprouter)

[See all benchmarks](/BENCHMARKS.md)

Benchmark name                              | (1)        | (2)         | (3) 		    | (4)
--------------------------------------------|-----------:|------------:|-----------:|---------:
**BenchmarkGin_GithubAll**                  | **30000**  |  **48375**  |     **0**  |   **0**
BenchmarkAce_GithubAll                      |   10000    |   134059    |   13792    |   167
BenchmarkBear_GithubAll                     |    5000    |   534445    |   86448    |   943
BenchmarkBeego_GithubAll                    |    3000    |   592444    |   74705    |   812
BenchmarkBone_GithubAll                     |     200    |  6957308    |  698784    |  8453
BenchmarkDenco_GithubAll                    |   10000    |   158819    |   20224    |   167
BenchmarkEcho_GithubAll                     |   10000    |   154700    |    6496    |   203
BenchmarkGocraftWeb_GithubAll               |    3000    |   570806    |  131656    |  1686
BenchmarkGoji_GithubAll                     |    2000    |   818034    |   56112    |   334
BenchmarkGojiv2_GithubAll                   |    2000    |  1213973    |  274768    |  3712
BenchmarkGoJsonRest_GithubAll               |    2000    |   785796    |  134371    |  2737
BenchmarkGoRestful_GithubAll                |     300    |  5238188    |  689672    |  4519
BenchmarkGorillaMux_GithubAll               |     100    | 10257726    |  211840    |  2272
BenchmarkHttpRouter_GithubAll               |   20000    |   105414    |   13792    |   167
BenchmarkHttpTreeMux_GithubAll              |   10000    |   319934    |   65856    |   671
BenchmarkKocha_GithubAll                    |   10000    |   209442    |   23304    |   843
BenchmarkLARS_GithubAll                     |   20000    |    62565    |       0    |     0
BenchmarkMacaron_GithubAll                  |    2000    |  1161270    |  204194    |  2000
BenchmarkMartini_GithubAll                  |     200    |  9991713    |  226549    |  2325
BenchmarkPat_GithubAll                      |     200    |  5590793    | 1499568    | 27435
BenchmarkPossum_GithubAll                   |   10000    |   319768    |   84448    |   609
BenchmarkR2router_GithubAll                 |   10000    |   305134    |   77328    |   979
BenchmarkRivet_GithubAll                    |   10000    |   132134    |   16272    |   167
BenchmarkTango_GithubAll                    |    3000    |   552754    |   63826    |  1618
BenchmarkTigerTonic_GithubAll               |    1000    |  1439483    |  239104    |  5374
BenchmarkTraffic_GithubAll                  |     100    | 11383067    | 2659329    | 21848
BenchmarkVulcan_GithubAll                   |    5000    |   394253    |   19894    |   609

- (1): Total Repetitions achieved in constant time, higher means more confident result
- (2): Single Repetition Duration (ns/op), lower is better
- (3): Heap Memory (B/op), lower is better
- (4): Average Allocations per Repetition (allocs/op), lower is better

## Gin v1. stable

- [x] Zero allocation router.
- [x] Still the fastest http router and framework. From routing to writing.
- [x] Complete suite of unit tests
- [x] Battle tested
- [x] API frozen, new releases will not break your code.

## Start using it

1. Download and install it:

```sh
$ go get github.com/gin-gonic/gin
```

2. Import it in your code:

```go
import "github.com/gin-gonic/gin"
```

3. (Optional) Import `net/http`. This is required for example if using constants such as `http.StatusOK`.

```go
import "net/http"
```

### Use a vendor tool like [Govendor](https://github.com/kardianos/govendor)

1. `go get` govendor

```sh
$ go get github.com/kardianos/govendor
```
2. Create your project folder and `cd` inside

```sh
$ mkdir -p $GOPATH/src/github.com/myusername/project && cd "$_"
```

3. Vendor init your project and add gin

```sh
$ govendor init
$ govendor fetch github.com/gin-gonic/gin@v1.2
```

4. Copy a starting template inside your project

```sh
$ curl https://raw.githubusercontent.com/gin-gonic/gin/master/examples/basic/main.go > main.go
```

5. Run your project

```sh
$ go run main.go
```

## Build with [jsoniter](https://github.com/json-iterator/go)

Gin use `encoding/json` as default json package but you can change to [jsoniter](https://github.com/json-iterator/go) by build from other tags.

```sh
$ go build -tags=jsoniter .
```

## API Examples

### Using GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
func main() {
	// Disable Console Color
	// gin.DisableConsoleColor()

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}
```

### Parameters in path

```go
func main() {
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// However, this one will match /user/john/ and also /user/john/send
	// If no other routers match /user/john, it will redirect to /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.Run(":8080")
}
```

### Querystring parameters

```go
func main() {
	router := gin.Default()

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

### Multipart/Urlencoded Form

```go
func main() {
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}
```

### Another example: query + post form

```
POST /post?id=1234&page=1 HTTP/1.1
Content-Type: application/x-www-form-urlencoded

name=manu&message=this_is_great
```

```go
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")

		fmt.Printf("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	router.Run(":8080")
}
```

```
id: 1234; page: 1; name: manu; message: this_is_great
```

### Upload files

#### Single file

References issue [#774](https://github.com/gin-gonic/gin/issues/774) and detail [example code](examples/upload-file/single).

```go
func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// single file
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	router.Run(":8080")
}
```

How to `curl`:

```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
  -H "Content-Type: multipart/form-data"
```

#### Multiple files

See the detail [example code](examples/upload-file/multiple).

```go
func main() {
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// Upload the file to specific dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":8080")
}
```

How to `curl`:

```bash
curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
```

### Grouping routes

```go
func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}
```

### Blank Gin without middleware by default

Use

```go
r := gin.New()
```

instead of

```go
// Default With the Logger and Recovery middleware already attached
r := gin.Default()
```


### Using middleware
```go
func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
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

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### How to write log file
```go
func main() {
    // Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console at the same time.
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

### Model binding and validation

To bind a request body into a type, use model binding. We currently support binding of JSON, XML and standard form values (foo=bar&boo=baz).

Gin uses [**go-playground/validator.v8**](https://github.com/go-playground/validator) for validation. Check the full docs on tags usage [here](http://godoc.org/gopkg.in/go-playground/validator.v8#hdr-Baked_In_Validators_and_Tags).

Note that you need to set the corresponding binding tag on all fields you want to bind. For example, when binding from JSON, set `json:"fieldname"`.

Also, Gin provides two sets of methods for binding:
- **Type** - Must bind
  - **Methods** - `Bind`, `BindJSON`, `BindQuery`
  - **Behavior** - These methods use `MustBindWith` under the hood. If there is a binding error, the request is aborted with `c.AbortWithError(400, err).SetType(ErrorTypeBind)`. This sets the response status code to 400 and the `Content-Type` header is set to `text/plain; charset=utf-8`. Note that if you try to set the response code after this, it will result in a warning `[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 400 with 422`. If you wish to have greater control over the behavior, consider using the `ShouldBind` equivalent method.
- **Type** - Should bind
  - **Methods** - `ShouldBind`, `ShouldBindJSON`, `ShouldBindQuery`
  - **Behavior** - These methods use `ShouldBindWith` under the hood. If there is a binding error, the error is returned and it is the developer's responsibility to handle the request and error appropriately.

When using the Bind-method, Gin tries to infer the binder depending on the Content-Type header. If you are sure what you are binding, you can use `MustBindWith` or `ShouldBindWith`.

You can also specify that specific fields are required. If a field is decorated with `binding:"required"` and has a empty value when binding, an error will be returned.

```go
// Binding from JSON
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	router := gin.Default()

	// Example for binding JSON ({"user": "manu", "password": "123"})
	router.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err == nil {
			if json.User == "manu" && json.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// Example for binding a HTML form (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err == nil {
			if form.User == "manu" && form.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
```

**Sample request**
```shell
$ curl -v -X POST \
  http://localhost:8080/loginJSON \
  -H 'content-type: application/json' \
  -d '{ "user": "manu" }'
> POST /loginJSON HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.51.0
> Accept: */*
> content-type: application/json
> Content-Length: 18
>
* upload completely sent off: 18 out of 18 bytes
< HTTP/1.1 400 Bad Request
< Content-Type: application/json; charset=utf-8
< Date: Fri, 04 Aug 2017 03:51:31 GMT
< Content-Length: 100
<
{"error":"Key: 'Login.Password' Error:Field validation for 'Password' failed on the 'required' tag"}
```

### Custom Validators

It is also possible to register custom validators. See the [example code](examples/custom-validation/server.go).

[embedmd]:# (examples/custom-validation/server.go go)
```go
package main

import (
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
)

type Booking struct {
	CheckIn  time.Time `form:"check_in" binding:"required,bookabledate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" binding:"required,gtfield=CheckIn" time_format:"2006-01-02"`
}

func bookableDate(
	v *validator.Validate, topStruct reflect.Value, currentStructOrField reflect.Value,
	field reflect.Value, fieldType reflect.Type, fieldKind reflect.Kind, param string,
) bool {
	if date, ok := field.Interface().(time.Time); ok {
		today := time.Now()
		if today.Year() > date.Year() || today.YearDay() > date.YearDay() {
			return false
		}
	}
	return true
}

func main() {
	route := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("bookabledate", bookableDate)
	}

	route.GET("/bookable", getBookable)
	route.Run(":8085")
}

func getBookable(c *gin.Context) {
	var b Booking
	if err := c.ShouldBindWith(&b, binding.Query); err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Booking dates are valid!"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
```

```console
$ curl "localhost:8085/bookable?check_in=2018-04-16&check_out=2018-04-17"
{"message":"Booking dates are valid!"}

$ curl "localhost:8085/bookable?check_in=2018-03-08&check_out=2018-03-09"
{"error":"Key: 'Booking.CheckIn' Error:Field validation for 'CheckIn' failed on the 'bookabledate' tag"}
```

[Struct level validations](https://github.com/go-playground/validator/releases/tag/v8.7) can also be registed this way.
See the [struct-lvl-validation example](examples/struct-lvl-validations) to learn more.

### Only Bind Query String

`ShouldBindQuery` function only binds the query params and not the post data. See the [detail information](https://github.com/gin-gonic/gin/issues/742#issuecomment-315953017).

```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name    string `form:"name"`
	Address string `form:"address"`
}

func main() {
	route := gin.Default()
	route.Any("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	if c.ShouldBindQuery(&person) == nil {
		log.Println("====== Only Bind By Query String ======")
		log.Println(person.Name)
		log.Println(person.Address)
	}
	c.String(200, "Success")
}

```

### Bind Query String or Post Data

See the [detail information](https://github.com/gin-gonic/gin/issues/742#issuecomment-264681292).

```go
package main

import "log"
import "github.com/gin-gonic/gin"
import "time"

type Person struct {
	Name     string    `form:"name"`
	Address  string    `form:"address"`
	Birthday time.Time `form:"birthday" time_format:"2006-01-02" time_utc:"1"`
}

func main() {
	route := gin.Default()
	route.GET("/testing", startPage)
	route.Run(":8085")
}

func startPage(c *gin.Context) {
	var person Person
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
	}

	c.String(200, "Success")
}
```

Test it with:
```sh
$ curl -X GET "localhost:8085/testing?name=appleboy&address=xyz&birthday=1992-03-15"
```

### Bind HTML checkboxes

See the [detail information](https://github.com/gin-gonic/gin/issues/129#issuecomment-124260092)

main.go

```go
...

type myForm struct {
    Colors []string `form:"colors[]"`
}

...

func formHandler(c *gin.Context) {
    var fakeForm myForm
    c.ShouldBind(&fakeForm)
    c.JSON(200, gin.H{"color": fakeForm.Colors})
}

...

```

form.html

```html
<form action="/" method="POST">
    <p>Check some colors</p>
    <label for="red">Red</label>
    <input type="checkbox" name="colors[]" value="red" id="red" />
    <label for="green">Green</label>
    <input type="checkbox" name="colors[]" value="green" id="green" />
    <label for="blue">Blue</label>
    <input type="checkbox" name="colors[]" value="blue" id="blue" />
    <input type="submit" />
</form>
```

result:

```
{"color":["red","green","blue"]}
```

### Multipart/Urlencoded binding

```go
package main

import (
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	User     string `form:"user" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func main() {
	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		// you can bind multipart form with explicit binding declaration:
		// c.ShouldBindWith(&form, binding.Form)
		// or you can simply use autobinding with ShouldBind method:
		var form LoginForm
		// in this case proper binding will be automatically selected
		if c.ShouldBind(&form) == nil {
			if form.User == "user" && form.Password == "password" {
				c.JSON(200, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(401, gin.H{"status": "unauthorized"})
			}
		}
	})
	router.Run(":8080")
}
```

Test it with:
```sh
$ curl -v --form user=user --form password=password http://localhost:8080/login
```

### XML, JSON and YAML rendering

```go
func main() {
	r := gin.Default()

	// gin.H is a shortcut for map[string]interface{}
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
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
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

#### SecureJSON

Using SecureJSON to prevent json hijacking. Default prepends `"while(1),"` to response body if the given struct is array values.

```go
func main() {
	r := gin.Default()

	// You can also use your own secure json prefix
	// r.SecureJsonPrefix(")]}',\n")

	r.GET("/someJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// Will output  :   while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```
#### JSONP

Using JSONP to request data from a server  in a different domain. Add callback to response body if the query parameter callback exists.

```go
func main() {
	r := gin.Default()

	r.GET("/JSONP?callback=x", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}
		
		//callback is x
		// Will output  :   x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### Serving static files

```go
func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// Listen and serve on 0.0.0.0:8080
	router.Run(":8080")
}
```

### HTML rendering

Using LoadHTMLGlob() or LoadHTMLFiles()

```go
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})
	router.Run(":8080")
}
```

templates/index.tmpl

```html
<html>
	<h1>
		{{ .title }}
	</h1>
</html>
```

Using templates with same name in different directories

```go
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/posts/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "posts/index.tmpl", gin.H{
			"title": "Posts",
		})
	})
	router.GET("/users/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users/index.tmpl", gin.H{
			"title": "Users",
		})
	})
	router.Run(":8080")
}
```

templates/posts/index.tmpl

```html
{{ define "posts/index.tmpl" }}
<html><h1>
	{{ .title }}
</h1>
<p>Using posts/index.tmpl</p>
</html>
{{ end }}
```

templates/users/index.tmpl

```html
{{ define "users/index.tmpl" }}
<html><h1>
	{{ .title }}
</h1>
<p>Using users/index.tmpl</p>
</html>
{{ end }}
```

#### Custom Template renderer

You can also use your own html template render

```go
import "html/template"

func main() {
	router := gin.Default()
	html := template.Must(template.ParseFiles("file1", "file2"))
	router.SetHTMLTemplate(html)
	router.Run(":8080")
}
```

#### Custom Delimiters

You may use custom delims

```go
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.LoadHTMLGlob("/path/to/templates"))
```

#### Custom Template Funcs

See the detail [example code](examples/template).

main.go

```go
import (
    "fmt"
    "html/template"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
    year, month, day := t.Date()
    return fmt.Sprintf("%d%02d/%02d", year, month, day)
}

func main() {
    router := gin.Default()
    router.Delims("{[{", "}]}")
    router.SetFuncMap(template.FuncMap{
        "formatAsDate": formatAsDate,
    })
    router.LoadHTMLFiles("./fixtures/basic/raw.tmpl")

    router.GET("/raw", func(c *gin.Context) {
        c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
            "now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
        })
    })

    router.Run(":8080")
}

```

raw.tmpl

```html
Date: {[{.now | formatAsDate}]}
```

Result:
```
Date: 2017/07/01
```

### Multitemplate

Gin allow by default use only one html.Template. Check [a multitemplate render](https://github.com/gin-contrib/multitemplate) for using features like go 1.6 `block template`.

### Redirects

Issuing a HTTP redirect is easy:

```go
r.GET("/test", func(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
})
```
Both internal and external locations are supported.


### Custom Middleware

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

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### Using BasicAuth() middleware

```go
// simulate some private data
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
		// get user, it was set by the BasicAuth middleware
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### Goroutines inside a middleware

When starting new Goroutines inside a middleware or handler, you **SHOULD NOT** use the original context inside it, you have to use a read-only copy.

```go
func main() {
	r := gin.Default()

	r.GET("/long_async", func(c *gin.Context) {
		// create copy to be used inside the goroutine
		cCp := c.Copy()
		go func() {
			// simulate a long task with time.Sleep(). 5 seconds
			time.Sleep(5 * time.Second)

			// note that you are using the copied context "cCp", IMPORTANT
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// simulate a long task with time.Sleep(). 5 seconds
		time.Sleep(5 * time.Second)

		// since we are NOT using a goroutine, we do not have to copy the context
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### Custom HTTP configuration

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

### Support Let's Encrypt

example for 1-line LetsEncrypt HTTPS servers.

[embedmd]:# (examples/auto-tls/example1/main.go go)
```go
package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
}
```

example for custom autocert manager.

[embedmd]:# (examples/auto-tls/example2/main.go go)
```go
package main

import (
	"log"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("example1.com", "example2.com"),
		Cache:      autocert.DirCache("/var/www/.cache"),
	}

	log.Fatal(autotls.RunWithManager(r, &m))
}
```

### Run multiple service using Gin

See the [question](https://github.com/gin-gonic/gin/issues/346) and try the following example:

[embedmd]:# (examples/multiple-service/main.go go)
```go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func router01() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 01",
			},
		)
	})

	return e
}

func router02() http.Handler {
	e := gin.New()
	e.Use(gin.Recovery())
	e.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":  http.StatusOK,
				"error": "Welcome server 02",
			},
		)
	})

	return e
}

func main() {
	server01 := &http.Server{
		Addr:         ":8080",
		Handler:      router01(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server02 := &http.Server{
		Addr:         ":8081",
		Handler:      router02(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return server01.ListenAndServe()
	})

	g.Go(func() error {
		return server02.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
```

### Graceful restart or stop

Do you want to graceful restart or stop your web server?
There are some ways this can be done.

We can use [fvbock/endless](https://github.com/fvbock/endless) to replace the default `ListenAndServe`. Refer issue [#296](https://github.com/gin-gonic/gin/issues/296) for more details.

```go
router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```

An alternative to endless:

* [manners](https://github.com/braintree/manners): A polite Go HTTP server that shuts down gracefully.
* [graceful](https://github.com/tylerb/graceful): Graceful is a Go package enabling graceful shutdown of an http.Handler server.
* [grace](https://github.com/facebookgo/grace): Graceful restart & zero downtime deploy for Go servers.

If you are using Go 1.8, you may not need to use this library! Consider using http.Server's built-in [Shutdown()](https://golang.org/pkg/net/http/#Server.Shutdown) method for graceful shutdowns. See the full [graceful-shutdown](./examples/graceful-shutdown) example with gin.

[embedmd]:# (examples/graceful-shutdown/graceful-shutdown/server.go go)
```go
// +build go1.8

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server")
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
```

### Build a single binary with templates

You can build a server into a single binary containing templates by using [go-assets][].

[go-assets]: https://github.com/jessevdk/go-assets

```go
func main() {
	r := gin.New()

	t, err := loadTemplate()
	if err != nil {
		panic(err)
	}
	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "/html/index.tmpl",nil)
	})
	r.Run(":8080")
}

// loadTemplate loads templates embedded by go-assets-builder
func loadTemplate() (*template.Template, error) {
	t := template.New("")
	for name, file := range Assets.Files {
		if file.IsDir() || !strings.HasSuffix(name, ".tmpl") {
			continue
		}
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(name).Parse(string(h))
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}
```

See a complete example in the `examples/assets-in-binary` directory.

### Bind form-data request with custom struct

The follow example using custom struct:

```go
type StructA struct {
    FieldA string `form:"field_a"`
}

type StructB struct {
    NestedStruct StructA
    FieldB string `form:"field_b"`
}

type StructC struct {
    NestedStructPointer *StructA
    FieldC string `form:"field_c"`
}

type StructD struct {
    NestedAnonyStruct struct {
        FieldX string `form:"field_x"`
    }
    FieldD string `form:"field_d"`
}

func GetDataB(c *gin.Context) {
    var b StructB
    c.Bind(&b)
    c.JSON(200, gin.H{
        "a": b.NestedStruct,
        "b": b.FieldB,
    })
}

func GetDataC(c *gin.Context) {
    var b StructC
    c.Bind(&b)
    c.JSON(200, gin.H{
        "a": b.NestedStructPointer,
        "c": b.FieldC,
    })
}

func GetDataD(c *gin.Context) {
    var b StructD
    c.Bind(&b)
    c.JSON(200, gin.H{
        "x": b.NestedAnonyStruct,
        "d": b.FieldD,
    })
}

func main() {
    r := gin.Default()
    r.GET("/getb", GetDataB)
    r.GET("/getc", GetDataC)
    r.GET("/getd", GetDataD)

    r.Run()
}
```

Using the command `curl` command result:

```
$ curl "http://localhost:8080/getb?field_a=hello&field_b=world"
{"a":{"FieldA":"hello"},"b":"world"}
$ curl "http://localhost:8080/getc?field_a=hello&field_c=world"
{"a":{"FieldA":"hello"},"c":"world"}
$ curl "http://localhost:8080/getd?field_x=hello&field_d=world"
{"d":"world","x":{"FieldX":"hello"}}
```

**NOTE**: NOT support the follow style struct:

```go
type StructX struct {
    X struct {} `form:"name_x"` // HERE have form
}

type StructY struct {
    Y StructX `form:"name_y"` // HERE hava form
}

type StructZ struct {
    Z *StructZ `form:"name_z"` // HERE hava form
}
```

In a word, only support nested custom struct which have no `form` now.

## Testing

The `net/http/httptest` package is preferable way for HTTP testing.

```go
package main

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8080")
}
```

Test for code example above:

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
```

## Users  [![Sourcegraph](https://sourcegraph.com/github.com/gin-gonic/gin/-/badge.svg)](https://sourcegraph.com/github.com/gin-gonic/gin?badge)

Awesome project lists using [Gin](https://github.com/gin-gonic/gin) web framework.

* [drone](https://github.com/drone/drone): Drone is a Continuous Delivery platform built on Docker, written in Go
* [gorush](https://github.com/appleboy/gorush): A push notification server written in Go.
