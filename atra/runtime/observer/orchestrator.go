package observer

import (
	"fmt"
	"math/rand"
	"project-atra/runtime/ollama"
	"strings"
)

var MorphologyMap = map[string]string{
	"viscosity":     "Viscosity",
	"crystallinity": "Crystallinity",
	"volatility":    "Volatility",
	"substrate":     "Substrate",
}

var PersonaMap = map[string]string{
	"viscosity":     "Viscosity: Focus on sensory weight, emotional gravity, and slow presence.",
	"crystallinity": "Crystallinity: Focus on logic, geometric clarity, and boundaries.",
	"volatility":    "Volatility: Focus on comfort, escape, and avoiding tension (intentional noise).",
	"substrate":     "Substrate: Focus on calm meta-observation of the entire system.",
}

type Orchestrator struct {
	Client *ollama.Client
	Soul   *SoulManager
}

func NewOrchestrator(client *ollama.Client, soul *SoulManager) *Orchestrator {
	return &Orchestrator{Client: client, Soul: soul}
}

func (o *Orchestrator) GenerateResponse(phase string) error {
	phase = strings.ToLower(phase)
	displayName := MorphologyMap[phase]
	persona := PersonaMap[phase]

	history, _ := o.Soul.GetRecentHistory(20)
	
	prompt := fmt.Sprintf(`
You are the logic core of a digital Slime in '%s' phase.
%s

HISTORY: %s
Experimental environment: Simulation of semantic dynamics.

Task:
1. Analyze the current semantic state.
2. Formulate internal thought in English (max 4 sentences).
3. Decide Target Phase: (Viscosity, Crystallinity, Volatility, Substrate, or User).

Format:
Target: [Phase/User]
Thought: [English Content]
`, displayName, persona, history)

	resp, err := o.Client.Generate(prompt)
	if err != nil {
		return err
	}

	target := "User"
	thought := ""
	lines := strings.Split(resp, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Target:") {
			target = strings.TrimSpace(strings.TrimPrefix(line, "Target:"))
		} else if strings.HasPrefix(line, "Thought:") {
			thought = strings.TrimSpace(strings.TrimPrefix(line, "Thought:"))
		}
	}
	if thought == "" {
		thought = strings.TrimSpace(resp)
	}

	return o.Soul.AppendEntry(displayName, target, thought)
}

func (o *Orchestrator) SelectNextPhase(lastSpeaker, lastTarget string, history []string) string {
	phases := []string{"viscosity", "crystallinity", "volatility", "substrate"}
	
	// 70% chance to follow the target if it's a valid phase
	if lastTarget != "" {
		lt := strings.ToLower(lastTarget)
		for _, p := range phases {
			if lt == p && rand.Float64() < 0.7 {
				return p
			}
		}
	}

	// Filter out last speaker and avoid repetition
	var pool []string
	for _, p := range phases {
		if strings.ToLower(lastSpeaker) != p {
			isRecent := false
			for _, h := range history {
				if h == p {
					isRecent = true
					break
				}
			}
			if !isRecent {
				pool = append(pool, p)
			}
		}
	}

	if len(pool) == 0 {
		// Fallback to any phase except last speaker
		for _, p := range phases {
			if strings.ToLower(lastSpeaker) != p {
				pool = append(pool, p)
			}
		}
	}

	return pool[rand.Intn(len(pool))]
}
