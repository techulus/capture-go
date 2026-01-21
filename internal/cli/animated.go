package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var animatedCmd = &cobra.Command{
	Use:   "animated <url>",
	Short: "Create an animated recording of a web page",
	Long: `Create an animated GIF or video recording of the specified URL.

Options are passed as key=value pairs using the -X flag.
See https://docs.capture.page/ for available options.

Examples:
  capture animated https://example.com -o recording.gif
  capture animated https://example.com -X duration=5 -X vw=1280 -o video.gif
  capture animated https://example.com -X format=mp4 -o recording.mp4
  capture animated https://example.com -X darkMode=true -X deviceScale=2 -o retina.gif`,
	Args: cobra.ExactArgs(1),
	RunE: runAnimated,
}

var (
	animatedOutput  string
	animatedOptions []string
)

func init() {
	rootCmd.AddCommand(animatedCmd)

	animatedCmd.Flags().StringVarP(&animatedOutput, "output", "o", "", "Output file (default: stdout)")
	animatedCmd.Flags().StringArrayVarP(&animatedOptions, "option", "X", nil, "API option as key=value (can be repeated)")
}

func runAnimated(cmd *cobra.Command, args []string) error {
	targetURL := args[0]

	opts, err := parseOptions(animatedOptions)
	if err != nil {
		return err
	}

	client := newCaptureClient()

	if dryRun {
		url, err := client.BuildAnimatedURL(targetURL, opts)
		if err != nil {
			return err
		}
		fmt.Println(url)
		return nil
	}

	verboseLog("Creating animated capture of %s", targetURL)

	data, err := client.FetchAnimated(targetURL, opts)
	if err != nil {
		return fmt.Errorf("failed to create animated capture: %w", err)
	}

	return writeOutput(data, animatedOutput)
}
