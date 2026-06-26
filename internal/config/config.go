// Package config loads runtime configuration from the environment, with an
// optional .env file for local development. Real environment variables always
// take precedence over .env entries.
package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime settings. Phase 1 uses Addr, Mock, and Dev; the HA
// fields are read now so the shape is stable for the live client (Phase 3).
type Config struct {
	Addr    string // TCP listen address, e.g. ":8080"
	HAHost  string // Home Assistant host[:port], e.g. "192.168.1.100:8123"
	HAToken string // long-lived access token (or Supervisor token in an add-on)
	HASSL   bool   // connect with wss:// instead of ws://
	Addon   bool   // running as a HA add-on: reach HA via the Supervisor proxy
	Mock    bool   // serve fabricated data instead of connecting to HA
	Dev     bool   // development mode (verbose logging; live template reload later)
}

// Load reads .env (best effort) then the environment and returns the Config.
func Load() Config {
	loadDotEnv(".env")

	// Add-on mode is explicit (ADDON) or inferred from the Supervisor token the
	// HA Supervisor injects into add-on containers. In that mode the token comes
	// from SUPERVISOR_TOKEN and HA is reached at ws://supervisor/core/websocket.
	addon := getbool("ADDON", os.Getenv("SUPERVISOR_TOKEN") != "")
	token := os.Getenv("HA_TOKEN")
	if token == "" {
		token = os.Getenv("SUPERVISOR_TOKEN")
	}

	return Config{
		Addr:    ":" + getenv("PORT", "8080"),
		HAHost:  os.Getenv("HA_HOST"),
		HAToken: token,
		HASSL:   getbool("HA_SSL", false),
		Addon:   addon,
		Mock:    getbool("MOCK", true),
		Dev:     getbool("DEV", false),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getbool(key string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}

// loadDotEnv parses simple KEY=VALUE lines from path and sets any variable not
// already present in the environment (so real env vars win). A missing file is
// not an error. Blank lines and lines beginning with # are ignored; surrounding
// single or double quotes on the value are stripped.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}
