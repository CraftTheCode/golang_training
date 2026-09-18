package main

import (
	"fmt"
	"golang_training/04-cli-applications/project-gomanager/internal/repository"
	"golang_training/04-cli-applications/project-gomanager/internal/service"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func printUsage() {
	fmt.Fprintf(os.Stderr, `gomanager - Clean Architecture Task Management CLI

Usage:
  gomanager <command> [subcommand] [arguments]

Available Commands:
  init                     Initialize storage and configuration
  config show              Display active configuration
  task create <title>      Create a new operational task
  task list                List all current tasks
  task delete <id>         Delete a task by ID
  help                     Show this help message
`)
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	storageDir := filepath.Join(".", ".gomanager")
	store, err := repository.NewFileStorage(storageDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing storage: %v\n", err)
		os.Exit(1)
	}

	svc := service.NewTaskService(store, store)

	cmd := strings.ToLower(os.Args[1])

	switch cmd {
	case "init":
		if err := svc.Initialize(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Initialized gomanager storage in %s\n", storageDir)

	case "config":
		if len(os.Args) < 3 || os.Args[2] != "show" {
			fmt.Fprintln(os.Stderr, "Usage: gomanager config show")
			os.Exit(1)
		}
		cfg, err := svc.GetConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("=== Active Configuration ===")
		fmt.Printf("Storage Path : %s\n", cfg.StoragePath)
		fmt.Printf("Version      : %s\n", cfg.Version)

	case "task":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: gomanager task [create|list|delete] ...")
			os.Exit(1)
		}
		subCmd := strings.ToLower(os.Args[2])

		switch subCmd {
		case "create":
			if len(os.Args) < 4 {
				fmt.Fprintln(os.Stderr, "Usage: gomanager task create <title>")
				os.Exit(1)
			}
			title := strings.Join(os.Args[3:], " ")
			task, err := svc.CreateTask(title)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating task: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Created Task #%d: %q [Status: %s]\n", task.ID, task.Title, task.Status)

		case "list":
			tasks, err := svc.ListTasks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing tasks: %v\n", err)
				os.Exit(1)
			}
			if len(tasks) == 0 {
				fmt.Println("No tasks found. Use 'gomanager task create <title>' to add one.")
				return
			}
			fmt.Println("----------------------------------------------------------------------")
			fmt.Printf("%-5s | %-12s | %-32s | %s\n", "ID", "STATUS", "TITLE", "CREATED AT")
			fmt.Println("----------------------------------------------------------------------")
			for _, t := range tasks {
				fmt.Printf("%-5d | %-12s | %-32s | %s\n",
					t.ID, t.Status, t.Title, t.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Println("----------------------------------------------------------------------")

		case "delete":
			if len(os.Args) < 4 {
				fmt.Fprintln(os.Stderr, "Usage: gomanager task delete <id>")
				os.Exit(1)
			}
			id, err := strconv.Atoi(os.Args[3])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid task ID %q: must be an integer\n", os.Args[3])
				os.Exit(1)
			}
			if err := svc.DeleteTask(id); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting task: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("✓ Task #%d deleted successfully.\n", id)

		default:
			fmt.Fprintf(os.Stderr, "Unknown task subcommand: %s\n", subCmd)
			os.Exit(1)
		}

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}
