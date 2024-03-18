package modregistry

import (
	"io"
	"os"
	"path/filepath"
)

func copyFile(src string, dest string) error {
	sf, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer sf.Close()
	return putFileContents(dest, sf)
}

func putFileContents(filename string, r io.Reader) error {
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	df, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, r)
	return nil
}

func getEnv(name string, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
