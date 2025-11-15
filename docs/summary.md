# Project Summary

## What We Built

A complete AI-powered presentation generator CLI (`pres`) that creates, updates, and generates reveal.js presentations through intelligent iterative Q&A.

## Architecture Overview

```
pres/
├── baml_src/
│   ├── presentations.baml    # AI functions for presentation generation
│   ├── clients.baml          # AI model configurations
│   └── generators.baml       # BAML code generation config
├── internal/
│   └── presentation/
│       ├── writer.go         # JSON storage and updates
│       └── generator.go      # HTML generation for reveal.js
├── cmd/
│   ├── root.go              # Root command
│   ├── create.go            # Create presentations
│   ├── update.go            # Update presentations
│   └── generate.go          # Generate HTML
├── examples/
│   └── input_components.go  # Example usage of agar/tui
└── baml_client/             # Generated BAML client code

Dependencies:
├── github.com/geoffjay/agar/tui  # Reusable TUI input components (SHARED)
└── github.com/boundaryml/baml    # Structured AI interactions
```

## Key Features

### 1. Iterative Q&A with Confidence Scoring

The AI asks questions iteratively, assessing confidence at each step:

```
Iteration 0: Audience, purpose, key message (confidence: 0.4)
→ More questions needed

Iteration 1: Structure, topics, depth (confidence: 0.7)
→ More questions needed

Iteration 2: Examples, visual preferences (confidence: 0.9)
→ Sufficient information! ✓
```

### 2. Natural Language Updates

```bash
# Instead of manual JSON editing:
pres update --path my-talk.json "Add an executive summary at the beginning"

# AI understands and applies changes intelligently
```

### 3. Professional HTML Output

Generates production-ready reveal.js presentations:

- Multiple themes
- Speaker notes
- Markdown content
- Two-column layouts
- Syntax highlighting

## Comparison: pres vs kb

Both projects share the same architectural patterns:

| Feature                | pres                  | kb                        |
| ---------------------- | --------------------- | ------------------------- |
| **Domain**             | Presentations         | Knowledge Base            |
| **Storage**            | JSON                  | Markdown + TOML           |
| **Output**             | reveal.js HTML        | Zola static site          |
| **Shared TUI**         | ✅ `pkg/tui`          | ✅ Can use `pkg/tui`      |
| **BAML Functions**     | Presentation-specific | Research/Journal-specific |
| **Iterative Q&A**      | ✅                    | ✅                        |
| **Confidence Scoring** | ✅                    | ✅                        |
| **Update Operations**  | Slide modifications   | N/A                       |

## BAML Functions

### Presentation Functions

1. **PrepareCreatePresentation** - Generate questions for creating presentations
2. **GeneratePresentation** - Create slides from Q&A responses
3. **PrepareUpdatePresentation** - Generate questions for updates
4. **GenerateUpdateOperations** - Create update operations

### KB Functions (for comparison)

1. **PrepareJournalEntry** - Generate questions for journal entries
2. **GenerateJournalEntry** - Create journal from Q&A
3. **PrepareExternalResearch** - Generate research questions
4. **ResearchExternal** - Conduct research from web sources

## Workflow Examples

### Create → Generate → Present

```bash
# Step 1: Create presentation
./pres create "Go Concurrency Workshop"
# Answers 3-4 iterations of questions
# Saves to: presentations/go-concurrency-workshop.json

# Step 2: Generate HTML
./pres generate --path presentations/go-concurrency-workshop.json
# Creates: presentations/go-concurrency-workshop.html

# Step 3: Present
open presentations/go-concurrency-workshop.html
```

### Create → Update → Generate

```bash
# Create initial presentation
./pres create "Product Roadmap Q1"

# Make updates
./pres update --path presentations/product-roadmap-q1.json \
    "Add a risk assessment slide before the conclusion"

./pres update --path presentations/product-roadmap-q1.json \
    "Change theme to 'night'"

# Generate final version
./pres generate --path presentations/product-roadmap-q1.json
```

## Shared Library Benefits

The `pkg/tui` package can be used by both projects:

### Before (Duplicated Code)

```
pres/internal/tui/iterative_form.go  (500 lines)
kb/internal/tui/iterative_form.go    (500 lines)
Total: 1000 lines, 2x maintenance
```

### After (Shared Library)

```
pres/pkg/tui/iterative_form.go       (350 lines, refined)
pres/cmd/* → uses pkg/tui
kb/cmd/* → uses pkg/tui
Total: 350 lines, 1x maintenance
```

## Integration Pattern

All commands follow this pattern:

```go
func runCommand(cmd *cobra.Command, args []string) error {
    // 1. Setup
    config := tui.IterationConfig{MaxIterations: 3}
    form := tui.NewIterativeForm("Title", config)

    // 2. Iterative Q&A
    for iteration := 0; iteration < maxIterations; iteration++ {
        // Ask AI to prepare questions
        prep, _ := baml_client.PrepareXXX(ctx, ...)

        // Add questions to form
        form.AddQuestions(convertQuestions(prep.Questions))

        // Run interactive TUI
        p := tea.NewProgram(form)
        finalModel, _ := p.Run()

        // Check confidence
        if !prep.Needs_more_info { break }
        if !form.NeedsMoreInfo() { break }

        form.NextIteration()
    }

    // 3. Generate output
    result, _ := baml_client.GenerateXXX(ctx, ...)

    // 4. Save
    writer.Save(result, path)

    return nil
}
```

## Next Steps

### For This Project

1. Add more slide layouts (e.g., image, quote, code)
2. Support importing from Markdown files
3. Add presentation templates
4. Implement slide transitions and animations
5. Add PDF export support

### For Shared Library

1. Extract `pkg/tui` to standalone Go module
2. Add more validation types (number, email, URL)
3. Support multi-select questions
4. Add progress indicators for long operations
5. Create plugin system for custom question types

### For Both Projects

1. Share BAML utility functions (confidence assessment, question generation)
2. Create shared storage interfaces
3. Build common CLI patterns (progress bars, error handling)
4. Shared configuration management
5. Common testing utilities

## File Statistics

```
Language                 Files       Lines       Code
────────────────────────────────────────────────────
Go                          10        1200        950
BAML                         3         400        350
Markdown                     3         500        500
────────────────────────────────────────────────────
Total                       16        2100       1800
```

## Dependencies

- **github.com/spf13/cobra** - CLI framework
- **github.com/charmbracelet/bubbletea** - TUI framework
- **github.com/charmbracelet/lipgloss** - TUI styling
- **github.com/boundaryml/baml** - Structured AI interactions
- **reveal.js** (via CDN) - Presentation framework

## Design Principles

1. **Iterative over exhaustive** - Ask questions progressively, not all at once
2. **Confidence over completion** - Stop when AI is confident, not after N iterations
3. **Structured over freeform** - Use BAML for structured AI output
4. **Reusable over custom** - Build shared components for common patterns
5. **Interactive over declarative** - Guide users through creation, don't require upfront knowledge

## Similar Patterns in kb Project

The `kb` project (in `contrib/kb`) uses identical patterns for:

- **Journal creation** - Iterative Q&A for daily journal entries
- **External research** - Web search and synthesis with Q&A
- **Hybrid research** - Combining KB and web sources

Both projects demonstrate how the same architectural patterns can be applied to different domains while sharing core components.

## Learning Resources

- **BAML Documentation**: https://docs.boundaryml.com/
- **bubbletea Tutorial**: https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- **reveal.js Documentation**: https://revealjs.com/
- **Cobra User Guide**: https://github.com/spf13/cobra/blob/main/site/content/user_guide.md

## Conclusion

This project demonstrates:

1. How to build AI-powered CLI tools with structured interactions
2. How to create reusable TUI components for iterative workflows
3. How to use BAML for confidence-based information gathering
4. How to share code between similar projects

The patterns established here can be applied to many other domains:

- Document generation
- Content creation
- Data transformation
- Configuration management
- Interactive setup wizards

The key insight is that **iterative Q&A with confidence scoring** is a powerful pattern for gathering information from users, and the shared TUI library makes it easy to implement consistently across projects.
