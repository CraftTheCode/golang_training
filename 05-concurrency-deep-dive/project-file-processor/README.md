# Project: Concurrent File Processor

A production-ready concurrent directory scanner and file processor in Go demonstrating:
- **Bounded Worker Pool**: Controlled concurrency limiting simultaneous disk/CPU usage.
- **Context-Aware Cancellation**: Clean termination on abort or timeout without leaking goroutines.
- **Atomic & Thread-Safe Progress Reporting**: Live progress callbacks.
- **Concurrent Error Accumulation**: Structured aggregation of per-file errors.
- **Race Condition Immunity**: Validated with `go test -race`.

## Running the Unit & Stress Tests
```bash
go test -v -race ./05-concurrency-deep-dive/project-file-processor/...
```
