// Code generated by eventcode -decl component_event.go -core  -emit_package core -export_emit=false; DO NOT EDIT.
package core

func emitEventComponentDestroySelf(event IEvent, comp Component) {
	if event == nil {
		panic("nil event")
	}
	event.Emit(func(delegate IFaceCache) bool {
		Cache2IFace[eventComponentDestroySelf](delegate).onComponentDestroySelf(comp)
		return true
	})
}
