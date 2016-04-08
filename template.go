package gitest

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"runtime"
)

type template struct {
	Folder string
}

func newTemplate(name string) (*template, error) {
	_, f, _, _ := runtime.Caller(1)
	folder, err := filepath.Abs(fmt.Sprintf("%s/data/%s", path.Dir(f), name))
	if err != nil {
		return nil, err
	}
	_, err = ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	return &template{folder}, nil
}
