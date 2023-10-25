package bits

import (
	"fmt"
	"io"
	"math"
	"unsafe"
)

type Value interface {
	Write(*Encoder)
	Read(*Decoder) error
}

func Write(w io.Writer, v Value) error {
	e := NewEncoder()
	v.Write(e)
	_, err := w.Write(e.Bytes())
	return err
}

func Read(r io.Reader, v Value) error {
	return v.Read(NewDecoder(r))
}

type Encoder struct {
	b []byte
	n uint64
}

type Decoder struct {
	r    io.Reader
	b    [9]byte
	j, n uint64
}

func NewEncoder() *Encoder {
	return &Encoder{
		b: make([]byte, 8),
	}
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
		n: math.MaxUint64,
	}
}

func (e *Encoder) i() uint64 { return e.n % 8 }
func (e *Encoder) x() *uint64 {
	return (*uint64)(unsafe.Pointer(&e.b[e.n/8]))
}

func (d *Decoder) i() uint64 { return d.j % 8 }
func (d *Decoder) x() uint64 {
	return *(*uint64)(unsafe.Pointer(&d.b[0]))
}

func (e *Encoder) grow(n int) {
	e.b = append(e.b, make([]byte, n)...)
}

func (d *Decoder) read(n uint64) error {
	start := (d.i() + 7) / 8
	_, err := d.r.Read(d.b[start : start+n])
	return err
}

func (e *Encoder) WriteBool(x bool) {
	if x {
		*e.x() |= 1 << e.i()
	}
	e.n++
	if e.i() == 0 {
		e.grow(1)
	}
}

func (d *Decoder) ReadBool(x *bool) error {
	if d.Remaining() < 1 {
		return io.ErrUnexpectedEOF
	}
	if d.i() == 0 {
		if err := d.read(1); err != nil {
			return err
		}
	}
	*x = d.x()>>d.i()&1 != 0
	d.j++
	return nil
}

func (e *Encoder) WriteUint32(x uint32) {
	*e.x() |= uint64(x) << e.i()
	e.n += 32
	e.grow(4)
}

func (d *Decoder) ReadUint32(x *uint32) error {
	if d.Remaining() < 32 {
		return io.ErrUnexpectedEOF
	}
	if err := d.read(4); err != nil {
		return err
	}
	*x = uint32(d.x() >> d.i())
	d.j += 32
	d.b[0] = d.b[4]
	return nil
}

func (e *Encoder) WriteUint64(x uint64) {
	*e.x() |= x << e.i()
	e.n += 64
	e.grow(8)
	*e.x() |= x >> (64 - e.i())
}

func (d *Decoder) ReadUint64(x *uint64) error {
	if d.Remaining() < 64 {
		return io.ErrUnexpectedEOF
	}
	if err := d.read(8); err != nil {
		return err
	}
	*x = d.x() >> d.i()
	d.j += 64
	d.b[0] = d.b[8]
	*x |= d.x() << (64 - d.i())
	return nil
}

func (e *Encoder) WriteInt64(x int64) { e.WriteUint64(uint64(x)) }
func (d *Decoder) ReadInt64(x *int64) error {
	var i uint64
	err := d.ReadUint64(&i)
	*x = int64(i)
	return err
}

func (e *Encoder) WriteVarUint(x uint64) {
	for {
		y := x & (1<<7 - 1)
		x >>= 7
		if x != 0 {
			y |= 1 << 7
		}
		*e.x() |= y << e.i()
		e.n += 8
		e.grow(1)
		if x == 0 {
			break
		}
	}
}

func (d *Decoder) ReadVarUint(x *uint64) error {
	for shift := 0; shift < 64; shift += 7 {
		if d.Remaining() < 8 {
			return io.ErrUnexpectedEOF
		}
		if err := d.read(1); err != nil {
			return err
		}

		y := d.x() >> d.i()
		d.j += 8
		d.b[0] = d.b[1]
		*x |= y & (1<<7 - 1) << shift
		if y&(1<<7) == 0 {
			return nil
		}
	}
	return fmt.Errorf("VarUint overflows %d bits", 64)
}

func (e *Encoder) WriteVarInt(x int64) { e.WriteVarUint(zigzag(x)) }

func (d *Decoder) ReadVarInt(x *int64) error {
	var u uint64
	err := d.ReadVarUint(&u)
	*x = int64(zagzig(u))
	return err
}

func (e *Encoder) WriteVarUint_4bit(x uint64) {
	for {
		y := x & (1<<3 - 1)
		x >>= 3
		if x != 0 {
			y |= 1 << 3
		}
		*e.x() |= y << e.i()
		e.n += 4
		if e.i() < 4 {
			e.grow(1)
		}
		if x == 0 {
			break
		}
	}
}

func (d *Decoder) ReadVarUint_4bit(x *uint64) error {
	for shift := 0; shift < 64; shift += 3 {
		if d.Remaining() < 4 {
			return io.ErrUnexpectedEOF
		}
		if (d.j+3)%8 < 4 {
			if err := d.read(1); err != nil {
				return err
			}
		}

		y := d.x() >> d.i()
		d.j += 4
		if d.i() < 4 {
			d.b[0] = d.b[1]
		}
		*x |= y & (1<<3 - 1) << shift
		if y&(1<<3) == 0 {
			return nil
		}
	}
	return fmt.Errorf("VarUint_4bit overflows uint64")
}

func zigzag(x int64) uint64 { return uint64((x >> 63) ^ (x << 1)) }
func zagzig(x uint64) int64 { return int64((x >> 1) ^ -(x & 1)) }

func (e *Encoder) WriteFloat32(x float32) { e.WriteUint32(math.Float32bits(x)) }
func (e *Encoder) WriteFloat64(x float64) { e.WriteUint64(math.Float64bits(x)) }

func (d *Decoder) ReadFloat32(x *float32) error {
	var i uint32
	if err := d.ReadUint32(&i); err != nil {
		return err
	}
	*x = math.Float32frombits(i)
	return nil
}

func (d *Decoder) ReadFloat64(x *float64) error {
	var i uint64
	if err := d.ReadUint64(&i); err != nil {
		return err
	}
	*x = math.Float64frombits(i)
	return nil
}

func (e *Encoder) WriteString(s string) {
	e.WriteVarUint(uint64(len(s)))
	e.WriteBytes([]byte(s))
}

func (d *Decoder) ReadString(s *string) error {
	var len uint64
	if err := d.ReadVarUint(&len); err != nil {
		return err
	}
	bb := make([]byte, len)
	if err := d.ReadBytes(bb); err != nil {
		return err
	}
	*s = string(bb)
	return nil
}

func (e *Encoder) WriteBytes(x []byte) {
	e.grow(len(x))
	x = append(x, make([]byte, 7)...)[:len(x)]
	for len(x) > 0 {
		*e.x() |= *(*uint64)(unsafe.Pointer(&x[0])) << e.i()
		if len(x) <= 7 {
			e.n += 8 * uint64(len(x))
			break
		}
		e.n += 56
		x = x[7:]
	}
}

func (d *Decoder) ReadBytes(x []byte) error {
	if d.Remaining() < 8*uint64(len(x)) {
		return io.ErrUnexpectedEOF
	}

	y := make([]byte, len(x)+7)[:len(x)]
	z := y
	for len(z) > 0 {
		n := 7
		if len(z) < 7 {
			n = len(z)
		}
		if err := d.read(uint64(n)); err != nil {
			return err
		}

		*(*uint64)(unsafe.Pointer(&z[0])) |= d.x() >> d.i() & (1<<56 - 1)
		if len(z) <= 7 {
			d.j += 8 * uint64(len(z))
			d.b[0] = d.b[len(z)]
			break
		}
		d.j += 56
		d.b[0] = d.b[7]
		z = z[7:]
	}

	copy(x, y)

	return nil
}

func (e *Encoder) WriteSize(f func()) {
	// WriteSize takes advantage of the fact that the size of the payload (a VarUint) occupies a whole number of bytes to avoid having to bit shift the payload, which could be expensive for large payloads.
	// The fact that the payload is written immediately with the right bit offset is also nice because it makes it possible for nested objects to do byte alignment, which would be good for large byte arrays.

	e2 := NewEncoder()
	e2.n = e.i()
	*e, *e2 = *e2, *e
	f()
	*e, *e2 = *e2, *e
	size := e2.n - e.i()
	e.WriteVarUint(size)
	if e.i() == 0 {
		e.b = append(e.Bytes(), e2.b...)
	} else {
		*e.x() |= uint64(e2.b[0])
		e.b = append(e.Bytes(), e2.b[1:]...)
	}
	e.n += size
}

func (d *Decoder) ReadSize(f func() error) error {
	var size uint64
	if err := d.ReadVarUint(&size); err != nil {
		return err
	}
	if d.Remaining() < size {
		return fmt.Errorf("size exceeds limit")
	}

	n := d.n
	d.n = d.j + size

	if err := f(); err != nil {
		return err
	}

	if k := (d.n-1)/8 - (d.j-1)/8; k > 0 {
		buf := make([]byte, k)
		if _, err := d.r.Read(buf); err != nil {
			return err
		}
		d.b[0] = buf[k-1]
	}
	d.j = d.n
	d.n = n

	return nil
}

func (e *Encoder) Bytes() []byte { return e.b[:(e.n+7)/8] }
func (e *Encoder) Size() uint64  { return e.n }

func (d *Decoder) SetLimit(n uint64) { d.n = n }
func (d *Decoder) Remaining() uint64 { return d.n - d.j }
