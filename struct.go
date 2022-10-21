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
	bits := NewBits()
	var buf bytes.Buffer
	e2 := NewEncoder(&buf)
	for _, f := range s.Fields {
		switch v := f.Value.(type) {
		case *BasicValue:
			if x, ok := v.X.(*bool); ok {
				bits.Append(*x)
				continue
			}
		case *OptionValue:
			bits.Append(v.Value != nil)
			if v, ok := v.Value.(*BasicValue); ok {
				if x, ok := v.X.(*bool); ok {
					bits.Append(*x)
					continue
				}
			}
			if !e2.encodeValue(v.Value) {
				return e2.err
			}
			continue
		}
		if !e2.encodeValue(f.Value) {
			return e2.err
		}
	}
	len := uint(buf.Len())
	_ = e.writeBinary(bits.Bytes()) && e.writeBinary(&len) && e.writeBinary(buf.Bytes())
	return e.err
}

func (d *Decoder) DecodeStructValue(s *StructValue) error {
	var bits Bits
	if d.err = bits.ReadFrom(d.r); d.err != nil {
		return d.err
	}
	var len uint
	if !d.readBinary(&len) {
		return d.err
	}
	lr := &io.LimitedReader{R: d.r, N: int64(len)}
	d2 := NewDecoder(lr)
fields:
	for _, f := range s.Fields {
		switch v := f.Value.(type) {
		case *BasicValue:
			if b, ok := v.X.(*bool); ok {
				if !d.DecodeBoolValueFromBits(b, &bits) {
					break fields
				}
				continue
			}
		case *OptionValue:
			if err := d2.DecodeOptionValueFromBits(v, &bits); err == io.EOF {
				// TODO: Translate EOF to UnexpectedEOF in decoder methods and distinguish them here.
				break fields
			} else if err != nil {
				return err
			}
			continue
		}
		if lr.N == 0 {
			// Ignore EOF as final fields may be omitted.
			break
		}
		if !d2.decodeValue(f.Value) {
			return d2.err
		}
	}
	if lr.N > 0 {
		_, d.err = d.r.Read(make([]byte, lr.N))
	}
	return d.err
}

func (d *Decoder) DecodeBoolValueFromBits(b *bool, bits *Bits) bool {
	set, ok := bits.Read()
	if !ok {
		return false
	}
	*b = set
	return true
}

func (d *Decoder) DecodeOptionValueFromBits(o *OptionValue, bits *Bits) error {
	if set, ok := bits.Read(); !ok {
		return io.EOF
	} else if !set {
		return nil
	}

	o.Value = NewValue(o.ValueType)
	switch v := o.Value.(type) {
	case *BasicValue:
		if b, ok := v.X.(*bool); ok {
			if !d.DecodeBoolValueFromBits(b, bits) {
				return io.EOF
			}
			return nil
		}
	case *OptionValue:
		if d.err = d.DecodeOptionValueFromBits(v, bits); d.err != nil {
			return d.err
		}
		return nil
	}
	if !d.decodeValue(o.Value) {
		return d.err
	}
	return nil
}
