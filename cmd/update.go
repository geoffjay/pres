package cmd

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/geoffjay/pres/baml_client"
	"github.com/geoffjay/pres/internal/presentation"
	"github.com/geoffjay/pres/pkg/tui"
	"github.com/spf13/cobra"
)

var (
	updatePath string
)

var updateCmd = &cobra.Command{
	Use:   "update [request]",
	Short: "Update an existing presentation",
	Long: `Update an existing presentation using an interactive Q&A process.

The command will:
1. Load the existing presentation
2. Gather contextual information about the changes
3. Apply updates to the presentation
4. Save the modified presentation

Examples:
  pres update --path presentations/my-talk.json "Add a slide at the beginning with an executive summary"
  pres update --path presentations/review.json "Change the theme to 'night'"
  pres update --path presentations/intro.json "Add more details to the goroutines slide"`,
	Args: cobra.ExactArgs(1),
	RunE: runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updatePath, "path", "p", "", "Path to presentation JSON file (required)")
	updateCmd.MarkFlagRequired("path")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	request := args[0]
	ctx := context.Background()

	fmt.Printf("🔄 Updating presentation: %s\n", updatePath)
	fmt.Printf("Request: %s\n\n", request)

	// Load existing presentation
	writer := presentation.NewWriter(".")
	existingData, err := writer.LoadPresentation(updatePath)
	if err != nil {
		return fmt.Errorf("failed to load presentation: %w", err)
	}

	fmt.Printf("Loaded: %s (%d slides)\n\n", existingData.Metadata.Title, len(existingData.Slides))

	// Generate presentation summary for context
	presentationSummary := existingData.GetSummary()

	const maxIterations = 3
	var allQAResponses []string

	// Iterative information gathering
	config := tui.IterationConfig{
		MaxIterations:    maxIterations,
		IterationPrompt:  "Gathering update context...",
		CompletionPrompt: "Do you need to provide more details about the update?",
	}

	form := tui.NewIterativeForm("Presentation Update", config)

	for iteration := 0; iteration < maxIterations; iteration++ {
		fmt.Printf("Preparing questions (iteration %d/%d)...\n", iteration+1, maxIterations)

		// Prepare questions using BAML
		preparation, err := baml_client.PrepareUpdatePresentation(ctx, request, presentationSummary, int64(iteration), allQAResponses)
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
			return fmt.Errorf("update cancelled")
		}

		// Collect responses
		iterationResponses := form.GetResponsesForIteration(iteration)
		for i, q := range preparation.Questions {
			if i < len(iterationResponses) {
				qa := fmt.Sprintf("Q: %s\nA: %s", q.Question, iterationResponses[i])
				allQAResponses = append(allQAResponses, qa)
			}
		}

		if !preparation.Needs_more_info {
			fmt.Printf("\n✓ Sufficient information gathered (confidence: %.2f)\n", preparation.Confidence_score)
			break
		}

		if iteration == maxIterations-1 {
			fmt.Println("\n⚠ Reached maximum iterations. Proceeding with available information...")
			break
		}

		if !form.NeedsMoreInfo() {
			break
		}

		form.NextIteration()
	}

	fmt.Println("\nGenerating update operations...")

	// Generate update operations
	updates, err := baml_client.GenerateUpdateOperations(ctx, request, presentationSummary, allQAResponses)
	if err != nil {
		return fmt.Errorf("failed to generate updates: %w", err)
	}

	if len(updates) == 0 {
		fmt.Println("⚠ No updates generated. Please try being more specific in your request.")
		return nil
	}

	// Display planned updates
	fmt.Printf("\nPlanned updates:\n")
	for i, update := range updates {
		fmt.Printf("  %d. %s: %s\n", i+1, update.Operation, update.Rationale)
	}

	// Apply updates
	fmt.Println("\nApplying updates...")
	if err := writer.UpdatePresentation(updatePath, updates); err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	// Reload to show summary
	updatedData, err := writer.LoadPresentation(updatePath)
	if err != nil {
		return fmt.Errorf("failed to reload presentation: %w", err)
	}

	fmt.Printf("\n✓ Presentation updated successfully!\n")
	fmt.Printf("  Location: %s\n", updatePath)
	fmt.Printf("  Slides: %d\n", len(updatedData.Slides))
	fmt.Printf("  Modified: %s\n", updatedData.Metadata.Modified.Format("2006-01-02 15:04:05"))

	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  • Review the changes: cat %s\n", updatePath)
	fmt.Printf("  • Generate HTML: pres generate --path %s\n", updatePath)
	fmt.Printf("  • Make more updates: pres update --path %s \"your next request\"\n", updatePath)

	return nil
}
