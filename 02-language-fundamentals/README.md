# Module 02: Go Language Fundamentals

## 1. Conceptual Foundation

Go's syntax is deliberate, compact, and designed for readability across teams. Unlike languages with multiple syntactic sugars or implicit type coercions, Go favors **explicitness**:
- Every Go program belongs to a `package`.
- Executable programs start with `package main` and execution begins at `func main()`.
- Unused variables and unused imports are compile-time errors.
- Variables are statically typed; type conversions are always explicit.

---

## 2. Internal Behavior: Memory Layouts & Data Structures

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
- Go has **no pointer arithmetic** (unless using the `unsafe` package), keeping memory safe.

---

## 3. Design Discussion: Value vs. Pointer Semantics

| Scenario | Recommendation | Rationale |
| :--- | :--- | :--- |
| **Primitives (`int`, `bool`, `float64`)** | **Value** | Passing by value is extremely cheap (register-passed). |
| **Small Structs (e.g. `time.Time`, `Point{X, Y int}`)** | **Value** | Safer, immutable, no GC overhead, fits in registers. |
| **Large Structs (> 64 bytes)** | **Pointer** | Avoids expensive memory copying on every function call. |
| **Need Mutability in Callee** | **Pointer** | Necessary for modifications to reflect in the caller. |
| **Slices and Maps** | **Value** | Slices and maps are already header references to underlying buffers. |

---

## 4. Module Directory Structure

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

## 5. Understanding Check

1. **What happens when you slice a slice: `b := a[1:3]`?**
   *Answer*: `b` shares the exact same underlying array as `a`. Modifying elements in `b` modifies elements in `a` unless `b` exceeds its capacity and triggers a reallocation via `append`.
2. **How do you safely check if a key exists in a map?**
   *Answer*: Use the comma-ok idiom: `val, ok := myMap[key]`. If `ok` is `false`, the key was not found.
3. **When is a `defer` statement evaluated?**
   *Answer*: The function arguments to the deferred call are evaluated **immediately** when the `defer` line is encountered, but the actual execution of the deferred function body occurs when the enclosing function returns.
