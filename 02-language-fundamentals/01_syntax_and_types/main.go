package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println("=== Variable Declarations ===")

	// 1. Explicit declaration with type
	var age int = 30

	// 2. Type-inferred declaration (package or function level)
	var name = "Gopher"

	// 3. Short variable declaration syntax := (function-scoped only)
	score := 99.58 // with .1f, it will round off to ceil value (99.6)

	// 4. Multiple variable declaration
	var (
		city    string = "Bengaluru"
		country string = "India"
		active  bool   = true
	)

	fmt.Printf("User: %s, Age: %d, Score: %.1f, Location: %s, %s, Active: %t\n",
		name, age, score, city, country, active)

	fmt.Println("\n=== Numeric Types and Zero Values ===")
	var defaultInt int
	var defaultFloat float64
	var defaultBool bool
	var defaultString string

	fmt.Printf("int zero value: %d\n", defaultInt)
	fmt.Printf("float64 zero value: %f\n", defaultFloat)
	fmt.Printf("bool zero value: %t\n", defaultBool)
	fmt.Printf("string zero value: %q\n", defaultString)

	fmt.Println("\n=== Explicit Type Conversions ===")
	// Go never coerces types implicitly. Conversions must be explicit.
	var integerVal int = 42
	var floatVal float64 = float64(integerVal)
	var unsignedVal uint = uint(integerVal)

	fmt.Printf("Converted: int(%d) -> float64(%.2f) -> uint(%d)\n",
		integerVal, floatVal, unsignedVal)

	// Math functions require float64
	hypot := math.Sqrt(float64(3*3 + 4*4))
	fmt.Printf("Hypotenuse of 3 and 4: %.2f\n", hypot)

	fmt.Println("\n=== Constants and Iota ===")
	// Constants are evaluated at compile-time — they CANNOT be changed at runtime.
	// An untyped constant (no explicit type) adapts to whatever context uses it.
	const Pi = 3.14159

	// iota is a compile-time integer generator:
	//   - Starts at 0 in each new const (...) block
	//   - Increments by 1 for EACH LINE (not each use)
	//   - The EXPRESSION from the first line is inherited by all subsequent lines
	//     (that's why StatusActive doesn't need "= iota" — it inherits it)
	const (
		StatusPending    = iota // 0 — iota starts at 0
		StatusActive            // 1 — inherits the expression "= iota", iota is now 1
		StatusSuspended         // 2
		StatusTerminated        // 3
	)

	fmt.Printf("Pi: %f\n", Pi)
	fmt.Printf("Statuses: Pending=%d, Active=%d, Suspended=%d, Terminated=%d\n",
		StatusPending, StatusActive, StatusSuspended, StatusTerminated)

	// ── Advanced iota: Bit Flags (Bitmask Permissions) ──
	// Each constant occupies a SINGLE BIT position using left-shift (<<).
	// This is the standard pattern for Unix-style permissions and feature flags.
	const (
		PermRead    = 1 << iota // 1 << 0 = 1   (binary: 001)
		PermWrite               // 1 << 1 = 2   (binary: 010)
		PermExecute             // 1 << 2 = 4   (binary: 100)
	)

	// Combine permissions using bitwise OR (|)
	readWrite := PermRead | PermWrite // 1 | 2 = 3 (binary: 011)
	fmt.Printf("\nBit Flags: Read=%d, Write=%d, Execute=%d\n", PermRead, PermWrite, PermExecute)
	fmt.Printf("Combined ReadWrite = %d (binary: %03b)\n", readWrite, readWrite)

	// Check if a specific permission is set using bitwise AND (&)
	hasRead := (readWrite & PermRead) != 0
	hasExec := (readWrite & PermExecute) != 0
	fmt.Printf("HasRead=%t, HasExecute=%t\n", hasRead, hasExec)

	// ── Advanced iota: Byte Sizes with expression inheritance ──
	// The expression "1 << (10 * iota)" is set on the FIRST line.
	// Every subsequent line inherits this same expression, but with iota incremented.
	// iota=0 → 1<<0  = 1      (but we skip this with _ since 1 byte isn't useful)
	// iota=1 → 1<<10 = 1024   (KB)
	// iota=2 → 1<<20 = 1,048,576 (MB)
	const (
		_  = iota             // 0 — skip with blank identifier (we don't need "1 byte")
		KB = 1 << (10 * iota) // 1 << 10 = 1,024
		MB                    // 1 << 20 = 1,048,576  (expression inherited!)
		GB                    // 1 << 30 = 1,073,741,824
		TB                    // 1 << 40
	)
	fmt.Printf("\nByte Sizes: KB=%d, MB=%d, GB=%d, TB=%d\n", KB, MB, GB, TB)
	fmt.Printf("5 GB in bytes = %d\n", 5*GB)

	// ── Strings, Bytes, and Runes ──
	fmt.Println("\n=== Strings, Bytes, and Runes ===")
	// Go strings are byte slices, NOT character arrays.
	// ASCII English characters take 1 byte each.
	// Accented letters and emojis take 2 to 4 bytes in UTF-8.
	s := "Café 🌍"
	fmt.Printf("String: %s\n", s)
	fmt.Printf("len() = %d bytes (counts raw bytes, NOT characters!)\n", len(s))

	// range over a string yields RUNES (Unicode code points), not bytes
	runeCount := 0
	for _, r := range s {
		_ = r // r is a rune (int32), representing one Unicode character
		runeCount++
	}
	fmt.Printf("Rune count = %d characters ('C','a','f'=1byte, 'é'=2bytes, ' '=1byte, '🌍'=4bytes)\n", runeCount)

	// byte = uint8 (1 byte, for ASCII and raw binary data)
	// rune = int32 (up to 4 bytes, represents any Unicode character/symbol/emoji)
	fmt.Printf("Type of 'A': %T, Type of 'é': %T, Type of '🌍': %T\n", 'A', 'é', '🌍')
}


