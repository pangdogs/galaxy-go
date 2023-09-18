package plugin

import (
	"fmt"
	"kit.golaxy.org/golaxy/util"
	"kit.golaxy.org/golaxy/util/iface"
)

// PluginResolver 插件解析器
type PluginResolver interface {
	// ResolvePlugin 解析插件
	ResolvePlugin(name string) (PluginInfo, bool)
}

// Fetch 获取插件。
//
//	@param pluginResolver 插件解析器。
//	@param name 插件名称。
func Fetch[T any](pluginResolver PluginResolver, name string) T {
	if pluginResolver == nil {
		panic(fmt.Errorf("%w: pluginResolver is nil", ErrPlugin))
	}

	pluginInfo, ok := pluginResolver.ResolvePlugin(name)
	if !ok {
		panic(fmt.Errorf("%w: %q not installed", ErrPlugin, name))
	}

	if !pluginInfo.Active {
		panic(fmt.Errorf("%w: %q not actived", ErrPlugin, name))
	}

	return iface.Cache2Iface[T](pluginInfo.Face.Cache)
}

// Access 访问插件
//
//	@param pluginResolver 插件解析器。
//	@param name 插件名称。
func Access[T any](pluginResolver PluginResolver, name string) (T, bool) {
	if pluginResolver == nil {
		return util.Zero[T](), false
	}

	pluginInfo, ok := pluginResolver.ResolvePlugin(name)
	if !ok {
		return util.Zero[T](), false
	}

	if !pluginInfo.Active {
		return util.Zero[T](), false
	}

	return iface.Cache2Iface[T](pluginInfo.Face.Cache), true
}
