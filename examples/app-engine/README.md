# Guide to run Gin under App Engine LOCAL Development Server

1. Download, install and setup Go in your computer. (That includes setting your `$GOPATH`.)
2. Download SDK for your platform from [here](https://cloud.google.com/appengine/docs/standard/go/download): `https://cloud.google.com/appengine/docs/standard/go/download`
3. Download Gin source code using: `$ go get github.com/gin-gonic/gin`
4. Navigate to examples folder: `$ cd $GOPATH/src/github.com/gin-gonic/gin/examples/app-engine/`
5. Run it: `$ dev_appserver.py .` (notice that you have to run this script by Python2)

