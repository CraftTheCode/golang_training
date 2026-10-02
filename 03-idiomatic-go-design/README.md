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
### Struct Tags: Field Annotations and Metadata

Struct tags are string literals placed in backticks after field declarations:
```go
type User struct {
    ID        int64     `json:"id" db:"user_id"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`                // "-" means never include in serialization
    Nickname  string    `json:"nickname,omitempty"` // omits field if zero-value ("")
    CreatedAt time.Time `json:"created_at"`
}
```

**How They Work**:
- At compile-time, the compiler stores tag strings inside the type's runtime metadata (`reflect.StructTag`).
- Packages like `encoding/json`, `encoding/xml`, ORMs (`gorm`), and validators (`go-playground/validator`) read these tags at runtime using the `reflect` package to determine custom field names, validation constraints, or column mappings.
- Struct tags eliminate the need for verbose mapping code or separate DTO configuration files.

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

## 4. Packages and Modules


### The Fundamental Rule: 1 Directory = 1 Package
The Go toolchain enforces that **all `.go` files in the same directory must declare the same `package` name**. A directory IS a package.

```text
myapp/
├── main.go           → package main     ← executable (entry point)
├── helpers.go         → package main     ← same package, no import needed
├── utils/
│   ├── strings.go    → package utils    ← separate package
│   └── math.go       → package utils    ← same package as strings.go
└── internal/
    └── auth/
        └── auth.go   → package auth     ← internal package (restricted access)
```

### Exported vs Unexported (Visibility)
Go has no `public`, `private`, or `protected` keywords. Visibility is controlled entirely by **capitalization**:

| Identifier | Starts With | Accessible Outside Package? |
| :--- | :--- | :--- |
| `UserName` | Uppercase | ✅ Yes — **Exported** |
| `userName` | Lowercase | ❌ No — **Unexported** |

This applies to functions, types, struct fields, constants, and variables.

### `go.mod` and Import Paths
Every Go project has a `go.mod` file at its root:
```
module golang_training

go 1.22
```

Import paths are built from the module name + directory path:
```go
import "golang_training/utils"           // imports package in ./utils/
import "golang_training/internal/auth"   // imports package in ./internal/auth/
```

### The `internal/` Convention
Any package under an `internal/` directory can only be imported by code **within the parent of `internal/`**. This is enforced by the Go compiler — not just convention.

```text
myapp/
├── cmd/server/main.go      ← CAN import myapp/internal/auth
├── internal/
│   └── auth/auth.go        ← restricted
└── pkg/utils/utils.go      ← CANNOT import myapp/internal/auth (compile error)
```

### `go run .` vs `go run main.go`
```bash
go run main.go    # Compiles ONLY main.go — misses other files in the package!
go run .          # Compiles ALL .go files in the directory as one package ✅
```

Always use `go run .` (or `go build`) when your package has multiple files.

### Same-Package vs Cross-Package Access

| Scenario | Import Needed? | Visibility Rule |
| :--- | :--- | :--- |
| Calling a function from **another file in the same directory** | ❌ No import | All identifiers visible (exported AND unexported) |
| Calling a function from **a different directory/package** | ✅ Yes, must import | Only **exported** (uppercase) identifiers visible |

---

## 5. Generics (Go 1.18+)

Go 1.18 introduced **type parameters**, allowing functions and types to work with multiple types without code duplication.

### Basic Syntax
```go
// Without generics: need separate functions for each type
func MinInt(a, b int) int { ... }
func MinFloat(a, b float64) float64 { ... }

// With generics: one function handles all comparable, ordered types
func Min[T constraints.Ordered](a, b T) T {
    if a < b {
        return a
    }
    return b
}
```

### Type Constraints
Type parameters must satisfy a **constraint** (an interface that specifies allowed types):
```go
// Built-in constraints from the "constraints" package (or "cmp" in Go 1.21+):
// - constraints.Ordered: types supporting <, <=, >, >=
// - comparable: types supporting == and !=
// - any: alias for interface{} — no restrictions

func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}
```

### When to Use Generics
| Use Generics For | Don't Use Generics For |
| :--- | :--- |
| Generic data structures (stacks, queues, trees) | Business logic with specific types |
| Utility functions (Min, Max, Contains, Filter) | Simple functions with 1-2 concrete types |
| Reducing copy-paste code across types | Making code "look clever" |

---

## 6. Module Directory Structure

```text
03-idiomatic-go-design/
├── 01_structs_and_methods/main.go          # Value vs pointer receivers, struct embedding
├── 02_interfaces_and_polymorphism/main.go  # Implicit satisfaction, type assertions, io.Reader
├── 03_robust_error_handling/main.go        # Custom errors, errors.Is/As, wrapping with %w
├── 04_packages_and_modules/                # Package organization demo
│   ├── main.go                             # Imports and uses the greetpkg package
│   └── greetpkg/greet.go                   # Exported vs unexported in a separate package
├── 05_generics_intro/main.go               # Type parameters, constraints, generic utilities
└── mini-project-config-manager/            # Reusable config package
    ├── README.md                           # Documentation
    ├── config.go                           # Env parser, validator, defaults
    └── config_test.go                      # Unit tests
```

---

## 7. Understanding Check

1. **Can two files in the same directory have different `package` names?**
   *Answer*: No. The Go toolchain enforces that all `.go` files in a directory must share the same package name. The only exception is `_test.go` files, which can use an external test package (e.g., `package foo_test`).

2. **If two directories both declare `package utils`, are they the same package?**
   *Answer*: No. They are separate, unrelated packages that happen to share a name. The directory path is the true package identity, not the name.

3. **When should you use generics vs interfaces?**
   *Answer*: Use generics when the operation is the same across types (e.g., sorting, finding min/max). Use interfaces when different types need different implementations of the same behavior (e.g., `Reader`, `Writer`).

