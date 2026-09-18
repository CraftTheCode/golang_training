package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func createTestDirectory(t *testing.T, count int) string {
	t.Helper()
	dir := t.TempDir()

	for i := 1; i <= count; i++ {
		filename := filepath.Join(dir, fmt.Sprintf("test_file_%d.txt", i))
		content := fmt.Sprintf("Line 1 in file %d\nLine 2 in file %d\nLine 3 in file %d\n", i, i, i)
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}
	return dir
}

func TestProcessDirectory_Success(t *testing.T) {
	fileCount := 20
	dir := createTestDirectory(t, fileCount)

	var progressUpdates atomic.Int32
	opts := Options{
		NumWorkers: 4,
		Extension:  ".txt",
		OnProgress: func(completed, total int) {
			progressUpdates.Add(1)
		},
	}

	summary, err := ProcessDirectory(context.Background(), dir, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if summary.TotalFiles != fileCount {
		t.Errorf("expected %d total files, got %d", fileCount, summary.TotalFiles)
	}
	if summary.Successful != fileCount {
		t.Errorf("expected %d successful files, got %d", fileCount, summary.Successful)
	}
	if summary.Failed != 0 {
		t.Errorf("expected 0 failed, got %d", summary.Failed)
	}
	if progressUpdates.Load() == 0 {
		t.Errorf("expected progress updates, got 0")
	}
}

func TestProcessDirectory_Cancellation(t *testing.T) {
	fileCount := 50
	dir := createTestDirectory(t, fileCount)

	ctx, cancel := context.WithCancel(context.Background())

	opts := Options{
		NumWorkers: 2,
		Extension:  ".txt",
		OnProgress: func(completed, total int) {
			if completed >= 5 {
				cancel() // Cancel early after 5 files
			}
		},
	}

	summary, err := ProcessDirectory(ctx, dir, opts)
	if err == nil && summary.Successful == fileCount {
		t.Errorf("expected cancellation error, but completed all files")
	}
}
