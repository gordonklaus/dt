package types

type ArrayType struct {
	Elem Type // not *EnumType or *StructType
}
