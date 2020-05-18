
# Benchmark System

**VM HOST:** Travis  
**Machine:** Ubuntu 16.04.6 LTS x64  
**Date:** May 04th, 2020  
**Version:** Gin v1.6.3
**Go Version:** 1.14.2 linux/amd64  
**Source:** [Go HTTP Router Benchmark](https://github.com/gin-gonic/go-http-routing-benchmark)
**Result:** [See the gist](https://gist.github.com/appleboy/b5f2ecfaf50824ae9c64dcfb9165ae5e) or [Travis result](https://travis-ci.org/github/gin-gonic/go-http-routing-benchmark/jobs/682947061)

## Static Routes: 157

```sh
Gin: 34936 Bytes

HttpServeMux: 14512 Bytes
Ace: 30680 Bytes
Aero: 34536 Bytes
Bear: 30456 Bytes
Beego: 98456 Bytes
Bone: 40224 Bytes
Chi: 83608 Bytes
Denco: 10216 Bytes
Echo: 80328 Bytes
GocraftWeb: 55288 Bytes
Goji: 29744 Bytes
Gojiv2: 105840 Bytes
GoJsonRest: 137496 Bytes
GoRestful: 816936 Bytes
GorillaMux: 585632 Bytes
GowwwRouter: 24968 Bytes
HttpRouter: 21712 Bytes
HttpTreeMux: 73448 Bytes
Kocha: 115472 Bytes
LARS: 30640 Bytes
Macaron: 38592 Bytes
Martini: 310864 Bytes
Pat: 19696 Bytes
Possum: 89920 Bytes
R2router: 23712 Bytes
Rivet: 24608 Bytes
Tango: 28264 Bytes
TigerTonic: 78768 Bytes
Traffic: 538976 Bytes
Vulcan: 369960 Bytes
```

## GithubAPI Routes: 203

```sh
Gin: 58512 Bytes

Ace: 48688 Bytes
Aero: 318568 Bytes
Bear: 84248 Bytes
Beego: 150936 Bytes
Bone: 100976 Bytes
Chi: 95112 Bytes
Denco: 36736 Bytes
Echo: 100296 Bytes
GocraftWeb: 95432 Bytes
Goji: 49680 Bytes
Gojiv2: 104704 Bytes
GoJsonRest: 141976 Bytes
GoRestful: 1241656 Bytes
GorillaMux: 1322784 Bytes
GowwwRouter: 80008 Bytes
HttpRouter: 37144 Bytes
HttpTreeMux: 78800 Bytes
Kocha: 785120 Bytes
LARS: 48600 Bytes
Macaron: 92784 Bytes
Martini: 485264 Bytes
Pat: 21200 Bytes
Possum: 85312 Bytes
R2router: 47104 Bytes
Rivet: 42840 Bytes
Tango: 54840 Bytes
TigerTonic: 95264 Bytes
Traffic: 921744 Bytes
Vulcan: 425992 Bytes
```

## GPlusAPI Routes: 13

```sh
Gin: 4384 Bytes

Ace: 3712 Bytes
Aero: 26056 Bytes
Bear: 7112 Bytes
Beego: 10272 Bytes
Bone: 6688 Bytes
Chi: 8024 Bytes
Denco: 3264 Bytes
Echo: 9688 Bytes
GocraftWeb: 7496 Bytes
Goji: 3152 Bytes
Gojiv2: 7376 Bytes
GoJsonRest: 11400 Bytes
GoRestful: 74328 Bytes
GorillaMux: 66208 Bytes
GowwwRouter: 5744 Bytes
HttpRouter: 2808 Bytes
HttpTreeMux: 7440 Bytes
Kocha: 128880 Bytes
LARS: 3656 Bytes
Macaron: 8656 Bytes
Martini: 23920 Bytes
Pat: 1856 Bytes
Possum: 7248 Bytes
R2router: 3928 Bytes
Rivet: 3064 Bytes
Tango: 5168 Bytes
TigerTonic: 9408 Bytes
Traffic: 46400 Bytes
Vulcan: 25544 Bytes
```

## ParseAPI Routes: 26

```sh
Gin: 7776 Bytes

Ace: 6704 Bytes
Aero: 28488 Bytes
Bear: 12320 Bytes
Beego: 19280 Bytes
Bone: 11440 Bytes
Chi: 9744 Bytes
Denco: 4192 Bytes
Echo: 11664 Bytes
GocraftWeb: 12800 Bytes
Goji: 5680 Bytes
Gojiv2: 14464 Bytes
GoJsonRest: 14072 Bytes
GoRestful: 116264 Bytes
GorillaMux: 105880 Bytes
GowwwRouter: 9344 Bytes
HttpRouter: 5072 Bytes
HttpTreeMux: 7848 Bytes
Kocha: 181712 Bytes
LARS: 6632 Bytes
Macaron: 13648 Bytes
Martini: 45888 Bytes
Pat: 2560 Bytes
Possum: 9200 Bytes
R2router: 7056 Bytes
Rivet: 5680 Bytes
Tango: 8920 Bytes
TigerTonic: 9840 Bytes
Traffic: 79096 Bytes
Vulcan: 44504 Bytes
```

## Static Routes

```sh
BenchmarkGin_StaticAll                   62169         19319 ns/op           0 B/op           0 allocs/op

BenchmarkAce_StaticAll                   65428         18313 ns/op           0 B/op           0 allocs/op
BenchmarkAero_StaticAll                 121132          9632 ns/op           0 B/op           0 allocs/op
BenchmarkHttpServeMux_StaticAll          52626         22758 ns/op           0 B/op           0 allocs/op
BenchmarkBeego_StaticAll                  9962        179058 ns/op       55264 B/op         471 allocs/op
BenchmarkBear_StaticAll                  14894         80966 ns/op       20272 B/op         469 allocs/op
BenchmarkBone_StaticAll                  18718         64065 ns/op           0 B/op           0 allocs/op
BenchmarkChi_StaticAll                   10000        149827 ns/op       67824 B/op         471 allocs/op
BenchmarkDenco_StaticAll                211393          5680 ns/op           0 B/op           0 allocs/op
BenchmarkEcho_StaticAll                  49341         24343 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_StaticAll            10000        126209 ns/op       46312 B/op         785 allocs/op
BenchmarkGoji_StaticAll                  27956         43174 ns/op           0 B/op           0 allocs/op
BenchmarkGojiv2_StaticAll                 3430        370718 ns/op      205984 B/op        1570 allocs/op
BenchmarkGoJsonRest_StaticAll             9134        188888 ns/op       51653 B/op        1727 allocs/op
BenchmarkGoRestful_StaticAll               706       1703330 ns/op      613280 B/op        2053 allocs/op
BenchmarkGorillaMux_StaticAll             1268        924083 ns/op      153233 B/op        1413 allocs/op
BenchmarkGowwwRouter_StaticAll           63374         18935 ns/op           0 B/op           0 allocs/op
BenchmarkHttpRouter_StaticAll           109938         10902 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_StaticAll          109166         10861 ns/op           0 B/op           0 allocs/op
BenchmarkKocha_StaticAll                 92258         12992 ns/op           0 B/op           0 allocs/op
BenchmarkLARS_StaticAll                  65200         18387 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_StaticAll                5671        291501 ns/op      115553 B/op        1256 allocs/op
BenchmarkMartini_StaticAll                 807       1460498 ns/op      125444 B/op        1717 allocs/op
BenchmarkPat_StaticAll                     513       2342396 ns/op      602832 B/op       12559 allocs/op
BenchmarkPossum_StaticAll                10000        128270 ns/op       65312 B/op         471 allocs/op
BenchmarkR2router_StaticAll              16726         71760 ns/op       22608 B/op         628 allocs/op
BenchmarkRivet_StaticAll                 41722         28723 ns/op           0 B/op           0 allocs/op
BenchmarkTango_StaticAll                  7606        205082 ns/op       39209 B/op        1256 allocs/op
BenchmarkTigerTonic_StaticAll            26247         45806 ns/op        7376 B/op         157 allocs/op
BenchmarkTraffic_StaticAll                 550       2284518 ns/op      754864 B/op       14601 allocs/op
BenchmarkVulcan_StaticAll                10000        131343 ns/op       15386 B/op         471 allocs/op
```

## Micro Benchmarks

```sh
BenchmarkGin_Param                    18785022          63.9 ns/op           0 B/op           0 allocs/op

BenchmarkAce_Param                    14689765          81.5 ns/op           0 B/op           0 allocs/op
BenchmarkAero_Param                   23094770          51.2 ns/op           0 B/op           0 allocs/op
BenchmarkBear_Param                    1417045           845 ns/op         456 B/op           5 allocs/op
BenchmarkBeego_Param                   1000000          1080 ns/op         352 B/op           3 allocs/op
BenchmarkBone_Param                    1000000          1463 ns/op         816 B/op           6 allocs/op
BenchmarkChi_Param                     1378756           885 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_Param                   8557899           143 ns/op          32 B/op           1 allocs/op
BenchmarkEcho_Param                   16433347          75.5 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_Param              1000000          1218 ns/op         648 B/op           8 allocs/op
BenchmarkGoji_Param                    1921248           617 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_Param                   561848          2156 ns/op        1328 B/op          11 allocs/op
BenchmarkGoJsonRest_Param              1000000          1358 ns/op         649 B/op          13 allocs/op
BenchmarkGoRestful_Param                224857          5307 ns/op        4192 B/op          14 allocs/op
BenchmarkGorillaMux_Param               498313          2459 ns/op        1280 B/op          10 allocs/op
BenchmarkGowwwRouter_Param             1864354           654 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_Param             26269074          47.7 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_Param             2109829           557 ns/op         352 B/op           3 allocs/op
BenchmarkKocha_Param                   5050216           243 ns/op          56 B/op           3 allocs/op
BenchmarkLARS_Param                   19811712          59.9 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_Param                  662746          2329 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_Param                  279902          4260 ns/op        1072 B/op          10 allocs/op
BenchmarkPat_Param                     1000000          1382 ns/op         536 B/op          11 allocs/op
BenchmarkPossum_Param                  1000000          1014 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_Param                1712559           707 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_Param                   6648086           182 ns/op          48 B/op           1 allocs/op
BenchmarkTango_Param                   1221504           994 ns/op         248 B/op           8 allocs/op
BenchmarkTigerTonic_Param               891661          2261 ns/op         776 B/op          16 allocs/op
BenchmarkTraffic_Param                  350059          3598 ns/op        1856 B/op          21 allocs/op
BenchmarkVulcan_Param                  2517823           472 ns/op          98 B/op           3 allocs/op
BenchmarkAce_Param5                    9214365           130 ns/op           0 B/op           0 allocs/op
BenchmarkAero_Param5                  15369013          77.9 ns/op           0 B/op           0 allocs/op
BenchmarkBear_Param5                   1000000          1113 ns/op         501 B/op           5 allocs/op
BenchmarkBeego_Param5                  1000000          1269 ns/op         352 B/op           3 allocs/op
BenchmarkBone_Param5                    986820          1873 ns/op         864 B/op           6 allocs/op
BenchmarkChi_Param5                    1000000          1156 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_Param5                  3036331           400 ns/op         160 B/op           1 allocs/op
BenchmarkEcho_Param5                   6447133           186 ns/op           0 B/op           0 allocs/op
BenchmarkGin_Param5                   10786068           110 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_Param5              844820          1944 ns/op         920 B/op          11 allocs/op
BenchmarkGoji_Param5                   1474965           827 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_Param5                  442820          2516 ns/op        1392 B/op          11 allocs/op
BenchmarkGoJsonRest_Param5              507555          2711 ns/op        1097 B/op          16 allocs/op
BenchmarkGoRestful_Param5               216481          6093 ns/op        4288 B/op          14 allocs/op
BenchmarkGorillaMux_Param5              314402          3628 ns/op        1344 B/op          10 allocs/op
BenchmarkGowwwRouter_Param5            1624660           733 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_Param5            13167324          92.0 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_Param5            1000000          1295 ns/op         576 B/op           6 allocs/op
BenchmarkKocha_Param5                  1000000          1138 ns/op         440 B/op          10 allocs/op
BenchmarkLARS_Param5                  11580613           105 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_Param5                 473596          2755 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_Param5                 230756          5111 ns/op        1232 B/op          11 allocs/op
BenchmarkPat_Param5                     469190          3370 ns/op         888 B/op          29 allocs/op
BenchmarkPossum_Param5                 1000000          1002 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_Param5               1422129           844 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_Param5                  2263789           539 ns/op         240 B/op           1 allocs/op
BenchmarkTango_Param5                  1000000          1256 ns/op         360 B/op           8 allocs/op
BenchmarkTigerTonic_Param5              175500          7492 ns/op        2279 B/op          39 allocs/op
BenchmarkTraffic_Param5                 233631          5816 ns/op        2208 B/op          27 allocs/op
BenchmarkVulcan_Param5                 1923416           629 ns/op          98 B/op           3 allocs/op
BenchmarkAce_Param20                   4321266           281 ns/op           0 B/op           0 allocs/op
BenchmarkAero_Param20                 31501641          35.2 ns/op           0 B/op           0 allocs/op
BenchmarkBear_Param20                   335204          3489 ns/op        1665 B/op           5 allocs/op
BenchmarkBeego_Param20                  503674          2860 ns/op         352 B/op           3 allocs/op
BenchmarkBone_Param20                   298922          4741 ns/op        2031 B/op           6 allocs/op
BenchmarkChi_Param20                    878181          1957 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_Param20                 1000000          1360 ns/op         640 B/op           1 allocs/op
BenchmarkEcho_Param20                  2104946           580 ns/op           0 B/op           0 allocs/op
BenchmarkGin_Param20                   4167204           290 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_Param20             173064          7514 ns/op        3796 B/op          15 allocs/op
BenchmarkGoji_Param20                   458778          2651 ns/op        1247 B/op           2 allocs/op
BenchmarkGojiv2_Param20                 364862          3178 ns/op        1632 B/op          11 allocs/op
BenchmarkGoJsonRest_Param20             125514          9760 ns/op        4485 B/op          20 allocs/op
BenchmarkGoRestful_Param20              101217         11964 ns/op        6715 B/op          18 allocs/op
BenchmarkGorillaMux_Param20             147654          8132 ns/op        3452 B/op          12 allocs/op
BenchmarkGowwwRouter_Param20           1000000          1225 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_Param20            4920895           247 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_Param20            173202          6605 ns/op        3196 B/op          10 allocs/op
BenchmarkKocha_Param20                  345988          3620 ns/op        1808 B/op          27 allocs/op
BenchmarkLARS_Param20                  4592326           262 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_Param20                166492          7286 ns/op        2924 B/op          12 allocs/op
BenchmarkMartini_Param20                122162         10653 ns/op        3595 B/op          13 allocs/op
BenchmarkPat_Param20                     78630         15239 ns/op        4424 B/op          93 allocs/op
BenchmarkPossum_Param20                1000000          1008 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_Param20               294981          4587 ns/op        2284 B/op           7 allocs/op
BenchmarkRivet_Param20                  691798          2090 ns/op        1024 B/op           1 allocs/op
BenchmarkTango_Param20                  842440          2505 ns/op         856 B/op           8 allocs/op
BenchmarkTigerTonic_Param20              38614         31509 ns/op        9870 B/op         119 allocs/op
BenchmarkTraffic_Param20                 57633         21107 ns/op        7853 B/op          47 allocs/op
BenchmarkVulcan_Param20                1000000          1178 ns/op          98 B/op           3 allocs/op
BenchmarkAce_ParamWrite                7330743           180 ns/op           8 B/op           1 allocs/op
BenchmarkAero_ParamWrite              13833598          86.7 ns/op           0 B/op           0 allocs/op
BenchmarkBear_ParamWrite               1363321           867 ns/op         456 B/op           5 allocs/op
BenchmarkBeego_ParamWrite              1000000          1104 ns/op         360 B/op           4 allocs/op
BenchmarkBone_ParamWrite               1000000          1475 ns/op         816 B/op           6 allocs/op
BenchmarkChi_ParamWrite                1320590           892 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_ParamWrite              7093605           172 ns/op          32 B/op           1 allocs/op
BenchmarkEcho_ParamWrite               8434424           161 ns/op           8 B/op           1 allocs/op
BenchmarkGin_ParamWrite               10377034           118 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_ParamWrite         1000000          1266 ns/op         656 B/op           9 allocs/op
BenchmarkGoji_ParamWrite               1874168           654 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_ParamWrite              459032          2352 ns/op        1360 B/op          13 allocs/op
BenchmarkGoJsonRest_ParamWrite          499434          2145 ns/op        1128 B/op          18 allocs/op
BenchmarkGoRestful_ParamWrite           241087          5470 ns/op        4200 B/op          15 allocs/op
BenchmarkGorillaMux_ParamWrite          425686          2522 ns/op        1280 B/op          10 allocs/op
BenchmarkGowwwRouter_ParamWrite         922172          1778 ns/op         976 B/op           8 allocs/op
BenchmarkHttpRouter_ParamWrite        15392049          77.7 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_ParamWrite        1973385           597 ns/op         352 B/op           3 allocs/op
BenchmarkKocha_ParamWrite              4262500           281 ns/op          56 B/op           3 allocs/op
BenchmarkLARS_ParamWrite              10764410           113 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_ParamWrite             486769          2726 ns/op        1176 B/op          14 allocs/op
BenchmarkMartini_ParamWrite             264804          4842 ns/op        1176 B/op          14 allocs/op
BenchmarkPat_ParamWrite                 735116          2047 ns/op         960 B/op          15 allocs/op
BenchmarkPossum_ParamWrite             1000000          1004 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_ParamWrite           1592136           768 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_ParamWrite              3582051           339 ns/op         112 B/op           2 allocs/op
BenchmarkTango_ParamWrite              2237337           534 ns/op         136 B/op           4 allocs/op
BenchmarkTigerTonic_ParamWrite          439608          3136 ns/op        1216 B/op          21 allocs/op
BenchmarkTraffic_ParamWrite             306979          4328 ns/op        2280 B/op          25 allocs/op
BenchmarkVulcan_ParamWrite             2529973           472 ns/op          98 B/op           3 allocs/op
```

## GitHub

```sh
BenchmarkGin_GithubStatic             15629472          76.7 ns/op           0 B/op           0 allocs/op

BenchmarkAce_GithubStatic             15542612          75.9 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GithubStatic            24777151          48.5 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GithubStatic             2788894           435 ns/op         120 B/op           3 allocs/op
BenchmarkBeego_GithubStatic            1000000          1064 ns/op         352 B/op           3 allocs/op
BenchmarkBone_GithubStatic               93507         12838 ns/op        2880 B/op          60 allocs/op
BenchmarkChi_GithubStatic              1387743           860 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_GithubStatic           39384996          30.4 ns/op           0 B/op           0 allocs/op
BenchmarkEcho_GithubStatic            12076382          99.1 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GithubStatic       1596495           756 ns/op         296 B/op           5 allocs/op
BenchmarkGoji_GithubStatic             6364876           189 ns/op           0 B/op           0 allocs/op
BenchmarkGojiv2_GithubStatic            550202          2098 ns/op        1312 B/op          10 allocs/op
BenchmarkGoRestful_GithubStatic         102183         12552 ns/op        4256 B/op          13 allocs/op
BenchmarkGoJsonRest_GithubStatic       1000000          1029 ns/op         329 B/op          11 allocs/op
BenchmarkGorillaMux_GithubStatic        255552          5190 ns/op         976 B/op           9 allocs/op
BenchmarkGowwwRouter_GithubStatic     15531916          77.1 ns/op           0 B/op           0 allocs/op
BenchmarkHttpRouter_GithubStatic      27920724          43.1 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GithubStatic     21448953          55.8 ns/op           0 B/op           0 allocs/op
BenchmarkKocha_GithubStatic           21405310          56.0 ns/op           0 B/op           0 allocs/op
BenchmarkLARS_GithubStatic            13625156          89.0 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GithubStatic          1000000          1747 ns/op         736 B/op           8 allocs/op
BenchmarkMartini_GithubStatic           187186          7326 ns/op         768 B/op           9 allocs/op
BenchmarkPat_GithubStatic               109143         11563 ns/op        3648 B/op          76 allocs/op
BenchmarkPossum_GithubStatic           1575898           770 ns/op         416 B/op           3 allocs/op
BenchmarkR2router_GithubStatic         3046231           404 ns/op         144 B/op           4 allocs/op
BenchmarkRivet_GithubStatic           11484826           105 ns/op           0 B/op           0 allocs/op
BenchmarkTango_GithubStatic            1000000          1153 ns/op         248 B/op           8 allocs/op
BenchmarkTigerTonic_GithubStatic       4929780           249 ns/op          48 B/op           1 allocs/op
BenchmarkTraffic_GithubStatic           106351         11819 ns/op        4664 B/op          90 allocs/op
BenchmarkVulcan_GithubStatic           1613271           722 ns/op          98 B/op           3 allocs/op
BenchmarkAce_GithubParam               8386032           143 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GithubParam             11816200           102 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GithubParam              1000000          1012 ns/op         496 B/op           5 allocs/op
BenchmarkBeego_GithubParam             1000000          1157 ns/op         352 B/op           3 allocs/op
BenchmarkBone_GithubParam               184653          6912 ns/op        1888 B/op          19 allocs/op
BenchmarkChi_GithubParam               1000000          1102 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_GithubParam             3484798           352 ns/op         128 B/op           1 allocs/op
BenchmarkEcho_GithubParam              6337380           189 ns/op           0 B/op           0 allocs/op
BenchmarkGin_GithubParam               9132032           131 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GithubParam        1000000          1446 ns/op         712 B/op           9 allocs/op
BenchmarkGoji_GithubParam              1248640           977 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_GithubParam             383233          2784 ns/op        1408 B/op          13 allocs/op
BenchmarkGoJsonRest_GithubParam        1000000          1991 ns/op         713 B/op          14 allocs/op
BenchmarkGoRestful_GithubParam           76414         16015 ns/op        4352 B/op          16 allocs/op
BenchmarkGorillaMux_GithubParam         150026          7663 ns/op        1296 B/op          10 allocs/op
BenchmarkGowwwRouter_GithubParam       1592044           751 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_GithubParam       10420628           115 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GithubParam       1403755           835 ns/op         384 B/op           4 allocs/op
BenchmarkKocha_GithubParam             2286170           533 ns/op         128 B/op           5 allocs/op
BenchmarkLARS_GithubParam              9540374           129 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GithubParam            533154          2742 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_GithubParam            119397          9638 ns/op        1152 B/op          11 allocs/op
BenchmarkPat_GithubParam                150675          8858 ns/op        2408 B/op          48 allocs/op
BenchmarkPossum_GithubParam            1000000          1001 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_GithubParam          1602886           761 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_GithubParam             2986579           409 ns/op          96 B/op           1 allocs/op
BenchmarkTango_GithubParam             1000000          1356 ns/op         344 B/op           8 allocs/op
BenchmarkTigerTonic_GithubParam         388899          3429 ns/op        1176 B/op          22 allocs/op
BenchmarkTraffic_GithubParam            123160          9734 ns/op        2816 B/op          40 allocs/op
BenchmarkVulcan_GithubParam            1000000          1138 ns/op          98 B/op           3 allocs/op
BenchmarkAce_GithubAll                   40543         29670 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GithubAll                  57632         20648 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GithubAll                   9234        216179 ns/op       86448 B/op         943 allocs/op
BenchmarkBeego_GithubAll                  7407        243496 ns/op       71456 B/op         609 allocs/op
BenchmarkBone_GithubAll                    420       2922835 ns/op      720160 B/op        8620 allocs/op
BenchmarkChi_GithubAll                    7620        238331 ns/op       87696 B/op         609 allocs/op
BenchmarkDenco_GithubAll                 18355         64494 ns/op       20224 B/op         167 allocs/op
BenchmarkEcho_GithubAll                  31251         38479 ns/op           0 B/op           0 allocs/op
BenchmarkGin_GithubAll                   43550         27364 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GithubAll             4117        300062 ns/op      131656 B/op        1686 allocs/op
BenchmarkGoji_GithubAll                   3274        416158 ns/op       56112 B/op         334 allocs/op
BenchmarkGojiv2_GithubAll                 1402        870518 ns/op      352720 B/op        4321 allocs/op
BenchmarkGoJsonRest_GithubAll             2976        401507 ns/op      134371 B/op        2737 allocs/op
BenchmarkGoRestful_GithubAll               410       2913158 ns/op      910144 B/op        2938 allocs/op
BenchmarkGorillaMux_GithubAll              346       3384987 ns/op      251650 B/op        1994 allocs/op
BenchmarkGowwwRouter_GithubAll           10000        143025 ns/op       72144 B/op         501 allocs/op
BenchmarkHttpRouter_GithubAll            55938         21360 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GithubAll           10000        153944 ns/op       65856 B/op         671 allocs/op
BenchmarkKocha_GithubAll                 10000        106315 ns/op       23304 B/op         843 allocs/op
BenchmarkLARS_GithubAll                  47779         25084 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GithubAll                3266        371907 ns/op      149409 B/op        1624 allocs/op
BenchmarkMartini_GithubAll                 331       3444706 ns/op      226551 B/op        2325 allocs/op
BenchmarkPat_GithubAll                     273       4381818 ns/op     1483152 B/op       26963 allocs/op
BenchmarkPossum_GithubAll                10000        164367 ns/op       84448 B/op         609 allocs/op
BenchmarkR2router_GithubAll              10000        160220 ns/op       77328 B/op         979 allocs/op
BenchmarkRivet_GithubAll                 14625         82453 ns/op       16272 B/op         167 allocs/op
BenchmarkTango_GithubAll                  6255        279611 ns/op       63826 B/op        1618 allocs/op
BenchmarkTigerTonic_GithubAll             2008        687874 ns/op      193856 B/op        4474 allocs/op
BenchmarkTraffic_GithubAll                 355       3478508 ns/op      820744 B/op       14114 allocs/op
BenchmarkVulcan_GithubAll                 6885        193333 ns/op       19894 B/op         609 allocs/op
```

## Google+

```sh
BenchmarkGin_GPlusStatic              19247326          62.2 ns/op           0 B/op           0 allocs/op

BenchmarkAce_GPlusStatic              20235060          59.2 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GPlusStatic             31978935          37.6 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GPlusStatic              3516523           341 ns/op         104 B/op           3 allocs/op
BenchmarkBeego_GPlusStatic             1212036           991 ns/op         352 B/op           3 allocs/op
BenchmarkBone_GPlusStatic              6736242           183 ns/op          32 B/op           1 allocs/op
BenchmarkChi_GPlusStatic               1490640           814 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_GPlusStatic            55006856          21.8 ns/op           0 B/op           0 allocs/op
BenchmarkEcho_GPlusStatic             17688258          67.9 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GPlusStatic        1829181           666 ns/op         280 B/op           5 allocs/op
BenchmarkGoji_GPlusStatic              9147451           130 ns/op           0 B/op           0 allocs/op
BenchmarkGojiv2_GPlusStatic             594015          2063 ns/op        1312 B/op          10 allocs/op
BenchmarkGoJsonRest_GPlusStatic        1264906           950 ns/op         329 B/op          11 allocs/op
BenchmarkGoRestful_GPlusStatic          231558          5341 ns/op        3872 B/op          13 allocs/op
BenchmarkGorillaMux_GPlusStatic         908418          1809 ns/op         976 B/op           9 allocs/op
BenchmarkGowwwRouter_GPlusStatic      40684604          29.5 ns/op           0 B/op           0 allocs/op
BenchmarkHttpRouter_GPlusStatic       46742804          25.7 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GPlusStatic      32567161          36.9 ns/op           0 B/op           0 allocs/op
BenchmarkKocha_GPlusStatic            33800060          35.3 ns/op           0 B/op           0 allocs/op
BenchmarkLARS_GPlusStatic             20431858          60.0 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GPlusStatic           1000000          1745 ns/op         736 B/op           8 allocs/op
BenchmarkMartini_GPlusStatic            442248          3619 ns/op         768 B/op           9 allocs/op
BenchmarkPat_GPlusStatic               4328004           292 ns/op          96 B/op           2 allocs/op
BenchmarkPossum_GPlusStatic            1570753           763 ns/op         416 B/op           3 allocs/op
BenchmarkR2router_GPlusStatic          3339474           355 ns/op         144 B/op           4 allocs/op
BenchmarkRivet_GPlusStatic            18570961          64.7 ns/op           0 B/op           0 allocs/op
BenchmarkTango_GPlusStatic             1388702           860 ns/op         200 B/op           8 allocs/op
BenchmarkTigerTonic_GPlusStatic        7803543           159 ns/op          32 B/op           1 allocs/op
BenchmarkTraffic_GPlusStatic            878605          2171 ns/op        1112 B/op          16 allocs/op
BenchmarkVulcan_GPlusStatic            2742446           437 ns/op          98 B/op           3 allocs/op
BenchmarkAce_GPlusParam               11626975           105 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GPlusParam              16914322          71.6 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GPlusParam               1405173           832 ns/op         480 B/op           5 allocs/op
BenchmarkBeego_GPlusParam              1000000          1075 ns/op         352 B/op           3 allocs/op
BenchmarkBone_GPlusParam               1000000          1557 ns/op         816 B/op           6 allocs/op
BenchmarkChi_GPlusParam                1347926           894 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_GPlusParam              5513000           212 ns/op          64 B/op           1 allocs/op
BenchmarkEcho_GPlusParam              11884383           101 ns/op           0 B/op           0 allocs/op
BenchmarkGin_GPlusParam               12898952          93.1 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GPlusParam         1000000          1194 ns/op         648 B/op           8 allocs/op
BenchmarkGoji_GPlusParam               1857229           645 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_GPlusParam              520939          2322 ns/op        1328 B/op          11 allocs/op
BenchmarkGoJsonRest_GPlusParam         1000000          1536 ns/op         649 B/op          13 allocs/op
BenchmarkGoRestful_GPlusParam           205449          5800 ns/op        4192 B/op          14 allocs/op
BenchmarkGorillaMux_GPlusParam          395310          3188 ns/op        1280 B/op          10 allocs/op
BenchmarkGowwwRouter_GPlusParam        1851798           667 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_GPlusParam        18420789          65.2 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GPlusParam        1878463           629 ns/op         352 B/op           3 allocs/op
BenchmarkKocha_GPlusParam              4495610           273 ns/op          56 B/op           3 allocs/op
BenchmarkLARS_GPlusParam              14615976          83.2 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GPlusParam             584145          2549 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_GPlusParam             250501          4583 ns/op        1072 B/op          10 allocs/op
BenchmarkPat_GPlusParam                1000000          1645 ns/op         576 B/op          11 allocs/op
BenchmarkPossum_GPlusParam             1000000          1008 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_GPlusParam           1708191           688 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_GPlusParam              5795014           211 ns/op          48 B/op           1 allocs/op
BenchmarkTango_GPlusParam              1000000          1091 ns/op         264 B/op           8 allocs/op
BenchmarkTigerTonic_GPlusParam          760221          2489 ns/op         856 B/op          16 allocs/op
BenchmarkTraffic_GPlusParam             309774          4039 ns/op        1872 B/op          21 allocs/op
BenchmarkVulcan_GPlusParam             1935730           623 ns/op          98 B/op           3 allocs/op
BenchmarkAce_GPlus2Params              9158314           134 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GPlus2Params            11300517           107 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GPlus2Params             1239238           961 ns/op         496 B/op           5 allocs/op
BenchmarkBeego_GPlus2Params            1000000          1202 ns/op         352 B/op           3 allocs/op
BenchmarkBone_GPlus2Params              335576          3725 ns/op        1168 B/op          10 allocs/op
BenchmarkChi_GPlus2Params              1000000          1014 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_GPlus2Params            4394598           280 ns/op          64 B/op           1 allocs/op
BenchmarkEcho_GPlus2Params             7851861           154 ns/op           0 B/op           0 allocs/op
BenchmarkGin_GPlus2Params              9958588           120 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GPlus2Params       1000000          1433 ns/op         712 B/op           9 allocs/op
BenchmarkGoji_GPlus2Params             1325134           909 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_GPlus2Params            405955          2870 ns/op        1408 B/op          14 allocs/op
BenchmarkGoJsonRest_GPlus2Params        977038          1987 ns/op         713 B/op          14 allocs/op
BenchmarkGoRestful_GPlus2Params         205018          6142 ns/op        4384 B/op          16 allocs/op
BenchmarkGorillaMux_GPlus2Params        205641          6015 ns/op        1296 B/op          10 allocs/op
BenchmarkGowwwRouter_GPlus2Params      1748542           684 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_GPlus2Params      14047102          87.7 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GPlus2Params      1418673           828 ns/op         384 B/op           4 allocs/op
BenchmarkKocha_GPlus2Params            2334562           520 ns/op         128 B/op           5 allocs/op
BenchmarkLARS_GPlus2Params            11954094           101 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GPlus2Params           491552          2890 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_GPlus2Params           120532          9545 ns/op        1200 B/op          13 allocs/op
BenchmarkPat_GPlus2Params               194739          6766 ns/op        2168 B/op          33 allocs/op
BenchmarkPossum_GPlus2Params           1201224          1009 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_GPlus2Params         1575535           756 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_GPlus2Params            3698930           325 ns/op          96 B/op           1 allocs/op
BenchmarkTango_GPlus2Params            1000000          1212 ns/op         344 B/op           8 allocs/op
BenchmarkTigerTonic_GPlus2Params        349350          3660 ns/op        1200 B/op          22 allocs/op
BenchmarkTraffic_GPlus2Params           169714          7862 ns/op        2248 B/op          28 allocs/op
BenchmarkVulcan_GPlus2Params           1222288           974 ns/op          98 B/op           3 allocs/op
BenchmarkAce_GPlusAll                   845606          1398 ns/op           0 B/op           0 allocs/op
BenchmarkAero_GPlusAll                 1000000          1009 ns/op           0 B/op           0 allocs/op
BenchmarkBear_GPlusAll                  103830         11386 ns/op        5488 B/op          61 allocs/op
BenchmarkBeego_GPlusAll                  82653         14784 ns/op        4576 B/op          39 allocs/op
BenchmarkBone_GPlusAll                   36601         33123 ns/op       11744 B/op         109 allocs/op
BenchmarkChi_GPlusAll                    95264         12831 ns/op        5616 B/op          39 allocs/op
BenchmarkDenco_GPlusAll                 567681          2950 ns/op         672 B/op          11 allocs/op
BenchmarkEcho_GPlusAll                  720366          1665 ns/op           0 B/op           0 allocs/op
BenchmarkGin_GPlusAll                  1000000          1185 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_GPlusAll             71575         16365 ns/op        8040 B/op         103 allocs/op
BenchmarkGoji_GPlusAll                  136352          9191 ns/op        3696 B/op          22 allocs/op
BenchmarkGojiv2_GPlusAll                 38006         31802 ns/op       17616 B/op         154 allocs/op
BenchmarkGoJsonRest_GPlusAll             57238         21561 ns/op        8117 B/op         170 allocs/op
BenchmarkGoRestful_GPlusAll              15147         79276 ns/op       55520 B/op         192 allocs/op
BenchmarkGorillaMux_GPlusAll             24446         48410 ns/op       16112 B/op         128 allocs/op
BenchmarkGowwwRouter_GPlusAll           150112          7770 ns/op        4752 B/op          33 allocs/op
BenchmarkHttpRouter_GPlusAll           1367820           878 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_GPlusAll           166628          8004 ns/op        4032 B/op          38 allocs/op
BenchmarkKocha_GPlusAll                 265694          4570 ns/op         976 B/op          43 allocs/op
BenchmarkLARS_GPlusAll                 1000000          1068 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_GPlusAll                54564         23305 ns/op        9568 B/op         104 allocs/op
BenchmarkMartini_GPlusAll                16274         73845 ns/op       14016 B/op         145 allocs/op
BenchmarkPat_GPlusAll                    27181         44478 ns/op       15264 B/op         271 allocs/op
BenchmarkPossum_GPlusAll                122587         10277 ns/op        5408 B/op          39 allocs/op
BenchmarkR2router_GPlusAll              130137          9297 ns/op        5040 B/op          63 allocs/op
BenchmarkRivet_GPlusAll                 532438          3323 ns/op         768 B/op          11 allocs/op
BenchmarkTango_GPlusAll                  86054         14531 ns/op        3656 B/op         104 allocs/op
BenchmarkTigerTonic_GPlusAll             33936         35356 ns/op       11600 B/op         242 allocs/op
BenchmarkTraffic_GPlusAll                17833         68181 ns/op       26248 B/op         341 allocs/op
BenchmarkVulcan_GPlusAll                120109          9861 ns/op        1274 B/op          39 allocs/op
```

## Parse.com

```sh
BenchmarkGin_ParseStatic              18877833          63.5 ns/op           0 B/op           0 allocs/op

BenchmarkAce_ParseStatic              19663731          60.8 ns/op           0 B/op           0 allocs/op
BenchmarkAero_ParseStatic             28967341          41.5 ns/op           0 B/op           0 allocs/op
BenchmarkBear_ParseStatic              3006984           402 ns/op         120 B/op           3 allocs/op
BenchmarkBeego_ParseStatic             1000000          1031 ns/op         352 B/op           3 allocs/op
BenchmarkBone_ParseStatic              1782482           675 ns/op         144 B/op           3 allocs/op
BenchmarkChi_ParseStatic               1453261           819 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_ParseStatic            45023595          26.5 ns/op           0 B/op           0 allocs/op
BenchmarkEcho_ParseStatic             17330470          69.3 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_ParseStatic        1644006           731 ns/op         296 B/op           5 allocs/op
BenchmarkGoji_ParseStatic              7026930           170 ns/op           0 B/op           0 allocs/op
BenchmarkGojiv2_ParseStatic             517618          2037 ns/op        1312 B/op          10 allocs/op
BenchmarkGoJsonRest_ParseStatic        1227080           975 ns/op         329 B/op          11 allocs/op
BenchmarkGoRestful_ParseStatic          192458          6659 ns/op        4256 B/op          13 allocs/op
BenchmarkGorillaMux_ParseStatic         744062          2109 ns/op         976 B/op           9 allocs/op
BenchmarkGowwwRouter_ParseStatic      37781062          31.8 ns/op           0 B/op           0 allocs/op
BenchmarkHttpRouter_ParseStatic       45311223          26.5 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_ParseStatic      21383475          56.1 ns/op           0 B/op           0 allocs/op
BenchmarkKocha_ParseStatic            29953290          40.1 ns/op           0 B/op           0 allocs/op
BenchmarkLARS_ParseStatic             20036196          62.7 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_ParseStatic           1000000          1740 ns/op         736 B/op           8 allocs/op
BenchmarkMartini_ParseStatic            404156          3801 ns/op         768 B/op           9 allocs/op
BenchmarkPat_ParseStatic               1547180           772 ns/op         240 B/op           5 allocs/op
BenchmarkPossum_ParseStatic            1608991           757 ns/op         416 B/op           3 allocs/op
BenchmarkR2router_ParseStatic          3177936           385 ns/op         144 B/op           4 allocs/op
BenchmarkRivet_ParseStatic            17783205          67.4 ns/op           0 B/op           0 allocs/op
BenchmarkTango_ParseStatic             1210777           990 ns/op         248 B/op           8 allocs/op
BenchmarkTigerTonic_ParseStatic        5316440           231 ns/op          48 B/op           1 allocs/op
BenchmarkTraffic_ParseStatic            496050          2539 ns/op        1256 B/op          19 allocs/op
BenchmarkVulcan_ParseStatic            2462798           488 ns/op          98 B/op           3 allocs/op
BenchmarkAce_ParseParam               13393669          89.6 ns/op           0 B/op           0 allocs/op
BenchmarkAero_ParseParam              19836619          60.4 ns/op           0 B/op           0 allocs/op
BenchmarkBear_ParseParam               1405954           864 ns/op         467 B/op           5 allocs/op
BenchmarkBeego_ParseParam              1000000          1065 ns/op         352 B/op           3 allocs/op
BenchmarkBone_ParseParam               1000000          1698 ns/op         896 B/op           7 allocs/op
BenchmarkChi_ParseParam                1356037           873 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_ParseParam              6241392           204 ns/op          64 B/op           1 allocs/op
BenchmarkEcho_ParseParam              14088100          85.1 ns/op           0 B/op           0 allocs/op
BenchmarkGin_ParseParam               17426064          68.9 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_ParseParam         1000000          1254 ns/op         664 B/op           8 allocs/op
BenchmarkGoji_ParseParam               1682574           713 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_ParseParam              502224          2333 ns/op        1360 B/op          12 allocs/op
BenchmarkGoJsonRest_ParseParam         1000000          1401 ns/op         649 B/op          13 allocs/op
BenchmarkGoRestful_ParseParam           182623          7097 ns/op        4576 B/op          14 allocs/op
BenchmarkGorillaMux_ParseParam          482332          2477 ns/op        1280 B/op          10 allocs/op
BenchmarkGowwwRouter_ParseParam        1834873           657 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_ParseParam        23593393          51.0 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_ParseParam        2100160           574 ns/op         352 B/op           3 allocs/op
BenchmarkKocha_ParseParam              4837220           252 ns/op          56 B/op           3 allocs/op
BenchmarkLARS_ParseParam              18411192          66.2 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_ParseParam             571870          2398 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_ParseParam             286262          4268 ns/op        1072 B/op          10 allocs/op
BenchmarkPat_ParseParam                 692906          2157 ns/op         992 B/op          15 allocs/op
BenchmarkPossum_ParseParam             1000000          1011 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_ParseParam           1722735           697 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_ParseParam              6058054           203 ns/op          48 B/op           1 allocs/op
BenchmarkTango_ParseParam              1000000          1061 ns/op         280 B/op           8 allocs/op
BenchmarkTigerTonic_ParseParam          890275          2277 ns/op         784 B/op          15 allocs/op
BenchmarkTraffic_ParseParam             351322          3543 ns/op        1896 B/op          21 allocs/op
BenchmarkVulcan_ParseParam             2076544           572 ns/op          98 B/op           3 allocs/op
BenchmarkAce_Parse2Params             11718074           101 ns/op           0 B/op           0 allocs/op
BenchmarkAero_Parse2Params            16264988          73.4 ns/op           0 B/op           0 allocs/op
BenchmarkBear_Parse2Params             1238322           973 ns/op         496 B/op           5 allocs/op
BenchmarkBeego_Parse2Params            1000000          1120 ns/op         352 B/op           3 allocs/op
BenchmarkBone_Parse2Params             1000000          1632 ns/op         848 B/op           6 allocs/op
BenchmarkChi_Parse2Params              1239477           955 ns/op         432 B/op           3 allocs/op
BenchmarkDenco_Parse2Params            4944133           245 ns/op          64 B/op           1 allocs/op
BenchmarkEcho_Parse2Params            10518286           114 ns/op           0 B/op           0 allocs/op
BenchmarkGin_Parse2Params             14505195          82.7 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_Parse2Params       1000000          1437 ns/op         712 B/op           9 allocs/op
BenchmarkGoji_Parse2Params             1689883           707 ns/op         336 B/op           2 allocs/op
BenchmarkGojiv2_Parse2Params            502334          2308 ns/op        1344 B/op          11 allocs/op
BenchmarkGoJsonRest_Parse2Params       1000000          1771 ns/op         713 B/op          14 allocs/op
BenchmarkGoRestful_Parse2Params         159092          7583 ns/op        4928 B/op          14 allocs/op
BenchmarkGorillaMux_Parse2Params        417548          2980 ns/op        1296 B/op          10 allocs/op
BenchmarkGowwwRouter_Parse2Params      1751737           686 ns/op         432 B/op           3 allocs/op
BenchmarkHttpRouter_Parse2Params      18089204          66.3 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_Parse2Params      1556986           777 ns/op         384 B/op           4 allocs/op
BenchmarkKocha_Parse2Params            2493082           485 ns/op         128 B/op           5 allocs/op
BenchmarkLARS_Parse2Params            15350108          78.5 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_Parse2Params           530974          2605 ns/op        1072 B/op          10 allocs/op
BenchmarkMartini_Parse2Params           247069          4673 ns/op        1152 B/op          11 allocs/op
BenchmarkPat_Parse2Params               816295          2126 ns/op         752 B/op          16 allocs/op
BenchmarkPossum_Parse2Params           1000000          1002 ns/op         496 B/op           5 allocs/op
BenchmarkR2router_Parse2Params         1569771           733 ns/op         432 B/op           5 allocs/op
BenchmarkRivet_Parse2Params            4080546           295 ns/op          96 B/op           1 allocs/op
BenchmarkTango_Parse2Params            1000000          1121 ns/op         312 B/op           8 allocs/op
BenchmarkTigerTonic_Parse2Params        399556          3470 ns/op        1168 B/op          22 allocs/op
BenchmarkTraffic_Parse2Params           314194          4159 ns/op        1944 B/op          22 allocs/op
BenchmarkVulcan_Parse2Params           1827559           664 ns/op          98 B/op           3 allocs/op
BenchmarkAce_ParseAll                   478395          2503 ns/op           0 B/op           0 allocs/op
BenchmarkAero_ParseAll                  715392          1658 ns/op           0 B/op           0 allocs/op
BenchmarkBear_ParseAll                   59191         20124 ns/op        8928 B/op         110 allocs/op
BenchmarkBeego_ParseAll                  45507         27266 ns/op        9152 B/op          78 allocs/op
BenchmarkBone_ParseAll                   29328         41459 ns/op       16208 B/op         147 allocs/op
BenchmarkChi_ParseAll                    48531         25053 ns/op       11232 B/op          78 allocs/op
BenchmarkDenco_ParseAll                 325532          4284 ns/op         928 B/op          16 allocs/op
BenchmarkEcho_ParseAll                  433771          2759 ns/op           0 B/op           0 allocs/op
BenchmarkGin_ParseAll                   576316          2082 ns/op           0 B/op           0 allocs/op
BenchmarkGocraftWeb_ParseAll             41500         29692 ns/op       13728 B/op         181 allocs/op
BenchmarkGoji_ParseAll                   80833         15563 ns/op        5376 B/op          32 allocs/op
BenchmarkGojiv2_ParseAll                 19836         60335 ns/op       34448 B/op         277 allocs/op
BenchmarkGoJsonRest_ParseAll             32210         38027 ns/op       13866 B/op         321 allocs/op
BenchmarkGoRestful_ParseAll               6644        190842 ns/op      117600 B/op         354 allocs/op
BenchmarkGorillaMux_ParseAll             12634         95894 ns/op       30288 B/op         250 allocs/op
BenchmarkGowwwRouter_ParseAll            98152         12159 ns/op        6912 B/op          48 allocs/op
BenchmarkHttpRouter_ParseAll            933208          1273 ns/op           0 B/op           0 allocs/op
BenchmarkHttpTreeMux_ParseAll           107191         11554 ns/op        5728 B/op          51 allocs/op
BenchmarkKocha_ParseAll                 184862          6225 ns/op        1112 B/op          54 allocs/op
BenchmarkLARS_ParseAll                  644546          1858 ns/op           0 B/op           0 allocs/op
BenchmarkMacaron_ParseAll                26145         46484 ns/op       19136 B/op         208 allocs/op
BenchmarkMartini_ParseAll                10000        121838 ns/op       25072 B/op         253 allocs/op
BenchmarkPat_ParseAll                    25417         47196 ns/op       15216 B/op         308 allocs/op
BenchmarkPossum_ParseAll                 58550         20735 ns/op       10816 B/op          78 allocs/op
BenchmarkR2router_ParseAll               72732         16584 ns/op        8352 B/op         120 allocs/op
BenchmarkRivet_ParseAll                 281365          4968 ns/op         912 B/op          16 allocs/op
BenchmarkTango_ParseAll                  42831         28668 ns/op        7168 B/op         208 allocs/op
BenchmarkTigerTonic_ParseAll             23774         49972 ns/op       16048 B/op         332 allocs/op
BenchmarkTraffic_ParseAll                10000        104679 ns/op       45520 B/op         605 allocs/op
BenchmarkVulcan_ParseAll                 64810         18108 ns/op        2548 B/op          78 allocs/op
```
