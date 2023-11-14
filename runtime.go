package golaxy

import (
	"fmt"
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/event"
	"kit.golaxy.org/golaxy/internal/concurrent"
	"kit.golaxy.org/golaxy/runtime"
	"kit.golaxy.org/golaxy/util/generic"
	"kit.golaxy.org/golaxy/util/iface"
	"kit.golaxy.org/golaxy/util/option"
	"kit.golaxy.org/golaxy/util/uid"
	"sync/atomic"
)

// NewRuntime 创建运行时
func NewRuntime(ctx runtime.Context, settings ...option.Setting[RuntimeOptions]) Runtime {
	return UnsafeNewRuntime(ctx, option.Make(_RuntimeOption{}.Default(), settings...))
}

// Deprecated: UnsafeNewRuntime 内部创建运行时
func UnsafeNewRuntime(ctx runtime.Context, options RuntimeOptions) Runtime {
	if !options.CompositeFace.IsNil() {
		options.CompositeFace.Iface.init(ctx, options)
		return options.CompositeFace.Iface
	}

	runtime := &RuntimeBehavior{}
	runtime.init(ctx, options)

	return runtime.opts.CompositeFace.Iface
}

// Runtime 运行时接口
type Runtime interface {
	_Runtime
	concurrent.CurrentContextResolver
	concurrent.ConcurrentContextResolver
	concurrent.Callee
	Running

	// GetContext 获取运行时上下文
	GetContext() runtime.Context
}

type _Runtime interface {
	init(ctx runtime.Context, opts RuntimeOptions)
	getOptions() *RuntimeOptions
}

// RuntimeBehavior 运行时行为，在需要扩展运行时能力时，匿名嵌入至运行时结构体中
type RuntimeBehavior struct {
	ctx             runtime.Context
	opts            RuntimeOptions
	started         atomic.Bool
	hooksMap        map[uid.Id][3]event.Hook
	processQueue    chan _Task
	eventUpdate     event.Event
	eventLateUpdate event.Event
}

// GetContext 获取运行时上下文
func (rt *RuntimeBehavior) GetContext() runtime.Context {
	return rt.ctx
}

// ResolveContext 解析上下文
func (rt *RuntimeBehavior) ResolveContext() iface.Cache {
	return rt.ctx.ResolveContext()
}

// ResolveCurrentContext 解析当前上下文
func (rt *RuntimeBehavior) ResolveCurrentContext() iface.Cache {
	return rt.ctx.ResolveCurrentContext()
}

// ResolveConcurrentContext 解析多线程安全的上下文
func (rt *RuntimeBehavior) ResolveConcurrentContext() iface.Cache {
	return rt.ctx.ResolveConcurrentContext()
}

func (rt *RuntimeBehavior) init(ctx runtime.Context, opts RuntimeOptions) {
	if ctx == nil {
		panic(fmt.Errorf("%w: %w: ctx is nil", ErrRuntime, ErrArgs))
	}

	if !concurrent.UnsafeContext(ctx).SetPaired(true) {
		panic(fmt.Errorf("%w: context already paired", ErrRuntime))
	}

	rt.ctx = ctx
	rt.opts = opts

	if rt.opts.CompositeFace.IsNil() {
		rt.opts.CompositeFace = iface.MakeFace[Runtime](rt)
	}

	rt.hooksMap = make(map[uid.Id][3]event.Hook)
	rt.processQueue = make(chan _Task, rt.opts.ProcessQueueCapacity)

	rt.eventUpdate.Init(ctx.GetAutoRecover(), ctx.GetReportError(), event.EventRecursion_Disallow, ctx.GetHookAllocator(), nil)
	rt.eventLateUpdate.Init(ctx.GetAutoRecover(), ctx.GetReportError(), event.EventRecursion_Disallow, ctx.GetHookAllocator(), nil)

	runtime.UnsafeContext(ctx).SetFrame(rt.opts.Frame)
	runtime.UnsafeContext(ctx).SetCallee(rt.opts.CompositeFace.Iface)

	rt.changeRunningState(runtime.RunningState_Birth)

	if rt.opts.AutoRun {
		rt.opts.CompositeFace.Iface.Run()
	}
}

func (rt *RuntimeBehavior) getOptions() *RuntimeOptions {
	return &rt.opts
}

// OnEntityMgrAddEntity 事件处理器：实体管理器添加实体
func (rt *RuntimeBehavior) OnEntityMgrAddEntity(entityMgr runtime.IEntityMgr, entity ec.Entity) {
	rt.connectEntity(entity)
	rt.initEntity(entity)
}

// OnEntityMgrRemoveEntity 事件处理器：实体管理器删除实体
func (rt *RuntimeBehavior) OnEntityMgrRemoveEntity(entityMgr runtime.IEntityMgr, entity ec.Entity) {
	rt.disconnectEntity(entity)
	rt.shutEntity(entity)
}

// OnEntityMgrEntityFirstAccessComponent 事件处理器：实体管理器中的实体首次访问组件
func (rt *RuntimeBehavior) OnEntityMgrEntityFirstAccessComponent(entityMgr runtime.IEntityMgr, entity ec.Entity, component ec.Component) {
	_comp := ec.UnsafeComponent(component)

	if _comp.GetState() != ec.ComponentState_Attach {
		return
	}

	_comp.SetState(ec.ComponentState_Awake)

	if compAwake, ok := component.(LifecycleComponentAwake); ok {
		generic.CastAction0(compAwake.Awake).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}

	_comp.SetState(ec.ComponentState_Start)
}

// OnEntityMgrEntityAddComponents 事件处理器：实体管理器中的实体添加组件
func (rt *RuntimeBehavior) OnEntityMgrEntityAddComponents(entityMgr runtime.IEntityMgr, entity ec.Entity, components []ec.Component) {
	rt.addComponents(entity, components)
}

// OnEntityMgrEntityRemoveComponent 事件处理器：实体管理器中的实体删除组件
func (rt *RuntimeBehavior) OnEntityMgrEntityRemoveComponent(entityMgr runtime.IEntityMgr, entity ec.Entity, component ec.Component) {
	rt.removeComponent(component)
}

// OnEntityDestroySelf 事件处理器：实体销毁自身
func (rt *RuntimeBehavior) OnEntityDestroySelf(entity ec.Entity) {
	rt.ctx.GetEntityMgr().RemoveEntity(entity.GetId())
}

// OnComponentDestroySelf 事件处理器：组件销毁自身
func (rt *RuntimeBehavior) OnComponentDestroySelf(comp ec.Component) {
	comp.GetEntity().RemoveComponentById(comp.GetId())
}

func (rt *RuntimeBehavior) addComponents(entity ec.Entity, components []ec.Component) {
	switch entity.GetState() {
	case ec.EntityState_Awake, ec.EntityState_Start, ec.EntityState_Living:
	default:
		return
	}

	for i := range components {
		rt.connectComponent(components[i])
	}

	for i := range components {
		_comp := ec.UnsafeComponent(components[i])

		if _comp.GetState() != ec.ComponentState_Awake {
			continue
		}

		if compAwake, ok := components[i].(LifecycleComponentAwake); ok {
			generic.CastAction0(compAwake.Awake).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		}

		_comp.SetState(ec.ComponentState_Start)
	}

	switch entity.GetState() {
	case ec.EntityState_Awake, ec.EntityState_Start, ec.EntityState_Living:
	default:
		return
	}

	for i := range components {
		_comp := ec.UnsafeComponent(components[i])

		if _comp.GetState() != ec.ComponentState_Start {
			continue
		}

		if compStart, ok := components[i].(LifecycleComponentStart); ok {
			generic.CastAction0(compStart.Start).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		}

		_comp.SetState(ec.ComponentState_Living)
	}
}

func (rt *RuntimeBehavior) removeComponent(component ec.Component) {
	rt.disconnectComponent(component)

	if component.GetState() != ec.ComponentState_Shut {
		return
	}

	if compShut, ok := component.(LifecycleComponentShut); ok {
		generic.CastAction0(compShut.Shut).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}

	ec.UnsafeComponent(component).SetState(ec.ComponentState_Death)

	if compDestroy, ok := component.(LifecycleComponentDestroy); ok {
		generic.CastAction0(compDestroy.Destroy).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}
}

func (rt *RuntimeBehavior) connectEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Entry {
		return
	}

	var hooks [3]event.Hook

	if entityUpdate, ok := entity.(LifecycleEntityUpdate); ok {
		hooks[0] = event.BindEvent[LifecycleEntityUpdate](&rt.eventUpdate, entityUpdate)
	}
	if entityLateUpdate, ok := entity.(LifecycleEntityLateUpdate); ok {
		hooks[1] = event.BindEvent[LifecycleEntityLateUpdate](&rt.eventLateUpdate, entityLateUpdate)
	}
	hooks[2] = ec.BindEventEntityDestroySelf(ec.UnsafeEntity(entity), rt)

	rt.hooksMap[entity.GetId()] = hooks

	entity.RangeComponents(func(comp ec.Component) bool {
		rt.connectComponent(comp)
		return true
	})

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Awake)
}

func (rt *RuntimeBehavior) disconnectEntity(entity ec.Entity) {
	entity.RangeComponents(func(comp ec.Component) bool {
		rt.disconnectComponent(comp)
		return true
	})

	entityId := entity.GetId()

	hooks, ok := rt.hooksMap[entityId]
	if ok {
		delete(rt.hooksMap, entityId)

		for i := range hooks {
			hooks[i].Unbind()
		}
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Shut)
}

func (rt *RuntimeBehavior) connectComponent(comp ec.Component) {
	if comp.GetState() != ec.ComponentState_Attach {
		return
	}

	var hooks [3]event.Hook

	if compUpdate, ok := comp.(LifecycleComponentUpdate); ok {
		hooks[0] = event.BindEvent[LifecycleComponentUpdate](&rt.eventUpdate, compUpdate)
	}
	if compLateUpdate, ok := comp.(LifecycleComponentLateUpdate); ok {
		hooks[1] = event.BindEvent[LifecycleComponentLateUpdate](&rt.eventLateUpdate, compLateUpdate)
	}
	hooks[2] = ec.BindEventComponentDestroySelf(ec.UnsafeComponent(comp), rt)

	rt.hooksMap[comp.GetId()] = hooks

	ec.UnsafeComponent(comp).SetState(ec.ComponentState_Awake)
}

func (rt *RuntimeBehavior) disconnectComponent(comp ec.Component) {
	compId := comp.GetId()

	hooks, ok := rt.hooksMap[compId]
	if ok {
		delete(rt.hooksMap, compId)

		for i := range hooks {
			hooks[i].Unbind()
		}
	}

	ec.UnsafeComponent(comp).SetState(ec.ComponentState_Shut)
}

func (rt *RuntimeBehavior) initEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Awake {
		return
	}

	if entityAwake, ok := entity.(LifecycleEntityAwake); ok {
		generic.CastAction0(entityAwake.Awake).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}

	if entity.GetState() != ec.EntityState_Awake {
		return
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Awake {
			return true
		}

		if compAwake, ok := comp.(LifecycleComponentAwake); ok {
			generic.CastAction0(compAwake.Awake).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		}

		_comp.SetState(ec.ComponentState_Start)

		return entity.GetState() == ec.EntityState_Awake
	})

	if entity.GetState() != ec.EntityState_Awake {
		return
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Start {
			return true
		}

		if compStart, ok := comp.(LifecycleComponentStart); ok {
			generic.CastAction0(compStart.Start).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		}

		_comp.SetState(ec.ComponentState_Living)

		return entity.GetState() == ec.EntityState_Awake
	})

	if entity.GetState() != ec.EntityState_Awake {
		return
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Start)

	if entityStart, ok := entity.(LifecycleEntityStart); ok {
		generic.CastAction0(entityStart.Start).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}

	if entity.GetState() != ec.EntityState_Start {
		return
	}

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Living)
}

func (rt *RuntimeBehavior) shutEntity(entity ec.Entity) {
	if entity.GetState() != ec.EntityState_Shut {
		return
	}

	if entityShut, ok := entity.(LifecycleEntityShut); ok {
		generic.CastAction0(entityShut.Shut).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}

	entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetState() != ec.ComponentState_Shut {
			return true
		}

		if compShut, ok := comp.(LifecycleComponentShut); ok {
			generic.CastAction0(compShut.Shut).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		}

		_comp.SetState(ec.ComponentState_Death)

		if compDestroy, ok := comp.(LifecycleComponentDestroy); ok {
			generic.CastAction0(compDestroy.Destroy).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		}

		return true
	})

	ec.UnsafeEntity(entity).SetState(ec.EntityState_Death)

	if entityDestroy, ok := entity.(LifecycleEntityDestroy); ok {
		generic.CastAction0(entityDestroy.Destroy).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}
}
