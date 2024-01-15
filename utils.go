package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Exit(code int) {
	fmt.Println("Press Enter key to exit...")
	fmt.Scanln()
	os.Exit(code)
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		fmt.Println("Extracting " + f.Name + " to " + path)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func Download(path string, url string) error {
	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func appendToPath(value string) {
	key := "Path"

	u, _ := user.Current()
	SID := u.Uid

	k, err := registry.OpenKey(registry.USERS, SID+`\Environment`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	currentPath, _, err := k.GetStringValue(key)
	if err != nil {
		panic(err)
	}

	if strings.Contains(currentPath, value) {
		fmt.Println(value + " already in the user's Path, skipping")
		return
	}

	err = k.SetStringValue(key, currentPath+";"+value)
	if err != nil {
		panic(err)
	}
}
