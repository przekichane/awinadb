package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
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

	fmt.Println("All done! Your platform tools are installed to C:\\platform-tools!")
	exit(0)
}
