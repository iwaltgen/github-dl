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
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// TarGz unarchives tar.gz(tgz) archive file.
type TarGz struct{}

// Unarchive unpacks the .zip file at source to destination.
func (t TarGz) Unarchive(source, destination string) error {
	sf, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source file error: %w", err)
	}
	defer sf.Close()

	gr, err := gzip.NewReader(sf)
	if err != nil {
		return fmt.Errorf("open gzip reader error: %w", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	if err := mkdir(destination, os.ModePerm); err != nil {
		return err
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reader next error: %w", err)
		}

		fpath := filepath.Join(destination, header.Name)
		fmode := os.FileMode(header.Mode)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := mkdir(fpath, fmode); err != nil {
				return err
			}

		case tar.TypeReg, tar.TypeRegA, tar.TypeChar, tar.TypeBlock, tar.TypeFifo, tar.TypeGNUSparse:
			if err := writeNewFile(fpath, tr, fmode); err != nil {
				return err
			}

		case tar.TypeXGlobalHeader, tar.TypeSymlink, tar.TypeLink: // ignore

		default:
			return fmt.Errorf("unknown type error: %v", header.Typeflag)
		}
	}
	return nil
}
