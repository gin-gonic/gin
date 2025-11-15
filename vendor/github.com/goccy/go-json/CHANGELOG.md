# v0.10.2 - 2023/03/20

### New features

* Support DebugDOT option for debugging encoder ( #440 )

### Fix bugs

* Fix combination of embedding structure and omitempty option ( #442 )

# v0.10.1 - 2023/03/13

### Fix bugs

* Fix checkptr error for array decoder ( #415 )
* Fix added buffer size check when decoding key ( #430 )
* Fix handling of anonymous fields other than struct ( #431 )
* Fix to not optimize when lower conversion can't handle byte-by-byte ( #432 )
* Fix a problem that MarshalIndent does not work when UnorderedMap is specified ( #435 )
* Fix mapDecoder.DecodeStream() for empty objects containing whitespace ( #425 )
* Fix an issue that could not set the correct NextField for fields in the embedded structure ( #438 )

# v0.10.0 - 2022/11/29

### New features

* Support JSON Path ( #250 )

### Fix bugs

* Fix marshaler for map's key ( #409 )

# v0.9.11 - 2022/08/18

### Fix bugs

* Fix unexpected behavior when buffer ends with backslash ( #383 )
* Fix stream decoding of escaped character ( #387 )

# v0.9.10 - 2022/07/15

### Fix bugs

* Fix boundary exception of type caching ( #382 )

# v0.9.9 - 2022/07/15

### Fix bugs

* Fix encoding of directed interface with typed nil ( #377 )
* Fix embedded primitive type encoding using alias ( #378 )
* Fix slice/array type encoding with types implementing MarshalJSON ( #379 )
* Fix unicode decoding when the expected buffer state is not met after reading ( #380 )

# v0.9.8 - 2022/06/30

### Fix bugs

* Fix decoding of surrogate-pair ( #365 )
* Fix handling of embedded primitive type ( #366 )
* Add validation of escape sequence for decoder ( #367 )
* Fix stream tokenizing respecting UseNumber ( #369 )
* Fix encoding when struct pointer type that implements Marshal JSON is embedded ( #375 )

### Improve performance

* Improve performance of linkRecursiveCode ( #368 )

# v0.9.7 - 2022/04/22

### Fix bugs

#### Encoder

* Add filtering process for encoding on slow path ( #355 )
* Fix encoding of interface{} with pointer type ( #363 )

#### Decoder

* Fix map key decoder that implements UnmarshalJSON ( #353 )
* Fix decoding of []uint8 type ( #361 )

### New features

* Add DebugWith option for encoder ( #356 )

# v0.9.6 - 2022/03/22

### Fix bugs

* Correct the handling of the minimum value of int type for decoder ( #344 )
* Fix bugs of stream decoder's bufferSize ( #349 )
* Add a guard to use typeptr more safely ( #351 )

### Improve decoder performance

* Improve escapeString's performance ( #345 )

### Others

* Update go version for CI ( #347 )

# v0.9.5 - 2022/03/04

### Fix bugs

* Fix panic when decoding time.Time with context ( #328 )
* Fix reading the next character in buffer to nul consideration ( #338 )
* Fix incorrect handling on skipValue ( #341 )

### Improve decoder performance

* Improve performance when a payload contains escape sequence ( #334 )

# v0.9.4 - 2022/01/21

* Fix IsNilForMarshaler for string type with omitempty ( #323 )
* Fix the case where the embedded field is at the end ( #326 )

# v0.9.3 - 2022/01/14

* Fix logic of removing struct field for decoder ( #322 )

# v0.9.2 - 2022/01/14

* Add invalid decoder to delay type error judgment at decode ( #321 )

# v0.9.1 - 2022/01/11

* Fix encoding of MarshalText/MarshalJSON operation with head offset ( #319 )

# v0.9.0 - 2022/01/05

### New feature

* Supports dynamic filtering of struct fields ( #314 )

### Improve encoding performance

* Improve map encoding performance ( #310 )
* Optimize encoding path for escaped string ( #311 )
* Add encoding option for performance ( #312 )

### Fix bugs

* Fix panic at encoding map value on 1.18 ( #310 )
* Fix MarshalIndent for interface type ( #317 )

# v0.8.1 - 2021/12/05

* Fix operation conversion from PtrHead to Head in Recursive type ( #305 )

# v0.8.0 - 2021/12/02

* Fix embedded field conflict behavior ( #300 )
* Refactor compiler for encoder ( #301 #302 )

# v0.7.10 - 2021/10/16

* Fix conversion from pointer to uint64  ( #294 )

# v0.7.9 - 2021/09/28

* Fix encoding of nil value about interface type that has method ( #291 )

# v0.7.8 - 2021/09/01

* Fix mapassign_faststr for indirect struct type ( #283 )
* Fix encoding of not empty interface type ( #284 )
* Fix encoding of empty struct interface type ( #286 )

# v0.7.7 - 2021/08/25

* Fix invalid utf8 on stream decoder ( #279 )
* Fix buffer length bug on string stream decoder ( #280 )

Thank you @orisano !!

# v0.7.6 - 2021/08/13

* Fix nil slice assignment ( #276 )
* Improve error message ( #277 )

# v0.7.5 - 2021/08/12

* Fix encoding of embedded struct with tags ( #265 )
* Fix encoding of embedded struct that isn't first field ( #272 )
* Fix decoding of binary type with escaped char ( #273 )

# v0.7.4 - 2021/07/06

* Fix encoding of indirect layout structure ( #264 )

# v0.7.3 - 2021/06/29

* Fix encoding of pointer type in empty interface ( #262 )

# v0.7.2 - 2021/06/26

### Fix decoder

* Add decoder for func type to fix decoding of nil function value ( #257 )
* Fix stream decoding of []byte type ( #258 )

### Performance

* Improve decoding performance of map[string]interface{} type ( use `mapassign_faststr` ) ( #256 )
* Improve encoding performance of empty interface type ( remove recursive calling of `vm.Run` ) ( #259 )

### Benchmark

* Add bytedance/sonic as benchmark target ( #254 )

# v0.7.1 - 2021/06/18

### Fix decoder

* Fix error when unmarshal empty array ( #253 )

# v0.7.0 - 2021/06/12

### Support context for MarshalJSON and UnmarshalJSON ( #248 )

* json.MarshalContext(context.Context, interface{}, ...json.EncodeOption) ([]byte, error)
* json.NewEncoder(io.Writer).EncodeContext(context.Context, interface{}, ...json.EncodeOption) error
* json.UnmarshalContext(context.Context, []byte, interface{}, ...json.DecodeOption) error
* json.NewDecoder(io.Reader).DecodeContext(context.Context, interface{}) error

```go
type MarshalerContext interface {
  MarshalJSON(context.Context) ([]byte, error)
}

type UnmarshalerContext interface {
  UnmarshalJSON(context.Context, []byte) error
}
```

### Add DecodeFieldPriorityFirstWin option ( #242 )

In the default behavior, go-json, like encoding/json, will reflect the result of the last evaluation when a field with the same name exists. I've added new options to allow you to change this behavior. `json.DecodeFieldPriorityFirstWin` option reflects the result of the first evaluation if a field with the same name exists. This behavior has a performance advantage as it allows the subsequent strings to be skipped if all fields have been evaluated.

### Fix encoder

* Fix indent number contains recursive type ( #249 )
* Fix encoding of using empty interface as map key ( #244 )

### Fix decoder

* Fix decoding fields containing escaped characters ( #237 )

### Refactor

* Move some tests to subdirectory ( #243 )
* Refactor package layout for decoder ( #238 )

# v0.6.1 - 2021/06/02

### Fix encoder

* Fix value of totalLength for encoding ( #236 )

# v0.6.0 - 2021/06/01

### Support Colorize option for encoding (#233)

```go
b, err := json.MarshalWithOption(v, json.Colorize(json.DefaultColorScheme))
if err != nil {
  ...
}
fmt.Println(string(b)) // print colored json
```

### Refactor

* Fix opcode layout - Adjust memory layout of the opcode to 128 bytes in a 64-bit environment ( #230 )
* Refactor encode option ( #231 )
* Refactor escape string ( #232 )

# v0.5.1 - 2021/5/20

### Optimization

* Add type addrShift to enable bigger encoder/decoder cache ( #213 )

### Fix decoder

* Keep original reference of slice element ( #229 )

### Refactor

* Refactor Debug mode for encoding ( #226 )
* Generate VM sources for encoding ( #227 )
* Refactor validator for null/true/false for decoding ( #221 )

# v0.5.0 - 2021/5/9

### Supports using omitempty and string tags at the same time ( #216 )

### Fix decoder

* Fix stream decoder for unicode char ( #215 )
* Fix decoding of slice element ( #219 )
* Fix calculating of buffer length for stream decoder ( #220 )

### Refactor

* replace skipWhiteSpace goto by loop ( #212 )

# v0.4.14 - 2021/5/4

### Benchmark

* Add valyala/fastjson to benchmark ( #193 )
* Add benchmark task for CI ( #211 )

### Fix decoder

* Fix decoding of slice with unmarshal json type ( #198 )
* Fix decoding of null value for interface type that does not implement Unmarshaler ( #205 )
* Fix decoding of null value to []byte by json.Unmarshal ( #206 )
* Fix decoding of backslash char at the end of string ( #207 )
* Fix stream decoder for null/true/false value ( #208 )
* Fix stream decoder for slow reader ( #211 )

### Performance

* If cap of slice is enough, reuse slice data for compatibility with encoding/json ( #200 )

# v0.4.13 - 2021/4/20

### Fix json.Compact and json.Indent

* Support validation the input buffer for json.Compact and json.Indent ( #189 )
* Optimize json.Compact and json.Indent ( improve memory footprint ) ( #190 )

# v0.4.12 - 2021/4/15

### Fix encoder

* Fix unnecessary indent for empty slice type ( #181 )
* Fix encoding of omitempty feature for the slice or interface type ( #183 )
* Fix encoding custom types zero values with omitempty when marshaller exists ( #187 )

### Fix decoder

* Fix decoder for invalid top level value ( #184 )
* Fix decoder for invalid number value ( #185 )

# v0.4.11 - 2021/4/3

* Improve decoder performance for interface type

# v0.4.10 - 2021/4/2

### Fix encoder

* Fixed a bug when encoding slice and map containing recursive structures
* Fixed a logic to determine if indirect reference

# v0.4.9 - 2021/3/29

### Add debug mode

If you use `json.MarshalWithOption(v, json.Debug())` and `panic` occurred in `go-json`, produces debug information to console.

### Support a new feature to compatible with encoding/json

- invalid UTF-8 is coerced to valid UTF-8 ( without performance down )

### Fix encoder

- Fixed handling of MarshalJSON of function type

### Fix decoding of slice of pointer type

If there is a pointer value, go-json will use it. (This behavior is necessary to achieve the ability to prioritize pre-filled values). However, since slices are reused internally, there was a bug that referred to the previous pointer value. Therefore, it is not necessary to refer to the pointer value in advance for the slice element, so we explicitly initialize slice element by `nil`.

# v0.4.8 - 2021/3/21

### Reduce memory usage at compile time

* go-json have used about 2GB of memory at compile time, but now it can compile with about less than 550MB.

### Fix any encoder's bug

* Add many test cases for encoder
* Fix composite type ( slice/array/map )
* Fix pointer types
* Fix encoding of MarshalJSON or MarshalText or json.Number type

### Refactor encoder

* Change package layout for reducing memory usage at compile
* Remove anonymous and only operation
* Remove root property from encodeCompileContext and opcode

### Fix CI

* Add Go 1.16
* Remove Go 1.13
* Fix `make cover` task

### Number/Delim/Token/RawMessage use the types defined in encoding/json by type alias

# v0.4.7 - 2021/02/22

### Fix decoder

* Fix decoding of deep recursive structure
* Fix decoding of embedded unexported pointer field
* Fix invalid test case
* Fix decoding of invalid value
* Fix decoding of prefilled value
* Fix not being able to return UnmarshalTypeError when it should be returned
* Fix decoding of null value
* Fix decoding of type of null string
* Use pre allocated pointer if exists it at decoding

### Reduce memory usage at compile

* Integrate int/int8/int16/int32/int64 and uint/uint8/uint16/uint32/uint64 operation to reduce memory usage at compile

### Remove unnecessary optype
