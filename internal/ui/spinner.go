package ui

import (
	"fmt"
	"os"
	"time"
)

var spinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ShowSpinner shows a simple spinner (simplified version)
func ShowSpinner(message string) {
	fmt.Fprintf(os.Stderr, "\r%s %s", spinnerChars[0], message)
}

// SpinnerWithFunc runs a function with a simple progress indicator
func SpinnerWithFunc(message string, fn func() error) error {
	fmt.Fprintf(os.Stderr, "⏳ %s...", message)
	
	if err := fn(); err != nil {
		fmt.Fprintf(os.Stderr, "\r❌ %s failed\n", message)
		return err
	}
	
	fmt.Fprintf(os.Stderr, "\r%s %s\n", SuccessText.Render("✅"), message)
	time.Sleep(100 * time.Millisecond)
	return nil
}

