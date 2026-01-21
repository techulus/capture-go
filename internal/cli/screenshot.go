package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot <url>",
	Short: "Take a screenshot of a web page",
	Long: `Take a screenshot of the specified URL.

Options are passed as key=value pairs using the -X flag.
See https://docs.capture.page/ for available options.

Examples:
  capture screenshot https://example.com -o screenshot.png
  capture screenshot https://example.com -X vw=1920 -X vh=1080 -o full.png
  capture screenshot https://example.com -X fullPage=true -X darkMode=true -o dark.png
  capture screenshot https://example.com -X selector=".main" -X format=webp -o element.webp`,
	Args: cobra.ExactArgs(1),
	RunE: runScreenshot,
}

var (
	screenshotOutput  string
	screenshotOptions []string
)

func init() {
	rootCmd.AddCommand(screenshotCmd)

	screenshotCmd.Flags().StringVarP(&screenshotOutput, "output", "o", "", "Output file (default: stdout)")
	screenshotCmd.Flags().StringArrayVarP(&screenshotOptions, "option", "X", nil, "API option as key=value (can be repeated)")
}

func runScreenshot(cmd *cobra.Command, args []string) error {
	targetURL := args[0]

	opts, err := parseOptions(screenshotOptions)
	if err != nil {
		return err
	}

	client := newCaptureClient()

	if dryRun {
		url, err := client.BuildImageURL(targetURL, opts)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	}

	verboseLog("Capturing screenshot of %s", targetURL)

	data, err := client.FetchImage(targetURL, opts)
	if err != nil {
		return fmt.Errorf("failed to capture screenshot: %w", err)
	}

	return writeOutput(data, screenshotOutput)
}
