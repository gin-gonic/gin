Package validator
=================
<img align="right" src="logo.png">[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/go-playground/validator)](https://github.com/go-playground/validator/releases)
[![Build Status](https://github.com/go-playground/validator/actions/workflows/workflow.yml/badge.svg)](https://github.com/go-playground/validator/actions)
[![Coverage Status](https://coveralls.io/repos/go-playground/validator/badge.svg?branch=master&service=github)](https://coveralls.io/github/go-playground/validator?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/validator)](https://goreportcard.com/report/github.com/go-playground/validator)
[![GoDoc](https://godoc.org/github.com/go-playground/validator?status.svg)](https://pkg.go.dev/github.com/go-playground/validator/v10)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Package validator implements value validations for structs and individual fields based on tags.

It has the following **unique** features:

-   Cross Field and Cross Struct validations by using validation tags or custom validators.
-   Slice, Array and Map diving, which allows any or all levels of a multidimensional field to be validated.
-   Ability to dive into both map keys and values for validation
-   Handles type interface by determining it's underlying type prior to validation.
-   Handles custom field types such as sql driver Valuer see [Valuer](https://golang.org/src/database/sql/driver/types.go?s=1210:1293#L29)
-   Alias validation tags, which allows for mapping of several validations to a single tag for easier defining of validations on structs
-   Extraction of custom defined Field Name e.g. can specify to extract the JSON name while validating and have it available in the resulting FieldError
-   Customizable i18n aware error messages.
-   Default validator for the [gin](https://github.com/gin-gonic/gin) web framework; upgrading from v8 to v9 in gin see [here](https://github.com/go-playground/validator/tree/master/_examples/gin-upgrading-overriding)

A Call for Maintainers
----------------------

Please read the discussiong started [here](https://github.com/go-playground/validator/discussions/1330) if you are interested in contributing/helping maintain this package.

Installation
------------

Use go get.

	go get github.com/go-playground/validator/v10

Then import the validator package into your own code.

	import "github.com/go-playground/validator/v10"

Error Return Value
-------

Validation functions return type error

They return type error to avoid the issue discussed in the following, where err is always != nil:

* http://stackoverflow.com/a/29138676/3158232
* https://github.com/go-playground/validator/issues/134

Validator returns only InvalidValidationError for bad validation input, nil or ValidationErrors as type error; so, in your code all you need to do is check if the error returned is not nil, and if it's not check if error is InvalidValidationError ( if necessary, most of the time it isn't ) type cast it to type ValidationErrors like so:

```go
err := validate.Struct(mystruct)
validationErrors := err.(validator.ValidationErrors)
 ```

Usage and documentation
------

Please see https://pkg.go.dev/github.com/go-playground/validator/v10 for detailed usage docs.

##### Examples:

- [Simple](https://github.com/go-playground/validator/blob/master/_examples/simple/main.go)
- [Custom Field Types](https://github.com/go-playground/validator/blob/master/_examples/custom/main.go)
- [Struct Level](https://github.com/go-playground/validator/blob/master/_examples/struct-level/main.go)
- [Translations & Custom Errors](https://github.com/go-playground/validator/blob/master/_examples/translations/main.go)
- [Gin upgrade and/or override validator](https://github.com/go-playground/validator/tree/v9/_examples/gin-upgrading-overriding)
- [wash - an example application putting it all together](https://github.com/bluesuncorp/wash)

Baked-in Validations
------

### Special Notes:
- If new to using validator it is highly recommended to initialize it using the `WithRequiredStructEnabled` option which is opt-in to new behaviour that will become the default behaviour in v11+. See documentation for more details.
```go
validate := validator.New(validator.WithRequiredStructEnabled())
```

### Fields:

| Tag | Description |
| - | - |
| eqcsfield | Field Equals Another Field (relative)|
| eqfield | Field Equals Another Field |
| fieldcontains | Check the indicated characters are present in the Field |
| fieldexcludes | Check the indicated characters are not present in the field |
| gtcsfield | Field Greater Than Another Relative Field |
| gtecsfield | Field Greater Than or Equal To Another Relative Field |
| gtefield | Field Greater Than or Equal To Another Field |
| gtfield | Field Greater Than Another Field |
| ltcsfield | Less Than Another Relative Field |
| ltecsfield | Less Than or Equal To Another Relative Field |
| ltefield | Less Than or Equal To Another Field |
| ltfield | Less Than Another Field |
| necsfield | Field Does Not Equal Another Field (relative) |
| nefield | Field Does Not Equal Another Field |

### Network:

| Tag | Description |
| - | - |
| cidr | Classless Inter-Domain Routing CIDR |
| cidrv4 | Classless Inter-Domain Routing CIDRv4 |
| cidrv6 | Classless Inter-Domain Routing CIDRv6 |
| datauri | Data URL |
| fqdn | Full Qualified Domain Name (FQDN) |
| hostname | Hostname RFC 952 |
| hostname_rfc1123 | Hostname RFC 1123 |
| hostname_port | HostPort |
| port | Port number |
| ip | Internet Protocol Address IP |
| ip4_addr | Internet Protocol Address IPv4 |
| ip6_addr | Internet Protocol Address IPv6 |
| ip_addr | Internet Protocol Address IP |
| ipv4 | Internet Protocol Address IPv4 |
| ipv6 | Internet Protocol Address IPv6 |
| mac | Media Access Control Address MAC |
| tcp4_addr | Transmission Control Protocol Address TCPv4 |
| tcp6_addr | Transmission Control Protocol Address TCPv6 |
| tcp_addr | Transmission Control Protocol Address TCP |
| udp4_addr | User Datagram Protocol Address UDPv4 |
| udp6_addr | User Datagram Protocol Address UDPv6 |
| udp_addr | User Datagram Protocol Address UDP |
| unix_addr | Unix domain socket end point Address |
| uri | URI String |
| url | URL String |
| http_url | HTTP(s) URL String |
| https_url | HTTPS-only URL String |
| url_encoded | URL Encoded |
| urn_rfc2141 | Urn RFC 2141 String |

### Strings:

| Tag | Description |
| - | - |
| alpha | Alpha Only |
| alphaspace | Alpha Space |
| alphanum | Alphanumeric |
| alphanumunicode | Alphanumeric Unicode |
| alphaunicode | Alpha Unicode |
| ascii | ASCII |
| boolean | Boolean |
| contains | Contains |
| containsany | Contains Any |
| containsrune | Contains Rune |
| endsnotwith | Ends Not With |
| endswith | Ends With |
| excludes | Excludes |
| excludesall | Excludes All |
| excludesrune | Excludes Rune |
| lowercase | Lowercase |
| multibyte | Multi-Byte Characters |
| number | Number |
| numeric | Numeric |
| printascii | Printable ASCII |
| startsnotwith | Starts Not With |
| startswith | Starts With |
| uppercase | Uppercase |

### Format:
| Tag | Description |
| - | - |
| base64 | Base64 String |
| base64url | Base64URL String |
| base64rawurl | Base64RawURL String |
| bic | Business Identifier Code (ISO 9362) |
| bcp47_language_tag | Language tag (BCP 47) |
| btc_addr | Bitcoin Address |
| btc_addr_bech32 | Bitcoin Bech32 Address (segwit) |
| credit_card | Credit Card Number |
| mongodb | MongoDB ObjectID |
| mongodb_connection_string | MongoDB Connection String |
| cron | Cron |
| spicedb | SpiceDb ObjectID/Permission/Type |
| datetime | Datetime |
| e164 | e164 formatted phone number |
| ein | U.S. Employeer Identification Number |
| email | E-mail String
| eth_addr | Ethereum Address |
| hexadecimal | Hexadecimal String |
| hexcolor | Hexcolor String |
| hsl | HSL String |
| hsla | HSLA String |
| html | HTML Tags |
| html_encoded | HTML Encoded |
| isbn | International Standard Book Number |
| isbn10 | International Standard Book Number 10 |
| isbn13 | International Standard Book Number 13 |
| issn | International Standard Serial Number |
| iso3166_1_alpha2 | Two-letter country code (ISO 3166-1 alpha-2) |
| iso3166_1_alpha3 | Three-letter country code (ISO 3166-1 alpha-3) |
| iso3166_1_alpha_numeric | Numeric country code (ISO 3166-1 numeric) |
| iso3166_2 | Country subdivision code (ISO 3166-2) |
| iso4217 | Currency code (ISO 4217) |
| json | JSON |
| jwt | JSON Web Token (JWT) |
| latitude | Latitude |
| longitude | Longitude |
| luhn_checksum | Luhn Algorithm Checksum (for strings and (u)int) |
| postcode_iso3166_alpha2 | Postcode |
| postcode_iso3166_alpha2_field | Postcode |
| rgb | RGB String |
| rgba | RGBA String |
| ssn | Social Security Number SSN |
| timezone | Timezone |
| uuid | Universally Unique Identifier UUID |
| uuid3 | Universally Unique Identifier UUID v3 |
| uuid3_rfc4122 | Universally Unique Identifier UUID v3 RFC4122 |
| uuid4 | Universally Unique Identifier UUID v4 |
| uuid4_rfc4122 | Universally Unique Identifier UUID v4 RFC4122 |
| uuid5 | Universally Unique Identifier UUID v5 |
| uuid5_rfc4122 | Universally Unique Identifier UUID v5 RFC4122 |
| uuid_rfc4122 | Universally Unique Identifier UUID RFC4122 |
| md4 | MD4 hash |
| md5 | MD5 hash |
| sha256 | SHA256 hash |
| sha384 | SHA384 hash |
| sha512 | SHA512 hash |
| ripemd128 | RIPEMD-128 hash |
| ripemd128 | RIPEMD-160 hash |
| tiger128 | TIGER128 hash |
| tiger160 | TIGER160 hash |
| tiger192 | TIGER192 hash |
| semver | Semantic Versioning 2.0.0 |
| ulid | Universally Unique Lexicographically Sortable Identifier ULID |
| cve | Common Vulnerabilities and Exposures Identifier (CVE id) |

### Comparisons:
| Tag | Description |
| - | - |
| eq | Equals |
| eq_ignore_case | Equals ignoring case |
| gt | Greater than|
| gte | Greater than or equal |
| lt | Less Than |
| lte | Less Than or Equal |
| ne | Not Equal |
| ne_ignore_case | Not Equal ignoring case |

### Other:
| Tag | Description |
| - | - |
| dir | Existing Directory |
| dirpath | Directory Path |
| file | Existing File |
| filepath | File Path |
| image | Image |
| isdefault | Is Default |
| len | Length |
| max | Maximum |
| min | Minimum |
| oneof | One Of |
| required | Required |
| required_if | Required If |
| required_unless | Required Unless |
| required_with | Required With |
| required_with_all | Required With All |
| required_without | Required Without |
| required_without_all | Required Without All |
| excluded_if | Excluded If |
| excluded_unless | Excluded Unless |
| excluded_with | Excluded With |
| excluded_with_all | Excluded With All |
| excluded_without | Excluded Without |
| excluded_without_all | Excluded Without All |
| unique | Unique |
| validateFn | Verify if the method `Validate() error` does not return an error (or any specified method) |


#### Aliases:
| Tag | Description |
| - | - |
| iscolor | hexcolor\|rgb\|rgba\|hsl\|hsla |
| country_code | iso3166_1_alpha2\|iso3166_1_alpha3\|iso3166_1_alpha_numeric |

Benchmarks
------
###### Run on MacBook Pro Max M3
```go
go version go1.23.3 darwin/arm64
goos: darwin
goarch: arm64
cpu: Apple M3 Max
pkg: github.com/go-playground/validator/v10
BenchmarkFieldSuccess-16                                                42461943                27.88 ns/op            0 B/op          0 allocs/op
BenchmarkFieldSuccessParallel-16                                        486632887                2.289 ns/op           0 B/op          0 allocs/op
BenchmarkFieldFailure-16                                                 9566167               121.3 ns/op           200 B/op          4 allocs/op
BenchmarkFieldFailureParallel-16                                        17551471                83.68 ns/op          200 B/op          4 allocs/op
BenchmarkFieldArrayDiveSuccess-16                                        7602306               155.6 ns/op            97 B/op          5 allocs/op
BenchmarkFieldArrayDiveSuccessParallel-16                               20664610                59.80 ns/op           97 B/op          5 allocs/op
BenchmarkFieldArrayDiveFailure-16                                        4659756               252.9 ns/op           301 B/op         10 allocs/op
BenchmarkFieldArrayDiveFailureParallel-16                                8010116               152.9 ns/op           301 B/op         10 allocs/op
BenchmarkFieldMapDiveSuccess-16                                          2834575               421.2 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveSuccessParallel-16                                  7179700               171.8 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveFailure-16                                          3081728               384.4 ns/op           376 B/op         13 allocs/op
BenchmarkFieldMapDiveFailureParallel-16                                  6058137               204.0 ns/op           377 B/op         13 allocs/op
BenchmarkFieldMapDiveWithKeysSuccess-16                                  2544975               464.8 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveWithKeysSuccessParallel-16                          6661954               181.4 ns/op           288 B/op         14 allocs/op
BenchmarkFieldMapDiveWithKeysFailure-16                                  2435484               490.7 ns/op           553 B/op         16 allocs/op
BenchmarkFieldMapDiveWithKeysFailureParallel-16                          4249617               282.0 ns/op           554 B/op         16 allocs/op
BenchmarkFieldCustomTypeSuccess-16                                      14943525                77.35 ns/op           32 B/op          2 allocs/op
BenchmarkFieldCustomTypeSuccessParallel-16                              64051954                20.61 ns/op           32 B/op          2 allocs/op
BenchmarkFieldCustomTypeFailure-16                                      10721384               107.1 ns/op           184 B/op          3 allocs/op
BenchmarkFieldCustomTypeFailureParallel-16                              18714495                69.77 ns/op          184 B/op          3 allocs/op
BenchmarkFieldOrTagSuccess-16                                            4063124               294.3 ns/op            16 B/op          1 allocs/op
BenchmarkFieldOrTagSuccessParallel-16                                   31903756                41.22 ns/op           18 B/op          1 allocs/op
BenchmarkFieldOrTagFailure-16                                            7748558               146.8 ns/op           216 B/op          5 allocs/op
BenchmarkFieldOrTagFailureParallel-16                                   13139854                92.05 ns/op          216 B/op          5 allocs/op
BenchmarkStructLevelValidationSuccess-16                                16808389                70.25 ns/op           16 B/op          1 allocs/op
BenchmarkStructLevelValidationSuccessParallel-16                        90686955                14.47 ns/op           16 B/op          1 allocs/op
BenchmarkStructLevelValidationFailure-16                                 5818791               200.2 ns/op           264 B/op          7 allocs/op
BenchmarkStructLevelValidationFailureParallel-16                        11115874               107.5 ns/op           264 B/op          7 allocs/op
BenchmarkStructSimpleCustomTypeSuccess-16                                7764956               151.9 ns/op            32 B/op          2 allocs/op
BenchmarkStructSimpleCustomTypeSuccessParallel-16                       52316265                30.37 ns/op           32 B/op          2 allocs/op
BenchmarkStructSimpleCustomTypeFailure-16                                4195429               277.2 ns/op           416 B/op          9 allocs/op
BenchmarkStructSimpleCustomTypeFailureParallel-16                        7305661               164.6 ns/op           432 B/op         10 allocs/op
BenchmarkStructFilteredSuccess-16                                        6312625               186.1 ns/op           216 B/op          5 allocs/op
BenchmarkStructFilteredSuccessParallel-16                               13684459                93.42 ns/op          216 B/op          5 allocs/op
BenchmarkStructFilteredFailure-16                                        6751482               171.2 ns/op           216 B/op          5 allocs/op
BenchmarkStructFilteredFailureParallel-16                               14146070                86.93 ns/op          216 B/op          5 allocs/op
BenchmarkStructPartialSuccess-16                                         6544448               177.3 ns/op           224 B/op          4 allocs/op
BenchmarkStructPartialSuccessParallel-16                                13951946                88.73 ns/op          224 B/op          4 allocs/op
BenchmarkStructPartialFailure-16                                         4075833               287.5 ns/op           440 B/op          9 allocs/op
BenchmarkStructPartialFailureParallel-16                                 7490805               161.3 ns/op           440 B/op          9 allocs/op
BenchmarkStructExceptSuccess-16                                          4107187               281.4 ns/op           424 B/op          8 allocs/op
BenchmarkStructExceptSuccessParallel-16                                 15979173                80.86 ns/op          208 B/op          3 allocs/op
BenchmarkStructExceptFailure-16                                          4434372               264.3 ns/op           424 B/op          8 allocs/op
BenchmarkStructExceptFailureParallel-16                                  8081367               154.1 ns/op           424 B/op          8 allocs/op
BenchmarkStructSimpleCrossFieldSuccess-16                                6459542               183.4 ns/op            56 B/op          3 allocs/op
BenchmarkStructSimpleCrossFieldSuccessParallel-16                       41013781                37.95 ns/op           56 B/op          3 allocs/op
BenchmarkStructSimpleCrossFieldFailure-16                                4034998               292.1 ns/op           272 B/op          8 allocs/op
BenchmarkStructSimpleCrossFieldFailureParallel-16                       11348446               115.3 ns/op           272 B/op          8 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldSuccess-16                     4448528               267.7 ns/op            64 B/op          4 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldSuccessParallel-16            26813619                48.33 ns/op           64 B/op          4 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldFailure-16                     3090646               384.5 ns/op           288 B/op          9 allocs/op
BenchmarkStructSimpleCrossStructCrossFieldFailureParallel-16             9870906               129.5 ns/op           288 B/op          9 allocs/op
BenchmarkStructSimpleSuccess-16                                         10675562               109.5 ns/op             0 B/op          0 allocs/op
BenchmarkStructSimpleSuccessParallel-16                                 131159784                8.932 ns/op           0 B/op          0 allocs/op
BenchmarkStructSimpleFailure-16                                          4094979               286.6 ns/op           416 B/op          9 allocs/op
BenchmarkStructSimpleFailureParallel-16                                  7606663               157.9 ns/op           416 B/op          9 allocs/op
BenchmarkStructComplexSuccess-16                                         2073470               576.0 ns/op           224 B/op          5 allocs/op
BenchmarkStructComplexSuccessParallel-16                                 7821831               161.3 ns/op           224 B/op          5 allocs/op
BenchmarkStructComplexFailure-16                                          576358              2001 ns/op            3042 B/op         48 allocs/op
BenchmarkStructComplexFailureParallel-16                                 1000000              1171 ns/op            3041 B/op         48 allocs/op
BenchmarkOneof-16                                                       22503973                52.82 ns/op            0 B/op          0 allocs/op
BenchmarkOneofParallel-16                                                8538474               140.4 ns/op             0 B/op          0 allocs/op
```

Complementary Software
----------------------

Here is a list of software that complements using this library either pre or post validation.

* [form](https://github.com/go-playground/form) - Decodes url.Values into Go value(s) and Encodes Go value(s) into url.Values. Dual Array and Full map support.
* [mold](https://github.com/go-playground/mold) - A general library to help modify or set data within data structures and other objects

How to Contribute
------

Make a pull request...

Maintenance and support for SDK major versions
----------------------------------------------

See prior discussion [here](https://github.com/go-playground/validator/discussions/1342) for more details.

This package is aligned with the [Go release policy](https://go.dev/doc/devel/release) in that support is guaranteed for 
the two most recent major versions.

This does not mean the package will not work with older versions of Go, only that we reserve the right to increase the 
MSGV(Minimum Supported Go Version) when the need arises to address Security issues/patches, OS issues & support or newly 
introduced functionality that would greatly benefit the maintenance and/or usage of this package.

If and when the MSGV is increased it will be done so in a minimum of a `Minor` release bump.

License
-------
Distributed under MIT License, please see license file within the code for more details.

Maintainers
-----------
This project has grown large enough that more than one person is required to properly support the community.
If you are interested in becoming a maintainer please reach out to me https://github.com/deankarn
