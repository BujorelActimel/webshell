package main

import (
	"fmt"
	"log"

	"webshell/internal/browser"
	"webshell/internal/ui"
	"webshell/internal/utils"
)

func main() {
	// Setup logging
	logFile, err := utils.SetupLogging("terminal_browser.log")
	if err != nil {
		fmt.Printf("Failed to set up logging: %v\n", err)
		return
	}
	defer logFile.Close()
	log.Printf("Starting terminal browser")

	// Enable terminal features
	if err := ui.EnableVirtualTerminalProcessing(); err != nil {
		log.Printf("Failed to enable virtual terminal processing: %v", err)
		fmt.Printf("Failed to enable virtual terminal processing: %v\n", err)
		return
	}

	// Setup keyboard
	if err := ui.InitKeyboard(); err != nil {
		log.Printf("Failed to initialize keyboard: %v", err)
		fmt.Printf("Failed to initialize keyboard: %v\n", err)
		return
	}
	defer ui.CloseKeyboard()

	// Setup browser
	ctx, cancelBrowser := browser.InitBrowser(log.Printf)
	defer cancelBrowser()

	// Setup application state
	state := browser.NewBrowserState(1024, 800)

	// Initialize viewport
	if err := browser.SetupViewport(ctx, state.ViewportWidth, state.ViewportHeight); err != nil {
		log.Fatal(err)
	}

	// Main application loop
	for {
		// Display the prompt and get URL input
		url := ui.DisplayPrompt()
		if url == "quit" {
			break
		}

		// Ensure URL has proper scheme
		if !browser.HasPrefix(url, "http://") && !browser.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		// Browse the URL
		browser.BrowseURL(ctx, url, state)
	}

	// Clear screen on exit
	fmt.Print("\x1b[2J\x1b[H")
}
