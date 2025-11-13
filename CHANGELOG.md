# Changelog

All notable changes to this project will be documented in this file.

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
