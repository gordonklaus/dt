package types

import (
	"fmt"

	"github.com/gordonklaus/data/bits"
)

type NamedType struct {
	Package PackageID
	Name    string
}

func (n *NamedType) Write(b *bits.Buffer) {
	b.WriteSize(func() {
		b.WriteVarUint_4bit(0)
		n.Package.Write(b)
		b.WriteString(n.Name)
	})
}

func (n *NamedType) Read(b *bits.Buffer) error {
	return b.ReadSize(func() error {
		pid, err := b.ReadVarUint_4bit()
		if err != nil {
			return err
		}
		switch pid {
		case 0:
			n.Package = &PackageID_Current{}
		default:
			panic(fmt.Sprintf("unknown package ID %d", pid))
		}
		if err := n.Package.Read(b); err != nil {
			return err
		}
		return b.ReadString(&n.Name)
	})
}
