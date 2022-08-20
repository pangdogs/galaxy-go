// Code generated by eventcode -decl ectree_event.go -core  -emit_package core -export_emit=false; DO NOT EDIT.
package core

func emitEventECTreeAddChild(event IEvent, ecTree IECTree, parent, child Entity) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate IFaceCache) bool {
		Cache2IFace[EventECTreeAddChild](delegate).OnAddChild(ecTree, parent, child)
		return true
	})
}

func emitEventECTreeRemoveChild(event IEvent, ecTree IECTree, parent, child Entity) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate IFaceCache) bool {
		Cache2IFace[EventECTreeRemoveChild](delegate).OnRemoveChild(ecTree, parent, child)
		return true
	})
}
