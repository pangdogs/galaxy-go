/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func genEvent() {
	declFile := viper.GetString("decl_file")
	packageEventAlias := viper.GetString("package_event_alias")
	defExport := viper.GetBool("default_export")
	defAuto := viper.GetBool("default_auto")
	fast := viper.Get("file_ast").(*ast.File)
	fset := viper.Get("file_set").(*token.FileSet)

	// 解析事件定义
	eventDeclTab := EventDeclTab{}
	eventDeclTab.Parse()

	code := &bytes.Buffer{}

	// 生成注释
	{
		program := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
		args := strings.Join(os.Args[1:], " ")
		copyright := copyrightNotice

		if !viper.GetBool("copyright") {
			copyright = ""
		}

		fmt.Fprintf(code, `%s// Code generated by %s %s; DO NOT EDIT.

package %s
`, copyright, program, args, eventDeclTab.Package)
	}

	// 生成import
	{
		importCode := &bytes.Buffer{}

		fmt.Fprintf(importCode, "\nimport (")

		fmt.Fprintf(importCode, `
	%s "%s"`, packageEventAlias, packageEventPath)

		for _, is := range fast.Imports {
			if is.Name != nil {
				if is.Name.Name == packageEventAlias {
					continue
				}
			} else if is.Path != nil {
				if path.Base(is.Path.Value) == packageEventAlias {
					continue
				}
			}

			var buf bytes.Buffer
			printer.Fprint(&buf, fset, is)

			fmt.Fprintf(importCode, "\n\t%s", buf.String())
		}

		fmt.Fprintf(importCode, "\n)\n")

		fmt.Fprintf(code, importCode.String())
	}

	// event包前缀
	eventPrefix := ""
	if packageEventAlias != "." {
		eventPrefix = packageEventAlias + "."
	}

	// 生成事件发送代码
	for _, eventDecl := range eventDeclTab.Events {
		// 是否导出事件发送代码
		exportEmitStr := "_Emit"
		if defExport {
			exportEmitStr = "Emit"
		}

		// 是否生成事件auto代码
		auto := defAuto

		// 解析atti
		atti := parseGenAtti(eventDecl.Comment, "+event-gen:")

		if atti.Has("export") {
			if b, err := strconv.ParseBool(atti.Get("export")); err == nil {
				if b {
					exportEmitStr = "Emit"
				} else {
					exportEmitStr = "_Emit"
				}
			}
		}

		if atti.Has("auto") {
			if b, err := strconv.ParseBool(atti.Get("auto")); err == nil {
				auto = b
			}
		}

		// 事件可见性
		var visibility string

		if unicode.IsLower(rune(eventDecl.Name[0])) {
			visibility = "_"
		}

		// 生成代码
		if auto {
			if eventDecl.FuncHasRet {
				fmt.Fprintf(code, `
type iAuto%[1]s interface {
	%[1]s() %[6]sIEvent
}

func Bind%[1]s(auto iAuto%[1]s, subscriber %[2]s%[8]s, priority ...int32) %[6]sHook {
	if auto == nil {
		%[6]sPanicf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	return %[6]sBind[%[2]s%[8]s](auto.%[1]s(), subscriber, priority...)
}

func %[9]s%[1]s%[7]s(auto iAuto%[1]s%[4]s) {
	if auto == nil {
		%[6]sPanicf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(subscriber %[6]sCache) bool {
		return %[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
	})
}

func %[9]s%[1]s%[7]sWithInterrupt(auto iAuto%[1]s, interrupt func(%[10]s) bool%[4]s) {
	if auto == nil {
		%[6]sPanicf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(subscriber %[6]sCache) bool {
		if interrupt != nil {
			if interrupt(%[5]s) {
				return false
			}
		}
		return %[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, strings.TrimLeft(eventDecl.FuncParamsDecl, ", "))

			} else {
				fmt.Fprintf(code, `
type iAuto%[1]s interface {
	%[1]s() %[6]sIEvent
}

func Bind%[1]s(auto iAuto%[1]s, subscriber %[2]s%[8]s, priority ...int32) %[6]sHook {
	if auto == nil {
		%[6]sPanicf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	return %[6]sBind[%[2]s%[8]s](auto.%[1]s(), subscriber, priority...)
}

func %[9]s%[1]s%[7]s(auto iAuto%[1]s%[4]s) {
	if auto == nil {
		%[6]sPanicf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(subscriber %[6]sCache) bool {
		%[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
		return true
	})
}

func %[9]s%[1]s%[7]sWithInterrupt(auto iAuto%[1]s, interrupt func(%[10]s) bool%[4]s) {
	if auto == nil {
		%[6]sPanicf("%%w: %%w: auto is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(auto.%[1]s()).Emit(func(subscriber %[6]sCache) bool {
		if interrupt != nil {
			if interrupt(%[5]s) {
				return false
			}
		}
		%[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, strings.TrimLeft(eventDecl.FuncParamsDecl, ", "))
			}
		} else {
			if eventDecl.FuncHasRet {
				fmt.Fprintf(code, `
func %[9]s%[1]s%[7]s(evt %[6]sIEvent%[4]s) {
	if evt == nil {
		%[6]sPanicf("%%w: %%w: evt is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(evt).Emit(func(subscriber %[6]sCache) bool {
		return %[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
	})
}

func %[9]s%[1]s%[7]sWithInterrupt(evt %[6]sIEvent, interrupt func(%[10]s) bool%[4]s) {
	if evt == nil {
		%[6]sPanicf("%%w: %%w: evt is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(evt).Emit(func(subscriber %[6]sCache) bool {
		if interrupt != nil {
			if interrupt(%[5]s) {
				return false
			}
		}
		return %[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, strings.TrimLeft(eventDecl.FuncParamsDecl, ", "))

			} else {
				fmt.Fprintf(code, `
func %[9]s%[1]s%[7]s(evt %[6]sIEvent%[4]s) {
	if evt == nil {
		%[6]sPanicf("%%w: %%w: evt is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(evt).Emit(func(subscriber %[6]sCache) bool {
		%[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
		return true
	})
}

func %[9]s%[1]s%[7]sWithInterrupt(evt %[6]sIEvent, interrupt func(%[10]s) bool%[4]s) {
	if evt == nil {
		%[6]sPanicf("%%w: %%w: evt is nil", %[6]sErrEvent, %[6]sErrArgs)
	}
	%[6]sUnsafeEvent(evt).Emit(func(subscriber %[6]sCache) bool {
		if interrupt != nil {
			if interrupt(%[5]s) {
				return false
			}
		}
		%[6]sCache2Iface[%[2]s%[8]s](subscriber).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventDecl.Name), eventDecl.Name, eventDecl.FuncName, eventDecl.FuncParamsDecl, eventDecl.FuncParams, eventPrefix, eventDecl.FuncTypeParamsDecl, eventDecl.FuncTypeParams, exportEmitStr, strings.TrimLeft(eventDecl.FuncParamsDecl, ", "))
			}
		}

		if eventDecl.FuncHasRet {
			fmt.Fprintf(code, `
func %[5]sHandle%[1]s(fun func(%[3]s) bool) %[5]s%[1]sHandler {
	return %[1]sHandler(fun)
}

type %[5]s%[1]sHandler func(%[3]s) bool

func (h %[5]s%[1]sHandler) %[2]s(%[3]s) bool {
	return h(%[4]s)
}
`, strings.Title(eventDecl.Name), eventDecl.FuncName, strings.TrimLeft(eventDecl.FuncParamsDecl, ", "), eventDecl.FuncParams, visibility)
		} else {
			fmt.Fprintf(code, `
func %[5]sHandle%[1]s(fun func(%[3]s)) %[5]s%[1]sHandler {
	return %[5]s%[1]sHandler(fun)
}

type %[5]s%[1]sHandler func(%[3]s)

func (h %[5]s%[1]sHandler) %[2]s(%[3]s) {
	h(%[4]s)
}
`, strings.Title(eventDecl.Name), eventDecl.FuncName, strings.TrimLeft(eventDecl.FuncParamsDecl, ", "), eventDecl.FuncParams, visibility)
		}

		log.Printf("Event: %s", eventDecl.Name)
	}

	// 输出文件
	outFile := filepath.Join(filepath.Dir(declFile), filepath.Base(strings.TrimSuffix(declFile, ".go"))+".gen.go")

	os.MkdirAll(filepath.Dir(outFile), os.ModePerm)

	if err := ioutil.WriteFile(outFile, code.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}
