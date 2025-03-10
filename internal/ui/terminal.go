package ui

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"golang.org/x/sys/windows"
	"golang.org/x/term"
)

const (
	colorReset = "\x1b[0m"
	colorCyan  = "\x1b[36m"
)

const logo = `
	 █     █░▓█████  ▄▄▄▄     ██████  ██░ ██ ▓█████  ██▓     ██▓    
	▓█░ █ ░█░▓█   ▀ ▓█████▄ ▒██    ▒ ▓██░ ██▒▓█   ▀ ▓██▒    ▓██▒    
	▒█░ █ ░█ ▒███   ▒██▒ ▄██░ ▓██▄   ▒██▀▀██░▒███   ▒██░    ▒██░    
	░█░ █ ░█ ▒▓█  ▄ ▒██░█▀    ▒   ██▒░▓█ ░██ ▒▓█  ▄ ▒██░    ▒██░    
	░░██▒██▓ ░▒████▒░▓█  ▀█▓▒██████▒▒░▓█▒░██▓░▒████▒░██████▒░██████▒
	░ ▓░▒ ▒  ░░ ▒░ ░░▒▓███▀▒▒ ▒▓▒ ▒ ░ ▒ ░░▒░▒░░ ▒░ ░░ ▒░▓  ░░ ▒░▓  ░
	  ▒ ░ ░   ░ ░  ░▒░▒   ░ ░ ░▒  ░ ░ ▒ ░▒░ ░ ░ ░  ░░ ░ ▒  ░░ ░ ▒  ░
	  ░   ░     ░    ░    ░ ░  ░  ░   ░  ░░ ░   ░     ░ ░     ░ ░   
  	  ░       ░  ░ ░            ░   ░  ░  ░   ░  ░    ░  ░    ░  ░`

// EnableVirtualTerminalProcessing enables ANSI escape sequence support on Windows
func EnableVirtualTerminalProcessing() error {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	err := windows.GetConsoleMode(stdout, &originalMode)
	if err != nil {
		return err
	}

	mode := originalMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING | windows.ENABLE_PROCESSED_OUTPUT

	err = windows.SetConsoleMode(stdout, mode)
	if err != nil {
		return err
	}

	return nil
}

// InitKeyboard initializes the keyboard input handler
func InitKeyboard() error {
	return keyboard.Open()
}

// CloseKeyboard closes the keyboard input handler
func CloseKeyboard() {
	keyboard.Close()
}

// DisplayPrompt shows the URL input prompt and returns the user input
func DisplayPrompt() string {
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	// Clear screen and move to top
	fmt.Print("\x1b[2J\x1b[H")

	// Print centered logo
	logoLines := strings.Split(logo, "\n")
	for _, line := range logoLines {
		padding := (termWidth - len(line)) / 2
		if padding > 0 {
			fmt.Print(strings.Repeat(" ", padding))
		}
		fmt.Printf("%s%s%s\n", colorCyan, line, colorReset)
	}
	fmt.Println() // Add some spacing
	fmt.Print("\tEnter URL (or 'quit' to exit): ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// ResetTerminal resets the terminal state completely
func ResetTerminal() {
	fmt.Print("\x1b[!p") // Soft reset
	fmt.Print("\x1b[3J") // Clear scrollback
	fmt.Print("\x1b[2J") // Clear screen
	fmt.Print("\x1bc")   // Full reset
	fmt.Print("\x1b[H")  // Home
}

// DisplayImage displays a base64 encoded image in the terminal
func DisplayImage(imageData []byte) {
	b64Image := base64.StdEncoding.EncodeToString(imageData)
	frame := fmt.Sprintf("\x1b]1337;File=inline=1;width=100%%:%s\x07", b64Image)
	fmt.Print(frame)
}

// DisplaySelectedLink shows information about the currently selected link
func DisplaySelectedLink(state interface{}) {
	// Use type assertion with interface to avoid circular imports
	selectedLink := -1
	totalLinks := 0
	var linkURL string

	// Check if state implements our required methods
	if s, ok := state.(interface {
		GetSelectedLink() int
		GetTotalLinks() int
		GetLinkURL(int) string
	}); ok {
		selectedLink = s.GetSelectedLink()
		totalLinks = s.GetTotalLinks()
		if selectedLink >= 0 && selectedLink < totalLinks {
			linkURL = s.GetLinkURL(selectedLink)
		}
	}

	if selectedLink >= 0 && selectedLink < totalLinks {
		fmt.Print("\x1b[?25l")                    // Hide cursor
		fmt.Print("\x1b[", os.Stdout.Fd(), ";0H") // Move to last line
		fmt.Print("\x1b[7m")                      // Inverse colors for status bar
		statusText := fmt.Sprintf(" [%d/%d] %s ",
			selectedLink+1,
			totalLinks,
			linkURL)
		fmt.Print("\x1b[2K") // Clear line
		fmt.Print(statusText)
		if termWidth, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
			padding := termWidth - len(statusText)
			if padding > 0 {
				fmt.Print(strings.Repeat(" ", padding))
			}
		}
		fmt.Print("\x1b[0m") // Reset colors
	}
}

// HandleInput processes keyboard input during browsing
func HandleInput(state interface{}) (bool, bool, error) {
	// Use type assertion with interface to avoid circular imports
	selectedLink := 0
	totalLinks := 0
	var setSelectedLink func(int)

	// Check if state implements our required methods
	if s, ok := state.(interface {
		GetSelectedLink() int
		SetSelectedLink(int)
		GetTotalLinks() int
	}); ok {
		selectedLink = s.GetSelectedLink()
		totalLinks = s.GetTotalLinks()
		setSelectedLink = s.SetSelectedLink
	}

	char, key, err := keyboard.GetKey()
	if err != nil {
		log.Printf("Error getting key: %v", err)
		return false, false, err
	}

	log.Printf("Input received - Char: %d, Key: %v", char, key)

	// Default values if interface isn't properly implemented
	if setSelectedLink == nil {
		setSelectedLink = func(int) {} // Empty function as fallback
	}

	switch key {
	case keyboard.KeyArrowUp:
		if selectedLink > 0 && setSelectedLink != nil {
			setSelectedLink(selectedLink - 1)
			DisplaySelectedLink(state)
		}
		return false, false, nil
	case keyboard.KeyArrowDown:
		if selectedLink < totalLinks-1 && setSelectedLink != nil {
			setSelectedLink(selectedLink + 1)
			DisplaySelectedLink(state)
		}
		return false, false, nil
	case keyboard.KeyEnter:
		if selectedLink >= 0 && selectedLink < totalLinks {
			fmt.Print("\x1b[2J\x1b[H") // Clear screen
			fmt.Printf("Navigating to link...")
			return false, true, nil
		}
	}

	if char == 'q' || char == 'Q' {
		log.Printf("Quit command received")
		return true, false, nil
	}

	return false, false, nil
}
