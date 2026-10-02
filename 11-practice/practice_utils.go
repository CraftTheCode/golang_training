package main

import (
	"fmt"
	"golang_training/12-practice-package"
)

func Experiment() {
	// ❌ testPrivate — lowercase type, CANNOT be accessed from outside the package
	// tprivate := practice.testPrivate{} // This would cause a compile error

	// ✅ TestPublic — uppercase type, exported. All fields are uppercase too, so all accessible.
	tpublic := practice.TestPublic{
		Name: "public",
		Age:  30,
	}
	fmt.Println("TestPublic:", tpublic)

	// ✅ TestMix — uppercase type, exported. But 'age' is lowercase (unexported).
	// We can only set the exported field 'Name', not the unexported 'age'.
	tmix := practice.TestMix{
		Name: "mix",
		// age: 25, // ❌ This would fail — 'age' is unexported
	}
	fmt.Println("TestMix:", tmix)
}
