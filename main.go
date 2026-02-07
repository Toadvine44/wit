package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "add":
		cmdAdd(args)
	case "commit":
		cmdCommit(args)
	case "status":
		cmdStatus(args)
	case "branch":
		cmdBranch(args)
	case "checkout":
		cmdCheckout(args)
	case "diff":
		cmdDiff(args)
	case "help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "wit: '%s' is not a wit command\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`wit - a personalized version control system

Usage: wit <command> [<args>]

Commands:
  add       Add file contents to the index
  commit    Record changes to the repository
  status    Show the working tree status
  branch    List, create, or delete branches
  checkout  Switch branches or restore working tree files
  diff      Show changes between commits, commit and working tree, etc.
  help      Show this help message`)
}

func cmdAdd(args []string) {
	if len(args) == 0 {
		fmt.Println("Nothing specified, nothing added.")
		return
	}
	fmt.Printf("Adding files: %v\n", args)
	// TODO: implement add logic
}

func cmdCommit(args []string) {
	message := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "-m" && i+1 < len(args) {
			message = args[i+1]
			break
		}
	}
	if message == "" {
		fmt.Println("wit commit: no commit message provided (use -m)")
		return
	}
	fmt.Printf("Committing with message: %s\n", message)
	// TODO: implement commit logic
}

func cmdStatus(args []string) {
	fmt.Println("On branch main")
	fmt.Println("nothing to commit, working tree clean")
	// TODO: implement status logic
}

func cmdBranch(args []string) {
	if len(args) == 0 {
		fmt.Println("* main")
		// TODO: list branches
		return
	}
	fmt.Printf("Creating branch: %s\n", args[0])
	// TODO: implement branch creation
}

func cmdCheckout(args []string) {
	if len(args) == 0 {
		fmt.Println("wit checkout: no branch or path specified")
		return
	}
	fmt.Printf("Switching to: %s\n", args[0])
	// TODO: implement checkout logic
}

func cmdDiff(args []string) {
	fmt.Println("No changes detected")
	// TODO: implement diff logic
}
