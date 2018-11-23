# How to build one effective middleware?

## Consitituent part

The middleware has two parts:

  - part one is what is executed once, when you initialize your middleware. That's where you set up all the global objects, logicals etc. Everything that happens one per application lifetime.

  - part two is what executes on every request. For example, a database middleware you simply inject your "global" database object into the context. Once it's inside the context, you can retrieve it from within other middlewares and your handler function.

```go
func funcName(params string) gin.HandlerFunc {
    // <---
    // This is part one
    // --->
    // The follow code is an example
    if err := check(params); err != nil {
        panic(err)
    }

    return func(c *gin.Context) {
        // <---
        // This is part two
        // --->
        // The follow code is an example
        c.Set("TestVar", params)
        c.Next()    
    }
}
```

## Execution process

Firstly, we have the follow example code:

```go
func main() {
	router := gin.Default()

	router.Use(globalMiddleware())

	router.GET("/rest/n/api/*some", mid1(), mid2(), handler)

	router.Run()
}

func globalMiddleware() gin.HandlerFunc {
	fmt.Println("globalMiddleware...1")

	return func(c *gin.Context) {
		fmt.Println("globalMiddleware...2")
		c.Next()
		fmt.Println("globalMiddleware...3")
	}
}

func handler(c *gin.Context) {
	fmt.Println("exec handler.")
}

func mid1() gin.HandlerFunc {
	fmt.Println("mid1...1")

	return func(c *gin.Context) {

		fmt.Println("mid1...2")
		c.Next()
		fmt.Println("mid1...3")
	}
}

func mid2() gin.HandlerFunc {
	fmt.Println("mid2...1")

	return func(c *gin.Context) {
		fmt.Println("mid2...2")
		c.Next()
		fmt.Println("mid2...3")
	}
}
```

According to [Consitituent part](#consitituent-part) said, when we run the gin process, **part one** will execute firstly and will print the follow information:

```go
globalMiddleware...1
mid1...1
mid2...1
```

And init order are:

```go
globalMiddleware...1
    |
    v
mid1...1
    |
    v
mid2...1
```

When we curl one request `curl -v localhost:8080/rest/n/api/some`, **part two** will execute their middleware and output the following information:

```go
globalMiddleware...2
mid1...2
mid2...2
exec handler.
mid2...3
mid1...3
globalMiddleware...3
```

In other words, run order are:

```go
globalMiddleware...2
    |
    v
mid1...2
    |
    v
mid2...2
    |
    v
exec handler.
    |
    v
mid2...3
    |
    v
mid1...3
    |
    v
globalMiddleware...3
```
