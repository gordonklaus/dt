package bits

import (
	"fmt"
	"io"
	"math"
	"unsafe"
)

type ReadWriter interface {
	Read(*Buffer) error
	Write(*Buffer)
}

type Buffer struct {
	b    []byte
	n, r uint64
}

func NewBuffer() *Buffer {
	b := &Buffer{}
	b.grow(8)
	return b
}

func NewReadBuffer(buf []byte) *Buffer {
	b := &Buffer{
		b: buf[:len(buf):len(buf)],
		n: 8 * uint64(len(buf)),
	}
	b.grow(8)
	return b
}

func (b *Buffer) i() uint64 { return b.n % 8 }
func (b *Buffer) x() *uint64 {
	return (*uint64)(unsafe.Pointer(&b.b[b.n/8]))
}

func (b *Buffer) ri() uint64 { return b.r % 8 }
func (b *Buffer) rx() uint64 {
	return *(*uint64)(unsafe.Pointer(&b.b[b.r/8]))
}

func (b *Buffer) grow(n int) {
	b.b = append(b.b, make([]byte, n)...)
}

func (b *Buffer) WriteBool(x bool) {
	if x {
		*b.x() |= 1 << b.i()
	}
	b.n++
	if b.i() == 0 {
		b.grow(1)
	}
}

func (b *Buffer) ReadBool(x *bool) error {
	if b.Remaining() < 1 {
		return io.ErrUnexpectedEOF
	}
	*x = b.rx()>>b.ri()&1 != 0
	b.r++
	return nil
}

func (b *Buffer) WriteUint32(x uint32) {
	*b.x() |= uint64(x) << b.i()
	b.grow(4)
	b.n += 32
}

func (b *Buffer) ReadUint32(x *uint32) error {
	if b.Remaining() < 32 {
		return io.ErrUnexpectedEOF
	}
	*x = uint32(b.rx() >> b.ri())
	b.r += 32
	return nil
}

func (b *Buffer) WriteUint64(x uint64) {
	*b.x() |= x << b.i()
	b.grow(8)
	*b.x() |= x >> (64 - b.i())
	b.n += 64
}

func (b *Buffer) ReadUint64(x *uint64) error {
	if b.Remaining() < 64 {
		return io.ErrUnexpectedEOF
	}
	*x = b.rx() >> b.ri()
	b.r += 64
	*x |= b.rx() << (64 - b.ri())
	return nil
}

func (b *Buffer) WriteInt64(x int64) { b.WriteUint64(uint64(x)) }
func (b *Buffer) ReadInt64(x *int64) error {
	var i uint64
	err := b.ReadUint64(&i)
	*x = int64(i)
	return err
}

func (b *Buffer) WriteVarUint(x uint64) {
	for {
		y := x & (1<<7 - 1)
		x >>= 7
		if x != 0 {
			y |= 1 << 7
		}
		*b.x() |= y << b.i()
		b.n += 8
		b.grow(1)
		if x == 0 {
			break
		}
	}
}

func (b *Buffer) ReadVarUint(x *uint64) error {
	for shift := 0; shift < 64; shift += 7 {
		if b.Remaining() < 8 {
			return io.ErrUnexpectedEOF
		}

		y := b.rx() >> b.ri()
		b.r += 8
		*x |= y & (1<<7 - 1) << shift
		if y&(1<<7) == 0 {
			return nil
		}
	}
	return fmt.Errorf("VarUint overflows %d bits", 64)
}

func (b *Buffer) WriteVarInt(x int64) { b.WriteVarUint(zigzag(x)) }

func (b *Buffer) ReadVarInt(x *int64) error {
	var u uint64
	err := b.ReadVarUint(&u)
	*x = int64(zagzig(u))
	return err
}

func (b *Buffer) WriteVarUint_4bit(x uint64) {
	for {
		y := x & (1<<3 - 1)
		x >>= 3
		if x != 0 {
			y |= 1 << 3
		}
		*b.x() |= y << b.i()
		b.n += 4
		if b.i() < 4 {
			b.grow(1)
		}
		if x == 0 {
			break
		}
	}
}

func (b *Buffer) ReadVarUint_4bit(x *uint64) error {
	for shift := 0; shift < 64; shift += 3 {
		if b.Remaining() < 4 {
			return io.ErrUnexpectedEOF
		}

		y := b.rx() >> b.ri()
		b.r += 4
		*x |= y & (1<<3 - 1) << shift
		if y&(1<<3) == 0 {
			return nil
		}
	}
	return fmt.Errorf("VarUint_4bit overflows uint64")
}

func zigzag(x int64) uint64 { return uint64((x >> 63) ^ (x << 1)) }
func zagzig(x uint64) int64 { return int64((x >> 1) ^ -(x & 1)) }

func (b *Buffer) WriteFloat32(x float32) { b.WriteUint32(math.Float32bits(x)) }
func (b *Buffer) WriteFloat64(x float64) { b.WriteUint64(math.Float64bits(x)) }

func (b *Buffer) ReadFloat32(x *float32) error {
	var i uint32
	if err := b.ReadUint32(&i); err != nil {
		return err
	}
	*x = math.Float32frombits(i)
	return nil
}

func (b *Buffer) ReadFloat64(x *float64) error {
	var i uint64
	if err := b.ReadUint64(&i); err != nil {
		return err
	}
	*x = math.Float64frombits(i)
	return nil
}

func (b *Buffer) WriteString(s string) {
	b.WriteVarUint(uint64(len(s)))
	b.WriteBytes([]byte(s))
}

func (b *Buffer) ReadString(s *string) error {
	var len uint64
	if err := b.ReadVarUint(&len); err != nil {
		return err
	}
	bb := make([]byte, len)
	if err := b.ReadBytes(bb); err != nil {
		return err
	}
	*s = string(bb)
	return nil
}

func (b *Buffer) WriteBytes(x []byte) {
	b.grow(len(x))
	x = append(x, make([]byte, 7)...)[:len(x)]
	for len(x) > 0 {
		*b.x() |= *(*uint64)(unsafe.Pointer(&x[0])) << b.i()
		if len(x) <= 7 {
			b.n += 8 * uint64(len(x))
			break
		}
		b.n += 56
		x = x[7:]
	}
}

func (b *Buffer) ReadBytes(x []byte) error {
	if b.Remaining() < 8*uint64(len(x)) {
		return io.ErrUnexpectedEOF
	}

	y := make([]byte, len(x)+7)[:len(x)]
	z := y
	for len(z) > 0 {
		*(*uint64)(unsafe.Pointer(&z[0])) |= b.rx() >> b.ri()
		if len(z) <= 7 {
			b.r += 8 * uint64(len(z))
			break
		}
		b.r += 56
		z = z[7:]
	}

	copy(x, y)

	return nil
}

func (b *Buffer) WriteSize(f func()) {
	// WriteSize takes advantage of the fact that the size of the payload (a VarUint) occupies a whole number of bytes to avoid having to bit shift the payload, which could be expensive for large payloads.
	// The fact that the payload is written immediately with the right bit offset is also nice because it makes it possible for nested objects to do byte alignment, which would be good for large byte arrays.

	b2 := NewBuffer()
	b2.n = b.i()
	*b, *b2 = *b2, *b
	f()
	*b, *b2 = *b2, *b
	size := b2.n - b.i()
	b.WriteVarUint(size)
	if b.i() == 0 {
		b.b = append(b.Bytes(), b2.b...)
	} else {
		*b.x() |= uint64(b2.b[0])
		b.b = append(b.Bytes(), b2.b[1:]...)
	}
	b.n += size
}

func (b *Buffer) ReadSize(f func() error) error {
	var size uint64
	if err := b.ReadVarUint(&size); err != nil {
		return err
	}
	if b.r+size > b.n {
		return fmt.Errorf("size overflows available space")
	}

	n := b.n
	b.n = b.r + size
	defer func() { b.n, b.r = n, b.n }()
	return f()
}

func (b *Buffer) Bytes() []byte { return b.b[:(b.n+7)/8] }

func (b *Buffer) Remaining() uint64 { return b.n - b.r }
