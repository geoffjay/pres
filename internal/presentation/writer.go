package presentation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/geoffjay/pres/baml_client/types"
)

// PresentationData represents the stored presentation format
type PresentationData struct {
	Metadata struct {
		Title    string    `json:"title"`
		Subtitle string    `json:"subtitle"`
		Author   string    `json:"author"`
		Date     string    `json:"date"`
		Theme    string    `json:"theme"`
		Tags     []string  `json:"tags"`
		Created  time.Time `json:"created"`
		Modified time.Time `json:"modified"`
	} `json:"metadata"`
	Slides []types.Slide `json:"slides"`
}

// Writer handles writing presentations to disk
type Writer struct {
	baseDir string
}

// NewWriter creates a new presentation writer
func NewWriter(baseDir string) *Writer {
	return &Writer{baseDir: baseDir}
}

// SavePresentation saves a presentation to a JSON file
func (w *Writer) SavePresentation(pres *types.Presentation, filename string) (string, error) {
	// Ensure filename has .json extension
	if filepath.Ext(filename) != ".json" {
		filename = filename + ".json"
	}

	// Full path to the file
	fullPath := filepath.Join(w.baseDir, filename)

	// Create directory for the file if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Create presentation data structure
	data := PresentationData{}
	data.Metadata.Title = pres.Title
	data.Metadata.Subtitle = pres.Subtitle
	data.Metadata.Author = pres.Author
	data.Metadata.Date = pres.Date
	data.Metadata.Theme = pres.Theme
	data.Metadata.Tags = pres.Tags
	data.Metadata.Created = time.Now()
	data.Metadata.Modified = time.Now()
	data.Slides = pres.Slides

	// Marshal to JSON with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fullPath, nil
}

// LoadPresentation loads a presentation from a JSON file
func (w *Writer) LoadPresentation(path string) (*PresentationData, error) {
	// Read file
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Try to unmarshal as PresentationData first (wrapped format)
	var data PresentationData
	if err := json.Unmarshal(jsonData, &data); err == nil {
		// Check if this is the wrapped format by seeing if metadata is populated
		if data.Metadata.Title != "" {
			return &data, nil
		}
	}

	// If that failed or metadata is empty, try raw Presentation format (BAML output)
	var pres types.Presentation
	if err := json.Unmarshal(jsonData, &pres); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON (tried both formats): %w", err)
	}

	// Convert to PresentationData format
	data = PresentationData{}
	data.Metadata.Title = pres.Title
	data.Metadata.Subtitle = pres.Subtitle
	data.Metadata.Author = pres.Author
	data.Metadata.Date = pres.Date
	data.Metadata.Theme = pres.Theme
	data.Metadata.Tags = pres.Tags
	data.Metadata.Created = time.Now()
	data.Metadata.Modified = time.Now()
	data.Slides = pres.Slides

	return &data, nil
}

// UpdatePresentation applies updates to an existing presentation
func (w *Writer) UpdatePresentation(path string, updates []types.PresentationUpdate) error {
	// Load existing presentation
	data, err := w.LoadPresentation(path)
	if err != nil {
		return err
	}

	// Apply each update operation
	for _, update := range updates {
		switch update.Operation {
		case "add_slide":
			data.Slides = w.addSlide(data.Slides, update.Slide_index, update.New_slide)
		case "modify_slide":
			if update.Slide_index >= 0 && update.Slide_index < int64(len(data.Slides)) {
				data.Slides[update.Slide_index] = update.New_slide
			}
		case "delete_slide":
			if update.Slide_index >= 0 && update.Slide_index < int64(len(data.Slides)) {
				data.Slides = append(data.Slides[:update.Slide_index], data.Slides[update.Slide_index+1:]...)
			}
		case "reorder_slides":
			data.Slides = w.reorderSlides(data.Slides, update.New_order)
		case "update_metadata":
			w.updateMetadata(&data.Metadata, update.Metadata_updates)
		}
	}

	// Update modification time
	data.Metadata.Modified = time.Now()

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// addSlide inserts a slide at the specified index
func (w *Writer) addSlide(slides []types.Slide, index int64, newSlide types.Slide) []types.Slide {
	if index < 0 {
		index = 0
	}
	if index > int64(len(slides)) {
		index = int64(len(slides))
	}

	// Insert slide at index
	result := make([]types.Slide, 0, len(slides)+1)
	result = append(result, slides[:index]...)
	result = append(result, newSlide)
	result = append(result, slides[index:]...)

	return result
}

// reorderSlides reorders slides based on new order indices
func (w *Writer) reorderSlides(slides []types.Slide, newOrder []int64) []types.Slide {
	if len(newOrder) != len(slides) {
		return slides // Invalid order, return unchanged
	}

	result := make([]types.Slide, len(slides))
	for i, oldIdx := range newOrder {
		if oldIdx >= 0 && oldIdx < int64(len(slides)) {
			result[i] = slides[oldIdx]
		}
	}

	return result
}

// updateMetadata updates presentation metadata
func (w *Writer) updateMetadata(metadata *struct {
	Title    string    `json:"title"`
	Subtitle string    `json:"subtitle"`
	Author   string    `json:"author"`
	Date     string    `json:"date"`
	Theme    string    `json:"theme"`
	Tags     []string  `json:"tags"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}, updates map[string]string) {
	for key, value := range updates {
		switch key {
		case "title":
			metadata.Title = value
		case "subtitle":
			metadata.Subtitle = value
		case "author":
			metadata.Author = value
		case "date":
			metadata.Date = value
		case "theme":
			metadata.Theme = value
		}
	}
}

// GetPresentationSummary generates a text summary of the presentation
func (data *PresentationData) GetSummary() string {
	return fmt.Sprintf(`Title: %s
Subtitle: %s
Author: %s
Date: %s
Theme: %s
Tags: %v
Number of Slides: %d
Created: %s
Modified: %s`,
		data.Metadata.Title,
		data.Metadata.Subtitle,
		data.Metadata.Author,
		data.Metadata.Date,
		data.Metadata.Theme,
		data.Metadata.Tags,
		len(data.Slides),
		data.Metadata.Created.Format("2006-01-02 15:04:05"),
		data.Metadata.Modified.Format("2006-01-02 15:04:05"),
	)
}
