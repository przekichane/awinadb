package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	CheckOS()

	fmt.Println("Installing latest google platform tools...")

	const PLATFORM_TOOLS_URL string = "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"
	const PLATFORM_TOOLS_PATH string = "C:\\platform-tools\\"
	const PLATFORM_TOOLS_NAME string = "platform-tools.zip"
	PLATFORM_TOOLS_FILE := filepath.Join(PLATFORM_TOOLS_PATH, PLATFORM_TOOLS_NAME)

	// Exit if directory already exists
	if _, err := os.Stat(PLATFORM_TOOLS_PATH); err == nil {
		fmt.Println("Directory " + PLATFORM_TOOLS_PATH + " already exist! Aborting...")
		Exit(1)
	}

	// Create inital directory
	os.Mkdir(PLATFORM_TOOLS_PATH, os.ModePerm)

	// Download platform tools archive
	fmt.Println("Downloading...")
	Download(PLATFORM_TOOLS_FILE, PLATFORM_TOOLS_URL)

	// Extract the platform tools
	fmt.Println("Extracting...")
	Unzip(PLATFORM_TOOLS_FILE, PLATFORM_TOOLS_PATH)

	// Remove archive
	fmt.Println("Cleaning up...")
	os.Remove(PLATFORM_TOOLS_FILE)

	fmt.Println("All done! Your platform tools are installed to C:\\platform-tools!")
	Exit(0)
}
