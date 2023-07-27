package types

type TypeName struct {
	Name, Doc string
	Type      Type // *EnumType or *StructType
}
