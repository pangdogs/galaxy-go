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

package runtime

import (
	"fmt"
	"git.golaxy.org/core/internal/gctx"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/exception"
	"git.golaxy.org/core/utils/iface"
	"git.golaxy.org/core/utils/uid"
)

// ConcurrentContextProvider 多线程安全的上下文提供者
type ConcurrentContextProvider = gctx.ConcurrentContextProvider

// ConcurrentContext 多线程安全的运行时上下文接口
type ConcurrentContext interface {
	gctx.ConcurrentContextProvider
	gctx.Context
	async.Caller
	fmt.Stringer

	// GetName 获取名称
	GetName() string
	// GetId 获取运行时Id
	GetId() uid.Id
}

// Concurrent 获取多线程安全的运行时上下文
func Concurrent(provider gctx.ConcurrentContextProvider) ConcurrentContext {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	return iface.Cache2Iface[Context](provider.GetConcurrentContext())
}

func getServiceContext(provider gctx.ConcurrentContextProvider) service.Context {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	ctx := iface.Cache2Iface[Context](provider.GetConcurrentContext())
	if ctx == nil {
		return nil
	}
	return ctx.getServiceCtx()
}

func getCaller(provider gctx.ConcurrentContextProvider) async.Caller {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	return iface.Cache2Iface[Context](provider.GetConcurrentContext())
}
