package data

import (
	"bytes"
	"io"

	"github.com/gordonklaus/data/types"
)

type StructValue struct {
	Type   *types.StructType
	Fields []*StructFieldValue
}

type StructFieldValue struct {
	*types.StructFieldType
	Value Value
}

func NewStructValue(t *types.StructType) *StructValue {
	s := &StructValue{
		Type:   t,
		Fields: make([]*StructFieldValue, len(t.Fields)),
	}
	for i, f := range t.Fields {
		s.Fields[i] = &StructFieldValue{
			StructFieldType: f,
			Value:           NewValue(f.Type),
		}
	}
	return s
}

func (e *Encoder) EncodeStructValue(s *StructValue) error {
	var bits Bits
	var buf bytes.Buffer
	e2 := NewEncoder(&buf)
	for _, f := range s.Fields {
		switch v := f.Value.(type) {
		case *BasicValue:
			if x, ok := v.X.(bool); ok {
				bits.AppendBit(x)
				continue
			}
		case *OptionValue:
			bits.AppendBit(v.Value != nil)
			if v, ok := v.Value.(*BasicValue); ok {
				if x, ok := v.X.(bool); ok {
					bits.AppendBit(x)
					continue
				}
			}
			if !e2.encodeValue(v) {
				return e2.err
			}
			continue
		}
		if !e2.encodeValue(f.Value) {
			return e2.err
		}
	}
	_ = e.writeBinary(bits.b) && e.writeBinary(uint(buf.Len())) && e.writeBinary(buf.Bytes())
	return e.err
}

func (d *Decoder) DecodeStructValue(s *StructValue) error {
	var fields, len uint
	if !d.readBinary(&fields) || !d.readBinary(&len) {
		return d.err
	}
	lr := &io.LimitedReader{R: d.r, N: int64(len)}
	d2 := NewDecoder(lr)
	for i, f := range s.Fields {
		if fields&(1<<i) != 0 {
			if !d2.decodeValue(f.Value) {
				return d2.err
			}
		}
	}
	if lr.N > 0 {
		_, d.err = d.r.Read(make([]byte, lr.N))
	}
	return d.err
}
