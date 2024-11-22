/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package define

import (
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/utils/generic"
)

// RuntimePluginInterface 定义运行时插件接口，支持运行时上下文，通常用于为同类插件的不同实现提供统一的接口
func RuntimePluginInterface[PLUGIN_IFACE any]() RuntimePluginInterfaceDefinition[PLUGIN_IFACE] {
	plug := definePluginInterface[PLUGIN_IFACE]()

	return RuntimePluginInterfaceDefinition[PLUGIN_IFACE]{
		Name:  plug.Name,
		Using: func(rtCtx runtime.Context) PLUGIN_IFACE { return plug.Using(rtCtx) },
	}
}

// RuntimePluginInterfaceDefinition 运行时插件接口定义
type RuntimePluginInterfaceDefinition[PLUGIN_IFACE any] struct {
	Name  string                                       // 插件名称
	Using generic.Func1[runtime.Context, PLUGIN_IFACE] // 使用插件
}
