package data

import "io"

type Bits struct {
	b []byte
	i int
}

func NewBits() *Bits {
	return &Bits{
		b: make([]byte, 1),
	}
}

func (b *Bits) ReadFrom(r io.ByteReader) error {
	for {
		c, err := r.ReadByte()
		if err != nil {
			return err
		}
		b.b = append(b.b, c)
		if c&(1<<7) == 0 {
			return nil
		}
	}
}

func (b *Bits) Append(set bool) {
	if b.i == 7 {
		b.set(7)
		b.b = append(b.b, 0)
	}

	if set {
		b.set(b.i)
	}
	b.i++
}

func (b *Bits) set(i int) {
	b.b[len(b.b)-1] |= 1 << i
}

// AppendN appends n of the bits (up to 8 per byte) to b.
func (b *Bits) AppendN(bits []byte, n int) {
	for i := 0; n > 0; n-- {
		b.Append(bits[0]&(1<<i) != 0)

		i++
		if i == 8 {
			i = 0
			bits = bits[1:]
		}
	}
}

func (b *Bits) Bytes() []byte {
	if len(b.b) == 1 || b.b[len(b.b)-1] != 0 {
		return b.b
	}

	bb := make([]byte, len(b.b)-1)
	copy(bb, b.b)
	for bb[len(bb)-1] == 0 {
		bb = bb[:len(bb)-1]
		bb[len(bb)-1] &^= 1 << 7
	}
	return bb
}

func (b *Bits) Read() (set bool, ok bool) {
	if len(b.b) == 0 {
		return false, false
	}
	set = b.b[0]&(1<<b.i) != 0
	b.i++
	if b.i == 7 {
		b.i = 0
		b.b = b.b[1:]
	}
	return set, true
}
