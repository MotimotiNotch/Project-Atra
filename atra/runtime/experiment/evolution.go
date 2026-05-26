package experiment

import (
	"fmt"
	"project-atra/runtime/ollama"
	"project-atra/runtime/observer"
	"strings"
	"time"
)

func RunEvolution(client *ollama.Client, soul *observer.SoulManager, orch *observer.Orchestrator) {
	fmt.Printf("--- STARTING LONG-TERM EVOLUTION EXPERIMENT (100 STEPS) ---\n")
	
	soul.Initialize()

	for i := 1; i <= 100; i++ {
		fmt.Printf(">>> Step %d/100\n", i)

		history, _ := soul.GetRecentHistory(10)
		
		genPrompt := fmt.Sprintf(`
You are the 'Observer' of a digital lifeform. 
Your task is to provide ONE distinct, short stimulus (question or comment) in English.
Explore deep themes: the nature of boundaries, the difference between 'knowing' and 'being', or sensory echoes.
Output ONLY the text. No quotes.

HISTORY:
%s
`, history)

		stimulus, err := client.Generate(genPrompt)
		if err != nil {
			fmt.Printf("Ollama Error: %v\n", err)
			continue
		}
		stimulus = strings.TrimPrefix(stimulus, "Observer:")
		stimulus = strings.TrimPrefix(stimulus, "User:")
		stimulus = strings.Trim(stimulus, "\"")

		fmt.Printf("Stimulus: %s\n", stimulus)
		
		err = soul.AppendEntry("Observer", "All", stimulus)
		if err != nil {
			fmt.Printf("Write Error: %v\n", err)
			continue
		}

		// Selection logic for response
		err = orch.GenerateResponse("viscosity") // Start with viscosity as a baseline or cycle through
		if err != nil {
			fmt.Printf("Response Error: %v\n", err)
		}

		time.Sleep(5 * time.Second)
	}

	fmt.Printf("\n=== LONG-TERM EXPERIMENT COMPLETE. ===\n")
}
