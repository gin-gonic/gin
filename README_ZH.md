# Gin Web Framework

<img align="right" width="159px" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">

[![Build Status](https://travis-ci.org/gin-gonic/gin.svg)](https://travis-ci.org/gin-gonic/gin)
 [![codecov](https://codecov.io/gh/gin-gonic/gin/branch/master/graph/badge.svg)](https://codecov.io/gh/gin-gonic/gin)
 [![Go Report Card](https://goreportcard.com/badge/github.com/gin-gonic/gin)](https://goreportcard.com/report/github.com/gin-gonic/gin)
 [![GoDoc](https://godoc.org/github.com/gin-gonic/gin?status.svg)](https://godoc.org/github.com/gin-gonic/gin)
 [![Join the chat at https://gitter.im/gin-gonic/gin](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/gin-gonic/gin?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Open Source Helpers](https://www.codetriage.com/gin-gonic/gin/badges/users.svg)](https://www.codetriage.com/gin-gonic/gin)

Gin 是一个用 Go 语言编写的 WEB 框架。它具有和 maritini 类似的 API 并拥有更好的性能， 感谢 [httprouter](https://github.com/julienschmidt/httprouter) 使它的速度提升了 40 倍。如果你需要性能和良好的生产力，你将会爱上 Gin 。

![Gin console logger](https://gin-gonic.github.io/gin/other/console.png)

## Contents

- [安装](#安装)
- [依赖](#依赖)
- [快速开始](#快速开始)
- [性能测试](#性能测试)
- [Gin v1.stable](#gin-v1-stable)
- [使用 jsoniter 构建](#使用-jsoniter-构建)
- [API 示例](#api-示例)
    - [使用 GET,POST,PUT,PATCH,DELETE and OPTIONS](#使用-get-post-put-patch-delete-and-options)
    - [path 中的参数](#path-中的参数)
    - [请求参数](#请求参数)
    - [Multipart/Urlencoded 表单](#Multipart/Urlencoded-表单)
    - [其他示例： 请求参数 + post form](#其他示例-请求参数--post-form)
    - [上传文件](#上传文件)
    - [组路由](#组路由)
    - [默认的没有中间件的空白 Gin](#默认的没有中间件的空白-Gin)
    - [使用中间件](#使用中间件)
    - [如何写入日志文件](#如何写入日志文件)
    - [模型绑定和验证](#模型绑定和验证)
    - [自定义验证器](#自定义验证器)
    - [只绑定查询字符串](#只绑定查询字符串)
    - [绑定查询字符串或 post 数据](#绑定查询字符串或-post-数据)
    - [绑定 HTML 复选框](#绑定-HTML-复选框)
    - [Multipart/Urlencoded 绑定](#Multipart/Urlencoded-绑定)
    - [XML, JSON 和 YAML 渲染](#XML-JSON-和-YAML-渲染)
    - [JSONP 渲染](#jsonp)
    - [静态文件服务](#静态文件服务)
    - [从 reader 提供数据](#从-reader-提供数据)
    - [HTML 渲染](#HTML-渲染)
    - [多模板](#多模板)
    - [重定向](#重定向)
    - [自定义中间件](#自定义中间件)
    - [使用 BasicAuth() 中间件](#使用-BasicAuth-中间件)
    - [在中间件中使用协成](#在中间件中使用协成)
    - [自定义 HTTP 配置](#自定义-HTTP-配置)
    - [支持 Let's Encrypt](#支持-Let's-Encrypt)
    - [使用 Gin 运行多个服务](#使用-Gin-运行多个服务)
    - [正常的重启或停止](#正常的重启或停止)
    - [使用模板构建单个二进制文件](#使用模板构建单个二进制文件)
    - [使用自定义结构绑定表单数据请求](#使用自定义结构绑定表单数据请求)
    - [尝试将 body 绑定到不同的结构中](#尝试将-body-绑定到不同的结构中)
    - [HTTP2 服务器推送](#HTTP2-服务器推送)
- [测试](#测试)
- [使用者](#使用者--)

## 安装

要安装 Gin 包，你需要先安装 Go 并且设置你的 Go 的工作工作空间。

1. 下载并安装它：

```sh
$ go get -u github.com/gin-gonic/gin
```

2. 在你的代码中导入它：

```go
import "github.com/gin-gonic/gin"
```

3. （可选的） 导入 `net/http` 。 如果你使用常量（如： `http.StatusOK` ） 的时候必须导入。

```go
import "net/http"
```

### 使用一个 vendor 工具，比如 [Govendor](https://github.com/kardianos/govendor)

1. `go get` govendor

```sh
$ go get github.com/kardianos/govendor
```
2. 创建你的项目目录并使用 `cd` 进入

```sh
$ mkdir -p $GOPATH/src/github.com/myusername/project && cd "$_"
```

3. 初始化你的项目的 Vendor 并添加 gin

```sh
$ govendor init
$ govendor fetch github.com/gin-gonic/gin@v1.2
```

4. 复制一个开始模板到你的项目中

```sh
$ curl https://raw.githubusercontent.com/gin-gonic/gin/master/examples/basic/main.go > main.go
```

5. 运行你的项目

```sh
$ go run main.go
```

## 依赖

现在 Gin 需要 Go 1.6 或更高版本，并且它就将会需要 Go 1.7 。

## 快速开始

```sh
# 假设下面代码在 example.go 文件中
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
	r.Run() // 在 0.0.0.0:8080 上监听并服务
}
```

```
# 运行 example.go 并在浏览器上访问 0.0.0.0:8080/ping
$ go run example.go
```

## 性能测试

Gin 采用一个 [HttpRouter](https://github.com/julienschmidt/httprouter) 的自定义版本

[查看所有性能测试](/BENCHMARKS.md)

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

- （1）： 在固定时间内重复完成的总数，数值越高的是越好的结果
- （2）： 单次重复持续时间（ns / op），越低越好
- （3）： 堆内存（B / op），越低越好
- （4）： 每次重复的平均分配（allocs / op），越低越好

## Gin v1. stable

- [x] 零分配路由
- [x] 仍然是最快的 http 路由器和框架， 从路由写入。
- [x] 完整的单元测试套件
- [x] 对战测试
- [x] API冻结，新版本不会破坏您的代码。

## 使用 [jsoniter](https://github.com/json-iterator/go) 构建

Gin 使用 `encoding/json` 作为默认的 json 包，但是你可以在构建的时候通过 tags 去使用 [jsoniter](https://github.com/json-iterator/go) 。

```sh
$ go build -tags=jsoniter .
```

## API 示例

### 使用 GET, POST, PUT, PATCH, DELETE and OPTIONS

```go
func main() {
	// 禁用控制台颜色
	// gin.DisableConsoleColor()

	// 使用默认中间件创建一个 gin 路由：
	// 日志与恢复中间件（无崩溃）。
	router := gin.Default()

	router.GET("/someGet", getting)
	router.POST("/somePost", posting)
	router.PUT("/somePut", putting)
	router.DELETE("/someDelete", deleting)
	router.PATCH("/somePatch", patching)
	router.HEAD("/someHead", head)
	router.OPTIONS("/someOptions", options)

	// 默认情况下，它使用：8080，除非定义了 PORT 环境变量。
	router.Run()
	// router.Run(":3000") 硬编码端口
}
```

### path 中的参数

```go
func main() {
	router := gin.Default()

	// 这个处理器将去匹配 /user/john ， 但是它不会去匹配 /user
	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	// 然而，这个将会匹配 /user/john 并且也会匹配 /user/john/send
	// 如果没有其他的路由匹配 /user/john ， 它将重定向到 /user/john/
	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.Run(":8080")
}
```

### 请求参数

```go
func main() {
	router := gin.Default()

	// 请求参数使用现有的底层 request 对象解析。
	// 请求响应匹配的 URL： /welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		// 这个是 c.Request.URL.Query().Get("lastname") 的快捷方式。
		lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

### Multipart/Urlencoded 表单

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

### 其他示例： 请求参数 + post form

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

### 上传文件

#### 单文件

参考 issue [#774](https://github.com/gin-gonic/gin/issues/774) 与详细的 [示例代码](https://github.com/gin-gonic/gin/tree/master/examples/upload-file/single) 。

```go
func main() {
	router := gin.Default()
	// 为 multipart 表单设置一个较低的内存限制（默认是 32 MiB）
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// 单文件
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		// 上传文件到指定的 dst 。
		// c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	router.Run(":8080")
}
```

如何使用 `curl`:

```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@/Users/appleboy/test.zip" \
  -H "Content-Type: multipart/form-data"
```

#### 多文件

See the detail [example code](examples/upload-file/multiple).

```go
func main() {
	router := gin.Default()
	// 为 multipart 表单设置一个较低的内存限制（默认是 32 MiB）
	// router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			// 上传文件到指定的 dst.
			// c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":8080")
}
```

如何使用 `curl`:

```bash
curl -X POST http://localhost:8080/upload \
  -F "upload[]=@/Users/appleboy/test1.zip" \
  -F "upload[]=@/Users/appleboy/test2.zip" \
  -H "Content-Type: multipart/form-data"
```

### 组路由

```go
func main() {
	router := gin.Default()

	// 简单组： v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// 简单组： v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}
```

### 默认的没有中间件的空白 Gin

Use

```go
r := gin.New()
```

代替

```go
// Default 已经连接了 Logger 和 Recovery 中间件
r := gin.Default()
```


### 使用中间件
```go
func main() {
	// 创建一个默认的没有任何中间件的路由
	r := gin.New()

	// 全局中间件
	// Logger 中间件将写日志到 gin.DefaultWriter ,即使你设置 GIN_MODE=release 。
	// 默认 gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery 中间件从任何 panic 恢复，如果出现 panic，它会写一个 500 错误。
	r.Use(gin.Recovery())

	// 每个路由的中间件, 你能添加任意数量的中间件
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// 授权组
	// authorized := r.Group("/", AuthRequired())
	// 也可以这样:
	authorized := r.Group("/")
	// 每个组的中间件！ 在这个实例中，我们只需要在 "authorized" 组中
	// 使用自定义创建的 AuthRequired() 中间件
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

### 如何写入日志文件
```go
func main() {
    // 禁用控制台颜色，当你将日志写入到文件的时候，你不需要控制台颜色。
    gin.DisableConsoleColor()

    // 写入日志的文件
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // 如果你需要同时写入日志文件和控制台上显示，使用下面代码
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

### 模型绑定和验证

绑定一个请求主体到一个类型，使用模型绑定。我们目前支持 JSON 、 XML 和标准表单的值（ foo=bar&boo=baz ）的绑定。

Gin 使用 [**go-playground/validator.v8**](https://github.com/go-playground/validator) 进行验证。 在 [这里](http://godoc.org/gopkg.in/go-playground/validator.v8#hdr-Baked_In_Validators_and_Tags) 查看标签使用的完整文档。

注意： 你需要在你想要绑定的所有字段上设置相应的绑定标签。例如：当要从 JSON 绑定的时候，设置 `json:"fieldname"` 。

此外， Gin 提供了两组绑定方法：

- **类型** - Must bind
  - **方法** - `Bind`, `BindJSON`, `BindQuery`
  - **行为** - 这些方法在 `MustBindWith` 引擎下面使用。如果存在绑定错误， 请求通过 `c.AbortWithError(400, err).SetType(ErrorTypeBind)` 被终止。 这组响应的状态吗被设置成 400 ，并将 `Content-Type` 头设置成 `text/plain; charset=utf-8` 。注意： 如果你尝试在这个之后去设置响应码，它会发出一个警告  `[GIN-debug] [WARNING] Headers were already written. Wanted to override status code 400 with 422` 。 如果你希望更好的控制行为， 请考虑使用 `ShouldBind` 等效的方法。
- **类型** - Should bind
  - **方法** - `ShouldBind`, `ShouldBindJSON`, `ShouldBindQuery`
  - **行为** - 这些方法在 `ShouldBindWith` 引擎下使用。 如果存在绑定错误，这个错误会被返回， 需要开发者去处理相应的请求和错误。

当使用绑定方式时， Gin 会尝试通过 Content-Type 推断出绑定器的依赖,如果你要明确你绑定的什么，可以使用 `MustBindWith` 或 `ShouldBindWith` 。

你也可以指定需要指定的字段。如果一个字段使用 `binding:"required"` 修饰，并且当绑定的时候是一个空值的时候，将会返回一个错误。

```go
// 从 JSON 绑定
type Login struct {
	User     string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func main() {
	router := gin.Default()

	// 绑定 JSON 的示例 ({"user": "manu", "password": "123"})
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

	// 一个 HTML 表单绑定的示例 (user=manu&password=123)
	router.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// 这个将通过 content-type 头去推断绑定器使用哪个依赖。
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

	// 监听并服务于 0.0.0.0:8080
	router.Run(":8080")
}
```

**请求样本**
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

**跳过验证**

当在命令行上使用 `curl` 运行上面的示例时，它会返回一个错误。因为示例给 `Password` 绑定了 `binding:"required"` 。如果 `Password` 使用 `binding:"-"` ，然后再次运行上面的示例，它将不会返回错误。

### 自定义验证器

也可以注册自定义验证器。 参见 [示例代码](examples/custom-validation/server.go) 。

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

[结构级别的验证](https://github.com/go-playground/validator/releases/tag/v8.7) 也可以这样注册。
查看 [示例 struct-lvl-validation ](examples/struct-lvl-validations) 学习更多。

### 只绑定查询字符串

`ShouldBindQuery` 函数只绑定查询参数并且没有 post 数据。查看 [详细信息](https://github.com/gin-gonic/gin/issues/742#issuecomment-315953017).

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

### 绑定查询字符串或 post 数据

查看 [详细信息](https://github.com/gin-gonic/gin/issues/742#issuecomment-264681292).

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
	// 如果是 `GET`, 只使用 `Form` 绑定引擎 (`query`) 。
	// 如果 `POST`, 首先检查 `content-type` 为 `JSON` 或 `XML`, 然后使用 `Form` (`form-data`) 。
	// 在这里查看更多信息 https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	if c.ShouldBind(&person) == nil {
		log.Println(person.Name)
		log.Println(person.Address)
		log.Println(person.Birthday)
	}

	c.String(200, "Success")
}
```

测试它：
```sh
$ curl -X GET "localhost:8085/testing?name=appleboy&address=xyz&birthday=1992-03-15"
```

### 绑定 HTML 复选框

查看 [详细信息](https://github.com/gin-gonic/gin/issues/129#issuecomment-124260092)

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

### Multipart/Urlencoded 绑定

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
		// 你可以使用显示绑定声明绑定 multipart 表单：
		// c.ShouldBindWith(&form, binding.Form
		// 或者你可以使用 ShouldBind 方法去简单的使用自动绑定：
		var form LoginForm
		// 在这种情况下，将自动选择适合的绑定
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

测试它：
```sh
$ curl -v --form user=user --form password=password http://localhost:8080/login
```

### XML, JSON 和 YAML 渲染

```go
func main() {
	r := gin.Default()

	// gin.H 是一个 map[string]interface{} 的快捷方式
	r.GET("/someJSON", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/moreJSON", func(c *gin.Context) {
		// 你也可以使用一个结构
		var msg struct {
			Name    string `json:"user"`
			Message string
			Number  int
		}
		msg.Name = "Lena"
		msg.Message = "hey"
		msg.Number = 123
		// 注意 msg.Name 在 JSON 中会变成 "user"
		// 将会输出： {"user": "Lena", "Message": "hey", "Number": 123}
		c.JSON(http.StatusOK, msg)
	})

	r.GET("/someXML", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	r.GET("/someYAML", func(c *gin.Context) {
		c.YAML(http.StatusOK, gin.H{"message": "hey", "status": http.StatusOK})
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

#### SecureJSON

使用 SecureJSON 来防止 json 劫持。如果给定的结构体是数组值，默认预置 `"while(1),"` 到 response body 。

```go
func main() {
	r := gin.Default()

	// 你也可以使用自己的安装 json 前缀
	// r.SecureJsonPrefix(")]}',\n")

	r.GET("/someJSON", func(c *gin.Context) {
		names := []string{"lena", "austin", "foo"}

		// 将会输出  :   while(1);["lena","austin","foo"]
		c.SecureJSON(http.StatusOK, names)
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```
#### JSONP

在不同的域中使用 JSONP 从一个服务器请求数据。如果请求参数中存在 callback，添加 callback 到 response body 。

```go
func main() {
	r := gin.Default()

	r.GET("/JSONP?callback=x", func(c *gin.Context) {
		data := map[string]interface{}{
			"foo": "bar",
		}

		//callback 是 x
		// 将会输出  :   x({\"foo\":\"bar\"})
		c.JSONP(http.StatusOK, data)
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

#### AsciiJSON

使用 AsciiJSON 为非 ASCII 字符生成仅有 ASCII 字符的 JSON 。

```go
func main() {
	r := gin.Default()

	r.GET("/someJSON", func(c *gin.Context) {
		data := map[string]interface{}{
			"lang": "GO语言",
			"tag":  "<br>",
		}

		// 将会输出 : {"lang":"GO\u8bed\u8a00","tag":"\u003cbr\u003e"}
		c.AsciiJSON(http.StatusOK, data)
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

### 静态文件服务

```go
func main() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.StaticFS("/more_static", http.Dir("my_file_system"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	// 监听并服务于 0.0.0.0:8080
	router.Run(":8080")
}
```

### 从 reader 提供数据

```go
func main() {
	router := gin.Default()
	router.GET("/someDataFromReader", func(c *gin.Context) {
		response, err := http.Get("https://raw.githubusercontent.com/gin-gonic/logo/master/color.png")
		if err != nil || response.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}

		reader := response.Body
		contentLength := response.ContentLength
		contentType := response.Header.Get("Content-Type")

		extraHeaders := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}

		c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
	})
	router.Run(":8080")
}
```

### HTML 渲染

使用 LoadHTMLGlob() 或 LoadHTMLFiles()

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

在不同的目录使用具有相同名称的模板

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

#### 自定义模板渲染器

你也可以使用你自己的 HTML 模板渲染

```go
import "html/template"

func main() {
	router := gin.Default()
	html := template.Must(template.ParseFiles("file1", "file2"))
	router.SetHTMLTemplate(html)
	router.Run(":8080")
}
```

#### 自定义分隔符

你可以使用自定义分隔符

```go
	r := gin.Default()
	r.Delims("{[{", "}]}")
	r.LoadHTMLGlob("/path/to/templates"))
```

#### 自定义模板函数

查看详细的 [示例代码](examples/template).

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

### 多模板

Gin 允许默认只使用一个 html.Template 。查看 [多模板渲染](https://github.com/gin-contrib/multitemplate) 的使用详情，类似 go 1.6 `block template`

### 重定向

发出一个 HTTP 重定向非常容易， 同时支持内部和外部地址。

```go
r.GET("/test", func(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://www.google.com/")
})
```

发出路由重定向，使用下面示例中的 `HandleContext` 。

``` go
r.GET("/test", func(c *gin.Context) {
    c.Request.URL.Path = "/test2"
    r.HandleContext(c)
})
r.GET("/test2", func(c *gin.Context) {
    c.JSON(200, gin.H{"hello": "world"})
})
```


### 自定义中间件

```go
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// 设置简单的变量
		c.Set("example", "12345")

		// 在请求之前

		c.Next()

		// 在请求之后
		latency := time.Since(t)
		log.Print(latency)

		// 记录我们的访问状态
		status := c.Writer.Status()
		log.Println(status)
	}
}

func main() {
	r := gin.New()
	r.Use(Logger())

	r.GET("/test", func(c *gin.Context) {
		example := c.MustGet("example").(string)

		// 它将打印： "12345"
		log.Println(example)
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

### 使用 BasicAuth() 中间件

```go
// 模拟一些私有的数据
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	r := gin.Default()

	// 在组中使用 gin.BasicAuth() 中间件
	// gin.Accounts 是 map[string]string 的快捷方式
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// /admin/secrets 结尾
	// 点击 "localhost:8080/admin/secrets
	authorized.GET("/secrets", func(c *gin.Context) {
		// 获取 user, 它是由 BasicAuth 中间件设置的
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

### 在中间件中使用协成

在一个中间件或处理器中启动一个新的协成时，你 **不应该** 使用它里面的原始的 context ，只能去使用它的只读副本。

```go
func main() {
	r := gin.Default()

	r.GET("/long_async", func(c *gin.Context) {
		// 创建在协成中使用的副本
		cCp := c.Copy()
		go func() {
			// 使用 time.Sleep() 休眠 5 秒，模拟一个用时长的任务。
			time.Sleep(5 * time.Second)

			// 注意，你使用的是复制的 context "cCp" ，重要
			log.Println("Done! in path " + cCp.Request.URL.Path)
		}()
	})

	r.GET("/long_sync", func(c *gin.Context) {
		// 使用 time.Sleep() 休眠 5 秒，模拟一个用时长的任务。
		time.Sleep(5 * time.Second)

		// 因为我们没有使用协成，我们不需要复制 context
		log.Println("Done! in path " + c.Request.URL.Path)
	})

	// 监听并服务于 0.0.0.0:8080
	r.Run(":8080")
}
```

### 自定义 HTTP 配置

直接使用 `http.ListenAndServe()` ，像这样：

```go
func main() {
	router := gin.Default()
	http.ListenAndServe(":8080", router)
}
```
或

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

### 支持 Let's Encrypt

一个 LetsEncrypt HTTPS 服务器的示例。

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

	// Ping 处理器
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	log.Fatal(autotls.Run(r, "example1.com", "example2.com"))
}
```

自定义 autocert 管理器示例。

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

### 使用 Gin 运行多个服务

查看 [问题](https://github.com/gin-gonic/gin/issues/346) 并尝试下面示例：

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

### 正常的重启或停止

你想正常的重启或停止你的 web 服务器吗？
有一些方法可以做到。

我们能使用 [fvbock/endless](https://github.com/fvbock/endless) 去替换默认的 `ListenAndServe`. 参考 issue [#296](https://github.com/gin-gonic/gin/issues/296) 了解更多细节。

```go
router := gin.Default()
router.GET("/", handler)
// [...]
endless.ListenAndServe(":4242", router)
```

另外一些替代方案：

* [manners](https://github.com/braintree/manners)　： 一个有礼貌的 Go HTTP服务器，它可以正常的关闭。
* [graceful](https://github.com/tylerb/graceful)　：　Graceful 是一个 Go 包，它可以正常的关闭一个　http.Handler 服务器。
* [grace](https://github.com/facebookgo/grace)　：　正常的重启 & Go　服务器零停机部署。

如果你使用的是　Go 1.8，你可能不需要使用这些库！考虑使用　http.Server　内置的　[Shutdown()](https://golang.org/pkg/net/http/#Server.Shutdown)　方法正常关闭。查看 Gin 中完整的　[graceful-shutdown](./examples/graceful-shutdown) 示例。

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
		// 连接服务器
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号超时　５　秒正常关闭服务器
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

### 使用模板构建单个二进制文件

你可以使用　[go-assets](https://github.com/jessevdk/go-assets)　将服务器构建到一个包含模板的单独的二进制文件中。

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

// loadTemplate 加载　go-assets-builder　嵌入的模板
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

在　`examples/assets-in-binary`　中查看一个完成的示例。

### 使用自定义结构绑定表单数据请求

下面示例使用自定义结构：

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

命令行中使用 `curl`　命令的结果：

```
$ curl "http://localhost:8080/getb?field_a=hello&field_b=world"
{"a":{"FieldA":"hello"},"b":"world"}
$ curl "http://localhost:8080/getc?field_a=hello&field_c=world"
{"a":{"FieldA":"hello"},"c":"world"}
$ curl "http://localhost:8080/getd?field_x=hello&field_d=world"
{"d":"world","x":{"FieldX":"hello"}}
```

**注意**: 不支持下面风格的结构：

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

总之，只支持当前没有　`form` 嵌套的自定义结构。

### 尝试将 body 绑定到不同的结构中

绑定 request body 的常规方法是使用　`c.Request.Body`　并且不能多次调用它们。

```go
type formA struct {
  Foo string `json:"foo" xml:"foo" binding:"required"`
}

type formB struct {
  Bar string `json:"bar" xml:"bar" binding:"required"`
}

func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // 这里 c.ShouldBind 使用 c.Request.Body 并且它不能被重复使用。
  if errA := c.ShouldBind(&objA); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // 这里总会出现一个错误，因为　c.Request.Body　现在是　EOF 。
  } else if errB := c.ShouldBind(&objB); errB == nil {
    c.String(http.StatusOK, `the body should be formB`)
  } else {
    ...
  }
}
```

对于这一点, 你可以使用 `c.ShouldBindBodyWith`　。

```go
func SomeHandler(c *gin.Context) {
  objA := formA{}
  objB := formB{}
  // 这里读取　c.Request.Body　并将结果存储到　context 中。
  if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
    c.String(http.StatusOK, `the body should be formA`)
  // At this time, it reuses body stored in the context.
  // 这是，它重用存储在　context 中的　body 。
  } else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
    c.String(http.StatusOK, `the body should be formB JSON`)
  // 并且它可以接受其他格式
  } else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
    c.String(http.StatusOK, `the body should be formB XML`)
  } else {
    ...
  }
}
```

* `c.ShouldBindBodyWith` 在绑定前存储 body 到 context 中。这对性能会有轻微的影响，所以如果你可以通过立即调用绑定, 不应该使用这个方法。

* 只有一些格式需要这个功能 -- `JSON` 、　`XML`　、 `MsgPack`、
`ProtoBuf`　。 对于其他格式， `Query`、 `Form`、 `FormPost`、 `FormMultipart`，
能被 `c.ShouldBind()`　多次调用，而不会对性能造成任何损害　（参见 [#1341](https://github.com/gin-gonic/gin/pull/1341)）。

### HTTP2 服务器推送

http.Pusher 仅仅被 **go1.8+**　支持。 在 [golang 官方博客](https://blog.golang.org/h2push) 中查看详细信息。

[embedmd]:# (examples/http-pusher/main.go go)
```go
package main

import (
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

var html = template.Must(template.New("https").Parse(`
<html>
<head>
  <title>Https Test</title>
  <script src="/assets/app.js"></script>
</head>
<body>
  <h1 style="color:red;">Welcome, Ginner!</h1>
</body>
</html>
`))

func main() {
	r := gin.Default()
	r.Static("/assets", "./assets")
	r.SetHTMLTemplate(html)

	r.GET("/", func(c *gin.Context) {
		if pusher := c.Writer.Pusher(); pusher != nil {
			// 使用 pusher.Push() 去进行服务器推送
			if err := pusher.Push("/assets/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		c.HTML(200, "https", gin.H{
			"status": "success",
		})
	})

	// 监听并服务于 https://127.0.0.1:8080
	r.RunTLS(":8080", "./testdata/server.pem", "./testdata/server.key")
}
```

## 测试

`net/http/httptest` 包是　HTTP 测试的首选方式。

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

测试上面代码的示例：

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

## 使用者  [![Sourcegraph](https://sourcegraph.com/github.com/gin-gonic/gin/-/badge.svg)](https://sourcegraph.com/github.com/gin-gonic/gin?badge)

使用 [Gin](https://github.com/gin-gonic/gin) web 框架的非常棒的项目列表。

* [drone](https://github.com/drone/drone)：　Drone　是一个用　Go 编写的基于　Docker 的持续交付平台。
* [gorush](https://github.com/appleboy/gorush)：　一个用　Go 编写的消息推送服务器。