// +build !amd64,!arm64 go1.26 !go1.17 arm64,!go1.20

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

package utf8

import (
	"unicode/utf8"

	"github.com/bytedance/sonic/internal/rt"
)

// ValidateFallback validates UTF-8 encoded bytes using standard library.
// This is used when native UTF-8 validation is not available.
func Validate(src []byte) bool {
	return utf8.Valid(src)
}

// ValidateStringFallback validates UTF-8 encoded string using standard library.
// This is used when native UTF-8 validation is not available.
func ValidateString(src string) bool {
	return utf8.Valid(rt.Str2Mem(src))
}
