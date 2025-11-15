/*
 * Copyright 2021 ByteDance Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ast

import (
    `github.com/bytedance/sonic/internal/rt`
    `github.com/bytedance/sonic/internal/native/types`
)

// SearchOptions controls Searcher's behavior
type SearchOptions struct {
    // ValidateJSON indicates the searcher to validate the entire JSON
    ValidateJSON bool

    // CopyReturn indicates the searcher to copy the result JSON instead of refer from the input
    // This can help to reduce memory usage if you cache the results
    CopyReturn bool

    // ConcurrentRead indicates the searcher to return a concurrently-READ-safe node,
    // including: GetByPath/Get/Index/GetOrIndex/Int64/Bool/Float64/String/Number/Interface/Array/Map/Raw/MarshalJSON
    ConcurrentRead bool
}

type Searcher struct {
    parser Parser
    SearchOptions
}

func NewSearcher(str string) *Searcher {
    return &Searcher{
        parser: Parser{
            s:      str,
            noLazy: false,
        },
        SearchOptions: SearchOptions{
            ValidateJSON: true,
        },
    }
}

// GetByPathCopy search in depth from top json and returns a **Copied** json node at the path location
func (self *Searcher) GetByPathCopy(path ...interface{}) (Node, error) {
    self.CopyReturn = true
    return self.getByPath(path...)
}

// GetByPathNoCopy search in depth from top json and returns a **Referenced** json node at the path location
//
// WARN: this search directly refer partial json from top json, which has faster speed,
// may consumes more memory.
func (self *Searcher) GetByPath(path ...interface{}) (Node, error) {
    return self.getByPath(path...)
}

func (self *Searcher) getByPath(path ...interface{}) (Node, error) {
    var err types.ParsingError
    var start int

    self.parser.p = 0
    start, err = self.parser.getByPath(self.ValidateJSON, path...)
    if err != 0 {
        // for compatibility with old version
        if err == types.ERR_NOT_FOUND {
            return Node{}, ErrNotExist
        }
        if err == types.ERR_UNSUPPORT_TYPE {
            panic("path must be either int(>=0) or string")
        }
        return Node{}, self.parser.syntaxError(err)
    }

    t := switchRawType(self.parser.s[start])
    if t == _V_NONE {
        return Node{}, self.parser.ExportError(err)
    }

    // copy string to reducing memory usage
    var raw string
    if self.CopyReturn {
        raw = rt.Mem2Str([]byte(self.parser.s[start:self.parser.p]))
    } else {
        raw = self.parser.s[start:self.parser.p]
    }
    return newRawNode(raw, t, self.ConcurrentRead), nil
}

// GetByPath searches a path and returns relaction and types of target
func _GetByPath(src string, path ...interface{}) (start int, end int, typ int, err error) {
	p := NewParserObj(src)
	s, e := p.getByPath(false, path...)
	if e != 0 {
		// for compatibility with old version
		if e == types.ERR_NOT_FOUND {
			return -1, -1, 0, ErrNotExist
		}
		if e == types.ERR_UNSUPPORT_TYPE {
			panic("path must be either int(>=0) or string")
		}
		return -1, -1, 0, p.syntaxError(e)
	}

	t := switchRawType(p.s[s])
	if t == _V_NONE {
		return -1, -1, 0, ErrNotExist
	}
    if t == _V_NUMBER {
        p.p = 1 + backward(p.s, p.p-1)
    }
	return s, p.p, int(t), nil
}

// ValidSyntax check if a json has a valid JSON syntax,
// while not validate UTF-8 charset
func _ValidSyntax(json string) bool {
	p := NewParserObj(json)
    _, e := p.skip()
	if e != 0 {
        return false
    }
   if skipBlank(p.s, p.p) != -int(types.ERR_EOF) {
        return false
   }
   return true
}

// SkipFast skip a json value in fast-skip algs, 
// while not strictly validate JSON syntax and UTF-8 charset.
func _SkipFast(src string, i int) (int, int, error) {
    p := NewParserObj(src)
    p.p = i
    s, e := p.skipFast()
    if e != 0 {
        return -1, -1, p.ExportError(e)
    }
    t := switchRawType(p.s[s])
	if t == _V_NONE {
		return -1, -1, ErrNotExist
	}
    if t == _V_NUMBER {
        p.p = 1 + backward(p.s, p.p-1)
    }
    return s, p.p, nil
}
