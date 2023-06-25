package data

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
	"github.com/gordonklaus/data/types"
)

type FloatValue struct {
	Type *types.FloatType
	x    float64
}

func NewFloatValue(t *types.FloatType) *FloatValue {
	return &FloatValue{Type: t}
}

func (f *FloatValue) SetFloat32(x float32) {
	if f.Type.Size != 32 {
		panic(fmt.Sprintf("float has %d bits", f.Type.Size))
	}
	f.x = float64(x)
}

func (f *FloatValue) SetFloat64(x float64) {
	if f.Type.Size != 64 {
		panic(fmt.Sprintf("float has %d bits", f.Type.Size))
	}
	f.x = x
}

func (f *FloatValue) Write(b *bits.Buffer) {
	if f.Type.Size == 32 {
		b.WriteFloat32(float32(f.x))
	} else {
		b.WriteFloat64(f.x)
	}
}

func (f *FloatValue) Read(b *bits.Buffer) error {
	if f.Type.Size == 32 {
		var x float32
		err := b.ReadFloat32(&x)
		f.x = float64(x)
		return err
	} else {
		return b.ReadFloat64(&f.x)
	}
}
