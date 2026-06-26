package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetbool(t *testing.T) {
	os.Unsetenv("X_BOOL")
	if !getbool("X_BOOL", true) {
		t.Fatal("missing var should return the default")
	}
	t.Setenv("X_BOOL", "false")
	if getbool("X_BOOL", true) {
		t.Fatal("should parse false")
	}
	t.Setenv("X_BOOL", "garbage")
	if !getbool("X_BOOL", true) {
		t.Fatal("unparseable value should fall back to the default")
	}
}

func TestLoadDotEnvRealEnvWins(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte("FOO=fromfile\nBAR=barfile\n# comment\n\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("FOO", "fromenv") // already set → file must not override
	os.Unsetenv("BAR")         // unset → file value should load

	loadDotEnv(p)

	if got := os.Getenv("FOO"); got != "fromenv" {
		t.Fatalf("real env should win, got %q", got)
	}
	if got := os.Getenv("BAR"); got != "barfile" {
		t.Fatalf("missing var should load from file, got %q", got)
	}
	os.Unsetenv("BAR")
}
