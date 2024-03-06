package types

import "fmt"

func ValidatePackage(p *Package) error {
	if err := ValidateName(p.Name, nil); err != nil {
		return err
	}
	ids := map[uint64]bool{}
	names := map[string]bool{}
	for _, t := range p.Types {
		if err := ValidateID(t.ID, ids); err != nil {
			return err
		}
		if err := ValidateName(t.Name, names); err != nil {
			return err
		}
		switch t.Type.(type) {
		case *EnumType, *StructType:
		default:
			return fmt.Errorf("invalid named type %T", t.Type)
		}
		if err := ValidateType(t.Type, ids); err != nil {
			return err
		}
	}
	return nil
}

func ValidateType(t Type, ids map[uint64]bool) error {
	switch t := t.(type) {
	case *FloatType:
		if t.Size != 32 && t.Size != 64 {
			return fmt.Errorf("invalid float size %d", t.Size)
		}
	case *OptionType:
		switch t.Elem.(type) {
		case *EnumType, *StructType:
			return fmt.Errorf("invalid option element %T", t.Elem)
		}
		return ValidateType(t.Elem, ids)
	case *ArrayType:
		switch t.Elem.(type) {
		case *EnumType, *StructType:
			return fmt.Errorf("invalid array element %T", t.Elem)
		}
		return ValidateType(t.Elem, ids)
	case *MapType:
		switch t.Key.(type) {
		case *IntType, *FloatType, *StringType:
		default:
			return fmt.Errorf("invalid map key %T", t.Key)
		}
		if err := ValidateType(t.Key, ids); err != nil {
			return err
		}
		switch t.Value.(type) {
		case *EnumType, *StructType:
			return fmt.Errorf("invalid map value %T", t.Value)
		}
		return ValidateType(t.Value, ids)
	case *EnumType:
		names := map[string]bool{}
		for _, e := range t.Elems {
			if err := ValidateID(e.ID, ids); err != nil {
				return err
			}
			if err := ValidateName(e.Name, names); err != nil {
				return err
			}
			switch e.Type.(type) {
			case *StructType:
			default:
				return fmt.Errorf("invalid enum element %T", e.Type)
			}
			if err := ValidateType(e.Type, ids); err != nil {
				return err
			}
		}
	case *StructType:
		names := map[string]bool{}
		for _, f := range t.Fields {
			if err := ValidateID(f.ID, ids); err != nil {
				return err
			}
			if err := ValidateName(f.Name, names); err != nil {
				return err
			}
			switch f.Type.(type) {
			case *EnumType, *StructType:
				return fmt.Errorf("invalid struct field %T", f.Type)
			}
			if err := ValidateType(f.Type, ids); err != nil {
				return err
			}
		}
	}
	return nil
}

func ValidateID(id uint64, ids map[uint64]bool) error {
	if ids[id] {
		return fmt.Errorf("duplicate id %d", id)
	}
	ids[id] = true
	return nil
}

func ValidateName(n string, names map[string]bool) error {
	if n == "" {
		return fmt.Errorf("invalid name %s", n)
	}
	if names != nil {
		if names[n] {
			return fmt.Errorf("duplicate name %s", n)
		}
		names[n] = true
	}
	return nil
}
