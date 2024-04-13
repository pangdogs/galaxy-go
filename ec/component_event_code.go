// Code generated by eventcode --decl_file=component_event.go gen_event --package=ec; DO NOT EDIT.

package ec

import (
	"fmt"
	event "git.golaxy.org/core/event"
	iface "git.golaxy.org/core/util/iface"
)

type iAutoEventComponentDestroySelf interface {
	EventComponentDestroySelf() event.IEvent
}

func BindEventComponentDestroySelf(auto iAutoEventComponentDestroySelf, subscriber EventComponentDestroySelf, priority ...int32) event.Hook {
	if auto == nil {
		panic(fmt.Errorf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs))
	}
	return event.BindEvent[EventComponentDestroySelf](auto.EventComponentDestroySelf(), subscriber, priority...)
}

func emitEventComponentDestroySelf(auto iAutoEventComponentDestroySelf, comp Component) {
	if auto == nil {
		panic(fmt.Errorf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs))
	}
	event.UnsafeEvent(auto.EventComponentDestroySelf()).Emit(func(subscriber iface.Cache) bool {
		iface.Cache2Iface[EventComponentDestroySelf](subscriber).OnComponentDestroySelf(comp)
		return true
	})
}

func HandleEventComponentDestroySelf(fun func(comp Component)) handleEventComponentDestroySelf {
	return handleEventComponentDestroySelf(fun)
}

type handleEventComponentDestroySelf func(comp Component)

func (handle handleEventComponentDestroySelf) OnComponentDestroySelf(comp Component) {
	handle(comp)
}
