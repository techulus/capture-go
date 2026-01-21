package cli

import (
	"fmt"
	"net/http"
	"os"
	"time"

	capture "github.com/techulus/capture-go"

	"github.com/spf13/cobra"
)

var (
	useEdge bool
	verbose bool
	timeout time.Duration
	dryRun  bool

	captureKey    string
	captureSecret string
)

var rootCmd = &cobra.Command{
	Use:   "capture",
	Short: "Capture CLI - Screenshots, PDFs, and content extraction",
	Long: `Capture CLI is a command-line tool for taking screenshots, generating PDFs,
and extracting content from web pages using the Capture API.

Documentation: https://docs.capture.page/

Authentication is done via environment variables:
  CAPTURE_KEY    - Your Capture API key
  CAPTURE_SECRET - Your Capture API secret`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "version" || cmd.Name() == "completion" || cmd.Name() == "help" {
			return nil
		}

		captureKey = os.Getenv("CAPTURE_KEY")
		captureSecret = os.Getenv("CAPTURE_SECRET")

		if captureKey == "" || captureSecret == "" {
			return fmt.Errorf("CAPTURE_KEY and CAPTURE_SECRET environment variables are required")
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&useEdge, "edge", false, "Use edge server for faster response")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Request timeout")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Print the request URL without executing")
}

func newCaptureClient() *capture.Capture {
	httpClient := &http.Client{
		Timeout: timeout,
	}

	var opts []capture.Option
	opts = append(opts, capture.WithHTTPClient(httpClient))
	if useEdge {
		opts = append(opts, capture.WithEdge())
	}
	return capture.New(captureKey, captureSecret, opts...)
}

func verboseLog(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[verbose] "+format+"\n", args...)
	}
}
