# Module 03: Idiomatic Go Design

## 1. Conceptual Foundation

Go avoids the complexity of traditional Object-Oriented Programming (OOP) hierarchies. It has no `class`, no `extends`, and no inheritance. Instead, it achieves modularity, reusability, and abstraction through:
1. **Composition**: Embedding structs within other structs.
2. **Implicit Interfaces**: Types satisfy an interface automatically simply by implementing its method set.
3. **Values over Exceptions**: Errors are first-class values inspected and wrapped explicitly.

---

## 2. Internal Behavior: How Interfaces Work at Runtime

An interface value in Go is a two-word pair:
```text
+----------------------+
|  _type (itab pointer)| ----> Concrete type metadata & method table
+----------------------+
|  data  (void pointer)| ----> Pointer to the concrete value
+----------------------+
```

- If both `_type` and `data` are `nil`, the interface itself is `nil`.
- **Trap**: An interface holding a non-nil type pointer with a nil concrete pointer is **NOT nil**!
  ```go
  var err *MyCustomError = nil
  var i error = err
  fmt.Println(i == nil) // FALSE! 'i' holds type information (*MyCustomError), so it is not nil.
  ```

---

## 3. Design Discussion: Idiomatic Principles

### Accept Interfaces, Return Structs
- Functions should accept interfaces as parameters: This decouples the function from concrete implementations and enables easy unit testing and mocking.
- Functions should return concrete types (structs or pointers to structs): This avoids premature abstraction and gives callers full flexibility.

### Keep Interfaces Small
Go standard library interfaces typically have only 1 to 2 methods:
```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
```
"The bigger the interface, the weaker the abstraction." — Rob Pike

### Error Wrapping and Inspection (Go 1.13+)
Always add context when returning errors:
```go
// Wrapping: %w preserves the original error
if err != nil {
    return fmt.Errorf("failed to load configuration: %w", err)
}

// Inspection:
if errors.Is(err, os.ErrNotExist) { ... } // Checks error chain for target sentinel
var pathErr *os.PathError
if errors.As(err, &pathErr) { ... }       // Extracts custom error type from chain
```

---

## 4. Module Directory Structure

```text
03-idiomatic-go-design/
├── 01_structs_and_methods/main.go          # Value vs pointer receivers, struct embedding
├── 02_interfaces_and_polymorphism/main.go  # Implicit satisfaction, type assertions, io.Reader
├── 03_robust_error_handling/main.go        # Custom errors, errors.Is/As, wrapping with %w
└── mini-project-config-manager/            # Reusable config package
    ├── README.md                           # Documentation
    ├── config.go                           # Env parser, validator, defaults
    └── config_test.go                      # Unit tests
```
