package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var (
		inputFile  = flag.String("input", "", "Input HCL file to read and edit")
		outputFile = flag.String("output", "services.hcl", "Output HCL file")
	)
	flag.Parse()

	// If no input file specified, try to use the output file as input if it exists
	if *inputFile == "" && *outputFile != "" {
		if _, err := os.Stat(*outputFile); err == nil {
			*inputFile = *outputFile
		}
	}

	// Parse existing HCL file if provided
	config, err := parseHCLFile(*inputFile)
	if err != nil {
		log.Fatalf("Error parsing HCL file: %v", err)
	}

	m := NewModelWithConfig(*outputFile, config)
	p := tea.NewProgram(m, tea.WithAltScreen())
	
	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	if model, ok := finalModel.(Model); ok && model.completed {
		hclContent, err := generateHCL(model.config)
		if err != nil {
			log.Fatalf("Error generating HCL: %v", err)
		}

		err = os.WriteFile(*outputFile, []byte(hclContent), 0644)
		if err != nil {
			log.Fatalf("Error writing HCL file: %v", err)
		}

		fmt.Printf("HCL file generated successfully: %s\n", *outputFile)
	}
}