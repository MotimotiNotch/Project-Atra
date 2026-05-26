package main

import (
	"flag"
	"fmt"
	"os"
	"project-atra/compiler"
	"project-atra/runtime/ollama"
	"project-atra/runtime/observer"
	"project-atra/runtime/experiment"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	subcommand := os.Args[1]
	switch subcommand {
	case "compile":
		compileCmd := flag.NewFlagSet("compile", flag.ExitOnError)
		file := compileCmd.String("file", "./experiments/slime_morphology/mindset.atra", "Path to .atra file")
		out := compileCmd.String("out", "./output/mindset.json", "Output JSON path")
		compileCmd.Parse(os.Args[2:])
		compiler.Compile(*file, *out)

	case "run":
		runCmd := flag.NewFlagSet("run", flag.ExitOnError)
		soul := runCmd.String("soul", "./logs/soul_log.md", "Path to soul log")
		model := runCmd.String("model", "gemma4:latest", "Ollama model name")
		url := runCmd.String("url", "http://localhost:11434/api/generate", "Ollama API URL")
		runCmd.Parse(os.Args[2:])

		client := ollama.NewClient(*url, *model)
		soulMgr := observer.NewSoulManager(*soul)
		orch := observer.NewOrchestrator(client, soulMgr)
		observer.Watch(orch)

	case "evolve":
		evolveCmd := flag.NewFlagSet("evolve", flag.ExitOnError)
		soul := evolveCmd.String("soul", "./logs/long_term_soul.md", "Path to soul log")
		model := evolveCmd.String("model", "gemma4:latest", "Ollama model name")
		url := evolveCmd.String("url", "http://localhost:11434/api/generate", "Ollama API URL")
		evolveCmd.Parse(os.Args[2:])

		client := ollama.NewClient(*url, *model)
		soulMgr := observer.NewSoulManager(*soul)
		orch := observer.NewOrchestrator(client, soulMgr)
		experiment.RunEvolution(client, soulMgr, orch)

	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: atra <subcommand> [flags]")
	fmt.Println("\nSubcommands:")
	fmt.Println("  compile    Compile .atra to .json")
	fmt.Println("  run        Start the observer runtime (watch mode)")
	fmt.Println("  evolve     Run a long-term evolution experiment (100 steps)")
}
