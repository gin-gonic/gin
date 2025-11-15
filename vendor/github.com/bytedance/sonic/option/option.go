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

package option

var (
    // DefaultDecoderBufferSize is the initial buffer size of StreamDecoder
    DefaultDecoderBufferSize  uint = 4 * 1024

    // DefaultEncoderBufferSize is the initial buffer size of Encoder
    DefaultEncoderBufferSize  uint = 4 * 1024

    // DefaultAstBufferSize is the initial buffer size of ast.Node.MarshalJSON()
    DefaultAstBufferSize  uint = 4 * 1024

    // LimitBufferSize indicates the max pool buffer size, in case of OOM.
    // See issue https://github.com/bytedance/sonic/issues/614
    LimitBufferSize uint = 1024 * 1024
)

// CompileOptions includes all options for encoder or decoder compiler.
type CompileOptions struct {
    // the maximum depth for compilation inline
    MaxInlineDepth int

    // the loop times for recursively pretouch
    RecursiveDepth int
}

var (
    // Default value(3) means the compiler only inline 3 layers of nested struct. 
    // when the depth exceeds, the compiler will recurse 
    // and compile subsequent structs when they are decoded 
    DefaultMaxInlineDepth = 3

    // Default value(1) means `Pretouch()` will be recursively executed once,
    // if any nested struct is left (depth exceeds MaxInlineDepth)
    DefaultRecursiveDepth = 1
)

// DefaultCompileOptions set default compile options.
func DefaultCompileOptions() CompileOptions {
    return CompileOptions{
        RecursiveDepth: DefaultRecursiveDepth,
        MaxInlineDepth: DefaultMaxInlineDepth,
    }
}

// CompileOption is a function used to change DefaultCompileOptions.
type CompileOption func(o *CompileOptions)

// WithCompileRecursiveDepth sets the loop times of recursive pretouch 
// in both decoder and encoder,
// for both concrete type and its pointer type.
//
// For deep nested struct (depth exceeds MaxInlineDepth), 
// try to set more loops to completely compile, 
// thus reduce JIT instability in the first hit.
func WithCompileRecursiveDepth(loop int) CompileOption {
    return func(o *CompileOptions) {
            if loop < 0 {
                panic("loop must be >= 0")
            }
            o.RecursiveDepth = loop
        }
}

// WithCompileMaxInlineDepth sets the max depth of inline compile 
// in decoder and encoder.
//
// For large nested struct, try to set smaller depth to reduce compiling time.
func WithCompileMaxInlineDepth(depth int) CompileOption {
    return func(o *CompileOptions) {
            if depth <= 0 {
                panic("depth must be > 0")
            }
            o.MaxInlineDepth = depth
        }
}
