package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func genEventTab(ctx *CommandContext) {
	eventTabFile := ctx.EventTabDir

	if eventTabFile == "" {
		eventTabFile = strings.TrimSuffix(ctx.DeclFile, ".go") + "_tab_code.go"
	} else {
		eventTabFile = strings.Join([]string{filepath.Dir(ctx.DeclFile), ctx.EventTabDir, filepath.Base(strings.TrimSuffix(ctx.DeclFile, ".go")) + "_tab_code.go"}, string(filepath.Separator))
	}

	code := &bytes.Buffer{}

	// 生成注释
	{
		program := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
		args := strings.Join(os.Args[1:], " ")

		fmt.Fprintf(code, `// Code generated by %s %s; DO NOT EDIT.

package %s
`, program, args, ctx.EventTabPackage)
	}

	// 生成import
	{
		importCode := &bytes.Buffer{}

		fmt.Fprintf(importCode, "\nimport (")

		fmt.Fprintf(importCode, `
	%s "%s"`, ctx.PackageEventAlias, packageEventPath)

		fmt.Fprintf(importCode, `
	container "git.golaxy.org/core/util/container"`)

		fmt.Fprintf(importCode, "\n)\n")

		fmt.Fprintf(code, importCode.String())
	}

	// 解析事件定义
	eventDeclTab := EventDeclTab{}
	eventDeclTab.Parse(ctx)

	// event包前缀
	eventPrefix := ""
	if ctx.PackageEventAlias != "." {
		eventPrefix = ctx.PackageEventAlias + "."
	}

	// 生成事件表接口
	{
		var eventsCode string

		for _, event := range eventDeclTab {
			eventsCode += fmt.Sprintf("\t%s() %sIEvent\n", event.Name, eventPrefix)
		}

		fmt.Fprintf(code, `
type I%[1]s interface {
%[2]s}
`, strings.Title(ctx.EventTabName), eventsCode)
	}

	// 生成事件表
	{
		var eventsRecursionCode string

		for i, event := range eventDeclTab {
			var eventRecursion string

			if strings.Contains(event.Comment, "[EventRecursion_Allow]") {
				eventRecursion = eventPrefix + "EventRecursion_Allow"
			} else if strings.Contains(event.Comment, "[EventRecursion_Disallow]") {
				eventRecursion = eventPrefix + "EventRecursion_Disallow"
			} else if strings.Contains(event.Comment, "[EventRecursion_Discard]") {
				eventRecursion = eventPrefix + "EventRecursion_Discard"
			} else if strings.Contains(event.Comment, "[EventRecursion_Truncate]") {
				eventRecursion = eventPrefix + "EventRecursion_Truncate"
			} else if strings.Contains(event.Comment, "[EventRecursion_Deepest]") {
				eventRecursion = eventPrefix + "EventRecursion_Deepest"
			} else {
				eventRecursion = "recursion"
			}

			eventsRecursionCode += fmt.Sprintf("\t(*eventTab)[%d].Init(autoRecover, reportError, %s, hookAllocator, gcCollector)\n", i, eventRecursion)
		}

		var eventsAccessCode string

		for i, event := range eventDeclTab {
			eventsAccessCode += fmt.Sprintf(`
const %[2]sId int = %[4]d

func (eventTab *%[1]s) %[2]s() %[3]sIEvent {
	return &(*eventTab)[%[2]sId]
}
`, ctx.EventTabName, event.Name, eventPrefix, i)
		}

		fmt.Fprintf(code, `
type %[1]s [%[2]d]%[4]sEvent

func (eventTab *%[1]s) Init(autoRecover bool, reportError chan error, recursion %[4]sEventRecursion, hookAllocator container.Allocator[%[4]sHook], gcCollector container.GCCollector) {
%[3]s}

func (eventTab *%[1]s) Get(id int) %[4]sIEvent {
	return &(*eventTab)[id]
}

func (eventTab *%[1]s) Open() {
	for i := range *eventTab {
		(*eventTab)[i].Open()
	}
}

func (eventTab *%[1]s) Close() {
	for i := range *eventTab {
		(*eventTab)[i].Close()
	}
}

func (eventTab *%[1]s) Clean() {
	for i := range *eventTab {
		(*eventTab)[i].Clean()
	}
}
%[5]s
`, ctx.EventTabName, len(eventDeclTab), eventsRecursionCode, eventPrefix, eventsAccessCode)
	}

	fmt.Printf("EventTab: %s\n", ctx.EventTabName)

	os.MkdirAll(filepath.Dir(eventTabFile), os.ModePerm)

	if err := ioutil.WriteFile(eventTabFile, code.Bytes(), os.ModePerm); err != nil {
		panic(err)
	}
}
