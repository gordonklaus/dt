package data

import (
	"fmt"
)

type BasicType struct {
	Kind Kind
}

func NewBasicType(kind Kind) *BasicType {
	return &BasicType{Kind: kind}
}

func (i BasicType) NewValue() Value {
	switch i.Kind {

	}
	panic(fmt.Sprintf("invalid Kind %d", i.Kind))
}
