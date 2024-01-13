// Code generated by eventcode --decl_file=entity_event.go gen_event --package=ec; DO NOT EDIT.

package ec

import (
	"fmt"
	event "git.golaxy.org/core/event"
	iface "git.golaxy.org/core/util/iface"
)

type iAutoEventEntityDestroySelf interface {
	EventEntityDestroySelf() event.IEvent
}

func BindEventEntityDestroySelf(auto iAutoEventEntityDestroySelf, subscriber EventEntityDestroySelf, priority ...int32) event.Hook {
	if auto == nil {
		panic(fmt.Errorf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs))
	}
	return event.BindEvent[EventEntityDestroySelf](auto.EventEntityDestroySelf(), subscriber, priority...)
}

func emitEventEntityDestroySelf(auto iAutoEventEntityDestroySelf, entity Entity) {
	if auto == nil {
		panic(fmt.Errorf("%w: %w: auto is nil", event.ErrEvent, event.ErrArgs))
	}
	event.UnsafeEvent(auto.EventEntityDestroySelf()).Emit(func(subscriber iface.Cache) bool {
		iface.Cache2Iface[EventEntityDestroySelf](subscriber).OnEntityDestroySelf(entity)
		return true
	})
}
