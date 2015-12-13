package file

import (
	"github.com/mattn/davfs"
	"golang.org/x/net/webdav"
	"os"
	"path/filepath"
)

func init() {
	davfs.Register("file", &Driver{})
}

type Driver struct {
}

func (d *Driver) Mount(source string) (webdav.FileSystem, error) {
	if source == "" {
		source = "."
	}
	if s, err := filepath.Abs(source); err == nil {
		source = s
	}
	return webdav.Dir(source), nil
}

func (d *Driver) CreateFS(source string) error {
	if source == "" {
		source = "."
	}
	if s, err := filepath.Abs(source); err == nil {
		source = s
	}
	return os.Mkdir(source, 0755)
}
