package gitest

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type template struct {
	Folder string
}

func newTemplate(name string) (*template, error) {
	_, f, _, _ := runtime.Caller(1)
	c, err := filepath.Abs(fmt.Sprintf("%s/%s.tgz", path.Dir(f), name))
	if err != nil {
		return nil, err
	}
	folder, err := decompress(c)
	if err != nil {
		return nil, err
	}

	fmt.Println(folder)
	return &template{folder}, nil
}

func decompress(tarPath string) (string, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileReader, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer fileReader.Close()

	tarBallReader := tar.NewReader(fileReader)

	destPath, err := ioutil.TempDir("", "decompressed_repo")
	if err != nil {
		return "", err
	}

	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		filename := path.Join(destPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(filename, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

		case tar.TypeReg:
			writer, err := os.OpenFile(filename, os.O_CREATE+os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

			io.Copy(writer, tarBallReader)
			if err != nil {
				return "", err
			}

			writer.Close()
		default:
			return "", fmt.Errorf("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}

	return destPath, nil
}
