package tq

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var version string

// GetVersion returns the current version of tq.
func GetVersion() string {
	result := strings.TrimSpace(version)
	if result == "" {
		return "v0.0.0-dev"
	}
	return result
}
