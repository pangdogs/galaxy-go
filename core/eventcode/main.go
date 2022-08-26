// Package eventcode 事件代码生成器，用于`go:generate`生成事件（Event）相关代码。
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	corePackage := flag.String("core", "core", "银河（Galaxy）框架的核心部分包名")
	eventRegexp := flag.String("regexp", "^[eE]vent.+", "匹配事件定义的正则表达式")
	declFile := flag.String("decl", "", "定义事件的源文件（*.go）")

	// 生成事件发送代码选项
	emitGOPackage := flag.String("emit_package", "", "生成事件发送代码的包名")
	emitGenDir := flag.String("gen_emit_dir", "", "生成的事件发送代码（*.go）的目录")
	exportEmit := flag.Bool("export_emit", true, "生成的事件发送代码的可见性")

	// 生成事件表代码选项
	genEventTabCode := flag.String("gen_eventtab_name", "", "生成的事件表名，不填或填空表示不生成事件表代码")
	eventTabGOPackage := flag.String("eventtab_package", "", "生成的事件表代码包名")
	eventTabGenDir := flag.String("gen_eventtab_dir", "", "生成的事件表代码（*.go）的目录")
	eventTabDefEventRecursion := flag.String("eventtab_default_event_recursion", "EventRecursion_Discard", "生成的事件表代码默认的事件递归处理方式")

	flag.Parse()

	if *declFile == "" || filepath.Ext(*declFile) != ".go" {
		flag.Usage()
		panic(flag.ErrHelp)
	}

	if *emitGOPackage == "" {
		flag.Usage()
		panic(flag.ErrHelp)
	}

	declFileData, err := ioutil.ReadFile(*declFile)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()

	fast, err := parser.ParseFile(fset, *declFile, declFileData, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	type EventInfo struct {
		Name    string
		Comment string
	}

	var events []EventInfo

	emitGenAbsDir, _ := filepath.Abs(*emitGenDir)
	eventTabGenAbsDir, _ := filepath.Abs(*eventTabGenDir)

	if *eventTabGOPackage == "" {
		*eventTabGOPackage = *emitGOPackage
	}

	{
		emitGenFile := *emitGenDir

		if emitGenFile == "" {
			emitGenFile = strings.TrimSuffix(*declFile, ".go") + "_emit_code.go"
		} else {
			emitGenFile = filepath.Dir(*declFile) + string(filepath.Separator) + *emitGenDir + string(filepath.Separator) + filepath.Base(strings.TrimSuffix(*declFile, ".go")) + "_emit_code.go"
		}

		genEmitCodeBuff := &bytes.Buffer{}

		fmt.Fprintf(genEmitCodeBuff, `// Code generated by %s%s; DO NOT EDIT.

package %s
`, strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])),
			func() (args string) {
				for _, arg := range os.Args[1:] {
					args += " " + arg
				}
				return
			}(),
			*emitGOPackage)

		emitImportCode := &bytes.Buffer{}

		fmt.Fprintf(emitImportCode, "\nimport (")

		if *corePackage != "" {
			fmt.Fprintf(emitImportCode, `
	%s "github.com/pangdogs/galaxy/core"`, *corePackage)
		}

		for _, imp := range fast.Imports {
			begin := fset.Position(imp.Pos())
			end := fset.Position(imp.End())

			impStr := string(declFileData[begin.Offset:end.Offset])

			if *corePackage != "" && strings.Contains(impStr, "github.com/pangdogs/galaxy/core") {
				continue
			}

			fmt.Fprintf(emitImportCode, "\n\t%s", impStr)
		}

		fmt.Fprintf(emitImportCode, "\n)\n")

		if emitImportCode.Len() > 12 {
			fmt.Fprintf(genEmitCodeBuff, emitImportCode.String())
		}

		exp, err := regexp.Compile(*eventRegexp)
		if err != nil {
			panic(err)
		}

		exportEmitStr := "emit"

		if *exportEmit {
			exportEmitStr = "Emit"
		}

		ast.Inspect(fast, func(node ast.Node) bool {
			ts, ok := node.(*ast.TypeSpec)
			if !ok {
				return true
			}

			eventName := ts.Name.Name
			var eventComment string

			for _, comment := range fast.Comments {
				if fset.Position(comment.End()).Line+1 == fset.Position(node.Pos()).Line {
					eventComment = comment.Text()
					break
				}
			}

			if !exp.MatchString(eventName) {
				return true
			}

			eventIFace, ok := ts.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}

			if eventIFace.Methods.NumFields() <= 0 {
				return true
			}

			eventFuncField := eventIFace.Methods.List[0]

			if len(eventFuncField.Names) <= 0 {
				return true
			}

			eventFuncName := eventFuncField.Names[0].Name

			eventFunc, ok := eventFuncField.Type.(*ast.FuncType)
			if !ok {
				return true
			}

			eventFuncParamsDecl := ""
			eventFuncParams := ""

			if eventFunc.Params != nil {
				for i, param := range eventFunc.Params.List {
					paramName := ""

					for _, pn := range param.Names {
						if paramName != "" {
							paramName += ", "
						}
						paramName += pn.Name
					}

					if paramName == "" {
						paramName = fmt.Sprintf("p%d", i)
					}

					if eventFuncParams != "" {
						eventFuncParams += ", "
					}
					eventFuncParams += paramName

					begin := fset.Position(param.Type.Pos())
					end := fset.Position(param.Type.End())

					eventFuncParamsDecl += fmt.Sprintf(", %s %s", paramName, declFileData[begin.Offset:end.Offset])
				}
			}

			eventFuncTypeParamsDecl := ""
			eventFuncTypeParams := ""

			if ts.TypeParams != nil {
				for i, typeParam := range ts.TypeParams.List {
					typeParamName := ""

					for _, pn := range typeParam.Names {
						if typeParamName != "" {
							typeParamName += ", "
						}
						typeParamName += pn.Name
					}

					if typeParamName == "" {
						typeParamName = fmt.Sprintf("p%d", i)
					}

					if eventFuncTypeParams != "" {
						eventFuncTypeParams += ", "
					}
					eventFuncTypeParams += typeParamName

					begin := fset.Position(typeParam.Type.Pos())
					end := fset.Position(typeParam.Type.End())

					if eventFuncTypeParamsDecl != "" {
						eventFuncTypeParamsDecl += ", "
					}
					eventFuncTypeParamsDecl += fmt.Sprintf("%s %s", typeParamName, declFileData[begin.Offset:end.Offset])
				}
			}

			if eventFuncTypeParamsDecl != "" {
				eventFuncTypeParamsDecl = fmt.Sprintf("[%s]", eventFuncTypeParamsDecl)
			}

			if eventFuncTypeParams != "" {
				eventFuncTypeParams = fmt.Sprintf("[%s]", eventFuncTypeParams)
			}

			_corePackage := ""
			if *corePackage != "" {
				_corePackage = *corePackage + "."
			}

			if eventFunc.Results.NumFields() > 0 {
				eventRet, ok := eventFunc.Results.List[0].Type.(*ast.Ident)
				if !ok {
					return true
				}

				if eventRet.Name != "bool" {
					return true
				}

				fmt.Fprintf(genEmitCodeBuff, `
func %[9]s%[1]s%[7]s(event %[6]sIEvent%[4]s) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate %[6]sIFaceCache) bool {
		return %[6]sCache2IFace[%[2]s%[8]s](delegate).%[3]s(%[5]s)
	})
}
`, strings.Title(eventName), eventName, eventFuncName, eventFuncParamsDecl, eventFuncParams, _corePackage, eventFuncTypeParamsDecl, eventFuncTypeParams, exportEmitStr)

			} else {

				fmt.Fprintf(genEmitCodeBuff, `
func %[9]s%[1]s%[7]s(event %[6]sIEvent%[4]s) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate %[6]sIFaceCache) bool {
		%[6]sCache2IFace[%[2]s%[8]s](delegate).%[3]s(%[5]s)
		return true
	})
}
`, strings.Title(eventName), eventName, eventFuncName, eventFuncParamsDecl, eventFuncParams, _corePackage, eventFuncTypeParamsDecl, eventFuncTypeParams, exportEmitStr)
			}

			fmt.Println(eventName)

			if emitGenAbsDir != eventTabGenAbsDir || *emitGOPackage != *eventTabGOPackage {
				eventName = strings.Title(eventName)
			}

			events = append(events, EventInfo{
				Name:    eventName,
				Comment: eventComment,
			})

			return true
		})

		os.MkdirAll(filepath.Dir(emitGenFile), os.ModePerm)

		if err := ioutil.WriteFile(emitGenFile, genEmitCodeBuff.Bytes(), os.ModePerm); err != nil {
			panic(err)
		}
	}

	if *genEventTabCode != "" {
		eventTabGenFile := *eventTabGenDir

		if eventTabGenFile == "" {
			eventTabGenFile = strings.TrimSuffix(*declFile, ".go") + "_eventtab_code.go"
		} else {
			eventTabGenFile = filepath.Dir(*declFile) + string(filepath.Separator) + *eventTabGenDir + string(filepath.Separator) + filepath.Base(strings.TrimSuffix(*declFile, ".go")) + "_eventTab_code.go"
		}

		genEventTabCodeBuff := &bytes.Buffer{}

		fmt.Fprintf(genEventTabCodeBuff, `// Code generated by %s%s; DO NOT EDIT.

package %s
`, strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0])),
			func() (args string) {
				for _, arg := range os.Args[1:] {
					args += " " + arg
				}
				return
			}(),
			*eventTabGOPackage)

		eventTabImportCode := &bytes.Buffer{}

		fmt.Fprintf(eventTabImportCode, "\nimport (")

		if *corePackage != "" {
			fmt.Fprintf(eventTabImportCode, `
	%s "github.com/pangdogs/galaxy/core"`, *corePackage)
		} else {
			if emitGenAbsDir != eventTabGenAbsDir || *emitGOPackage != *eventTabGOPackage {
				fmt.Fprintf(eventTabImportCode, `
	core "github.com/pangdogs/galaxy/core"`)
			}
		}

		fmt.Fprintf(eventTabImportCode, `
	"github.com/pangdogs/galaxy/core/container"`)

		fmt.Fprintf(eventTabImportCode, "\n)\n")

		if eventTabImportCode.Len() > 12 {
			fmt.Fprintf(genEventTabCodeBuff, eventTabImportCode.String())
		}

		var eventsCode string
		var eventsRecursionCode string

		_corePackage := ""
		if *corePackage != "" {
			_corePackage = *corePackage + "."
		} else {
			if emitGenAbsDir != eventTabGenAbsDir || *emitGOPackage != *eventTabGOPackage {
				_corePackage = "core."
			}
		}

		for _, event := range events {
			eventsCode += fmt.Sprintf("\t%s() %sIEvent\n", event.Name, _corePackage)
		}

		fmt.Fprintf(genEventTabCodeBuff, `
type I%[1]s interface {
%[2]s}
`, *genEventTabCode, eventsCode)

		for i, event := range events {
			eventRecursion := *eventTabDefEventRecursion

			switch eventRecursion {
			case "EventRecursion_Allow", "EventRecursion_Discard", "EventRecursion_Deepest", "EventRecursion_Disallow":
			default:
				eventRecursion = "EventRecursion_Discard"
			}

			if strings.Contains(event.Comment, "[EventRecursion_Allow]") {
				eventRecursion = "EventRecursion_Allow"
			} else if strings.Contains(event.Comment, "[EventRecursion_Discard]") {
				eventRecursion = "EventRecursion_Discard"
			} else if strings.Contains(event.Comment, "[EventRecursion_Deepest]") {
				eventRecursion = "EventRecursion_Deepest"
			} else if strings.Contains(event.Comment, "[EventRecursion_Disallow]") {
				eventRecursion = "EventRecursion_Disallow"
			}

			eventsRecursionCode += fmt.Sprintf("\t(*eventTab)[%d].Init(autoRecover, reportError, %s%s, hookCache, gcCollector)\n", i, _corePackage, eventRecursion)
		}

		var eventsAccessCode string

		for i, event := range events {
			eventsAccessCode += fmt.Sprintf(`
const %[2]sID int = %[4]d

func (eventTab *%[1]s) %[2]s() %[3]sIEvent {
	return &(*eventTab)[%[2]sID]
}
`, *genEventTabCode, event.Name, _corePackage, i)
		}

		fmt.Fprintf(genEventTabCodeBuff, `
type %[1]s [%[2]d]%[4]sEvent

func (eventTab *%[1]s) Init(autoRecover bool, reportError chan error, hookCache *container.Cache[%[4]sHook], gcCollector container.GCCollector) {
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
`, *genEventTabCode, len(events), eventsRecursionCode, _corePackage, eventsAccessCode)

		os.MkdirAll(filepath.Dir(eventTabGenFile), os.ModePerm)

		if err := ioutil.WriteFile(eventTabGenFile, genEventTabCodeBuff.Bytes(), os.ModePerm); err != nil {
			panic(err)
		}
	}
}
