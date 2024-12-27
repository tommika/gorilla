// Copyright (c) 2024 Thomas Mikalsen. Subject to the MIT License
package util

import (
	"runtime"
	"strings"
)

func CallerFuncName(skip int) (packageName, funcName string) {
	f := CallerFunc(1 + skip)
	if f != nil {
		packageName, funcName = ParseFuncName(f.Name())
	}
	return
}

func CallerFunc(skip int) (fn *runtime.Func) {
	if pc, _, _, ok := runtime.Caller(1 + skip); ok {
		fn = runtime.FuncForPC(pc)
	}
	return
}

func ParseFuncName(fqFuncName string) (packageName, funcName string) {
	// fully qualified func names look like so:
	//   example.com/my/pack.MyFunc
	//   example.com/my/pack.init.0
	//   main/MyFunc
	//   example.com/MyFunc
	//   pack.MyFunc
	split := strings.LastIndexByte(fqFuncName, '/')
	if split < 0 {
		split = strings.IndexByte(fqFuncName, '.')
	} else {
		split += strings.IndexByte(fqFuncName[split+1:], '.') + 1
	}
	packageName = fqFuncName[:split]
	funcName = fqFuncName[split+1:]
	return
}
