//go:build integration

package integration

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func testdataDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "testdata")
}

func LoadTestEvent[T any](t *testing.T, filename string) T {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(testdataDir(), filename))
	if err != nil {
		t.Fatalf("failed to read testdata %s: %v", filename, err)
	}
	var event T
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatalf("failed to unmarshal testdata %s: %v", filename, err)
	}
	return event
}
