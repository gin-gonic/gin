## universal-translator
<img align="right" src="https://raw.githubusercontent.com/go-playground/universal-translator/master/logo.png">![Project status](https://img.shields.io/badge/version-0.18.1-green.svg)
[![Coverage Status](https://coveralls.io/repos/github/go-playground/universal-translator/badge.svg)](https://coveralls.io/github/go-playground/universal-translator)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-playground/universal-translator)](https://goreportcard.com/report/github.com/go-playground/universal-translator)
[![GoDoc](https://godoc.org/github.com/go-playground/universal-translator?status.svg)](https://godoc.org/github.com/go-playground/universal-translator)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Universal Translator is an i18n Translator for Go/Golang using CLDR data + pluralization rules

Why another i18n library?
--------------------------
Because none of the plural rules seem to be correct out there, including the previous implementation of this package,
so I took it upon myself to create [locales](https://github.com/go-playground/locales) for everyone to use; this package 
is a thin wrapper around [locales](https://github.com/go-playground/locales) in order to store and translate text for 
use in your applications.

Features
--------
- [x] Rules generated from the [CLDR](http://cldr.unicode.org/index/downloads) data, v36.0.1
- [x] Contains Cardinal, Ordinal and Range Plural Rules
- [x] Contains Month, Weekday and Timezone translations built in
- [x] Contains Date & Time formatting functions
- [x] Contains Number, Currency, Accounting and Percent formatting functions
- [x] Supports the "Gregorian" calendar only ( my time isn't unlimited, had to draw the line somewhere )
- [x] Support loading translations from files
- [x] Exporting translations to file(s), mainly for getting them professionally translated
- [ ] Code Generation for translation files -> Go code.. i.e. after it has been professionally translated
- [ ] Tests for all languages, I need help with this, please see [here](https://github.com/go-playground/locales/issues/1)

Installation
-----------

Use go get 

```shell
go get github.com/go-playground/universal-translator
```

Usage & Documentation
-------

Please see https://godoc.org/github.com/go-playground/universal-translator for usage docs

##### Examples:

- [Basic](https://github.com/go-playground/universal-translator/tree/master/_examples/basic)
- [Full - no files](https://github.com/go-playground/universal-translator/tree/master/_examples/full-no-files)
- [Full - with files](https://github.com/go-playground/universal-translator/tree/master/_examples/full-with-files)

File formatting
--------------
All types, Plain substitution, Cardinal, Ordinal and Range translations can all be contained within the same file(s);
they are only separated for easy viewing.

##### Examples:

- [Formats](https://github.com/go-playground/universal-translator/tree/master/_examples/file-formats)

##### Basic Makeup
NOTE: not all fields are needed for all translation types, see [examples](https://github.com/go-playground/universal-translator/tree/master/_examples/file-formats)
```json
{
    "locale": "en",
    "key": "days-left",
    "trans": "You have {0} day left.",
    "type": "Cardinal",
    "rule": "One",
    "override": false
}
```
|Field|Description|
|---|---|
|locale|The locale for which the translation is for.|
|key|The translation key that will be used to store and lookup each translation; normally it is a string or integer.|
|trans|The actual translation text.|
|type|The type of translation Cardinal, Ordinal, Range or "" for a plain substitution(not required to be defined if plain used)|
|rule|The plural rule for which the translation is for eg. One, Two, Few, Many or Other.(not required to be defined if plain used)|
|override|If you wish to override an existing translation that has already been registered, set this to 'true'. 99% of the time there is no need to define it.|

Help With Tests
---------------
To anyone interesting in helping or contributing, I sure could use some help creating tests for each language.
Please see issue [here](https://github.com/go-playground/locales/issues/1) for details.

License
------
Distributed under MIT License, please see license file in code for more details.
