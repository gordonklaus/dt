package bits

import "testing"

func TestZigZag(t *testing.T) {
	for i := int64(-256); i < 256; i++ {
		if zagzig(zigzag(i)) != i {
			t.Fatal()
		}
	}
}
