package analyzer

import (
	"project-atra/ast"
	"fmt"
)

// NodeEntity は意味空間における「天体」の実体です。
// 同じ識別子は、場所が違っても同じ実体として扱われます（@参照を含む）。
type NodeEntity struct {
	Name string
	Mass float64
}

// Connection は二点間の「張力（力学的な繋がり）」を表現します。
type Connection struct {
	Source string
	Target string
	Strength float64 // 正の値は引力、負の値は斥力
}

// Space は Atra の「意味空間」全体の状態を保持します。
type Space struct {
	Nodes       map[string]*NodeEntity
	Connections []Connection
}

func NewSpace() *Space {
	return &Space{
		Nodes:       make(map[string]*NodeEntity),
		Connections: []Connection{},
	}
}

// AnalyzeUniverse は AST を巡回し、Space 内に物理パラメータを蓄積します。
func (s *Space) AnalyzeUniverse(u *ast.Universe) {
	for _, stmt := range u.Statements {
		s.analyzeStatement(stmt)
	}
}

func (s *Space) analyzeStatement(stmt ast.Statement) {
	switch v := stmt.(type) {
	case *ast.ExpressionStatement:
		s.analyzeExpression(v.Expression)
	case *ast.FlowStatement:
		s.analyzeFlow(v)
	case *ast.ForceStatement:
		s.analyzeForce(v)
	}
}

func (s *Space) analyzeExpression(exp ast.Expression) {
	switch v := exp.(type) {
	case *ast.CelestialNode:
		s.registerNode(v)
		if v.Body != nil {
			for _, child := range v.Body.Statements {
				s.analyzeStatement(child)
			}
		}
	case *ast.RhizomeExpression:
		s.registerNodeByName(v.Target, 0)
	case *ast.FlowStatement:
		s.analyzeFlow(v)
	case *ast.ForceStatement:
		s.analyzeForce(v)
	}
}

func (s *Space) analyzeFlow(f *ast.FlowStatement) {
	s.analyzeExpression(f.Source)
	s.analyzeExpression(f.Target)
	// Flow 固有の解析（将来の拡張用）
}

func (s *Space) registerNodeByName(name string, mass float64) {
	if _, ok := s.Nodes[name]; !ok {
		s.Nodes[name] = &NodeEntity{Name: name, Mass: 0}
	}
	s.Nodes[name].Mass += mass
}

func (s *Space) registerNode(n *ast.CelestialNode) {
	// 質量の加算 (デフォルト 1.0 + 星の数)
	mass := 1.0 + float64(n.Mass)
	if n.IsSingularity {
		mass += 100.0 // 特異点定数
	}
	s.registerNodeByName(n.Value, mass)
}

func (s *Space) analyzeForce(f *ast.ForceStatement) {
	leftName := getExpressionName(f.Left)
	rightName := getExpressionName(f.Right)

	if leftName == "" || rightName == "" {
		return
	}

	strength := 0.0
	if f.Operator == "+++" {
		strength = 10.0 // 引力定数
	} else if f.Operator == "-!- " || f.Operator == "-!-" {
		strength = -10.0 // 斥力定数
	}

	s.Connections = append(s.Connections, Connection{
		Source:   leftName,
		Target:   rightName,
		Strength: strength,
	})

	// 両方のノードを登録しておく
	s.analyzeExpression(f.Left)
	s.analyzeExpression(f.Right)
}

func getExpressionName(exp ast.Expression) string {
	switch v := exp.(type) {
	case *ast.CelestialNode:
		return v.Value
	case *ast.RhizomeExpression:
		return v.Target
	}
	return ""
}

func (s *Space) Dump() {
	fmt.Println("--- Space Analysis Dump ---")
	fmt.Println("Nodes (Mass):")
	for name, node := range s.Nodes {
		fmt.Printf("  [%s]: %.1f\n", name, node.Mass)
	}
	fmt.Println("Connections (Tension):")
	for _, conn := range s.Connections {
		typeStr := "Attract"
		if conn.Strength < 0 {
			typeStr = "Repel"
		}
		fmt.Printf("  [%s] < %s > [%s] : %.1f\n", conn.Source, typeStr, conn.Target, conn.Strength)
	}
}
