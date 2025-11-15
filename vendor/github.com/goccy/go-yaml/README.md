# YAML support for the Go language

[![PkgGoDev](https://pkg.go.dev/badge/github.com/goccy/go-yaml)](https://pkg.go.dev/github.com/goccy/go-yaml)
![Go](https://github.com/goccy/go-yaml/workflows/Go/badge.svg)
[![codecov](https://codecov.io/gh/goccy/go-yaml/branch/master/graph/badge.svg)](https://codecov.io/gh/goccy/go-yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/goccy/go-yaml)](https://goreportcard.com/report/github.com/goccy/go-yaml)

<img width="300px" src="https://user-images.githubusercontent.com/209884/67159116-64d94b80-f37b-11e9-9b28-f8379636a43c.png"></img>

## This library has **NO** relation to the go-yaml/yaml library

> [!IMPORTANT]
> This library is developed from scratch to replace [`go-yaml/yaml`](https://github.com/go-yaml/yaml).
> If you're looking for a better YAML library, this one should be helpful.

# Why a new library?

As of this writing, there already exists a de facto standard library for YAML processing for Go: [https://github.com/go-yaml/yaml](https://github.com/go-yaml/yaml). However, we believe that a new YAML library is necessary for the following reasons:

- Not actively maintained
- `go-yaml/yaml` has ported the libyaml written in C to Go, so the source code is not written in Go style
- There is a lot of content that cannot be parsed
- YAML is often used for configuration, and it is common to include validation along with it. However, the errors in `go-yaml/yaml` are not intuitive, and it is difficult to provide meaningful validation errors
- When creating tools that use YAML, there are cases where reversible transformation of YAML is required. However, to perform reversible transformations of content that includes Comments or Anchors/Aliases, manipulating the AST is the only option
- Non-intuitive [Marshaler](https://pkg.go.dev/gopkg.in/yaml.v3#Marshaler) / [Unmarshaler](https://pkg.go.dev/gopkg.in/yaml.v3#Unmarshaler)

By the way, libraries such as [ghodss/yaml](https://github.com/ghodss/yaml) and [sigs.k8s.io/yaml](https://github.com/kubernetes-sigs/yaml) also depend on go-yaml/yaml, so if you are using these libraries, the same issues apply: they cannot parse things that go-yaml/yaml cannot parse, and they inherit many of the problems that go-yaml/yaml has.

# Features

- No dependencies
- A better parser than `go-yaml/yaml`. 
  - [Support recursive processing](https://github.com/apple/device-management/blob/release/docs/schema.yaml)
  - Higher coverage in the [YAML Test Suite](https://github.com/yaml/yaml-test-suite?tab=readme-ov-file)
    - YAML Test Suite consists of 402 cases in total, of which `gopkg.in/yaml.v3` passes `295`. In addition to passing all those test cases, `goccy/go-yaml` successfully passes nearly 60 additional test cases ( 2024/12/15 )
    - The test code is [here](https://github.com/goccy/go-yaml/blob/master/yaml_test_suite_test.go#L77)
- Ease and sustainability of maintenance
  - The main maintainer is [@goccy](https://github.com/goccy), but we are also building a system to develop as a team with trusted developers
  - Since it is written from scratch, the code is easy to read for Gophers
- An API structure that allows the use of not only `Encoder`/`Decoder` but also `Tokenizer` and `Parser` functionalities.
  - [lexer.Tokenize](https://pkg.go.dev/github.com/goccy/go-yaml@v1.15.4/lexer#Tokenize)
  - [parser.Parse](https://pkg.go.dev/github.com/goccy/go-yaml@v1.15.4/parser#Parse)
- Filtering, replacing, and merging YAML content using YAML Path
- Reversible transformation without using the AST for YAML that includes Anchors, Aliases, and Comments
- Customize the Marshal/Unmarshal behavior for primitive types and third-party library types ([RegisterCustomMarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#RegisterCustomMarshaler), [RegisterCustomUnmarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#RegisterCustomUnmarshaler))
- Respects `encoding/json` behavior
  - Accept the `json` tag. Note that not all options from the `json` tag will have significance when parsing YAML documents. If both tags exist, `yaml` tag will take precedence.
  - [json.Marshaler](https://pkg.go.dev/encoding/json#Marshaler) style [marshaler](https://pkg.go.dev/github.com/goccy/go-yaml#BytesMarshaler)
  - [json.Unmarshaler](https://pkg.go.dev/encoding/json#Unmarshaler) style [unmarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#BytesUnmarshaler)
  - Options for using `MarshalJSON` and `UnmarshalJSON` ([UseJSONMarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#UseJSONMarshaler), [UseJSONUnmarshaler](https://pkg.go.dev/github.com/goccy/go-yaml#UseJSONUnmarshaler))
- Pretty format for error notifications
- Smart validation processing combined with [go-playground/validator](https://github.com/go-playground/validator)
  - [example test code is here](https://github.com/goccy/go-yaml/blob/45889c98b0a0967240eb595a1bd6896e2f575106/testdata/validate_test.go#L12)
- Allow referencing elements declared in another file via anchors

# Users

The repositories that use goccy/go-yaml are listed here.

- https://github.com/goccy/go-yaml/wiki/Users

The source data is [here](https://github.com/goccy/go-yaml/network/dependents). 
It is already being used in many repositories. Now it's your turn 😄

# Playground

The Playground visualizes how go-yaml processes YAML text. Use it to assist with your debugging or issue reporting.

https://goccy.github.io/go-yaml

# Installation

```sh
go get github.com/goccy/go-yaml
```

# Synopsis

## 1. Simple Encode/Decode

Has an interface like `go-yaml/yaml` using `reflect`

```go
var v struct {
	A int
	B string
}
v.A = 1
v.B = "hello"
bytes, err := yaml.Marshal(v)
if err != nil {
	//...
}
fmt.Println(string(bytes)) // "a: 1\nb: hello\n"
```

```go
	yml := `
%YAML 1.2
---
a: 1
b: c
`
var v struct {
	A int
	B string
}
if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
	//...
}
```

To control marshal/unmarshal behavior, you can use the `yaml` tag.

```go
	yml := `---
foo: 1
bar: c
`
var v struct {
	A int    `yaml:"foo"`
	B string `yaml:"bar"`
}
if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
	//...
}
```

For convenience, we also accept the `json` tag. Note that not all options from
the `json` tag will have significance when parsing YAML documents. If both
tags exist, `yaml` tag will take precedence.

```go
	yml := `---
foo: 1
bar: c
`
var v struct {
	A int    `json:"foo"`
	B string `json:"bar"`
}
if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
	//...
}
```

For custom marshal/unmarshaling, implement either `Bytes` or `Interface` variant of marshaler/unmarshaler. The difference is that while `BytesMarshaler`/`BytesUnmarshaler` behaves like [`encoding/json`](https://pkg.go.dev/encoding/json) and `InterfaceMarshaler`/`InterfaceUnmarshaler` behaves like [`gopkg.in/yaml.v2`](https://pkg.go.dev/gopkg.in/yaml.v2).

Semantically both are the same, but they differ in performance. Because indentation matters in YAML, you cannot simply accept a valid YAML fragment from a Marshaler, and expect it to work when it is attached to the parent container's serialized form. Therefore when we receive use the `BytesMarshaler`, which returns `[]byte`, we must decode it once to figure out how to make it work in the given context. If you use the `InterfaceMarshaler`, we can skip the decoding.

If you are repeatedly marshaling complex objects, the latter is always better
performance wise. But if you are, for example, just providing a choice between
a config file format that is read only once, the former is probably easier to
code.

## 2. Reference elements declared in another file

`testdata` directory contains `anchor.yml` file:

```shell
├── testdata
   └── anchor.yml
```

And `anchor.yml` is defined as follows:

```yaml
a: &a
  b: 1
  c: hello
```

Then, if `yaml.ReferenceDirs("testdata")` option is passed to `yaml.Decoder`, 
 `Decoder` tries to find the anchor definition from YAML files the under `testdata` directory.
 
```go
buf := bytes.NewBufferString("a: *a\n")
dec := yaml.NewDecoder(buf, yaml.ReferenceDirs("testdata"))
var v struct {
	A struct {
		B int
		C string
	}
}
if err := dec.Decode(&v); err != nil {
	//...
}
fmt.Printf("%+v\n", v) // {A:{B:1 C:hello}}
```

## 3. Encode with `Anchor` and `Alias`

### 3.1. Explicitly declared `Anchor` name and `Alias` name

If you want to use `anchor`, you can define it as a struct tag.
If the value specified for an anchor is a pointer type and the same address as the pointer is found, the value is automatically set to alias.
If an explicit alias name is specified, an error is raised if its value is different from the value specified in the anchor.

```go
type T struct {
  A int
  B string
}
var v struct {
  C *T `yaml:"c,anchor=x"`
  D *T `yaml:"d,alias=x"`
}
v.C = &T{A: 1, B: "hello"}
v.D = v.C
bytes, err := yaml.Marshal(v)
if err != nil {
  panic(err)
}
fmt.Println(string(bytes))
/*
c: &x
  a: 1
  b: hello
d: *x
*/
```

### 3.2. Implicitly declared `Anchor` and `Alias` names

If you do not explicitly declare the anchor name, the default behavior is to
use the equivalent of `strings.ToLower($FieldName)` as the name of the anchor.
If the value specified for an anchor is a pointer type and the same address as the pointer is found, the value is automatically set to alias.

```go
type T struct {
	I int
	S string
}
var v struct {
	A *T `yaml:"a,anchor"`
	B *T `yaml:"b,anchor"`
	C *T `yaml:"c"`
	D *T `yaml:"d"`
}
v.A = &T{I: 1, S: "hello"}
v.B = &T{I: 2, S: "world"}
v.C = v.A // C has same pointer address to A
v.D = v.B // D has same pointer address to B
bytes, err := yaml.Marshal(v)
if err != nil {
	//...
}
fmt.Println(string(bytes)) 
/*
a: &a
  i: 1
  s: hello
b: &b
  i: 2
  s: world
c: *a
d: *b
*/
```

### 3.3 MergeKey and Alias

Merge key and alias ( `<<: *alias` ) can be used by embedding a structure with the `inline,alias` tag.

```go
type Person struct {
	*Person `yaml:",omitempty,inline,alias"` // embed Person type for default value
	Name    string `yaml:",omitempty"`
	Age     int    `yaml:",omitempty"`
}
defaultPerson := &Person{
	Name: "John Smith",
	Age:  20,
}
people := []*Person{
	{
		Person: defaultPerson, // assign default value
		Name:   "Ken",         // override Name property
		Age:    10,            // override Age property
	},
	{
		Person: defaultPerson, // assign default value only
	},
}
var doc struct {
	Default *Person   `yaml:"default,anchor"`
	People  []*Person `yaml:"people"`
}
doc.Default = defaultPerson
doc.People = people
bytes, err := yaml.Marshal(doc)
if err != nil {
	//...
}
fmt.Println(string(bytes))
/*
default: &default
  name: John Smith
  age: 20
people:
- <<: *default
  name: Ken
  age: 10
- <<: *default
*/
```

## 4. Pretty Formatted Errors

Error values produced during parsing have two extra features over regular
error values.

First, by default, they contain extra information on the location of the error
from the source YAML document, to make it easier to find the error location.

Second, the error messages can optionally be colorized.

If you would like to control exactly how the output looks like, consider
using  `yaml.FormatError`, which accepts two boolean values to
control turning these features on or off.

<img src="https://user-images.githubusercontent.com/209884/67358124-587f0980-f59a-11e9-96fc-7205aab77695.png"></img>

## 5. Use YAMLPath

```go
yml := `
store:
  book:
    - author: john
      price: 10
    - author: ken
      price: 12
  bicycle:
    color: red
    price: 19.95
`
path, err := yaml.PathString("$.store.book[*].author")
if err != nil {
  //...
}
var authors []string
if err := path.Read(strings.NewReader(yml), &authors); err != nil {
  //...
}
fmt.Println(authors)
// [john ken]
```

### 5.1 Print customized error with YAML source code

```go
package main

import (
  "fmt"

  "github.com/goccy/go-yaml"
)

func main() {
  yml := `
a: 1
b: "hello"
`
  var v struct {
    A int
    B string
  }
  if err := yaml.Unmarshal([]byte(yml), &v); err != nil {
    panic(err)
  }
  if v.A != 2 {
    // output error with YAML source
    path, err := yaml.PathString("$.a")
    if err != nil {
      panic(err)
    }
    source, err := path.AnnotateSource([]byte(yml), true)
    if err != nil {
      panic(err)
    }
    fmt.Printf("a value expected 2 but actual %d:\n%s\n", v.A, string(source))
  }
}
```

output result is the following:

<img src="https://user-images.githubusercontent.com/209884/84148813-7aca8680-aa9a-11ea-8fc9-37dece2ebdac.png"></img>


# Tools

## ycat

print yaml file with color

<img width="713" alt="ycat" src="https://user-images.githubusercontent.com/209884/66986084-19b00600-f0f9-11e9-9f0e-1f91eb072fe0.png">

### Installation

```sh
git clone https://github.com/goccy/go-yaml.git
cd go-yaml/cmd/ycat && go install .
```


# For Developers

> [!NOTE]
> In this project, we manage such test code under the `testdata` directory to avoid adding dependencies  on libraries that are only needed for testing to the top `go.mod` file. Therefore, if you want to add test cases that use 3rd party libraries, please add the test code to the `testdata` directory.

# Looking for Sponsors

I'm looking for sponsors this library. This library is being developed as a personal project in my spare time. If you want a quick response or problem resolution when using this library in your project, please register as a [sponsor](https://github.com/sponsors/goccy). I will cooperate as much as possible. Of course, this library is developed as an MIT license, so you can use it freely for free.

# License

MIT
