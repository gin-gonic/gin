[![Build](https://img.shields.io/circleci/build/github/leodido/go-urn?style=for-the-badge)](https://app.circleci.com/pipelines/github/leodido/go-urn) [![Coverage](https://img.shields.io/codecov/c/github/leodido/go-urn.svg?style=for-the-badge)](https://codecov.io/gh/leodido/go-urn) [![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/leodido/go-urn)

**A parser for URNs**.

> As seen on [RFC 2141](https://datatracker.ietf.org/doc/html/rfc2141), [RFC 7643](https://datatracker.ietf.org/doc/html/rfc7643#section-10), and on [RFC 8141](https://datatracker.ietf.org/doc/html/rfc8141).

[API documentation](https://godoc.org/github.com/leodido/go-urn).

Starting with version 1.3 this library also supports [RFC 7643 SCIM URNs](https://datatracker.ietf.org/doc/html/rfc7643#section-10).

Starting with version 1.4 this library also supports [RFC 8141 URNs (2017)](https://datatracker.ietf.org/doc/html/rfc8141).

## Installation

```
go get github.com/leodido/go-urn
```

## Features

1. RFC 2141 URNs parsing (default)
2. RFC 8141 URNs parsing (supersedes RFC 2141)
3. RFC 7643 SCIM URNs parsing
4. Normalization as per RFCs
5. Lexical equivalence as per RFCs
6. Precise, fine-grained errors

## Performances

This implementation results to be really fast.

Usually below 400 ns on my machine<sup>[1](#mymachine)</sup>.

Notice it also performs, while parsing:

1. fine-grained and informative erroring
2. specific-string normalization

```
ok/00/urn:a:b______________________________________/-10    51372006    109.0 ns/op    275 B/op    3 allocs/op
ok/01/URN:foo:a123,456_____________________________/-10    36024072    160.8 ns/op    296 B/op    6 allocs/op
ok/02/urn:foo:a123%2C456___________________________/-10    31901007    188.4 ns/op    320 B/op    7 allocs/op
ok/03/urn:ietf:params:scim:schemas:core:2.0:User___/-10    22736756    266.6 ns/op    376 B/op    6 allocs/op
ok/04/urn:ietf:params:scim:schemas:extension:enterp/-10    18291859    335.2 ns/op    408 B/op    6 allocs/op
ok/05/urn:ietf:params:scim:schemas:extension:enterp/-10    15283087    379.4 ns/op    440 B/op    6 allocs/op
ok/06/urn:burnout:nss______________________________/-10    39407593    155.1 ns/op    288 B/op    6 allocs/op
ok/07/urn:abcdefghilmnopqrstuvzabcdefghilm:x_______/-10    27832718    211.4 ns/op    307 B/op    4 allocs/op
ok/08/urn:urnurnurn:urn____________________________/-10    33269596    168.1 ns/op    293 B/op    6 allocs/op
ok/09/urn:ciao:!!*_________________________________/-10    41100675    148.8 ns/op    288 B/op    6 allocs/op
ok/10/urn:ciao:=@__________________________________/-10    37214253    149.7 ns/op    284 B/op    6 allocs/op
ok/11/urn:ciao:@!=%2C(xyz)+a,b.*@g=$_'_____________/-10    26534240    229.8 ns/op    336 B/op    7 allocs/op
ok/12/URN:x:abc%1Dz%2F%3az_________________________/-10    28166396    211.8 ns/op    336 B/op    7 allocs/op
no/13/URN:---xxx:x_________________________________/-10    23635159    255.6 ns/op    419 B/op    5 allocs/op
no/14/urn::colon:nss_______________________________/-10    23594779    258.4 ns/op    419 B/op    5 allocs/op
no/15/URN:@,:x_____________________________________/-10    23742535    261.5 ns/op    419 B/op    5 allocs/op
no/16/URN:URN:NSS__________________________________/-10    27432714    223.3 ns/op    371 B/op    5 allocs/op
no/17/urn:UrN:NSS__________________________________/-10    26922117    224.9 ns/op    371 B/op    5 allocs/op
no/18/urn:a:%______________________________________/-10    24926733    224.6 ns/op    371 B/op    5 allocs/op
no/19/urn:urn:NSS__________________________________/-10    27652641    220.7 ns/op    371 B/op    5 allocs/op
```

* <a name="mymachine">[1]</a>: Apple M1 Pro


## Example

For more examples take a look at the [examples file](examples_test.go).


```go
package main

import (
	"fmt"
	"github.com/leodido/go-urn"
)

func main() {
	var uid = "URN:foo:a123,456"

    // Parse the input string as a RFC 2141 URN only
	u, e := urn.NewMachine().Parse(uid)
	if e != nil {
		fmt.Errorf(err)

		return
	}

	fmt.Println(u.ID)
	fmt.Println(u.SS)

	// Output:
	// foo
	// a123,456
}
```

```go
package main

import (
	"fmt"
	"github.com/leodido/go-urn"
)

func main() {
	var uid = "URN:foo:a123,456"

    // Parse the input string as a RFC 2141 URN only
	u, ok := urn.Parse([]byte(uid))
	if !ok {
		panic("error parsing urn")
	}

	fmt.Println(u.ID)
	fmt.Println(u.SS)

	// Output:
	// foo
	// a123,456
}
```

```go
package main

import (
	"fmt"
	"github.com/leodido/go-urn"
)

func main() {
	input := "urn:ietf:params:scim:api:messages:2.0:ListResponse"

	// Parsing the input string as a RFC 7643 SCIM URN
	u, ok := urn.Parse([]byte(input), urn.WithParsingMode(urn.RFC7643Only))
	if !ok {
		panic("error parsing urn")
	}

	fmt.Println(u.IsSCIM())
	scim := u.SCIM()
	fmt.Println(scim.Type.String())
	fmt.Println(scim.Name)
	fmt.Println(scim.Other)

	// Output:
	// true
	// api
	// messages
	// 2.0:ListResponse
}
```