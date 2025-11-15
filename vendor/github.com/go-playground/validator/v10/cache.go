package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

type tagType uint8

const (
	typeDefault tagType = iota
	typeOmitEmpty
	typeIsDefault
	typeNoStructLevel
	typeStructOnly
	typeDive
	typeOr
	typeKeys
	typeEndKeys
	typeOmitNil
	typeOmitZero
)

const (
	invalidValidation   = "Invalid validation tag on field '%s'"
	undefinedValidation = "Undefined validation function '%s' on field '%s'"
	keysTagNotDefined   = "'" + endKeysTag + "' tag encountered without a corresponding '" + keysTag + "' tag"
)

type structCache struct {
	lock sync.Mutex
	m    atomic.Value // map[reflect.Type]*cStruct
}

func (sc *structCache) Get(key reflect.Type) (c *cStruct, found bool) {
	c, found = sc.m.Load().(map[reflect.Type]*cStruct)[key]
	return
}

func (sc *structCache) Set(key reflect.Type, value *cStruct) {
	m := sc.m.Load().(map[reflect.Type]*cStruct)
	nm := make(map[reflect.Type]*cStruct, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}
	nm[key] = value
	sc.m.Store(nm)
}

type tagCache struct {
	lock sync.Mutex
	m    atomic.Value // map[string]*cTag
}

func (tc *tagCache) Get(key string) (c *cTag, found bool) {
	c, found = tc.m.Load().(map[string]*cTag)[key]
	return
}

func (tc *tagCache) Set(key string, value *cTag) {
	m := tc.m.Load().(map[string]*cTag)
	nm := make(map[string]*cTag, len(m)+1)
	for k, v := range m {
		nm[k] = v
	}
	nm[key] = value
	tc.m.Store(nm)
}

type cStruct struct {
	name   string
	fields []*cField
	fn     StructLevelFuncCtx
}

type cField struct {
	idx        int
	name       string
	altName    string
	namesEqual bool
	cTags      *cTag
}

type cTag struct {
	tag                  string
	aliasTag             string
	actualAliasTag       string
	param                string
	keys                 *cTag // only populated when using tag's 'keys' and 'endkeys' for map key validation
	next                 *cTag
	fn                   FuncCtx
	typeof               tagType
	hasTag               bool
	hasAlias             bool
	hasParam             bool // true if parameter used eg. eq= where the equal sign has been set
	isBlockEnd           bool // indicates the current tag represents the last validation in the block
	runValidationWhenNil bool
}

func (v *Validate) extractStructCache(current reflect.Value, sName string) *cStruct {
	v.structCache.lock.Lock()
	defer v.structCache.lock.Unlock() // leave as defer! because if inner panics, it will never get unlocked otherwise!

	typ := current.Type()

	// could have been multiple trying to access, but once first is done this ensures struct
	// isn't parsed again.
	cs, ok := v.structCache.Get(typ)
	if ok {
		return cs
	}

	cs = &cStruct{name: sName, fields: make([]*cField, 0), fn: v.structLevelFuncs[typ]}

	numFields := current.NumField()
	rules := v.rules[typ]

	var ctag *cTag
	var fld reflect.StructField
	var tag string
	var customName string

	for i := 0; i < numFields; i++ {
		fld = typ.Field(i)

		if !v.privateFieldValidation && !fld.Anonymous && len(fld.PkgPath) > 0 {
			continue
		}

		if rtag, ok := rules[fld.Name]; ok {
			tag = rtag
		} else {
			tag = fld.Tag.Get(v.tagName)
		}

		if tag == skipValidationTag {
			continue
		}

		customName = fld.Name

		if v.hasTagNameFunc {
			name := v.tagNameFunc(fld)
			if len(name) > 0 {
				customName = name
			}
		}

		// NOTE: cannot use shared tag cache, because tags may be equal, but things like alias may be different
		// and so only struct level caching can be used instead of combined with Field tag caching

		if len(tag) > 0 {
			ctag, _ = v.parseFieldTagsRecursive(tag, fld.Name, "", false)
		} else {
			// even if field doesn't have validations need cTag for traversing to potential inner/nested
			// elements of the field.
			ctag = new(cTag)
		}

		cs.fields = append(cs.fields, &cField{
			idx:        i,
			name:       fld.Name,
			altName:    customName,
			cTags:      ctag,
			namesEqual: fld.Name == customName,
		})
	}
	v.structCache.Set(typ, cs)
	return cs
}

func (v *Validate) parseFieldTagsRecursive(tag string, fieldName string, alias string, hasAlias bool) (firstCtag *cTag, current *cTag) {
	var t string
	noAlias := len(alias) == 0
	tags := strings.Split(tag, tagSeparator)

	for i := 0; i < len(tags); i++ {
		t = tags[i]
		if noAlias {
			alias = t
		}

		// check map for alias and process new tags, otherwise process as usual
		if tagsVal, found := v.aliases[t]; found {
			if i == 0 {
				firstCtag, current = v.parseFieldTagsRecursive(tagsVal, fieldName, t, true)
			} else {
				next, curr := v.parseFieldTagsRecursive(tagsVal, fieldName, t, true)
				current.next, current = next, curr
			}
			continue
		}

		var prevTag tagType

		if i == 0 {
			current = &cTag{aliasTag: alias, hasAlias: hasAlias, hasTag: true, typeof: typeDefault}
			firstCtag = current
		} else {
			prevTag = current.typeof
			current.next = &cTag{aliasTag: alias, hasAlias: hasAlias, hasTag: true}
			current = current.next
		}

		switch t {
		case diveTag:
			current.typeof = typeDive

		case keysTag:
			current.typeof = typeKeys

			if i == 0 || prevTag != typeDive {
				panic(fmt.Sprintf("'%s' tag must be immediately preceded by the '%s' tag", keysTag, diveTag))
			}

			// need to pass along only keys tag
			// need to increment i to skip over the keys tags
			b := make([]byte, 0, 64)

			i++

			for ; i < len(tags); i++ {
				b = append(b, tags[i]...)
				b = append(b, ',')

				if tags[i] == endKeysTag {
					break
				}
			}

			current.keys, _ = v.parseFieldTagsRecursive(string(b[:len(b)-1]), fieldName, "", false)

		case endKeysTag:
			current.typeof = typeEndKeys

			// if there are more in tags then there was no keysTag defined
			// and an error should be thrown
			if i != len(tags)-1 {
				panic(keysTagNotDefined)
			}
			return

		case omitzero:
			current.typeof = typeOmitZero
			continue

		case omitempty:
			current.typeof = typeOmitEmpty

		case omitnil:
			current.typeof = typeOmitNil

		case structOnlyTag:
			current.typeof = typeStructOnly

		case noStructLevelTag:
			current.typeof = typeNoStructLevel

		default:
			if t == isdefault {
				current.typeof = typeIsDefault
			}
			// if a pipe character is needed within the param you must use the utf8Pipe representation "0x7C"
			orVals := strings.Split(t, orSeparator)

			for j := 0; j < len(orVals); j++ {
				vals := strings.SplitN(orVals[j], tagKeySeparator, 2)
				if noAlias {
					alias = vals[0]
					current.aliasTag = alias
				} else {
					current.actualAliasTag = t
				}

				if j > 0 {
					current.next = &cTag{aliasTag: alias, actualAliasTag: current.actualAliasTag, hasAlias: hasAlias, hasTag: true}
					current = current.next
				}
				current.hasParam = len(vals) > 1

				current.tag = vals[0]
				if len(current.tag) == 0 {
					panic(strings.TrimSpace(fmt.Sprintf(invalidValidation, fieldName)))
				}

				if wrapper, ok := v.validations[current.tag]; ok {
					current.fn = wrapper.fn
					current.runValidationWhenNil = wrapper.runValidationOnNil
				} else {
					panic(strings.TrimSpace(fmt.Sprintf(undefinedValidation, current.tag, fieldName)))
				}

				if len(orVals) > 1 {
					current.typeof = typeOr
				}

				if len(vals) > 1 {
					current.param = strings.ReplaceAll(strings.ReplaceAll(vals[1], utf8HexComma, ","), utf8Pipe, "|")
				}
			}
			current.isBlockEnd = true
		}
	}
	return
}

func (v *Validate) fetchCacheTag(tag string) *cTag {
	// find cached tag
	ctag, found := v.tagCache.Get(tag)
	if !found {
		v.tagCache.lock.Lock()
		defer v.tagCache.lock.Unlock()

		// could have been multiple trying to access, but once first is done this ensures tag
		// isn't parsed again.
		ctag, found = v.tagCache.Get(tag)
		if !found {
			ctag, _ = v.parseFieldTagsRecursive(tag, "", "", false)
			v.tagCache.Set(tag, ctag)
		}
	}
	return ctag
}
