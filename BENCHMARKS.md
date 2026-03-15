# Gin Benchmark Report

**Machine:** Apple M4 Pro
**OS:** macOS (Darwin 25.3.0), arm64
**Date:** March 15th, 2026
**Gin Version:** v1.12.0
**Go Version:** 1.25.8 darwin/arm64
**Source:** [Go HTTP Router Benchmark](https://github.com/gin-gonic/go-http-routing-benchmark)

---

## Table of Contents

- [Memory Consumption](#memory-consumption)
- [Benchmark Summary](#benchmark-summary)
  - [GitHub API (203 routes)](#github-api-all-routes)
  - [Google+ API (13 routes)](#google-api-all-routes)
  - [Parse API (26 routes)](#parse-api-all-routes)
  - [Static Routes (157 routes)](#static-routes-all)
- [Micro Benchmarks](#micro-benchmarks)
- [Full Benchmark Results](#full-benchmark-results)

---

## Memory Consumption

Memory required for loading the routing structure (lower is better).

### Static Routes: 157

| Router | Bytes |
|:-------|------:|
| **Denco** | **9,928** |
| **Pat** | **20,200** |
| **HttpRouter** | **21,680** |
| GowwwRouter | 24,512 |
| R2router | 23,256 |
| Rivet | 24,608 |
| Aero | 28,728 |
| Tango | 28,296 |
| Bear | 29,984 |
| Goji | 27,216 |
| Ace | 30,648 |
| LARS | 30,640 |
| Gin | 34,408 |
| Macaron | 36,976 |
| Bone | 40,480 |
| GocraftWeb | 55,088 |
| HttpServeMux | 71,728 |
| HttpTreeMux | 73,448 |
| TigerTonic | 78,328 |
| Chi | 83,160 |
| Possum | 88,672 |
| Echo | 91,976 |
| Beego | 98,824 |
| Gojiv2 | 117,952 |
| Kocha | 123,160 |
| GoJsonRest | 136,296 |
| Martini | 314,360 |
| Vulcan | 367,064 |
| Traffic | 555,248 |
| GorillaMux | 599,496 |
| GoRestful | 819,704 |

### GitHub API Routes: 203

| Router | Bytes |
|:-------|------:|
| **Pat** | **19,400** |
| **Denco** | **35,360** |
| **HttpRouter** | **37,072** |
| Rivet | 42,840 |
| Goji | 46,416 |
| R2router | 46,616 |
| Ace | 48,616 |
| LARS | 48,592 |
| Tango | 54,872 |
| Gin | 58,840 |
| Bear | 80,528 |
| Possum | 81,832 |
| GowwwRouter | 86,656 |
| HttpTreeMux | 78,800 |
| Macaron | 90,632 |
| GocraftWeb | 93,896 |
| Chi | 94,888 |
| TigerTonic | 95,104 |
| Bone | 100,080 |
| Echo | 117,784 |
| Gojiv2 | 118,640 |
| GoJsonRest | 141,120 |
| Beego | 150,840 |
| Martini | 474,792 |
| Aero | 524,992 |
| Kocha | 784,144 |
| Traffic | 916,720 |
| GoRestful | 1,270,720 |
| GorillaMux | 1,319,696 |
| Vulcan | 423,216 |

### Google+ API Routes: 13

| Router | Bytes |
|:-------|------:|
| **Pat** | **1,848** |
| **HttpRouter** | **2,776** |
| **Goji** | **2,928** |
| Rivet | 3,064 |
| Denco | 3,224 |
| Ace | 3,680 |
| LARS | 3,656 |
| R2router | 3,864 |
| Gin | 4,576 |
| Tango | 5,200 |
| GowwwRouter | 6,168 |
| Bone | 6,688 |
| Bear | 7,096 |
| GocraftWeb | 7,480 |
| HttpTreeMux | 7,440 |
| Gojiv2 | 8,096 |
| Chi | 8,008 |
| Macaron | 8,672 |
| TigerTonic | 9,392 |
| Echo | 10,968 |
| Beego | 10,256 |
| GoJsonRest | 11,432 |
| Martini | 24,656 |
| Aero | 25,744 |
| Vulcan | 25,448 |
| Traffic | 48,176 |
| GorillaMux | 68,000 |
| GoRestful | 72,536 |
| Kocha | 128,880 |
| Possum | 7,192 |

### Parse API Routes: 26

| Router | Bytes |
|:-------|------:|
| **Pat** | **2,552** |
| **Denco** | **4,080** |
| **HttpRouter** | **5,024** |
| Goji | 5,248 |
| Rivet | 5,680 |
| Ace | 6,656 |
| LARS | 6,632 |
| R2router | 6,920 |
| Gin | 7,896 |
| HttpTreeMux | 7,848 |
| Possum | 8,104 |
| Bear | 12,256 |
| Tango | 8,952 |
| Chi | 9,656 |
| GowwwRouter | 10,048 |
| TigerTonic | 9,808 |
| Bone | 11,440 |
| GocraftWeb | 12,752 |
| Echo | 13,816 |
| Macaron | 13,704 |
| GoJsonRest | 14,112 |
| Gojiv2 | 16,064 |
| Beego | 19,256 |
| Aero | 28,152 |
| Martini | 43,760 |
| Traffic | 78,832 |
| Vulcan | 44,024 |
| GorillaMux | 105,384 |
| GoRestful | 121,200 |
| Kocha | 181,712 |

---

## Benchmark Summary

The tables below highlight the **top performers** across each major API benchmark.
Gin consistently ranks among the fastest routers with **zero memory allocations**.

### GitHub API (All Routes)

Routing all 203 GitHub API endpoints per operation.

| Rank | Router | Time (ns/op) | Bytes (B/op) | Allocs/op |
|:----:|:-------|-------------:|-------------:|----------:|
| 1 | **Aero** | 6,646 | 0 | 0 |
| 2 | **LARS** | 8,714 | 0 | 0 |
| 3 | **Gin** | 9,058 | 0 | 0 |
| 4 | Echo | 12,030 | 0 | 0 |
| 5 | HttpRouter | 13,840 | 13,792 | 167 |
| 6 | Ace | 21,354 | 13,792 | 167 |
| 7 | Denco | 22,401 | 20,224 | 167 |
| 8 | Rivet | 23,337 | 16,272 | 167 |
| 9 | Kocha | 33,213 | 20,592 | 504 |
| 10 | GowwwRouter | 35,511 | 61,456 | 334 |
| ... | ... | ... | ... | ... |
| 28 | GorillaMux | 1,282,012 | 225,666 | 1,588 |
| 29 | Pat | 1,432,291 | 1,421,792 | 23,019 |
| 30 | Martini | 1,483,537 | 231,418 | 2,731 |

### Google+ API (All Routes)

Routing all 13 Google+ API endpoints per operation.

| Rank | Router | Time (ns/op) | Bytes (B/op) | Allocs/op |
|:----:|:-------|-------------:|-------------:|----------:|
| 1 | **Aero** | 323.2 | 0 | 0 |
| 2 | **LARS** | 382.6 | 0 | 0 |
| 3 | **Gin** | 412.3 | 0 | 0 |
| 4 | Echo | 426.0 | 0 | 0 |
| 5 | HttpRouter | 649.1 | 640 | 11 |
| 6 | Denco | 910.9 | 672 | 11 |
| 7 | Ace | 912.1 | 640 | 11 |
| 8 | Rivet | 928.1 | 768 | 11 |
| 9 | Kocha | 1,323 | 848 | 27 |
| 10 | GowwwRouter | 1,858 | 4,048 | 22 |
| ... | ... | ... | ... | ... |
| 28 | Traffic | 22,350 | 26,616 | 366 |
| 29 | GoRestful | 23,142 | 60,720 | 193 |
| 30 | Martini | 24,600 | 14,328 | 171 |

### Parse API (All Routes)

Routing all 26 Parse API endpoints per operation.

| Rank | Router | Time (ns/op) | Bytes (B/op) | Allocs/op |
|:----:|:-------|-------------:|-------------:|----------:|
| 1 | **Aero** | 546.0 | 0 | 0 |
| 2 | **LARS** | 669.4 | 0 | 0 |
| 3 | **Gin** | 719.8 | 0 | 0 |
| 4 | Echo | 723.0 | 0 | 0 |
| 5 | HttpRouter | 921.5 | 640 | 16 |
| 6 | Denco | 1,291 | 928 | 16 |
| 7 | Rivet | 1,385 | 912 | 16 |
| 8 | Ace | 1,445 | 640 | 16 |
| 9 | Kocha | 1,813 | 960 | 35 |
| 10 | GowwwRouter | 2,924 | 5,888 | 32 |
| ... | ... | ... | ... | ... |
| 28 | Traffic | 34,164 | 45,760 | 630 |
| 29 | Martini | 39,785 | 25,696 | 305 |
| 30 | GoRestful | 52,633 | 131,728 | 380 |

### Static Routes (All)

Routing all 157 static routes per operation.

| Rank | Router | Time (ns/op) | Bytes (B/op) | Allocs/op |
|:----:|:-------|-------------:|-------------:|----------:|
| 1 | **Denco** | 2,240 | 0 | 0 |
| 2 | **Aero** | 2,872 | 0 | 0 |
| 3 | **HttpTreeMux** | 4,031 | 0 | 0 |
| 4 | HttpRouter | 4,095 | 0 | 0 |
| 5 | LARS | 5,272 | 0 | 0 |
| 6 | Gin | 5,566 | 0 | 0 |
| 7 | Ace | 5,981 | 0 | 0 |
| 8 | Kocha | 6,013 | 0 | 0 |
| 9 | GowwwRouter | 6,622 | 0 | 0 |
| 10 | Echo | 6,714 | 0 | 0 |
| ... | ... | ... | ... | ... |
| 28 | Martini | 577,268 | 129,209 | 2,031 |
| 29 | Traffic | 699,470 | 749,837 | 14,444 |
| 30 | Pat | 813,157 | 602,832 | 12,559 |

---

## Micro Benchmarks

### Single Parameter

```
BenchmarkGin_Param                	50441539	        24.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_Param                	24363132	        50.00 ns/op	      32 B/op	       1 allocs/op
BenchmarkAero_Param               	77342221	        15.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param               	 4424601	       271.1 ns/op	     456 B/op	       5 allocs/op
BenchmarkBeego_Param              	 3209640	       381.4 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Param               	 2987578	       432.5 ns/op	     752 B/op	       5 allocs/op
BenchmarkChi_Param                	 3111004	       358.4 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_Param              	28276489	        40.28 ns/op	      32 B/op	       1 allocs/op
BenchmarkEcho_Param               	65189638	        18.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param         	 3556015	       340.4 ns/op	     624 B/op	       7 allocs/op
BenchmarkGoji_Param               	 6142947	       198.9 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Param             	 2446671	       540.1 ns/op	    1136 B/op	       8 allocs/op
BenchmarkGoJsonRest_Param         	 2730858	       419.8 ns/op	     617 B/op	      13 allocs/op
BenchmarkGoRestful_Param          	  784556	      1549 ns/op	    4600 B/op	      15 allocs/op
BenchmarkGorillaMux_Param         	 1857164	       650.0 ns/op	    1152 B/op	       8 allocs/op
BenchmarkGowwwRouter_Param        	 8161994	       147.8 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_Param         	35229486	        31.71 ns/op	      32 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Param        	 6598603	       186.5 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_Param              	18681621	        61.98 ns/op	      48 B/op	       2 allocs/op
BenchmarkLARS_Param               	60355720	        20.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param            	 1669972	       725.6 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_Param            	  993966	      1343 ns/op	    1096 B/op	      12 allocs/op
BenchmarkPat_Param                	 2232280	       481.9 ns/op	     552 B/op	      12 allocs/op
BenchmarkPossum_Param             	 3422949	       349.7 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Param           	 6040192	       198.3 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_Param              	25981898	        46.35 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_Param              	 4682972	       249.3 ns/op	     192 B/op	       6 allocs/op
BenchmarkTigerTonic_Param         	 1715374	       688.0 ns/op	     744 B/op	      16 allocs/op
BenchmarkTraffic_Param            	  938739	      1240 ns/op	    1888 B/op	      23 allocs/op
BenchmarkVulcan_Param             	 9498867	       124.1 ns/op	      98 B/op	       3 allocs/op
```

### 5 Parameters

```
BenchmarkGin_Param5               	30261015	        38.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_Param5               	11764478	       100.2 ns/op	     160 B/op	       1 allocs/op
BenchmarkAero_Param5              	33756498	        34.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param5              	 3066855	       387.5 ns/op	     501 B/op	       5 allocs/op
BenchmarkBeego_Param5             	 2459563	       495.6 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Param5              	 1000000	      1042 ns/op	    1232 B/op	      10 allocs/op
BenchmarkChi_Param5               	 2523159	       469.0 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_Param5             	 8588504	       130.9 ns/op	     160 B/op	       1 allocs/op
BenchmarkEcho_Param5              	22098444	        49.58 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param5        	 2255988	       522.1 ns/op	     920 B/op	      11 allocs/op
BenchmarkGoji_Param5              	 4107170	       283.3 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Param5            	 1691694	       706.7 ns/op	    1248 B/op	      10 allocs/op
BenchmarkGoJsonRest_Param5        	 1579816	       736.0 ns/op	    1073 B/op	      16 allocs/op
BenchmarkGoRestful_Param5         	  495588	      2283 ns/op	    5400 B/op	      15 allocs/op
BenchmarkGorillaMux_Param5        	  734078	      1488 ns/op	    1184 B/op	       8 allocs/op
BenchmarkGowwwRouter_Param5       	 7023253	       171.0 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_Param5        	16128580	        71.82 ns/op	     160 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Param5       	 2799802	       421.4 ns/op	     576 B/op	       6 allocs/op
BenchmarkKocha_Param5             	 5397548	       216.4 ns/op	     304 B/op	       6 allocs/op
BenchmarkLARS_Param5              	31268811	        38.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param5           	 1246534	       934.2 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_Param5           	  523628	      2384 ns/op	    1256 B/op	      13 allocs/op
BenchmarkPat_Param5               	  726936	      1592 ns/op	     944 B/op	      32 allocs/op
BenchmarkPossum_Param5            	 3320758	       356.0 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Param5          	 3726697	       312.7 ns/op	     432 B/op	       5 allocs/op
BenchmarkRivet_Param5             	 8050082	       145.0 ns/op	     240 B/op	       1 allocs/op
BenchmarkTango_Param5             	 3265436	       366.0 ns/op	     392 B/op	       6 allocs/op
BenchmarkTigerTonic_Param5        	  469060	      2632 ns/op	    2408 B/op	      41 allocs/op
BenchmarkTraffic_Param5           	  425434	      2722 ns/op	    2280 B/op	      31 allocs/op
BenchmarkVulcan_Param5            	 5213118	       241.1 ns/op	      98 B/op	       3 allocs/op
```

### 20 Parameters

```
BenchmarkGin_Param20              	13291521	        89.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_Param20              	 4432400	       265.2 ns/op	     640 B/op	       1 allocs/op
BenchmarkAero_Param20             	10534245	       111.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Param20             	 1000000	      1027 ns/op	    1665 B/op	       5 allocs/op
BenchmarkBeego_Param20            	 1398088	       856.2 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Param20             	  275770	      4416 ns/op	    4368 B/op	      32 allocs/op
BenchmarkChi_Param20              	 1000000	      1111 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_Param20            	 2717076	       430.2 ns/op	     640 B/op	       1 allocs/op
BenchmarkEcho_Param20             	 7671782	       152.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Param20       	  672906	      1633 ns/op	    3795 B/op	      15 allocs/op
BenchmarkGoji_Param20             	 1000000	      1001 ns/op	    1247 B/op	       2 allocs/op
BenchmarkGojiv2_Param20           	  760102	      1528 ns/op	    1568 B/op	      12 allocs/op
BenchmarkGoJsonRest_Param20       	  516908	      2099 ns/op	    3539 B/op	      24 allocs/op
BenchmarkGoRestful_Param20        	  196784	      5920 ns/op	    8520 B/op	      15 allocs/op
BenchmarkGorillaMux_Param20       	  307082	      3793 ns/op	    3483 B/op	      10 allocs/op
BenchmarkGowwwRouter_Param20      	 5098538	       254.1 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_Param20       	 5192563	       230.2 ns/op	     640 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Param20      	  471916	      3085 ns/op	    3168 B/op	      10 allocs/op
BenchmarkKocha_Param20            	 1429558	       831.6 ns/op	    1264 B/op	      12 allocs/op
BenchmarkLARS_Param20             	13176316	        90.49 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Param20          	  498850	      2459 ns/op	    2680 B/op	      12 allocs/op
BenchmarkMartini_Param20          	  124274	     10258 ns/op	    3668 B/op	      15 allocs/op
BenchmarkPat_Param20              	  454987	      2592 ns/op	    1640 B/op	      44 allocs/op
BenchmarkPossum_Param20           	 3251416	       383.8 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Param20         	  908521	      1291 ns/op	    2176 B/op	       7 allocs/op
BenchmarkRivet_Param20            	 1000000	      1027 ns/op	    1024 B/op	       1 allocs/op
BenchmarkTango_Param20            	  992810	      1123 ns/op	     952 B/op	       6 allocs/op
BenchmarkTigerTonic_Param20       	  121791	      9949 ns/op	   10248 B/op	     111 allocs/op
BenchmarkTraffic_Param20          	  119523	     10037 ns/op	    7784 B/op	      56 allocs/op
BenchmarkVulcan_Param20           	 2427007	       504.5 ns/op	      98 B/op	       3 allocs/op
```

### Param Write

```
BenchmarkGin_ParamWrite           	23316990	        50.84 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_ParamWrite           	13474160	        92.09 ns/op	      40 B/op	       2 allocs/op
BenchmarkAero_ParamWrite          	34476930	        33.00 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParamWrite          	 3806036	       310.0 ns/op	     456 B/op	       5 allocs/op
BenchmarkBeego_ParamWrite         	 2834390	       419.9 ns/op	     360 B/op	       4 allocs/op
BenchmarkBone_ParamWrite          	 2627498	       466.8 ns/op	     752 B/op	       5 allocs/op
BenchmarkChi_ParamWrite           	 3128284	       371.5 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_ParamWrite         	16741923	        71.49 ns/op	      32 B/op	       1 allocs/op
BenchmarkEcho_ParamWrite          	29455958	        37.94 ns/op	       8 B/op	       1 allocs/op
BenchmarkGocraftWeb_ParamWrite    	 2992538	       393.9 ns/op	     648 B/op	       8 allocs/op
BenchmarkGoji_ParamWrite          	 4816736	       250.7 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_ParamWrite        	 2044604	       590.6 ns/op	    1152 B/op	       9 allocs/op
BenchmarkGoJsonRest_ParamWrite    	 2149220	       544.9 ns/op	    1064 B/op	      18 allocs/op
BenchmarkGoRestful_ParamWrite     	  633780	      1853 ns/op	    4616 B/op	      16 allocs/op
BenchmarkGorillaMux_ParamWrite    	 1652036	       727.3 ns/op	    1160 B/op	       9 allocs/op
BenchmarkGowwwRouter_ParamWrite   	 4195808	       279.1 ns/op	     392 B/op	       3 allocs/op
BenchmarkHttpRouter_ParamWrite    	21606003	        54.14 ns/op	      32 B/op	       1 allocs/op
BenchmarkHttpTreeMux_ParamWrite   	 5398851	       224.9 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_ParamWrite         	14068706	        83.18 ns/op	      48 B/op	       2 allocs/op
BenchmarkLARS_ParamWrite          	25768928	        47.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParamWrite       	 1355762	       894.2 ns/op	    1152 B/op	      13 allocs/op
BenchmarkMartini_ParamWrite       	  514706	      2544 ns/op	    1224 B/op	      16 allocs/op
BenchmarkPat_ParamWrite           	 1285564	       880.5 ns/op	    1040 B/op	      17 allocs/op
BenchmarkPossum_ParamWrite        	 3127676	       383.7 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_ParamWrite      	 4948783	       234.5 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_ParamWrite         	 9561326	       125.0 ns/op	     112 B/op	       2 allocs/op
BenchmarkTango_ParamWrite         	 3855218	       308.3 ns/op	     280 B/op	       7 allocs/op
BenchmarkTigerTonic_ParamWrite    	 1000000	      1090 ns/op	    1024 B/op	      18 allocs/op
BenchmarkTraffic_ParamWrite       	  742009	      1523 ns/op	    2248 B/op	      25 allocs/op
BenchmarkVulcan_ParamWrite        	 6929622	       170.7 ns/op	      98 B/op	       3 allocs/op
```

---

## Full Benchmark Results

### GitHub API

#### Static

```
BenchmarkGin_GithubStatic         	43895277	        28.10 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GithubStatic         	42779731	        28.31 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GithubStatic        	83284414	        14.50 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubStatic        	 9311720	       131.2 ns/op	     120 B/op	       3 allocs/op
BenchmarkBeego_GithubStatic       	 2897301	       395.7 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GithubStatic        	  308800	      3903 ns/op	    2880 B/op	      60 allocs/op
BenchmarkChi_GithubStatic         	 5401872	       229.6 ns/op	     368 B/op	       2 allocs/op
BenchmarkDenco_GithubStatic       	98722773	        11.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GithubStatic        	40625690	        27.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubStatic  	 5341020	       229.8 ns/op	     288 B/op	       5 allocs/op
BenchmarkGoji_GithubStatic        	15548287	        76.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_GithubStatic      	 2520218	       490.3 ns/op	    1120 B/op	       7 allocs/op
BenchmarkGoJsonRest_GithubStatic  	 3861292	       321.3 ns/op	     297 B/op	      11 allocs/op
BenchmarkGoRestful_GithubStatic   	  330211	      3671 ns/op	    4792 B/op	      14 allocs/op
BenchmarkGorillaMux_GithubStatic  	  938425	      1474 ns/op	     848 B/op	       7 allocs/op
BenchmarkGowwwRouter_GithubStatic 	38324576	        32.24 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GithubStatic  	69731800	        17.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GithubStatic 	55391751	        22.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_GithubStatic       	52384971	        24.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_GithubStatic        	45034478	        25.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubStatic     	 2014342	       578.0 ns/op	     728 B/op	       8 allocs/op
BenchmarkMartini_GithubStatic     	  554347	      2686 ns/op	     792 B/op	      11 allocs/op
BenchmarkPat_GithubStatic         	  288867	      4049 ns/op	    3648 B/op	      76 allocs/op
BenchmarkPossum_GithubStatic      	 3928141	       306.4 ns/op	     416 B/op	       3 allocs/op
BenchmarkR2router_GithubStatic    	11773146	        97.06 ns/op	     112 B/op	       3 allocs/op
BenchmarkRivet_GithubStatic       	37206849	        32.03 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_GithubStatic       	 3806214	       321.1 ns/op	     192 B/op	       6 allocs/op
BenchmarkTigerTonic_GithubStatic  	15421950	        82.96 ns/op	      48 B/op	       1 allocs/op
BenchmarkTraffic_GithubStatic     	  307382	      3618 ns/op	    4632 B/op	      89 allocs/op
BenchmarkVulcan_GithubStatic      	 6107649	       198.2 ns/op	      98 B/op	       3 allocs/op
```

#### Param

```
BenchmarkGin_GithubParam          	25399108	        50.29 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GithubParam          	12033406	       102.1 ns/op	      96 B/op	       1 allocs/op
BenchmarkAero_GithubParam         	34218457	        35.12 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubParam         	 3491167	       340.3 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_GithubParam        	 2462625	       492.9 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GithubParam         	  479170	      2435 ns/op	    1824 B/op	      18 allocs/op
BenchmarkChi_GithubParam          	 2392147	       439.7 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_GithubParam        	10566573	       108.5 ns/op	     128 B/op	       1 allocs/op
BenchmarkEcho_GithubParam         	22438311	        53.46 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubParam   	 2686350	       425.8 ns/op	     656 B/op	       7 allocs/op
BenchmarkGoji_GithubParam         	 3468987	       331.6 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GithubParam       	 1713819	       693.3 ns/op	    1216 B/op	      10 allocs/op
BenchmarkGoJsonRest_GithubParam   	 2105818	       561.9 ns/op	     681 B/op	      14 allocs/op
BenchmarkGoRestful_GithubParam    	  268806	      4732 ns/op	    4696 B/op	      15 allocs/op
BenchmarkGorillaMux_GithubParam   	  484622	      2436 ns/op	    1168 B/op	       8 allocs/op
BenchmarkGowwwRouter_GithubParam  	 6445003	       195.4 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_GithubParam   	14665429	        88.59 ns/op	      96 B/op	       1 allocs/op
BenchmarkHttpTreeMux_GithubParam  	 4199426	       311.2 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_GithubParam        	 8478230	       139.3 ns/op	     112 B/op	       3 allocs/op
BenchmarkLARS_GithubParam         	27976429	        43.36 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubParam      	 1475805	       805.7 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_GithubParam      	  354076	      3459 ns/op	    1176 B/op	      13 allocs/op
BenchmarkPat_GithubParam          	  422454	      2919 ns/op	    2360 B/op	      45 allocs/op
BenchmarkPossum_GithubParam       	 3345200	       349.5 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_GithubParam     	 5887364	       205.1 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_GithubParam        	 9056540	       118.9 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_GithubParam        	 3217730	       376.9 ns/op	     296 B/op	       6 allocs/op
BenchmarkTigerTonic_GithubParam   	 1000000	      1043 ns/op	    1072 B/op	      21 allocs/op
BenchmarkTraffic_GithubParam      	  375636	      3171 ns/op	    2840 B/op	      43 allocs/op
BenchmarkVulcan_GithubParam       	 3770610	       321.0 ns/op	      98 B/op	       3 allocs/op
```

#### All Routes

```
BenchmarkGin_GithubAll            	  135811	      9058 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GithubAll            	   61230	     21354 ns/op	   13792 B/op	     167 allocs/op
BenchmarkAero_GithubAll           	  176716	      6646 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GithubAll           	   16941	     69598 ns/op	   86448 B/op	     943 allocs/op
BenchmarkBeego_GithubAll          	   12163	     97593 ns/op	   71456 B/op	     609 allocs/op
BenchmarkBone_GithubAll           	    1176	   1040342 ns/op	  709472 B/op	    8453 allocs/op
BenchmarkChi_GithubAll            	   12346	     98943 ns/op	  130816 B/op	     740 allocs/op
BenchmarkDenco_GithubAll          	   59050	     22401 ns/op	   20224 B/op	     167 allocs/op
BenchmarkEcho_GithubAll           	  104632	     12030 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GithubAll     	   13927	     85845 ns/op	  123552 B/op	    1400 allocs/op
BenchmarkGoji_GithubAll           	    8512	    158755 ns/op	   56112 B/op	     334 allocs/op
BenchmarkGojiv2_GithubAll         	    5108	    226986 ns/op	  313744 B/op	    3712 allocs/op
BenchmarkGoJsonRest_GithubAll     	    9859	    120437 ns/op	  127875 B/op	    2737 allocs/op
BenchmarkGoRestful_GithubAll      	    1359	    863892 ns/op	 1006744 B/op	    3009 allocs/op
BenchmarkGorillaMux_GithubAll     	     974	   1282012 ns/op	  225666 B/op	    1588 allocs/op
BenchmarkGowwwRouter_GithubAll    	   32976	     35511 ns/op	   61456 B/op	     334 allocs/op
BenchmarkHttpRouter_GithubAll     	   89264	     13840 ns/op	   13792 B/op	     167 allocs/op
BenchmarkHttpTreeMux_GithubAll    	   25005	     49430 ns/op	   65856 B/op	     671 allocs/op
BenchmarkKocha_GithubAll          	   38254	     33213 ns/op	   20592 B/op	     504 allocs/op
BenchmarkLARS_GithubAll           	  144225	      8714 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GithubAll        	    8706	    131574 ns/op	  147784 B/op	    1624 allocs/op
BenchmarkMartini_GithubAll        	     825	   1483537 ns/op	  231418 B/op	    2731 allocs/op
BenchmarkPat_GithubAll            	     838	   1432291 ns/op	 1421792 B/op	   23019 allocs/op
BenchmarkPossum_GithubAll         	   18913	     62291 ns/op	   84448 B/op	     609 allocs/op
BenchmarkR2router_GithubAll       	   26276	     45938 ns/op	   70832 B/op	     776 allocs/op
BenchmarkRivet_GithubAll          	   49792	     23337 ns/op	   16272 B/op	     167 allocs/op
BenchmarkTango_GithubAll          	   14306	     84094 ns/op	   53850 B/op	    1215 allocs/op
BenchmarkTigerTonic_GithubAll     	    5797	    209921 ns/op	  188584 B/op	    4300 allocs/op
BenchmarkTraffic_GithubAll        	    1044	   1213864 ns/op	  829175 B/op	   14582 allocs/op
BenchmarkVulcan_GithubAll         	   19022	     60788 ns/op	   19894 B/op	     609 allocs/op
```

### Google+ API

#### Static

```
BenchmarkGin_GPlusStatic          	50430675	        22.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GPlusStatic          	51961077	        23.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_GPlusStatic         	86196060	        13.97 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusStatic         	11848935	        99.99 ns/op	     104 B/op	       3 allocs/op
BenchmarkBeego_GPlusStatic        	 3394909	       360.2 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GPlusStatic         	21614740	        54.72 ns/op	      32 B/op	       1 allocs/op
BenchmarkChi_GPlusStatic          	 5905422	       202.2 ns/op	     368 B/op	       2 allocs/op
BenchmarkDenco_GPlusStatic        	100000000	        11.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_GPlusStatic         	63868430	        18.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusStatic   	 6419139	       183.5 ns/op	     272 B/op	       5 allocs/op
BenchmarkGoji_GPlusStatic         	26074863	        46.67 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_GPlusStatic       	 2529218	       474.7 ns/op	    1120 B/op	       7 allocs/op
BenchmarkGoJsonRest_GPlusStatic   	 4662106	       254.6 ns/op	     297 B/op	      11 allocs/op
BenchmarkGoRestful_GPlusStatic    	  782607	      1489 ns/op	    4272 B/op	      14 allocs/op
BenchmarkGorillaMux_GPlusStatic   	 2741367	       417.2 ns/op	     848 B/op	       7 allocs/op
BenchmarkGowwwRouter_GPlusStatic  	100000000	        11.61 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_GPlusStatic   	123925137	         9.716 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_GPlusStatic  	84271881	        14.55 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_GPlusStatic        	96201384	        12.78 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_GPlusStatic         	60936904	        19.79 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusStatic      	 2281216	       523.1 ns/op	     728 B/op	       8 allocs/op
BenchmarkMartini_GPlusStatic      	 1000000	      1158 ns/op	     792 B/op	      11 allocs/op
BenchmarkPat_GPlusStatic          	12473952	        94.32 ns/op	      96 B/op	       2 allocs/op
BenchmarkPossum_GPlusStatic       	 4341033	       283.0 ns/op	     416 B/op	       3 allocs/op
BenchmarkR2router_GPlusStatic     	14272594	        80.44 ns/op	     112 B/op	       3 allocs/op
BenchmarkRivet_GPlusStatic        	64791607	        18.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_GPlusStatic        	 5208178	       231.0 ns/op	     152 B/op	       6 allocs/op
BenchmarkTigerTonic_GPlusStatic   	26052549	        45.54 ns/op	      32 B/op	       1 allocs/op
BenchmarkTraffic_GPlusStatic      	 1655068	       664.6 ns/op	    1080 B/op	      15 allocs/op
BenchmarkVulcan_GPlusStatic       	10144974	       117.5 ns/op	      98 B/op	       3 allocs/op
```

#### Param

```
BenchmarkGin_GPlusParam           	38126504	        31.22 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GPlusParam           	15315734	        70.21 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_GPlusParam          	55314088	        22.34 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusParam          	 3841524	       285.5 ns/op	     472 B/op	       5 allocs/op
BenchmarkBeego_GPlusParam         	 2792283	       414.8 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GPlusParam          	 2710533	       440.2 ns/op	     752 B/op	       5 allocs/op
BenchmarkChi_GPlusParam           	 3291210	       369.3 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_GPlusParam         	19789339	        64.83 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_GPlusParam          	40679114	        30.04 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusParam    	 3499014	       343.7 ns/op	     624 B/op	       7 allocs/op
BenchmarkGoji_GPlusParam          	 5488542	       232.6 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GPlusParam        	 2221576	       543.0 ns/op	    1136 B/op	       8 allocs/op
BenchmarkGoJsonRest_GPlusParam    	 2624800	       447.0 ns/op	     617 B/op	      13 allocs/op
BenchmarkGoRestful_GPlusParam     	  728508	      1712 ns/op	    4616 B/op	      15 allocs/op
BenchmarkGorillaMux_GPlusParam    	 1370820	       887.2 ns/op	    1152 B/op	       8 allocs/op
BenchmarkGowwwRouter_GPlusParam   	 7895398	       147.6 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_GPlusParam    	25334516	        48.01 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_GPlusParam   	 5918101	       198.7 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_GPlusParam         	15647581	        77.64 ns/op	      48 B/op	       2 allocs/op
BenchmarkLARS_GPlusParam          	41545491	        28.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusParam       	 1619770	       742.6 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_GPlusParam       	  927007	      1412 ns/op	    1096 B/op	      12 allocs/op
BenchmarkPat_GPlusParam           	 2187026	       545.2 ns/op	     600 B/op	      12 allocs/op
BenchmarkPossum_GPlusParam        	 3570588	       341.9 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_GPlusParam      	 6142230	       195.6 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_GPlusParam         	21743265	        56.36 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_GPlusParam         	 4237848	       285.6 ns/op	     216 B/op	       6 allocs/op
BenchmarkTigerTonic_GPlusParam    	 1625415	       752.7 ns/op	     824 B/op	      16 allocs/op
BenchmarkTraffic_GPlusParam       	  852801	      1389 ns/op	    1904 B/op	      23 allocs/op
BenchmarkVulcan_GPlusParam        	 7105036	       167.8 ns/op	      98 B/op	       3 allocs/op
```

#### 2 Params

```
BenchmarkGin_GPlus2Params         	30045693	        39.13 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GPlus2Params         	15562932	        79.43 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_GPlus2Params        	33729063	        34.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlus2Params        	 3764409	       318.7 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_GPlus2Params       	 2423888	       494.5 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_GPlus2Params        	 1000000	      1090 ns/op	    1104 B/op	       9 allocs/op
BenchmarkChi_GPlus2Params         	 3077922	       392.4 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_GPlus2Params       	14066016	        85.80 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_GPlus2Params        	29279718	        41.98 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlus2Params  	 3111852	       369.1 ns/op	     656 B/op	       7 allocs/op
BenchmarkGoji_GPlus2Params        	 3919934	       298.6 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_GPlus2Params      	 1724162	       698.0 ns/op	    1216 B/op	      11 allocs/op
BenchmarkGoJsonRest_GPlus2Params  	 2149428	       557.3 ns/op	     681 B/op	      14 allocs/op
BenchmarkGoRestful_GPlus2Params   	  573634	      1893 ns/op	    4712 B/op	      15 allocs/op
BenchmarkGorillaMux_GPlus2Params  	  682111	      1824 ns/op	    1168 B/op	       8 allocs/op
BenchmarkGowwwRouter_GPlus2Params 	 7580608	       156.0 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_GPlus2Params  	20394765	        59.07 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_GPlus2Params 	 5053549	       238.5 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_GPlus2Params       	 8661169	       137.7 ns/op	     112 B/op	       3 allocs/op
BenchmarkLARS_GPlus2Params        	32620809	        37.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlus2Params     	 1484382	       828.6 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_GPlus2Params     	  375764	      3118 ns/op	    1224 B/op	      15 allocs/op
BenchmarkPat_GPlus2Params         	  593236	      1964 ns/op	    2136 B/op	      31 allocs/op
BenchmarkPossum_GPlus2Params      	 3535449	       338.7 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_GPlus2Params    	 6382530	       189.1 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_GPlus2Params       	13734822	        86.79 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_GPlus2Params       	 3824863	       306.1 ns/op	     296 B/op	       6 allocs/op
BenchmarkTigerTonic_GPlus2Params  	 1000000	      1078 ns/op	    1152 B/op	      21 allocs/op
BenchmarkTraffic_GPlus2Params     	  468492	      2531 ns/op	    2296 B/op	      31 allocs/op
BenchmarkVulcan_GPlus2Params      	 4796979	       248.9 ns/op	      98 B/op	       3 allocs/op
```

#### All Routes

```
BenchmarkGin_GPlusAll             	 2893287	       412.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_GPlusAll             	 1329721	       912.1 ns/op	     640 B/op	      11 allocs/op
BenchmarkAero_GPlusAll            	 3728048	       323.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_GPlusAll            	  310212	      3637 ns/op	    5488 B/op	      61 allocs/op
BenchmarkBeego_GPlusAll           	  215176	      5681 ns/op	    4576 B/op	      39 allocs/op
BenchmarkBone_GPlusAll            	  117496	     10217 ns/op	   11040 B/op	      98 allocs/op
BenchmarkChi_GPlusAll             	  238798	      4876 ns/op	    8480 B/op	      48 allocs/op
BenchmarkDenco_GPlusAll           	 1313409	       910.9 ns/op	     672 B/op	      11 allocs/op
BenchmarkEcho_GPlusAll            	 2705866	       426.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_GPlusAll      	  263034	      4459 ns/op	    7600 B/op	      87 allocs/op
BenchmarkGoji_GPlusAll            	  399655	      2891 ns/op	    3696 B/op	      22 allocs/op
BenchmarkGojiv2_GPlusAll          	  154381	      7584 ns/op	   15120 B/op	     115 allocs/op
BenchmarkGoJsonRest_GPlusAll      	  181494	      6304 ns/op	    7701 B/op	     170 allocs/op
BenchmarkGoRestful_GPlusAll       	   52003	     23142 ns/op	   60720 B/op	     193 allocs/op
BenchmarkGorillaMux_GPlusAll      	   87273	     13911 ns/op	   14448 B/op	     102 allocs/op
BenchmarkGowwwRouter_GPlusAll     	  574188	      1858 ns/op	    4048 B/op	      22 allocs/op
BenchmarkHttpRouter_GPlusAll      	 1877773	       649.1 ns/op	     640 B/op	      11 allocs/op
BenchmarkHttpTreeMux_GPlusAll     	  473188	      2373 ns/op	    4032 B/op	      38 allocs/op
BenchmarkKocha_GPlusAll           	 1000000	      1323 ns/op	     848 B/op	      27 allocs/op
BenchmarkLARS_GPlusAll            	 3073754	       382.6 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_GPlusAll         	  178746	      6983 ns/op	    9464 B/op	     104 allocs/op
BenchmarkMartini_GPlusAll         	   48831	     24600 ns/op	   14328 B/op	     171 allocs/op
BenchmarkPat_GPlusAll             	   76422	     15692 ns/op	   15272 B/op	     269 allocs/op
BenchmarkPossum_GPlusAll          	  316477	      3773 ns/op	    5408 B/op	      39 allocs/op
BenchmarkR2router_GPlusAll        	  520474	      2423 ns/op	    4624 B/op	      50 allocs/op
BenchmarkRivet_GPlusAll           	 1295100	       928.1 ns/op	     768 B/op	      11 allocs/op
BenchmarkTango_GPlusAll           	  298516	      3843 ns/op	    3008 B/op	      78 allocs/op
BenchmarkTigerTonic_GPlusAll      	  101194	     11505 ns/op	   11160 B/op	     237 allocs/op
BenchmarkTraffic_GPlusAll         	   53414	     22350 ns/op	   26616 B/op	     366 allocs/op
BenchmarkVulcan_GPlusAll          	  463639	      2531 ns/op	    1274 B/op	      39 allocs/op
```

### Parse API

#### Static

```
BenchmarkGin_ParseStatic          	53518070	        22.62 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_ParseStatic          	52354974	        23.33 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_ParseStatic         	87150981	        13.87 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseStatic         	 9958563	       121.5 ns/op	     120 B/op	       3 allocs/op
BenchmarkBeego_ParseStatic        	 3357967	       355.5 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_ParseStatic         	 6466970	       185.2 ns/op	     144 B/op	       3 allocs/op
BenchmarkChi_ParseStatic          	 5862630	       198.8 ns/op	     368 B/op	       2 allocs/op
BenchmarkDenco_ParseStatic        	100000000	        12.05 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_ParseStatic         	62838466	        18.40 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseStatic   	 6035592	       197.0 ns/op	     288 B/op	       5 allocs/op
BenchmarkGoji_ParseStatic         	18933330	        63.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_ParseStatic       	 2645055	       447.3 ns/op	    1120 B/op	       7 allocs/op
BenchmarkGoJsonRest_ParseStatic   	 4464429	       271.8 ns/op	     297 B/op	      11 allocs/op
BenchmarkGoRestful_ParseStatic    	  653388	      1786 ns/op	    4792 B/op	      14 allocs/op
BenchmarkGorillaMux_ParseStatic   	 2294398	       508.1 ns/op	     848 B/op	       7 allocs/op
BenchmarkGowwwRouter_ParseStatic  	85855849	        13.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_ParseStatic   	100000000	        10.71 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_ParseStatic  	60829818	        20.59 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_ParseStatic        	86461821	        13.88 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_ParseStatic         	60537139	        19.89 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseStatic      	 2327062	       511.8 ns/op	     728 B/op	       8 allocs/op
BenchmarkMartini_ParseStatic      	 1000000	      1223 ns/op	     792 B/op	      11 allocs/op
BenchmarkPat_ParseStatic          	 5077684	       234.6 ns/op	     240 B/op	       5 allocs/op
BenchmarkPossum_ParseStatic       	 4378826	       276.9 ns/op	     416 B/op	       3 allocs/op
BenchmarkR2router_ParseStatic     	13918509	        86.61 ns/op	     112 B/op	       3 allocs/op
BenchmarkRivet_ParseStatic        	59392341	        19.74 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_ParseStatic        	 4877336	       253.6 ns/op	     192 B/op	       6 allocs/op
BenchmarkTigerTonic_ParseStatic   	17529486	        68.49 ns/op	      48 B/op	       1 allocs/op
BenchmarkTraffic_ParseStatic      	 1513068	       781.4 ns/op	    1224 B/op	      18 allocs/op
BenchmarkVulcan_ParseStatic       	 8955946	       134.6 ns/op	      98 B/op	       3 allocs/op
```

#### Param

```
BenchmarkGin_ParseParam           	46613631	        25.30 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_ParseParam           	19263414	        62.50 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_ParseParam          	65469127	        18.25 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseParam          	 4148041	       282.4 ns/op	     467 B/op	       5 allocs/op
BenchmarkBeego_ParseParam         	 3173386	       382.8 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_ParseParam          	 2519083	       472.7 ns/op	     832 B/op	       6 allocs/op
BenchmarkChi_ParseParam           	 3320446	       368.9 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_ParseParam         	20528611	        56.96 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_ParseParam          	53122124	        22.99 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseParam    	 3454389	       346.5 ns/op	     640 B/op	       7 allocs/op
BenchmarkGoji_ParseParam          	 5256088	       235.8 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_ParseParam        	 2296279	       519.1 ns/op	    1168 B/op	       9 allocs/op
BenchmarkGoJsonRest_ParseParam    	 2735005	       433.4 ns/op	     617 B/op	      13 allocs/op
BenchmarkGoRestful_ParseParam     	  559504	      1927 ns/op	    5112 B/op	      15 allocs/op
BenchmarkGorillaMux_ParseParam    	 1812528	       634.8 ns/op	    1152 B/op	       8 allocs/op
BenchmarkGowwwRouter_ParseParam   	 8530036	       142.9 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_ParseParam    	29872574	        40.20 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_ParseParam   	 6568846	       185.7 ns/op	     352 B/op	       3 allocs/op
BenchmarkKocha_ParseParam         	16294209	        71.62 ns/op	      48 B/op	       2 allocs/op
BenchmarkLARS_ParseParam          	56620798	        21.32 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseParam       	 1664428	       732.1 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_ParseParam       	  981326	      1318 ns/op	    1096 B/op	      12 allocs/op
BenchmarkPat_ParseParam           	 1675594	       712.1 ns/op	     992 B/op	      15 allocs/op
BenchmarkPossum_ParseParam        	 3543426	       336.8 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_ParseParam      	 6078027	       198.4 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_ParseParam         	23278633	        52.40 ns/op	      48 B/op	       1 allocs/op
BenchmarkTango_ParseParam         	 4534116	       268.0 ns/op	     224 B/op	       6 allocs/op
BenchmarkTigerTonic_ParseParam    	 1706557	       715.8 ns/op	     752 B/op	      15 allocs/op
BenchmarkTraffic_ParseParam       	  963790	      1191 ns/op	    1928 B/op	      23 allocs/op
BenchmarkVulcan_ParseParam        	 8465382	       139.8 ns/op	      98 B/op	       3 allocs/op
```

#### 2 Params

```
BenchmarkGin_Parse2Params         	39753472	        30.02 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_Parse2Params         	18333952	        68.63 ns/op	      64 B/op	       1 allocs/op
BenchmarkAero_Parse2Params        	52058101	        22.95 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_Parse2Params        	 3673090	       321.8 ns/op	     496 B/op	       5 allocs/op
BenchmarkBeego_Parse2Params       	 2765089	       431.9 ns/op	     352 B/op	       3 allocs/op
BenchmarkBone_Parse2Params        	 2676722	       458.0 ns/op	     784 B/op	       5 allocs/op
BenchmarkChi_Parse2Params         	 3197473	       378.5 ns/op	     704 B/op	       4 allocs/op
BenchmarkDenco_Parse2Params       	17205883	        68.76 ns/op	      64 B/op	       1 allocs/op
BenchmarkEcho_Parse2Params        	43142511	        28.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_Parse2Params  	 3212872	       372.8 ns/op	     656 B/op	       7 allocs/op
BenchmarkGoji_Parse2Params        	 5678137	       208.5 ns/op	     336 B/op	       2 allocs/op
BenchmarkGojiv2_Parse2Params      	 2286504	       527.8 ns/op	    1152 B/op	       8 allocs/op
BenchmarkGoJsonRest_Parse2Params  	 2387936	       501.1 ns/op	     681 B/op	      14 allocs/op
BenchmarkGoRestful_Parse2Params   	  507566	      2179 ns/op	    5536 B/op	      15 allocs/op
BenchmarkGorillaMux_Parse2Params  	 1516276	       791.5 ns/op	    1168 B/op	       8 allocs/op
BenchmarkGowwwRouter_Parse2Params 	 7861905	       149.0 ns/op	     368 B/op	       2 allocs/op
BenchmarkHttpRouter_Parse2Params  	25059865	        47.14 ns/op	      64 B/op	       1 allocs/op
BenchmarkHttpTreeMux_Parse2Params 	 5467455	       218.1 ns/op	     384 B/op	       4 allocs/op
BenchmarkKocha_Parse2Params       	10478869	       114.7 ns/op	     112 B/op	       3 allocs/op
BenchmarkLARS_Parse2Params        	44204389	        26.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_Parse2Params     	 1570063	       758.0 ns/op	    1064 B/op	      10 allocs/op
BenchmarkMartini_Parse2Params     	  885129	      1407 ns/op	    1176 B/op	      13 allocs/op
BenchmarkPat_Parse2Params         	 1701120	       706.4 ns/op	     768 B/op	      17 allocs/op
BenchmarkPossum_Parse2Params      	 3545026	       341.2 ns/op	     496 B/op	       5 allocs/op
BenchmarkR2router_Parse2Params    	 6219794	       199.8 ns/op	     400 B/op	       4 allocs/op
BenchmarkRivet_Parse2Params       	15446226	        77.50 ns/op	      96 B/op	       1 allocs/op
BenchmarkTango_Parse2Params       	 4123059	       285.0 ns/op	     264 B/op	       6 allocs/op
BenchmarkTigerTonic_Parse2Params  	 1000000	      1053 ns/op	    1120 B/op	      21 allocs/op
BenchmarkTraffic_Parse2Params     	  870003	      1372 ns/op	    1992 B/op	      25 allocs/op
BenchmarkVulcan_Parse2Params      	 7048611	       168.3 ns/op	      98 B/op	       3 allocs/op
```

#### All Routes

```
BenchmarkGin_ParseAll             	 1681276	       719.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_ParseAll             	 1000000	      1445 ns/op	     640 B/op	      16 allocs/op
BenchmarkAero_ParseAll            	 2210972	       546.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_ParseAll            	  179937	      6452 ns/op	    8920 B/op	     110 allocs/op
BenchmarkBeego_ParseAll           	  114976	     10182 ns/op	    9152 B/op	      78 allocs/op
BenchmarkBone_ParseAll            	   99656	     11868 ns/op	   15184 B/op	     131 allocs/op
BenchmarkChi_ParseAll             	  137953	      8453 ns/op	   14944 B/op	      84 allocs/op
BenchmarkDenco_ParseAll           	 1000000	      1291 ns/op	     928 B/op	      16 allocs/op
BenchmarkEcho_ParseAll            	 1656105	       723.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_ParseAll      	  142503	      8165 ns/op	   13168 B/op	     162 allocs/op
BenchmarkGoji_ParseAll            	  240613	      4792 ns/op	    5376 B/op	      32 allocs/op
BenchmarkGojiv2_ParseAll          	   85312	     13700 ns/op	   29456 B/op	     199 allocs/op
BenchmarkGoJsonRest_ParseAll      	  105716	     11300 ns/op	   13034 B/op	     321 allocs/op
BenchmarkGoRestful_ParseAll       	   22760	     52633 ns/op	  131728 B/op	     380 allocs/op
BenchmarkGorillaMux_ParseAll      	   46082	     26107 ns/op	   26960 B/op	     198 allocs/op
BenchmarkGowwwRouter_ParseAll     	  363032	      2924 ns/op	    5888 B/op	      32 allocs/op
BenchmarkHttpRouter_ParseAll      	 1292341	       921.5 ns/op	     640 B/op	      16 allocs/op
BenchmarkHttpTreeMux_ParseAll     	  322304	      3484 ns/op	    5728 B/op	      51 allocs/op
BenchmarkKocha_ParseAll           	  834361	      1813 ns/op	     960 B/op	      35 allocs/op
BenchmarkLARS_ParseAll            	 1779954	       669.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_ParseAll         	   87348	     13595 ns/op	   18928 B/op	     208 allocs/op
BenchmarkMartini_ParseAll         	   30399	     39785 ns/op	   25696 B/op	     305 allocs/op
BenchmarkPat_ParseAll             	   70246	     17221 ns/op	   15376 B/op	     323 allocs/op
BenchmarkPossum_ParseAll          	  159716	      7413 ns/op	   10816 B/op	      78 allocs/op
BenchmarkR2router_ParseAll        	  258717	      4338 ns/op	    7520 B/op	      94 allocs/op
BenchmarkRivet_ParseAll           	 1000000	      1385 ns/op	     912 B/op	      16 allocs/op
BenchmarkTango_ParseAll           	  157684	      7464 ns/op	    5864 B/op	     156 allocs/op
BenchmarkTigerTonic_ParseAll      	   76174	     16031 ns/op	   15488 B/op	     329 allocs/op
BenchmarkTraffic_ParseAll         	   34507	     34164 ns/op	   45760 B/op	     630 allocs/op
BenchmarkVulcan_ParseAll          	  253556	      5019 ns/op	    2548 B/op	      78 allocs/op
```

### Static Routes

```
BenchmarkGin_StaticAll            	  214844	      5566 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpServeMux_StaticAll   	   63081	     18896 ns/op	       0 B/op	       0 allocs/op
BenchmarkAce_StaticAll            	  199526	      5981 ns/op	       0 B/op	       0 allocs/op
BenchmarkAero_StaticAll           	  420883	      2872 ns/op	       0 B/op	       0 allocs/op
BenchmarkBear_StaticAll           	   54873	     22053 ns/op	   19960 B/op	     469 allocs/op
BenchmarkBeego_StaticAll          	   17143	     69371 ns/op	   55264 B/op	     471 allocs/op
BenchmarkBone_StaticAll           	   50510	     23765 ns/op	       0 B/op	       0 allocs/op
BenchmarkChi_StaticAll            	   29125	     41246 ns/op	   57776 B/op	     314 allocs/op
BenchmarkDenco_StaticAll          	  528363	      2240 ns/op	       0 B/op	       0 allocs/op
BenchmarkEcho_StaticAll           	  167440	      6714 ns/op	       0 B/op	       0 allocs/op
BenchmarkGocraftWeb_StaticAll     	   32882	     36485 ns/op	   45056 B/op	     785 allocs/op
BenchmarkGoji_StaticAll           	   65331	     18211 ns/op	       0 B/op	       0 allocs/op
BenchmarkGojiv2_StaticAll         	   13905	     86141 ns/op	  175840 B/op	    1099 allocs/op
BenchmarkGoJsonRest_StaticAll     	   20342	     59899 ns/op	   46629 B/op	    1727 allocs/op
BenchmarkGoRestful_StaticAll      	    2685	    447104 ns/op	  677824 B/op	    2193 allocs/op
BenchmarkGorillaMux_StaticAll     	    3967	    300555 ns/op	  133137 B/op	    1099 allocs/op
BenchmarkGowwwRouter_StaticAll    	  179246	      6622 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpRouter_StaticAll     	  291186	      4095 ns/op	       0 B/op	       0 allocs/op
BenchmarkHttpTreeMux_StaticAll    	  295596	      4031 ns/op	       0 B/op	       0 allocs/op
BenchmarkKocha_StaticAll          	  198813	      6013 ns/op	       0 B/op	       0 allocs/op
BenchmarkLARS_StaticAll           	  220419	      5272 ns/op	       0 B/op	       0 allocs/op
BenchmarkMacaron_StaticAll        	   14395	     82544 ns/op	  114296 B/op	    1256 allocs/op
BenchmarkMartini_StaticAll        	    2061	    577268 ns/op	  129209 B/op	    2031 allocs/op
BenchmarkPat_StaticAll            	    1453	    813157 ns/op	  602832 B/op	   12559 allocs/op
BenchmarkPossum_StaticAll         	   26097	     45967 ns/op	   65312 B/op	     471 allocs/op
BenchmarkR2router_StaticAll       	   61281	     19281 ns/op	   17584 B/op	     471 allocs/op
BenchmarkRivet_StaticAll          	  132847	      9082 ns/op	       0 B/op	       0 allocs/op
BenchmarkTango_StaticAll          	   20103	     59732 ns/op	   31274 B/op	     942 allocs/op
BenchmarkTigerTonic_StaticAll     	   83690	     14425 ns/op	    7376 B/op	     157 allocs/op
BenchmarkTraffic_StaticAll        	    1669	    699470 ns/op	  749837 B/op	   14444 allocs/op
BenchmarkVulcan_StaticAll         	   31063	     38764 ns/op	   15386 B/op	     471 allocs/op
```
