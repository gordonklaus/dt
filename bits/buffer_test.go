package bits

import "testing"

func TestVarUint(t *testing.T) {
	b := NewBuffer()
	x := uint64(1<<64 - 1)
	b.WriteVarUint(x)
	var y uint64
	if err := b.ReadVarUint(&y); err != nil {
		t.Fatal(err)
	}
	if x != y {
		t.Fatalf("expected %d, got %d", x, y)
	}
}

func TestVarInt(t *testing.T) {
	b := NewBuffer()
	x := int64(1<<63 - 1)
	b.WriteVarInt(x)
	var y int64
	if err := b.ReadVarInt(&y); err != nil {
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
