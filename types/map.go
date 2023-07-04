package types

import "github.com/gordonklaus/data/bits"

type MapType struct {
	Key   Type // *IntType, *FloatType, or *StringType
	Value Type // not *EnumType or *StructType
}

func (m *MapType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		WriteType(b, m.Key)
		WriteType(b, m.Value)
	})
}

func (m *MapType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := ReadType(b, &m.Key); err != nil {
			return err
		}
		return ReadType(b, &m.Value)
	})
}
