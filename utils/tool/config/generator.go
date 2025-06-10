// Package config provides configuration generators for shared test configurations.
package config

import (
	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
)

const (
	// DefaultOutputDir is the default output directory
	DefaultOutputDir = "config/etc/api-test/config"
)

// Generator is a configuration generator used to generate shared configurations.
type Generator struct {
	*generator.Generator
}

// NewGenerator creates a new configuration generator with the specified output directory.
func NewGenerator(outputDir string) *Generator {
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	return &Generator{
		Generator: generator.NewGenerator(outputDir),
	}
}

// GenerateConfig generates the configuration file using the provided values.
func (g *Generator) GenerateConfig() error {
	// Create configuration structure
	cfg := config.Config{
		ConfigValues: []config.Value{
			// Merchant and business related IDs
			{
				Name:  "mid",
				Value: "6552f99b99821c568c0115cc",
			},
		},
	}

	// Write configuration file
	return g.WriteYAML("config", &cfg)
}
