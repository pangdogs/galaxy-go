// Code generated by eventcode --decl_file=entitymgr_event.go gen_emit --package=runtime; DO NOT EDIT.

package runtime

import (
	localevent "github.com/pangdogs/galaxy/localevent"
	"github.com/pangdogs/galaxy/ec"
	"github.com/pangdogs/galaxy/util"
)

func emitEventEntityMgrAddEntity(event localevent.IEvent, entityMgr IEntityMgr, entity ec.Entity) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[EventEntityMgrAddEntity](delegate).OnEntityMgrAddEntity(entityMgr, entity)
		return true
	})
}

func emitEventEntityMgrRemoveEntity(event localevent.IEvent, entityMgr IEntityMgr, entity ec.Entity) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[EventEntityMgrRemoveEntity](delegate).OnEntityMgrRemoveEntity(entityMgr, entity)
		return true
	})
}

func emitEventEntityMgrEntityAddComponents(event localevent.IEvent, entityMgr IEntityMgr, entity ec.Entity, components []ec.Component) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[EventEntityMgrEntityAddComponents](delegate).OnEntityMgrEntityAddComponents(entityMgr, entity, components)
		return true
	})
}

func emitEventEntityMgrEntityRemoveComponent(event localevent.IEvent, entityMgr IEntityMgr, entity ec.Entity, component ec.Component) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[EventEntityMgrEntityRemoveComponent](delegate).OnEntityMgrEntityRemoveComponent(entityMgr, entity, component)
		return true
	})
}

func emitEventEntityMgrEntityFirstAccessComponent(event localevent.IEvent, entityMgr IEntityMgr, entity ec.Entity, component ec.Component) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[EventEntityMgrEntityFirstAccessComponent](delegate).OnEntityMgrEntityFirstAccessComponent(entityMgr, entity, component)
		return true
	})
}

func emitEventEntityMgrNotifyECTreeRemoveEntity(event localevent.IEvent, entityMgr IEntityMgr, entity ec.Entity) {
	if event == nil {
		panic("nil event")
	}
	localevent.UnsafeEvent(event).Emit(func(delegate util.IfaceCache) bool {
		util.Cache2Iface[eventEntityMgrNotifyECTreeRemoveEntity](delegate).onEntityMgrNotifyECTreeRemoveEntity(entityMgr, entity)
		return true
	})
}
