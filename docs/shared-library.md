# Shared Library Documentation

This document describes the shared components that can be abstracted out for reuse between `pres` and `kb` (and other CLI projects).

## Overview

Both `pres` and `kb` use similar patterns:

1. Iterative Q&A with confidence scoring
2. BAML-powered AI interactions
3. bubbletea TUI for user input
4. Structured data storage

The shared library (`agar/tui`) provides reusable components for these patterns.

## Shared TUI Components

### `agar/tui/iterative_form.go`

The `IterativeFormModel` provides a complete iterative Q&A system with:

- Multiple iterations (configurable)
- Confidence-based progression
- Response collection and display
- User-controlled iteration advancement

**Example Usage:**

```go
import "github.com/geoffjay/agar/tui"

config := tui.IterationConfig{
    MaxIterations:    3,
    IterationPrompt:  "Gathering context...",
    CompletionPrompt: "Do you want to provide more information?",
}

form := tui.NewIterativeForm("My Interactive Form", config)

// Add questions for current iteration
questions := []tui.IterativeQuestion{
    {
        Question:  "What is your goal?",
        HelpText:  "Describe what you want to achieve",
        Iteration: 0,
    },
}
form.AddQuestions(questions)

// Run the TUI
p := tea.NewProgram(form)
finalModel, _ := p.Run()

// Get responses
responses := finalModel.(tui.IterativeFormModel).GetResponses()
```

### Shared Styles

The TUI package exports standard lipgloss styles:

```go
tui.TitleStyle    // Bold, prominent titles
tui.QuestionStyle // Bold, colored questions
tui.HelpStyle     // Italic, subdued help text
tui.InputStyle    // Colored user input
tui.ErrorStyle    // Bold, red errors
tui.SuccessStyle  // Bold, green success messages
```

## Integration with BAML

### Typical Pattern

Both projects follow this pattern:

1. **Preparation Function** - BAML function generates questions

   ```go
   preparation, err := baml_client.PrepareXXX(ctx, input, iteration, previousResponses)
   ```

2. **Convert to TUI Questions**

   ```go
   var questions []tui.IterativeQuestion
   for _, q := range preparation.Questions {
       questions = append(questions, tui.IterativeQuestion{
           Question:  q.Question,
           HelpText:  q.Help_text,
           Iteration: int(q.Iteration),
       })
   }
   ```

3. **Run Interactive Form**

   ```go
   form.AddQuestions(questions)
   p := tea.NewProgram(form)
   finalModel, _ := p.Run()
   ```

4. **Check Confidence & Iterate**

   ```go
   if !preparation.Needs_more_info {
       break // Sufficient info
   }
   if form.NeedsMoreInfo() {
       form.NextIteration()
   }
   ```

5. **Generation Function** - BAML function generates output
   ```go
   result, err := baml_client.GenerateXXX(ctx, input, allResponses)
   ```

## BAML Function Patterns

### Preparation Functions

Preparation functions should return:

```baml
class PreparationResult {
  questions Question[] @description("2-5 questions to ask")
  rationale string @description("Why these questions help")
  confidence_score float @description("Confidence 0.0-1.0")
  confidence_reasoning string @description("Why this score")
  needs_more_info bool @description("Whether to iterate again")
}
```

### Question Structure

```baml
class Question {
  question string @description("The question text")
  help_text string @description("Optional help text")
  iteration int @description("Which iteration")
}
```

### Generation Functions

Generation functions should:

- Accept `qa_responses: string[]` containing formatted Q&A pairs
- Return structured output (e.g., `Presentation`, `JournalEntry`)
- Use confidence scoring from preparation phase

## File Organization

```
project/
├── baml_src/
│   ├── clients.baml          # Shared AI clients
│   ├── generators.baml       # BAML generator config
│   └── domain.baml           # Domain-specific functions
├── pkg/
│   └── tui/
│       └── iterative_form.go # Shared TUI components
├── internal/
│   └── domain/               # Domain-specific logic
│       ├── writer.go         # Storage operations
│       └── generator.go      # Output generation
└── cmd/
    └── commands.go           # CLI commands
```

## Creating a New Project with Shared Components

### 1. Set up BAML

```bash
# Copy shared BAML files
cp ../pres/baml_src/clients.baml baml_src/
cp ../pres/baml_src/generators.baml baml_src/

# Create domain-specific BAML
vim baml_src/my_domain.baml
```

### 2. Add agar Dependency

```bash
# Add the agar library
go get github.com/geoffjay/agar/tui
```

### 3. Create Domain Logic

```go
// internal/mydomain/writer.go
type Writer struct {
    baseDir string
}

func (w *Writer) Save(data MyData, path string) error {
    // Storage logic
}
```

### 4. Implement Commands

```go
// cmd/create.go
func runCreate(cmd *cobra.Command, args []string) error {
    // Use shared TUI pattern
    config := tui.IterationConfig{...}
    form := tui.NewIterativeForm("Create", config)

    // Iterate with BAML
    for iteration := 0; iteration < maxIterations; iteration++ {
        prep, _ := baml_client.PrepareCreate(...)

        // Convert and add questions
        form.AddQuestions(convertQuestions(prep.Questions))

        // Run TUI
        p := tea.NewProgram(form)
        finalModel, _ := p.Run()

        // Check if done
        if !prep.Needs_more_info {
            break
        }
    }

    // Generate output
    result, _ := baml_client.GenerateOutput(...)

    return nil
}
```

## Benefits of Shared Library

1. **Consistency** - Same UX across projects
2. **Reusability** - Write iterative Q&A once, use everywhere
3. **Maintainability** - Fix bugs in one place
4. **Extensibility** - Easy to add new question types or validation
5. **Best Practices** - Confidence scoring, iteration management built-in

## Future Enhancements

### Additional TUI Components

Consider adding:

- **ValidationForm** - Form with custom validators
- **SelectForm** - Multiple-choice questions
- **NumberForm** - Numeric input with ranges
- **MultilineForm** - Long-form text input

### BAML Utilities

Shared BAML utilities:

```baml
// baml_src/shared_utilities.baml

// Generic confidence scoring
function AssessConfidence(
  domain: string,
  information_gathered: string[],
  required_fields: string[]
) -> ConfidenceAssessment {
  // Generic confidence logic
}

// Generic question generation
function GenerateQuestions(
  context: string,
  iteration: int,
  focus_areas: string[]
) -> Question[] {
  // Generic question logic
}
```

### Configuration Abstraction

```go
// pkg/config/iterative.go
type IterativeConfig struct {
    MaxIterations    int
    MinConfidence    float64
    DefaultClient    string
    Prompts          PromptConfig
}

func LoadConfig(path string) (*IterativeConfig, error) {
    // Load from YAML/JSON
}
```

### Storage Abstraction

```go
// pkg/storage/interface.go
type Storage interface {
    Save(id string, data interface{}) error
    Load(id string) (interface{}, error)
    Update(id string, updates interface{}) error
    List(filter Filter) ([]string, error)
}

// Implementations for different backends
type JSONStorage struct { ... }
type SQLiteStorage struct { ... }
```

## Migration Guide

### Migrating `kb` to Use Shared Library

1. Copy `agar/tui` from `pres`
2. Update imports:

   ```go
   // Old
   import "github.com/geoffjay/kb/internal/tui"

   // New
   import "github.com/geoffjay/pres/agar/tui"
   ```

3. Update style references:

   ```go
   // Old
   titleStyle := lipgloss.NewStyle()...

   // New
   tui.TitleStyle
   ```

4. Keep domain-specific logic in `internal/tui` if needed

### Publishing as Standalone Module

To use across multiple projects:

1. Extract to separate repository:

   ```bash
   git init go-iterative-tui
   cp -r agar/tui/* .
   go mod init github.com/yourname/go-iterative-tui
   ```

2. Use in projects:

   ```bash
   go get github.com/yourname/go-iterative-tui
   ```

3. Import:
   ```go
   import tui "github.com/yourname/go-iterative-tui"
   ```

## Examples

See:

- `pres/cmd/create.go` - Presentation creation with iterative Q&A
- `kb/cmd/journal.go` - Journal entry with iterative Q&A
- `kb/cmd/research.go` - External research with iterative Q&A

All follow the same shared pattern.
