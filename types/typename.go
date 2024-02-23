package types

type TypeName struct {
	ID        uint64
	Name, Doc string
	Type      Type // *EnumType or *StructType
}
