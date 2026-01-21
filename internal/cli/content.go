package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var contentCmd = &cobra.Command{
	Use:   "content <url>",
	Short: "Extract content from a web page",
	Long: `Extract content (HTML, text, or markdown) from the specified URL.

Options are passed as key=value pairs using the -X flag.
See https://docs.capture.page/ for available options.

Examples:
  capture content https://example.com --format markdown
  capture content https://example.com --format html -o page.html
  capture content https://example.com --format text
  capture content https://example.com --json
  capture content https://example.com -X delay=1000 --format markdown

Output formats (--format flag):
  html      Raw HTML content
  text      Plain text content
  markdown  Markdown formatted content`,
	Args: cobra.ExactArgs(1),
	RunE: runContent,
}

var (
	contentOutput  string
	contentFormat  string
	contentJSON    bool
	contentOptions []string
)

func init() {
	rootCmd.AddCommand(contentCmd)

	contentCmd.Flags().StringVarP(&contentOutput, "output", "o", "", "Output file (default: stdout)")
	contentCmd.Flags().StringVar(&contentFormat, "format", "markdown", "Output format: html, text, markdown")
	contentCmd.Flags().BoolVar(&contentJSON, "json", false, "Output raw JSON response")
	contentCmd.Flags().StringArrayVarP(&contentOptions, "option", "X", nil, "API option as key=value (can be repeated)")
}

func runContent(cmd *cobra.Command, args []string) error {
	targetURL := args[0]

	opts, err := parseOptions(contentOptions)
	if err != nil {
		return err
	}

	client := newCaptureClient()

	if dryRun {
		url, err := client.BuildContentURL(targetURL, opts)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	}

	verboseLog("Extracting content from %s", targetURL)

	content, err := client.FetchContent(targetURL, opts)
	if err != nil {
		return fmt.Errorf("failed to extract content: %w", err)
	}

	if contentJSON {
		data, err := json.MarshalIndent(content, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		return writeOutput(data, contentOutput)
	}

	var output string
	switch contentFormat {
	case "html":
		output = content.HTML
	case "text":
		output = content.TextContent
	case "markdown":
		output = content.Markdown
	default:
		return fmt.Errorf("invalid format: %s (use html, text, or markdown)", contentFormat)
	}

	return writeStringOutput(output, contentOutput)
}
