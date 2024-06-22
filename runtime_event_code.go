// Code generated by eventc event --default_export=false --default_auto=false; DO NOT EDIT.

package core

import (
	"fmt"
	event "git.golaxy.org/core/event"
	iface "git.golaxy.org/core/utils/iface"
)

func _EmitEventUpdate(evt event.IEvent) {
	if evt == nil {
		panic(fmt.Errorf("%w: %w: evt is nil", event.ErrEvent, event.ErrArgs))
	}
	event.UnsafeEvent(evt).Emit(func(subscriber iface.Cache) bool {
		iface.Cache2Iface[eventUpdate](subscriber).Update()
		return true
	})
}

func _EmitEventUpdateWithInterrupt(evt event.IEvent, interrupt func() bool) {
	if evt == nil {
		panic(fmt.Errorf("%w: %w: evt is nil", event.ErrEvent, event.ErrArgs))
	}
	event.UnsafeEvent(evt).Emit(func(subscriber iface.Cache) bool {
		if interrupt != nil {
			if interrupt() {
				return false
			}
		}
		iface.Cache2Iface[eventUpdate](subscriber).Update()
		return true
	})
}

func _HandleEventUpdate(fun func()) _EventUpdateHandler {
	return _EventUpdateHandler(fun)
}

type _EventUpdateHandler func()

func (h _EventUpdateHandler) Update() {
	h()
}

func _EmitEventLateUpdate(evt event.IEvent) {
	if evt == nil {
		panic(fmt.Errorf("%w: %w: evt is nil", event.ErrEvent, event.ErrArgs))
	}
	event.UnsafeEvent(evt).Emit(func(subscriber iface.Cache) bool {
		iface.Cache2Iface[eventLateUpdate](subscriber).LateUpdate()
		return true
	})
}

func _EmitEventLateUpdateWithInterrupt(evt event.IEvent, interrupt func() bool) {
	if evt == nil {
		panic(fmt.Errorf("%w: %w: evt is nil", event.ErrEvent, event.ErrArgs))
	}
	event.UnsafeEvent(evt).Emit(func(subscriber iface.Cache) bool {
		if interrupt != nil {
			if interrupt() {
				return false
			}
		}
		iface.Cache2Iface[eventLateUpdate](subscriber).LateUpdate()
		return true
	})
}

func _HandleEventLateUpdate(fun func()) _EventLateUpdateHandler {
	return _EventLateUpdateHandler(fun)
}

type _EventLateUpdateHandler func()

func (h _EventLateUpdateHandler) LateUpdate() {
	h()
}
