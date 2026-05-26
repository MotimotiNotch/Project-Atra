package parser

import (
	"project-atra/ast"
	"project-atra/lexer"
	"testing"
)

func TestCelestialNodes(t *testing.T) {
	input := `
*star
!singularity
***Viscosity
`
	l := lexer.New(input)
	p := New(l)
	universe := p.ParseUniverse()
	checkParserErrors(t, p)

	if universe == nil {
		t.Fatal("ParseUniverse() returned nil")
	}
	if len(universe.Statements) != 3 {
		t.Fatalf("universe.Statements does not contain 3 statements. got=%d", len(universe.Statements))
	}

	tests := []struct {
		expectedMass          int
		expectedIsSingularity bool
		expectedValue         string
	}{
		{1, false, "star"},
		{0, true, "singularity"},
		{3, false, "Viscosity"},
	}

	for i, tt := range tests {
		stmt := universe.Statements[i]
		// CelestialNode は式として解析されることを想定
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		node, ok := expStmt.Expression.(*ast.CelestialNode)
		if !ok {
			t.Fatalf("expStmt.Expression is not *ast.CelestialNode. got=%T", expStmt.Expression)
		}

		if node.Mass != tt.expectedMass {
			t.Errorf("test[%d] - mass wrong. expected=%d, got=%d", i, tt.expectedMass, node.Mass)
		}
		if node.IsSingularity != tt.expectedIsSingularity {
			t.Errorf("test[%d] - isSingularity wrong. expected=%t, got=%t", i, tt.expectedIsSingularity, node.IsSingularity)
		}
		if node.Value != tt.expectedValue {
			t.Errorf("test[%d] - value wrong. expected=%s, got=%s", i, tt.expectedValue, node.Value)
		}
	}
}

func TestFlowStatements(t *testing.T) {
	input := `
self -> world
dream ~> memory
void => truth
`
	l := lexer.New(input)
	p := New(l)
	universe := p.ParseUniverse()
	checkParserErrors(t, p)

	if len(universe.Statements) != 3 {
		t.Fatalf("universe.Statements does not contain 3 statements. got=%d", len(universe.Statements))
	}

	tests := []struct {
		expectedSource   string
		expectedOperator string
		expectedTarget   string
	}{
		{"self", "->", "world"},
		{"dream", "~>", "memory"},
		{"void", "=>", "truth"},
	}

	for i, tt := range tests {
		stmt := universe.Statements[i]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		flow, ok := expStmt.Expression.(*ast.FlowStatement)
		if !ok {
			t.Fatalf("expStmt.Expression is not *ast.FlowStatement. got=%T", expStmt.Expression)
		}

		if flow.Source.String() != tt.expectedSource {
			t.Errorf("test[%d] - source wrong. expected=%s, got=%s", i, tt.expectedSource, flow.Source.String())
		}
		if flow.Operator != tt.expectedOperator {
			t.Errorf("test[%d] - operator wrong. expected=%s, got=%s", i, tt.expectedOperator, flow.Operator)
		}
		if flow.Target.String() != tt.expectedTarget {
			t.Errorf("test[%d] - target wrong. expected=%s, got=%s", i, tt.expectedTarget, flow.Target.String())
		}
	}
}

func TestForceStatements(t *testing.T) {
	input := `
Viscosity +++ Silent
Volatility -!- social
`
	l := lexer.New(input)
	p := New(l)
	universe := p.ParseUniverse()
	checkParserErrors(t, p)

	if len(universe.Statements) != 2 {
		t.Fatalf("universe.Statements does not contain 2 statements. got=%d", len(universe.Statements))
	}

	tests := []struct {
		expectedLeft     string
		expectedOperator string
		expectedRight    string
	}{
		{"Viscosity", "+++", "Silent"},
		{"Volatility", "-!-", "social"},
	}

	for i, tt := range tests {
		stmt := universe.Statements[i]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		force, ok := expStmt.Expression.(*ast.ForceStatement)
		if !ok {
			t.Fatalf("expStmt.Expression is not *ast.ForceStatement. got=%T", expStmt.Expression)
		}

		if force.Left.String() != tt.expectedLeft {
			t.Errorf("test[%d] - left wrong. expected=%s, got=%s", i, tt.expectedLeft, force.Left.String())
		}
		if force.Operator != tt.expectedOperator {
			t.Errorf("test[%d] - operator wrong. expected=%s, got=%s", i, tt.expectedOperator, force.Operator)
		}
		if force.Right.String() != tt.expectedRight {
			t.Errorf("test[%d] - right wrong. expected=%s, got=%s", i, tt.expectedRight, force.Right.String())
		}
	}
}

func TestRhizomeExpressions(t *testing.T) {
	input := `
@memory
@dream
`
	l := lexer.New(input)
	p := New(l)
	universe := p.ParseUniverse()
	checkParserErrors(t, p)

	if len(universe.Statements) != 2 {
		t.Fatalf("universe.Statements does not contain 2 statements. got=%d", len(universe.Statements))
	}

	tests := []struct {
		expectedTarget string
	}{
		{"memory"},
		{"dream"},
	}

	for i, tt := range tests {
		stmt := universe.Statements[i]
		expStmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", stmt)
		}

		ref, ok := expStmt.Expression.(*ast.RhizomeExpression)
		if !ok {
			t.Fatalf("expStmt.Expression is not *ast.RhizomeExpression. got=%T", expStmt.Expression)
		}

		if ref.Target != tt.expectedTarget {
			t.Errorf("test[%d] - target wrong. expected=%s, got=%s", i, tt.expectedTarget, ref.Target)
		}
	}
}

func TestBlockStatements(t *testing.T) {
	input := `
Viscosity:
  Silent
  *star
`
	l := lexer.New(input)
	p := New(l)
	universe := p.ParseUniverse()
	checkParserErrors(t, p)

	if len(universe.Statements) != 1 {
		t.Fatalf("universe.Statements does not contain 1 statement. got=%d", len(universe.Statements))
	}

	stmt, ok := universe.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not *ast.ExpressionStatement. got=%T", universe.Statements[0])
	}

	node, ok := stmt.Expression.(*ast.CelestialNode)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.CelestialNode. got=%T", stmt.Expression)
	}

	if node.Value != "Viscosity" {
		t.Errorf("node.Value wrong. expected=Viscosity, got=%s", node.Value)
	}

	if node.Body == nil {
		t.Fatal("node.Body is nil")
	}

	if len(node.Body.Statements) != 2 {
		t.Fatalf("node.Body.Statements does not contain 2 statements. got=%d", len(node.Body.Statements))
	}

	expectedInBody := []string{"Silent", "*star"}
	for i, val := range expectedInBody {
		if node.Body.Statements[i].String() != val {
			t.Errorf("node.Body.Statements[%d] wrong. expected=%s, got=%s", i, val, node.Body.Statements[i].String())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
