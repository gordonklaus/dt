package typec

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gordonklaus/dt/types"
)

func BuildC(loader *types.Loader, pkg *types.Package, out string) {
	if err := os.MkdirAll(out, os.ModePerm); err != nil {
		fmt.Println(err)
		return
	}

	w := &cWriter{}
	w.writePackage(pkg)
	if err := os.WriteFile(filepath.Join(out, snake(pkg.Name)+".h"), w.hbuf.Bytes(), fs.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
	if err := os.WriteFile(filepath.Join(out, snake(pkg.Name)+".c"), w.cbuf.Bytes(), fs.ModePerm); err != nil {
		fmt.Println(err)
		return
	}
}

type cWriter struct {
	currentPkg       *types.Package
	hbuf, cbuf       bytes.Buffer
	hindent, cindent int
}

func (w *cWriter) writePackage(p *types.Package) {
	w.currentPkg = p

	w.hln("// Code generated by github.com/gordonklaus/dt.  DO NOT EDIT\n")
	w.hln(`#pragma once"`)
	w.hln("")
	w.hln(`#include "dt/bits/c/codec.h"`)
	w.hln("")

	w.cln("// Code generated by github.com/gordonklaus/dt.  DO NOT EDIT\n")
	w.cln(`#include "%s.h"`, snake(p.Name))
	w.cln("")
	fmt.Fprintln(&w.cbuf, "#define chk(x) { dt_error err = (x); if (err != dt_ok) return err; }")
	w.cln("")

	typs, ok := toposort(p.Types)
	if !ok {
		fmt.Println("cycle detected")
		return
	}

	for _, n := range typs {
		name := word(snake(p.Name) + "__" + snake(n.Name))
		w.hln("typedef struct %s %s;", name, name)
	}
	w.hln("")

	for _, n := range typs {
		name := word(snake(p.Name) + "__" + snake(n.Name))
		if n.Doc != "" {
			w.hln("// %s", n.Doc)
		}
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

func toposort(typs []*types.TypeName) ([]*types.TypeName, bool) {
	sorted := []*types.TypeName{}
	done := map[*types.TypeName]bool{}
	var visit func(n *types.TypeName) bool
	visit = func(n *types.TypeName) bool {
		if done, ok := done[n]; ok {
			return done
		}
		done[n] = false
		switch t := n.Type.(type) {
		case *types.EnumType:
			for _, e := range t.Elems {
				typs = append(typs, &types.TypeName{
					Name: n.Name + "__" + e.Name,
					Doc:  e.Doc,
					Type: e.Type,
				})
			}
		case *types.StructType:
			for _, f := range t.Fields {
				if n, ok := f.Type.(*types.NamedType); ok {
					if _, ok := n.Package.(*types.PackageID_Current); ok && !visit(n.TypeName) {
						return false
					}
				}
			}
		}
		sorted = append(sorted, n)
		done[n] = true
		return true
	}
	for i := 0; i < len(typs); i++ {
		if !visit(typs[i]) {
			return nil, false
		}
	}
	return sorted, true
}

func snake(s string) string { return strings.ReplaceAll(s, " ", "_") }

func word(s string) string {
	switch s {
	case "alignas", "alignof", "auto", "bool", "break", "case", "char", "const", "constexpr", "continue", "default", "do", "double", "else", "enum", "extern", "false", "float", "for", "goto", "if", "inline", "int", "long", "nullptr", "register", "restrict", "return", "short", "signed", "sizeof", "static", "static_assert", "struct", "switch", "thread_local", "true", "typedef", "typeof", "typeof_unqual", "union", "unsigned", "void", "volatile", "while":
		return s + "_"
	}
	return s
}

func (w *cWriter) writeEnum(t *types.EnumType, name string) {
	ename := make([]string, len(t.Elems))
	enameshort := make([]string, len(t.Elems))
	for i, e := range t.Elems {
		ename[i] = word(name + "__" + snake(e.Name))
		enameshort[i] = word(snake(e.Name))
	}

	w.hln("struct {")
	w.hln("enum {")
	for i := range t.Elems {
		w.hln("%s__tag,", ename[i])
	}
	w.hln("} tag;")
	w.hln("union {")
	for i := range t.Elems {
		w.hln("%s *%s;", ename[i], enameshort[i])
	}
	w.hln("} value;")
	w.hln("} %s;", name)
	w.hln("dt_error %s___write(dt_encoder *e, %s *x);", name, name)

	w.cln("dt_error %s___write(dt_encoder *e, %s *x) {", name, name)
	w.cln("if (x->%s == NULL) return dt_error_invalid_enum;", enameshort[0])
	w.cln("chk(dt_write_var_uint_4bit(e, x->tag));")
	w.cln("switch (x->tag) {")
	for i := range t.Elems {
		w.cln("case %s__tag:", ename[i])
		w.cln("return %s___write(e, x->%s);", ename[i], enameshort[i])
	}
	w.cln("}")
	w.cln("return dt_error_invalid_enum;")
	w.cln("}\n")

	w.hln("dt_error %s___read(dt_decoder *d, %s *x);", name, name)
	w.cln("dt_error %s___read(dt_decoder *d, %s *x) {", name, name)
	w.cln("%s___reset(x);", name)
	w.cln("chk(dt_read_var_uint_4bit(d, &x->tag));")
	w.cln("switch (x->tag) {")
	for i := range t.Elems {
		w.cln("case %s__tag:", ename[i])
		w.cln("x->%s = calloc(1, sizeof(*x->%s));", enameshort[i], enameshort[i])
		w.cln("if (x->%s == NULL) return dt_error_out_of_memory;", enameshort[i])
		w.cln("return %s___read(d, x->%s);", ename[i], enameshort[i])
	}
	w.cln("}")
	w.cln("x->tag = dt_unknown_enum_tag;")
	w.cln("return dt_read_size(d, NULL, NULL);")
	w.cln("}\n")

	w.hln("void %s___reset(%s *x);\n", name, name)
	w.cln("void %s___reset(%s *x) {", name, name)
	w.cln("switch (x->tag) {")
	for i := range t.Elems {
		w.cln("case %s__tag:", ename[i])
		w.cln("%s___reset(x->%s);", ename[i], enameshort[i])
		w.cln("free(x->%s);", enameshort[i])
		w.cln("x->%s = NULL;", enameshort[i])
		w.cln("break;")
	}
	w.cln("}")
	w.cln("x->tag = 0;")
	w.cln("}\n")
}

func (w *cWriter) writeStruct(t *types.StructType, name string) {
	fname := make([]string, len(t.Fields))

	w.hln("struct {")
	for i, f := range t.Fields {
		fname[i] = snake(f.Name)
		if f.Doc != "" {
			w.hln("// %s", f.Doc)
		}
		w.hln("%s %s;", w.typ(f.Type), fname[i])
	}
	w.hln("} %s;", name)

	w.cln("static dt_error %s___write_fields(dt_encoder *e, %s *x) {", name, name)
	for i, f := range t.Fields {
		w.writeTypeWriter(f.Type, "x->"+fname[i])
	}
	w.cln("return dt_ok;")
	w.cln("}\n")

	w.hln("dt_error %s___write(dt_encoder *e, %s *x);", name, name)
	w.cln("dt_error %s___write(dt_encoder *e, %s *x) {", name, name)
	w.cln("return dt_write_size(e, %s___write_fields, x);", name)
	w.cln("}\n")

	w.cln("static dt_error %s___read_fields(dt_decoder *d, %s *x) {", name, name)
	for i, f := range t.Fields {
		// TODO: Return if remaining == 0.
		w.writeTypeReader(f.Type, "&x->"+fname[i])
	}
	w.cln("return dt_ok;")
	w.cln("}\n")

	w.hln("dt_error %s___read(dt_decoder *d, %s *x);", name, name)
	w.cln("dt_error %s___read(dt_decoder *d, %s *x) {", name, name)
	w.cln("return dt_read_size(d, %s___read_fields, x);", name)
	w.cln("}\n")

	w.hln("void %s___reset(%s *x);\n", name, name)
	w.cln("void %s___reset(%s *x) {", name, name)
	for i, f := range t.Fields {
		w.writeTypeResetter(f.Type, "&x->"+fname[i])
	}
	w.cln("}\n")
}

func (w *cWriter) typ(t types.Type) string {
	switch t := t.(type) {
	case *types.BoolType:
		return "bool"
	case *types.IntType:
		if t.Unsigned {
			return "uint64_t"
		} else {
			return "int64_t"
		}
	case *types.FloatType:
		if t.Size == 64 {
			return "double"
		}
		return "float"

	case *types.ArrayType:
		return fmt.Sprintf("struct { int len; %s *data; }", w.typ(t.Elem))
	case *types.MapType:
		return fmt.Sprintf("struct { int len; struct { %s key; %s value; } *data; }", w.typ(t.Key), w.typ(t.Value))

	case *types.OptionType:
		return w.typ(t.Elem) + "*"
	case *types.StringType:
		return "dt_string"
	case *types.NamedType:
		pkg := "TODO"
		if _, ok := t.Package.(*types.PackageID_Current); ok {
			pkg = w.currentPkg.Name
		}
		return snake(pkg) + "__" + snake(t.TypeName.Name)
	}
	panic("unreached")
}

func (w *cWriter) writeTypeWriter(t types.Type, v string) {
	switch t := t.(type) {
	case *types.ArrayType:
		w.cln("chk(dt_write_var_uint(e, %s.len));", v)
		w.cln("for (int i = 0; i < %s.len; i++) {", v)
		w.writeTypeWriter(t.Elem, v+".data[i]")
		w.cln("}")
		return
	case *types.MapType:
		w.cln("{")
		w.cln("chk(dt_write_var_uint(e, %s.len));", v)
		w.cln("for (int i = 0; i < %s.len; i++) {", v)
		w.writeTypeWriter(t.Key, v+".data[i].key")
		w.writeTypeWriter(t.Value, v+".data[i].value")
		w.cln("}")
		w.cln("}")
		return

	case *types.OptionType:
		w.cln("chk(dt_write_bool(e, %s != NULL));", v)
		w.cln("if (%s != NULL) {", v)
		w.writeTypeWriter(t.Elem, "*"+v)
		w.cln("}")
		return
	}

	switch t := t.(type) {
	case *types.BoolType:
		w.cln("chk(dt_write_bool(e, %s));", v)
	case *types.IntType:
		if t.Unsigned {
			w.cln("chk(dt_write_var_uint(e, %s));", v)
		} else {
			w.cln("chk(dt_write_var_int(e, %s));", v)
		}
	case *types.FloatType:
		w.cln("chk(dt_write_float%d(e, %s));", t.Size, v)

	case *types.StringType:
		w.cln("chk(dt_write_bytes(e, %s));", v)
	case *types.NamedType:
		w.cln("chk(%s___write(e, %s));", w.typ(t), v)
	}
}

func (w *cWriter) writeTypeReader(t types.Type, v string) {
	indirect := func(v string) string {
		if v[0] == '&' {
			return v[1:]
		}
		return "(*" + v + ")"
	}

	switch t := t.(type) {
	case *types.ArrayType:
		v = indirect(v)
		w.cln("chk(dt_read_var_uint(d, &%s.len));", v)
		w.cln("%s.data = calloc(%s.len, sizeof(*%s.data));", v, v, v)
		w.cln("if (%s.data == NULL && %s.len > 0) return dt_error_out_of_memory;", v, v)
		w.cln("for (int i = 0; i < %s.len; i++) {", v)
		w.writeTypeReader(t.Elem, "&"+v+".data[i]")
		w.cln("}")
		return
	case *types.MapType:
		v = indirect(v)
		w.cln("chk(dt_read_var_uint(d, &%s.len));", v)
		w.cln("%s.data = calloc(%s.len, sizeof(*%s.data));", v, v, v)
		w.cln("if (%s.data == NULL && %s.len > 0) return dt_error_out_of_memory;", v, v)
		w.cln("for (int i = 0; i < %s.len; i++) {", v)
		w.writeTypeReader(t.Key, "&"+v+".data[i].key")
		w.writeTypeReader(t.Value, "&"+v+".data[i].value")
		w.cln("}")
		return

	case *types.OptionType:
		v = indirect(v)
		w.cln("{")
		w.cln("bool ok = false;")
		w.cln("chk(dt_read_bool(d, &ok));")
		w.cln("if (ok) {")
		w.cln("%s = calloc(1, sizeof(*%s));", v, v)
		w.writeTypeReader(t.Elem, v)
		w.cln("}")
		w.cln("}")
		return
	}

	switch t := t.(type) {
	case *types.BoolType:
		w.cln("chk(dt_read_bool(d, %s));", v)
	case *types.IntType:
		if t.Unsigned {
			w.cln("chk(dt_read_var_uint(d, %s));", v)
		} else {
			w.cln("chk(dt_read_var_int(d, %s));", v)
		}
	case *types.FloatType:
		w.cln("chk(dt_read_float%d(d, )%s);", t.Size, v)

	case *types.StringType:
		w.cln("chk(dt_read_bytes(d, %s));", v)
	case *types.NamedType:
		w.cln("chk(%s___read(d, %s));", w.typ(t), v)
	}
}

func (w *cWriter) writeTypeResetter(t types.Type, v string) {
	indirect := func(v string) string {
		if v[0] == '&' {
			return v[1:]
		}
		return "(*" + v + ")"
	}

	switch t := t.(type) {
	case *types.ArrayType:
		v = indirect(v)
		w.cln("for (int i = 0; i < %s.len; i++) {", v)
		w.writeTypeResetter(t.Elem, v+".data[i]")
		w.cln("}")
		w.cln("%s.len = 0;", v)
		w.cln("free(%s.data);", v)
		w.cln("%s.data = NULL;", v)
	case *types.MapType:
		v = indirect(v)
		w.cln("for (int i = 0; i < %s.len; i++) {", v)
		w.writeTypeResetter(t.Key, v+".data[i].key")
		w.writeTypeResetter(t.Value, v+".data[i].value")
		w.cln("}")
		w.cln("%s.len = 0;", v)
		w.cln("free(%s.data);", v)
		w.cln("%s.data = NULL;", v)

	case *types.OptionType:
		w.cln("if (%s != NULL) {", v)
		w.writeTypeResetter(t.Elem, "*"+v)
		w.cln("free(%s);", v)
		w.cln("%s = NULL;", v)
		w.cln("}")
	case *types.StringType:
		w.cln("dt_bytes_delete(%s);", v)
	case *types.NamedType:
		w.cln("%s___reset(%s);", w.typ(t), v)
	}
}

func (w *cWriter) hln(format string, a ...any) {
	w.h(format+"\n", a...)
}

func (w *cWriter) h(format string, a ...any) {
	if strings.ContainsRune(format, '}') {
		w.hindent--
	}
	fmt.Fprint(&w.hbuf, strings.Repeat("\t", w.hindent))
	fmt.Fprintf(&w.hbuf, format, a...)
	if strings.ContainsRune(format, '{') {
		w.hindent++
	}
}

func (w *cWriter) cln(format string, a ...any) {
	w.c(format+"\n", a...)
}

func (w *cWriter) c(format string, a ...any) {
	if strings.ContainsRune(format, '}') {
		w.cindent--
	}
	indent := w.cindent
	if strings.HasPrefix(format, "case") || strings.HasPrefix(format, "default") {
		indent--
	}
	fmt.Fprint(&w.cbuf, strings.Repeat("\t", indent))
	fmt.Fprintf(&w.cbuf, format, a...)
	if strings.ContainsRune(format, '{') {
		w.cindent++
	}
}
