# Module 08: Testing, Benchmarking & Production Practices

## 1. Conceptual Foundation

Testing in Go is built directly into the language toolchain (`go test`). You don't need third-party test runners or heavy assertion frameworks. The Go philosophy emphasizes:
- **Table-Driven Tests**: Grouping test cases as structured slice tables to eliminate repetitive boilerplate and test edge cases exhaustively.
- **Interface-Based Mocking**: Creating lightweight fakes or mocks adhering to single-method interfaces rather than heavy bytecode-manipulating mock libraries.
- **First-Class Benchmarks & Allocations**: Built-in microbenchmarks measuring execution speed (ns/op) and memory allocations (B/op, allocs/op).

---

## 2. Core Testing Concepts & Patterns

### A. Anatomy of a Table-Driven Test
Table-driven testing is Go's standard idiom for covering multiple scenarios without duplicating test harness code:
```go
func TestCalculate(t *testing.T) {
    // 1. Define slice of anonymous test case structs
    tests := []struct {
        name     string  // Human-readable scenario name
        input    int     // Input values
        expected int     // Expected outcome
        wantErr  bool    // Whether an error is expected
    }{
        {name: "positive value", input: 5, expected: 25, wantErr: false},
        {name: "zero value", input: 0, expected: 0, wantErr: false},
        {name: "negative value fails", input: -1, expected: 0, wantErr: true},
    }

    // 2. Iterate through test cases
    for _, tc := range tests {
        // 3. t.Run creates an isolated subtest
        t.Run(tc.name, func(t *testing.T) {
            got, err := Calculate(tc.input)
            if (err != nil) != tc.wantErr {
                t.Fatalf("Calculate(%d) unexpected error status: %v", tc.input, err)
            }
            if got != tc.expected {
                t.Errorf("Calculate(%d) = %d; want %d", tc.input, got, tc.expected)
            }
        })
    }
}
```

### B. Why `t.Run` Subtests Matter
1. **Targeted Execution**: You can run a single subtest from the terminal using regular expressions:
   ```bash
   go test -v -run TestCalculate/negative
   ```
2. **Failure Isolation**: If one test case fails, remaining cases in the table still execute.
3. **Parallel Execution**: Adding `t.Parallel()` inside `t.Run` lets independent table cases run concurrently across CPU cores.

### C. `t.Errorf` vs `t.Fatalf`
- **`t.Errorf(...)`**: Marks the test as failed, prints diagnostic message, but **continues executing the rest of the test**. Use this for checking individual output fields.
- **`t.Fatalf(...)`**: Marks the test as failed, prints message, and **immediately aborts execution of that test function**. Use this when subsequent code depends on the result (e.g., checking if `err != nil` before dereferencing a pointer that would otherwise panic).

### D. Writing Benchmarks (`*testing.B`)
Go benchmarks measure operations per second and heap allocations:
```go
func BenchmarkBuilder(b *testing.B) {
    b.ReportAllocs() // Tracks heap allocations (B/op and allocs/op)
    b.ResetTimer()   // Excludes setup time from benchmark measurements

    for i := 0; i < b.N; i++ {
        // b.N is dynamically adjusted by the Go test runner until
        // the benchmark runs for a statistically stable duration (default 1s).
        BuildString(100)
    }
}
```

---

## 3. Testing Toolkit Commands

```bash
# 1. Run all tests in the current package and subpackages with verbose output
go test -v ./...

# 2. Run a specific test or subtest
go test -v -run TestCalculate/negative ./...

# 3. Run tests with the Race Detector (mandatory in CI)
go test -race ./...

# 4. Run benchmarks and report memory allocations
go test -bench=. -benchmem ./...

# 5. Generate and view HTML test coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 6. Static code analysis
go vet ./...
```

---


## 4. Module Directory Structure

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
