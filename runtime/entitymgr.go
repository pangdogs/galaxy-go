package runtime

import (
	"errors"
	"fmt"
	"github.com/pangdogs/galaxy/ec"
	"github.com/pangdogs/galaxy/localevent"
	"github.com/pangdogs/galaxy/util"
	"github.com/pangdogs/galaxy/util/container"
)

// IEntityMgr 实体管理器接口
type IEntityMgr interface {
	// GetRuntimeCtx 获取运行时上下文
	GetRuntimeCtx() Context

	// GetEntity 查询实体
	GetEntity(id int64) (ec.Entity, bool)

	// RangeEntities 遍历所有实体
	RangeEntities(func(entity ec.Entity) bool)

	// ReverseRangeEntities 反向遍历所有实体
	ReverseRangeEntities(func(entity ec.Entity) bool)

	// GetEntityCount 获取实体数量
	GetEntityCount() int

	// AddEntity 添加实体
	AddEntity(entity ec.Entity) error

	// RemoveEntity 删除实体
	RemoveEntity(id int64)

	// EventEntityMgrAddEntity 事件：实体管理器中添加实体
	EventEntityMgrAddEntity() localevent.IEvent

	// EventEntityMgrRemoveEntity 事件：实体管理器中删除实体
	EventEntityMgrRemoveEntity() localevent.IEvent

	// EventEntityMgrEntityAddComponents 事件：实体管理器中的实体添加组件
	EventEntityMgrEntityAddComponents() localevent.IEvent

	// EventEntityMgrEntityRemoveComponent 事件：实体管理器中的实体删除组件
	EventEntityMgrEntityRemoveComponent() localevent.IEvent

	// EventEntityMgrEntityFirstAccessComponent 事件：实体管理器中的实体首次访问组件
	EventEntityMgrEntityFirstAccessComponent() localevent.IEvent

	eventEntityMgrNotifyECTreeRemoveEntity() localevent.IEvent
}

// EntityMgr 实体管理器
type EntityMgr struct {
	runtimeCtx                               Context
	entityMap                                map[int64]_EntityInfo
	entityList                               container.List[util.FaceAny]
	eventEntityMgrAddEntity                  localevent.Event
	eventEntityMgrRemoveEntity               localevent.Event
	eventEntityMgrEntityAddComponents        localevent.Event
	eventEntityMgrEntityRemoveComponent      localevent.Event
	eventEntityMgrEntityFirstAccessComponent localevent.Event
	_eventEntityMgrNotifyECTreeRemoveEntity  localevent.Event
	inited                                   bool
}

// Init 初始化实体管理器
func (entityMgr *EntityMgr) Init(runtimeCtx Context) {
	if runtimeCtx == nil {
		panic("nil runtimeCtx")
	}

	if entityMgr.inited {
		panic("repeated init entity manager")
	}

	entityMgr.entityList.Init(runtimeCtx.GetFaceCache(), runtimeCtx)
	entityMgr.entityMap = map[int64]_EntityInfo{}

	entityMgr.eventEntityMgrAddEntity.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Discard, runtimeCtx.GetHookCache(), runtimeCtx)
	entityMgr.eventEntityMgrRemoveEntity.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Discard, runtimeCtx.GetHookCache(), runtimeCtx)
	entityMgr.eventEntityMgrEntityAddComponents.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Discard, runtimeCtx.GetHookCache(), runtimeCtx)
	entityMgr.eventEntityMgrEntityRemoveComponent.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Discard, runtimeCtx.GetHookCache(), runtimeCtx)
	entityMgr.eventEntityMgrEntityFirstAccessComponent.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Discard, runtimeCtx.GetHookCache(), runtimeCtx)
	entityMgr._eventEntityMgrNotifyECTreeRemoveEntity.Init(runtimeCtx.GetAutoRecover(), runtimeCtx.GetReportError(), localevent.EventRecursion_Discard, runtimeCtx.GetHookCache(), runtimeCtx)

	entityMgr.inited = true
}

// GetRuntimeCtx 获取运行时上下文
func (entityMgr *EntityMgr) GetRuntimeCtx() Context {
	return entityMgr.runtimeCtx
}

// GetEntity 查询实体
func (entityMgr *EntityMgr) GetEntity(id int64) (ec.Entity, bool) {
	e, ok := entityMgr.entityMap[id]
	if !ok {
		return nil, false
	}

	if e.Element.Escaped() {
		return nil, false
	}

	return util.Cache2Iface[ec.Entity](e.Element.Value.Cache), true
}

// RangeEntities 遍历所有实体
func (entityMgr *EntityMgr) RangeEntities(fun func(entity ec.Entity) bool) {
	if fun == nil {
		return
	}

	entityMgr.entityList.Traversal(func(e *container.Element[util.FaceAny]) bool {
		return fun(util.Cache2Iface[ec.Entity](e.Value.Cache))
	})
}

// ReverseRangeEntities 反向遍历所有实体
func (entityMgr *EntityMgr) ReverseRangeEntities(fun func(entity ec.Entity) bool) {
	if fun == nil {
		return
	}

	entityMgr.entityList.ReverseTraversal(func(e *container.Element[util.FaceAny]) bool {
		return fun(util.Cache2Iface[ec.Entity](e.Value.Cache))
	})
}

// GetEntityCount 获取实体数量
func (entityMgr *EntityMgr) GetEntityCount() int {
	return entityMgr.entityList.Len()
}

// AddEntity 添加实体
func (entityMgr *EntityMgr) AddEntity(entity ec.Entity) error {
	if entity == nil {
		return errors.New("nil entity")
	}

	_entity := ec.UnsafeEntity(entity)

	if _entity.GetContext() != util.NilIfaceCache {
		return errors.New("entity context has already been setup")
	}

	if entity.GetID() <= 0 {
		_entity.SetID(entityMgr.runtimeCtx.GetServiceCtx().GenUID())
	}

	_entity.SetContext(util.Iface2Cache[Context](entityMgr.runtimeCtx))
	_entity.SetGCCollector(entityMgr.runtimeCtx)

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetID() <= 0 {
			_comp.SetID(entityMgr.runtimeCtx.GetServiceCtx().GenUID())
		}

		_comp.SetPrimary(true)

		return true
	})

	if _, ok := entityMgr.entityMap[entity.GetID()]; ok {
		return fmt.Errorf("repeated entity '%d' in this runtime context", entity.GetID())
	}

	_entity.SetAdding(true)
	defer _entity.SetAdding(false)

	entityInfo := _EntityInfo{}
	entityInfo.Hooks[0] = localevent.BindEvent[ec.EventCompMgrAddComponents](entity.EventCompMgrAddComponents(), entityMgr)
	entityInfo.Hooks[1] = localevent.BindEvent[ec.EventCompMgrRemoveComponent](entity.EventCompMgrRemoveComponent(), entityMgr)

	if _entity.GetOptions().EnableComponentAwakeByAccess {
		entityInfo.Hooks[2] = localevent.BindEvent[ec.EventCompMgrFirstAccessComponent](entity.EventCompMgrFirstAccessComponent(), entityMgr)
	}

	entityInfo.Element = entityMgr.entityList.PushBack(util.FaceAny{
		Iface: entity,
		Cache: util.Iface2Cache(entity),
	})

	entityMgr.entityMap[entity.GetID()] = entityInfo
	entityMgr.runtimeCtx.CollectGC(_entity.GetInnerGC())

	emitEventEntityMgrAddEntity(&entityMgr.eventEntityMgrAddEntity, entityMgr, entity)

	return nil
}

// RemoveEntity 删除实体
func (entityMgr *EntityMgr) RemoveEntity(id int64) {
	entityInfo, ok := entityMgr.entityMap[id]
	if !ok {
		return
	}

	entity := ec.UnsafeEntity(util.Cache2Iface[ec.Entity](entityInfo.Element.Value.Cache))
	if entity.GetAdding() || entity.GetRemoving() || entity.GetInitialing() || entity.GetShutting() {
		return
	}

	entity.SetRemoving(true)
	defer entity.SetRemoving(false)

	emitEventEntityMgrNotifyECTreeRemoveEntity(&entityMgr._eventEntityMgrNotifyECTreeRemoveEntity, entityMgr, entity.Entity)
	entityMgr.runtimeCtx.GetECTree().RemoveChild(id)

	delete(entityMgr.entityMap, id)
	entityInfo.Element.Escape()

	for i := range entityInfo.Hooks {
		entityInfo.Hooks[i].Unbind()
	}

	emitEventEntityMgrRemoveEntity(&entityMgr.eventEntityMgrRemoveEntity, entityMgr, entity.Entity)
}

// EventEntityMgrAddEntity 事件：运行时上下文添加实体
func (entityMgr *EntityMgr) EventEntityMgrAddEntity() localevent.IEvent {
	return &entityMgr.eventEntityMgrAddEntity
}

// EventEntityMgrRemoveEntity 事件：运行时上下文删除实体
func (entityMgr *EntityMgr) EventEntityMgrRemoveEntity() localevent.IEvent {
	return &entityMgr.eventEntityMgrRemoveEntity
}

// EventEntityMgrEntityAddComponents 事件：运行时上下文中的实体添加组件
func (entityMgr *EntityMgr) EventEntityMgrEntityAddComponents() localevent.IEvent {
	return &entityMgr.eventEntityMgrEntityAddComponents
}

// EventEntityMgrEntityRemoveComponent 事件：运行时上下文中的实体删除组件
func (entityMgr *EntityMgr) EventEntityMgrEntityRemoveComponent() localevent.IEvent {
	return &entityMgr.eventEntityMgrEntityRemoveComponent
}

// EventEntityMgrEntityFirstAccessComponent 事件：运行时上下文中的实体首次访问组件
func (entityMgr *EntityMgr) EventEntityMgrEntityFirstAccessComponent() localevent.IEvent {
	return &entityMgr.eventEntityMgrEntityFirstAccessComponent
}

func (entityMgr *EntityMgr) eventEntityMgrNotifyECTreeRemoveEntity() localevent.IEvent {
	return &entityMgr._eventEntityMgrNotifyECTreeRemoveEntity
}

// OnCompMgrAddComponents 事件回调：实体的组件管理器加入一些组件
func (entityMgr *EntityMgr) OnCompMgrAddComponents(entity ec.Entity, components []ec.Component) {
	for i := range components {
		if components[i].GetID() <= 0 {
			ec.UnsafeComponent(components[i]).SetID(entityMgr.runtimeCtx.GetServiceCtx().GenUID())
		}
	}
	emitEventEntityMgrEntityAddComponents(&entityMgr.eventEntityMgrEntityAddComponents, entityMgr, entity, components)
}

// OnCompMgrRemoveComponent 事件回调：实体的组件管理器删除组件
func (entityMgr *EntityMgr) OnCompMgrRemoveComponent(entity ec.Entity, component ec.Component) {
	emitEventEntityMgrEntityRemoveComponent(&entityMgr.eventEntityMgrEntityRemoveComponent, entityMgr, entity, component)
}

// OnCompMgrFirstAccessComponent 事件回调：实体的组件管理器首次访问组件
func (entityMgr *EntityMgr) OnCompMgrFirstAccessComponent(entity ec.Entity, component ec.Component) {
	emitEventEntityMgrEntityFirstAccessComponent(&entityMgr.eventEntityMgrEntityFirstAccessComponent, entityMgr, entity, component)
}
