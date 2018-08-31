# How to build one effective middleware?

## Consitituent part

The middleware has two parts:

  - part one is what is executed once, when you initalize your middleware. That's where you set up all the global objects, logicals etc. Everything that happens one per application lifetime.

  - part two is what executes on every request. For example, a database middleware you simply inject your "global" database object into the context. Once it's inside the context, you can retrieve it from within other middlewares and your handler furnction.

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

	fmt.Println("end")

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
