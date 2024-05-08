- With issues:
  - Use the search tool before opening a new issue.
  - Please provide source code and commit sha if you found a bug.
  - Review existing issues and provide feedback or react to them.

## Description

<!-- Description of a problem -->

## How to reproduce

<!-- The smallest possible code example to show the problem that can be compiled, like -->
```
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.Default()
	g.GET("/hello/:name", func(c *gin.Context) {
		c.String(200, "Hello %s", c.Param("name"))
	})
	g.Run(":9000")
}
```

## Expectations

<!-- Your expectation result of 'curl' command, like -->
```
$ curl http://localhost:9000/hello/world
Hello world
```

## Actual result

<!-- Actual result showing the problem -->
```
$ curl -i http://localhost:9000/hello/world
<YOUR RESULT>
```

## Environment

- go version:
- gin version (or commit ref):
- operating system:
