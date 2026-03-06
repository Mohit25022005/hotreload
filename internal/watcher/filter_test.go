package watcher

import "testing"

func TestShouldIgnore(t *testing.T) {

	tests := []struct {
		path     string
		expected bool
	}{
		{"project/.git/config", true},
		{"project/node_modules/lib.js", true},
		{"project/bin/server", true},
		{"project/tmp/file.tmp", true},
		{"project/main.go", false},
	}

	for _, test := range tests {

		result := shouldIgnore(test.path)

		if result != test.expected {
			t.Fatalf("expected %v for path %s, got %v",
				test.expected,
				test.path,
				result)
		}
	}
}