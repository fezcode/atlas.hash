package main

import (
	"fmt"
	"os"

	"atlas.hash/internal/hash"
	"atlas.hash/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("atlas.hash v%s\n", Version)
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help") {
		showHelp()
		return
	}

	if len(os.Args) > 1 {
		// Non-interactive mode
		filePath := os.Args[1]
		res, err := hash.Compute(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error computing hashes: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Target: %s\n\n", filePath)
		fmt.Printf("%-10s %s\n", "MD5", res.MD5)
		fmt.Printf("%-10s %s\n", "SHA1", res.SHA1)
		fmt.Printf("%-10s %s\n", "SHA256", res.SHA256)
		fmt.Printf("%-10s %s\n", "SHA512", res.SHA512)
		return
	}

	// Interactive TUI Mode
	p := tea.NewProgram(ui.NewModel(""))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Atlas Hash - A fast, minimalist hash utility for your terminal.")
	fmt.Println("\nUsage:")
	fmt.Println("  atlas.hash               Start the interactive TUI to select a file")
	fmt.Println("  atlas.hash <file>        Compute and display hashes for the given file (non-interactive)")
	fmt.Println("  atlas.hash -v, --version Show version information")
	fmt.Println("  atlas.hash -h, --help    Show this help information")
	fmt.Println("\nTUI Controls:")
	fmt.Println("  Type/Paste     Enter file path or hash to compare")
	fmt.Println("  Enter          Confirm input")
	fmt.Println("  Esc            Quit the application")
}
