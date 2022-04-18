package data

type EnumType struct {
	Elems []EnumElem
}

type EnumElem struct {
	Name string
	Type Type
}

type EnumValue struct {
	Type  *EnumType
	Elem  int
	Value Value
}

// func (e *Enum) WriteTo(w io.Writer) (int64, error) {
// 	var c counterr
// 	if !c.write(w, e.Elem) {
// 		return c.error
// 	}
// 	if !c.write(w, e.Value) {
// 		return c.error
// 	}
// 	return c.error
// }

// func (e *Enum) ReadFrom(r io.Reader) (int64, error) {
// 	var c counterr
// 	if !c.read(r, e.Elem) {
// 		return c.error
// 	}
// 	e.Value = e.Type.Elems[e.Elem].Type.NewValue()
// 	if !c.read(r, e.Value) {
// 		return c.error
// 	}
// 	return c.error
// }
