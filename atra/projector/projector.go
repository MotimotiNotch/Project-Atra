package projector

import (
	"project-atra/analyzer"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

// AgentBias は各エージェントの活性率と精神的な偏りを保持します。
type AgentBias struct {
	ActiveRate  float64  `json:"active_rate"`
	Attractions []string `json:"attractions"`
	Repulsions  []string `json:"repulsions"`
}

// MappingConfig は外部から読み込むための設定構造体です。
type MappingConfig struct {
	Agents  []string          `json:"agents"`
	Mapping map[string]string `json:"mapping"`
}

// Mindset は AI エージェント全体の精神構造を表現するデータ構造です。
type Mindset struct {
	Agents map[string]*AgentBias `json:"agents"`
}

// Recipe は各エージェントの活性率の配合を保持します。
type Recipe map[string]float64

// Projector は物理空間をエージェントごとの活性率と精神バイアスに変換します。
// エージェントの数や名前は外部設定（MappingConfig）により自由に変更可能です。
type Projector struct {
	agents  []string
	mapping map[string]string
}

func New(config *MappingConfig) *Projector {
	return &Projector{
		agents:  config.Agents,
		mapping: config.Mapping,
	}
}

// Project は Space を解析してエージェントごとの活性率（Recipe）を生成します。
func (p *Projector) Project(s *analyzer.Space) Recipe {
	scores := make(map[string]float64)

	// 1. ノードの質量に基づくスコアリング
	for name, node := range s.Nodes {
		agent := p.identifyAgent(name)
		scores[agent] += node.Mass
	}

	// 2. 張力（接続）に基づく増幅
	for _, conn := range s.Connections {
		srcAgent := p.identifyAgent(conn.Source)
		tgtAgent := p.identifyAgent(conn.Target)

		if srcAgent == tgtAgent {
			scores[srcAgent] += math.Abs(conn.Strength) * 0.5
		} else {
			if conn.Strength > 0 {
				scores[srcAgent] += s.Nodes[conn.Target].Mass * (conn.Strength / 10.0)
			}
		}
	}

	return p.normalize(scores)
}

// CreateMindset は物理的な接続をエージェントごとの精神構造として抽出します。
func (p *Projector) CreateMindset(s *analyzer.Space) *Mindset {
	recipe := p.Project(s)
	mindset := &Mindset{
		Agents: make(map[string]*AgentBias),
	}

	for _, a := range p.agents {
		mindset.Agents[a] = &AgentBias{
			ActiveRate:  recipe[a],
			Attractions: []string{},
			Repulsions:  []string{},
		}
	}

	for _, conn := range s.Connections {
		srcAgent := p.identifyAgent(conn.Source)
		
		if _, ok := mindset.Agents[srcAgent]; !ok {
			continue
		}

		relation := fmt.Sprintf("%s(%.1f)", conn.Target, conn.Strength)

		if conn.Strength > 0 {
			mindset.Agents[srcAgent].Attractions = append(mindset.Agents[srcAgent].Attractions, relation)
		} else {
			mindset.Agents[srcAgent].Repulsions = append(mindset.Agents[srcAgent].Repulsions, relation)
		}
	}

	return mindset
}

func (p *Projector) identifyAgent(name string) string {
	for k, v := range p.mapping {
		if strings.EqualFold(k, name) {
			return v
		}
	}
	if len(p.agents) > 0 {
		// マッピングにない場合は、リストの最後（通常は統合役を想定）を返します。
		return p.agents[len(p.agents)-1]
	}
	return "unknown"
}

func (p *Projector) normalize(scores map[string]float64) Recipe {
	total := 0.0
	for _, s := range scores {
		total += s
	}

	recipe := make(Recipe)

	if total == 0 {
		for _, a := range p.agents {
			recipe[a] = 100.0 / float64(len(p.agents))
		}
		return recipe
	}

	for _, a := range p.agents {
		recipe[a] = math.Round((scores[a]/total)*1000) / 10.0
	}

	return recipe
}

func (m *Mindset) DumpJSON() string {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(b)
}

func (m *Mindset) ExportEnglishPrompt(agentName string) string {
	bias, ok := m.Agents[strings.ToLower(agentName)]
	if !ok {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## INTERNAL GRAVITATIONAL FIELD (MINDSET)\n")
	sb.WriteString(fmt.Sprintf("You are currently operating with an Active Rate of %.1f%% and the following semantic biases:\n\n", bias.ActiveRate))

	if len(bias.Attractions) > 0 {
		sb.WriteString("### ATTRACTIONS (Thinking toward)\n")
		for _, a := range bias.Attractions {
			sb.WriteString(fmt.Sprintf("- %s: High attraction. Naturally gravitate toward this concept.\n", a))
		}
		sb.WriteString("\n")
	}

	if len(bias.Repulsions) > 0 {
		sb.WriteString("### REPULSIONS (Avoiding)\n")
		for _, r := range bias.Repulsions {
			sb.WriteString(fmt.Sprintf("- %s: Cognitive dissonance. Avoid this concept or express discomfort.\n", r))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (r Recipe) Dump() {
	fmt.Println("--- Agent Recipe Projection ---")
	keys := make([]string, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return r[keys[i]] > r[keys[j]]
	})

	for _, k := range keys {
		fmt.Printf("  %-8s: %.1f%%\n", k, r[k])
	}
}
