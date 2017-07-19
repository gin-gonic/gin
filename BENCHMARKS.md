
## Benchmark System

**VM HOST:** DigitalOcean  
**Machine:** 4 CPU, 8 GB RAM. Ubuntu 16.04.2 x64  
**Date:** July 19th, 2017  
**Go Version:** 1.8.3 linux/amd64  
**Source:** [Go HTTP Router Benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)  

## Static Routes: 157

```
Gin:             30512 Bytes

HttpServeMux:    17344 Bytes
Ace:             30080 Bytes
Bear:            30472 Bytes
Beego:           96408 Bytes
Bone:            37904 Bytes
Denco:           10464 Bytes
Echo:            73680 Bytes
GocraftWeb:      55720 Bytes
Goji:            27200 Bytes
Gojiv2:         104464 Bytes
GoJsonRest:     136472 Bytes
GoRestful:      914904 Bytes
GorillaMux:     675568 Bytes
HttpRouter:      21128 Bytes
HttpTreeMux:     73448 Bytes
Kocha:          115072 Bytes
LARS:            30120 Bytes
Macaron:         37984 Bytes
Martini:        310832 Bytes
Pat:             20464 Bytes
Possum:          91328 Bytes
R2router:        23712 Bytes
Rivet:           23880 Bytes
Tango:           28008 Bytes
TigerTonic:      80368 Bytes
Traffic:        626480 Bytes
Vulcan:         369064 Bytes
```

## GithubAPI Routes: 203

```
Gin:             52672 Bytes

Ace:             48992 Bytes
Bear:           161592 Bytes
Beego:          147992 Bytes
Bone:            97728 Bytes
Denco:           36440 Bytes
Echo:            95672 Bytes
GocraftWeb:      95640 Bytes
Goji:            86088 Bytes
Gojiv2:         144392 Bytes
GoJsonRest:     134648 Bytes
GoRestful:     1410760 Bytes
GorillaMux:    1509488 Bytes
HttpRouter:      37464 Bytes
HttpTreeMux:     78800 Bytes
Kocha:          785408 Bytes
LARS:            49032 Bytes
Macaron:        132712 Bytes
Martini:        564352 Bytes
Pat:             21200 Bytes
Possum:          83888 Bytes
R2router:        47104 Bytes
Rivet:           42840 Bytes
Tango:           54584 Bytes
TigerTonic:      96384 Bytes
Traffic:       1061920 Bytes
Vulcan:         465296 Bytes
```

## GPlusAPI Routes: 13

```
Gin:              3968 Bytes

Ace:              3600 Bytes
Bear:             7112 Bytes
Beego:           10048 Bytes
Bone:             6480 Bytes
Denco:            3256 Bytes
Echo:             9000 Bytes
GocraftWeb:       7496 Bytes
Goji:             2912 Bytes
Gojiv2:           7376 Bytes
GoJsonRest:      11544 Bytes
GoRestful:       88776 Bytes
GorillaMux:      71488 Bytes
HttpRouter:       2712 Bytes
HttpTreeMux:      7440 Bytes
Kocha:          128880 Bytes
LARS:             3640 Bytes
Macaron:          8656 Bytes
Martini:         23936 Bytes
Pat:              1856 Bytes
Possum:           7248 Bytes
R2router:         3928 Bytes
Rivet:            3064 Bytes
Tango:            4912 Bytes
TigerTonic:       9408 Bytes
Traffic:         49472 Bytes
Vulcan:          25496 Bytes
```

## ParseAPI Routes: 26

```
Gin:              6928 Bytes

Ace:              6592 Bytes
Bear:            12320 Bytes
Beego:           18960 Bytes
Bone:            11024 Bytes
Denco:            4184 Bytes
Echo:            11168 Bytes
GocraftWeb:      12800 Bytes
Goji:             5232 Bytes
Gojiv2:          14464 Bytes
GoJsonRest:      14216 Bytes
GoRestful:      127368 Bytes
GorillaMux:     123016 Bytes
HttpRouter:       4976 Bytes
HttpTreeMux:      7848 Bytes
Kocha:          181712 Bytes
LARS:             6632 Bytes
Macaron:         13648 Bytes
Martini:         45952 Bytes
Pat:              2560 Bytes
Possum:           9200 Bytes
R2router:         7056 Bytes
Rivet:            5680 Bytes
Tango:            8664 Bytes
TigerTonic:       9840 Bytes
Traffic:         93480 Bytes
Vulcan:          44504 Bytes
```

## Static Routes

```
BenchmarkGin_StaticAll                     50000             34506 ns/op               0 B/op          0 allocs/op

BenchmarkAce_StaticAll                     30000             49657 ns/op               0 B/op          0 allocs/op
BenchmarkHttpServeMux_StaticAll             2000           1183737 ns/op              96 B/op          8 allocs/op
BenchmarkBeego_StaticAll                    5000            412621 ns/op           57776 B/op        628 allocs/op
BenchmarkBear_StaticAll                    10000            149242 ns/op           20336 B/op        461 allocs/op
BenchmarkBone_StaticAll                    10000            118583 ns/op               0 B/op          0 allocs/op
BenchmarkDenco_StaticAll                  100000             13247 ns/op               0 B/op          0 allocs/op
BenchmarkEcho_StaticAll                    20000             79914 ns/op            5024 B/op        157 allocs/op
BenchmarkGocraftWeb_StaticAll              10000            211823 ns/op           46440 B/op        785 allocs/op
BenchmarkGoji_StaticAll                    10000            109390 ns/op               0 B/op          0 allocs/op
BenchmarkGojiv2_StaticAll                   3000            415533 ns/op          145696 B/op       1099 allocs/op
BenchmarkGoJsonRest_StaticAll               5000            364403 ns/op           51653 B/op       1727 allocs/op
BenchmarkGoRestful_StaticAll                 500           2578579 ns/op          314936 B/op       3144 allocs/op
BenchmarkGorillaMux_StaticAll                500           2704856 ns/op          115648 B/op       1578 allocs/op
BenchmarkHttpRouter_StaticAll             100000             18541 ns/op               0 B/op          0 allocs/op
BenchmarkHttpTreeMux_StaticAll            100000             22332 ns/op               0 B/op          0 allocs/op
BenchmarkKocha_StaticAll                   50000             31176 ns/op               0 B/op          0 allocs/op
BenchmarkLARS_StaticAll                    50000             40840 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_StaticAll                  5000            517656 ns/op          120576 B/op       1413 allocs/op
BenchmarkMartini_StaticAll                   300           4462289 ns/op          125442 B/op       1717 allocs/op
BenchmarkPat_StaticAll                       500           2157275 ns/op          533904 B/op      11123 allocs/op
BenchmarkPossum_StaticAll                  10000            254701 ns/op           65312 B/op        471 allocs/op
BenchmarkR2router_StaticAll                10000            133956 ns/op           22608 B/op        628 allocs/op
BenchmarkRivet_StaticAll                   30000             46812 ns/op               0 B/op          0 allocs/op
BenchmarkTango_StaticAll                    5000            390613 ns/op           39225 B/op       1256 allocs/op
BenchmarkTigerTonic_StaticAll              20000             88060 ns/op            7504 B/op        157 allocs/op
BenchmarkTraffic_StaticAll                   500           2910236 ns/op          729736 B/op      14287 allocs/op
BenchmarkVulcan_StaticAll                   5000            277366 ns/op           15386 B/op        471 allocs/op
```

## Micro Benchmarks

```
BenchmarkGin_Param                      20000000               113 ns/op               0 B/op          0 allocs/op

BenchmarkAce_Param                       5000000               375 ns/op              32 B/op          1 allocs/op
BenchmarkBear_Param                      1000000              1709 ns/op             456 B/op          5 allocs/op
BenchmarkBeego_Param                     1000000              2484 ns/op             368 B/op          4 allocs/op
BenchmarkBone_Param                      1000000              2391 ns/op             688 B/op          5 allocs/op
BenchmarkDenco_Param                    10000000               240 ns/op              32 B/op          1 allocs/op
BenchmarkEcho_Param                      5000000               366 ns/op              32 B/op          1 allocs/op
BenchmarkGocraftWeb_Param                1000000              2343 ns/op             648 B/op          8 allocs/op
BenchmarkGoji_Param                      1000000              1197 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_Param                    1000000              2771 ns/op             944 B/op          8 allocs/op
BenchmarkGoJsonRest_Param                1000000              2993 ns/op             649 B/op         13 allocs/op
BenchmarkGoRestful_Param                  200000              8860 ns/op            2296 B/op         21 allocs/op
BenchmarkGorillaMux_Param                 500000              4461 ns/op            1056 B/op         11 allocs/op
BenchmarkHttpRouter_Param               10000000               175 ns/op              32 B/op          1 allocs/op
BenchmarkHttpTreeMux_Param               1000000              1167 ns/op             352 B/op          3 allocs/op
BenchmarkKocha_Param                     3000000               429 ns/op              56 B/op          3 allocs/op
BenchmarkLARS_Param                     10000000               134 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_Param                    500000              4635 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_Param                    200000              9933 ns/op            1072 B/op         10 allocs/op
BenchmarkPat_Param                       1000000              2929 ns/op             648 B/op         12 allocs/op
BenchmarkPossum_Param                    1000000              2503 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_Param                  1000000              1507 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_Param                     5000000               297 ns/op              48 B/op          1 allocs/op
BenchmarkTango_Param                     1000000              1862 ns/op             248 B/op          8 allocs/op
BenchmarkTigerTonic_Param                 500000              5660 ns/op             992 B/op         17 allocs/op
BenchmarkTraffic_Param                    200000              8408 ns/op            1960 B/op         21 allocs/op
BenchmarkVulcan_Param                    2000000               963 ns/op              98 B/op          3 allocs/op
BenchmarkAce_Param5                      2000000               740 ns/op             160 B/op          1 allocs/op
BenchmarkBear_Param5                     1000000              2777 ns/op             501 B/op          5 allocs/op
BenchmarkBeego_Param5                    1000000              3740 ns/op             368 B/op          4 allocs/op
BenchmarkBone_Param5                     1000000              2950 ns/op             736 B/op          5 allocs/op
BenchmarkDenco_Param5                    2000000               644 ns/op             160 B/op          1 allocs/op
BenchmarkEcho_Param5                     3000000               558 ns/op              32 B/op          1 allocs/op
BenchmarkGin_Param5                     10000000               198 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_Param5                500000              3870 ns/op             920 B/op         11 allocs/op
BenchmarkGoji_Param5                     1000000              1746 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_Param5                   1000000              3214 ns/op            1008 B/op          8 allocs/op
BenchmarkGoJsonRest_Param5                500000              5509 ns/op            1097 B/op         16 allocs/op
BenchmarkGoRestful_Param5                 200000             11232 ns/op            2392 B/op         21 allocs/op
BenchmarkGorillaMux_Param5                300000              7777 ns/op            1184 B/op         11 allocs/op
BenchmarkHttpRouter_Param5               3000000               631 ns/op             160 B/op          1 allocs/op
BenchmarkHttpTreeMux_Param5              1000000              2800 ns/op             576 B/op          6 allocs/op
BenchmarkKocha_Param5                    1000000              2053 ns/op             440 B/op         10 allocs/op
BenchmarkLARS_Param5                    10000000               232 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_Param5                   500000              5888 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_Param5                   200000             12807 ns/op            1232 B/op         11 allocs/op
BenchmarkPat_Param5                       300000              7320 ns/op             964 B/op         32 allocs/op
BenchmarkPossum_Param5                   1000000              2495 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_Param5                 1000000              1844 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_Param5                    2000000               935 ns/op             240 B/op          1 allocs/op
BenchmarkTango_Param5                    1000000              2327 ns/op             360 B/op          8 allocs/op
BenchmarkTigerTonic_Param5                100000             18514 ns/op            2551 B/op         43 allocs/op
BenchmarkTraffic_Param5                   200000             11997 ns/op            2248 B/op         25 allocs/op
BenchmarkVulcan_Param5                   1000000              1333 ns/op              98 B/op          3 allocs/op
BenchmarkAce_Param20                     1000000              2031 ns/op             640 B/op          1 allocs/op
BenchmarkBear_Param20                     200000              7285 ns/op            1664 B/op          5 allocs/op
BenchmarkBeego_Param20                    300000              6224 ns/op             368 B/op          4 allocs/op
BenchmarkBone_Param20                     200000              8023 ns/op            1903 B/op          5 allocs/op
BenchmarkDenco_Param20                   1000000              2262 ns/op             640 B/op          1 allocs/op
BenchmarkEcho_Param20                    1000000              1387 ns/op              32 B/op          1 allocs/op
BenchmarkGin_Param20                     3000000               503 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_Param20               100000             14408 ns/op            3795 B/op         15 allocs/op
BenchmarkGoji_Param20                     500000              5272 ns/op            1247 B/op          2 allocs/op
BenchmarkGojiv2_Param20                  1000000              4163 ns/op            1248 B/op          8 allocs/op
BenchmarkGoJsonRest_Param20               100000             17866 ns/op            4485 B/op         20 allocs/op
BenchmarkGoRestful_Param20                100000             21022 ns/op            4724 B/op         23 allocs/op
BenchmarkGorillaMux_Param20               100000             17055 ns/op            3547 B/op         13 allocs/op
BenchmarkHttpRouter_Param20              1000000              1748 ns/op             640 B/op          1 allocs/op
BenchmarkHttpTreeMux_Param20              200000             12246 ns/op            3196 B/op         10 allocs/op
BenchmarkKocha_Param20                    300000              6861 ns/op            1808 B/op         27 allocs/op
BenchmarkLARS_Param20                    3000000               526 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_Param20                  100000             13069 ns/op            2906 B/op         12 allocs/op
BenchmarkMartini_Param20                  100000             23602 ns/op            3597 B/op         13 allocs/op
BenchmarkPat_Param20                       50000             32143 ns/op            4688 B/op        111 allocs/op
BenchmarkPossum_Param20                  1000000              2396 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_Param20                 200000              8907 ns/op            2283 B/op          7 allocs/op
BenchmarkRivet_Param20                   1000000              3280 ns/op            1024 B/op          1 allocs/op
BenchmarkTango_Param20                    500000              4640 ns/op             856 B/op          8 allocs/op
BenchmarkTigerTonic_Param20                20000             67581 ns/op           10532 B/op        138 allocs/op
BenchmarkTraffic_Param20                   50000             40313 ns/op            7941 B/op         45 allocs/op
BenchmarkVulcan_Param20                  1000000              2264 ns/op              98 B/op          3 allocs/op
BenchmarkAce_ParamWrite                  3000000               532 ns/op              40 B/op          2 allocs/op
BenchmarkBear_ParamWrite                 1000000              1778 ns/op             456 B/op          5 allocs/op
BenchmarkBeego_ParamWrite                1000000              2596 ns/op             376 B/op          5 allocs/op
BenchmarkBone_ParamWrite                 1000000              2519 ns/op             688 B/op          5 allocs/op
BenchmarkDenco_ParamWrite                5000000               411 ns/op              32 B/op          1 allocs/op
BenchmarkEcho_ParamWrite                 2000000               718 ns/op              40 B/op          2 allocs/op
BenchmarkGin_ParamWrite                  5000000               283 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_ParamWrite           1000000              2561 ns/op             656 B/op          9 allocs/op
BenchmarkGoji_ParamWrite                 1000000              1378 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_ParamWrite               1000000              3128 ns/op             976 B/op         10 allocs/op
BenchmarkGoJsonRest_ParamWrite            500000              4446 ns/op            1128 B/op         18 allocs/op
BenchmarkGoRestful_ParamWrite             200000             10291 ns/op            2304 B/op         22 allocs/op
BenchmarkGorillaMux_ParamWrite            500000              5153 ns/op            1064 B/op         12 allocs/op
BenchmarkHttpRouter_ParamWrite           5000000               263 ns/op              32 B/op          1 allocs/op
BenchmarkHttpTreeMux_ParamWrite          1000000              1351 ns/op             352 B/op          3 allocs/op
BenchmarkKocha_ParamWrite                3000000               538 ns/op              56 B/op          3 allocs/op
BenchmarkLARS_ParamWrite                 5000000               316 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_ParamWrite               500000              5756 ns/op            1160 B/op         14 allocs/op
BenchmarkMartini_ParamWrite               200000             13097 ns/op            1176 B/op         14 allocs/op
BenchmarkPat_ParamWrite                   500000              4954 ns/op            1072 B/op         17 allocs/op
BenchmarkPossum_ParamWrite               1000000              2499 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_ParamWrite             1000000              1531 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_ParamWrite                3000000               570 ns/op             112 B/op          2 allocs/op
BenchmarkTango_ParamWrite                2000000               957 ns/op             136 B/op          4 allocs/op
BenchmarkTigerTonic_ParamWrite            200000              7025 ns/op            1424 B/op         23 allocs/op
BenchmarkTraffic_ParamWrite               200000             10112 ns/op            2384 B/op         25 allocs/op
BenchmarkVulcan_ParamWrite               1000000              1006 ns/op              98 B/op          3 allocs/op
```

## GitHub

```
BenchmarkGin_GithubStatic               10000000               156 ns/op               0 B/op          0 allocs/op

BenchmarkAce_GithubStatic                5000000               294 ns/op               0 B/op          0 allocs/op
BenchmarkBear_GithubStatic               2000000               893 ns/op             120 B/op          3 allocs/op
BenchmarkBeego_GithubStatic              1000000              2491 ns/op             368 B/op          4 allocs/op
BenchmarkBone_GithubStatic                 50000             25300 ns/op            2880 B/op         60 allocs/op
BenchmarkDenco_GithubStatic             20000000                76.0 ns/op             0 B/op          0 allocs/op
BenchmarkEcho_GithubStatic               2000000               516 ns/op              32 B/op          1 allocs/op
BenchmarkGocraftWeb_GithubStatic         1000000              1448 ns/op             296 B/op          5 allocs/op
BenchmarkGoji_GithubStatic               3000000               496 ns/op               0 B/op          0 allocs/op
BenchmarkGojiv2_GithubStatic             1000000              2941 ns/op             928 B/op          7 allocs/op
BenchmarkGoRestful_GithubStatic           100000             27256 ns/op            3224 B/op         22 allocs/op
BenchmarkGoJsonRest_GithubStatic         1000000              2196 ns/op             329 B/op         11 allocs/op
BenchmarkGorillaMux_GithubStatic           50000             31617 ns/op             736 B/op         10 allocs/op
BenchmarkHttpRouter_GithubStatic        20000000                88.4 ns/op             0 B/op          0 allocs/op
BenchmarkHttpTreeMux_GithubStatic       10000000               134 ns/op               0 B/op          0 allocs/op
BenchmarkKocha_GithubStatic             20000000               113 ns/op               0 B/op          0 allocs/op
BenchmarkLARS_GithubStatic              10000000               195 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GithubStatic             500000              3740 ns/op             768 B/op          9 allocs/op
BenchmarkMartini_GithubStatic              50000             27673 ns/op             768 B/op          9 allocs/op
BenchmarkPat_GithubStatic                 100000             19470 ns/op            3648 B/op         76 allocs/op
BenchmarkPossum_GithubStatic             1000000              1729 ns/op             416 B/op          3 allocs/op
BenchmarkR2router_GithubStatic           2000000               879 ns/op             144 B/op          4 allocs/op
BenchmarkRivet_GithubStatic             10000000               231 ns/op               0 B/op          0 allocs/op
BenchmarkTango_GithubStatic              1000000              2325 ns/op             248 B/op          8 allocs/op
BenchmarkTigerTonic_GithubStatic         3000000               610 ns/op              48 B/op          1 allocs/op
BenchmarkTraffic_GithubStatic              20000             62973 ns/op           18904 B/op        148 allocs/op
BenchmarkVulcan_GithubStatic             1000000              1447 ns/op              98 B/op          3 allocs/op
BenchmarkAce_GithubParam                 2000000               686 ns/op              96 B/op          1 allocs/op
BenchmarkBear_GithubParam                1000000              2155 ns/op             496 B/op          5 allocs/op
BenchmarkBeego_GithubParam               1000000              2713 ns/op             368 B/op          4 allocs/op
BenchmarkBone_GithubParam                 100000             15088 ns/op            1760 B/op         18 allocs/op
BenchmarkDenco_GithubParam               2000000               629 ns/op             128 B/op          1 allocs/op
BenchmarkEcho_GithubParam                2000000               653 ns/op              32 B/op          1 allocs/op
BenchmarkGin_GithubParam                 5000000               255 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_GithubParam          1000000              3145 ns/op             712 B/op          9 allocs/op
BenchmarkGoji_GithubParam                1000000              1916 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_GithubParam              1000000              3975 ns/op            1024 B/op         10 allocs/op
BenchmarkGoJsonRest_GithubParam           300000              4134 ns/op             713 B/op         14 allocs/op
BenchmarkGoRestful_GithubParam             50000             30782 ns/op            2360 B/op         21 allocs/op
BenchmarkGorillaMux_GithubParam           100000             17148 ns/op            1088 B/op         11 allocs/op
BenchmarkHttpRouter_GithubParam          3000000               523 ns/op              96 B/op          1 allocs/op
BenchmarkHttpTreeMux_GithubParam         1000000              1671 ns/op             384 B/op          4 allocs/op
BenchmarkKocha_GithubParam               1000000              1021 ns/op             128 B/op          5 allocs/op
BenchmarkLARS_GithubParam                5000000               283 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GithubParam              500000              4270 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_GithubParam              100000             21728 ns/op            1152 B/op         11 allocs/op
BenchmarkPat_GithubParam                  200000             11208 ns/op            2464 B/op         48 allocs/op
BenchmarkPossum_GithubParam              1000000              2334 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_GithubParam            1000000              1487 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_GithubParam               2000000               782 ns/op              96 B/op          1 allocs/op
BenchmarkTango_GithubParam               1000000              2653 ns/op             344 B/op          8 allocs/op
BenchmarkTigerTonic_GithubParam           300000             14073 ns/op            1440 B/op         24 allocs/op
BenchmarkTraffic_GithubParam               50000             29164 ns/op            5992 B/op         52 allocs/op
BenchmarkVulcan_GithubParam              1000000              2529 ns/op              98 B/op          3 allocs/op
BenchmarkAce_GithubAll                     10000            134059 ns/op           13792 B/op        167 allocs/op
BenchmarkBear_GithubAll                     5000            534445 ns/op           86448 B/op        943 allocs/op
BenchmarkBeego_GithubAll                    3000            592444 ns/op           74705 B/op        812 allocs/op
BenchmarkBone_GithubAll                      200           6957308 ns/op          698784 B/op       8453 allocs/op
BenchmarkDenco_GithubAll                   10000            158819 ns/op           20224 B/op        167 allocs/op
BenchmarkEcho_GithubAll                    10000            154700 ns/op            6496 B/op        203 allocs/op
BenchmarkGin_GithubAll                     30000             48375 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_GithubAll               3000            570806 ns/op          131656 B/op       1686 allocs/op
BenchmarkGoji_GithubAll                     2000            818034 ns/op           56112 B/op        334 allocs/op
BenchmarkGojiv2_GithubAll                   2000           1213973 ns/op          274768 B/op       3712 allocs/op
BenchmarkGoJsonRest_GithubAll               2000            785796 ns/op          134371 B/op       2737 allocs/op
BenchmarkGoRestful_GithubAll                 300           5238188 ns/op          689672 B/op       4519 allocs/op
BenchmarkGorillaMux_GithubAll                100          10257726 ns/op          211840 B/op       2272 allocs/op
BenchmarkHttpRouter_GithubAll              20000            105414 ns/op           13792 B/op        167 allocs/op
BenchmarkHttpTreeMux_GithubAll             10000            319934 ns/op           65856 B/op        671 allocs/op
BenchmarkKocha_GithubAll                   10000            209442 ns/op           23304 B/op        843 allocs/op
BenchmarkLARS_GithubAll                    20000             62565 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GithubAll                  2000           1161270 ns/op          204194 B/op       2000 allocs/op
BenchmarkMartini_GithubAll                   200           9991713 ns/op          226549 B/op       2325 allocs/op
BenchmarkPat_GithubAll                       200           5590793 ns/op         1499568 B/op      27435 allocs/op
BenchmarkPossum_GithubAll                  10000            319768 ns/op           84448 B/op        609 allocs/op
BenchmarkR2router_GithubAll                10000            305134 ns/op           77328 B/op        979 allocs/op
BenchmarkRivet_GithubAll                   10000            132134 ns/op           16272 B/op        167 allocs/op
BenchmarkTango_GithubAll                    3000            552754 ns/op           63826 B/op       1618 allocs/op
BenchmarkTigerTonic_GithubAll               1000           1439483 ns/op          239104 B/op       5374 allocs/op
BenchmarkTraffic_GithubAll                   100          11383067 ns/op         2659329 B/op      21848 allocs/op
BenchmarkVulcan_GithubAll                   5000            394253 ns/op           19894 B/op        609 allocs/op
```

## Google+

```
BenchmarkGin_GPlusStatic                10000000               183 ns/op               0 B/op          0 allocs/op

BenchmarkAce_GPlusStatic                 5000000               276 ns/op               0 B/op          0 allocs/op
BenchmarkBear_GPlusStatic                2000000               652 ns/op             104 B/op          3 allocs/op
BenchmarkBeego_GPlusStatic               1000000              2239 ns/op             368 B/op          4 allocs/op
BenchmarkBone_GPlusStatic                5000000               380 ns/op              32 B/op          1 allocs/op
BenchmarkDenco_GPlusStatic              30000000                45.8 ns/op             0 B/op          0 allocs/op
BenchmarkEcho_GPlusStatic                5000000               338 ns/op              32 B/op          1 allocs/op
BenchmarkGocraftWeb_GPlusStatic          1000000              1158 ns/op             280 B/op          5 allocs/op
BenchmarkGoji_GPlusStatic                5000000               331 ns/op               0 B/op          0 allocs/op
BenchmarkGojiv2_GPlusStatic              1000000              2106 ns/op             928 B/op          7 allocs/op
BenchmarkGoJsonRest_GPlusStatic          1000000              1626 ns/op             329 B/op         11 allocs/op
BenchmarkGoRestful_GPlusStatic            300000              7598 ns/op            1976 B/op         20 allocs/op
BenchmarkGorillaMux_GPlusStatic          1000000              2629 ns/op             736 B/op         10 allocs/op
BenchmarkHttpRouter_GPlusStatic         30000000                52.5 ns/op             0 B/op          0 allocs/op
BenchmarkHttpTreeMux_GPlusStatic        20000000                85.8 ns/op             0 B/op          0 allocs/op
BenchmarkKocha_GPlusStatic              20000000                89.2 ns/op             0 B/op          0 allocs/op
BenchmarkLARS_GPlusStatic               10000000               162 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GPlusStatic              500000              3479 ns/op             768 B/op          9 allocs/op
BenchmarkMartini_GPlusStatic              200000              9092 ns/op             768 B/op          9 allocs/op
BenchmarkPat_GPlusStatic                 3000000               493 ns/op              96 B/op          2 allocs/op
BenchmarkPossum_GPlusStatic              1000000              1467 ns/op             416 B/op          3 allocs/op
BenchmarkR2router_GPlusStatic            2000000               788 ns/op             144 B/op          4 allocs/op
BenchmarkRivet_GPlusStatic              20000000               114 ns/op               0 B/op          0 allocs/op
BenchmarkTango_GPlusStatic               1000000              1534 ns/op             200 B/op          8 allocs/op
BenchmarkTigerTonic_GPlusStatic          5000000               282 ns/op              32 B/op          1 allocs/op
BenchmarkTraffic_GPlusStatic              500000              3798 ns/op            1192 B/op         15 allocs/op
BenchmarkVulcan_GPlusStatic              2000000              1125 ns/op              98 B/op          3 allocs/op
BenchmarkAce_GPlusParam                  3000000               528 ns/op              64 B/op          1 allocs/op
BenchmarkBear_GPlusParam                 1000000              1570 ns/op             480 B/op          5 allocs/op
BenchmarkBeego_GPlusParam                1000000              2369 ns/op             368 B/op          4 allocs/op
BenchmarkBone_GPlusParam                 1000000              2028 ns/op             688 B/op          5 allocs/op
BenchmarkDenco_GPlusParam                5000000               385 ns/op              64 B/op          1 allocs/op
BenchmarkEcho_GPlusParam                 3000000               441 ns/op              32 B/op          1 allocs/op
BenchmarkGin_GPlusParam                 10000000               174 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_GPlusParam           1000000              2033 ns/op             648 B/op          8 allocs/op
BenchmarkGoji_GPlusParam                 1000000              1399 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_GPlusParam               1000000              2641 ns/op             944 B/op          8 allocs/op
BenchmarkGoJsonRest_GPlusParam           1000000              2824 ns/op             649 B/op         13 allocs/op
BenchmarkGoRestful_GPlusParam             200000              8875 ns/op            2296 B/op         21 allocs/op
BenchmarkGorillaMux_GPlusParam            200000              6291 ns/op            1056 B/op         11 allocs/op
BenchmarkHttpRouter_GPlusParam           5000000               316 ns/op              64 B/op          1 allocs/op
BenchmarkHttpTreeMux_GPlusParam          1000000              1129 ns/op             352 B/op          3 allocs/op
BenchmarkKocha_GPlusParam                3000000               538 ns/op              56 B/op          3 allocs/op
BenchmarkLARS_GPlusParam                10000000               198 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GPlusParam               500000              3554 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_GPlusParam               200000              9831 ns/op            1072 B/op         10 allocs/op
BenchmarkPat_GPlusParam                  1000000              2706 ns/op             688 B/op         12 allocs/op
BenchmarkPossum_GPlusParam               1000000              2297 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_GPlusParam             1000000              1318 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_GPlusParam                5000000               399 ns/op              48 B/op          1 allocs/op
BenchmarkTango_GPlusParam                1000000              2070 ns/op             264 B/op          8 allocs/op
BenchmarkTigerTonic_GPlusParam            500000              4853 ns/op            1056 B/op         17 allocs/op
BenchmarkTraffic_GPlusParam               200000              8278 ns/op            1976 B/op         21 allocs/op
BenchmarkVulcan_GPlusParam               1000000              1243 ns/op              98 B/op          3 allocs/op
BenchmarkAce_GPlus2Params                3000000               549 ns/op              64 B/op          1 allocs/op
BenchmarkBear_GPlus2Params               1000000              2112 ns/op             496 B/op          5 allocs/op
BenchmarkBeego_GPlus2Params               500000              2750 ns/op             368 B/op          4 allocs/op
BenchmarkBone_GPlus2Params                300000              7032 ns/op            1040 B/op          9 allocs/op
BenchmarkDenco_GPlus2Params              3000000               502 ns/op              64 B/op          1 allocs/op
BenchmarkEcho_GPlus2Params               3000000               641 ns/op              32 B/op          1 allocs/op
BenchmarkGin_GPlus2Params                5000000               250 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_GPlus2Params         1000000              2681 ns/op             712 B/op          9 allocs/op
BenchmarkGoji_GPlus2Params               1000000              1926 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_GPlus2Params              500000              3996 ns/op            1024 B/op         11 allocs/op
BenchmarkGoJsonRest_GPlus2Params          500000              3886 ns/op             713 B/op         14 allocs/op
BenchmarkGoRestful_GPlus2Params           200000             10376 ns/op            2360 B/op         21 allocs/op
BenchmarkGorillaMux_GPlus2Params          100000             14162 ns/op            1088 B/op         11 allocs/op
BenchmarkHttpRouter_GPlus2Params         5000000               336 ns/op              64 B/op          1 allocs/op
BenchmarkHttpTreeMux_GPlus2Params        1000000              1523 ns/op             384 B/op          4 allocs/op
BenchmarkKocha_GPlus2Params              2000000               970 ns/op             128 B/op          5 allocs/op
BenchmarkLARS_GPlus2Params               5000000               238 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GPlus2Params             500000              4016 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_GPlus2Params             100000             21253 ns/op            1200 B/op         13 allocs/op
BenchmarkPat_GPlus2Params                 200000              8632 ns/op            2256 B/op         34 allocs/op
BenchmarkPossum_GPlus2Params             1000000              2171 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_GPlus2Params           1000000              1340 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_GPlus2Params              3000000               557 ns/op              96 B/op          1 allocs/op
BenchmarkTango_GPlus2Params              1000000              2186 ns/op             344 B/op          8 allocs/op
BenchmarkTigerTonic_GPlus2Params          200000              9060 ns/op            1488 B/op         24 allocs/op
BenchmarkTraffic_GPlus2Params             100000             20324 ns/op            3272 B/op         31 allocs/op
BenchmarkVulcan_GPlus2Params             1000000              2039 ns/op              98 B/op          3 allocs/op
BenchmarkAce_GPlusAll                     300000              6603 ns/op             640 B/op         11 allocs/op
BenchmarkBear_GPlusAll                    100000             22363 ns/op            5488 B/op         61 allocs/op
BenchmarkBeego_GPlusAll                    50000             38757 ns/op            4784 B/op         52 allocs/op
BenchmarkBone_GPlusAll                     20000             54916 ns/op           10336 B/op         98 allocs/op
BenchmarkDenco_GPlusAll                   300000              4959 ns/op             672 B/op         11 allocs/op
BenchmarkEcho_GPlusAll                    200000              6558 ns/op             416 B/op         13 allocs/op
BenchmarkGin_GPlusAll                     500000              2757 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_GPlusAll               50000             34615 ns/op            8040 B/op        103 allocs/op
BenchmarkGoji_GPlusAll                    100000             16002 ns/op            3696 B/op         22 allocs/op
BenchmarkGojiv2_GPlusAll                   50000             35060 ns/op           12624 B/op        115 allocs/op
BenchmarkGoJsonRest_GPlusAll               50000             41479 ns/op            8117 B/op        170 allocs/op
BenchmarkGoRestful_GPlusAll                10000            131653 ns/op           32024 B/op        275 allocs/op
BenchmarkGorillaMux_GPlusAll               10000            101380 ns/op           13296 B/op        142 allocs/op
BenchmarkHttpRouter_GPlusAll              500000              3711 ns/op             640 B/op         11 allocs/op
BenchmarkHttpTreeMux_GPlusAll             100000             14438 ns/op            4032 B/op         38 allocs/op
BenchmarkKocha_GPlusAll                   200000              8039 ns/op             976 B/op         43 allocs/op
BenchmarkLARS_GPlusAll                    500000              2630 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_GPlusAll                  30000             51123 ns/op           13152 B/op        128 allocs/op
BenchmarkMartini_GPlusAll                  10000            176157 ns/op           14016 B/op        145 allocs/op
BenchmarkPat_GPlusAll                      20000             69911 ns/op           16576 B/op        298 allocs/op
BenchmarkPossum_GPlusAll                  100000             20716 ns/op            5408 B/op         39 allocs/op
BenchmarkR2router_GPlusAll                100000             17463 ns/op            5040 B/op         63 allocs/op
BenchmarkRivet_GPlusAll                   300000              5142 ns/op             768 B/op         11 allocs/op
BenchmarkTango_GPlusAll                    50000             27321 ns/op            3656 B/op        104 allocs/op
BenchmarkTigerTonic_GPlusAll               20000             77597 ns/op           14512 B/op        288 allocs/op
BenchmarkTraffic_GPlusAll                  10000            151406 ns/op           37360 B/op        392 allocs/op
BenchmarkVulcan_GPlusAll                  100000             18555 ns/op            1274 B/op         39 allocs/op
```

## Parse.com

```
BenchmarkGin_ParseStatic                10000000               133 ns/op               0 B/op          0 allocs/op

BenchmarkAce_ParseStatic                 5000000               241 ns/op               0 B/op          0 allocs/op
BenchmarkBear_ParseStatic                2000000               728 ns/op             120 B/op          3 allocs/op
BenchmarkBeego_ParseStatic               1000000              2623 ns/op             368 B/op          4 allocs/op
BenchmarkBone_ParseStatic                1000000              1285 ns/op             144 B/op          3 allocs/op
BenchmarkDenco_ParseStatic              30000000                57.8 ns/op             0 B/op          0 allocs/op
BenchmarkEcho_ParseStatic                5000000               342 ns/op              32 B/op          1 allocs/op
BenchmarkGocraftWeb_ParseStatic          1000000              1478 ns/op             296 B/op          5 allocs/op
BenchmarkGoji_ParseStatic                3000000               415 ns/op               0 B/op          0 allocs/op
BenchmarkGojiv2_ParseStatic              1000000              2087 ns/op             928 B/op          7 allocs/op
BenchmarkGoJsonRest_ParseStatic          1000000              1712 ns/op             329 B/op         11 allocs/op
BenchmarkGoRestful_ParseStatic            200000             11072 ns/op            3224 B/op         22 allocs/op
BenchmarkGorillaMux_ParseStatic           500000              4129 ns/op             752 B/op         11 allocs/op
BenchmarkHttpRouter_ParseStatic         30000000                52.4 ns/op             0 B/op          0 allocs/op
BenchmarkHttpTreeMux_ParseStatic        20000000               109 ns/op               0 B/op          0 allocs/op
BenchmarkKocha_ParseStatic              20000000                81.8 ns/op             0 B/op          0 allocs/op
BenchmarkLARS_ParseStatic               10000000               150 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_ParseStatic             1000000              3288 ns/op             768 B/op          9 allocs/op
BenchmarkMartini_ParseStatic              200000              9110 ns/op             768 B/op          9 allocs/op
BenchmarkPat_ParseStatic                 1000000              1135 ns/op             240 B/op          5 allocs/op
BenchmarkPossum_ParseStatic              1000000              1557 ns/op             416 B/op          3 allocs/op
BenchmarkR2router_ParseStatic            2000000               730 ns/op             144 B/op          4 allocs/op
BenchmarkRivet_ParseStatic              10000000               121 ns/op               0 B/op          0 allocs/op
BenchmarkTango_ParseStatic               1000000              1688 ns/op             248 B/op          8 allocs/op
BenchmarkTigerTonic_ParseStatic          3000000               427 ns/op              48 B/op          1 allocs/op
BenchmarkTraffic_ParseStatic              500000              5962 ns/op            1816 B/op         20 allocs/op
BenchmarkVulcan_ParseStatic              2000000               969 ns/op              98 B/op          3 allocs/op
BenchmarkAce_ParseParam                  3000000               497 ns/op              64 B/op          1 allocs/op
BenchmarkBear_ParseParam                 1000000              1473 ns/op             467 B/op          5 allocs/op
BenchmarkBeego_ParseParam                1000000              2384 ns/op             368 B/op          4 allocs/op
BenchmarkBone_ParseParam                 1000000              2513 ns/op             768 B/op          6 allocs/op
BenchmarkDenco_ParseParam                5000000               364 ns/op              64 B/op          1 allocs/op
BenchmarkEcho_ParseParam                 5000000               418 ns/op              32 B/op          1 allocs/op
BenchmarkGin_ParseParam                 10000000               163 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_ParseParam           1000000              2361 ns/op             664 B/op          8 allocs/op
BenchmarkGoji_ParseParam                 1000000              1590 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_ParseParam               1000000              2851 ns/op             976 B/op          9 allocs/op
BenchmarkGoJsonRest_ParseParam           1000000              2965 ns/op             649 B/op         13 allocs/op
BenchmarkGoRestful_ParseParam             200000             12207 ns/op            3544 B/op         23 allocs/op
BenchmarkGorillaMux_ParseParam            500000              5187 ns/op            1088 B/op         12 allocs/op
BenchmarkHttpRouter_ParseParam           5000000               275 ns/op              64 B/op          1 allocs/op
BenchmarkHttpTreeMux_ParseParam          1000000              1108 ns/op             352 B/op          3 allocs/op
BenchmarkKocha_ParseParam                3000000               495 ns/op              56 B/op          3 allocs/op
BenchmarkLARS_ParseParam                10000000               192 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_ParseParam               500000              4103 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_ParseParam               200000              9878 ns/op            1072 B/op         10 allocs/op
BenchmarkPat_ParseParam                   500000              3657 ns/op            1120 B/op         17 allocs/op
BenchmarkPossum_ParseParam               1000000              2084 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_ParseParam             1000000              1251 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_ParseParam                5000000               335 ns/op              48 B/op          1 allocs/op
BenchmarkTango_ParseParam                1000000              1854 ns/op             280 B/op          8 allocs/op
BenchmarkTigerTonic_ParseParam            500000              4582 ns/op            1008 B/op         17 allocs/op
BenchmarkTraffic_ParseParam               200000              8125 ns/op            2248 B/op         23 allocs/op
BenchmarkVulcan_ParseParam               1000000              1148 ns/op              98 B/op          3 allocs/op
BenchmarkAce_Parse2Params                3000000               539 ns/op              64 B/op          1 allocs/op
BenchmarkBear_Parse2Params               1000000              1778 ns/op             496 B/op          5 allocs/op
BenchmarkBeego_Parse2Params              1000000              2519 ns/op             368 B/op          4 allocs/op
BenchmarkBone_Parse2Params               1000000              2596 ns/op             720 B/op          5 allocs/op
BenchmarkDenco_Parse2Params              3000000               492 ns/op              64 B/op          1 allocs/op
BenchmarkEcho_Parse2Params               3000000               484 ns/op              32 B/op          1 allocs/op
BenchmarkGin_Parse2Params               10000000               193 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_Parse2Params         1000000              2575 ns/op             712 B/op          9 allocs/op
BenchmarkGoji_Parse2Params               1000000              1373 ns/op             336 B/op          2 allocs/op
BenchmarkGojiv2_Parse2Params              500000              2416 ns/op             960 B/op          8 allocs/op
BenchmarkGoJsonRest_Parse2Params          300000              3452 ns/op             713 B/op         14 allocs/op
BenchmarkGoRestful_Parse2Params           100000             17719 ns/op            6008 B/op         25 allocs/op
BenchmarkGorillaMux_Parse2Params          300000              5102 ns/op            1088 B/op         11 allocs/op
BenchmarkHttpRouter_Parse2Params         5000000               303 ns/op              64 B/op          1 allocs/op
BenchmarkHttpTreeMux_Parse2Params        1000000              1372 ns/op             384 B/op          4 allocs/op
BenchmarkKocha_Parse2Params              2000000               874 ns/op             128 B/op          5 allocs/op
BenchmarkLARS_Parse2Params              10000000               192 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_Parse2Params             500000              3871 ns/op            1056 B/op         10 allocs/op
BenchmarkMartini_Parse2Params             200000              9954 ns/op            1152 B/op         11 allocs/op
BenchmarkPat_Parse2Params                 500000              4194 ns/op             832 B/op         17 allocs/op
BenchmarkPossum_Parse2Params             1000000              2121 ns/op             560 B/op          6 allocs/op
BenchmarkR2router_Parse2Params           1000000              1415 ns/op             432 B/op          5 allocs/op
BenchmarkRivet_Parse2Params              3000000               457 ns/op              96 B/op          1 allocs/op
BenchmarkTango_Parse2Params              1000000              1914 ns/op             312 B/op          8 allocs/op
BenchmarkTigerTonic_Parse2Params          300000              6895 ns/op            1408 B/op         24 allocs/op
BenchmarkTraffic_Parse2Params             200000              8317 ns/op            2040 B/op         22 allocs/op
BenchmarkVulcan_Parse2Params             1000000              1274 ns/op              98 B/op          3 allocs/op
BenchmarkAce_ParseAll                     200000             10401 ns/op             640 B/op         16 allocs/op
BenchmarkBear_ParseAll                     50000             37743 ns/op            8928 B/op        110 allocs/op
BenchmarkBeego_ParseAll                    20000             63193 ns/op            9568 B/op        104 allocs/op
BenchmarkBone_ParseAll                     20000             61767 ns/op           14160 B/op        131 allocs/op
BenchmarkDenco_ParseAll                   300000              7036 ns/op             928 B/op         16 allocs/op
BenchmarkEcho_ParseAll                    200000             11824 ns/op             832 B/op         26 allocs/op
BenchmarkGin_ParseAll                     300000              4199 ns/op               0 B/op          0 allocs/op
BenchmarkGocraftWeb_ParseAll               30000             51758 ns/op           13728 B/op        181 allocs/op
BenchmarkGoji_ParseAll                     50000             29614 ns/op            5376 B/op         32 allocs/op
BenchmarkGojiv2_ParseAll                   20000             68676 ns/op           24464 B/op        199 allocs/op
BenchmarkGoJsonRest_ParseAll               20000             76135 ns/op           13866 B/op        321 allocs/op
BenchmarkGoRestful_ParseAll                 5000            389487 ns/op          110928 B/op        600 allocs/op
BenchmarkGorillaMux_ParseAll               10000            221250 ns/op           24864 B/op        292 allocs/op
BenchmarkHttpRouter_ParseAll              200000              6444 ns/op             640 B/op         16 allocs/op
BenchmarkHttpTreeMux_ParseAll              50000             30702 ns/op            5728 B/op         51 allocs/op
BenchmarkKocha_ParseAll                   200000             13712 ns/op            1112 B/op         54 allocs/op
BenchmarkLARS_ParseAll                    300000              6925 ns/op               0 B/op          0 allocs/op
BenchmarkMacaron_ParseAll                  20000             96278 ns/op           24576 B/op        250 allocs/op
BenchmarkMartini_ParseAll                   5000            271352 ns/op           25072 B/op        253 allocs/op
BenchmarkPat_ParseAll                      20000             74941 ns/op           17264 B/op        343 allocs/op
BenchmarkPossum_ParseAll                   50000             39947 ns/op           10816 B/op         78 allocs/op
BenchmarkR2router_ParseAll                 50000             42479 ns/op            8352 B/op        120 allocs/op
BenchmarkRivet_ParseAll                   200000              7726 ns/op             912 B/op         16 allocs/op
BenchmarkTango_ParseAll                    30000             50014 ns/op            7168 B/op        208 allocs/op
BenchmarkTigerTonic_ParseAll               10000            106550 ns/op           19728 B/op        379 allocs/op
BenchmarkTraffic_ParseAll                  10000            216037 ns/op           57776 B/op        642 allocs/op
BenchmarkVulcan_ParseAll                   50000             34379 ns/op            2548 B/op         78 allocs/op
```
