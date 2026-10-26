package subdomain

import (
	"context"
	"reflect"
	"testing"
)

func TestRecorder_AddAndResults(t *testing.T) {
	r := NewRecorder()
	r.Add("API.Example.com.", []string{"203.0.113.2", "203.0.113.1"})
	r.Add("api.example.com", []string{"203.0.113.1"})
	r.Add("empty.example.com", nil)

	want := map[string][]string{"api.example.com": {"203.0.113.1", "203.0.113.2"}}
	if got := r.Results(); !reflect.DeepEqual(got, want) {
		t.Errorf("Results() = %v, want %v", got, want)
	}
}

func TestRecord_ContextWithoutRecorder(t *testing.T) {
	// Must be a no-op, not a panic
	Record(context.Background(), "api.example.com", []string{"203.0.113.1"})
}

func TestRecord_ContextWithRecorder(t *testing.T) {
	r := NewRecorder()
	ctx := WithRecorder(context.Background(), r)
	Record(ctx, "dev.example.com", []string{"203.0.113.3"})

	if got := r.Results()["dev.example.com"]; len(got) != 1 || got[0] != "203.0.113.3" {
		t.Errorf("recorded IPs = %v", got)
	}
}

func TestApexDomain(t *testing.T) {
	tests := map[string]string{
		"adminapi.sarmayegozarebartar.com": "sarmayegozarebartar.com",
		"example.com":                      "example.com",
		"a.b.example.co.uk":                "example.co.uk",
		"WWW.Example.COM.":                 "example.com",
	}
	for in, want := range tests {
		if got := ApexDomain(in); got != want {
			t.Errorf("ApexDomain(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAllIn(t *testing.T) {
	set := map[string]bool{"203.0.113.1": true}
	if !allIn([]string{"203.0.113.1"}, set) {
		t.Error("expected wildcard match")
	}
	if allIn([]string{"203.0.113.1", "203.0.113.9"}, set) {
		t.Error("expected no match when an IP is outside the wildcard set")
	}
	if allIn([]string{"203.0.113.1"}, map[string]bool{}) {
		t.Error("empty wildcard set must match nothing")
	}
}
