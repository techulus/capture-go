package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	capture "github.com/techulus/capture-go"
)

func writeOutput(data []byte, outputFile string) error {
	if outputFile == "" || outputFile == "-" {
		if _, err := os.Stdout.Write(data); err != nil {
			return fmt.Errorf("failed to write to stdout: %w", err)
		}
		return nil
	}

	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write to %s: %w", outputFile, err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Written to %s (%d bytes)\n", outputFile, len(data))
	}

	return nil
}

func writeStringOutput(data string, outputFile string) error {
	return writeOutput([]byte(data), outputFile)
}

func parseOptions(options []string) (capture.RequestOptions, error) {
	opts := capture.RequestOptions{}

	for _, opt := range options {
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid option format: %s (expected key=value)", opt)
		}

		key := parts[0]
		value := parts[1]

		if intVal, err := strconv.Atoi(value); err == nil {
			opts[key] = intVal
		} else if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			opts[key] = floatVal
		} else if boolVal, err := strconv.ParseBool(value); err == nil {
			opts[key] = boolVal
		} else {
			opts[key] = value
		}
	}

	return opts, nil
}
