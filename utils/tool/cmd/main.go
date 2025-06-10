// Package main provides the command-line interface for generating API test configurations and workflows.
package main

import (
	"fmt"
	"log"

	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/config"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/workflow"
)

func main() {
	// Generate configuration
	fmt.Println("Generating config...")
	configGenerator := config.NewGenerator(config.DefaultOutputDir)
	if err := configGenerator.GenerateConfig(); err != nil {
		log.Fatalf("Failed to generate config: %v", err)
	}

	// Generate workflow
	fmt.Println("\nGenerating all workflow test cases...")
	workflowGenerator := workflow.NewGenerator(generator.DefaultOutputDir)
	if err := workflowGenerator.GenerateAllWorkflows(); err != nil {
		log.Fatalf("Failed to generate workflows: %v", err)
	}

	fmt.Println("\nAll test cases have been generated successfully!")
	fmt.Println("Generated files are located in: config/etc/api-test/")
}
