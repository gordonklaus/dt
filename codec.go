package data

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Encoder struct {
	w   io.Writer
	err error
}

type Decoder struct {
	r   *bufio.Reader
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: bufio.NewReader(r)}
}

func (e *Encoder) encodeValue(v Value) bool {
	e.err = e.EncodeValue(v)
	return e.err == nil
}

func (d *Decoder) decodeValue(v Value) bool {
	d.err = d.DecodeValue(v)
	return d.err == nil
}

func (e *Encoder) writeString(s string) bool {
	if !e.writeBinary(int64(len(s))) {
		return false
	}
	_, e.err = e.w.Write([]byte(s))
	return e.err == nil
}

func (d *Decoder) readString(s *string) bool {
	var len int64
	if !d.readBinary(&len) {
		return false
	}
	b := make([]byte, len)
	_, d.err = d.r.Read(b)
	if d.err == nil {
		*s = string(b)
	}
	return d.err == nil
}

func (e *Encoder) writeBinary(x any) bool {
	var buf [binary.MaxVarintLen64]byte
	switch x := x.(type) {
	case int:
		n := binary.PutVarint(buf[:], int64(x))
		e.writeBinary(buf[:n])
	case uint:
		n := binary.PutUvarint(buf[:], uint64(x))
		e.writeBinary(buf[:n])
	default:
		e.err = binary.Write(e.w, binary.LittleEndian, x)
	}
	return e.err == nil
}

func (d *Decoder) readBinary(x any) bool {
	switch x := x.(type) {
	case *int:
		i, err := binary.ReadVarint(d.r)
		*x = int(i)
		d.err = err
	case *uint:
		i, err := binary.ReadUvarint(d.r)
		*x = uint(i)
		d.err = err
	default:
		d.err = binary.Read(d.r, binary.LittleEndian, x)
	}
	return d.err == nil
}

func writeSlice[T any](e *Encoder, encode func(T) error, s []T) bool {
	if !e.writeBinary(int64(len(s))) {
		return false
	}
	for _, x := range s {
		if e.err = encode(x); e.err != nil {
			return false
		}
	}
	return true
}

func readSlice[T any](d *Decoder, decode func(*T) error, s *[]*T) bool {
	var len int64
	if !d.readBinary(&len) {
		return false
	}
	*s = make([]*T, len)
	for i := range *s {
		x := new(T)
		(*s)[i] = x
		if d.err = decode(x); d.err != nil {
			return false
		}
	}
	return true
}
