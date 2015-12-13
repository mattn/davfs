package davfs

import (
	"os"

	"golang.org/x/net/webdav"
)

type Driver interface {
	Mount(source string) (webdav.FileSystem, error)
	CreateFS(source string) error
}

var drivers = map[string]Driver{}

func Register(name string, driver Driver) {
	drivers[name] = driver
}

func NewFS(driver, source string) (webdav.FileSystem, error) {
	if d, ok := drivers[driver]; ok {
		return d.Mount(source)
	}
	return nil, os.ErrNotExist
}

func CreateFS(driver, source string) error {
	if d, ok := drivers[driver]; ok {
		return d.CreateFS(source)
	}
	return os.ErrNotExist
}
