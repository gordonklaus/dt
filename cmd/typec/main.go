package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gordonklaus/data/types"
)

func main() {
	out := flag.String("out", ".", "output directory")
	flag.Parse()

	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	loader := types.NewLoader(types.NewStorage(dir))
	pkg, err := loader.Load(&types.PackageID_Current{}) // TODO: Resolve current package ID based on current directory and source control/module configuration.
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(*out, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	w := &writer{}
	w.writePackage(pkg)
	buf := gofmt(gofmt(w.buf.Bytes())) // twice because gofmt isn't quite idempotent
	if err := ioutil.WriteFile(filepath.Join(*out, "pkg.dt.go"), buf, fs.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func gofmt(src []byte) []byte {
	buf, err := format.Source(src)
	if err != nil {
		fmt.Println(string(src))
		log.Fatal(err)
	}
	return buf
}

type writer struct {
	buf bytes.Buffer
}

func (w *writer) writePackage(p *types.Package) {
	w.writeln("package %s", p.Name)
	w.writeln(`import (`)
	w.writeln(`"fmt"`)
	w.writeln(``)
	w.writeln(`"github.com/gordonklaus/data/bits"`)
	w.writeln(`"golang.org/x/exp/maps"`)
	w.writeln(`"golang.org/x/exp/slices"`)
	w.writeln(`)`)
	w.writeln("var (")
	w.writeln("_ = fmt.Print")
	w.writeln("_ = bits.NewBuffer")
	w.writeln("_ = maps.Keys[map[int]int]")
	w.writeln("_ = slices.Sort[int]")
	w.writeln(")")
	for _, n := range p.Types {
		name := camel(n.Name)
		w.writeln("// %s", n.Doc)
		switch t := n.Type.(type) {
		case *types.EnumType:
			w.writeEnum(t, name)
		case *types.StructType:
			w.writeStruct(t, name, false)
		default:
			panic(fmt.Sprintf("unexpected type %T", t))
		}
	}
}

func camel(s string) string { return strings.ReplaceAll(strings.Title(s), " ", "") }

func (w *writer) writeEnum(t *types.EnumType, name string) {
	w.writeln("type %s struct { %s interface { is%s(); bits.ReadWriter } }", name, name, name)

	ename := make([]string, len(t.Elems))
	for i, e := range t.Elems {
		ename[i] = name + "_" + camel(e.Name)
		w.writeln("func (*%s) is%s() {}", ename[i], name)
	}

	w.writeln("\nfunc (x *%s) Write(b *bits.Buffer) {", name)
	w.writeln("switch x.%s.(type) {", name)
	for i := range t.Elems {
		w.writeln("case *%s: b.WriteVarUint_4bit(%d)", ename[i], i)
	}
	w.writeln(`default: panic(fmt.Sprintf("invalid %s enum value %%T", x))`, name)
	w.writeln("}")
	w.writeln("x.%s.Write(b)", name)
	w.writeln("}\n")

	w.writeln("func (x *%s) Read(b *bits.Buffer) error {", name)
	w.writeln("var i uint64")
	w.writeln("if err := b.ReadVarUint_4bit(&i); err != nil {")
	w.writeln("return err")
	w.writeln("}")
	w.writeln("switch i {")
	for i := range t.Elems {
		w.writeln("case %d: x.%s = new(%s)", i, name, ename[i])
	}
	w.writeln("default: x.%s = nil // TODO: &%s__Unknown{i}", name, name)
	w.writeln("}")
	w.writeln("return x.%s.Read(b)", name)
	w.writeln("}\n")

	for i, e := range t.Elems {
		isStruct := false
		if nt, ok := e.Type.(*types.NamedType); ok {
			_, isStruct = nt.Type.(*types.StructType)
		}
		w.writeStruct(&types.StructType{Fields: []*types.StructFieldType{
			{Name: e.Name, Type: e.Type},
		}}, ename[i], isStruct)
	}
}

func (w *writer) writeStruct(t *types.StructType, name string, omitSize bool) {
	fname := make([]string, len(t.Fields))

	w.write("type %s struct {", name)
	for i, f := range t.Fields {
		fname[i] = camel(f.Name)
		w.writeln("\n// %s", f.Doc)
		w.write("%s ", fname[i])
		w.writeType(f.Type)
	}
	w.writeln("}")

	w.writeln("func (x *%s) Write(b *bits.Buffer) {", name)
	if !omitSize {
		w.writeln("b.WriteSize(func() {")
	}
	for i, f := range t.Fields {
		w.writeTypeWriter(f.Type, "x."+fname[i])
	}
	if !omitSize {
		w.writeln("})")
	}
	w.writeln("}\n")

	w.writeln("func (x *%s) Read(b *bits.Buffer) error {", name)
	if !omitSize {
		w.writeln("return b.ReadSize(func() error {")
	}
	for i, f := range t.Fields {
		w.writeTypeReader(f.Type, "&x."+fname[i])
	}
	w.writeln("return nil")
	if !omitSize {
		w.writeln("})")
	}
	w.writeln("}\n")
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
		w.write(camel(t.Name))
	}
}

func (w *writer) writeTypeWriter(t types.Type, v string) {
	switch t := t.(type) {
	case *types.BoolType:
		w.writeln("b.WriteBool(%s)", v)
	case *types.IntType:
		if t.Unsigned {
			w.writeln("b.WriteVarUint(%s)", v)
		} else {
			w.writeln("b.WriteVarInt(%s)", v)
		}
	case *types.FloatType:
		w.writeln("b.WriteFloat%d(%s)", t.Size, v)

	case *types.ArrayType:
		w.writeln("b.WriteVarUint(uint64(len(%s)))", v)
		w.writeln("for _, x := range %s {", v)
		w.writeTypeWriter(t.Elem, "x")
		w.writeln("}")
	case *types.MapType:
		w.writeln("{")
		w.writeln("b.WriteVarUint(uint64(len(%s)))", v)
		w.writeln("keys := maps.Keys(%s)", v)
		w.writeln("slices.Sort(keys)")
		w.writeln("for _, k := range keys {")
		w.writeTypeWriter(t.Key, "k")
		w.writeTypeWriter(t.Value, v+"[k]")
		w.writeln("}")
		w.writeln("}")

	case *types.OptionType:
		w.writeln("b.WriteBool(%s != nil)", v)
		w.writeln("if %s != nil {", v)
		w.writeTypeWriter(t.Elem, "*"+v)
		w.writeln("}")
	case *types.StringType:
		w.writeln("b.WriteString(%s)", v)
	case *types.NamedType:
		w.writeln("(%s).Write(b)", v)
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
		w.writeln("if err := b.ReadVarUint(&len); err != nil { return err }")
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
		w.writeln("if err := b.ReadVarUint(&len); err != nil { return err }")
		w.write("%s = make(", v)
		w.writeType(t)
		w.writeln(", len)")
		w.writeln("for i := len; i > 0; i-- {")
		w.write("var k ")
		w.writeType(t.Key)
		w.writeln("")
		w.writeTypeReader(t.Key, "&k")
		w.write("var v ")
		w.writeType(t.Value)
		w.writeln("")
		w.writeTypeReader(t.Value, "&v")
		w.writeln("%s[k]=v}}", v)
		return

	case *types.OptionType:
		v = indirect(v)
		w.writeln("{var ok bool")
		w.writeln("if err := b.ReadBool(&ok); err != nil { return err }")
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
		w.write("b.ReadBool(%s)", v)
	case *types.IntType:
		if t.Unsigned {
			w.write("b.ReadVarUint(%s)", v)
		} else {
			w.write("b.ReadVarInt(%s)", v)
		}
	case *types.FloatType:
		w.write("b.ReadFloat%d(%s)", t.Size, v)

	case *types.StringType:
		w.write("b.ReadString(%s)", v)
	case *types.NamedType:
		w.write("(%s).Read(b)", v)
	}
	w.writeln("; err != nil { return err }")
}

func (w *writer) writeln(format string, a ...any) {
	w.write(format+"\n", a...)
}

func (w *writer) write(format string, a ...any) {
	fmt.Fprintf(&w.buf, format, a...)
}
