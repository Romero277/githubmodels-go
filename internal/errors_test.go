package internal

import "testing"

func TestAPIError(t *testing.T) {
	err := NewAPIError(404, "not found")
	expected := "GitHub Models API error: 404 - not found"
	if err.Error() != expected {
		t.Errorf("unexpected error string: got %q, want %q", err.Error(), expected)
	}
}
