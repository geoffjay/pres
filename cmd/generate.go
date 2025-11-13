package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/geoffjay/pres/internal/presentation"
	"github.com/spf13/cobra"
)

var (
	generatePath   string
	generateOutput string
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate HTML output from a presentation",
	Long: `Generate a reveal.js HTML file from a presentation JSON file.

The command will:
1. Load the presentation from JSON
2. Generate a reveal.js HTML file with all slides
3. Include speaker notes and styling

The generated HTML file can be opened directly in a browser.

Examples:
  pres generate --path presentations/my-talk.json
  pres generate --path presentations/review.json --output output/review.html`,
	RunE: runGenerate,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generatePath, "path", "p", "", "Path to presentation JSON file (required)")
	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "Output path for HTML file (default: same name as JSON with .html extension)")
	generateCmd.MarkFlagRequired("path")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	fmt.Printf("📄 Generating HTML from: %s\n", generatePath)

	// Load presentation
	writer := presentation.NewWriter(".")
	data, err := writer.LoadPresentation(generatePath)
	if err != nil {
		return fmt.Errorf("failed to load presentation: %w", err)
	}

	fmt.Printf("Loaded: %s (%d slides)\n", data.Metadata.Title, len(data.Slides))

	// Determine output path
	outputPath := generateOutput
	if outputPath == "" {
		// Use same directory and name as input, but with .html extension
		dir := filepath.Dir(generatePath)
		base := filepath.Base(generatePath)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		outputPath = filepath.Join(dir, name+".html")
	}

	// Generate HTML
	fmt.Println("\nGenerating reveal.js HTML...")
	generator := presentation.NewGenerator()
	if err := generator.GenerateHTML(data, outputPath); err != nil {
		return fmt.Errorf("failed to generate HTML: %w", err)
	}

	fmt.Printf("\n✓ HTML generated successfully!\n")
	fmt.Printf("  Location: %s\n", outputPath)
	fmt.Printf("  Title: %s\n", data.Metadata.Title)
	fmt.Printf("  Theme: %s\n", data.Metadata.Theme)
	fmt.Printf("  Slides: %d\n", len(data.Slides))

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  • Open in browser: open %s\n", outputPath)
	fmt.Printf("  • Or start a local server: python3 -m http.server 8000\n")
	fmt.Printf("    Then visit: http://localhost:8000/%s\n", outputPath)

	return nil
}
