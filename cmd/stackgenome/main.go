package main

import (
	"fmt"
	"os"

	"stackgenome/internal/cli"

	// Import detectors for initialization
	_ "stackgenome/internal/detectors/cicd"
	_ "stackgenome/internal/detectors/editor"
	_ "stackgenome/internal/detectors/infra"
	_ "stackgenome/internal/detectors/language"
	_ "stackgenome/internal/detectors/platform"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
	case "analyze":
		if err := cli.RunAnalyze(args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "version", "--version", "-v":
		fmt.Printf("StackGenome version %s\n", cli.Version)
	case "help", "--help", "-h":
		printUsage()
	case "completion":
		fmt.Println(cli.GetCompletionScript())
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("StackGenome - Local-first repository analyzer and tool recommender.")
	fmt.Println("\nUsage:")
	fmt.Println("  stackgenome <command> [arguments]")
	fmt.Println("\nCommands:")
	fmt.Println("  analyze     Analyze a repository and optionally get recommendations")
	fmt.Println("  completion  Generate shell completion script (e.g., bash/zsh)")
	fmt.Println("  version     Print version information")
	fmt.Println("  help        Show this help message")
	fmt.Println("\nRun 'stackgenome <command> -h' for more information on a command.")
}
