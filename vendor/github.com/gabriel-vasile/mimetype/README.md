<h1 align="center">
  mimetype
</h1>

<h4 align="center">
  A package for detecting MIME types and extensions based on magic numbers
</h4>
<h6 align="center">
  Goroutine safe, extensible, no C bindings
</h6>

<p align="center">
  <a href="https://pkg.go.dev/github.com/gabriel-vasile/mimetype">
    <img alt="Go Reference" src="https://pkg.go.dev/badge/github.com/gabriel-vasile/mimetype.svg">
  </a>
  <a href="https://goreportcard.com/report/github.com/gabriel-vasile/mimetype">
    <img alt="Go report card" src="https://goreportcard.com/badge/github.com/gabriel-vasile/mimetype">
  </a>
  <a href="LICENSE">
    <img alt="License" src="https://img.shields.io/badge/License-MIT-green.svg">
  </a>
</p>

## Features
- fast and precise MIME type and file extension detection
- long list of [supported MIME types](supported_mimes.md)
- possibility to [extend](https://pkg.go.dev/github.com/gabriel-vasile/mimetype#example-package-Extend) with other file formats
- common file formats are prioritized
- [text vs. binary files differentiation](https://pkg.go.dev/github.com/gabriel-vasile/mimetype#example-package-TextVsBinary)
- no external dependencies
- safe for concurrent usage

## Install
```bash
go get github.com/gabriel-vasile/mimetype
```

## Usage
```go
mtype := mimetype.Detect([]byte)
// OR
mtype, err := mimetype.DetectReader(io.Reader)
// OR
mtype, err := mimetype.DetectFile("/path/to/file")
fmt.Println(mtype.String(), mtype.Extension())
```
See the [runnable Go Playground examples](https://pkg.go.dev/github.com/gabriel-vasile/mimetype#pkg-overview).

Caution: only use libraries like **mimetype** as a last resort. Content type detection
using magic numbers is slow, inaccurate, and non-standard. Most of the times
protocols have methods for specifying such metadata; e.g., `Content-Type` header
in HTTP and SMTP.

## FAQ
Q: My file is in the list of [supported MIME types](supported_mimes.md) but
it is not correctly detected. What should I do?

A: Some file formats (often Microsoft Office documents) keep their signatures
towards the end of the file. Try increasing the number of bytes used for detection
with:
```go
mimetype.SetLimit(1024*1024) // Set limit to 1MB.
// or
mimetype.SetLimit(0) // No limit, whole file content used.
mimetype.DetectFile("file.doc")
```
If increasing the limit does not help, please
[open an issue](https://github.com/gabriel-vasile/mimetype/issues/new?assignees=&labels=&template=mismatched-mime-type-detected.md&title=).

## Tests
In addition to unit tests,
[mimetype_tests](https://github.com/gabriel-vasile/mimetype_tests) compares the
library with the [Unix file utility](https://en.wikipedia.org/wiki/File_(command))
for around 50 000 sample files. Check the latest comparison results
[here](https://github.com/gabriel-vasile/mimetype_tests/actions).

## Benchmarks
Benchmarks for each file format are performed when a PR is open. The results can
be seen on the [workflows page](https://github.com/gabriel-vasile/mimetype/actions/workflows/benchmark.yml).
Performance improvements are welcome but correctness is prioritized.

## Structure
**mimetype** uses a hierarchical structure to keep the MIME type detection logic.
This reduces the number of calls needed for detecting the file type. The reason
behind this choice is that there are file formats used as containers for other
file formats. For example, Microsoft Office files are just zip archives,
containing specific metadata files. Once a file has been identified as a
zip, there is no need to check if it is a text file, but it is worth checking if
it is an Microsoft Office file.

To prevent loading entire files into memory, when detecting from a
[reader](https://pkg.go.dev/github.com/gabriel-vasile/mimetype#DetectReader)
or from a [file](https://pkg.go.dev/github.com/gabriel-vasile/mimetype#DetectFile)
**mimetype** limits itself to reading only the header of the input.
<div align="center">
  <img alt="how project is structured" src="https://raw.githubusercontent.com/gabriel-vasile/mimetype/master/testdata/gif.gif" width="88%">
</div>

## Contributing
Contributions are unexpected but welcome. When submitting a PR for detection of
a new file format, please make sure to add a record to the list of testcases
from [mimetype_test.go](mimetype_test.go). For complex files a record can be added
in the [testdata](testdata) directory.
