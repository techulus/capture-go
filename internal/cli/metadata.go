package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata <url>",
	Short: "Extract metadata from a web page",
	Long: `Extract metadata (title, description, Open Graph, etc.) from the specified URL.

Options are passed as key=value pairs using the -X flag.
See https://docs.capture.page/ for available options.

Examples:
  capture metadata https://example.com
  capture metadata https://example.com --pretty
  capture metadata https://example.com -o metadata.json
  capture metadata https://example.com -X delay=1000 --pretty`,
	Args: cobra.ExactArgs(1),
	RunE: runMetadata,
}

var (
	metadataOutput  string
	metadataPretty  bool
	metadataOptions []string
)

func init() {
	rootCmd.AddCommand(metadataCmd)

	metadataCmd.Flags().StringVarP(&metadataOutput, "output", "o", "", "Output file (default: stdout)")
	metadataCmd.Flags().BoolVar(&metadataPretty, "pretty", false, "Pretty print JSON output")
	metadataCmd.Flags().StringArrayVarP(&metadataOptions, "option", "X", nil, "API option as key=value (can be repeated)")
}

func runMetadata(cmd *cobra.Command, args []string) error {
	targetURL := args[0]

	opts, err := parseOptions(metadataOptions)
	if err != nil {
		return err
	}

	client := newCaptureClient()

	if dryRun {
		url, err := client.BuildMetadataURL(targetURL, opts)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	}

	verboseLog("Extracting metadata from %s", targetURL)

	metadata, err := client.FetchMetadata(targetURL, opts)
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %w", err)
	}

	var data []byte
	if metadataPretty {
		data, err = json.MarshalIndent(metadata, "", "  ")
	} else {
		data, err = json.Marshal(metadata)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return writeOutput(data, metadataOutput)
}
