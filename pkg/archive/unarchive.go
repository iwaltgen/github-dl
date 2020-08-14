/*
Copyright © 2020 iwaltgen

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package archive

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrNotSupportFile does not support file extensions.
	ErrNotSupportFile = errors.New("not supoort file extension")
)

// Unarchiver is a type that can extract archive files into a folder.
type Unarchiver interface {
	Unarchive(source, destination string) error
}

// Unarchive unarchives the given archive file into the destination folder.
// The archive format is selected implicitly.
func Unarchive(source, destination string) error {
	unarchiver, err := byExtension(source)
	if err != nil {
		return fmt.Errorf("unarchive `%s` error: %w", source, err)
	}

	if err := unarchiver.Unarchive(source, destination); err != nil {
		return fmt.Errorf("unarchive `%s` error: %w", source, err)
	}
	return nil
}

// Support check for handle the archive file format.
func Support(fpath string) bool {
	_, err := byExtension(fpath)
	return err == nil
}

func byExtension(fpath string) (Unarchiver, error) {
	switch {
	case strings.HasSuffix(fpath, ".zip"):
		return Zip{}, nil

	case strings.HasSuffix(fpath, ".tar.gz"),
		strings.HasSuffix(fpath, ".tgz"):
		return TarGz{}, nil

	default:
		return nil, ErrNotSupportFile
	}
}

func mkdir(dpath string, mode os.FileMode) error {
	err := os.MkdirAll(dpath, mode)
	if err != nil {
		return fmt.Errorf("mkdir `%s` error: %w", dpath, err)
	}
	return nil
}

func writeNewFile(fpath string, in io.Reader, mode os.FileMode) error {
	err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir `%s` for file error: %w", fpath, err)
	}

	out, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("create file `%s` error: %w", fpath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("write file `%s`: %w", fpath, err)
	}
	return nil
}
