package presentation

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/geoffjay/pres/baml_client/types"
)

// Generator handles generating HTML output from presentations
type Generator struct {
	templatePath string
}

// NewGenerator creates a new HTML generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateHTML generates a reveal.js HTML file from presentation data
func (g *Generator) GenerateHTML(data *PresentationData, outputPath string) error {
	// Create output directory if it doesn't exist
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate HTML content
	html := g.buildHTML(data)

	// Write to file
	if err := os.WriteFile(outputPath, []byte(html), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file: %w", err)
	}

	return nil
}

// buildHTML constructs the complete HTML document
func (g *Generator) buildHTML(data *PresentationData) string {
	var sb strings.Builder

	// HTML header
	sb.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>`)
	sb.WriteString(template.HTMLEscapeString(data.Metadata.Title))
	sb.WriteString(`</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/dist/reset.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/dist/reveal.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/dist/theme/`)
	sb.WriteString(data.Metadata.Theme)
	sb.WriteString(`.css">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/plugin/highlight/monokai.css">
    <style>
        .reveal .slides section {
            text-align: left;
        }
        .reveal h1, .reveal h2, .reveal h3 {
            text-transform: none;
        }
        .two-column {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 2rem;
        }
    </style>
</head>
<body>
    <div class="reveal">
        <div class="slides">
`)

	// Generate slides
	for _, slide := range data.Slides {
		g.writeSlide(&sb, slide)
	}

	// HTML footer
	sb.WriteString(`        </div>
    </div>
    <script src="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/dist/reveal.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/plugin/notes/notes.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/plugin/markdown/markdown.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/reveal.js@5.1.0/plugin/highlight/highlight.js"></script>
    <script>
        Reveal.initialize({
            hash: true,
            slideNumber: true,
            plugins: [ RevealMarkdown, RevealHighlight, RevealNotes ]
        });
    </script>
</body>
</html>
`)

	return sb.String()
}

// writeSlide writes a single slide to the HTML
func (g *Generator) writeSlide(sb *strings.Builder, slide types.Slide) {
	// Start section with optional background color
	sb.WriteString("            <section")
	if slide.Background_color != "" {
		sb.WriteString(` data-background-color="`)
		sb.WriteString(template.HTMLEscapeString(slide.Background_color))
		sb.WriteString(`"`)
	}
	sb.WriteString(">\n")

	// Add slide title if present
	if slide.Title != "" {
		// Determine heading level based on layout
		headingLevel := "h2"
		if slide.Layout == "title" {
			headingLevel = "h1"
		}

		sb.WriteString("                <")
		sb.WriteString(headingLevel)
		sb.WriteString(">")
		sb.WriteString(template.HTMLEscapeString(slide.Title))
		sb.WriteString("</")
		sb.WriteString(headingLevel)
		sb.WriteString(">\n")
	}

	// Add content based on layout
	switch slide.Layout {
	case "two-column":
		g.writeTwoColumnContent(sb, slide.Content)
	default:
		// Standard content or blank slide
		if slide.Content != "" {
			sb.WriteString("                <div data-markdown>\n")
			sb.WriteString("                    <textarea data-template>\n")
			sb.WriteString(slide.Content)
			sb.WriteString("\n                    </textarea>\n")
			sb.WriteString("                </div>\n")
		}
	}

	// Add speaker notes if present
	if slide.Notes != "" {
		sb.WriteString("                <aside class=\"notes\">\n")
		sb.WriteString("                    ")
		sb.WriteString(template.HTMLEscapeString(slide.Notes))
		sb.WriteString("\n                </aside>\n")
	}

	// End section
	sb.WriteString("            </section>\n")
}

// writeTwoColumnContent writes content in a two-column layout
func (g *Generator) writeTwoColumnContent(sb *strings.Builder, content string) {
	// Split content by a delimiter (e.g., "---" or "|||")
	columns := strings.Split(content, "|||")
	if len(columns) < 2 {
		columns = strings.Split(content, "---")
	}

	sb.WriteString("                <div class=\"two-column\">\n")

	for i, col := range columns {
		if i >= 2 {
			break // Only support two columns
		}
		sb.WriteString("                    <div data-markdown>\n")
		sb.WriteString("                        <textarea data-template>\n")
		sb.WriteString(strings.TrimSpace(col))
		sb.WriteString("\n                        </textarea>\n")
		sb.WriteString("                    </div>\n")
	}

	sb.WriteString("                </div>\n")
}

// GetRevealJSThemes returns the list of available reveal.js themes
func GetRevealJSThemes() []string {
	return []string{
		"black",
		"white",
		"league",
		"beige",
		"sky",
		"night",
		"serif",
		"simple",
		"solarized",
	}
}
