// Code generated by eventcode --decl_file=component_event.go gen_emit --package=ec; DO NOT EDIT.

package ec

import (
	localevent "github.com/pangdogs/galaxy/localevent"
	"github.com/pangdogs/galaxy/util"
)

func emitEventComponentDestroySelf(event localevent.IEvent, comp Component) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[EventComponentDestroySelf](delegate).OnComponentDestroySelf(comp)
		return true
	})
}
