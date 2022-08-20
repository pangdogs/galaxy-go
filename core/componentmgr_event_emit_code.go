// Code generated by eventcode -decl componentmgr_event.go -core  -emit_package core -export_emit=false; DO NOT EDIT.
package core

func emitEventCompMgrAddComponents[T any](event IEvent, compMgr T, components []Component) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate IFaceCache) bool {
		Cache2IFace[EventCompMgrAddComponents[T]](delegate).OnCompMgrAddComponents(compMgr, components)
		return true
	})
}

func emitEventCompMgrRemoveComponent[T any](event IEvent, compMgr T, component Component) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate IFaceCache) bool {
		Cache2IFace[EventCompMgrRemoveComponent[T]](delegate).OnCompMgrRemoveComponent(compMgr, component)
		return true
	})
}
