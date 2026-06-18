package cli

import "testing"

func TestSessionsCreateCommandHasCDPFlag(t *testing.T) {
	flag := sessionsCreateCmd.Flags().Lookup("cdp")
	if flag == nil {
		t.Fatal("expected sessions create --cdp flag")
	}
	if flag.DefValue != "false" {
		t.Fatalf("expected --cdp default false, got %q", flag.DefValue)
	}
}
