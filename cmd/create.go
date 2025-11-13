package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/geoffjay/pres/baml_client"
	"github.com/geoffjay/pres/internal/presentation"
	"github.com/geoffjay/pres/pkg/tui"
	"github.com/spf13/cobra"
)

var (
	createOutput string
	createAuthor string
)

var createCmd = &cobra.Command{
	Use:   "create [description]",
	Short: "Create a new presentation",
	Long: `Create a new presentation using an interactive Q&A process.

The command will:
1. Gather contextual information through questions
2. Generate presentation slides based on your responses
3. Save the presentation to a JSON file

Examples:
  pres create "Introduction to Go concurrency patterns"
  pres create "Q4 Business Review" --author "Jane Doe"
  pres create "Product Launch" --output presentations/launch.json`,
	Args: cobra.ExactArgs(1),
	RunE: runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createOutput, "output", "o", "", "Output path for presentation (default: generated from title)")
	createCmd.Flags().StringVar(&createAuthor, "author", "", "Author name (default: from environment or empty)")
}

func runCreate(cmd *cobra.Command, args []string) error {
	description := args[0]
	ctx := context.Background()

	fmt.Printf("📊 Creating presentation: %s\n\n", description)

	const maxIterations = 3
	var allQAResponses []string

	// Iterative information gathering with confidence scoring
	config := tui.IterationConfig{
		MaxIterations:    maxIterations,
		IterationPrompt:  "Gathering presentation context...",
		CompletionPrompt: "Do you want to provide more context for the presentation?",
	}

	form := tui.NewIterativeForm("Presentation Creation", config)

	for iteration := 0; iteration < maxIterations; iteration++ {
		fmt.Printf("Preparing questions (iteration %d/%d)...\n", iteration+1, maxIterations)

		// Prepare questions using BAML
		preparation, err := baml_client.PrepareCreatePresentation(ctx, description, int64(iteration), allQAResponses)
		if err != nil {
			return fmt.Errorf("failed to prepare questions: %w", err)
		}

		if len(preparation.Questions) == 0 {
			break
		}

		fmt.Printf("\n%s\n", preparation.Rationale)
		fmt.Printf("Confidence: %.2f/1.0 - %s\n\n", preparation.Confidence_score, preparation.Confidence_reasoning)

		// Convert BAML questions to TUI questions
		var questions []tui.IterativeQuestion
		for _, q := range preparation.Questions {
			questions = append(questions, tui.IterativeQuestion{
				Question:  q.Question,
				HelpText:  q.Help_text,
				Iteration: int(q.Iteration),
			})
		}

		form.AddQuestions(questions)

		// Run interactive TUI
		p := tea.NewProgram(form)
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("error running interactive form: %w", err)
		}

		form = finalModel.(tui.IterativeFormModel)

		if !form.IsDone() && !form.NeedsMoreInfo() {
			return fmt.Errorf("presentation creation cancelled")
		}

		// Collect responses from this iteration
		iterationResponses := form.GetResponsesForIteration(iteration)
		for i, q := range preparation.Questions {
			if i < len(iterationResponses) {
				qa := fmt.Sprintf("Q: %s\nA: %s", q.Question, iterationResponses[i])
				allQAResponses = append(allQAResponses, qa)
			}
		}

		// Check if we need more information based on AI confidence
		if !preparation.Needs_more_info {
			fmt.Printf("\n✓ Sufficient information gathered (confidence: %.2f)\n", preparation.Confidence_score)
			break
		}

		// If not enough info but at max iterations
		if iteration == maxIterations-1 {
			fmt.Println("\n⚠ Reached maximum iterations. Proceeding with available information...")
			break
		}

		if !form.NeedsMoreInfo() {
			break
		}

		form.NextIteration()
	}

	fmt.Println("\nGenerating presentation from your responses...")

	// Generate presentation from all Q&A
	today := time.Now().Format("2006-01-02")
	result, err := baml_client.GeneratePresentation(ctx, description, allQAResponses, today)
	if err != nil {
		return fmt.Errorf("failed to generate presentation: %w", err)
	}

	// Override author if provided
	if createAuthor != "" {
		result.Author = createAuthor
	}

	// Determine output path
	outputPath := createOutput
	if outputPath == "" {
		// Generate filename from title
		filename := strings.ToLower(result.Title)
		filename = strings.ReplaceAll(filename, " ", "-")
		// Remove non-alphanumeric characters except hyphens
		var cleanName strings.Builder
		for _, r := range filename {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
				cleanName.WriteRune(r)
			}
		}
		filename = cleanName.String()
		// Remove duplicate hyphens
		for strings.Contains(filename, "--") {
			filename = strings.ReplaceAll(filename, "--", "-")
		}
		filename = strings.Trim(filename, "-")
		outputPath = "presentations/" + filename + ".json"
	}

	// Save presentation
	writer := presentation.NewWriter(".")
	savedPath, err := writer.SavePresentation(&result, outputPath)
	if err != nil {
		return fmt.Errorf("failed to save presentation: %w", err)
	}

	// Display summary
	fmt.Printf("\n✓ Presentation created successfully!\n")
	fmt.Printf("  Location: %s\n", savedPath)
	fmt.Printf("  Title: %s\n", result.Title)
	if result.Subtitle != "" {
		fmt.Printf("  Subtitle: %s\n", result.Subtitle)
	}
	fmt.Printf("  Author: %s\n", result.Author)
	fmt.Printf("  Theme: %s\n", result.Theme)
	fmt.Printf("  Slides: %d\n", len(result.Slides))
	if len(result.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", strings.Join(result.Tags, ", "))
	}

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  • Review the presentation: cat %s\n", savedPath)
	fmt.Printf("  • Generate HTML: pres generate --path %s\n", savedPath)
	fmt.Printf("  • Update content: pres update --path %s \"your update request\"\n", savedPath)

	return nil
}
