package data

import (
	"math"

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
	b.WriteUint32(math.Float32bits(f.X))
}

func (f *Float32Value) Read(b *bits.Buffer) error {
	x, err := b.ReadUint32()
	f.X = math.Float32frombits(x)
	return err
}

type Float64Value struct {
	Type *types.Float64Type
	X    float64
}

func NewFloat64Value(t *types.Float64Type) *Float64Value {
	return &Float64Value{Type: t}
}

func (f *Float64Value) Write(b *bits.Buffer) {
	b.WriteUint64(math.Float64bits(f.X))
}

func (f *Float64Value) Read(b *bits.Buffer) error {
	x, err := b.ReadUint64()
	f.X = math.Float64frombits(x)
	return err
}
