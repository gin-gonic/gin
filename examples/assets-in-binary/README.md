# Building a single binary containing templates

This is a complete example to create a single binary with the
[gin-gonic/gin][gin] Web Server with HTML templates.

[gin]: https://github.com/gin-gonic/gin

## How to use

### Prepare Packages

```
go get github.com/gin-gonic/gin
go get github.com/jessevdk/go-assets-builder
```

### Generate assets.go

```
go-assets-builder html -o assets.go
```

### Build the server

```
go build -o assets-in-binary
```

### Run

```
./assets-in-binary
```
