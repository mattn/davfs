package memory

import (
	"github.com/mattn/davfs"
	"golang.org/x/net/webdav"
)

func init() {
	davfs.Register("memory", &Driver{})
}

type Driver struct {
}

func (d *Driver) Mount(source string) (webdav.FileSystem, error) {
	return webdav.NewMemFS(), nil
}

func (d *Driver) CreateFS(source string) error {
	return nil
}
