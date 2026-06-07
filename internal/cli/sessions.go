package cli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	capture "github.com/techulus/capture-go"
)

var sessionsCmd = &cobra.Command{
	Use:   "sessions",
	Short: "Manage browser sessions",
	Long: `Manage stateful Capture browser sessions.

Sessions use the Capture Sessions API and bearer authentication derived from
CAPTURE_KEY and CAPTURE_SECRET.`,
}

var sessionsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a browser session",
	RunE:  runSessionsCreate,
}

var sessionsGetCmd = &cobra.Command{
	Use:   "get <session-id>",
	Short: "Get browser session metadata",
	Args:  cobra.ExactArgs(1),
	RunE:  runSessionsGet,
}

var sessionsCloseCmd = &cobra.Command{
	Use:   "close <session-id>",
	Short: "Close a browser session",
	Args:  cobra.ExactArgs(1),
	RunE:  runSessionsClose,
}

var sessionsActionCmd = &cobra.Command{
	Use:   "action <session-id> <type>",
	Short: "Execute a generic browser session action",
	Args:  cobra.ExactArgs(2),
	RunE:  runSessionsAction,
}

var (
	sessionMaxTTLSeconds      int
	sessionProxy              bool
	sessionBypassBotDetection bool
	sessionsPretty            bool
	sessionActionOptions      []string
	sessionActionPayloadJSON  string
)

func init() {
	rootCmd.AddCommand(sessionsCmd)
	sessionsCmd.AddCommand(sessionsCreateCmd, sessionsGetCmd, sessionsCloseCmd, sessionsActionCmd)

	sessionsCreateCmd.Flags().IntVar(&sessionMaxTTLSeconds, "max-ttl-seconds", 0, "Maximum session lifetime in seconds")
	sessionsCreateCmd.Flags().BoolVar(&sessionProxy, "proxy", false, "Use the authenticated user's configured browser proxy")
	sessionsCreateCmd.Flags().BoolVar(&sessionBypassBotDetection, "bypass-bot-detection", false, "Use Capture's bot-detection bypass browser when available")
	sessionsCreateCmd.Flags().BoolVar(&sessionsPretty, "pretty", false, "Pretty print JSON output")

	sessionsGetCmd.Flags().BoolVar(&sessionsPretty, "pretty", false, "Pretty print JSON output")
	sessionsCloseCmd.Flags().BoolVar(&sessionsPretty, "pretty", false, "Pretty print JSON output")

	sessionsActionCmd.Flags().StringArrayVarP(&sessionActionOptions, "option", "X", nil, "Action payload option as key=value (can be repeated)")
	sessionsActionCmd.Flags().StringVar(&sessionActionPayloadJSON, "payload-json", "", "Action payload as a JSON object")
	sessionsActionCmd.Flags().BoolVar(&sessionsPretty, "pretty", false, "Pretty print JSON output")
}

func runSessionsCreate(cmd *cobra.Command, args []string) error {
	client := newCaptureClient()
	options := &capture.CreateSessionOptions{
		MaxTtlSeconds:      sessionMaxTTLSeconds,
		Proxy:              sessionProxy,
		BypassBotDetection: sessionBypassBotDetection,
	}

	response, err := client.CreateSession(options)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return emitJSON(response, sessionsPretty)
}

func runSessionsGet(cmd *cobra.Command, args []string) error {
	client := newCaptureClient()
	response, err := client.GetSession(args[0])
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	return emitJSON(response, sessionsPretty)
}

func runSessionsClose(cmd *cobra.Command, args []string) error {
	client := newCaptureClient()
	response, err := client.CloseSession(args[0])
	if err != nil {
		return fmt.Errorf("failed to close session: %w", err)
	}

	return emitJSON(response, sessionsPretty)
}

func runSessionsAction(cmd *cobra.Command, args []string) error {
	payload, err := parseActionPayload(sessionActionPayloadJSON, sessionActionOptions)
	if err != nil {
		return err
	}

	client := newCaptureClient()
	response, err := client.ExecuteAction(args[0], args[1], payload)
	if err != nil {
		return fmt.Errorf("failed to execute action: %w", err)
	}

	return emitJSON(response, sessionsPretty)
}

func parseActionPayload(payloadJSON string, optionPairs []string) (capture.SessionActionPayload, error) {
	payload := capture.SessionActionPayload{}

	if payloadJSON != "" {
		if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
			return nil, fmt.Errorf("invalid --payload-json: %w", err)
		}
	}

	options, err := parseOptions(optionPairs)
	if err != nil {
		return nil, err
	}
	for key, value := range options {
		payload[key] = value
	}

	return payload, nil
}

func emitJSON(value interface{}, pretty bool) error {
	var (
		data []byte
		err  error
	)
	if pretty {
		data, err = json.MarshalIndent(value, "", "  ")
	} else {
		data, err = json.Marshal(value)
	}
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(data))
	return nil
}
