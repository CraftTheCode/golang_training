// 04_packages_and_modules demonstrates Go's package system:
//   - How to import and use packages from other directories
//   - Exported (uppercase) vs unexported (lowercase) visibility
//   - The "1 directory = 1 package" rule
//
// Run this example from the repository root:
//
//	go run ./03-idiomatic-go-design/04_packages_and_modules/
package main

import (
	"fmt"

	// Import using the module path (from go.mod) + directory path.
	// The module is "golang_training", and the package directory is
	// "03-idiomatic-go-design/04_packages_and_modules/greetpkg".
	"golang_training/03-idiomatic-go-design/04_packages_and_modules/greetpkg"
)

func main() {
	fmt.Println("=== 1. Calling Exported Functions from Another Package ===")
	// greetpkg.Greet is exported (uppercase G) — we can call it from here.
	msg := greetpkg.Greet("Divya")
	fmt.Println(msg)

	msg2 := greetpkg.Greet("Gopher")
	fmt.Println(msg2)

	fmt.Println("\n=== 2. Accessing Exported Constants ===")
	// greetpkg.DefaultGreeting is exported — we can read it.
	fmt.Println("Default greeting from package:", greetpkg.DefaultGreeting)

	fmt.Println("\n=== 3. Unexported Identifiers Are Hidden ===")
	// The following lines would cause COMPILE ERRORS if uncommented:
	//
	// greetpkg.formatName("test")      // ERROR: formatName is unexported (lowercase f)
	// greetpkg.internalCounter         // ERROR: internalCounter is unexported (lowercase i)
	//
	// This is Go's encapsulation: the package controls what is part of its public API.
	fmt.Println("greetpkg.formatName() and greetpkg.internalCounter are not accessible here.")
	fmt.Println("Only uppercase identifiers cross package boundaries.")

	fmt.Println("\n=== 4. Key Takeaways ===")
	fmt.Println("• 1 directory = 1 package (enforced by the Go toolchain)")
	fmt.Println("• Import path = module name + directory path")
	fmt.Println("• Uppercase = exported (public), lowercase = unexported (private)")
	fmt.Println("• Files in the same directory share everything without imports")
}
