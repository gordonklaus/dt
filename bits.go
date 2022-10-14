package data

type Bits struct {
	b   []byte
	pos int
}

func (b *Bits) AppendBit(set bool) {
	if set {
		b.Append([]byte{1}, 1)
	} else {
		b.Append([]byte{0}, 1)
	}
}

// Append appends n of the bits (up to 8 per byte) to b.
func (b *Bits) Append(bits []byte, n int) {
	for pos := 0; n > 0; n-- {
		if b.pos == 7 {
			b.set(7)
			b.pos = 0
		}
		if b.pos == 0 {
			b.b = append(b.b, 0)
		}

		if bits[0]&(1<<pos) != 0 {
			b.set(b.pos)
		}
		b.pos++

		pos++
		if pos == 8 {
			pos = 0
			bits = bits[1:]
		}
	}
}

func (b *Bits) set(i int) {
	b.b[len(b.b)-1] |= 1 << i
}
