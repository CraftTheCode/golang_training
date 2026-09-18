# Module 08: Testing, Benchmarking & Production Practices

## 1. Conceptual Foundation

Testing in Go is built directly into the language toolchain (`go test`). You don't need third-party test runners or heavy assertion frameworks. The Go philosophy emphasizes:
- **Table-Driven Tests**: Grouping test cases as structured slice tables to eliminate repetitive boilerplate and test edge cases exhaustively.
- **Interface-Based Mocking**: Creating lightweight fakes or mocks adhering to single-method interfaces rather than heavy bytecode-manipulating mock libraries.
- **First-Class Benchmarks & Allocations**: Built-in microbenchmarks measuring execution speed (ns/op) and memory allocations (B/op, allocs/op).

---

## 2. Testing Toolkit Commands

```bash
# 1. Run all tests in the current package and subpackages
go test -v ./...

# 2. Run tests with the Race Detector (mandatory in CI)
go test -race ./...

# 3. Run benchmarks and report memory allocations
go test -bench=. -benchmem ./...

# 4. Generate and view HTML test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 5. Static code analysis
go vet ./...
```

---

## 3. Module Directory Structure

```text
08-testing-and-benchmarking/
├── 01_table_driven_tests/              # Canonical table test pattern with t.Run
│   ├── calc.go
│   └── calc_test.go
├── 02_mocking_and_fakes/               # Interface-driven mocking without third-party frameworks
│   ├── notifier.go
│   └── notifier_test.go
├── 03_http_testing/                    # net/http/httptest response recording
│   ├── handler.go
│   └── handler_test.go
└── 04_benchmarking_and_profiling/      # b.ResetTimer, b.ReportAllocs, string concat comparison
    ├── string_builder.go
    └── string_builder_test.go
```
