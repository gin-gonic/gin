
# Benchmark System

**VM HOST:** Travis  
**Machine:** Red Hat 4.8.5-36
**Date:** Dec 23th, 2020
**Version:** Gin v1.6.3
**Go Version:** go1.15.6 linux/amd64
**Source:** [Go HTTP Router Benchmark](https://github.com/gin-gonic/go-http-routing-benchmark)
**Result:** [See the gist](https://gist.github.com/appleboy/b5f2ecfaf50824ae9c64dcfb9165ae5e) or [Travis result](https://travis-ci.org/github/gin-gonic/go-http-routing-benchmark/jobs/682947061)

## Static Routes: 157

```sh
Gin: 34984 Bytes

HttpServeMux: 14512 Bytes
Ace: 30680 Bytes
Aero: 34536 Bytes
Bear: 30456 Bytes
Beego: 98456 Bytes
Bone: 40224 Bytes
Chi: 83608 Bytes
Denco: 10216 Bytes
Echo: 80280 Bytes
GocraftWeb: 55496 Bytes
Goji: 29744 Bytes
Gojiv2: 105840 Bytes
GoJsonRest: 137512 Bytes
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
Traffic: 538992 Bytes
Vulcan: 369960 Bytes
```

## GithubAPI Routes: 203

```sh
Gin: 59328 Bytes

Ace: 48688 Bytes
Aero: 90912 Bytes
Bear: 82744 Bytes
Beego: 150872 Bytes
Bone: 100976 Bytes
Chi: 95112 Bytes
Denco: 36448 Bytes
Echo: 100456 Bytes
GocraftWeb: 95640 Bytes
Goji: 49680 Bytes
Gojiv2: 104704 Bytes
GoJsonRest: 142200 Bytes
GoRestful: 1241656 Bytes
GorillaMux: 1322784 Bytes
GowwwRouter: 80008 Bytes
HttpRouter: 37144 Bytes
HttpTreeMux: 78800 Bytes
Kocha: 784976 Bytes
LARS: 48600 Bytes
Macaron: 93680 Bytes
Martini: 485264 Bytes
Pat: 21200 Bytes
Possum: 85312 Bytes
R2router: 47104 Bytes
Rivet: 42840 Bytes
Tango: 54840 Bytes
TigerTonic: 95664 Bytes
Traffic: 921744 Bytes
Vulcan: 425368 Bytes
```

## GPlusAPI Routes: 13

```sh
Gin: 4464 Bytes

Ace: 3712 Bytes
Aero: 26056 Bytes
Bear: 7112 Bytes
Beego: 10272 Bytes
Bone: 6688 Bytes
Chi: 8024 Bytes
Denco: 3264 Bytes
Echo: 9640 Bytes
GocraftWeb: 7496 Bytes
Goji: 3152 Bytes
Gojiv2: 7376 Bytes
GoJsonRest: 11416 Bytes
GoRestful: 74328 Bytes
GorillaMux: 66208 Bytes
GowwwRouter: 5744 Bytes
HttpRouter: 2808 Bytes
HttpTreeMux: 7440 Bytes
Kocha: 128880 Bytes
LARS: 3656 Bytes
Macaron: 8864 Bytes
Martini: 23920 Bytes
Pat: 1856 Bytes
Possum: 7728 Bytes
R2router: 3928 Bytes
Rivet: 3064 Bytes
Tango: 5168 Bytes
TigerTonic: 9408 Bytes
Traffic: 46400 Bytes
Vulcan: 25544 Bytes
```

## ParseAPI Routes: 26

```sh
Gin: 7808 Bytes

Ace: 6704 Bytes
Aero: 28488 Bytes
Bear: 12320 Bytes
Beego: 19184 Bytes
Bone: 11440 Bytes
Chi: 9744 Bytes
Denco: 4192 Bytes
Echo: 11824 Bytes
GocraftWeb: 12800 Bytes
Goji: 5680 Bytes
Gojiv2: 14464 Bytes
GoJsonRest: 14200 Bytes
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
BenchmarkGin_StaticAll            	   49336	     24033 ns/op	       0 B/op	       0 allocs/op

BenchmarkAce_StaticAll            	   50060	     23107 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_StaticAll           	  100614	     11642 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpServeMux_StaticAll   	   42536	     28392 ns/op	       0 B/op	       0 allocs/op
BenchmarkBeego_StaticAll          	    7972	    236801 ns/op	   55264 B/op	     471 allocs/op
BenchmarkBear_StaticAll           	   12056	    102094 ns/op	   20272 B/op	     469 allocs/op
BenchmarkBone_StaticAll           	   15552	     77261 ns/op	       0 B/op	       0 allocs/op
BenchmarkChi_StaticAll            	   10000	    208043 ns/op	   70336 B/op	     471 allocs/op
BenchmarkDenco_StaticAll          	  164359	      7487 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_StaticAll           	   40406	     29253 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_StaticAll     	   10000	    153266 ns/op	   46312 B/op	     785 allocs/op
BenchmarkGoji_StaticAll           	   24103	     49582 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_StaticAll         	    2949	    481678 ns/op	  213520 B/op	    1570 allocs/op
BenchmarkGoJsonRest_StaticAll     	    7600	    232376 ns/op	   51653 B/op	    1727 allocs/op
BenchmarkGoRestful_StaticAll      	     618	   2063453 ns/op	  613280 B/op	    2053 allocs/op
BenchmarkGorillaMux_StaticAll     	     988	   1225690 ns/op	  158257 B/op	    1413 allocs/op
BenchmarkGowwwRouter_StaticAll    	   51604	     22902 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_StaticAll     	   87294	     13587 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_StaticAll    	   86989	     13886 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_StaticAll          	   74659	     16005 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_StaticAll           	   55255	     21618 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_StaticAll        	    5268	    354738 ns/op	  115553 B/op	    1256 allocs/op
BenchmarkMartini_StaticAll        	     667	   1777210 ns/op	  125444 B/op	    1717 allocs/op
BenchmarkPat_StaticAll            	     458	   2707947 ns/op	  602832 B/op	   12559 allocs/op
BenchmarkPossum_StaticAll         	   10000	    164330 ns/op	   65312 B/op	     471 allocs/op
BenchmarkR2router_StaticAll       	   13470	     90500 ns/op	   22608 B/op	     628 allocs/op
BenchmarkRivet_StaticAll          	   37128	     32519 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_StaticAll          	    6505	    247626 ns/op	   37953 B/op	    1099 allocs/op
BenchmarkTigerTonic_StaticAll     	   21229	     56544 ns/op	    7376 B/op	     157 allocs/op
BenchmarkTraffic_StaticAll        	     427	   2890102 ns/op	  754862 B/op	   14601 allocs/op
BenchmarkVulcan_StaticAll         	    8528	    161477 ns/op	   15386 B/op	     471 allocs/op
```

## Micro Benchmarks

```sh
BenchmarkGin_Param                	14482286	        84.3 ns/op	       0 B/op	       0 allocs/op

BenchmarkAce_Param                	11408986	       104 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param               	19642081	        62.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param               	 1000000	      1107 ns/op	     456 B/op	       5 allocs/op
BenchmarkBeego_Param              	 1000000	      1318 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Param               	 1000000	      1813 ns/op	     832 B/op	       6 allocs/op
BenchmarkChi_Param                	 1000000	      1142 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_Param              	 7226908	       163 ns/op	      32 B/op	       1 allocs/op
BenchmarkEcho_Param               	13762222	        88.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param         	 1000000	      1662 ns/op	     648 B/op	       8 allocs/op
BenchmarkGoji_Param               	 1464320	       819 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Param             	  387694	      2792 ns/op	    1376 B/op	      11 allocs/op
BenchmarkGoJsonRest_Param         	 1000000	      1818 ns/op	     649 B/op	      13 allocs/op
BenchmarkGoRestful_Param          	  187965	      6963 ns/op	    4192 B/op	      14 allocs/op
BenchmarkGorillaMux_Param         	  389700	      3250 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGowwwRouter_Param        	 1405456	       842 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_Param         	21498526	        57.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_Param        	 1694540	       712 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_Param              	 4053416	       302 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_Param               	15720844	        81.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param            	  541933	      3152 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_Param            	  259760	      5084 ns/op	    1072 B/op	      10 allocs/op
BenchmarkPat_Param                	 1000000	      1726 ns/op	     536 B/op	      11 allocs/op
BenchmarkPossum_Param             	 1000000	      1383 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Param           	 1294652	       868 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Param              	 5570419	       218 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_Param              	 1000000	      1177 ns/op	     240 B/op	       7 allocs/op
BenchmarkTigerTonic_Param         	  705588	      2686 ns/op	     776 B/op	      16 allocs/op
BenchmarkTraffic_Param            	  277063	      4672 ns/op	    1856 B/op	      21 allocs/op
BenchmarkVulcan_Param             	 2075674	       588 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_Param5               	 7781608	       156 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param5              	13826379	        86.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param5              	 1000000	      1499 ns/op	     501 B/op	       5 allocs/op
BenchmarkBeego_Param5             	 1000000	      1652 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Param5              	 1000000	      2563 ns/op	     880 B/op	       6 allocs/op
BenchmarkChi_Param5               	 1000000	      1529 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_Param5             	 2493051	       482 ns/op	     160 B/op	       1 allocs/op
BenchmarkEcho_Param5              	 5406230	       227 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param5               	 8970267	       135 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param5        	  923665	      2594 ns/op	     920 B/op	      11 allocs/op
BenchmarkGoji_Param5              	 1000000	      1064 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Param5            	  336223	      3247 ns/op	    1440 B/op	      11 allocs/op
BenchmarkGoJsonRest_Param5        	  426156	      3383 ns/op	    1097 B/op	      16 allocs/op
BenchmarkGoRestful_Param5         	  175233	      7518 ns/op	    4288 B/op	      14 allocs/op
BenchmarkGorillaMux_Param5        	  272763	      4304 ns/op	    1376 B/op	      10 allocs/op
BenchmarkGowwwRouter_Param5       	 1300671	       924 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_Param5        	10856960	       108 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_Param5       	 1000000	      1723 ns/op	     576 B/op	       6 allocs/op
BenchmarkKocha_Param5             	 1000000	      1375 ns/op	     440 B/op	      10 allocs/op
BenchmarkLARS_Param5              	 9627297	       129 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param5           	  397291	      3444 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_Param5           	  209604	      6112 ns/op	    1232 B/op	      11 allocs/op
BenchmarkPat_Param5               	  372729	      4142 ns/op	     888 B/op	      29 allocs/op
BenchmarkPossum_Param5            	 1000000	      1417 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Param5          	 1000000	      1155 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Param5             	 1721068	       722 ns/op	     240 B/op	       1 allocs/op
BenchmarkTango_Param5             	 1000000	      1557 ns/op	     352 B/op	       7 allocs/op
BenchmarkTigerTonic_Param5        	  147892	      9669 ns/op	    2279 B/op	      39 allocs/op
BenchmarkTraffic_Param5           	  194479	      6882 ns/op	    2208 B/op	      27 allocs/op
BenchmarkVulcan_Param5            	 1542697	       784 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_Param20              	 3478586	       344 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Param20             	24969956	        45.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param20             	  271641	      4855 ns/op	    1665 B/op	       5 allocs/op
BenchmarkBeego_Param20            	  441153	      3397 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Param20             	  241766	      6238 ns/op	    2047 B/op	       6 allocs/op
BenchmarkChi_Param20              	  721503	      2580 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_Param20            	 1000000	      1719 ns/op	     640 B/op	       1 allocs/op
BenchmarkEcho_Param20             	 1724497	       679 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Param20              	 3782846	       326 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param20       	  139232	      9401 ns/op	    3796 B/op	      15 allocs/op
BenchmarkGoji_Param20             	  393200	      3474 ns/op	    1247 B/op	       2 allocs/op
BenchmarkGojiv2_Param20           	  280774	      4394 ns/op	    1680 B/op	      11 allocs/op
BenchmarkGoJsonRest_Param20       	   93594	     13209 ns/op	    4485 B/op	      20 allocs/op
BenchmarkGoRestful_Param20        	   76891	     16265 ns/op	    6716 B/op	      18 allocs/op
BenchmarkGorillaMux_Param20       	   97863	     10830 ns/op	    3484 B/op	      12 allocs/op
BenchmarkGowwwRouter_Param20      	 1000000	      1583 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_Param20       	 4062788	       295 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_Param20      	  142550	      8804 ns/op	    3196 B/op	      10 allocs/op
BenchmarkKocha_Param20            	  280284	      4654 ns/op	    1808 B/op	      27 allocs/op
BenchmarkLARS_Param20             	 3995473	       302 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param20          	  136138	      8806 ns/op	    2924 B/op	      12 allocs/op
BenchmarkMartini_Param20          	  102451	     12173 ns/op	    3596 B/op	      13 allocs/op
BenchmarkPat_Param20              	   64381	     18743 ns/op	    4423 B/op	      93 allocs/op
BenchmarkPossum_Param20           	 1000000	      1360 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Param20         	  243445	      5862 ns/op	    2284 B/op	       7 allocs/op
BenchmarkRivet_Param20            	  590210	      2589 ns/op	    1024 B/op	       1 allocs/op
BenchmarkTango_Param20            	  668371	      2989 ns/op	     848 B/op	       7 allocs/op
BenchmarkTigerTonic_Param20       	   30514	     39525 ns/op	    9865 B/op	     119 allocs/op
BenchmarkTraffic_Param20          	   46384	     26394 ns/op	    7857 B/op	      47 allocs/op
BenchmarkVulcan_Param20           	  893665	      1501 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_ParamWrite           	 5379445	       228 ns/op	       8 B/op	       1 allocs/op
BenchmarkAero_ParamWrite          	11855682	       105 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParamWrite          	 1000000	      1254 ns/op	     456 B/op	       5 allocs/op
BenchmarkBeego_ParamWrite         	 1000000	      1463 ns/op	     360 B/op	       4 allocs/op
BenchmarkBone_ParamWrite          	 1000000	      2015 ns/op	     832 B/op	       6 allocs/op
BenchmarkChi_ParamWrite           	 1000000	      1225 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_ParamWrite         	 5793268	       210 ns/op	      32 B/op	       1 allocs/op
BenchmarkEcho_ParamWrite          	 6736492	       203 ns/op	       8 B/op	       1 allocs/op
BenchmarkGin_ParamWrite           	 8405484	       150 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParamWrite    	 1000000	      1719 ns/op	     656 B/op	       9 allocs/op
BenchmarkGoji_ParamWrite          	 1409599	       846 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_ParamWrite        	  400083	      2968 ns/op	    1408 B/op	      13 allocs/op
BenchmarkGoJsonRest_ParamWrite    	  494056	      2814 ns/op	    1128 B/op	      18 allocs/op
BenchmarkGoRestful_ParamWrite     	  183366	      7116 ns/op	    4200 B/op	      15 allocs/op
BenchmarkGorillaMux_ParamWrite    	  386100	      3510 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGowwwRouter_ParamWrite   	 1000000	      2389 ns/op	    1008 B/op	       8 allocs/op
BenchmarkHttpRouter_ParamWrite    	11336546	       105 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_ParamWrite   	 1534879	       843 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_ParamWrite         	 3425571	       351 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_ParamWrite          	 8510874	       134 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParamWrite       	  406435	      3497 ns/op	    1176 B/op	      14 allocs/op
BenchmarkMartini_ParamWrite       	  197914	      5774 ns/op	    1176 B/op	      14 allocs/op
BenchmarkPat_ParamWrite           	  632290	      2798 ns/op	     960 B/op	      15 allocs/op
BenchmarkPossum_ParamWrite        	 1000000	      1360 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_ParamWrite      	 1280834	       948 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_ParamWrite         	 2897456	       407 ns/op	     112 B/op	       2 allocs/op
BenchmarkTango_ParamWrite         	 1829428	       678 ns/op	     136 B/op	       4 allocs/op
BenchmarkTigerTonic_ParamWrite    	  342619	      4186 ns/op	    1216 B/op	      21 allocs/op
BenchmarkTraffic_ParamWrite       	  236800	      5970 ns/op	    2280 B/op	      25 allocs/op
BenchmarkVulcan_ParamWrite        	 2045410	       593 ns/op	      98 B/op	       3 allocs/op
```

## GitHub

```sh
BenchmarkGin_GithubStatic         	12447151	        98.2 ns/op	       0 B/op	       0 allocs/op

BenchmarkAce_GithubStatic         	12695574	        94.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubStatic        	20914966	        57.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubStatic        	 2138509	       560 ns/op	     120 B/op	       3 allocs/op
BenchmarkBeego_GithubStatic       	 1000000	      1330 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GithubStatic        	   79369	     15663 ns/op	    2880 B/op	      60 allocs/op
BenchmarkChi_GithubStatic         	 1000000	      1168 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_GithubStatic       	31023784	        38.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GithubStatic        	10384629	       119 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubStatic  	 1255374	       979 ns/op	     296 B/op	       5 allocs/op
BenchmarkGoji_GithubStatic        	 5014126	       220 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_GithubStatic      	  428376	      2803 ns/op	    1360 B/op	      10 allocs/op
BenchmarkGoRestful_GithubStatic   	   81528	     15813 ns/op	    4256 B/op	      13 allocs/op
BenchmarkGoJsonRest_GithubStatic  	 1000000	      1290 ns/op	     329 B/op	      11 allocs/op
BenchmarkGorillaMux_GithubStatic  	  203433	      6529 ns/op	    1008 B/op	       9 allocs/op
BenchmarkGowwwRouter_GithubStatic 	12604948	        94.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GithubStatic  	23026732	        51.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GithubStatic 	17519814	        68.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_GithubStatic       	18998347	        63.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_GithubStatic        	12535339	        95.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubStatic     	  932761	      2234 ns/op	     736 B/op	       8 allocs/op
BenchmarkMartini_GithubStatic     	  154873	      8815 ns/op	     768 B/op	       9 allocs/op
BenchmarkPat_GithubStatic         	   94911	     13394 ns/op	    3648 B/op	      76 allocs/op
BenchmarkPossum_GithubStatic      	 1201784	      1012 ns/op	     416 B/op	       3 allocs/op
BenchmarkR2router_GithubStatic    	 2526429	       477 ns/op	     144 B/op	       4 allocs/op
BenchmarkRivet_GithubStatic       	10158662	       119 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_GithubStatic       	 1000000	      1435 ns/op	     240 B/op	       7 allocs/op
BenchmarkTigerTonic_GithubStatic  	 3931868	       317 ns/op	      48 B/op	       1 allocs/op
BenchmarkTraffic_GithubStatic     	   81194	     14835 ns/op	    4664 B/op	      90 allocs/op
BenchmarkVulcan_GithubStatic      	 1313310	       893 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_GithubParam          	 6994286	       174 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubParam         	10592709	       115 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubParam         	 1000000	      1379 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_GithubParam        	 1000000	      1493 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GithubParam         	  158872	      8510 ns/op	    1904 B/op	      19 allocs/op
BenchmarkChi_GithubParam          	 1000000	      1445 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_GithubParam        	 2979142	       410 ns/op	     128 B/op	       1 allocs/op
BenchmarkEcho_GithubParam         	 5409882	       228 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubParam          	 7978224	       153 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubParam   	 1000000	      1956 ns/op	     712 B/op	       9 allocs/op
BenchmarkGoji_GithubParam         	 1000000	      1198 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GithubParam       	  320469	      3728 ns/op	    1456 B/op	      13 allocs/op
BenchmarkGoJsonRest_GithubParam   	  855488	      2273 ns/op	     713 B/op	      14 allocs/op
BenchmarkGoRestful_GithubParam    	   62221	     19740 ns/op	    4352 B/op	      16 allocs/op
BenchmarkGorillaMux_GithubParam   	  119104	      9406 ns/op	    1328 B/op	      10 allocs/op
BenchmarkGowwwRouter_GithubParam  	 1250324	       959 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_GithubParam   	 9249931	       136 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GithubParam  	 1000000	      1153 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_GithubParam        	 1762138	       647 ns/op	     128 B/op	       5 allocs/op
BenchmarkLARS_GithubParam         	 7947790	       152 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubParam      	  446935	      3322 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_GithubParam      	  105180	     10987 ns/op	    1152 B/op	      11 allocs/op
BenchmarkPat_GithubParam          	  126670	     10816 ns/op	    2408 B/op	      48 allocs/op
BenchmarkPossum_GithubParam       	 1000000	      1382 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_GithubParam     	 1268685	       961 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_GithubParam        	 2503034	       478 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_GithubParam        	 1000000	      1659 ns/op	     336 B/op	       7 allocs/op
BenchmarkTigerTonic_GithubParam   	  309370	      4395 ns/op	    1176 B/op	      22 allocs/op
BenchmarkTraffic_GithubParam      	   95318	     12518 ns/op	    2816 B/op	      40 allocs/op
BenchmarkVulcan_GithubParam       	  933199	      1385 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_GithubAll            	   33668	     35010 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubAll           	   48855	     24599 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubAll           	    7350	    274253 ns/op	   86448 B/op	     943 allocs/op
BenchmarkBeego_GithubAll          	    5942	    310119 ns/op	   71456 B/op	     609 allocs/op
BenchmarkBone_GithubAll           	     344	   3447934 ns/op	  722832 B/op	    8620 allocs/op
BenchmarkChi_GithubAll            	    7850	    297949 ns/op	   90944 B/op	     609 allocs/op
BenchmarkDenco_GithubAll          	   14900	     80650 ns/op	   20224 B/op	     167 allocs/op
BenchmarkEcho_GithubAll           	   26797	     44659 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GithubAll            	   37629	     31788 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubAll     	    3061	    397318 ns/op	  131656 B/op	    1686 allocs/op
BenchmarkGoji_GithubAll           	    2940	    479530 ns/op	   56112 B/op	     334 allocs/op
BenchmarkGojiv2_GithubAll         	    1201	   1099221 ns/op	  362464 B/op	    4321 allocs/op
BenchmarkGoJsonRest_GithubAll     	    2488	    464032 ns/op	  134371 B/op	    2737 allocs/op
BenchmarkGoRestful_GithubAll      	     320	   3766819 ns/op	  910144 B/op	    2938 allocs/op
BenchmarkGorillaMux_GithubAll     	     273	   4317768 ns/op	  258146 B/op	    1994 allocs/op
BenchmarkGowwwRouter_GithubAll    	   10000	    187451 ns/op	   74816 B/op	     501 allocs/op
BenchmarkHttpRouter_GithubAll     	   48434	     24781 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GithubAll    	   10000	    196925 ns/op	   65856 B/op	     671 allocs/op
BenchmarkKocha_GithubAll          	   10000	    131909 ns/op	   23304 B/op	     843 allocs/op
BenchmarkLARS_GithubAll           	   40170	     29822 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubAll        	    2650	    464097 ns/op	  149409 B/op	    1624 allocs/op
BenchmarkMartini_GithubAll        	     290	   3951758 ns/op	  226551 B/op	    2325 allocs/op
BenchmarkPat_GithubAll            	     232	   5204752 ns/op	 1483152 B/op	   26963 allocs/op
BenchmarkPossum_GithubAll         	    9512	    220569 ns/op	   84448 B/op	     609 allocs/op
BenchmarkR2router_GithubAll       	   10000	    212878 ns/op	   77328 B/op	     979 allocs/op
BenchmarkRivet_GithubAll          	   10000	    101505 ns/op	   16272 B/op	     167 allocs/op
BenchmarkTango_GithubAll          	    5066	    351082 ns/op	   62225 B/op	    1418 allocs/op
BenchmarkTigerTonic_GithubAll     	    1575	    857892 ns/op	  193856 B/op	    4474 allocs/op
BenchmarkTraffic_GithubAll        	     276	   4500196 ns/op	  820743 B/op	   14114 allocs/op
BenchmarkVulcan_GithubAll         	    5672	    241725 ns/op	   19894 B/op	     609 allocs/op
```

## Google+

```sh
BenchmarkGin_GPlusStatic          	15277999	        78.5 ns/op	       0 B/op	       0 allocs/op

BenchmarkAce_GPlusStatic          	15926042	        75.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusStatic         	26459589	        45.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusStatic         	 2781216	       447 ns/op	     104 B/op	       3 allocs/op
BenchmarkBeego_GPlusStatic        	 1000000	      1252 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GPlusStatic         	 5681490	       225 ns/op	      32 B/op	       1 allocs/op
BenchmarkChi_GPlusStatic          	 1000000	      1096 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_GPlusStatic        	43606446	        27.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlusStatic         	14120419	        85.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusStatic   	 1415466	       864 ns/op	     280 B/op	       5 allocs/op
BenchmarkGoji_GPlusStatic         	 7619877	       158 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_GPlusStatic       	  441046	      2724 ns/op	    1360 B/op	      10 allocs/op
BenchmarkGoJsonRest_GPlusStatic   	 1000000	      1090 ns/op	     329 B/op	      11 allocs/op
BenchmarkGoRestful_GPlusStatic    	  183778	      6814 ns/op	    3872 B/op	      13 allocs/op
BenchmarkGorillaMux_GPlusStatic   	  599338	      2459 ns/op	    1008 B/op	       9 allocs/op
BenchmarkGowwwRouter_GPlusStatic  	32987973	        36.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GPlusStatic   	40541908	        31.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GPlusStatic  	28913025	        42.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_GPlusStatic        	27474925	        43.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_GPlusStatic         	16731944	        71.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusStatic      	  958616	      2217 ns/op	     736 B/op	       8 allocs/op
BenchmarkMartini_GPlusStatic      	  367564	      4354 ns/op	     768 B/op	       9 allocs/op
BenchmarkPat_GPlusStatic          	 3692040	       325 ns/op	      96 B/op	       2 allocs/op
BenchmarkPossum_GPlusStatic       	 1228786	       969 ns/op	     416 B/op	       3 allocs/op
BenchmarkR2router_GPlusStatic     	 2750990	       446 ns/op	     144 B/op	       4 allocs/op
BenchmarkRivet_GPlusStatic        	16919646	        71.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_GPlusStatic        	 1000000	      1060 ns/op	     192 B/op	       7 allocs/op
BenchmarkTigerTonic_GPlusStatic   	 6248180	       199 ns/op	      32 B/op	       1 allocs/op
BenchmarkTraffic_GPlusStatic      	  523586	      2841 ns/op	    1112 B/op	      16 allocs/op
BenchmarkVulcan_GPlusStatic       	 2270104	       540 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_GPlusParam           	 9733280	       125 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusParam          	13036228	        77.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusParam          	 1000000	      1096 ns/op	     480 B/op	       5 allocs/op
BenchmarkBeego_GPlusParam         	 1000000	      1401 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GPlusParam          	 1000000	      2034 ns/op	     832 B/op	       6 allocs/op
BenchmarkChi_GPlusParam           	 1000000	      1299 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_GPlusParam         	 4716946	       259 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_GPlusParam          	10018851	       120 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusParam           	11831241	       104 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusParam    	 1000000	      1619 ns/op	     648 B/op	       8 allocs/op
BenchmarkGoji_GPlusParam          	 1402820	       850 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GPlusParam        	  388882	      2982 ns/op	    1376 B/op	      11 allocs/op
BenchmarkGoJsonRest_GPlusParam    	 1000000	      1939 ns/op	     649 B/op	      13 allocs/op
BenchmarkGoRestful_GPlusParam     	  171267	      7495 ns/op	    4192 B/op	      14 allocs/op
BenchmarkGorillaMux_GPlusParam    	  282151	      4069 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGowwwRouter_GPlusParam   	 1364444	       891 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_GPlusParam    	14464407	        82.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GPlusParam   	 1598323	       767 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_GPlusParam         	 3630046	       334 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_GPlusParam          	11398248	       105 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusParam       	  540242	      3177 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_GPlusParam       	  244903	      5486 ns/op	    1072 B/op	      10 allocs/op
BenchmarkPat_GPlusParam           	  876195	      2053 ns/op	     576 B/op	      11 allocs/op
BenchmarkPossum_GPlusParam        	 1000000	      1406 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_GPlusParam      	 1382778	       871 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_GPlusParam         	 4733049	       259 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_GPlusParam         	 1000000	      1305 ns/op	     256 B/op	       7 allocs/op
BenchmarkTigerTonic_GPlusParam    	  635600	      3036 ns/op	     856 B/op	      16 allocs/op
BenchmarkTraffic_GPlusParam       	  252438	      5520 ns/op	    1872 B/op	      21 allocs/op
BenchmarkVulcan_GPlusParam        	 1576790	       765 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_GPlus2Params         	 8219248	       147 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlus2Params        	10168023	       118 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlus2Params        	 1000000	      1350 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_GPlus2Params       	 1000000	      1533 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GPlus2Params        	  276637	      4481 ns/op	    1184 B/op	      10 allocs/op
BenchmarkChi_GPlus2Params         	 1000000	      1371 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_GPlus2Params       	 3639152	       327 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_GPlus2Params        	 6661582	       180 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlus2Params         	 9361000	       128 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlus2Params  	 1000000	      1921 ns/op	     712 B/op	       9 allocs/op
BenchmarkGoji_GPlus2Params        	 1000000	      1142 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GPlus2Params      	  316833	      3882 ns/op	    1456 B/op	      14 allocs/op
BenchmarkGoJsonRest_GPlus2Params  	  835053	      2300 ns/op	     713 B/op	      14 allocs/op
BenchmarkGoRestful_GPlus2Params   	  163738	      8066 ns/op	    4384 B/op	      16 allocs/op
BenchmarkGorillaMux_GPlus2Params  	  166983	      7577 ns/op	    1328 B/op	      10 allocs/op
BenchmarkGowwwRouter_GPlus2Params 	 1337229	       926 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_GPlus2Params  	11574624	       103 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GPlus2Params 	 1000000	      1110 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_GPlus2Params       	 1868533	       650 ns/op	     128 B/op	       5 allocs/op
BenchmarkLARS_GPlus2Params        	 9156480	       131 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlus2Params     	  429186	      3382 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_GPlus2Params     	  114469	     10684 ns/op	    1200 B/op	      13 allocs/op
BenchmarkPat_GPlus2Params         	  168004	      7965 ns/op	    2168 B/op	      33 allocs/op
BenchmarkPossum_GPlus2Params      	 1000000	      1385 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_GPlus2Params    	 1279712	       981 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_GPlus2Params       	 3064262	       403 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_GPlus2Params       	 1000000	      1490 ns/op	     336 B/op	       7 allocs/op
BenchmarkTigerTonic_GPlus2Params  	  278316	      4775 ns/op	    1200 B/op	      22 allocs/op
BenchmarkTraffic_GPlus2Params     	  133639	     10417 ns/op	    2248 B/op	      28 allocs/op
BenchmarkVulcan_GPlus2Params      	 1000000	      1213 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_GPlusAll             	  717478	      1673 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusAll            	 1000000	      1161 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusAll            	   85567	     15313 ns/op	    5488 B/op	      61 allocs/op
BenchmarkBeego_GPlusAll           	   64579	     18277 ns/op	    4576 B/op	      39 allocs/op
BenchmarkBone_GPlusAll            	   31428	     38878 ns/op	   11920 B/op	     109 allocs/op
BenchmarkChi_GPlusAll             	   70621	     17240 ns/op	    5824 B/op	      39 allocs/op
BenchmarkDenco_GPlusAll           	  472082	      3651 ns/op	     672 B/op	      11 allocs/op
BenchmarkEcho_GPlusAll            	  618861	      2075 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_GPlusAll             	  833515	      1432 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusAll      	   55903	     22055 ns/op	    8040 B/op	     103 allocs/op
BenchmarkGoji_GPlusAll            	  110674	     11959 ns/op	    3696 B/op	      22 allocs/op
BenchmarkGojiv2_GPlusAll          	   28976	     42417 ns/op	   18240 B/op	     154 allocs/op
BenchmarkGoJsonRest_GPlusAll      	   47835	     26060 ns/op	    8117 B/op	     170 allocs/op
BenchmarkGoRestful_GPlusAll       	   10000	    102604 ns/op	   55520 B/op	     192 allocs/op
BenchmarkGorillaMux_GPlusAll      	   19669	     63014 ns/op	   16528 B/op	     128 allocs/op
BenchmarkGowwwRouter_GPlusAll     	  113444	     10431 ns/op	    4928 B/op	      33 allocs/op
BenchmarkHttpRouter_GPlusAll      	 1000000	      1099 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GPlusAll     	  125463	     11146 ns/op	    4032 B/op	      38 allocs/op
BenchmarkKocha_GPlusAll           	  237786	      5654 ns/op	     976 B/op	      43 allocs/op
BenchmarkLARS_GPlusAll            	  901672	      1369 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusAll         	   39888	     30349 ns/op	    9568 B/op	     104 allocs/op
BenchmarkMartini_GPlusAll         	   13508	     87842 ns/op	   14016 B/op	     145 allocs/op
BenchmarkPat_GPlusAll             	   21909	     54992 ns/op	   15264 B/op	     271 allocs/op
BenchmarkPossum_GPlusAll          	   92986	     14108 ns/op	    5408 B/op	      39 allocs/op
BenchmarkR2router_GPlusAll        	  102181	     12429 ns/op	    5040 B/op	      63 allocs/op
BenchmarkRivet_GPlusAll           	  444087	      4039 ns/op	     768 B/op	      11 allocs/op
BenchmarkTango_GPlusAll           	   71474	     17633 ns/op	    3552 B/op	      91 allocs/op
BenchmarkTigerTonic_GPlusAll      	   27124	     45531 ns/op	   11600 B/op	     242 allocs/op
BenchmarkTraffic_GPlusAll         	   13905	     85842 ns/op	   26248 B/op	     341 allocs/op
BenchmarkVulcan_GPlusAll          	   90640	     12404 ns/op	    1274 B/op	      39 allocs/op
```

## Parse.com

```sh
BenchmarkGin_ParseStatic          	15009693	        80.4 ns/op	       0 B/op	       0 allocs/op

BenchmarkAce_ParseStatic          	15957213	        76.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseStatic         	22928330	        52.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseStatic         	 2340534	       540 ns/op	     120 B/op	       3 allocs/op
BenchmarkBeego_ParseStatic        	 1000000	      1281 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_ParseStatic         	 1505856	       796 ns/op	     144 B/op	       3 allocs/op
BenchmarkChi_ParseStatic          	 1000000	      1097 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_ParseStatic        	34671934	        35.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_ParseStatic         	14155250	        85.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseStatic   	 1296708	       928 ns/op	     296 B/op	       5 allocs/op
BenchmarkGoji_ParseStatic         	 5961315	       201 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_ParseStatic       	  435138	      2663 ns/op	    1360 B/op	      10 allocs/op
BenchmarkGoJsonRest_ParseStatic   	 1000000	      1179 ns/op	     329 B/op	      11 allocs/op
BenchmarkGoRestful_ParseStatic    	  152389	      8508 ns/op	    4256 B/op	      13 allocs/op
BenchmarkGorillaMux_ParseStatic   	  502206	      2842 ns/op	    1008 B/op	       9 allocs/op
BenchmarkGowwwRouter_ParseStatic  	32175411	        37.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_ParseStatic   	40009770	        30.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_ParseStatic  	18656181	        64.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_ParseStatic        	25166415	        48.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_ParseStatic         	16176346	        74.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseStatic      	  924776	      2274 ns/op	     736 B/op	       8 allocs/op
BenchmarkMartini_ParseStatic      	  335364	      4575 ns/op	     768 B/op	       9 allocs/op
BenchmarkPat_ParseStatic          	 1271038	       928 ns/op	     240 B/op	       5 allocs/op
BenchmarkPossum_ParseStatic       	 1234796	      1000 ns/op	     416 B/op	       3 allocs/op
BenchmarkR2router_ParseStatic     	 2423956	       452 ns/op	     144 B/op	       4 allocs/op
BenchmarkRivet_ParseStatic        	15375571	        77.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_ParseStatic        	 1000000	      1194 ns/op	     240 B/op	       7 allocs/op
BenchmarkTigerTonic_ParseStatic   	 4268598	       307 ns/op	      48 B/op	       1 allocs/op
BenchmarkTraffic_ParseStatic      	  384829	      3422 ns/op	    1256 B/op	      19 allocs/op
BenchmarkVulcan_ParseStatic       	 2023200	       598 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_ParseParam           	11223447	       109 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseParam          	14840814	        73.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseParam          	 1000000	      1114 ns/op	     467 B/op	       5 allocs/op
BenchmarkBeego_ParseParam         	 1000000	      1399 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_ParseParam          	  820448	      2300 ns/op	     912 B/op	       7 allocs/op
BenchmarkChi_ParseParam           	 1000000	      1195 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_ParseParam         	 4753371	       259 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_ParseParam          	11635582	       103 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseParam           	13721174	        88.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseParam    	 1000000	      1604 ns/op	     664 B/op	       8 allocs/op
BenchmarkGoji_ParseParam          	 1259976	       969 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_ParseParam        	  382359	      2962 ns/op	    1408 B/op	      12 allocs/op
BenchmarkGoJsonRest_ParseParam    	 1000000	      1913 ns/op	     649 B/op	      13 allocs/op
BenchmarkGoRestful_ParseParam     	  137076	      8919 ns/op	    4576 B/op	      14 allocs/op
BenchmarkGorillaMux_ParseParam    	  333202	      3315 ns/op	    1312 B/op	      10 allocs/op
BenchmarkGowwwRouter_ParseParam   	 1405658	       845 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_ParseParam    	19034205	        64.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_ParseParam   	 1661428	       713 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_ParseParam         	 3950151	       308 ns/op	      56 B/op	       3 allocs/op
BenchmarkLARS_ParseParam          	14987967	        80.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseParam       	  541587	      3195 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_ParseParam       	  242671	      5145 ns/op	    1072 B/op	      10 allocs/op
BenchmarkPat_ParseParam           	  608030	      2838 ns/op	     992 B/op	      15 allocs/op
BenchmarkPossum_ParseParam        	 1000000	      1359 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_ParseParam      	 1371126	       893 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_ParseParam         	 4926261	       243 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_ParseParam         	 1000000	      1282 ns/op	     272 B/op	       7 allocs/op
BenchmarkTigerTonic_ParseParam    	  720100	      2750 ns/op	     784 B/op	      15 allocs/op
BenchmarkTraffic_ParseParam       	  284826	      4579 ns/op	    1896 B/op	      21 allocs/op
BenchmarkVulcan_ParseParam        	 1716049	       713 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_Parse2Params         	 9270141	       125 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_Parse2Params        	14511698	        83.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Parse2Params        	 1000000	      1328 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_Parse2Params       	 1000000	      1459 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Parse2Params        	 1000000	      2165 ns/op	     864 B/op	       6 allocs/op
BenchmarkChi_Parse2Params         	 1000000	      1348 ns/op	     448 B/op	       3 allocs/op
BenchmarkDenco_Parse2Params       	 4107226	       299 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_Parse2Params        	 8701603	       137 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_Parse2Params         	11868313	       102 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Parse2Params  	 1000000	      1956 ns/op	     712 B/op	       9 allocs/op
BenchmarkGoji_Parse2Params        	 1256845	       956 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Parse2Params      	  357972	      3040 ns/op	    1392 B/op	      11 allocs/op
BenchmarkGoJsonRest_Parse2Params  	  960513	      2154 ns/op	     713 B/op	      14 allocs/op
BenchmarkGoRestful_Parse2Params   	  129975	      9684 ns/op	    4928 B/op	      14 allocs/op
BenchmarkGorillaMux_Parse2Params  	  296415	      4013 ns/op	    1328 B/op	      10 allocs/op
BenchmarkGowwwRouter_Parse2Params 	 1341662	       893 ns/op	     448 B/op	       3 allocs/op
BenchmarkHttpRouter_Parse2Params  	14590942	        82.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_Parse2Params 	 1000000	      1021 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_Parse2Params       	 2072586	       587 ns/op	     128 B/op	       5 allocs/op
BenchmarkLARS_Parse2Params        	12111204	        98.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Parse2Params     	  450412	      3309 ns/op	    1072 B/op	      10 allocs/op
BenchmarkMartini_Parse2Params     	  232142	      5715 ns/op	    1152 B/op	      11 allocs/op
BenchmarkPat_Parse2Params         	  688216	      2687 ns/op	     752 B/op	      16 allocs/op
BenchmarkPossum_Parse2Params      	 1000000	      1355 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Parse2Params    	 1281830	       930 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Parse2Params       	 3238431	       359 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_Parse2Params       	 1000000	      1361 ns/op	     304 B/op	       7 allocs/op
BenchmarkTigerTonic_Parse2Params  	  341338	      4479 ns/op	    1168 B/op	      22 allocs/op
BenchmarkTraffic_Parse2Params     	  249236	      5482 ns/op	    1944 B/op	      22 allocs/op
BenchmarkVulcan_Parse2Params      	 1485796	       807 ns/op	      98 B/op	       3 allocs/op
BenchmarkAce_ParseAll             	  405536	      2981 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseAll            	  637077	      1879 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseAll            	   46346	     27078 ns/op	    8928 B/op	     110 allocs/op
BenchmarkBeego_ParseAll           	   34216	     33786 ns/op	    9152 B/op	      78 allocs/op
BenchmarkBone_ParseAll            	   23863	     51539 ns/op	   16464 B/op	     147 allocs/op
BenchmarkChi_ParseAll             	   36600	     33419 ns/op	   11648 B/op	      78 allocs/op
BenchmarkDenco_ParseAll           	  199340	      5305 ns/op	     928 B/op	      16 allocs/op
BenchmarkEcho_ParseAll            	  350238	      3353 ns/op	       0 B/op	       0 allocs/op
BenchmarkGin_ParseAll             	  487788	      2481 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseAll      	   31798	     40336 ns/op	   13728 B/op	     181 allocs/op
BenchmarkGoji_ParseAll            	   64731	     18814 ns/op	    5376 B/op	      32 allocs/op
BenchmarkGojiv2_ParseAll          	   15840	     78120 ns/op	   35696 B/op	     277 allocs/op
BenchmarkGoJsonRest_ParseAll      	   26872	     44424 ns/op	   13866 B/op	     321 allocs/op
BenchmarkGoRestful_ParseAll       	    5437	    235182 ns/op	  117600 B/op	     354 allocs/op
BenchmarkGorillaMux_ParseAll      	   10000	    120964 ns/op	   31120 B/op	     250 allocs/op
BenchmarkGowwwRouter_ParseAll     	   79274	     15462 ns/op	    7168 B/op	      48 allocs/op
BenchmarkHttpRouter_ParseAll      	  748496	      1573 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_ParseAll     	   78631	     15675 ns/op	    5728 B/op	      51 allocs/op
BenchmarkKocha_ParseAll           	  147765	      7598 ns/op	    1112 B/op	      54 allocs/op
BenchmarkLARS_ParseAll            	  520750	      2312 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseAll         	   20241	     58433 ns/op	   19136 B/op	     208 allocs/op
BenchmarkMartini_ParseAll         	   10000	    146952 ns/op	   25072 B/op	     253 allocs/op
BenchmarkPat_ParseAll             	   20721	     58816 ns/op	   15216 B/op	     308 allocs/op
BenchmarkPossum_ParseAll          	   45824	     27271 ns/op	   10816 B/op	      78 allocs/op
BenchmarkR2router_ParseAll        	   56907	     21497 ns/op	    8352 B/op	     120 allocs/op
BenchmarkRivet_ParseAll           	  238440	      5908 ns/op	     912 B/op	      16 allocs/op
BenchmarkTango_ParseAll           	   35584	     34872 ns/op	    6960 B/op	     182 allocs/op
BenchmarkTigerTonic_ParseAll      	   19269	     62525 ns/op	   16048 B/op	     332 allocs/op
BenchmarkTraffic_ParseAll         	   10000	    132594 ns/op	   45520 B/op	     605 allocs/op
BenchmarkVulcan_ParseAll          	   53130	     21850 ns/op	    2548 B/op	      78 allocs/op
```
