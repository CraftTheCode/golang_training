// 01_flags_and_io demonstrates production CLI patterns in Go:
//   - Command-line flag parsing with the standard library `flag` package
//   - Why flags return pointers (*bool, *int, *string)
//   - How to detect if input is piped (e.g., `echo "hi" | go run .`) vs typed interactively
//   - Stream separation: stdout for output, stderr for logs and usage errors
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// ── 1. Defining Flags ──
	// Note: flag.Bool(), flag.Int(), flag.String() return POINTERS (*bool, *int, *string), NOT direct values!
	// WHY POINTERS?
	// Because flag declarations only REGISTER what flags exist.
	// The actual command-line arguments are NOT read until flag.Parse() is executed below.
	// Returning pointers allows the flag package to populate the values in-place when Parse() runs.
	uppercase := flag.Bool("upper", false, "Convert output to uppercase")
	repeat := flag.Int("repeat", 1, "Number of times to repeat the output")
	prefix := flag.String("prefix", "[OUTPUT]", "Prefix string before text")

	// Custom usage message shown when user runs with -h or provides invalid flags
	flag.Usage = func() {
		// Rule: Diagnostic messages and help text ALWAYS go to os.Stderr, not os.Stdout
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [text]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	// flag.Parse() processes os.Args[1:] up to the first non-flag argument.
	// After this line, *uppercase, *repeat, and *prefix contain the user's values.
	flag.Parse()

	var input string

	// ── 2. Detecting Piped Input (Terminal vs Unix Pipe) ──
	// In Unix/Windows, standard input (os.Stdin) can either be an interactive keyboard (Character Device)
	// or a stream redirected from a file or pipe (e.g., `cat file.txt | mycli`).
	//
	// os.Stdin.Stat() returns FileInfo describing the input descriptor.
	// `stat.Mode() & os.ModeCharDevice` performs a bitwise test:
	//   - If non-zero: Stdin is an interactive terminal character device (keyboard).
	//   - If == 0: Stdin is a pipe or redirected file!
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Input is being piped into stdin!
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(string(bytes))
	} else if flag.NArg() > 0 {
		// ── 3. Positional (Non-Flag) Arguments ──
		// flag.NArg() returns the count of non-flag arguments remaining after flags.
		// flag.Args() returns []string of those remaining arguments.
		// Example: `mycli -upper=true hello world` -> flag.Args() = ["hello", "world"]
		input = strings.Join(flag.Args(), " ")
	} else {
		// ── 4. Interactive Keyboard Prompt Fallback ──
		fmt.Print("Enter text to process: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}
	}

	// Dereference the pointer with `*uppercase` to read the boolean value
	if *uppercase {
		input = strings.ToUpper(input)
	}

	// ── 5. Standard Output ──
	// Normal program output goes strictly to os.Stdout (fmt.Printf defaults to stdout)
	for i := 0; i < *repeat; i++ {
		fmt.Printf("%s %s\n", *prefix, input)
	}
}

