package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultOutputDir is the default output directory
	DefaultOutputDir = "config/etc/api-test/workflows"
)

// Generator is a base generator that provides YAML conversion and saving functionality
type Generator struct {
	outputDir string
}

// NewGenerator creates a new base generator
func NewGenerator(outputDir string) *Generator {
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	return &Generator{
		outputDir: outputDir,
	}
}

// WriteYAML writes data to YAML file
func (g *Generator) WriteYAML(name string, data interface{}) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename, convert to lowercase and replace spaces with underscores
	filename := filepath.Join(g.outputDir, fmt.Sprintf("%s.yaml", strings.ToLower(strings.ReplaceAll(name, " ", "_"))))

	// Convert struct to YAML
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// GetOutputDir gets the output directory
func (g *Generator) GetOutputDir() string {
	return g.outputDir
}
