# Module 02: Go Language Fundamentals & Syntax Guide

Welcome to Go Language Fundamentals. This guide starts with the exact anatomy of a Go program—explaining what every keyword means—before diving into syntax, memory layouts, and data structures.

---

## 1. Anatomy of a Go Program (Line-by-Line Breakdown)

Every Go source file follows a strict, deliberate structure:

```go
// Line 1: Package Declaration
package main

// Line 2: Import Declarations
import (
    "fmt"
    "time"
)

// Line 3: Entrypoint Function
func main() {
    fmt.Println("Hello, World! Time:", time.Now())
}
```

### A. What is `package main`?
- Every Go file **must** declare its package on line 1 (`package <name>`).
- The name **`main`** is special:
  - It indicates this file belongs to an **executable program** (not a shared library).
  - The Go compiler looks for an entrypoint function named `func main()` and builds a standalone executable binary.
  - If a file starts with any other name (e.g. `package utils`), it compiles into a **reusable library package**.

### B. What is `import (...)`?
- Declares the libraries or packages this file needs.
- Go's compiler is famously strict: **If you import a package and do not use it, compilation will fail with an error.** This prevents dependency bloat and dead code.
- Standard library packages (like `"fmt"` for formatted I/O, `"time"`, `"os"`) require no installation—they are built into Go.

### C. What is `func main()`?
- `func main()` is the entrypoint where execution begins when your compiled binary is launched.
- It takes **no parameters** and returns **no values**.
- Command-line arguments are accessed using the `os.Args` slice from the standard library `os` package.

### D. The Semicolon & Curly Brace Rule
Go has no semicolons at the ends of lines—the Go lexer automatically inserts them for you. Because of this automatic semicolon insertion, **opening braces `{` MUST stay on the same line**:
```go
// CORRECT:
func calculate() {
}

// SYNTAX ERROR in Go (compiler inserts semicolon after func header!):
func calculate()
{
}
```

---

## 2. Core Syntax: Variables, Types & Visibility

### A. Three Ways to Declare Variables
Go is statically typed, but provides expressive type inference:

```go
// 1. Short variable declaration := (Used 90% of the time inside functions)
name := "Divya"      // Infers type string
score := 98.5        // Infers type float64
count := 10          // Infers type int
isActive := true     // Infers type bool

// 2. Explicit var declaration (Useful when declaring without an initial value)
var age int = 30

// 3. Zero Value declaration (Guaranteed safe, default initial state)
var total int        // Automatically initialized to 0
var title string     // Automatically initialized to "" (empty string)
var flag bool        // Automatically initialized to false
var ptr *int         // Automatically initialized to nil
```

### B. Capitalization Rule: Exported (Public) vs. Unexported (Private)
Go does **not** have `public`, `private`, or `protected` keywords. Visibility is controlled entirely by the **first letter's case**:
- **Starts with an UPPERCASE letter (`User`, `CalculateTotal`, `Version`)**:
  - **Exported (Public)**: Visible and accessible to other packages importing this package.
- **Starts with a LOWERCASE letter (`user`, `calculateTotal`, `internalID`)**:
  - **Unexported (Private)**: Accessible **only** within the current package.

### C. Explicit Type Conversions (No Silent Coercion)
Go will **never** implicitly convert types. Passing an `int32` to an `int64` parameter is a compile error unless converted explicitly:
```go
var a int32 = 10
var b int64 = int64(a) // Explicit conversion mandatory
```

---

## 3. How Functions Work (Parameters & Return Types)

In Go, parameter types are written **after** parameter names.

### A. Basic Function Signature
```go
//          name   type,  name   type   --> return type
func Multiply(x     int,   y      int)          int {
    return x * y
}

// Shorthand when consecutive parameters share the same type:
func Add(x, y int) int {
    return x + y
}
```

### B. Multiple Return Values & The `val, err` Pattern
In Go, functions that can fail return both the result and an `error`. Go does not use exceptions (`try/catch`):
```go
// Returns a float64 AND an error
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0.0, errors.New("cannot divide by zero")
    }
    return a / b, nil // 'nil' means no error occurred
}

// Caller checks error explicitly:
result, err := Divide(10.0, 2.0)
if err != nil {
    fmt.Println("Error occurred:", err)
    return
}
fmt.Println("Result:", result)
```

---

## 4. Structs: How to Model Data

Go has no classes. Custom types and structured records are created with **`struct`**.

### A. Defining a Struct
```go
type Contact struct {
    ID    int
    Name  string
    Email string
}
```

### B. Instantiating a Struct
```go
// 1. Named-field instantiation (Recommended for clarity)
c1 := Contact{
    ID:    1,
    Name:  "Alice",
    Email: "alice@example.com",
}

// 2. Zero-value instantiation
var c2 Contact // c2.ID = 0, c2.Name = "", c2.Email = ""
c2.Name = "Bob"

// 3. Pointer to a Struct
// 'c3' holds memory address (*Contact)
c3 := &Contact{
    ID:   3,
    Name: "Charlie",
}
// Go automatically dereferences struct pointers (no c3->Name syntax needed!):
fmt.Println(c3.Name)
```

---

## 5. Internal Behavior: Memory Layouts & Data Structures

### Slice Internals
A slice in Go is not an array; it is a **slice header** (a 24-byte struct on 64-bit systems):
```text
+-------------------+
| Pointer (8 bytes) | ----> Pointer to element in underlying array
+-------------------+
| Length  (8 bytes) | ----> Current number of elements (len(s))
+-------------------+
| Capacity(8 bytes) | ----> Max elements before reallocation (cap(s))
+-------------------+
```

When you append beyond `cap`, Go allocates a new, larger array (typically doubling capacity for small slices) and copies existing elements over.

### Map Internals
A Go map is a pointer to an `hmap` struct managing a bucket array:
- Each bucket holds up to 8 key-value pairs.
- Keys are hashed with a randomized hash seed (preventing HashDoS attacks).
- Iteration order over a map is intentionally **randomized** across program runs.
- Maps are **not thread-safe** by default. Concurrent reads and writes trigger a fatal runtime crash.

### Pointer Mechanics
Pointers hold the memory address of a value.
- `&x` produces a pointer to `x` (`*T`).
- `*p` dereferences the pointer to access the underlying value.
- Go has **no pointer arithmetic** (keeping memory safe).

---

## 6. Design Discussion: Value vs. Pointer Semantics

| Scenario | Recommendation | Rationale |
| :--- | :--- | :--- |
| **Primitives (`int`, `bool`, `float64`)** | **Value** | Passing by value is extremely cheap (register-passed). |
| **Small Structs (e.g. `time.Time`, `Point{X, Y int}`)** | **Value** | Safer, immutable, no GC overhead, fits in registers. |
| **Large Structs (> 64 bytes)** | **Pointer** | Avoids expensive memory copying on every function call. |
| **Need Mutability in Callee** | **Pointer** | Necessary for modifications to reflect in the caller. |
| **Slices and Maps** | **Value** | Slices and maps are already header references to underlying buffers. |

---

## 7. Module Directory Structure

```text
02-language-fundamentals/
├── 01_syntax_and_types/main.go         # Declarations, types, zero values, conversions
├── 02_control_flow/main.go             # if, for, switch, defer behavior
├── 03_functions_and_closures/main.go   # Multiple returns, variadic, closures
├── 04_slices_arrays_maps/main.go       # Headers, append, maps, comma-ok idiom
├── 05_pointers_and_memory/main.go      # Value vs pointer semantics
└── exercises/
    ├── 01_temp_converter/main.go       # CLI Celsius/Fahrenheit conversion
    ├── 02_log_line_counter/main.go     # Parse and count log levels (INFO, WARN, ERROR)
    ├── 03_word_frequency/main.go       # Tokenizer and word frequency counter
    ├── 04_contact_manager/main.go      # In-memory CRUD contact manager
    └── 05_config_parser/main.go        # Key-value INI/env style parser
```

---

## 8. Understanding Check

1. **Why does Go use `package main` instead of naming packages after filenames like in Java or Python?**
   *Answer*: In Go, package names represent namespaces and compilation targets. Multiple `.go` files in the same directory share the same package name. `package main` specifically flags the directory as an executable program with a `main()` entrypoint.
2. **How does Go determine whether a struct field or function is public or private?**
   *Answer*: Purely by capitalization. Identifiers starting with a capital letter are exported (public); identifiers starting with a lowercase letter are unexported (private).
3. **What is the comma-ok idiom for maps?**
   *Answer*: `val, ok := myMap[key]`. If `ok` is `true`, the key exists. If `false`, the key is absent and `val` contains the type's default zero value.
