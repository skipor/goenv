// Copyright 2017 Vladimir Skipor
// Copyright 2014 Chris Hines
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package goenv based on https://github.com/go-stack/stack unexported
// features. Implements utilities to get go compile time $GOROOT, $GOPATH,
// and some other useful runtime information.
package goenv

import (
	"runtime"
	"strings"
)

// PkgIndex returns the index that results in file[index:] being the path of
// file relative to the compile time GOPATH, and file[:index] being the
// $GOPATH/src/ portion of file. funcName must be the name of a function in
// file as returned by runtime.Func.Name.
func PkgIndex(file, funcName string) int {
	// As of Go 1.6.2 there is no direct way to know the compile time GOPATH
	// at runtime, but we can infer the number of path segments in the GOPATH.
	// We note that runtime.Func.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/sub.Type.Method
	//
	// We want to produce:
	//
	//    file[:idx] == /home/user/src/
	//    file[idx:] == pkg/sub/file.go
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired result for file[idx:]. We count separators from the
	// end of the file path until it finds two more than in the function name
	// and then move one character forward to preserve the initial path
	// segment without a leading separator.
	const sep = "/"
	i := len(file)
	for n := strings.Count(funcName, sep) + 2; n > 0; n-- {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			i = -len(sep)
			break
		}
	}
	// get back to 0 or trim the leading separator
	return i + len(sep)
}

func Goroot() string {
	return goroot
}

func Gopath() string {
	return gopath
}

func Gosrc() string {
	return gosrc
}

// TrimGosrc trims compile time $GOPATH/src/
func TrimGosrc(path string) string {
	return strings.TrimPrefix(path, Gopath())
}

func TrimGoroot(path string) string {
	return strings.TrimPrefix(path, Goroot())
}

// InGoroot returns true if file unknown, under GOROOT, or _testmain.go.
func InGoroot(file string) bool {
	if len(file) == 0 || file[0] == '?' {
		return true
	}
	if runtime.GOOS == "windows" {
		file = strings.ToLower(file)
	}
	return strings.HasPrefix(file, goroot) || strings.HasSuffix(file, "/_testmain.go")
}

// Sigpanic returns runtime.sigpanic *runtime.Func.
func Sigpanic() *runtime.Func {
	return sigpanic
}

// Compile time variables.
var (
	goroot string // $GOROOT
	gopath string // $GOPATH
	gosrc  string // $GOPATH/src/
)

func init() {
	// TODO: set gopath, gosrc
	var pcs [1]uintptr
	runtime.Callers(0, pcs[:])
	fn := runtime.FuncForPC(pcs[0])
	file, _ := fn.FileLine(pcs[0])

	idx := PkgIndex(file, fn.Name())

	goroot = file[:idx]
	if runtime.GOOS == "windows" {
		goroot = strings.ToLower(goroot)
	}
}

// findSigpanic intentionally executes faulting code to generate a stack trace
// containing an entry for runtime.sigpanic.
func findSigpanic() *runtime.Func {
	var fn *runtime.Func
	var p *int
	func() int {
		defer func() {
			if p := recover(); p != nil {
				var pcs [512]uintptr
				n := runtime.Callers(2, pcs[:])
				for _, pc := range pcs[:n] {
					f := runtime.FuncForPC(pc)
					if f.Name() == "runtime.sigpanic" {
						fn = f
						break
					}
				}
			}
		}()
		// intentional nil pointer dereference to trigger sigpanic
		return *p
	}()
	return fn
}

var sigpanic = findSigpanic()
