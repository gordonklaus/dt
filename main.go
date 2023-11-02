package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/gordonklaus/dt/internal/typec"
	"github.com/gordonklaus/dt/internal/typed"
	"github.com/gordonklaus/dt/types"
)

func main() {
	flag.CommandLine.Usage = printUsage
	if len(os.Args) < 2 {
		printUsage()
		return
	}
	cmd := os.Args[1]
	os.Args = slices.Delete(os.Args, 1, 2)
	var outdir string
	switch cmd {
	default:
		fmt.Printf("unknown command %q\n", cmd)
		return
	case "edit":
	case "build":
		flag.StringVar(&outdir, "out", ".", "output directory")
	}

	flag.Parse()
	dir := "."
	if flag.NArg() > 0 {
		dir = flag.Arg(0)
	}
	dir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	loader := types.NewLoader(types.NewStorage(dir))
	pkg, err := loader.Load(&types.PackageID_Current{}) // TODO: Resolve current package ID based on current directory and source control/module configuration.
	if cmd == "edit" && errors.Is(err, fs.ErrNotExist) {
		pkg = &types.Package{Name: filepath.Base(dir)}
		loader.Packages[&types.PackageID_Current{}] = pkg
	} else if err != nil {
		fmt.Println(err)
		return
	}

	switch cmd {
	case "edit":
		typed.Run(loader, pkg)
	case "build":
		typec.Run(loader, pkg, outdir)
	}
}

func printUsage() {
	fmt.Println(usage)
}

var usage = `dt is a tool for working with dt data.

Usage:

	dt <command> [arguments]

The commands are:

	edit    edit data
	build   build a package
`
