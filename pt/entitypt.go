package pt

import (
	"kit.golaxy.org/golaxy/ec"
	"kit.golaxy.org/golaxy/uid"
)

// EntityPt 实体原型
type EntityPt struct {
	Prototype string // 实体原型名称
	compPts   []ComponentPt
}

// Construct 创建实体
func (pt *EntityPt) Construct(options ...EntityOption) ec.Entity {
	opts := EntityOptions{}
	WithOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return pt.UnsafeConstruct(opts)
}

// UnsafeConstruct 不安全的创建实体，需要自己初始化所有选项
func (pt *EntityPt) UnsafeConstruct(options EntityOptions) ec.Entity {
	options.Prototype = pt.Prototype
	return pt.Assemble(ec.UnsafeNewEntity(options.EntityOptions), options.AssignCompId)
}

// Assemble 向实体安装组件
func (pt *EntityPt) Assemble(entity ec.Entity, assignCompId func(entity ec.Entity, compPt ComponentPt) uid.Id) ec.Entity {
	if entity == nil {
		return nil
	}

	if assignCompId == nil {
		assignCompId = func(entity ec.Entity, compPt ComponentPt) uid.Id { return uid.Nil }
	}

	for i := range pt.compPts {
		compPt := pt.compPts[i]
		entity.AddComponent(compPt.Name, compPt.Construct(assignCompId(entity, compPt)))
	}

	return entity
}
