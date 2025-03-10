package browser

import (
	"context"
	"fmt"
	"log"
	"os"

	"webshell/internal/ui"

	"github.com/chromedp/chromedp"
)

// Link represents a hyperlink on a webpage
type Link struct {
	X1, Y1, X2, Y2 int
	URL            string
}

// BrowserState holds the current state of the browser
type BrowserState struct {
	ViewportWidth  int
	ViewportHeight int
	ImageWidth     int
	ImageHeight    int
	Links          []Link
	SelectedLink   int
}

// NewBrowserState creates a new browser state with default values
func NewBrowserState(width, height int) *BrowserState {
	return &BrowserState{
		ViewportWidth:  width,
		ViewportHeight: height,
		SelectedLink:   0,
	}
}

// GetSelectedLink returns the currently selected link index
func (s *BrowserState) GetSelectedLink() int {
	return s.SelectedLink
}

// SetSelectedLink updates the selected link index
func (s *BrowserState) SetSelectedLink(index int) {
	s.SelectedLink = index
}

// GetTotalLinks returns the total number of links
func (s *BrowserState) GetTotalLinks() int {
	return len(s.Links)
}

// GetLinkURL returns the URL of the link at the specified index
func (s *BrowserState) GetLinkURL(index int) string {
	if index >= 0 && index < len(s.Links) {
		return s.Links[index].URL
	}
	return ""
}

// InitBrowser sets up a new ChromeDP browser instance
func InitBrowser(logf func(format string, v ...interface{})) (context.Context, context.CancelFunc) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(logf),
	)
	return ctx, cancel
}

// SetupViewport configures the browser viewport size
func SetupViewport(ctx context.Context, width, height int) error {
	return chromedp.Run(ctx,
		chromedp.EmulateViewport(int64(width), int64(height)),
	)
}

// BrowseURL handles the browsing session for a specific URL
func BrowseURL(ctx context.Context, url string, state *BrowserState) {
	log.Printf("Navigating to URL: %s", url)

	// Navigate to URL
	if err := chromedp.Run(ctx, chromedp.Navigate(url)); err != nil {
		log.Printf("Error navigating: %v", err)
		fmt.Printf("Error navigating: %v\n", err)
		return
	}

	// Process page and display it
	if err := CaptureAndDisplayWebpage(ctx, state); err != nil {
		log.Printf("Error: %v", err)
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Process user input while browsing
	var shouldQuit bool
	for {
		quit, refresh, err := ui.HandleInput(state)
		if err != nil {
			log.Printf("Error handling input: %v", err)
			continue
		}

		if quit {
			shouldQuit = true
			break
		}

		if refresh {
			var currentURL string
			if err := chromedp.Run(ctx,
				chromedp.Location(&currentURL),
			); err != nil {
				log.Printf("Error getting current URL: %v", err)
				break
			}
			url = currentURL
			break
		}
	}

	// Quit if requested
	if shouldQuit {
		return
	}
}

// CaptureAndDisplayWebpage takes a screenshot and displays it
func CaptureAndDisplayWebpage(ctx context.Context, state *BrowserState) error {
	var buf []byte
	if err := chromedp.Run(ctx,
		chromedp.WaitReady("body"),
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		return fmt.Errorf("failed to capture screenshot: %v", err)
	}

	tempFile := "temp_screenshot.png"
	if err := os.WriteFile(tempFile, buf, 0644); err != nil {
		return fmt.Errorf("failed to save screenshot: %v", err)
	}
	defer os.Remove(tempFile)

	imageData, err := os.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("failed to read screenshot: %v", err)
	}

	// Reset terminal state
	ui.ResetTerminal()

	// Get page links
	links, err := GetPageLinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get links: %v", err)
	}
	state.Links = links
	state.SelectedLink = 0

	// Display the image
	ui.DisplayImage(imageData)

	// Show selected link if available
	if len(state.Links) > 0 {
		ui.DisplaySelectedLink(state)
	}

	return nil
}

// GetPageLinks extracts all hyperlinks from the current page
func GetPageLinks(ctx context.Context) ([]Link, error) {
	log.Printf("Getting page links")
	var links []Link
	var js = `
    Array.from(document.links).map(link => {
        const rect = link.getBoundingClientRect();
        return {
            url: link.href,
            x1: Math.round(rect.left),
            y1: Math.round(rect.top),
            x2: Math.round(rect.right),
            y2: Math.round(rect.bottom)
        };
    })
    `
	var rawLinks []map[string]interface{}
	if err := chromedp.Run(ctx,
		chromedp.Evaluate(js, &rawLinks),
	); err != nil {
		log.Printf("Failed to evaluate JavaScript for links: %v", err)
		return nil, fmt.Errorf("failed to get links: %v", err)
	}

	for _, raw := range rawLinks {
		link := Link{
			X1:  int(raw["x1"].(float64)),
			Y1:  int(raw["y1"].(float64)),
			X2:  int(raw["x2"].(float64)),
			Y2:  int(raw["y2"].(float64)),
			URL: raw["url"].(string),
		}
		links = append(links, link)
		log.Printf("Found link: %+v", link)
	}
	log.Printf("Total links found: %d", len(links))
	return links, nil
}

// Helper function to check string prefix
func HasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[0:len(prefix)] == prefix
}
