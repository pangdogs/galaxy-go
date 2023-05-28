package ec

import (
	"database/sql/driver"
	"errors"
	"kit.golaxy.org/golaxy/internal"
)

func (comp *ComponentBehavior) MarshalText() ([]byte, error) {
	return ComponentSerializer.MarshalText(comp.composite)
}

func (comp *ComponentBehavior) UnmarshalText(b []byte) error {
	return ComponentSerializer.UnmarshalText(comp.composite, b)
}

func (comp *ComponentBehavior) MarshalBinary() ([]byte, error) {
	return ComponentSerializer.MarshalBinary(comp.composite)
}

func (comp *ComponentBehavior) UnmarshalBinary(b []byte) error {
	return ComponentSerializer.UnmarshalBinary(comp.composite, b)
}

func (comp *ComponentBehavior) Value() (driver.Value, error) {
	return ComponentSerializer.Value(comp.composite)
}

func (comp *ComponentBehavior) Scan(src interface{}) error {
	return ComponentSerializer.Scan(comp.composite, src)
}

var ComponentSerializer internal.Serializer[Component] = DefaultComponentSerializer{}

type DefaultComponentSerializer struct{}

func (DefaultComponentSerializer) MarshalText(comp Component) ([]byte, error) {
	return nil, errors.New("not support")
}

func (DefaultComponentSerializer) UnmarshalText(comp Component, b []byte) error {
	return errors.New("not support")
}

func (DefaultComponentSerializer) MarshalBinary(comp Component) ([]byte, error) {
	return nil, errors.New("not support")
}

func (DefaultComponentSerializer) UnmarshalBinary(comp Component, b []byte) error {
	return errors.New("not support")
}

func (DefaultComponentSerializer) Value(comp Component) (driver.Value, error) {
	return nil, errors.New("not support")
}

func (DefaultComponentSerializer) Scan(comp Component, src interface{}) error {
	return errors.New("not support")
}
