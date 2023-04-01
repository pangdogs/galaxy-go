package service

import "kit.golaxy.org/golaxy/util"

// GetComposite 获取服务上下文的扩展者
func GetComposite[T any](ctx Context) T {
	return util.Cache2Iface[T](ctx.getOptions().CompositeFace.Cache)
}
