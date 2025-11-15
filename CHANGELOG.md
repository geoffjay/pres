# Changelog

All notable changes to this project will be documented in this file.

## [0.6.0] - 2025-11-14

### Changed
- **Migrated to agar library**: Replaced local `pkg/tui` with external dependency `github.com/geoffjay/agar/tui`
  - Removed local `pkg/tui/` directory
  - Added `github.com/geoffjay/agar` dependency to `go.mod`
  - Updated all imports from `github.com/geoffjay/pres/pkg/tui` to `github.com/geoffjay/agar/tui`
  - Updated examples and documentation

### Benefits
- **Shared library**: TUI components now available to all Go projects
- **Versioned dependency**: Proper semantic versioning and release management
- **Easier updates**: Pull latest components via `go get -u`
- **Reduced duplication**: One canonical source for TUI components

### Migration
For projects using the local `pkg/tui`:
```bash
# Add agar dependency
go get github.com/geoffjay/agar/tui

# Update imports
# From: import "yourproject/pkg/tui"
# To:   import "github.com/geoffjay/agar/tui"
```

## [0.5.0] - 2025-11-13

### Changed
- **Renamed to Input Components**: Renamed all "question" terminology to "input" for better genericity
  - Files: `question_*.go` → `input_*.go`
  - Functions: `NewYesNoQuestion` → `NewYesNoInput`, etc.
  - Internal fields: `question` → `prompt`
  - Examples: `question_components.go` → `input_components.go`
  - Binary: `question_demo` → `input_demo`
- Updated all documentation to reflect input-centric terminology

### Why?
The components are generic input collectors, not specifically "questions". The new naming makes them more reusable across different contexts (forms, wizards, configuration, etc.)

## [0.4.0] - 2025-11-13

### Added
- **Reusable Input Components**: Created three new TUI components for building interactive CLI apps
  - `YesNoModel` - Yes/No input with checkbox-style selection and keyboard shortcuts
  - `TextModel` - Single-line and multi-line text input with automatic wrapping
  - `OptionsModel` - Multiple choice input with radio button-style interface
- **Component Documentation**: Comprehensive README in `pkg/tui/` with usage examples
- **Example Program**: `examples/input_components.go` demonstrating all component types
- **Makefile**: Comprehensive build automation with 20+ targets
  - Build commands: `build`, `examples`, `build-all`, `rebuild`
  - Development: `dev`, `watch`, `run-examples`
  - Testing: `test`, `test-coverage`, `check`, `fmt`, `lint`
  - Maintenance: `clean`, `clean-all`, `deps`, `update-deps`
  - Code generation: `baml`, `tidy`
  - Installation: `install`
  - Colored output and helpful documentation
  - See `make help` for all available targets

### Features

Each component:
- Implements `tea.Model` interface for bubbletea compatibility
- Supports keyboard shortcuts (arrow keys, vim keys, quick keys)
- Provides consistent UX with shared styling
- Handles terminal resize automatically
- Can be used standalone or composed into larger forms

**Question Types:**
1. **Yes/No**: Navigate with arrows/vim keys, quick answer with y/n
2. **Text Single-line**: Enter to submit, automatic wrapping
3. **Text Multi-line**: Enter for new line, Ctrl+D to submit, preserves formatting
4. **Options**: Navigate with arrows/vim keys, number keys (1-9) for quick select

See `pkg/tui/README.md` for detailed documentation.

## [0.3.0] - 2025-11-13

### Fixed
- **Text Input Wrapping**: Fixed text overflow issue in TUI input fields
  - Text now wraps properly to the next line when typing long responses
  - Added dynamic terminal width tracking via `tea.WindowSizeMsg`
  - Implemented word-boundary wrapping for better readability
  - Added `renderWrappedInput()` and `wrapLineAtWords()` helper functions
  - Matches behavior from `contrib/kb` project

### Technical Details

#### Changes to `pkg/tui/iterative_form.go`:
- Added `width` field to `IterativeFormModel` to track terminal width
- Handle `tea.WindowSizeMsg` to dynamically update width
- Replaced simple input rendering with `renderWrappedInput()` calls
- Implemented intelligent word-wrapping at maxWidth - 6 characters

#### Before:
```go
b.WriteString(InputStyle.Render("> " + m.input + "█"))
```

#### After:
```go
b.WriteString(renderWrappedInput(m.input, m.width))
```

The wrapping function:
- Splits input by newlines (preserves user's line breaks)
- Wraps long lines at word boundaries
- Indents continuation lines with proper spacing
- Places cursor at the end of wrapped text

## [0.2.0] - 2025-11-13

### Fixed

- **Directory Creation Bug**: Fixed issue where `presentations/` directory was not created automatically

  - Now creates full directory path for output files
  - Changed from creating only base directory to creating complete path using `filepath.Dir()`

- **Format Compatibility**: Added support for both JSON formats
  - Can now load raw BAML `Presentation` format (what AI generates)
  - Can also load wrapped `PresentationData` format (what SavePresentation creates)
  - Automatically detects and converts between formats

### Technical Details

#### Before:

```go
// Only created base directory
if err := os.MkdirAll(w.baseDir, 0755); err != nil {
    return "", fmt.Errorf("failed to create directory: %w", err)
}
```

#### After:

```go
// Creates full directory path including subdirectories
fullPath := filepath.Join(w.baseDir, filename)
dir := filepath.Dir(fullPath)
if err := os.MkdirAll(dir, 0755); err != nil {
    return "", fmt.Errorf("failed to create directory: %w", err)
}
```

#### Format Handling:

- `LoadPresentation` now tries wrapped format first
- Falls back to raw BAML format if metadata is empty
- Automatically converts raw format to wrapped format internally

## [0.1.0] - 2025-11-13

### Added

- Initial release with complete feature set
- `pres create` - Create presentations with iterative Q&A
- `pres update` - Update presentations with natural language
- `pres generate` - Generate reveal.js HTML output
- Shared TUI library (`pkg/tui`) for iterative forms
- BAML integration for AI-powered generation
- Support for all reveal.js themes
- Speaker notes in presentations
- Multiple slide layouts (title, content, two-column, blank)

### Architecture

- BAML functions for presentation generation
- Reusable TUI components with bubbletea
- JSON storage format
- HTML generation with reveal.js
- Confidence-based iterative Q&A system
