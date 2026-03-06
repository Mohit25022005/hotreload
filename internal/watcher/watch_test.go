package watcher

import "testing"

func TestShouldWatch(t *testing.T) {

	tests := []struct {
		file     string
		expected bool
	}{
		{"main.go", true},
		{"go.mod", true},
		{"go.sum", true},
		{"index.html", false},
		{"script.js", false},
	}

	for _, test := range tests {

		result := shouldWatch(test.file)

		if result != test.expected {
			t.Fatalf("expected %v for file %s, got %v",
				test.expected,
				test.file,
				result)
		}
	}
}