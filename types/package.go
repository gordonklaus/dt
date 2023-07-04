package types

import (
	"github.com/gordonklaus/data/bits"
)

type Package struct {
	Name, Doc string
	Types     []*TypeName
}

func (p *Package) Type(name string) *TypeName {
	for _, t := range p.Types {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (p *Package) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteString(p.Name)
		b.WriteString(p.Doc)
		b.WriteVarUint(uint64(len(p.Types)))
		for _, t := range p.Types {
			t.Write(b)
		}
	})
}

func (p *Package) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		if err := b.ReadString(&p.Name); err != nil {
			return err
		}
		if err := b.ReadString(&p.Doc); err != nil {
			return err
		}
		var len uint64
		if err := b.ReadVarUint(&len); err != nil {
			return err
		}
		p.Types = make([]*TypeName, len)
		for i := range p.Types {
			p.Types[i] = &TypeName{}
			if err := p.Types[i].Read(b); err != nil {
				return err
			}
		}
		return nil
	})
}

type PackageID interface {
	Write(*bits.Buffer)
	Read(*bits.Buffer) error
}

type PackageID_Current struct{}

func (*PackageID_Current) Write(b *bits.Buffer) {
	b.WriteSize(func() {})
}

func (*PackageID_Current) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error { return nil })
}
