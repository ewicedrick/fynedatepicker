package copyFiles

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func copyFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(destination)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Sync()
}

func copytxtFiles(dstDir, srcDir string) error {
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return err
	}

	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".txt" {
			dstPath := filepath.Join(dstDir, filepath.Base(path))
			err = copyFile(path, dstPath)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
