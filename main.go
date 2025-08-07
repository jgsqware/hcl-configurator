package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var (
		outputFile = flag.String("output", "terraform.hcl", "Output HCL file")
	)
	flag.Parse()

	config, err := buildConfigInteractively()
	if err != nil {
		log.Fatalf("Error building config: %v", err)
	}

	hclContent, err := generateHCL(config)
	if err != nil {
		log.Fatalf("Error generating HCL: %v", err)
	}

	err = os.WriteFile(*outputFile, []byte(hclContent), 0644)
	if err != nil {
		log.Fatalf("Error writing HCL file: %v", err)
	}

	fmt.Printf("HCL file generated successfully: %s\n", *outputFile)
}