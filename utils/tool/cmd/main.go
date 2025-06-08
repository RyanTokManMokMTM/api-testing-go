package main

import (
	"fmt"
	"log"

	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/workflow"
)

func main() {
	// 創建工作流生成器
	generator := workflow.NewWorkflowGenerator(generator.DefaultOutputDir)

	// 生成所有工作流
	fmt.Println("Generating all workflow test cases...")
	if err := generator.GenerateAllWorkflows(); err != nil {
		log.Fatalf("Failed to generate workflows: %v", err)
	}

	fmt.Println("\nAll test cases have been generated successfully!")
	fmt.Println("Generated files are located in: config/etc/api-test/")
}
