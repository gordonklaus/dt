package types

type NamedType struct {
	Package PackageID
	Name    string
	Type    Type // *EnumType or *StructType
}
