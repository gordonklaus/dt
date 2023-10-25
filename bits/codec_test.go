package bits

import (
	"bytes"
	"math"
	"testing"
)

func TestCodec(t *testing.T) {
	testValues(t, func(*Encoder, bool) {}, func(*Decoder, *bool) error { return nil }, false)
	testValues(t, (*Encoder).WriteBool, (*Decoder).ReadBool, false, true)
	testValues(t, (*Encoder).WriteUint32, (*Decoder).ReadUint32, 0, 1, 2, 3, math.MaxUint32)
	testValues(t, (*Encoder).WriteUint64, (*Decoder).ReadUint64, 0, 1, 2, 3, math.MaxUint64)
	testValues(t, (*Encoder).WriteInt64, (*Decoder).ReadInt64, 0, 1, 2, 3, math.MaxInt64, -1, -2, -3, math.MinInt64)
	testValues(t, (*Encoder).WriteVarUint, (*Decoder).ReadVarUint, 0, 1, 2, 3, math.MaxUint64)
	testValues(t, (*Encoder).WriteVarInt, (*Decoder).ReadVarInt, 0, 1, 2, 3, math.MaxInt64, -1, -2, -3, math.MinInt64)
	testValues(t, (*Encoder).WriteVarUint_4bit, (*Decoder).ReadVarUint_4bit, 0, 1, 2, 3, math.MaxUint64)
	testValues(t, (*Encoder).WriteFloat32, (*Decoder).ReadFloat32, 0, 1, 2, 3, math.MaxFloat32, -1, -2, -3, -math.MaxFloat32)
	testValues(t, (*Encoder).WriteFloat64, (*Decoder).ReadFloat64, 0, 1, 2, 3, math.MaxFloat64, -1, -2, -3, -math.MaxFloat64)
	testValues(t, (*Encoder).WriteString, (*Decoder).ReadString, "", "a", "847fqh938", "øˍðº,ßœ≥p«®£ª¢º˝ð-")
}

func testValues[T comparable](t *testing.T, w func(*Encoder, T), r func(*Decoder, *T) error, x ...T) {
	for i := range [8]int{} {
		for _, x := range x {
			e := NewEncoder()
			e.WriteSize(func() {
				for i := i; i > 0; i-- {
					e.WriteBool(true)
				}
				w(e, x)
			})
			d := NewDecoder(bytes.NewBuffer(e.Bytes()))
			d.SetLimit(e.Size())
			var y T
			if err := d.ReadSize(func() error {
				for i := i; i > 0; i-- {
					var b bool
					if err := d.ReadBool(&b); err != nil {
						return err
					}
					if !b {
						t.Fatalf("expected true from ReadBool (value=%v, offset=%d)", x, i)
					}
				}
				return r(d, &y)
			}); err != nil {
				t.Fatalf("ReadSize failed (value=%v, offset=%d): %v", x, i, err)
			}
			if d.Remaining() != 0 {
				t.Fatalf("%d bits remaining (value=%v, offset=%d)", d.Remaining(), x, i)
			}
			if x != y {
				t.Fatalf("expected %v, got %v (offset=%d)", x, y, i)
			}
		}
	}
}

func TestReadSizeSkipsExtraBits(t *testing.T) {
	e := NewEncoder()
	e.WriteSize(func() {
		for i := 0; i < 37; i++ {
			e.WriteBool(false)
		}
	})
	e.WriteBool(true)

	d := NewDecoder(bytes.NewBuffer(e.Bytes()))
	d.SetLimit(e.Size())
	if err := d.ReadSize(func() error {
		for i := 0; i < 11; i++ {
			var b bool
			if err := d.ReadBool(&b); err != nil {
				return err
			}
			if b {
				t.Fatalf("expected false from ReadBool")
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("ReadSize failed: %v", err)
	}
	var b bool
	if err := d.ReadBool(&b); err != nil {
		t.Fatalf("ReadBool failed: %v", err)
	}
	if !b {
		t.Fatalf("expected true from ReadBool")
	}
	if d.Remaining() != 0 {
		t.Fatalf("%d bits remaining", d.Remaining())
	}
}

func TestZigZag(t *testing.T) {
	for i := int64(-256); i < 256; i++ {
		if zagzig(zigzag(i)) != i {
			t.Fatal()
		}
	}
}
