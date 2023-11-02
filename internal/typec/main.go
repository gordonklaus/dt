package typec

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gordonklaus/dt/types"
)

func Run(loader *types.Loader, pkg *types.Package, out string) {
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	w := &writer{}
	w.writePackage(pkg)
	buf := gofmt(gofmt(w.buf.Bytes())) // twice because gofmt isn't quite idempotent
	if err := os.WriteFile(filepath.Join(out, "pkg.dt.go"), buf, fs.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
}

func gofmt(src []byte) []byte {
	buf, err := format.Source(src)
	if err != nil {
		fmt.Println(string(src))
		fmt.Println(err)
		os.Exit(1)
	}
	return buf
}

type writer struct {
	buf   bytes.Buffer
	varID int
}

func (w *writer) writePackage(p *types.Package) {
	w.writeln("package %s", p.Name)
	w.writeln(`import (`)
	w.writeln(`"fmt"`)
	w.writeln(`"slices"`)
	w.writeln(``)
	w.writeln(`"github.com/gordonklaus/dt/bits"`)
	w.writeln(`"golang.org/x/exp/maps"`)
	w.writeln(`)`)
	w.writeln("var (")
	w.writeln("_ = fmt.Print")
	w.writeln("_ = bits.NewEncoder")
	w.writeln("_ = maps.Keys[map[int]int]")
	w.writeln("_ = slices.Sort[[]int]")
	w.writeln(")")
	for _, n := range p.Types {
		name := camel(n.Name)
		w.writeln("// %s", n.Doc)
		switch t := n.Type.(type) {
		case *types.EnumType:
			w.writeEnum(t, name)
		case *types.StructType:
			w.writeStruct(t, name)
		default:
			panic(fmt.Sprintf("unexpected type %T", t))
		}
	}
}

func camel(s string) string { return strings.ReplaceAll(strings.Title(s), " ", "") }

func (w *writer) writeEnum(t *types.EnumType, name string) {
	w.writeln("type %s struct { %s %s__Enum }", name, name, name)
	w.writeln("type %s__Enum interface { is%s(); bits.Value }", name, name)

	ename := make([]string, len(t.Elems))
	for i, e := range t.Elems {
		ename[i] = name + "_" + camel(e.Name)
		w.writeln("func (*%s) is%s() {}", ename[i], name)
	}

	w.writeln("\nfunc (x *%s) Write(e *bits.Encoder) {", name)
	w.writeln("switch x.%s.(type) {", name)
	for i := range t.Elems {
		w.writeln("case *%s: e.WriteVarUint_4bit(%d)", ename[i], i)
	}
	w.writeln(`default: panic(fmt.Sprintf("invalid %s enum value %%T", x))`, name)
	w.writeln("}")
	w.writeln("x.%s.Write(e)", name)
	w.writeln("}\n")

	w.writeln("func (x *%s) Read(d *bits.Decoder) error {", name)
	w.writeln("var i uint64")
	w.writeln("if err := d.ReadVarUint_4bit(&i); err != nil {")
	w.writeln("return err")
	w.writeln("}")
	w.writeln("switch i {")
	for i := range t.Elems {
		w.writeln("case %d: x.%s = new(%s)", i, name, ename[i])
	}
	w.writeln("default: x.%s = nil // TODO: &%s__Unknown{i}", name, name)
	w.writeln("}")
	w.writeln("return x.%s.Read(d)", name)
	w.writeln("}\n")

	for i, e := range t.Elems {
		w.writeStruct(e.Type.(*types.StructType), ename[i])
	}
}

func (w *writer) writeStruct(t *types.StructType, name string) {
	fname := make([]string, len(t.Fields))

	w.write("type %s struct {", name)
	for i, f := range t.Fields {
		fname[i] = camel(f.Name)
		w.writeln("\n// %s", f.Doc)
		w.write("%s ", fname[i])
		w.writeType(f.Type)
	}
	w.writeln("}")

	w.writeln("func (x *%s) Write(e *bits.Encoder) {", name)
	w.writeln("e.WriteSize(func() {")
	for i, f := range t.Fields {
		w.writeTypeWriter(f.Type, "x."+fname[i])
	}
	w.writeln("})}\n")

	w.varID = 0
	w.writeln("func (x *%s) Read(d *bits.Decoder) error {", name)
	w.writeln("return d.ReadSize(func() error {")
	for i, f := range t.Fields {
		w.writeTypeReader(f.Type, "&x."+fname[i])
	}
	w.writeln("return nil})}\n")
}

func (w *writer) writeType(t types.Type) {
	switch t := t.(type) {
	case *types.BoolType:
		w.write("bool")
	case *types.IntType:
		if t.Unsigned {
			w.write("uint64")
		} else {
			w.write("int64")
		}
	case *types.FloatType:
		w.write("float%d", t.Size)

	case *types.ArrayType:
		w.write("[]")
		w.writeType(t.Elem)
	case *types.MapType:
		w.write("map[")
		w.writeType(t.Key)
		w.write("]")
		w.writeType(t.Value)

	case *types.OptionType:
		w.write("*")
		w.writeType(t.Elem)
	case *types.StringType:
		w.write("string")
	case *types.NamedType:
		w.write(camel(t.TypeName.Name))
	}
}

func (w *writer) writeTypeWriter(t types.Type, v string) {
	switch t := t.(type) {
	case *types.BoolType:
		w.writeln("e.WriteBool(%s)", v)
	case *types.IntType:
		if t.Unsigned {
			w.writeln("e.WriteVarUint(%s)", v)
		} else {
			w.writeln("e.WriteVarInt(%s)", v)
		}
	case *types.FloatType:
		w.writeln("e.WriteFloat%d(%s)", t.Size, v)

	case *types.ArrayType:
		w.writeln("e.WriteVarUint(uint64(len(%s)))", v)
		w.writeln("for _, x := range %s {", v)
		w.writeTypeWriter(t.Elem, "x")
		w.writeln("}")
	case *types.MapType:
		w.writeln("{")
		w.writeln("e.WriteVarUint(uint64(len(%s)))", v)
		w.writeln("keys := maps.Keys(%s)", v)
		w.writeln("slices.Sort(keys)")
		w.writeln("for _, k := range keys {")
		w.writeTypeWriter(t.Key, "k")
		w.writeTypeWriter(t.Value, v+"[k]")
		w.writeln("}")
		w.writeln("}")

	case *types.OptionType:
		w.writeln("e.WriteBool(%s != nil)", v)
		w.writeln("if %s != nil {", v)
		w.writeTypeWriter(t.Elem, "*"+v)
		w.writeln("}")
	case *types.StringType:
		w.writeln("e.WriteString(%s)", v)
	case *types.NamedType:
		w.writeln("(%s).Write(e)", v)
	}
}

func (w *writer) writeTypeReader(t types.Type, v string) {
	indirect := func(v string) string {
		if v[0] == '&' {
			return v[1:]
		}
		return "*" + v
	}

	switch t := t.(type) {
	case *types.ArrayType:
		v = indirect(v)
		w.writeln("{var len uint64")
		w.writeln("if err := d.ReadVarUint(&len); err != nil { return err }")
		w.write("%s = make([]", v)
		w.writeType(t.Elem)
		w.writeln(", len)")
		w.writeln("for i := range %s {", v)
		w.writeTypeReader(t.Elem, "&("+v+")[i]")
		w.writeln("}}")
		return
	case *types.MapType:
		v = indirect(v)
		w.writeln("{var len uint64")
		w.writeln("if err := d.ReadVarUint(&len); err != nil { return err }")
		w.write("%s = make(", v)
		w.writeType(t)
		w.writeln(", len)")
		w.writeln("for i := len; i > 0; i-- {")
		w.write("var k ")
		w.writeType(t.Key)
		w.writeln("")
		w.writeTypeReader(t.Key, "&k")
		varID := w.varID
		w.varID++
		w.write("var v%d ", varID)
		w.writeType(t.Value)
		w.writeln("")
		w.writeTypeReader(t.Value, fmt.Sprintf("&v%d", varID))
		w.writeln("%s[k]=v%d}}", v, varID)
		return

	case *types.OptionType:
		v = indirect(v)
		w.writeln("{var ok bool")
		w.writeln("if err := d.ReadBool(&ok); err != nil { return err }")
		w.writeln("if ok {")
		w.write("%s = new(", v)
		w.writeType(t.Elem)
		w.writeln(")")
		w.writeTypeReader(t.Elem, v)
		w.writeln("}}")
		return
	}

	w.write("if err := ")
	switch t := t.(type) {
	case *types.BoolType:
		w.write("d.ReadBool(%s)", v)
	case *types.IntType:
		if t.Unsigned {
			w.write("d.ReadVarUint(%s)", v)
		} else {
			w.write("d.ReadVarInt(%s)", v)
		}
	case *types.FloatType:
		w.write("d.ReadFloat%d(%s)", t.Size, v)

	case *types.StringType:
		w.write("d.ReadString(%s)", v)
	case *types.NamedType:
		w.write("(%s).Read(d)", v)
	}
	w.writeln("; err != nil { return err }")
}

func (w *writer) writeln(format string, a ...any) {
	w.write(format+"\n", a...)
}

func (w *writer) write(format string, a ...any) {
	fmt.Fprintf(&w.buf, format, a...)
}
