package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pdfCmd = &cobra.Command{
	Use:   "pdf <url>",
	Short: "Generate a PDF from a web page",
	Long: `Generate a PDF document from the specified URL.

Options are passed as key=value pairs using the -X flag.
See https://docs.capture.page/ for available options.

Examples:
  capture pdf https://example.com -o document.pdf
  capture pdf https://example.com -X format=A4 -X landscape=true -o landscape.pdf
  capture pdf https://example.com -X printBackground=true -o styled.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: runPDF,
}

var (
	pdfOutput  string
	pdfOptions []string
)

func init() {
	rootCmd.AddCommand(pdfCmd)

	pdfCmd.Flags().StringVarP(&pdfOutput, "output", "o", "", "Output file (default: stdout)")
	pdfCmd.Flags().StringArrayVarP(&pdfOptions, "option", "X", nil, "API option as key=value (can be repeated)")
}

func runPDF(cmd *cobra.Command, args []string) error {
	targetURL := args[0]

	opts, err := parseOptions(pdfOptions)
	if err != nil {
		return err
	}

	client := newCaptureClient()

	if dryRun {
		url, err := client.BuildPDFURL(targetURL, opts)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	}

	verboseLog("Generating PDF from %s", targetURL)

	data, err := client.FetchPDF(targetURL, opts)
	if err != nil {
		return fmt.Errorf("failed to generate PDF: %w", err)
	}

	return writeOutput(data, pdfOutput)
}
