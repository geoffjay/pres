# pres - AI-Powered Presentation Generator

`pres` is a CLI utility for creating, updating, and generating presentations using AI-powered iterative Q&A and
reveal.js.

## Features

- **Interactive Creation**: Create presentations through an intelligent Q&A process
- **Smart Updates**: Update existing presentations with natural language requests
- **HTML Generation**: Generate reveal.js HTML presentations ready for the browser
- **Iterative Refinement**: AI-powered confidence scoring ensures you provide the right amount of context
- **Multiple Themes**: Support for all reveal.js themes
- **Speaker Notes**: Include speaker notes with your slides

## Installation

```bash
go build -o pres .
```

## Fixed Issues

- **v0.2.0**: Fixed directory creation bug - now creates full directory path for presentations
- **v0.2.0**: Added support for loading both wrapped (`PresentationData`) and raw BAML (`Presentation`) JSON formats

## Quick Start

### 1. Create a Presentation

```bash
./pres create "Introduction to Go concurrency patterns"
```

This will:

1. Ask you contextual questions about your presentation
2. Use AI to determine if more information is needed
3. Generate slides based on your responses
4. Save to `presentations/<title>.json`

### 2. Update a Presentation

```bash
./pres update --path presentations/my-talk.json "Add a slide about context.Context at the beginning"
```

This will:

1. Load your existing presentation
2. Ask clarifying questions about the update
3. Apply the changes intelligently
4. Save the updated presentation

### 3. Generate HTML

```bash
./pres generate --path presentations/my-talk.json
```

This will:

1. Load your presentation JSON
2. Generate a reveal.js HTML file
3. Output to `presentations/my-talk.html`

Then open in your browser:

```bash
open presentations/my-talk.html
```

Or serve with a local server:

```bash
python3 -m http.server 8000
# Visit http://localhost:8000/presentations/my-talk.html
```

## Commands

### `pres create [description]`

Create a new presentation with an interactive Q&A process.

**Flags:**

- `--author string` - Author name (default: empty)
- `--output string` - Output path (default: auto-generated from title)

**Examples:**

```bash
pres create "Introduction to Go concurrency patterns"
pres create "Q4 Business Review" --author "Jane Doe"
pres create "Product Launch" --output presentations/launch.json
```

### `pres update [request]`

Update an existing presentation using natural language.

**Flags:**

- `--path string` - Path to presentation JSON (required)

**Examples:**

```bash
pres update --path presentations/my-talk.json "Add an executive summary slide at the beginning"
pres update --path presentations/review.json "Change the theme to 'night'"
pres update --path presentations/intro.json "Add more code examples to the goroutines slide"
```

### `pres generate`

Generate reveal.js HTML from a presentation JSON file.

**Flags:**

- `--path string` - Path to presentation JSON (required)
- `--output string` - Output HTML path (default: same as input with .html extension)

**Examples:**

```bash
pres generate --path presentations/my-talk.json
pres generate --path presentations/review.json --output output/review.html
```

## Presentation Format

Presentations are stored as JSON files with the following structure:

```json
{
  "metadata": {
    "title": "Introduction to Go Concurrency",
    "subtitle": "Patterns and Best Practices",
    "author": "Jane Doe",
    "date": "2025-01-15",
    "theme": "black",
    "tags": ["go", "concurrency", "programming"],
    "created": "2025-01-15T10:00:00Z",
    "modified": "2025-01-15T10:00:00Z"
  },
  "slides": [
    {
      "title": "Introduction",
      "content": "# Welcome\n\nToday we'll explore...",
      "notes": "Start with a warm welcome...",
      "layout": "title",
      "background_color": ""
    }
  ]
}
```

## Slide Layouts

- `title` - Large centered text for section introductions
- `content` - Standard content slide with title and bullet points
- `two-column` - Split content into two columns (use `|||` to separate)
- `blank` - Minimal slide for images or quotes

## reveal.js Themes

Available themes:

- `black` - Dark background, white text (modern, professional)
- `white` - White background, dark text (clean, minimal)
- `league` - Gray background (neutral, versatile)
- `beige` - Beige background (warm, approachable)
- `sky` - Sky blue background (calm, friendly)
- `night` - Black background with orange highlights (bold, energetic)
- `serif` - Serif fonts (classic, formal)
- `simple` - Simple and minimal (understated)
- `solarized` - Solarized colors (eye-friendly, technical)

## Iterative Q&A Process

The CLI uses an intelligent iterative Q&A process:

1. **Iteration 0**: Core information (audience, purpose, key message)
2. **Iteration 1**: Structure and details (topics, depth, constraints)
3. **Iteration 2**: Refinements (examples, visual preferences)

The AI assigns a confidence score at each iteration. If confidence is high enough, it proceeds. Otherwise, it asks follow-up questions.

## Shared Library

The `pkg/tui` package contains reusable TUI components that can be used in other projects:

```go
import "github.com/geoffjay/pres/pkg/tui"

config := tui.IterationConfig{
    MaxIterations:    3,
    CompletionPrompt: "Need more info?",
}

form := tui.NewIterativeForm("My Form", config)
// Add questions and run...
```

## Architecture

- **BAML Functions** (`baml_src/presentations.baml`) - AI prompts for generation
- **Shared TUI** (`pkg/tui/`) - Reusable interactive form components
- **Internal Packages** (`internal/presentation/`) - Core logic
  - `writer.go` - JSON storage and updates
  - `generator.go` - HTML generation
- **CLI Commands** (`cmd/`) - Command implementations

## BAML Integration

This project uses [BAML](https://www.boundaryml.com/) for structured AI interactions. The BAML functions define:

- Question generation with confidence scoring
- Presentation structure and content
- Update operations

To regenerate the BAML client after modifying `.baml` files:

```bash
baml-cli generate
```

## Environment Variables

Set `ANTHROPIC_API_KEY` for Claude models:

```bash
export ANTHROPIC_API_KEY=your_key_here
```

## Examples

### Create a Technical Presentation

```bash
./pres create "Advanced Rust Async Programming" --author "John Smith"
```

Sample Q&A:

- **Q**: Who is your target audience?
- **A**: Intermediate Rust developers familiar with basics

- **Q**: What's the main goal of this presentation?
- **A**: Teach tokio patterns and best practices

- **Q**: How long should it be?
- **A**: 45 minutes with code examples

### Update an Existing Presentation

```bash
./pres update --path presentations/rust-async.json "Add a slide about error handling patterns after the tokio basics slide"
```

### Generate Multiple Formats

```bash
# Generate default HTML
./pres generate --path presentations/rust-async.json

# Generate with custom output
./pres generate --path presentations/rust-async.json --output public/rust-async.html
```

## Contributing

This project demonstrates patterns for building AI-powered CLI tools with:

- Iterative Q&A flows
- Confidence-based information gathering
- Structured output generation
- Reusable TUI components

Feel free to use the shared library (`pkg/tui`) in your own projects.

## License

MIT

## Related Projects

- [kb](contrib/kb) - Knowledge base CLI using similar patterns
- [reveal.js](https://revealjs.com/) - The presentation framework we generate for
- [BAML](https://www.boundaryml.com/) - Structured AI interactions
- [bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
