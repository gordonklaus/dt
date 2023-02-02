package data

import (
	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type Float32Value struct {
	Type *types.Float32Type
	X    float32
}

func NewFloat32Value(t *types.Float32Type) *Float32Value {
	return &Float32Value{Type: t}
}

func (f *Float32Value) Write(b *bits.Buffer) {
	b.WriteFloat32(f.X)
}

func (f *Float32Value) Read(b *bits.Buffer) error {
	return b.ReadFloat32(&f.X)
}

type Float64Value struct {
	Type *types.Float64Type
	X    float64
}

func NewFloat64Value(t *types.Float64Type) *Float64Value {
	return &Float64Value{Type: t}
}

func (f *Float64Value) Write(b *bits.Buffer) {
	b.WriteFloat64(f.X)
}

func (f *Float64Value) Read(b *bits.Buffer) error {
	return b.ReadFloat64(&f.X)
}
