# 1.11.2 - 2023-09-15

### Fix bugs

- Fix quoted comments ( #370 )
- Fix handle of space at start or last ( #376 )
- Fix sequence with comment ( #390 )

# 1.11.1 - 2023-09-14

### Fix bugs

- Handle `\r` in a double-quoted string the same as `\n` ( #372 )
- Replace loop with n.Values = append(n.Values, target.Values...) ( #380 )
- Skip encoding an inline field if it is null ( #386 )
- Fix comment parsing with null value ( #388 )

# 1.11.0 - 2023-04-03

### Features

- Supports dynamically switch encode and decode processing for a given type

# 1.10.1 - 2023-03-28

### Features

- Quote YAML 1.1 bools at encoding time for compatibility with other legacy parsers
- Add support of 32-bit architecture

### Fix bugs

- Don't trim all space characters in block style sequence
- Support strings starting with `@`

# 1.10.0 - 2023-03-01

### Fix bugs

Reversible conversion of comments was not working in various cases, which has been corrected.
**Breaking Change** exists in the comment map interface. However, if you are dealing with CommentMap directly, there is no problem.


# 1.9.8 - 2022-12-19

### Fix feature

- Append new line at the end of file ( #329 )

### Fix bugs

- Fix custom marshaler ( #333, #334 )
- Fix behavior when struct fields conflicted( #335 )
- Fix position calculation for literal, folded and raw folded strings ( #330 )

# 1.9.7 - 2022-12-03

### Fix bugs

- Fix handling of quoted map key ( #328 )
- Fix resusing process of scanning context ( #322 )

## v1.9.6 - 2022-10-26

### New Features

- Introduce MapKeyNode interface to limit node types for map key ( #312 )

### Fix bugs

- Quote strings with special characters in flow mode ( #270 )
- typeError implements PrettyPrinter interface ( #280 )
- Fix incorrect const type ( #284 )
- Fix large literals type inference on 32 bits ( #293 )
- Fix UTF-8 characters ( #294 )
- Fix decoding of unknown aliases ( #317 )
- Fix stream encoder for insert a separator between each encoded document ( #318 )

### Update

- Update golang.org/x/sys ( #289 )
- Update Go version in CI ( #295 )
- Add test cases for missing keys to struct literals ( #300 )

## v1.9.5 - 2022-01-12

### New Features

* Add UseSingleQuote option ( #265 )

### Fix bugs

* Preserve defaults while decoding nested structs ( #260 )
* Fix minor typo in decodeInit error ( #264 )
* Handle empty sequence entries ( #275 )
* Fix encoding of sequence with multiline string ( #276 )
* Fix encoding of BytesMarshaler type ( #277 )
* Fix indentState logic for multi-line value ( #278 )

## v1.9.4 - 2021-10-12

### Fix bugs

* Keep prev/next reference between tokens containing comments when filtering comment tokens ( #257 )
* Supports escaping reserved keywords in PathBuilder ( #258 )

## v1.9.3 - 2021-09-07

### New Features

* Support encoding and decoding `time.Duration` fields ( #246 )
* Allow reserved characters for key name in YAMLPath ( #251 )
* Support getting YAMLPath from ast.Node ( #252 )
* Support CommentToMap option ( #253 )

### Fix bugs

* Fix encoding nested sequences with `yaml.IndentSequence` ( #241 )
* Fix error reporting on inline structs in strict mode ( #244, #245 )
* Fix encoding of large floats ( #247 )

### Improve workflow

* Migrate CI from CircleCI to GitHub Action ( #249 )
* Add workflow for ycat ( #250 )

## v1.9.2 - 2021-07-26

### Support WithComment option ( #238 )

`yaml.WithComment` is a option for encoding with comment.
The position where you want to add a comment is represented by YAMLPath, and it is the key of `yaml.CommentMap`.
Also, you can select `Head` comment or `Line` comment as the comment type.

## v1.9.1 - 2021-07-20

### Fix DecodeFromNode ( #237 )

- Fix YAML handling where anchor exists

## v1.9.0 - 2021-07-19

### New features

- Support encoding of comment node ( #233 )
- Support `yaml.NodeToValue(ast.Node, interface{}, ...DecodeOption) error` ( #236 )
  - Can convert a AST node to a value directly

### Fix decoder for comment

- Fix parsing of literal with comment ( #234 )

### Rename API ( #235 )

- Rename `MarshalWithContext` to `MarshalContext`
- Rename `UnmarshalWithContext` to `UnmarshalContext`

## v1.8.10 - 2021-07-02

### Fixed bugs

- Fix searching anchor by alias name ( #212 )
- Fixing Issue 186, scanner should account for newline characters when processing multi-line text. Without this source annotations line/column number (for this and all subsequent tokens) is inconsistent with plain text editors. e.g. https://github.com/goccy/go-yaml/issues/186. This addresses the issue specifically for single and double quote text only. ( #210 )
- Add error for unterminated flow mapping node ( #213 )
- Handle missing required field validation ( #221 )
- Nicely format unexpected node type errors ( #229 )
- Support to encode map which has defined type key ( #231 )

### New features

- Support sequence indentation by EncodeOption ( #232 )

## v1.8.9 - 2021-03-01

### Fixed bugs

- Fix origin buffer for DocumentHeader and DocumentEnd and Directive
- Fix origin buffer for anchor value
- Fix syntax error about map value
- Fix parsing MergeKey ('<<') characters
- Fix encoding of float value
- Fix incorrect column annotation when single or double quotes are used

### New features

- Support to encode/decode of ast.Node directly
