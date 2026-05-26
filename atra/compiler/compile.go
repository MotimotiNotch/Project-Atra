package compiler

import (
	"project-atra/analyzer"
	"project-atra/lexer"
	"project-atra/parser"
	"project-atra/projector"
	"fmt"
	"os"
	"path/filepath"
)

func Compile(atraPath, outputPath string) {
	// Read Atra Source
	content, err := os.ReadFile(atraPath)
	if err != nil {
		fmt.Printf("Failed to read atra file: %v\n", err)
		os.Exit(1)
	}

	input := string(content)

	l := lexer.New(input)
	p := parser.New(l)
	universe := p.ParseUniverse()

	if len(p.Errors()) != 0 {
		fmt.Printf("Parser errors:\n")
		for _, msg := range p.Errors() {
			fmt.Printf("\t%s\n", msg)
		}
		os.Exit(1)
	}

	fmt.Println("--- Atra Universe Observation ---")
	fmt.Printf("Source File: %s\n", atraPath)
	
	fmt.Println("\nParsed Structure (Graph View):")
	for _, stmt := range universe.Statements {
		printStatement(stmt, 0)
	}

	// Structural Analysis
	space := analyzer.NewSpace()
	space.AnalyzeUniverse(universe)

	fmt.Println()
	space.Dump()

	// Generic Phase Mapping
	config := &projector.MappingConfig{
		Agents: []string{"viscosity", "crystallinity", "volatility", "substrate"},
		Mapping: map[string]string{
			"Viscosity":        "viscosity",
			"Body":             "viscosity",
			"Emotion":          "viscosity",
			"Visceral":         "viscosity",
			"Crystallinity":    "crystallinity",
			"Structure":        "crystallinity",
			"Logic":            "crystallinity",
			"Boundary":         "crystallinity",
			"Volatility":       "volatility",
			"Play":             "volatility",
			"Buffer":           "volatility",
			"Noise":            "volatility",
			"Substrate":        "substrate",
			"Meta":             "substrate",
			"Universe":         "substrate",
			"Space":            "substrate",
			"Flow":             "substrate",
		},
	}

	proj := projector.New(config)
	mindset := proj.CreateMindset(space)
	jsonStr := mindset.DumpJSON()

	// Ensure output directory exists
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
	}

	err = os.WriteFile(outputPath, []byte(jsonStr), 0644)
	if err != nil {
		fmt.Printf("Failed to save mindset: %v\n", err)
	} else {
		fmt.Printf("\nMindset saved to: %s\n", outputPath)
	}

	recipe := proj.Project(space)
	fmt.Println()
	recipe.Dump()
}

func printStatement(stmt interface{}, indent int) {
	indentStr := ""
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	switch s := stmt.(type) {
	default:
		fmt.Printf("%s%s\n", indentStr, s)
	}
}
