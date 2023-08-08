package runtime

import (
	"context"
	"kit.golaxy.org/golaxy/localevent"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/uid"
	"kit.golaxy.org/golaxy/util"
	"kit.golaxy.org/golaxy/util/container"
)

// WithOption 所有选项设置器
type WithOption struct{}

type (
	Callback = func(ctx Context) // 回调函数
)

// ContextOptions 创建运行时上下文的所有选项
type ContextOptions struct {
	CompositeFace    util.Face[Context]                   // 扩展者，需要扩展运行时上下文自身能力时需要使用
	Context          context.Context                      // 父Context
	AutoRecover      bool                                 // 是否开启panic时自动恢复
	ReportError      chan error                           // panic时错误写入的error channel
	Name             string                               // 运行时名称
	PersistId        uid.Id                               // 运行时持久化Id
	PluginBundle     plugin.PluginBundle                  // 插件包
	StartedCb        Callback                             // 启动运行时回调函数
	StoppingCb       Callback                             // 开始停止运行时回调函数
	StoppedCb        Callback                             // 完全停止运行时回调函数
	FrameBeginCb     Callback                             // 帧开始时的回调函数
	FrameEndCb       Callback                             // 帧结束时的回调函数
	FaceAnyAllocator container.Allocator[util.FaceAny]    // 自定义FaceAny内存分配器，用于提高性能
	HookAllocator    container.Allocator[localevent.Hook] // 自定义Hook内存分配器，用于提高性能
}

// ContextOption 创建运行时上下文的选项设置器
type ContextOption func(o *ContextOptions)

// Default 默认值
func (WithOption) Default() ContextOption {
	return func(o *ContextOptions) {
		WithOption{}.Composite(util.Face[Context]{})(o)
		WithOption{}.Context(nil)(o)
		WithOption{}.AutoRecover(false)(o)
		WithOption{}.ReportError(nil)(o)
		WithOption{}.Name("")(o)
		WithOption{}.PersistId(util.Zero[uid.Id]())(o)
		WithOption{}.PluginBundle(nil)(o)
		WithOption{}.StartedCb(nil)(o)
		WithOption{}.StoppingCb(nil)(o)
		WithOption{}.StoppedCb(nil)(o)
		WithOption{}.FrameBeginCb(nil)(o)
		WithOption{}.FrameEndCb(nil)(o)
		WithOption{}.FaceAnyAllocator(container.DefaultAllocator[util.FaceAny]())(o)
		WithOption{}.HookAllocator(container.DefaultAllocator[localevent.Hook]())(o)
	}
}

// Composite 扩展者，需要扩展运行时上下文自身功能时需要使用
func (WithOption) Composite(face util.Face[Context]) ContextOption {
	return func(o *ContextOptions) {
		o.CompositeFace = face
	}
}

// Context 父Context
func (WithOption) Context(ctx context.Context) ContextOption {
	return func(o *ContextOptions) {
		o.Context = ctx
	}
}

// AutoRecover 是否开启panic时自动恢复
func (WithOption) AutoRecover(b bool) ContextOption {
	return func(o *ContextOptions) {
		o.AutoRecover = b
	}
}

// ReportError panic时错误写入的error channel
func (WithOption) ReportError(ch chan error) ContextOption {
	return func(o *ContextOptions) {
		o.ReportError = ch
	}
}

// Name 运行时名称
func (WithOption) Name(name string) ContextOption {
	return func(o *ContextOptions) {
		o.Name = name
	}
}

// PersistId 运行时持久化Id
func (WithOption) PersistId(id uid.Id) ContextOption {
	return func(o *ContextOptions) {
		o.PersistId = id
	}
}

// PluginBundle 插件包
func (WithOption) PluginBundle(bundle plugin.PluginBundle) ContextOption {
	return func(o *ContextOptions) {
		o.PluginBundle = bundle
	}
}

// StartedCb 启动运行时回调函数
func (WithOption) StartedCb(fn Callback) ContextOption {
	return func(o *ContextOptions) {
		o.StartedCb = fn
	}
}

// StoppingCb 开始停止运行时回调函数
func (WithOption) StoppingCb(fn Callback) ContextOption {
	return func(o *ContextOptions) {
		o.StoppingCb = fn
	}
}

// StoppedCb 完全停止运行时回调函数
func (WithOption) StoppedCb(fn Callback) ContextOption {
	return func(o *ContextOptions) {
		o.StoppedCb = fn
	}
}

// FrameBeginCb 帧更新开始时的回调函数
func (WithOption) FrameBeginCb(fn Callback) ContextOption {
	return func(o *ContextOptions) {
		o.FrameBeginCb = fn
	}
}

// FrameEndCb 帧更新结束时的回调函数
func (WithOption) FrameEndCb(fn Callback) ContextOption {
	return func(o *ContextOptions) {
		o.FrameEndCb = fn
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能
func (WithOption) FaceAnyAllocator(allocator container.Allocator[util.FaceAny]) ContextOption {
	return func(o *ContextOptions) {
		if allocator == nil {
			panic("nil allocator")
		}
		o.FaceAnyAllocator = allocator
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能
func (WithOption) HookAllocator(allocator container.Allocator[localevent.Hook]) ContextOption {
	return func(o *ContextOptions) {
		if allocator == nil {
			panic("nil allocator")
		}
		o.HookAllocator = allocator
	}
}
