# Sonic

English | [中文](README_ZH_CN.md)

A blazingly fast JSON serializing &amp; deserializing library, accelerated by JIT (just-in-time compiling) and SIMD (single-instruction-multiple-data).

## Requirement

- Go: 1.18~1.25
  - Notice: Go1.24.0 is not supported due to the [issue](https://github.com/golang/go/issues/71672), please use higher go version or add build tag `--ldflags="-checklinkname=0"` 
- OS: Linux / MacOS / Windows
- CPU: AMD64 / (ARM64, need go1.20 above)

## Features

- Runtime object binding without code generation
- Complete APIs for JSON value manipulation
- Fast, fast, fast!

## APIs

see [go.dev](https://pkg.go.dev/github.com/bytedance/sonic)

## Benchmarks

For **all sizes** of json and **all scenarios** of usage, **Sonic performs best**.

- [Medium](https://github.com/bytedance/sonic/blob/main/decoder/testdata_test.go#L19) (13KB, 300+ key, 6 layers)

```powershell
goversion: 1.17.1
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkEncoder_Generic_Sonic-16                      32393 ns/op         402.40 MB/s       11965 B/op          4 allocs/op
BenchmarkEncoder_Generic_Sonic_Fast-16                 21668 ns/op         601.57 MB/s       10940 B/op          4 allocs/op
BenchmarkEncoder_Generic_JsonIter-16                   42168 ns/op         309.12 MB/s       14345 B/op        115 allocs/op
BenchmarkEncoder_Generic_GoJson-16                     65189 ns/op         199.96 MB/s       23261 B/op         16 allocs/op
BenchmarkEncoder_Generic_StdLib-16                    106322 ns/op         122.60 MB/s       49136 B/op        789 allocs/op
BenchmarkEncoder_Binding_Sonic-16                       6269 ns/op        2079.26 MB/s       14173 B/op          4 allocs/op
BenchmarkEncoder_Binding_Sonic_Fast-16                  5281 ns/op        2468.16 MB/s       12322 B/op          4 allocs/op
BenchmarkEncoder_Binding_JsonIter-16                   20056 ns/op         649.93 MB/s        9488 B/op          2 allocs/op
BenchmarkEncoder_Binding_GoJson-16                      8311 ns/op        1568.32 MB/s        9481 B/op          1 allocs/op
BenchmarkEncoder_Binding_StdLib-16                     16448 ns/op         792.52 MB/s        9479 B/op          1 allocs/op
BenchmarkEncoder_Parallel_Generic_Sonic-16              6681 ns/op        1950.93 MB/s       12738 B/op          4 allocs/op
BenchmarkEncoder_Parallel_Generic_Sonic_Fast-16         4179 ns/op        3118.99 MB/s       10757 B/op          4 allocs/op
BenchmarkEncoder_Parallel_Generic_JsonIter-16           9861 ns/op        1321.84 MB/s       14362 B/op        115 allocs/op
BenchmarkEncoder_Parallel_Generic_GoJson-16            18850 ns/op         691.52 MB/s       23278 B/op         16 allocs/op
BenchmarkEncoder_Parallel_Generic_StdLib-16            45902 ns/op         283.97 MB/s       49174 B/op        789 allocs/op
BenchmarkEncoder_Parallel_Binding_Sonic-16              1480 ns/op        8810.09 MB/s       13049 B/op          4 allocs/op
BenchmarkEncoder_Parallel_Binding_Sonic_Fast-16         1209 ns/op        10785.23 MB/s      11546 B/op          4 allocs/op
BenchmarkEncoder_Parallel_Binding_JsonIter-16           6170 ns/op        2112.58 MB/s        9504 B/op          2 allocs/op
BenchmarkEncoder_Parallel_Binding_GoJson-16             3321 ns/op        3925.52 MB/s        9496 B/op          1 allocs/op
BenchmarkEncoder_Parallel_Binding_StdLib-16             3739 ns/op        3486.49 MB/s        9480 B/op          1 allocs/op

BenchmarkDecoder_Generic_Sonic-16                      66812 ns/op         195.10 MB/s       57602 B/op        723 allocs/op
BenchmarkDecoder_Generic_Sonic_Fast-16                 54523 ns/op         239.07 MB/s       49786 B/op        313 allocs/op
BenchmarkDecoder_Generic_StdLib-16                    124260 ns/op         104.90 MB/s       50869 B/op        772 allocs/op
BenchmarkDecoder_Generic_JsonIter-16                   91274 ns/op         142.81 MB/s       55782 B/op       1068 allocs/op
BenchmarkDecoder_Generic_GoJson-16                     88569 ns/op         147.17 MB/s       66367 B/op        973 allocs/op
BenchmarkDecoder_Binding_Sonic-16                      32557 ns/op         400.38 MB/s       28302 B/op        137 allocs/op
BenchmarkDecoder_Binding_Sonic_Fast-16                 28649 ns/op         455.00 MB/s       24999 B/op         34 allocs/op
BenchmarkDecoder_Binding_StdLib-16                    111437 ns/op         116.97 MB/s       10576 B/op        208 allocs/op
BenchmarkDecoder_Binding_JsonIter-16                   35090 ns/op         371.48 MB/s       14673 B/op        385 allocs/op
BenchmarkDecoder_Binding_GoJson-16                     28738 ns/op         453.59 MB/s       22039 B/op         49 allocs/op
BenchmarkDecoder_Parallel_Generic_Sonic-16             12321 ns/op        1057.91 MB/s       57233 B/op        723 allocs/op
BenchmarkDecoder_Parallel_Generic_Sonic_Fast-16        10644 ns/op        1224.64 MB/s       49362 B/op        313 allocs/op
BenchmarkDecoder_Parallel_Generic_StdLib-16            57587 ns/op         226.35 MB/s       50874 B/op        772 allocs/op
BenchmarkDecoder_Parallel_Generic_JsonIter-16          38666 ns/op         337.12 MB/s       55789 B/op       1068 allocs/op
BenchmarkDecoder_Parallel_Generic_GoJson-16            30259 ns/op         430.79 MB/s       66370 B/op        974 allocs/op
BenchmarkDecoder_Parallel_Binding_Sonic-16              5965 ns/op        2185.28 MB/s       27747 B/op        137 allocs/op
BenchmarkDecoder_Parallel_Binding_Sonic_Fast-16         5170 ns/op        2521.31 MB/s       24715 B/op         34 allocs/op
BenchmarkDecoder_Parallel_Binding_StdLib-16            27582 ns/op         472.58 MB/s       10576 B/op        208 allocs/op
BenchmarkDecoder_Parallel_Binding_JsonIter-16          13571 ns/op         960.51 MB/s       14685 B/op        385 allocs/op
BenchmarkDecoder_Parallel_Binding_GoJson-16            10031 ns/op        1299.51 MB/s       22111 B/op         49 allocs/op

BenchmarkGetOne_Sonic-16                                3276 ns/op        3975.78 MB/s          24 B/op          1 allocs/op
BenchmarkGetOne_Gjson-16                                9431 ns/op        1380.81 MB/s           0 B/op          0 allocs/op
BenchmarkGetOne_Jsoniter-16                            51178 ns/op         254.46 MB/s       27936 B/op        647 allocs/op
BenchmarkGetOne_Parallel_Sonic-16                      216.7 ns/op       60098.95 MB/s          24 B/op          1 allocs/op
BenchmarkGetOne_Parallel_Gjson-16                       1076 ns/op        12098.62 MB/s          0 B/op          0 allocs/op
BenchmarkGetOne_Parallel_Jsoniter-16                   17741 ns/op         734.06 MB/s       27945 B/op        647 allocs/op
BenchmarkSetOne_Sonic-16                               9571 ns/op         1360.61 MB/s        1584 B/op         17 allocs/op
BenchmarkSetOne_Sjson-16                               36456 ns/op         357.22 MB/s       52180 B/op          9 allocs/op
BenchmarkSetOne_Jsoniter-16                            79475 ns/op         163.86 MB/s       45862 B/op        964 allocs/op
BenchmarkSetOne_Parallel_Sonic-16                      850.9 ns/op       15305.31 MB/s        1584 B/op         17 allocs/op
BenchmarkSetOne_Parallel_Sjson-16                      18194 ns/op         715.77 MB/s       52247 B/op          9 allocs/op
BenchmarkSetOne_Parallel_Jsoniter-16                   33560 ns/op         388.05 MB/s       45892 B/op        964 allocs/op
BenchmarkLoadNode/LoadAll()-16                         11384 ns/op        1143.93 MB/s        6307 B/op         25 allocs/op
BenchmarkLoadNode_Parallel/LoadAll()-16                 5493 ns/op        2370.68 MB/s        7145 B/op         25 allocs/op
BenchmarkLoadNode/Interface()-16                       17722 ns/op         734.85 MB/s       13323 B/op         88 allocs/op
BenchmarkLoadNode_Parallel/Interface()-16              10330 ns/op        1260.70 MB/s       15178 B/op         88 allocs/op
```

- [Small](https://github.com/bytedance/sonic/blob/main/testdata/small.go) (400B, 11 keys, 3 layers)
![small benchmarks](./docs/imgs/bench-small.png)
- [Large](https://github.com/bytedance/sonic/blob/main/testdata/twitter.json) (635KB, 10000+ key, 6 layers)
![large benchmarks](./docs/imgs/bench-large.png)

See [bench.sh](https://github.com/bytedance/sonic/blob/main/scripts/bench.sh) for benchmark codes.

## How it works

See [INTRODUCTION.md](./docs/INTRODUCTION.md).

## Usage

### Marshal/Unmarshal

Default behaviors are mostly consistent with `encoding/json`, except HTML escaping form (see [Escape HTML](https://github.com/bytedance/sonic/blob/main/README.md#escape-html)) and `SortKeys` feature (optional support see [Sort Keys](https://github.com/bytedance/sonic/blob/main/README.md#sort-keys)) that is **NOT** in conformity to [RFC8259](https://datatracker.ietf.org/doc/html/rfc8259).

 ```go
import "github.com/bytedance/sonic"

var data YourSchema
// Marshal
output, err := sonic.Marshal(&data)
// Unmarshal
err := sonic.Unmarshal(output, &data)
 ```

### Streaming IO

Sonic supports decoding json from `io.Reader` or encoding objects into `io.Writer`, aims at handling multiple values as well as reducing memory consumption.

- encoder

```go
var o1 = map[string]interface{}{
    "a": "b",
}
var o2 = 1
var w = bytes.NewBuffer(nil)
var enc = sonic.ConfigDefault.NewEncoder(w)
enc.Encode(o1)
enc.Encode(o2)
fmt.Println(w.String())
// Output:
// {"a":"b"}
// 1
```

- decoder

```go
var o =  map[string]interface{}{}
var r = strings.NewReader(`{"a":"b"}{"1":"2"}`)
var dec = sonic.ConfigDefault.NewDecoder(r)
dec.Decode(&o)
dec.Decode(&o)
fmt.Printf("%+v", o)
// Output:
// map[1:2 a:b]
```

### Use Number/Use Int64

 ```go
import "github.com/bytedance/sonic/decoder"

var input = `1`
var data interface{}

// default float64
dc := decoder.NewDecoder(input)
dc.Decode(&data) // data == float64(1)
// use json.Number
dc = decoder.NewDecoder(input)
dc.UseNumber()
dc.Decode(&data) // data == json.Number("1")
// use int64
dc = decoder.NewDecoder(input)
dc.UseInt64()
dc.Decode(&data) // data == int64(1)

root, err := sonic.GetFromString(input)
// Get json.Number
jn := root.Number()
jm := root.InterfaceUseNumber().(json.Number) // jn == jm
// Get float64
fn := root.Float64()
fm := root.Interface().(float64) // jn == jm
 ```

### Sort Keys

On account of the performance loss from sorting (roughly 10%), sonic doesn't enable this feature by default. If your component depends on it to work (like [zstd](https://github.com/facebook/zstd)), Use it like this:

```go
import "github.com/bytedance/sonic"
import "github.com/bytedance/sonic/encoder"

// Binding map only
m := map[string]interface{}{}
v, err := encoder.Encode(m, encoder.SortMapKeys)

// Or ast.Node.SortKeys() before marshal
var root := sonic.Get(JSON)
err := root.SortKeys()
```

### Escape HTML

On account of the performance loss (roughly 15%), sonic doesn't enable this feature by default. You can use `encoder.EscapeHTML` option to open this feature (align with `encoding/json.HTMLEscape`).

```go
import "github.com/bytedance/sonic"

v := map[string]string{"&&":"<>"}
ret, err := Encode(v, EscapeHTML) // ret == `{"\u0026\u0026":{"X":"\u003c\u003e"}}`
```

### Compact Format

Sonic encodes primitive objects (struct/map...) as compact-format JSON by default, except marshaling `json.RawMessage` or `json.Marshaler`: sonic ensures validating their output JSON but **DO NOT** compacting them for performance concerns. We provide the option `encoder.CompactMarshaler` to add compacting process.

### Print Error

If there invalid syntax in input JSON, sonic will return `decoder.SyntaxError`, which supports pretty-printing of error position

```go
import "github.com/bytedance/sonic"
import "github.com/bytedance/sonic/decoder"

var data interface{}
err := sonic.UnmarshalString("[[[}]]", &data)
if err != nil {
    /* One line by default */
    println(e.Error()) // "Syntax error at index 3: invalid char\n\n\t[[[}]]\n\t...^..\n"
    /* Pretty print */
    if e, ok := err.(decoder.SyntaxError); ok {
        /*Syntax error at index 3: invalid char

            [[[}]]
            ...^..
        */
        print(e.Description())
    } else if me, ok := err.(*decoder.MismatchTypeError); ok {
        // decoder.MismatchTypeError is new to Sonic v1.6.0
        print(me.Description())
    }
}
```

#### Mismatched Types [Sonic v1.6.0]

If there a **mismatch-typed** value for a given key, sonic will report `decoder.MismatchTypeError` (if there are many, report the last one), but still skip wrong the value and keep decoding next JSON.

```go
import "github.com/bytedance/sonic"
import "github.com/bytedance/sonic/decoder"

var data = struct{
    A int
    B int
}{}
err := UnmarshalString(`{"A":"1","B":1}`, &data)
println(err.Error())    // Mismatch type int with value string "at index 5: mismatched type with value\n\n\t{\"A\":\"1\",\"B\":1}\n\t.....^.........\n"
fmt.Printf("%+v", data) // {A:0 B:1}
```

### Ast.Node

Sonic/ast.Node is a completely self-contained AST for JSON. It implements serialization and deserialization both and provides robust APIs for obtaining and modification of generic data.

#### Get/Index

Search partial JSON by given paths, which must be non-negative integer or string, or nil

```go
import "github.com/bytedance/sonic"

input := []byte(`{"key1":[{},{"key2":{"key3":[1,2,3]}}]}`)

// no path, returns entire json
root, err := sonic.Get(input)
raw := root.Raw() // == string(input)

// multiple paths
root, err := sonic.Get(input, "key1", 1, "key2")
sub := root.Get("key3").Index(2).Int64() // == 3
```

**Tip**: since `Index()` uses offset to locate data, which is much faster than scanning like `Get()`, we suggest you use it as much as possible. And sonic also provides another API `IndexOrGet()` to underlying use offset as well as ensure the key is matched.

#### SearchOption

`Searcher` provides some options for user to meet different needs:

```go
opts := ast.SearchOption{ CopyReturn: true ... }
val, err := sonic.GetWithOptions(JSON, opts, "key")
```

- CopyReturn
Indicate the searcher to copy the result JSON string instead of refer from the input. This can help to reduce memory usage if you cache the results
- ConcurentRead
Since `ast.Node` use `Lazy-Load` design, it doesn't support Concurrently-Read by default. If you want to read it concurrently, please specify it.
- ValidateJSON
Indicate the searcher to validate the entire JSON. This option is enabled by default, which slow down the search speed a little.

#### Set/Unset

Modify the json content by Set()/Unset()

```go
import "github.com/bytedance/sonic"

// Set
exist, err := root.Set("key4", NewBool(true)) // exist == false
alias1 := root.Get("key4")
println(alias1.Valid()) // true
alias2 := root.Index(1)
println(alias1 == alias2) // true

// Unset
exist, err := root.UnsetByIndex(1) // exist == true
println(root.Get("key4").Check()) // "value not exist"
```

#### Serialize

To encode `ast.Node` as json, use `MarshalJson()` or `json.Marshal()` (MUST pass the node's pointer)

```go
import (
    "encoding/json"
    "github.com/bytedance/sonic"
)

buf, err := root.MarshalJson()
println(string(buf))                // {"key1":[{},{"key2":{"key3":[1,2,3]}}]}
exp, err := json.Marshal(&root)     // WARN: use pointer
println(string(buf) == string(exp)) // true
```

#### APIs

- validation: `Check()`, `Error()`, `Valid()`, `Exist()`
- searching: `Index()`, `Get()`, `IndexPair()`, `IndexOrGet()`, `GetByPath()`
- go-type casting: `Int64()`, `Float64()`, `String()`, `Number()`, `Bool()`, `Map[UseNumber|UseNode]()`, `Array[UseNumber|UseNode]()`, `Interface[UseNumber|UseNode]()`
- go-type packing: `NewRaw()`, `NewNumber()`, `NewNull()`, `NewBool()`, `NewString()`, `NewObject()`, `NewArray()`
- iteration: `Values()`, `Properties()`, `ForEach()`, `SortKeys()`
- modification: `Set()`, `SetByIndex()`, `Add()`

### Ast.Visitor

Sonic provides an advanced API for fully parsing JSON into non-standard types (neither `struct` not `map[string]interface{}`) without using any intermediate representation (`ast.Node` or `interface{}`). For example, you might have the following types which are like `interface{}` but actually not `interface{}`:

```go
type UserNode interface {}

// the following types implement the UserNode interface.
type (
    UserNull    struct{}
    UserBool    struct{ Value bool }
    UserInt64   struct{ Value int64 }
    UserFloat64 struct{ Value float64 }
    UserString  struct{ Value string }
    UserObject  struct{ Value map[string]UserNode }
    UserArray   struct{ Value []UserNode }
)
```

Sonic provides the following API to return **the preorder traversal of a JSON AST**. The `ast.Visitor` is a SAX style interface which is used in some C++ JSON library. You should implement `ast.Visitor` by yourself and pass it to `ast.Preorder()` method. In your visitor you can make your custom types to represent JSON values. There may be an O(n) space container (such as stack) in your visitor to record the object / array hierarchy.

```go
func Preorder(str string, visitor Visitor, opts *VisitorOptions) error

type Visitor interface {
    OnNull() error
    OnBool(v bool) error
    OnString(v string) error
    OnInt64(v int64, n json.Number) error
    OnFloat64(v float64, n json.Number) error
    OnObjectBegin(capacity int) error
    OnObjectKey(key string) error
    OnObjectEnd() error
    OnArrayBegin(capacity int) error
    OnArrayEnd() error
}
```

See [ast/visitor.go](https://github.com/bytedance/sonic/blob/main/ast/visitor.go) for detailed usage. We also implement a demo visitor for `UserNode` in [ast/visitor_test.go](https://github.com/bytedance/sonic/blob/main/ast/visitor_test.go).

## Compatibility

For developers who want to use sonic to meet different scenarios, we provide some integrated configs as `sonic.API`

- `ConfigDefault`: the sonic's default config (`EscapeHTML=false`,`SortKeys=false`...) to run sonic fast meanwhile ensure security.
- `ConfigStd`: the std-compatible config (`EscapeHTML=true`,`SortKeys=true`...)
- `ConfigFastest`: the fastest config (`NoQuoteTextMarshaler=true`) to run on sonic as fast as possible.
Sonic **DOES NOT** ensure to support all environments, due to the difficulty of developing high-performance codes. On non-sonic-supporting environment, the implementation will fall back to `encoding/json`. Thus below configs will all equal to `ConfigStd`.

## Tips

### Pretouch

Since Sonic uses [golang-asm](https://github.com/twitchyliquid64/golang-asm) as a JIT assembler, which is NOT very suitable for runtime compiling, first-hit running of a huge schema may cause request-timeout or even process-OOM. For better stability, we advise **using `Pretouch()` for huge-schema or compact-memory applications** before `Marshal()/Unmarshal()`.

```go
import (
    "reflect"
    "github.com/bytedance/sonic"
    "github.com/bytedance/sonic/option"
)

func init() {
    var v HugeStruct

    // For most large types (nesting depth <= option.DefaultMaxInlineDepth)
    err := sonic.Pretouch(reflect.TypeOf(v))

    // with more CompileOption...
    err := sonic.Pretouch(reflect.TypeOf(v),
        // If the type is too deep nesting (nesting depth > option.DefaultMaxInlineDepth),
        // you can set compile recursive loops in Pretouch for better stability in JIT.
        option.WithCompileRecursiveDepth(loop),
        // For a large nested struct, try to set a smaller depth to reduce compiling time.
        option.WithCompileMaxInlineDepth(depth),
    )
}
```

### Copy string

When decoding **string values without any escaped characters**, sonic references them from the origin JSON buffer instead of mallocing a new buffer to copy. This helps a lot for CPU performance but may leave the whole JSON buffer in memory as long as the decoded objects are being used. In practice, we found the extra memory introduced by referring JSON buffer is usually 20% ~ 80% of decoded objects. Once an application holds these objects for a long time (for example, cache the decoded objects for reusing), its in-use memory on the server may go up. - `Config.CopyString`/`decoder.CopyString()`: We provide the option for `Decode()` / `Unmarshal()` users to choose not to reference the JSON buffer, which may cause a decline in CPU performance to some degree.

- `GetFromStringNoCopy()`: For memory safety, `sonic.Get()` / `sonic.GetFromString()` now copies return JSON. If users want to get json more quickly and not care about memory usage, you can use `GetFromStringNoCopy()` to return a JSON directly referenced from source.

### Pass string or []byte?

For alignment to `encoding/json`, we provide API to pass `[]byte` as an argument, but the string-to-bytes copy is conducted at the same time considering safety, which may lose performance when the origin JSON is huge. Therefore, you can use `UnmarshalString()` and `GetFromString()` to pass a string, as long as your origin data is a string or **nocopy-cast** is safe for your []byte. We also provide API `MarshalString()` for convenient **nocopy-cast** of encoded JSON []byte, which is safe since sonic's output bytes is always duplicated and unique.

### Accelerate `encoding.TextMarshaler`

To ensure data security, sonic.Encoder quotes and escapes string values from `encoding.TextMarshaler` interfaces by default, which may degrade performance much if most of your data is in form of them. We provide `encoder.NoQuoteTextMarshaler` to skip these operations, which means you **MUST** ensure their output string escaped and quoted following [RFC8259](https://datatracker.ietf.org/doc/html/rfc8259).

### Better performance for generic data

In **fully-parsed** scenario, `Unmarshal()` performs better than `Get()`+`Node.Interface()`. But if you only have a part of the schema for specific json, you can combine `Get()` and `Unmarshal()` together:

```go
import "github.com/bytedance/sonic"

node, err := sonic.GetFromString(_TwitterJson, "statuses", 3, "user")
var user User // your partial schema...
err = sonic.UnmarshalString(node.Raw(), &user)
```

Even if you don't have any schema, use `ast.Node` as the container of generic values instead of `map` or `interface`:

```go
import "github.com/bytedance/sonic"

root, err := sonic.GetFromString(_TwitterJson)
user := root.GetByPath("statuses", 3, "user")  // === root.Get("status").Index(3).Get("user")
err = user.Check()

// err = user.LoadAll() // only call this when you want to use 'user' concurrently...
go someFunc(user)
```

Why? Because `ast.Node` stores its children using `array`:

- `Array`'s performance is **much better** than `Map` when Inserting (Deserialize) and Scanning (Serialize) data;
- **Hashing** (`map[x]`) is not as efficient as **Indexing** (`array[x]`), which `ast.Node` can conduct on **both array and object**;
- Using `Interface()`/`Map()` means Sonic must parse all the underlying values, while `ast.Node` can parse them **on demand**.

**CAUTION:** `ast.Node` **DOESN'T** ensure concurrent security directly, due to its **lazy-load** design. However, you can call `Node.Load()`/`Node.LoadAll()` to achieve that, which may bring performance reduction while it still works faster than converting to `map` or `interface{}`

### Ast.Node or Ast.Visitor?

For generic data, `ast.Node` should be enough for your needs in most cases.

However, `ast.Node` is designed for partially processing JSON string. It has some special designs such as lazy-load which might not be suitable for directly parsing the whole JSON string like `Unmarshal()`. Although `ast.Node` is better then `map` or `interface{}`, it's also a kind of intermediate representation after all if your final types are customized and you have to convert the above types to your custom types after parsing.

For better performance, in previous case the `ast.Visitor` will be the better choice. It performs JSON decoding like `Unmarshal()` and you can directly use your final types to represents a JSON AST without any intermediate representations.

But `ast.Visitor` is not a very handy API. You might need to write a lot of code to implement your visitor and carefully maintain the tree hierarchy during decoding. Please read the comments in [ast/visitor.go](https://github.com/bytedance/sonic/blob/main/ast/visitor.go) carefully if you decide to use this API.

### Buffer Size

Sonic use memory pool in many places like `encoder.Encode`, `ast.Node.MarshalJSON` to improve performance, which may produce more memory usage (in-use) when server's load is high. See [issue 614](https://github.com/bytedance/sonic/issues/614). Therefore, we introduce some options to let user control the behavior of memory pool. See [option](https://pkg.go.dev/github.com/bytedance/sonic@v1.11.9/option#pkg-variables) package.

### Faster JSON Skip

For security, sonic use [FSM](native/skip_one.c) algorithm to  validate JSON when decoding raw JSON or encoding `json.Marshaler`, which is much slower (1~10x) than [SIMD-searching-pair](native/skip_one_fast.c) algorithm. If user has many redundant JSON value and DO NOT NEED to strictly validate JSON correctness, you can enable below options:

- `Config.NoValidateSkipJSON`: for faster skipping JSON when decoding, such as unknown fields, json.Unmarshaler(json.RawMessage), mismatched values, and redundant array elements
- `Config.NoValidateJSONMarshaler`: avoid validating JSON when encoding `json.Marshaler`
- `SearchOption.ValidateJSON`: indicates if validate located JSON value when `Get`

## JSON-Path Support (GJSON)

[tidwall/gjson](https://github.com/tidwall/gjson) has provided a comprehensive and popular JSON-Path API, and
 a lot of older codes heavily relies on it. Therefore, we provides a wrapper library, which combines gjson's API with sonic's SIMD algorithm to boost up the performance. See [cloudwego/gjson](https://github.com/cloudwego/gjson).

## Community

Sonic is a subproject of [CloudWeGo](https://www.cloudwego.io/). We are committed to building a cloud native ecosystem.
