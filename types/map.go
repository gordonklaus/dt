package types

type MapType struct {
	Key   Type // *IntType, *FloatType, or *StringType
	Value Type // not *EnumType or *StructType
}
