// Presentation Generation Functions
// These functions help create, update, and generate presentations using reveal.js

// ============================================================================
// DATA MODELS
// ============================================================================

// Represents a single slide in a presentation
class Slide {
  title string @description("Slide title, can be empty for title slides")
  content string @description("Markdown content for the slide")
  notes string @description("Speaker notes for the slide")
  layout string @description("Layout type: title, content, two-column, or blank")
  background_color string @description("Optional background color (e.g., #1a1a1a)")
}

// Represents a complete presentation
class Presentation {
  title string @description("Presentation title")
  subtitle string @description("Presentation subtitle")
  author string @description("Author name")
  date string @description("Presentation date")
  theme string @description("reveal.js theme: black, white, league, beige, sky, night, serif, simple, solarized")
  slides Slide[] @description("Array of slides in the presentation")
  tags string[] @description("Tags for categorization")
}

// Represents contextual questions for gathering information
class PresentationQuestion {
  question string @description("The question to ask the user")
  help_text string @description("Optional help text explaining the question")
  iteration int @description("Which iteration this question belongs to")
}

// Represents the preparation phase for creating/updating a presentation
class PresentationPreparation {
  questions PresentationQuestion[] @description("3-5 questions to gather context")
  rationale string @description("Why these questions will help create a better presentation")
  confidence_score float @description("Confidence that we have enough information (0.0-1.0)")
  confidence_reasoning string @description("Why this confidence score was assigned")
  needs_more_info bool @description("Whether another iteration is recommended")
}

// Represents an update operation on an existing presentation
class PresentationUpdate {
  operation string @description("Type of update: add_slide, modify_slide, delete_slide, reorder_slides, update_metadata")
  slide_index int @description("Index of slide to modify/delete (0-based), -1 for add/reorder/metadata operations")
  new_slide Slide @description("New slide content for add/modify operations")
  new_order int[] @description("New slide order for reorder operation (array of indices)")
  metadata_updates map<string, string> @description("Metadata updates for update_metadata operation")
  rationale string @description("Explanation of the update")
}

// ============================================================================
// PRESENTATION CREATION
// ============================================================================

// Prepare questions to gather context for creating a presentation
function PrepareCreatePresentation(
  description: string,
  iteration: int,
  previous_responses: string[]
) -> PresentationPreparation {
  client CustomHaiku
  prompt #"
    You are helping create a presentation by gathering contextual information.

    Presentation description: {{ description }}
    Current iteration: {{ iteration }}
    Max iterations: 3

    {% if previous_responses %}
    Previous responses from user:
    {{ previous_responses }}
    {% endif %}

    Generate 3-5 thoughtful questions that will help gather the information needed
    to create an effective presentation.

    Iteration focus:
    - Iteration 0: Audience, purpose, key message, desired outcome
    - Iteration 1: Main topics, structure, level of detail, time constraints
    - Iteration 2: Visual preferences, specific examples, supporting data

    Questions should:
    1. Build on previous responses when provided
    2. Gather specific information about audience and context
    3. Understand the key message and takeaways
    4. Identify the structure and flow
    5. Determine appropriate depth and complexity
    6. NOT be redundant with previous iterations

    After generating questions, assign a confidence score (0.0-1.0):
    - 0.0-0.4: Need much more information
    - 0.4-0.8: Have basic info, more details would help
    - 0.8-1.0: Have sufficient information to create presentation

    Consider:
    - Do we understand the audience and their needs?
    - Is the main message and structure clear?
    - Do we have enough detail to create meaningful slides?
    - Are there gaps that would make the presentation generic?

    Set needs_more_info to true if confidence < 0.8 OR if this is iteration 0 or 1.
    Set needs_more_info to false if confidence >= 0.8 AND iteration >= 2.

    {{ ctx.output_format }}
  "#
}

// Generate a complete presentation from user responses
function GeneratePresentation(
  description: string,
  qa_responses: string[],
  today_date: string
) -> Presentation {
  client AnthropicFallback
  prompt #"
    You are creating a reveal.js presentation based on user-provided information.

    IMPORTANT: Today's date is {{ today_date }}.

    Presentation description: {{ description }}

    User's responses to contextual questions:
    {{ qa_responses }}

    Generate a complete, well-structured presentation that:
    - Creates an engaging title and subtitle
    - Includes a title slide with author and date
    - Organizes content into logical, focused slides
    - Uses appropriate slide layouts (title, content, two-column)
    - Keeps each slide focused and not overwhelming (3-5 points max per slide)
    - Uses markdown formatting effectively (lists, emphasis, code blocks)
    - Includes speaker notes with additional context
    - Follows presentation best practices:
      * One main idea per slide
      * Clear visual hierarchy
      * Concise bullet points
      * Smooth narrative flow
    - Chooses an appropriate reveal.js theme
    - Suggests relevant tags for categorization

    Available reveal.js themes:
    - black: Dark background, white text (modern, professional)
    - white: White background, dark text (clean, minimal)
    - league: Gray background (neutral, versatile)
    - beige: Beige background (warm, approachable)
    - sky: Sky blue background (calm, friendly)
    - night: Black background with orange highlights (bold, energetic)
    - serif: Serif fonts (classic, formal)
    - simple: Simple and minimal (understated)
    - solarized: Solarized colors (eye-friendly, technical)

    Slide layouts:
    - title: For section introductions (large centered text)
    - content: Standard content slide with title and bullet points
    - two-column: Split content into two columns
    - blank: Minimal slide for images or quotes

    Use ONLY the information provided by the user. Create 8-15 slides for a
    complete presentation. Format slide content in markdown.

    {{ ctx.output_format }}
  "#
}

// ============================================================================
// PRESENTATION UPDATES
// ============================================================================

// Prepare questions to gather context for updating a presentation
function PrepareUpdatePresentation(
  update_request: string,
  current_presentation: string,
  iteration: int,
  previous_responses: string[]
) -> PresentationPreparation {
  client CustomHaiku
  prompt #"
    You are helping update an existing presentation by gathering contextual information.

    Update request: {{ update_request }}
    Current iteration: {{ iteration }}
    Max iterations: 3

    Current presentation summary:
    {{ current_presentation }}

    {% if previous_responses %}
    Previous responses from user:
    {{ previous_responses }}
    {% endif %}

    Generate 2-4 thoughtful questions that will help understand exactly what
    changes the user wants to make.

    Iteration focus:
    - Iteration 0: What specifically to change, where in the presentation, why
    - Iteration 1: Specific content details, placement preferences
    - Iteration 2: Visual preferences, final clarifications

    Questions should:
    1. Build on previous responses
    2. Clarify the specific changes needed
    3. Understand the rationale for changes
    4. Determine placement and structure
    5. NOT be redundant with previous iterations

    Confidence scoring (0.0-1.0):
    - 0.0-0.4: Don't understand what to change yet
    - 0.4-0.8: Have general idea, need specific details
    - 0.8-1.0: Clear on exactly what changes to make

    {{ ctx.output_format }}
  "#
}

// Generate update operations for an existing presentation
function GenerateUpdateOperations(
  update_request: string,
  current_presentation: string,
  qa_responses: string[]
) -> PresentationUpdate[] {
  client AnthropicFallback
  prompt #"
    You are updating an existing presentation based on user requests.

    Update request: {{ update_request }}

    Current presentation:
    {{ current_presentation }}

    User's responses to contextual questions:
    {{ qa_responses }}

    Generate the specific update operations needed to fulfill the user's request.

    Available operations:
    - add_slide: Add a new slide at a specific position
      * Set slide_index to where to insert (0 = beginning)
      * Provide complete new_slide content
    - modify_slide: Change content of an existing slide
      * Set slide_index to the slide to modify
      * Provide updated new_slide content
    - delete_slide: Remove a slide
      * Set slide_index to the slide to remove
    - reorder_slides: Change slide order
      * Provide new_order array with reordered indices
    - update_metadata: Change presentation title, author, theme, etc.
      * Provide metadata_updates map with key-value changes

    Guidelines:
    - Make minimal, focused changes to address the request
    - Maintain the presentation's overall structure and flow
    - Ensure slide indices are correct (0-based)
    - Provide clear rationale for each operation
    - If adding multiple slides, create separate operations for each
    - When modifying slides, preserve good formatting and structure

    Return an array of operations to apply in sequence.

    {{ ctx.output_format }}
  "#
}

// ============================================================================
// TESTS
// ============================================================================

test prepare_create_iter0 {
  functions [PrepareCreatePresentation]
  args {
    description "Introduction to Go concurrency patterns"
    iteration 0
    previous_responses []
  }
}

test prepare_create_iter1 {
  functions [PrepareCreatePresentation]
  args {
    description "Introduction to Go concurrency patterns"
    iteration 1
    previous_responses [
      "Q: Who is your target audience?\nA: Intermediate Go developers who are new to concurrency",
      "Q: What's the main goal of this presentation?\nA: Help them understand goroutines, channels, and common patterns",
      "Q: How long should the presentation be?\nA: About 30 minutes with examples"
    ]
  }
}

test generate_presentation {
  functions [GeneratePresentation]
  args {
    description "Introduction to Go concurrency patterns"
    qa_responses [
      "Q: Who is your target audience?\nA: Intermediate Go developers new to concurrency",
      "Q: What's the main goal?\nA: Understand goroutines, channels, and patterns",
      "Q: How long?\nA: 30 minutes with examples",
      "Q: What level of depth?\nA: Practical examples, not too theoretical",
      "Q: Any specific patterns to cover?\nA: Worker pools, fan-out/fan-in, pipelines"
    ]
    today_date "2025-01-15"
  }
}

test prepare_update_iter0 {
  functions [PrepareUpdatePresentation]
  args {
    update_request "Add a slide at the beginning with an executive summary"
    current_presentation #"
      Title: Introduction to Go Concurrency
      Slides: 12
      Topics: Goroutines, Channels, Select, Patterns
    "#
    iteration 0
    previous_responses []
  }
}

test generate_updates {
  functions [GenerateUpdateOperations]
  args {
    update_request "Add an executive summary at the beginning and a Q&A slide at the end"
    current_presentation #"
      Title: Introduction to Go Concurrency
      Author: John Doe
      Theme: black
      Slides:
      1. Title slide
      2. What is concurrency?
      3. Goroutines basics
      ...
      12. Conclusion
    "#
    qa_responses [
      "Q: What should the executive summary include?\nA: Key takeaways, who should attend, time estimate",
      "Q: What about the Q&A slide?\nA: Just a simple slide inviting questions"
    ]
  }
}
