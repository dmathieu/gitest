package gitest

import (
	"fmt"
	"io/ioutil"
)

type template struct {
	Folder string
}

func newTemplate(name string) (*template, error) {
	folder := fmt.Sprintf("data/%s", name)
	_, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	return &template{folder}, nil
}
