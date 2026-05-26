package observer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func Watch(orch *Orchestrator) {
	fmt.Printf("=== Project-Atra: Observer Watcher Started ===\n")
	fmt.Printf("Watching file: %s\n", orch.Soul.Path)

	orch.Soul.Initialize()

	lastMtime := getMtime(orch.Soul.Path)
	history := []string{}

	for {
		time.Sleep(3 * time.Second)
		currentMtime := getMtime(orch.Soul.Path)

		if currentMtime != lastMtime {
			fmt.Printf("\n[Watcher] Detected changes in logs.\n")
			lastMtime = currentMtime

			lastSpeaker, lastTarget := parseLastEntry(orch.Soul.Path)
			nextPhase := orch.SelectNextPhase(lastSpeaker, lastTarget, history)

			fmt.Printf("[Watcher] Selected Phase: %s\n", nextPhase)
			history = append(history, nextPhase)
			if len(history) > 2 {
				history = history[1:]
			}

			time.Sleep(2 * time.Second)
			err := orch.GenerateResponse(nextPhase)
			if err != nil {
				fmt.Printf("[Watcher] Error: %v\n", err)
			}

			// Update mtime to avoid re-triggering on its own write
			lastMtime = getMtime(orch.Soul.Path)
		}
	}
}

func getMtime(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}

func parseLastEntry(path string) (string, string) {
	file, err := os.Open(path)
	if err != nil {
		return "User", "All"
	}
	defer file.Close()

	var lastHeader string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "## [") {
			lastHeader = line
		}
	}

	if lastHeader == "" {
		return "User", "All"
	}

	// Extract [Speaker -> Target]
	content := lastHeader[strings.Index(lastHeader, "[")+1 : strings.Index(lastHeader, "]")]
	if strings.Contains(content, "->") {
		parts := strings.Split(content, "->")
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(content), "All"
}
