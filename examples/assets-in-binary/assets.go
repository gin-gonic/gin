package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assetsbfa8d115ce0617d89507412d5393a462f8e9b003 = "<!doctype html>\n<body>\n  <p>Can you see this? → {{.Bar}}</p>\n</body>\n"
var _Assets3737a75b5254ed1f6d588b40a3449721f9ea86c2 = "<!doctype html>\n<body>\n  <p>Hello, {{.Foo}}</p>\n</body>\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": {"html"}, "/html": {"bar.tmpl", "index.tmpl"}}, map[string]*assets.File{
	"/": {
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1524365738, 1524365738517125470),
		Data:     nil,
	}, "/html": {
		Path:     "/html",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1524365491, 1524365491289799093),
		Data:     nil,
	}, "/html/bar.tmpl": {
		Path:     "/html/bar.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1524365491, 1524365491289611557),
		Data:     []byte(_Assetsbfa8d115ce0617d89507412d5393a462f8e9b003),
	}, "/html/index.tmpl": {
		Path:     "/html/index.tmpl",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1524365491, 1524365491289995821),
		Data:     []byte(_Assets3737a75b5254ed1f6d588b40a3449721f9ea86c2),
	}}, "")
