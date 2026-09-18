package processor

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileResult holds the processing outcome for a single file.
type FileResult struct {
	Path      string
	BytesRead int64
	LineCount int
	Error     error
}

// ProgressCallback is invoked periodically as files are processed.
type ProgressCallback func(completed, total int)

// Options holds configuration for the processor.
type Options struct {
	NumWorkers int
	Extension  string
	OnProgress ProgressCallback
}

// Summary summarizes the entire batch processing run.
type Summary struct {
	TotalFiles int
	Successful int
	Failed     int
	TotalBytes int64
	TotalLines int
	Duration   time.Duration
	Results    []FileResult
}

// ProcessDirectory scans rootDir for files matching ext and processes them concurrently.
func ProcessDirectory(ctx context.Context, rootDir string, opts Options) (*Summary, error) {
	startTime := time.Now()

	if opts.NumWorkers <= 0 {
		opts.NumWorkers = 4
	}

	// 1. Scan directory to collect matching file paths
	var files []string
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if opts.Extension == "" || strings.HasSuffix(strings.ToLower(d.Name()), strings.ToLower(opts.Extension)) {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory %s: %w", rootDir, err)
	}

	total := len(files)
	if total == 0 {
		return &Summary{Duration: time.Since(startTime)}, nil
	}

	// 2. Setup channels & synchronization
	jobs := make(chan string, total)
	results := make(chan FileResult, total)
	var wg sync.WaitGroup

	// 3. Start worker pool
	for i := 0; i < opts.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return // Stop immediately on context cancellation
				case path, ok := <-jobs:
					if !ok {
						return // Jobs channel closed
					}
					results <- processSingleFile(ctx, path)
				}
			}
		}()
	}

	// 4. Feed jobs to workers in background
	go func() {
		for _, f := range files {
			select {
			case <-ctx.Done():
				break
			case jobs <- f:
			}
		}
		close(jobs)
	}()

	// 5. Close results once all workers exit
	go func() {
		wg.Wait()
		close(results)
	}()

	// 6. Collect results
	summary := &Summary{
		TotalFiles: total,
		Results:    make([]FileResult, 0, total),
	}

	for res := range results {
		summary.Results = append(summary.Results, res)
		if res.Error != nil {
			summary.Failed++
		} else {
			summary.Successful++
			summary.TotalBytes += res.BytesRead
			summary.TotalLines += res.LineCount
		}

		if opts.OnProgress != nil {
			opts.OnProgress(len(summary.Results), total)
		}
	}

	summary.Duration = time.Since(startTime)

	if ctx.Err() != nil {
		return summary, ctx.Err()
	}

	return summary, nil
}

// processSingleFile reads and computes basic metrics for a file.
func processSingleFile(ctx context.Context, path string) FileResult {
	select {
	case <-ctx.Done():
		return FileResult{Path: path, Error: ctx.Err()}
	default:
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return FileResult{Path: path, Error: err}
	}

	// Count lines
	lines := 0
	for _, b := range data {
		if b == '\n' {
			lines++
		}
	}
	if len(data) > 0 && lines == 0 {
		lines = 1
	}

	return FileResult{
		Path:      path,
		BytesRead: int64(len(data)),
		LineCount: lines,
		Error:     nil,
	}
}
