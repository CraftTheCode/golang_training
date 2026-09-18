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
	// Define CLI flags
	uppercase := flag.Bool("upper", false, "Convert output to uppercase")
	repeat := flag.Int("repeat", 1, "Number of times to repeat the output")
	prefix := flag.String("prefix", "[OUTPUT]", "Prefix string before text")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [text]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	var input string

	// Check if data is piped via stdin
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Data is being piped into stdin!
		bytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(string(bytes))
	} else if flag.NArg() > 0 {
		// Positional arguments passed
		input = strings.Join(flag.Args(), " ")
	} else {
		// Interactive prompt if nothing passed
		fmt.Print("Enter text to process: ")
		scanner := bufio.NewScanner(os.Stdin)
		if scanner.Scan() {
			input = scanner.Text()
		}
	}

	if *uppercase {
		input = strings.ToUpper(input)
	}

	// Output to stdout
	for i := 0; i < *repeat; i++ {
		fmt.Printf("%s %s\n", *prefix, input)
	}
}
