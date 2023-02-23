package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func exit(code int) {
	fmt.Println("Press Enter key to exit...")
	fmt.Scanln()
	os.Exit(code)
}

func checkOS() {
	if runtime.GOOS != "windows" {
		fmt.Println("This program works only on windows! Aborting")
		exit(1)
	}

}

func unzip(src, dest string) error {
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

func download(path string, url string) error {
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

func main() {
	checkOS()

	fmt.Println("Installing latest google platform tools...")

	const PLATFORM_TOOLS_PATH string = "C:\\platform-tools\\"

	// Exit if directory already exists
	if _, err := os.Stat(PLATFORM_TOOLS_PATH); err == nil {
		fmt.Println("Directory " + PLATFORM_TOOLS_PATH + " already exist! Aborting...")
		exit(1)
	}

	// Create inital directory
	os.Mkdir(PLATFORM_TOOLS_PATH, os.ModePerm)

	const PLATFORM_TOOLS_URL string = "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"

	// Download platform tools archive
	download(PLATFORM_TOOLS_PATH+"platform-tools.zip", PLATFORM_TOOLS_URL)

	// Extract the platform tools
	unzip(PLATFORM_TOOLS_PATH+"platform-tools.zip", PLATFORM_TOOLS_PATH)

	// Remove archive
	os.Remove(PLATFORM_TOOLS_PATH + "platform-tools.zip")

	fmt.Println("All done! Your platform tools are installed to C:\\platform-tools!")
	exit(0)
}
