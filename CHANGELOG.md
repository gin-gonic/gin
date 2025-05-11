# Gin ChangeLog

## Gin v1.10.0

### Features

* feat(auth): add proxy-server authentication (#3877) (@EndlessParadox1)
* feat(bind): ShouldBindBodyWith shortcut and change doc (#3871) (@RedCrazyGhost)
* feat(binding): Support custom BindUnmarshaler for binding. (#3933) (@dkkb)
* feat(binding): support override default binding implement (#3514) (@ssfyn)
* feat(engine): Added `OptionFunc` and `With` (#3572) (@flc1125)
* feat(logger): ability to skip logs based on user-defined logic (#3593) (@palvaneh)

### Bug fixes

* Revert "fix(uri): query binding bug (#3236)" (#3899) (@appleboy)
* fix(binding): binding error while not upload file (#3819) (#3820) (@clearcodecn)
* fix(binding): dereference pointer to struct (#3199) (@echovl)
* fix(context): make context Value method adhere to Go standards (#3897) (@FarmerChillax)
* fix(engine): fix unit test (#3878) (@flc1125)
* fix(header): Allow header according to RFC 7231 (HTTP 405) (#3759) (@Crocmagnon)
* fix(route): Add fullPath in context copy (#3784) (@KarthikReddyPuli)
* fix(router): catch-all conflicting wildcard (#3812) (@FirePing32)
* fix(sec): upgrade golang.org/x/crypto to 0.17.0 (#3832) (@chncaption)
* fix(tree): correctly expand the capacity of params (#3502) (@georgijd-form3)
* fix(uri): query binding bug (#3236) (@illiafox)
* fix: Add pointer support for url query params (#3659) (#3666) (@omkar-foss)
* fix: protect Context.Keys map when call Copy method (#3873) (@kingcanfish)
 
### Enhancements

* chore(CI): update release args (#3595) (@qloog)
* chore(IP): add TrustedPlatform constant for Fly.io. (#3839) (@ab)
* chore(debug): add ability to override the debugPrint statement (#2337) (@josegonzalez)
* chore(deps): update dependencies to latest versions (#3835) (@appleboy)
* chore(header): Add support for RFC 9512: application/yaml (#3851) (@vincentbernat)
* chore(http): use white color for HTTP 1XX (#3741) (@viralparmarme)
* chore(optimize): the ShouldBindUri method of the Context struct (#3911) (@1911860538)
* chore(perf): Optimize the Copy method of the Context struct (#3859) (@1911860538)
* chore(refactor): modify interface check way (#3855) (@demoManito)
* chore(request): check reader if it's nil before reading (#3419) (@noahyao1024)
* chore(security): upgrade Protobuf for CVE-2024-24786 (#3893) (@Fotkurz)
* chore: refactor CI and update dependencies (#3848) (@appleboy)
* chore: refactor configuration files for better readability (#3951) (@appleboy)
* chore: update GitHub Actions configuration (#3792) (@appleboy)
* chore: update changelog categories and improve documentation (#3917) (@appleboy)
* chore: update dependencies to latest versions (#3694) (@appleboy)
* chore: update external dependencies to latest versions (#3950) (@appleboy)
* chore: update various Go dependencies to latest versions (#3901) (@appleboy)

### Build process updates

* build(codecov): Added a codecov configuration (#3891) (@flc1125)
* ci(Makefile): vet command add .PHONY (#3915) (@imalasong)
* ci(lint): update tooling and workflows for consistency (#3834) (@appleboy)
* ci(release): refactor changelog regex patterns and exclusions (#3914) (@appleboy)
* ci(testing): add go1.22 version (#3842) (@appleboy)

### Documentation updates

* docs(context): Added deprecation comments to BindWith (#3880) (@flc1125)
* docs(middleware): comments to function `BasicAuthForProxy` (#3881) (@EndlessParadox1)
* docs: Add document  to constant `AuthProxyUserKey` and  `BasicAuthForProxy`. (#3887) (@EndlessParadox1)
* docs: fix typo in comment (#3868) (@testwill)
* docs: fix typo in function documentation (#3872) (@TotomiEcio)
* docs: remove redundant comments (#3765) (@WeiTheShinobi)
* feat: update version constant to v1.10.0 (#3952) (@appleboy)

### Others

* Upgrade golang.org/x/net -> v0.13.0 (#3684) (@cpcf)
* test(git): gitignore add develop tools (#3370) (@demoManito)
* test(http): use constant instead of numeric literal (#3863) (@testwill)
* test(path): Optimize unit test execution results (#3883) (@flc1125)
* test(render): increased unit tests coverage (#3691) (@araujo88)

## Gin v1.9.1

### BUG FIXES

* fix Request.Context() checks [#3512](https://github.com/gin-gonic/gin/pull/3512)

### SECURITY

* fix lack of escaping of filename in Content-Disposition [#3556](https://github.com/gin-gonic/gin/pull/3556) 

### ENHANCEMENTS

* refactor: use bytes.ReplaceAll directly [#3455](https://github.com/gin-gonic/gin/pull/3455)
* convert strings and slices using the officially recommended way [#3344](https://github.com/gin-gonic/gin/pull/3344)
* improve render code coverage [#3525](https://github.com/gin-gonic/gin/pull/3525)

### DOCS

* docs: changed documentation link for trusted proxies [#3575](https://github.com/gin-gonic/gin/pull/3575)
* chore: improve linting, testing, and GitHub Actions setup [#3583](https://github.com/gin-gonic/gin/pull/3583)

## Gin v1.9.0

### BREAK CHANGES

* Stop useless panicking in context and render [#2150](https://github.com/gin-gonic/gin/pull/2150)

### BUG FIXES

* fix(router): tree bug where loop index is not decremented. [#3460](https://github.com/gin-gonic/gin/pull/3460)
* fix(context): panic on NegotiateFormat - index out of range [#3397](https://github.com/gin-gonic/gin/pull/3397)
* Add escape logic for header [#3500](https://github.com/gin-gonic/gin/pull/3500) and [#3503](https://github.com/gin-gonic/gin/pull/3503)

### SECURITY

* Fix the GO-2022-0969 and GO-2022-0288 vulnerabilities [#3333](https://github.com/gin-gonic/gin/pull/3333)
* fix(security): vulnerability GO-2023-1571 [#3505](https://github.com/gin-gonic/gin/pull/3505)

### ENHANCEMENTS

* feat: add sonic json support [#3184](https://github.com/gin-gonic/gin/pull/3184)
* chore(file): Creates a directory named path [#3316](https://github.com/gin-gonic/gin/pull/3316)
* fix: modify interface check way [#3327](https://github.com/gin-gonic/gin/pull/3327)
* remove deprecated of package io/ioutil [#3395](https://github.com/gin-gonic/gin/pull/3395)
* refactor: avoid calling strings.ToLower twice [#3343](https://github.com/gin-gonic/gin/pull/3433)
* console logger HTTP status code bug fixed [#3453](https://github.com/gin-gonic/gin/pull/3453)
* chore(yaml): upgrade dependency to v3 version [#3456](https://github.com/gin-gonic/gin/pull/3456)
* chore(router): match method added to routergroup for multiple HTTP methods supporting [#3464](https://github.com/gin-gonic/gin/pull/3464)
* chore(http): add support for go1.20 http.rwUnwrapper to gin.responseWriter [#3489](https://github.com/gin-gonic/gin/pull/3489)

### DOCS

* docs: update markdown format [#3260](https://github.com/gin-gonic/gin/pull/3260)
* docs(readme): Add the TOML rendering example [#3400](https://github.com/gin-gonic/gin/pull/3400)
* docs(readme): move more example to docs/doc.md [#3449](https://github.com/gin-gonic/gin/pull/3449)
* docs: update markdown format [#3446](https://github.com/gin-gonic/gin/pull/3446)

## Gin v1.8.2

### BUG FIXES

* fix(route): redirectSlash bug ([#3227]((https://github.com/gin-gonic/gin/pull/3227)))
* fix(engine): missing route params for CreateTestContext ([#2778]((https://github.com/gin-gonic/gin/pull/2778))) ([#2803]((https://github.com/gin-gonic/gin/pull/2803)))

### SECURITY

* Fix the GO-2022-1144 vulnerability ([#3432]((https://github.com/gin-gonic/gin/pull/3432)))

## Gin v1.8.1

### ENHANCEMENTS

* feat(context): add ContextWithFallback feature flag [#3172](https://github.com/gin-gonic/gin/pull/3172)

## Gin v1.8.0

### BREAK CHANGES

* TrustedProxies: Add default IPv6 support and refactor [#2967](https://github.com/gin-gonic/gin/pull/2967). Please replace `RemoteIP() (net.IP, bool)` with `RemoteIP() net.IP`
* gin.Context with fallback value from gin.Context.Request.Context() [#2751](https://github.com/gin-gonic/gin/pull/2751)

### BUG FIXES

* Fixed SetOutput() panics on go 1.17 [#2861](https://github.com/gin-gonic/gin/pull/2861)
* Fix: wrong when wildcard follows named param [#2983](https://github.com/gin-gonic/gin/pull/2983)
* Fix: missing sameSite when do context.reset() [#3123](https://github.com/gin-gonic/gin/pull/3123)

### ENHANCEMENTS

* Use Header() instead of deprecated HeaderMap [#2694](https://github.com/gin-gonic/gin/pull/2694)
* RouterGroup.Handle regular match optimization of http method [#2685](https://github.com/gin-gonic/gin/pull/2685)
* Add support go-json, another drop-in json replacement [#2680](https://github.com/gin-gonic/gin/pull/2680)
* Use errors.New to replace fmt.Errorf will much better [#2707](https://github.com/gin-gonic/gin/pull/2707)
* Use Duration.Truncate for truncating precision [#2711](https://github.com/gin-gonic/gin/pull/2711)
* Get client IP when using Cloudflare [#2723](https://github.com/gin-gonic/gin/pull/2723)
* Optimize code adjust [#2700](https://github.com/gin-gonic/gin/pull/2700/files)
* Optimize code and reduce code cyclomatic complexity [#2737](https://github.com/gin-gonic/gin/pull/2737)
* Improve sliceValidateError.Error performance [#2765](https://github.com/gin-gonic/gin/pull/2765)
* Support custom struct tag [#2720](https://github.com/gin-gonic/gin/pull/2720)
* Improve router group tests [#2787](https://github.com/gin-gonic/gin/pull/2787)
* Fallback Context.Deadline() Context.Done() Context.Err() to Context.Request.Context() [#2769](https://github.com/gin-gonic/gin/pull/2769)
* Some codes optimize [#2830](https://github.com/gin-gonic/gin/pull/2830) [#2834](https://github.com/gin-gonic/gin/pull/2834) [#2838](https://github.com/gin-gonic/gin/pull/2838) [#2837](https://github.com/gin-gonic/gin/pull/2837) [#2788](https://github.com/gin-gonic/gin/pull/2788) [#2848](https://github.com/gin-gonic/gin/pull/2848) [#2851](https://github.com/gin-gonic/gin/pull/2851) [#2701](https://github.com/gin-gonic/gin/pull/2701)
* TrustedProxies: Add default IPv6 support and refactor [#2967](https://github.com/gin-gonic/gin/pull/2967)
* Test(route): expose performRequest func [#3012](https://github.com/gin-gonic/gin/pull/3012)
* Support h2c with prior knowledge [#1398](https://github.com/gin-gonic/gin/pull/1398)
* Feat attachment filename support utf8 [#3071](https://github.com/gin-gonic/gin/pull/3071)
* Feat: add StaticFileFS [#2749](https://github.com/gin-gonic/gin/pull/2749)
* Feat(context): return GIN Context from Value method [#2825](https://github.com/gin-gonic/gin/pull/2825)
* Feat: automatically SetMode to TestMode when run go test [#3139](https://github.com/gin-gonic/gin/pull/3139)
* Add TOML bining for gin [#3081](https://github.com/gin-gonic/gin/pull/3081)
* IPv6 add default trusted proxies [#3033](https://github.com/gin-gonic/gin/pull/3033)

### DOCS

* Add note about nomsgpack tag to the readme [#2703](https://github.com/gin-gonic/gin/pull/2703)

## Gin v1.7.7

### BUG FIXES

* Fixed X-Forwarded-For unsafe handling of CVE-2020-28483 [#2844](https://github.com/gin-gonic/gin/pull/2844), closed issue [#2862](https://github.com/gin-gonic/gin/issues/2862).
* Tree: updated the code logic for `latestNode` [#2897](https://github.com/gin-gonic/gin/pull/2897), closed issue [#2894](https://github.com/gin-gonic/gin/issues/2894) [#2878](https://github.com/gin-gonic/gin/issues/2878).
* Tree: fixed the misplacement of adding slashes [#2847](https://github.com/gin-gonic/gin/pull/2847), closed issue [#2843](https://github.com/gin-gonic/gin/issues/2843).
* Tree: fixed tsr with mixed static and wildcard paths [#2924](https://github.com/gin-gonic/gin/pull/2924), closed issue [#2918](https://github.com/gin-gonic/gin/issues/2918).

### ENHANCEMENTS

* TrustedProxies: make it backward-compatible [#2887](https://github.com/gin-gonic/gin/pull/2887), closed issue [#2819](https://github.com/gin-gonic/gin/issues/2819).
* TrustedPlatform: provide custom options for another CDN services [#2906](https://github.com/gin-gonic/gin/pull/2906).

### DOCS

* NoMethod: added usage annotation ([#2832](https://github.com/gin-gonic/gin/pull/2832#issuecomment-929954463)).

## Gin v1.7.6

### BUG FIXES

* bump new release to fix v1.7.5 release error by using v1.7.4 codes.

## Gin v1.7.4

### BUG FIXES

* bump new release to fix checksum mismatch

## Gin v1.7.3

### BUG FIXES

* fix level 1 router match [#2767](https://github.com/gin-gonic/gin/issues/2767), [#2796](https://github.com/gin-gonic/gin/issues/2796)

## Gin v1.7.2

### BUG FIXES

* Fix conflict between param and exact path [#2706](https://github.com/gin-gonic/gin/issues/2706). Close issue [#2682](https://github.com/gin-gonic/gin/issues/2682) [#2696](https://github.com/gin-gonic/gin/issues/2696).

## Gin v1.7.1

### BUG FIXES

* fix: data race with trustedCIDRs from [#2674](https://github.com/gin-gonic/gin/issues/2674)([#2675](https://github.com/gin-gonic/gin/pull/2675))

## Gin v1.7.0

### BUG FIXES

* fix compile error from [#2572](https://github.com/gin-gonic/gin/pull/2572) ([#2600](https://github.com/gin-gonic/gin/pull/2600))
* fix: print headers without Authorization header on broken pipe ([#2528](https://github.com/gin-gonic/gin/pull/2528))
* fix(tree): reassign fullpath when register new node ([#2366](https://github.com/gin-gonic/gin/pull/2366))

### ENHANCEMENTS

* Support params and exact routes without creating conflicts ([#2663](https://github.com/gin-gonic/gin/pull/2663))
* chore: improve render string performance ([#2365](https://github.com/gin-gonic/gin/pull/2365))
* Sync route tree to httprouter latest code ([#2368](https://github.com/gin-gonic/gin/pull/2368))
* chore: rename getQueryCache/getFormCache to initQueryCache/initFormCa ([#2375](https://github.com/gin-gonic/gin/pull/2375))
* chore(performance): improve countParams ([#2378](https://github.com/gin-gonic/gin/pull/2378))
* Remove some functions that have the same effect as the bytes package ([#2387](https://github.com/gin-gonic/gin/pull/2387))
* update:SetMode function ([#2321](https://github.com/gin-gonic/gin/pull/2321))
* remove an unused type SecureJSONPrefix ([#2391](https://github.com/gin-gonic/gin/pull/2391))
* Add a redirect sample for POST method ([#2389](https://github.com/gin-gonic/gin/pull/2389))
* Add CustomRecovery builtin middleware ([#2322](https://github.com/gin-gonic/gin/pull/2322))
* binding: avoid 2038 problem on 32-bit architectures ([#2450](https://github.com/gin-gonic/gin/pull/2450))
* Prevent panic in Context.GetQuery() when there is no Request ([#2412](https://github.com/gin-gonic/gin/pull/2412))
* Add GetUint and GetUint64 method on gin.context ([#2487](https://github.com/gin-gonic/gin/pull/2487))
* update content-disposition header to MIME-style ([#2512](https://github.com/gin-gonic/gin/pull/2512))
* reduce allocs and improve the render `WriteString` ([#2508](https://github.com/gin-gonic/gin/pull/2508))
* implement ".Unwrap() error" on Error type ([#2525](https://github.com/gin-gonic/gin/pull/2525)) ([#2526](https://github.com/gin-gonic/gin/pull/2526))
* Allow bind with a map[string]string ([#2484](https://github.com/gin-gonic/gin/pull/2484))
* chore: update tree ([#2371](https://github.com/gin-gonic/gin/pull/2371))
* Support binding for slice/array obj [Rewrite] ([#2302](https://github.com/gin-gonic/gin/pull/2302))
* basic auth: fix timing oracle ([#2609](https://github.com/gin-gonic/gin/pull/2609))
* Add mixed param and non-param paths (port of httprouter[#329](https://github.com/gin-gonic/gin/pull/329)) ([#2663](https://github.com/gin-gonic/gin/pull/2663))
* feat(engine): add trustedproxies and remoteIP ([#2632](https://github.com/gin-gonic/gin/pull/2632))

## Gin v1.6.3

### ENHANCEMENTS

  * Improve performance: Change `*sync.RWMutex` to `sync.RWMutex` in context. [#2351](https://github.com/gin-gonic/gin/pull/2351)

## Gin v1.6.2

### BUG FIXES

  * fix missing initial sync.RWMutex [#2305](https://github.com/gin-gonic/gin/pull/2305)

### ENHANCEMENTS

  * Add set samesite in cookie. [#2306](https://github.com/gin-gonic/gin/pull/2306)

## Gin v1.6.1

### BUG FIXES

  * Revert "fix accept incoming network connections" [#2294](https://github.com/gin-gonic/gin/pull/2294)

## Gin v1.6.0

### BREAKING

  * chore(performance): Improve performance for adding RemoveExtraSlash flag [#2159](https://github.com/gin-gonic/gin/pull/2159)
  * drop support govendor [#2148](https://github.com/gin-gonic/gin/pull/2148)
  * Added support for SameSite cookie flag [#1615](https://github.com/gin-gonic/gin/pull/1615)

### FEATURES

  * add yaml negotiation [#2220](https://github.com/gin-gonic/gin/pull/2220)
  * FileFromFS [#2112](https://github.com/gin-gonic/gin/pull/2112)

### BUG FIXES

  * Unix Socket Handling [#2280](https://github.com/gin-gonic/gin/pull/2280)
  * Use json marshall in context json to fix breaking new line issue. Fixes #2209 [#2228](https://github.com/gin-gonic/gin/pull/2228)
  * fix accept incoming network connections [#2216](https://github.com/gin-gonic/gin/pull/2216)
  * Fixed a bug in the calculation of the maximum number of parameters [#2166](https://github.com/gin-gonic/gin/pull/2166)
  * [FIX] allow empty headers on DataFromReader [#2121](https://github.com/gin-gonic/gin/pull/2121)
  * Add mutex for protect Context.Keys map [#1391](https://github.com/gin-gonic/gin/pull/1391)

### ENHANCEMENTS

  * Add mitigation for log injection [#2277](https://github.com/gin-gonic/gin/pull/2277)
  * tree: range over nodes values [#2229](https://github.com/gin-gonic/gin/pull/2229)
  * tree: remove duplicate assignment [#2222](https://github.com/gin-gonic/gin/pull/2222)
  * chore: upgrade go-isatty and json-iterator/go [#2215](https://github.com/gin-gonic/gin/pull/2215)
  * path: sync code with httprouter [#2212](https://github.com/gin-gonic/gin/pull/2212)
  * Use zero-copy approach to convert types between string and byte slice [#2206](https://github.com/gin-gonic/gin/pull/2206)
  * Reuse bytes when cleaning the URL paths [#2179](https://github.com/gin-gonic/gin/pull/2179)
  * tree: remove one else statement [#2177](https://github.com/gin-gonic/gin/pull/2177)
  * tree: sync httprouter update (#2173) (#2172) [#2171](https://github.com/gin-gonic/gin/pull/2171)
  * tree: sync part httprouter codes and reduce if/else [#2163](https://github.com/gin-gonic/gin/pull/2163)
  * use http method constant [#2155](https://github.com/gin-gonic/gin/pull/2155)
  * upgrade go-validator to v10 [#2149](https://github.com/gin-gonic/gin/pull/2149)
  * Refactor redirect request in gin.go [#1970](https://github.com/gin-gonic/gin/pull/1970)
  * Add build tag nomsgpack [#1852](https://github.com/gin-gonic/gin/pull/1852)

### DOCS

  * docs(path): improve comments [#2223](https://github.com/gin-gonic/gin/pull/2223)
  * Renew README to fit the modification of SetCookie method [#2217](https://github.com/gin-gonic/gin/pull/2217)
  * Fix spelling [#2202](https://github.com/gin-gonic/gin/pull/2202)
  * Remove broken link from README. [#2198](https://github.com/gin-gonic/gin/pull/2198)
  * Update docs on Context.Done(), Context.Deadline() and Context.Err() [#2196](https://github.com/gin-gonic/gin/pull/2196)
  * Update validator to v10 [#2190](https://github.com/gin-gonic/gin/pull/2190)
  * upgrade go-validator to v10 for README [#2189](https://github.com/gin-gonic/gin/pull/2189)
  * Update to currently output [#2188](https://github.com/gin-gonic/gin/pull/2188)
  * Fix "Custom Validators" example [#2186](https://github.com/gin-gonic/gin/pull/2186)
  * Add project to README [#2165](https://github.com/gin-gonic/gin/pull/2165)
  * docs(benchmarks): for gin v1.5 [#2153](https://github.com/gin-gonic/gin/pull/2153)
  * Changed wording for clarity in README.md [#2122](https://github.com/gin-gonic/gin/pull/2122)

### MISC

  * ci support go1.14 [#2262](https://github.com/gin-gonic/gin/pull/2262)
  * chore: upgrade depend version [#2231](https://github.com/gin-gonic/gin/pull/2231)
  * Drop support go1.10 [#2147](https://github.com/gin-gonic/gin/pull/2147)
  * fix comment in `mode.go` [#2129](https://github.com/gin-gonic/gin/pull/2129)

## Gin v1.5.0

- [FIX] Use DefaultWriter and DefaultErrorWriter for debug messages [#1891](https://github.com/gin-gonic/gin/pull/1891)
- [NEW] Now you can parse the inline lowercase start structure [#1893](https://github.com/gin-gonic/gin/pull/1893)
- [FIX] Some code improvements [#1909](https://github.com/gin-gonic/gin/pull/1909)
- [FIX] Use encode replace json marshal increase json encoder speed [#1546](https://github.com/gin-gonic/gin/pull/1546)
- [NEW] Hold matched route full path in the Context [#1826](https://github.com/gin-gonic/gin/pull/1826)
- [FIX] Fix context.Params race condition on Copy() [#1841](https://github.com/gin-gonic/gin/pull/1841)
- [NEW] Add context param query cache [#1450](https://github.com/gin-gonic/gin/pull/1450)
- [FIX] Improve GetQueryMap performance [#1918](https://github.com/gin-gonic/gin/pull/1918)
- [FIX] Improve get post data [#1920](https://github.com/gin-gonic/gin/pull/1920)
- [FIX] Use context instead of x/net/context [#1922](https://github.com/gin-gonic/gin/pull/1922)
- [FIX] Attempt to fix PostForm cache bug [#1931](https://github.com/gin-gonic/gin/pull/1931)
- [NEW] Add support of multipart multi files [#1949](https://github.com/gin-gonic/gin/pull/1949)
- [NEW] Support bind http header param [#1957](https://github.com/gin-gonic/gin/pull/1957)
- [FIX] Drop support for go1.8 and go1.9 [#1933](https://github.com/gin-gonic/gin/pull/1933)
- [FIX] Bugfix for the FullPath feature [#1919](https://github.com/gin-gonic/gin/pull/1919)
- [FIX] Gin1.5 bytes.Buffer to strings.Builder [#1939](https://github.com/gin-gonic/gin/pull/1939)
- [FIX] Upgrade github.com/ugorji/go/codec [#1969](https://github.com/gin-gonic/gin/pull/1969)
- [NEW] Support bind unix time [#1980](https://github.com/gin-gonic/gin/pull/1980)
- [FIX] Simplify code [#2004](https://github.com/gin-gonic/gin/pull/2004)
- [NEW] Support negative Content-Length in DataFromReader [#1981](https://github.com/gin-gonic/gin/pull/1981)
- [FIX] Identify terminal on a RISC-V architecture for auto-colored logs [#2019](https://github.com/gin-gonic/gin/pull/2019)
- [BREAKING] `Context.JSONP()` now expects a semicolon (`;`) at the end [#2007](https://github.com/gin-gonic/gin/pull/2007)
- [BREAKING] Upgrade default `binding.Validator` to v9 (see [its changelog](https://github.com/go-playground/validator/releases/tag/v9.0.0)) [#1015](https://github.com/gin-gonic/gin/pull/1015)
- [NEW] Add `DisallowUnknownFields()` in `Context.BindJSON()` [#2028](https://github.com/gin-gonic/gin/pull/2028)
- [NEW] Use specific `net.Listener` with `Engine.RunListener()` [#2023](https://github.com/gin-gonic/gin/pull/2023)
- [FIX] Fix some typo [#2079](https://github.com/gin-gonic/gin/pull/2079) [#2080](https://github.com/gin-gonic/gin/pull/2080)
- [FIX] Relocate binding body tests [#2086](https://github.com/gin-gonic/gin/pull/2086)
- [FIX] Use Writer in Context.Status [#1606](https://github.com/gin-gonic/gin/pull/1606)
- [FIX] `Engine.RunUnix()` now returns the error if it can't change the file mode [#2093](https://github.com/gin-gonic/gin/pull/2093)
- [FIX] `RouterGroup.StaticFS()` leaked files. Now it closes them. [#2118](https://github.com/gin-gonic/gin/pull/2118)
- [FIX] `Context.Request.FormFile` leaked file. Now it closes it. [#2114](https://github.com/gin-gonic/gin/pull/2114)
- [FIX] Ignore walking on `form:"-"` mapping [#1943](https://github.com/gin-gonic/gin/pull/1943)

### Gin v1.4.0

- [NEW] Support for [Go Modules](https://github.com/golang/go/wiki/Modules)  [#1569](https://github.com/gin-gonic/gin/pull/1569)
- [NEW] Refactor of form mapping multipart request [#1829](https://github.com/gin-gonic/gin/pull/1829)
- [FIX] Truncate Latency precision in long running request [#1830](https://github.com/gin-gonic/gin/pull/1830)
- [FIX] IsTerm flag should not be affected by DisableConsoleColor method. [#1802](https://github.com/gin-gonic/gin/pull/1802)
- [NEW] Supporting file binding [#1264](https://github.com/gin-gonic/gin/pull/1264)
- [NEW] Add support for mapping arrays [#1797](https://github.com/gin-gonic/gin/pull/1797)
- [FIX] Readme updates [#1793](https://github.com/gin-gonic/gin/pull/1793) [#1788](https://github.com/gin-gonic/gin/pull/1788) [1789](https://github.com/gin-gonic/gin/pull/1789)
- [FIX] StaticFS: Fixed Logging two log lines on 404.  [#1805](https://github.com/gin-gonic/gin/pull/1805), [#1804](https://github.com/gin-gonic/gin/pull/1804)
- [NEW] Make context.Keys available as LogFormatterParams [#1779](https://github.com/gin-gonic/gin/pull/1779)
- [NEW] Use internal/json for Marshal/Unmarshal [#1791](https://github.com/gin-gonic/gin/pull/1791)
- [NEW] Support mapping time.Duration [#1794](https://github.com/gin-gonic/gin/pull/1794)
- [NEW] Refactor form mappings [#1749](https://github.com/gin-gonic/gin/pull/1749)
- [NEW] Added flag to context.Stream indicates if client disconnected in middle of stream [#1252](https://github.com/gin-gonic/gin/pull/1252)
- [FIX] Moved [examples](https://github.com/gin-gonic/examples) to stand alone Repo [#1775](https://github.com/gin-gonic/gin/pull/1775)
- [NEW] Extend context.File to allow for the content-disposition attachments via a new method context.Attachment [#1260](https://github.com/gin-gonic/gin/pull/1260)
- [FIX] Support HTTP content negotiation wildcards [#1112](https://github.com/gin-gonic/gin/pull/1112)
- [NEW] Add prefix from X-Forwarded-Prefix in redirectTrailingSlash [#1238](https://github.com/gin-gonic/gin/pull/1238)
- [FIX] context.Copy() race condition [#1020](https://github.com/gin-gonic/gin/pull/1020)
- [NEW] Add context.HandlerNames() [#1729](https://github.com/gin-gonic/gin/pull/1729)
- [FIX] Change color methods to public in the defaultLogger. [#1771](https://github.com/gin-gonic/gin/pull/1771)
- [FIX] Update writeHeaders method to use http.Header.Set [#1722](https://github.com/gin-gonic/gin/pull/1722)
- [NEW] Add response size to LogFormatterParams [#1752](https://github.com/gin-gonic/gin/pull/1752)
- [NEW] Allow ignoring field on form mapping [#1733](https://github.com/gin-gonic/gin/pull/1733)
- [NEW] Add a function to force color in console output. [#1724](https://github.com/gin-gonic/gin/pull/1724)
- [FIX] Context.Next() - recheck len of handlers on every iteration. [#1745](https://github.com/gin-gonic/gin/pull/1745)
- [FIX] Fix all errcheck warnings [#1739](https://github.com/gin-gonic/gin/pull/1739) [#1653](https://github.com/gin-gonic/gin/pull/1653)
- [NEW] context: inherits context cancellation and deadline from http.Request context for Go>=1.7 [#1690](https://github.com/gin-gonic/gin/pull/1690)
- [NEW] Binding for URL Params [#1694](https://github.com/gin-gonic/gin/pull/1694)
- [NEW] Add LoggerWithFormatter method [#1677](https://github.com/gin-gonic/gin/pull/1677)
- [FIX] CI testing updates [#1671](https://github.com/gin-gonic/gin/pull/1671) [#1670](https://github.com/gin-gonic/gin/pull/1670) [#1682](https://github.com/gin-gonic/gin/pull/1682) [#1669](https://github.com/gin-gonic/gin/pull/1669)
- [FIX] StaticFS(): Send 404 when path does not exist [#1663](https://github.com/gin-gonic/gin/pull/1663)
- [FIX] Handle nil body for JSON binding [#1638](https://github.com/gin-gonic/gin/pull/1638)
- [FIX] Support bind uri param [#1612](https://github.com/gin-gonic/gin/pull/1612)
- [FIX] recovery: fix issue with syscall import on google app engine [#1640](https://github.com/gin-gonic/gin/pull/1640)
- [FIX] Make sure the debug log contains line breaks [#1650](https://github.com/gin-gonic/gin/pull/1650)
- [FIX] Panic stack trace being printed during recovery of broken pipe [#1089](https://github.com/gin-gonic/gin/pull/1089) [#1259](https://github.com/gin-gonic/gin/pull/1259)
- [NEW] RunFd method to run http.Server through a file descriptor [#1609](https://github.com/gin-gonic/gin/pull/1609)
- [NEW] Yaml binding support [#1618](https://github.com/gin-gonic/gin/pull/1618)
- [FIX] Pass MaxMultipartMemory when FormFile is called [#1600](https://github.com/gin-gonic/gin/pull/1600)
- [FIX] LoadHTML* tests [#1559](https://github.com/gin-gonic/gin/pull/1559)
- [FIX] Removed use of sync.pool from HandleContext [#1565](https://github.com/gin-gonic/gin/pull/1565)
- [FIX] Format output log to os.Stderr [#1571](https://github.com/gin-gonic/gin/pull/1571)
- [FIX] Make logger use a yellow background and a darkgray text for legibility [#1570](https://github.com/gin-gonic/gin/pull/1570)
- [FIX] Remove sensitive request information from panic log. [#1370](https://github.com/gin-gonic/gin/pull/1370)
- [FIX] log.Println() does not print timestamp [#829](https://github.com/gin-gonic/gin/pull/829) [#1560](https://github.com/gin-gonic/gin/pull/1560)
- [NEW] Add PureJSON renderer [#694](https://github.com/gin-gonic/gin/pull/694)
- [FIX] Add missing copyright and update if/else [#1497](https://github.com/gin-gonic/gin/pull/1497)
- [FIX] Update msgpack usage [#1498](https://github.com/gin-gonic/gin/pull/1498)
- [FIX] Use protobuf on render [#1496](https://github.com/gin-gonic/gin/pull/1496)
- [FIX] Add support for Protobuf format response [#1479](https://github.com/gin-gonic/gin/pull/1479)
- [NEW] Set default time format in form binding [#1487](https://github.com/gin-gonic/gin/pull/1487)
- [FIX] Add BindXML and ShouldBindXML [#1485](https://github.com/gin-gonic/gin/pull/1485)
- [NEW] Upgrade dependency libraries [#1491](https://github.com/gin-gonic/gin/pull/1491)


## Gin v1.3.0

- [NEW] Add [`func (*Context) QueryMap`](https://godoc.org/github.com/gin-gonic/gin#Context.QueryMap), [`func (*Context) GetQueryMap`](https://godoc.org/github.com/gin-gonic/gin#Context.GetQueryMap), [`func (*Context) PostFormMap`](https://godoc.org/github.com/gin-gonic/gin#Context.PostFormMap) and [`func (*Context) GetPostFormMap`](https://godoc.org/github.com/gin-gonic/gin#Context.GetPostFormMap) to support `type map[string]string` as query string or form parameters, see [#1383](https://github.com/gin-gonic/gin/pull/1383)
- [NEW] Add [`func (*Context) AsciiJSON`](https://godoc.org/github.com/gin-gonic/gin#Context.AsciiJSON), see [#1358](https://github.com/gin-gonic/gin/pull/1358)
- [NEW] Add `Pusher()` in [`type ResponseWriter`](https://godoc.org/github.com/gin-gonic/gin#ResponseWriter) for supporting http2 push, see [#1273](https://github.com/gin-gonic/gin/pull/1273)
- [NEW] Add [`func (*Context) DataFromReader`](https://godoc.org/github.com/gin-gonic/gin#Context.DataFromReader) for serving dynamic data, see [#1304](https://github.com/gin-gonic/gin/pull/1304)
- [NEW] Add [`func (*Context) ShouldBindBodyWith`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBindBodyWith) allowing to call binding multiple times, see [#1341](https://github.com/gin-gonic/gin/pull/1341)
- [NEW] Support pointers in form binding, see [#1336](https://github.com/gin-gonic/gin/pull/1336)
- [NEW] Add [`func (*Context) JSONP`](https://godoc.org/github.com/gin-gonic/gin#Context.JSONP), see [#1333](https://github.com/gin-gonic/gin/pull/1333)
- [NEW] Support default value in form binding, see [#1138](https://github.com/gin-gonic/gin/pull/1138)
- [NEW] Expose validator engine in [`type StructValidator`](https://godoc.org/github.com/gin-gonic/gin/binding#StructValidator), see [#1277](https://github.com/gin-gonic/gin/pull/1277)
- [NEW] Add [`func (*Context) ShouldBind`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBind), [`func (*Context) ShouldBindQuery`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBindQuery) and [`func (*Context) ShouldBindJSON`](https://godoc.org/github.com/gin-gonic/gin#Context.ShouldBindJSON), see [#1047](https://github.com/gin-gonic/gin/pull/1047)
- [NEW] Add support for `time.Time` location in form binding, see [#1117](https://github.com/gin-gonic/gin/pull/1117)
- [NEW] Add [`func (*Context) BindQuery`](https://godoc.org/github.com/gin-gonic/gin#Context.BindQuery), see [#1029](https://github.com/gin-gonic/gin/pull/1029)
- [NEW] Make [jsonite](https://github.com/json-iterator/go) optional with build tags, see [#1026](https://github.com/gin-gonic/gin/pull/1026)
- [NEW] Show query string in logger, see [#999](https://github.com/gin-gonic/gin/pull/999)
- [NEW] Add [`func (*Context) SecureJSON`](https://godoc.org/github.com/gin-gonic/gin#Context.SecureJSON), see [#987](https://github.com/gin-gonic/gin/pull/987) and [#993](https://github.com/gin-gonic/gin/pull/993)
- [DEPRECATE] `func (*Context) GetCookie` for [`func (*Context) Cookie`](https://godoc.org/github.com/gin-gonic/gin#Context.Cookie)
- [FIX] Don't display color tags if [`func DisableConsoleColor`](https://godoc.org/github.com/gin-gonic/gin#DisableConsoleColor) called, see [#1072](https://github.com/gin-gonic/gin/pull/1072)
- [FIX] Gin Mode `""` when calling [`func Mode`](https://godoc.org/github.com/gin-gonic/gin#Mode) now returns `const DebugMode`, see [#1250](https://github.com/gin-gonic/gin/pull/1250)
- [FIX] `Flush()` now doesn't overwrite `responseWriter` status code, see [#1460](https://github.com/gin-gonic/gin/pull/1460)

## Gin 1.2.0

- [NEW] Switch from godeps to govendor
- [NEW] Add support for Let's Encrypt via gin-gonic/autotls
- [NEW] Improve README examples and add extra at examples folder
- [NEW] Improved support with App Engine
- [NEW] Add custom template delimiters, see #860
- [NEW] Add Template Func Maps, see #962
- [NEW] Add \*context.Handler(), see #928
- [NEW] Add \*context.GetRawData()
- [NEW] Add \*context.GetHeader() (request)
- [NEW] Add \*context.AbortWithStatusJSON() (JSON content type)
- [NEW] Add \*context.Keys type cast helpers
- [NEW] Add \*context.ShouldBindWith()
- [NEW] Add \*context.MustBindWith()
- [NEW] Add \*engine.SetFuncMap()
- [DEPRECATE] On next release: \*context.BindWith(), see #855
- [FIX] Refactor render
- [FIX] Reworked tests
- [FIX] logger now supports cygwin
- [FIX] Use X-Forwarded-For before X-Real-IP
- [FIX] time.Time binding (#904)

## Gin 1.1.4

- [NEW] Support google appengine for IsTerminal func

## Gin 1.1.3

- [FIX] Reverted Logger: skip ANSI color commands

## Gin 1.1

- [NEW] Implement QueryArray and PostArray methods
- [NEW] Refactor GetQuery and GetPostForm
- [NEW] Add contribution guide
- [FIX] Corrected typos in README
- [FIX] Removed additional Iota
- [FIX] Changed imports to gopkg instead of github in README (#733)
- [FIX] Logger: skip ANSI color commands if output is not a tty

## Gin 1.0rc2 (...)

- [PERFORMANCE] Fast path for writing Content-Type.
- [PERFORMANCE] Much faster 404 routing
- [PERFORMANCE] Allocation optimizations
- [PERFORMANCE] Faster root tree lookup
- [PERFORMANCE] Zero overhead, String() and JSON() rendering.
- [PERFORMANCE] Faster ClientIP parsing
- [PERFORMANCE] Much faster SSE implementation
- [NEW] Benchmarks suite
- [NEW] Bind validation can be disabled and replaced with custom validators.
- [NEW] More flexible HTML render
- [NEW] Multipart and PostForm bindings
- [NEW] Adds method to return all the registered routes
- [NEW] Context.HandlerName() returns the main handler's name
- [NEW] Adds Error.IsType() helper
- [FIX] Binding multipart form
- [FIX] Integration tests
- [FIX] Crash when binding non struct object in Context.
- [FIX] RunTLS() implementation
- [FIX] Logger() unit tests
- [FIX] Adds SetHTMLTemplate() warning
- [FIX] Context.IsAborted()
- [FIX] More unit tests
- [FIX] JSON, XML, HTML renders accept custom content-types
- [FIX] gin.AbortIndex is unexported
- [FIX] Better approach to avoid directory listing in StaticFS()
- [FIX] Context.ClientIP() always returns the IP with trimmed spaces.
- [FIX] Better warning when running in debug mode.
- [FIX] Google App Engine integration. debugPrint does not use os.Stdout
- [FIX] Fixes integer overflow in error type
- [FIX] Error implements the json.Marshaller interface
- [FIX] MIT license in every file


## Gin 1.0rc1 (May 22, 2015)

- [PERFORMANCE] Zero allocation router
- [PERFORMANCE] Faster JSON, XML and text rendering
- [PERFORMANCE] Custom hand optimized HttpRouter for Gin
- [PERFORMANCE] Misc code optimizations. Inlining, tail call optimizations
- [NEW] Built-in support for golang.org/x/net/context
- [NEW] Any(path, handler). Create a route that matches any path
- [NEW] Refactored rendering pipeline (faster and static typed)
- [NEW] Refactored errors API
- [NEW] IndentedJSON() prints pretty JSON
- [NEW] Added gin.DefaultWriter
- [NEW] UNIX socket support
- [NEW] RouterGroup.BasePath is exposed
- [NEW] JSON validation using go-validate-yourself (very powerful options)
- [NEW] Completed suite of unit tests
- [NEW] HTTP streaming with c.Stream()
- [NEW] StaticFile() creates a router for serving just one file.
- [NEW] StaticFS() has an option to disable directory listing.
- [NEW] StaticFS() for serving static files through virtual filesystems
- [NEW] Server-Sent Events native support
- [NEW] WrapF() and WrapH() helpers for wrapping http.HandlerFunc and http.Handler
- [NEW] Added LoggerWithWriter() middleware
- [NEW] Added RecoveryWithWriter() middleware
- [NEW] Added DefaultPostFormValue()
- [NEW] Added DefaultFormValue()
- [NEW] Added DefaultParamValue()
- [FIX] BasicAuth() when using custom realm
- [FIX] Bug when serving static files in nested routing group
- [FIX] Redirect using built-in http.Redirect()
- [FIX] Logger when printing the requested path
- [FIX] Documentation typos
- [FIX] Context.Engine renamed to Context.engine
- [FIX] Better debugging messages
- [FIX] ErrorLogger
- [FIX] Debug HTTP render
- [FIX] Refactored binding and render modules
- [FIX] Refactored Context initialization
- [FIX] Refactored BasicAuth()
- [FIX] NoMethod/NoRoute handlers
- [FIX] Hijacking http
- [FIX] Better support for Google App Engine (using log instead of fmt)


## Gin 0.6 (Mar 9, 2015)

- [NEW] Support multipart/form-data
- [NEW] NoMethod handler
- [NEW] Validate sub structures
- [NEW] Support for HTTP Realm Auth
- [FIX] Unsigned integers in binding
- [FIX] Improve color logger


## Gin 0.5 (Feb 7, 2015)

- [NEW] Content Negotiation
- [FIX] Solved security bug that allow a client to spoof ip
- [FIX] Fix unexported/ignored fields in binding


## Gin 0.4 (Aug 21, 2014)

- [NEW] Development mode
- [NEW] Unit tests
- [NEW] Add Content.Redirect()
- [FIX] Deferring WriteHeader()
- [FIX] Improved documentation for model binding


## Gin 0.3 (Jul 18, 2014)

- [PERFORMANCE] Normal log and error log are printed in the same call.
- [PERFORMANCE] Improve performance of NoRouter()
- [PERFORMANCE] Improve context's memory locality, reduce CPU cache faults.
- [NEW] Flexible rendering API
- [NEW] Add Context.File()
- [NEW] Add shortcut RunTLS() for http.ListenAndServeTLS
- [FIX] Rename NotFound404() to NoRoute()
- [FIX] Errors in context are purged
- [FIX] Adds HEAD method in Static file serving
- [FIX] Refactors Static() file serving
- [FIX] Using keyed initialization to fix app-engine integration
- [FIX] Can't unmarshal JSON array, #63
- [FIX] Renaming Context.Req to Context.Request
- [FIX] Check application/x-www-form-urlencoded when parsing form


## Gin 0.2b (Jul 08, 2014)
- [PERFORMANCE] Using sync.Pool to allocatio/gc overhead
- [NEW] Travis CI integration
- [NEW] Completely new logger
- [NEW] New API for serving static files. gin.Static()
- [NEW] gin.H() can be serialized into XML
- [NEW] Typed errors. Errors can be typed. Internet/external/custom.
- [NEW] Support for Godeps
- [NEW] Travis/Godocs badges in README
- [NEW] New Bind() and BindWith() methods for parsing request body.
- [NEW] Add Content.Copy()
- [NEW] Add context.LastError()
- [NEW] Add shortcut for OPTIONS HTTP method
- [FIX] Tons of README fixes
- [FIX] Header is written before body
- [FIX] BasicAuth() and changes API a little bit
- [FIX] Recovery() middleware only prints panics
- [FIX] Context.Get() does not panic anymore. Use MustGet() instead.
- [FIX] Multiple http.WriteHeader() in NotFound handlers
- [FIX] Engine.Run() panics if http server can't be set up
- [FIX] Crash when route path doesn't start with '/'
- [FIX] Do not update header when status code is negative
- [FIX] Setting response headers before calling WriteHeader in context.String()
- [FIX] Add MIT license
- [FIX] Changes behaviour of ErrorLogger() and Logger()
