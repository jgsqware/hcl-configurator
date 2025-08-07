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
		outputFile = flag.String("output", "terraform.hcl", "Output HCL file")
	)
	flag.Parse()

	m := NewModel(*outputFile)
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