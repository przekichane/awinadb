package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	CheckOS()

	fmt.Println("Installing latest google platform tools...")

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	const PLATFORM_TOOLS_URL string = "https://dl.google.com/android/repository/platform-tools-latest-windows.zip"
	const PLATFORM_TOOLS_NAME string = "platform-tools.zip"
	PLATFORM_TOOLS_PATH := filepath.Join(home, "platform-tools")
	PLATFORM_TOOLS_FILE := filepath.Join(cache, PLATFORM_TOOLS_NAME)

	// Handle the case where the directory already exists
	if _, err := os.Stat(PLATFORM_TOOLS_PATH); err == nil {
		fmt.Println("Directory " + PLATFORM_TOOLS_PATH + " already exist!")

		// Ask user whether to overwrite
		var response string
		fmt.Print("Overwrite? (y/n) ")
		fmt.Scanln(&response)

		if response == "y" || response == "Y" {
			err = os.RemoveAll(PLATFORM_TOOLS_PATH)
			if err != nil {
				panic("Failed to remove " + PLATFORM_TOOLS_PATH + "! Aborting...")
			}
		} else {
			panic("Aborting...")
		}
	}

	// Download platform tools archive
	fmt.Println("Downloading...")
	Download(PLATFORM_TOOLS_FILE, PLATFORM_TOOLS_URL)

	// Extract the platform tools
	fmt.Println("Extracting...")
	Unzip(PLATFORM_TOOLS_FILE, home)

	fmt.Println("Appending " + PLATFORM_TOOLS_PATH + " to user's Path...")
	appendToPath(PLATFORM_TOOLS_PATH)

	// Remove archive
	fmt.Println("Cleaning up...")
	os.Remove(PLATFORM_TOOLS_FILE)

	fmt.Println("All done! Your platform tools are installed to " + PLATFORM_TOOLS_PATH + "!")
	Exit(0)
}
