//go:generate go run kit.golaxy.org/golaxy/event/eventcode --decl_file=$GOFILE gen_emit --package=$GOPACKAGE --default_auto=true

package ec

// EventComponentDestroySelf [EmitUnExport] 事件：组件销毁自身
type EventComponentDestroySelf interface {
	OnComponentDestroySelf(comp Component)
}
