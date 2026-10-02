// greetpkg demonstrates a reusable package with exported and unexported identifiers.
//
// KEY RULES:
// - This file is in directory "greetpkg/", so it MUST declare "package greetpkg".
// - Exported identifiers (uppercase): accessible from other packages that import this one.
// - Unexported identifiers (lowercase): accessible ONLY within this package.
package greetpkg

import "fmt"

// DefaultGreeting is an exported constant — accessible from any package that imports greetpkg.
const DefaultGreeting = "Hello"

// internalCounter is unexported — only code inside the greetpkg directory can access it.
// This is how Go achieves encapsulation without private/public keywords.
var internalCounter int

// Greet is an exported function (uppercase G).
// Other packages call it as: greetpkg.Greet("Alice")
func Greet(name string) string {
	internalCounter++ // We can access unexported vars from within the same package
	return fmt.Sprintf("%s, %s! (greeting #%d)", DefaultGreeting, name, internalCounter)
}

// formatName is unexported (lowercase f) — it can only be called within this package.
// Attempting to call greetpkg.formatName() from outside will cause a compile error.
func formatName(name string) string {
	return fmt.Sprintf("[%s]", name)
}
