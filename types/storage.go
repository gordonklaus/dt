package types

import (
	"fmt"
	"os"
	"path/filepath"
)

type Storage struct {
	workingDirectory string
}

func NewStorage(workingDirectory string) *Storage {
	return &Storage{
		workingDirectory: workingDirectory,
	}
}

func (s *Storage) Load(id PackageID) ([]byte, error) {
	switch id.(type) {
	case *PackageID_Current:
		return os.ReadFile(filepath.Join(s.workingDirectory, "pkg.dt"))
	}
	panic(fmt.Errorf("unknown package ID %#v", id))
}

func (s *Storage) Store(id PackageID, buf []byte) error {
	switch id.(type) {
	case *PackageID_Current:
		return os.WriteFile(filepath.Join(s.workingDirectory, "pkg.dt"), buf, os.ModePerm)
	}
	panic(fmt.Errorf("unknown package ID %#v", id))
}
