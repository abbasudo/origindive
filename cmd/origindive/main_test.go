package main

import (
	"testing"
)

// TestMain_BasicUsage tests basic command invocation
func TestMain_BasicUsage(t *testing.T) {
	// Test that main package compiles and can be imported
	t.Log("Main package test placeholder")
}

// TestVersion_Flag tests version flag
func TestVersion_Flag(t *testing.T) {
	// This would normally test --version flag
	// Since main() calls os.Exit(), we can't test it directly
	// Instead, test flag parsing logic if extracted
	t.Log("Version flag test placeholder")
}

// TestHelp_Flag tests help flag
func TestHelp_Flag(t *testing.T) {
	// This would normally test --help flag
	t.Log("Help flag test placeholder")
}

// Note: Testing main() directly is challenging because it calls os.Exit()
// Best practice is to extract logic into testable functions and test those
// For now, these placeholder tests ensure the package compiles

// TestFormatResolvedSubdomains tests the --show-subdomains output
func TestFormatResolvedSubdomains(t *testing.T) {
	results := map[string][]string{
		"mail.example.com": {"203.0.113.12"},
		"api.example.com":  {"203.0.113.10", "203.0.113.20"},
		"dev.example.com":  {"203.0.113.11"},
	}

	got := formatResolvedSubdomains(results)
	want := "  → Resolved subdomains (3):\n" +
		"    api.example.com → 203.0.113.10, 203.0.113.20\n" +
		"    dev.example.com → 203.0.113.11\n" +
		"    mail.example.com → 203.0.113.12\n"

	if got != want {
		t.Errorf("formatResolvedSubdomains() =\n%s\nwant:\n%s", got, want)
	}
}

// TestFormatResolvedSubdomains_Empty tests output when nothing resolved
func TestFormatResolvedSubdomains_Empty(t *testing.T) {
	got := formatResolvedSubdomains(map[string][]string{})
	if got != "  → Resolved subdomains: none found\n" {
		t.Errorf("formatResolvedSubdomains(empty) = %q", got)
	}
}
