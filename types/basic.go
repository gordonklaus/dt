package types

type BasicType struct {
	Kind Kind
}

func NewBasicType(kind Kind) *BasicType {
	return &BasicType{Kind: kind}
}
