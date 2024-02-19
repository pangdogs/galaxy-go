//go:generate stringer -type RunningState
package service

import (
	"context"
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/core/util/iface"
	"git.golaxy.org/core/util/option"
	"git.golaxy.org/core/util/uid"
)

// Option 所有选项设置器
type Option struct{}

type (
	RunningHandler = generic.DelegateAction2[Context, RunningState] // 运行状态变化处理器
)

// ContextOptions 创建服务上下文的所有选项
type ContextOptions struct {
	CompositeFace  iface.Face[Context] // 扩展者，在扩展服务上下文自身能力时使用
	Context        context.Context     // 父Context
	AutoRecover    bool                // 是否开启panic时自动恢复
	ReportError    chan error          // panic时错误写入的error channel
	Name           string              // 服务名称
	PersistId      uid.Id              // 服务持久化Id
	EntityLib      pt.EntityLib        // 实体原型库
	PluginBundle   plugin.PluginBundle // 插件包
	RunningHandler RunningHandler      // 运行状态变化处理器
}

// Default 默认值
func (Option) Default() option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		Option{}.CompositeFace(iface.Face[Context]{})(o)
		Option{}.Context(nil)(o)
		Option{}.PanicHandling(false, nil)(o)
		Option{}.Name("")(o)
		Option{}.PersistId(uid.Nil)(o)
		Option{}.EntityLib(pt.DefaultEntityLib())(o)
		Option{}.PluginBundle(plugin.NewPluginBundle())(o)
		Option{}.RunningHandler(nil)(o)
	}
}

// CompositeFace 扩展者，在扩展服务上下文自身能力时使用
func (Option) CompositeFace(face iface.Face[Context]) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.CompositeFace = face
	}
}

// Context 父Context
func (Option) Context(ctx context.Context) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.Context = ctx
	}
}

// PanicHandling panic时的处理方式
func (Option) PanicHandling(autoRecover bool, reportError chan error) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.AutoRecover = autoRecover
		o.ReportError = reportError
	}
}

// Name 服务名称
func (Option) Name(name string) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.Name = name
	}
}

// PersistId 服务持久化Id
func (Option) PersistId(id uid.Id) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.PersistId = id
	}
}

// EntityLib 实体原型库
func (Option) EntityLib(lib pt.EntityLib) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.EntityLib = lib
	}
}

// PluginBundle 插件包
func (Option) PluginBundle(bundle plugin.PluginBundle) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.PluginBundle = bundle
	}
}

// RunningHandler 运行状态变化处理器
func (Option) RunningHandler(handler RunningHandler) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.RunningHandler = handler
	}
}
