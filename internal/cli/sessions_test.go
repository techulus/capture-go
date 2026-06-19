package cli

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestSessionsCreateCommandHasCDPFlag(t *testing.T) {
	flag := sessionsCreateCmd.Flags().Lookup("cdp")
	if flag == nil {
		t.Fatal("expected sessions create --cdp flag")
	}
	if flag.DefValue != "false" {
		t.Fatalf("expected --cdp default false, got %q", flag.DefValue)
	}
}

func TestSessionsCommandsInheritDryRun(t *testing.T) {
	cmds := []*cobra.Command{sessionsCreateCmd, sessionsGetCmd, sessionsCloseCmd, sessionsActionCmd}
	for _, cmd := range cmds {
		if cmd.InheritedFlags().Lookup("dry-run") == nil {
			t.Fatalf("expected sessions %s to inherit --dry-run flag", cmd.Name())
		}
	}
}

func TestIsSessionsCommand(t *testing.T) {
	if !isSessionsCommand(sessionsActionCmd) {
		t.Error("expected sessions action to be recognised as a sessions command")
	}
	if isSessionsCommand(screenshotCmd) {
		t.Error("did not expect screenshot to be a sessions command")
	}
}

func TestDryRunSkipsCredentialRequirementForSessions(t *testing.T) {
	t.Setenv("CAPTURE_KEY", "")
	t.Setenv("CAPTURE_SECRET", "")

	prev := dryRun
	dryRun = true
	defer func() { dryRun = prev }()

	if err := rootCmd.PersistentPreRunE(sessionsCreateCmd, nil); err != nil {
		t.Fatalf("expected sessions dry-run to skip credential check, got %v", err)
	}
	if err := rootCmd.PersistentPreRunE(screenshotCmd, nil); err == nil {
		t.Fatal("expected non-session dry-run to still require credentials")
	}
}
