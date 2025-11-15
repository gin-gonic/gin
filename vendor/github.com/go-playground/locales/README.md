## locales
<img align="right" src="https://raw.githubusercontent.com/go-playground/locales/master/logo.png">![Project status](https://img.shields.io/badge/version-0.14.1-green.svg)
[![Build Status](https://travis-ci.org/go-playground/locales.svg?branch=master)](https://travis-ci.org/go-playground/locales)
[![GoDoc](https://godoc.org/github.com/go-playground/locales?status.svg)](https://godoc.org/github.com/go-playground/locales)
![License](https://img.shields.io/dub/l/vibe-d.svg)

Locales is a set of locales generated from the [Unicode CLDR Project](http://cldr.unicode.org/) which can be used independently or within
an i18n package; these were built for use with, but not exclusive to, [Universal Translator](https://github.com/go-playground/universal-translator).

Features
--------
- [x] Rules generated from the latest [CLDR](http://cldr.unicode.org/index/downloads) data, v36.0.1
- [x] Contains Cardinal, Ordinal and Range Plural Rules
- [x] Contains Month, Weekday and Timezone translations built in
- [x] Contains Date & Time formatting functions
- [x] Contains Number, Currency, Accounting and Percent formatting functions
- [x] Supports the "Gregorian" calendar only ( my time isn't unlimited, had to draw the line somewhere )

Full Tests
--------------------
I could sure use your help adding tests for every locale, it is a huge undertaking and I just don't have the free time to do it all at the moment;
any help would be **greatly appreciated!!!!** please see [issue](https://github.com/go-playground/locales/issues/1) for details.

Installation
-----------

Use go get 

```shell
go get github.com/go-playground/locales
```  

NOTES
--------
You'll notice most return types are []byte, this is because most of the time the results will be concatenated with a larger body
of text and can avoid some allocations if already appending to a byte array, otherwise just cast as string.

Usage
-------
```go
package main

import (
	"fmt"
	"time"

	"github.com/go-playground/locales/currency"
	"github.com/go-playground/locales/en_CA"
)

func main() {

	loc, _ := time.LoadLocation("America/Toronto")
	datetime := time.Date(2016, 02, 03, 9, 0, 1, 0, loc)

	l := en_CA.New()

	// Dates
	fmt.Println(l.FmtDateFull(datetime))
	fmt.Println(l.FmtDateLong(datetime))
	fmt.Println(l.FmtDateMedium(datetime))
	fmt.Println(l.FmtDateShort(datetime))

	// Times
	fmt.Println(l.FmtTimeFull(datetime))
	fmt.Println(l.FmtTimeLong(datetime))
	fmt.Println(l.FmtTimeMedium(datetime))
	fmt.Println(l.FmtTimeShort(datetime))

	// Months Wide
	fmt.Println(l.MonthWide(time.January))
	fmt.Println(l.MonthWide(time.February))
	fmt.Println(l.MonthWide(time.March))
	// ...

	// Months Abbreviated
	fmt.Println(l.MonthAbbreviated(time.January))
	fmt.Println(l.MonthAbbreviated(time.February))
	fmt.Println(l.MonthAbbreviated(time.March))
	// ...

	// Months Narrow
	fmt.Println(l.MonthNarrow(time.January))
	fmt.Println(l.MonthNarrow(time.February))
	fmt.Println(l.MonthNarrow(time.March))
	// ...

	// Weekdays Wide
	fmt.Println(l.WeekdayWide(time.Sunday))
	fmt.Println(l.WeekdayWide(time.Monday))
	fmt.Println(l.WeekdayWide(time.Tuesday))
	// ...

	// Weekdays Abbreviated
	fmt.Println(l.WeekdayAbbreviated(time.Sunday))
	fmt.Println(l.WeekdayAbbreviated(time.Monday))
	fmt.Println(l.WeekdayAbbreviated(time.Tuesday))
	// ...

	// Weekdays Short
	fmt.Println(l.WeekdayShort(time.Sunday))
	fmt.Println(l.WeekdayShort(time.Monday))
	fmt.Println(l.WeekdayShort(time.Tuesday))
	// ...

	// Weekdays Narrow
	fmt.Println(l.WeekdayNarrow(time.Sunday))
	fmt.Println(l.WeekdayNarrow(time.Monday))
	fmt.Println(l.WeekdayNarrow(time.Tuesday))
	// ...

	var f64 float64

	f64 = -10356.4523

	// Number
	fmt.Println(l.FmtNumber(f64, 2))

	// Currency
	fmt.Println(l.FmtCurrency(f64, 2, currency.CAD))
	fmt.Println(l.FmtCurrency(f64, 2, currency.USD))

	// Accounting
	fmt.Println(l.FmtAccounting(f64, 2, currency.CAD))
	fmt.Println(l.FmtAccounting(f64, 2, currency.USD))

	f64 = 78.12

	// Percent
	fmt.Println(l.FmtPercent(f64, 0))

	// Plural Rules for locale, so you know what rules you must cover
	fmt.Println(l.PluralsCardinal())
	fmt.Println(l.PluralsOrdinal())

	// Cardinal Plural Rules
	fmt.Println(l.CardinalPluralRule(1, 0))
	fmt.Println(l.CardinalPluralRule(1.0, 0))
	fmt.Println(l.CardinalPluralRule(1.0, 1))
	fmt.Println(l.CardinalPluralRule(3, 0))

	// Ordinal Plural Rules
	fmt.Println(l.OrdinalPluralRule(21, 0)) // 21st
	fmt.Println(l.OrdinalPluralRule(22, 0)) // 22nd
	fmt.Println(l.OrdinalPluralRule(33, 0)) // 33rd
	fmt.Println(l.OrdinalPluralRule(34, 0)) // 34th

	// Range Plural Rules
	fmt.Println(l.RangePluralRule(1, 0, 1, 0)) // 1-1
	fmt.Println(l.RangePluralRule(1, 0, 2, 0)) // 1-2
	fmt.Println(l.RangePluralRule(5, 0, 8, 0)) // 5-8
}
```

NOTES:
-------
These rules were generated from the [Unicode CLDR Project](http://cldr.unicode.org/), if you encounter any issues
I strongly encourage contributing to the CLDR project to get the locale information corrected and the next time 
these locales are regenerated the fix will come with.

I do however realize that time constraints are often important and so there are two options:

1. Create your own locale, copy, paste and modify, and ensure it complies with the `Translator` interface.
2. Add an exception in the locale generation code directly and once regenerated, fix will be in place.

Please to not make fixes inside the locale files, they WILL get overwritten when the locales are regenerated.

License
------
Distributed under MIT License, please see license file in code for more details.
