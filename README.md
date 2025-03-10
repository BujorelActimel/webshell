# WebShell - Terminal Web Browser

WebShell is a terminal-based web browser that renders websites as images directly in your terminal window and allows for keyboard-based navigation. It uses ChromeDP for rendering and iTerm2-compatible image protocols for display.

## Features

* **Visual Browsing** : Renders full website screenshots directly in the terminal
* **Keyboard Navigation** : Browse websites using keyboard shortcuts
* **Link Detection** : Automatically detects and highlights clickable links
* **Clean Interface** : Minimalist UI with a focus on content
* **Multi-platform** : Works on Windows with terminal supporting inline images

## Demo

<video src="https://github.com/BujorelActimel/webshell/raw/refs/heads/main/demo/demo.mp4" controls="controls">
</video>

## Requirements

* Go 1.24 or higher
* Windows operating system
* Terminal with support for inline images (like WezTerm)
* Chrome or Chromium browser installed (used by ChromeDP)
* Internet connection

## Installation

### From Source

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/webshell.git
   cd webshell
   ```
2. Build the application:
   ```
   go build -o bin/webshell.exe ./cmd/webshell
   ```
3. Run WebShell:
   ```
   ./bin/webshell.exe
   ```

## Usage

Start WebShell by running the executable:

```
webshell
```

### Navigation

* Enter a URL at the prompt or type 'quit' to exit
* Use arrow keys to navigate between links
* Press Enter to follow the selected link
* Press 'q' to return to the URL prompt

## How It Works

WebShell uses ChromeDP to control a headless Chrome browser, navigating to web pages and capturing screenshots. These screenshots are then rendered directly in the terminal using a special image protocol compatible with modern terminal emulators. Link detection is performed by executing JavaScript in the browser context to identify all clickable elements.

## Known Limitations

* Some websites may not render correctly in headless mode
* Terminal must support inline image display (iTerm2 protocol)
* Interaction is limited to clicking links (no forms or JavaScript interaction)
* Screen size may affect website layout
