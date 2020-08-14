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
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
)

// Zip unarchives zip archive file.
type Zip struct{}

// Unarchive unpacks the .zip file at source to destination.
func (z Zip) Unarchive(source, destination string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("open reader error: %w", err)
	}
	defer r.Close()

	if err := mkdir(destination, os.ModePerm); err != nil {
		return err
	}

	for _, zf := range r.File {
		f, err := zf.Open()
		if err != nil {
			return fmt.Errorf("open file `%s` error: %w", zf.Name, err)
		}
		defer f.Close()

		fileinfo := zf.FileInfo()
		fpath := filepath.Join(destination, zf.Name)
		if fileinfo.IsDir() {
			if err := mkdir(fpath, fileinfo.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := writeNewFile(fpath, f, fileinfo.Mode()); err != nil {
			return err
		}
	}
	return nil
}
