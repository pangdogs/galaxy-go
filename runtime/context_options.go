package runtime

import (
	"context"
	"kit.golaxy.org/golaxy/localevent"
	"kit.golaxy.org/golaxy/plugin"
	"kit.golaxy.org/golaxy/uid"
	"kit.golaxy.org/golaxy/util"
	"kit.golaxy.org/golaxy/util/container"
)

// Options 创建运行时上下文的所有选项
type Options struct {
	CompositeFace      util.Face[Context]                   // 扩展者，需要扩展运行时上下文自身能力时需要使用
	Context            context.Context                      // 父Context
	AutoRecover        bool                                 // 是否开启panic时自动恢复
	ReportError        chan error                           // panic时错误写入的error channel
	Name               string                               // 运行时名称
	PersistId          uid.Id                               // 运行时持久化Id
	PluginBundle       plugin.PluginBundle                  // 插件包
	StartedCallback    func(ctx Context)                    // 启动运行时回调函数
	StoppingCallback   func(ctx Context)                    // 开始停止运行时回调函数
	StoppedCallback    func(ctx Context)                    // 完全停止运行时回调函数
	FrameBeginCallback func(ctx Context)                    // 帧开始时的回调函数
	FrameEndCallback   func(ctx Context)                    // 帧结束时的回调函数
	FaceAnyAllocator   container.Allocator[util.FaceAny]    // 自定义FaceAny内存分配器，用于提高性能
	HookAllocator      container.Allocator[localevent.Hook] // 自定义Hook内存分配器，用于提高性能
}

// Option 创建运行时上下文的选项设置器
type Option func(o *Options)

// WithOption 创建运行时上下文的所有选项设置器
type WithOption struct{}

// Default 默认值
func (WithOption) Default() Option {
	return func(o *Options) {
		WithOption{}.Composite(util.Face[Context]{})(o)
		WithOption{}.Context(nil)(o)
		WithOption{}.AutoRecover(false)(o)
		WithOption{}.ReportError(nil)(o)
		WithOption{}.Name("")(o)
		WithOption{}.PersistId(util.Zero[uid.Id]())(o)
		WithOption{}.PluginBundle(nil)(o)
		WithOption{}.StartedCallback(nil)(o)
		WithOption{}.StoppingCallback(nil)(o)
		WithOption{}.StoppedCallback(nil)(o)
		WithOption{}.FrameBeginCallback(nil)(o)
		WithOption{}.FrameEndCallback(nil)(o)
		WithOption{}.FaceAnyAllocator(container.DefaultAllocator[util.FaceAny]())(o)
		WithOption{}.HookAllocator(container.DefaultAllocator[localevent.Hook]())(o)
	}
}

// Composite 扩展者，需要扩展运行时上下文自身功能时需要使用
func (WithOption) Composite(face util.Face[Context]) Option {
	return func(o *Options) {
		o.CompositeFace = face
	}
}

// Context 父Context
func (WithOption) Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

// AutoRecover 是否开启panic时自动恢复
func (WithOption) AutoRecover(b bool) Option {
	return func(o *Options) {
		o.AutoRecover = b
	}
}

// ReportError panic时错误写入的error channel
func (WithOption) ReportError(ch chan error) Option {
	return func(o *Options) {
		o.ReportError = ch
	}
}

// Name 运行时名称
func (WithOption) Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

// PersistId 运行时持久化Id
func (WithOption) PersistId(id uid.Id) Option {
	return func(o *Options) {
		o.PersistId = id
	}
}

// PluginBundle 插件包
func (WithOption) PluginBundle(bundle plugin.PluginBundle) Option {
	return func(o *Options) {
		o.PluginBundle = bundle
	}
}

// StartedCallback 启动运行时回调函数
func (WithOption) StartedCallback(fn func(ctx Context)) Option {
	return func(o *Options) {
		o.StartedCallback = fn
	}
}

// StoppingCallback 开始停止运行时回调函数
func (WithOption) StoppingCallback(fn func(ctx Context)) Option {
	return func(o *Options) {
		o.StoppingCallback = fn
	}
}

// StoppedCallback 完全停止运行时回调函数
func (WithOption) StoppedCallback(fn func(ctx Context)) Option {
	return func(o *Options) {
		o.StoppedCallback = fn
	}
}

// FrameBeginCallback 帧更新开始时的回调函数
func (WithOption) FrameBeginCallback(fn func(ctx Context)) Option {
	return func(o *Options) {
		o.FrameBeginCallback = fn
	}
}

// FrameEndCallback 帧更新结束时的回调函数
func (WithOption) FrameEndCallback(fn func(ctx Context)) Option {
	return func(o *Options) {
		o.FrameEndCallback = fn
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能
func (WithOption) FaceAnyAllocator(allocator container.Allocator[util.FaceAny]) Option {
	return func(o *Options) {
		if allocator == nil {
			panic("nil allocator")
		}
		o.FaceAnyAllocator = allocator
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能
func (WithOption) HookAllocator(allocator container.Allocator[localevent.Hook]) Option {
	return func(o *Options) {
		if allocator == nil {
			panic("nil allocator")
		}
		o.HookAllocator = allocator
	}
}
