package runtime

import (
	"context"
	"fmt"
	"git.golaxy.org/core/event"
	"git.golaxy.org/core/internal/exception"
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/util/container"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/core/util/iface"
	"git.golaxy.org/core/util/option"
	"git.golaxy.org/core/util/uid"
)

type _ContextOption struct{}

type (
	RunningHandler = generic.DelegateAction2[Context, RunningState] // 运行状态变化处理器
)

// ContextOptions 创建运行时上下文的所有选项
type ContextOptions struct {
	CompositeFace    iface.Face[Context]                // 扩展者，在扩展运行时上下文自身能力时使用
	Context          context.Context                    // 父Context
	AutoRecover      bool                               // 是否开启panic时自动恢复
	ReportError      chan error                         // panic时错误写入的error channel
	Name             string                             // 运行时名称
	PersistId        uid.Id                             // 运行时持久化Id
	PluginBundle     plugin.PluginBundle                // 插件包
	RunningHandler   RunningHandler                     // 运行状态变化处理器
	FaceAnyAllocator container.Allocator[iface.FaceAny] // 自定义FaceAny内存分配器，用于提高性能
	HookAllocator    container.Allocator[event.Hook]    // 自定义Hook内存分配器，用于提高性能
}

// Default 默认值
func (_ContextOption) Default() option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		_ContextOption{}.CompositeFace(iface.Face[Context]{})(o)
		_ContextOption{}.Context(nil)(o)
		_ContextOption{}.AutoRecover(false)(o)
		_ContextOption{}.ReportError(nil)(o)
		_ContextOption{}.Name("")(o)
		_ContextOption{}.PersistId(uid.Nil)(o)
		_ContextOption{}.PluginBundle(plugin.NewPluginBundle())(o)
		_ContextOption{}.RunningHandler(nil)(o)
		_ContextOption{}.FaceAnyAllocator(container.DefaultAllocator[iface.FaceAny]())(o)
		_ContextOption{}.HookAllocator(container.DefaultAllocator[event.Hook]())(o)
	}
}

// CompositeFace 扩展者，在扩展运行时上下文自身能力时使用
func (_ContextOption) CompositeFace(face iface.Face[Context]) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.CompositeFace = face
	}
}

// Context 父Context
func (_ContextOption) Context(ctx context.Context) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.Context = ctx
	}
}

// AutoRecover 是否开启panic时自动恢复
func (_ContextOption) AutoRecover(b bool) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.AutoRecover = b
	}
}

// ReportError panic时错误写入的error channel
func (_ContextOption) ReportError(ch chan error) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.ReportError = ch
	}
}

// Name 运行时名称
func (_ContextOption) Name(name string) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.Name = name
	}
}

// PersistId 运行时持久化Id
func (_ContextOption) PersistId(id uid.Id) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.PersistId = id
	}
}

// PluginBundle 插件包
func (_ContextOption) PluginBundle(bundle plugin.PluginBundle) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.PluginBundle = bundle
	}
}

// RunningHandler 运行状态变化处理器
func (_ContextOption) RunningHandler(handler RunningHandler) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.RunningHandler = handler
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能
func (_ContextOption) FaceAnyAllocator(allocator container.Allocator[iface.FaceAny]) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		if allocator == nil {
			panic(fmt.Errorf("%w: %w: allocator is nil", ErrContext, exception.ErrArgs))
		}
		o.FaceAnyAllocator = allocator
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能
func (_ContextOption) HookAllocator(allocator container.Allocator[event.Hook]) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		if allocator == nil {
			panic(fmt.Errorf("%w: %w: allocator is nil", ErrContext, exception.ErrArgs))
		}
		o.HookAllocator = allocator
	}
}
