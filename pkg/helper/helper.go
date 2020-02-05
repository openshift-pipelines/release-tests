package helper

import (
	"path"
	"path/filepath"
	"runtime"
)

// RootDir returns you the root directory of this package as string
func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
