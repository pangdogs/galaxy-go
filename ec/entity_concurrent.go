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

package ec

import (
	"context"
	"fmt"
	"git.golaxy.org/core/internal/ictx"
	"git.golaxy.org/core/utils/uid"
)

// ConcurrentEntity 多线程安全的实体接口
type ConcurrentEntity interface {
	iConcurrentEntity
	iContext
	ictx.ConcurrentContextProvider
	fmt.Stringer

	// GetId 获取实体Id
	GetId() uid.Id
	// GetPrototype 获取实体原型
	GetPrototype() string
}

type iContext interface {
	context.Context

	// Terminated 已停止
	Terminated() <-chan struct{}
}

type iConcurrentEntity interface {
	getEntity() Entity
}

// Terminated 已停止
func (entity *EntityBehavior) Terminated() <-chan struct{} {
	return entity.terminated
}

func (entity *EntityBehavior) getEntity() Entity {
	return entity.opts.CompositeFace.Iface
}
