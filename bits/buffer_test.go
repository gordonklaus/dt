package bits

import (
	"bytes"
	"testing"
)

func TestVarUint(t *testing.T) {
	e := NewEncoder()
	x := uint64(1<<64 - 1)
	e.WriteVarUint(x)
	var y uint64
	if err := NewDecoder(bytes.NewBuffer(e.Bytes())).ReadVarUint(&y); err != nil {
		t.Fatal(err)
	}
	if x != y {
		t.Fatalf("expected %d, got %d", x, y)
	}
}

func TestVarInt(t *testing.T) {
	e := NewEncoder()
	x := int64(1<<63 - 1)
	e.WriteVarInt(x)
	var y int64
	if err := NewDecoder(bytes.NewBuffer(e.Bytes())).ReadVarInt(&y); err != nil {
		t.Fatal(err)
	}
	if x != y {
		t.Fatalf("expected %d, got %d", x, y)
	}
}

func TestZigZag(t *testing.T) {
	for i := int64(-256); i < 256; i++ {
		if zagzig(zigzag(i)) != i {
			t.Fatal()
		}
	}
}
