package config

import (
	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
)

const (
	// DefaultOutputDir is the default output directory
	DefaultOutputDir = "config/etc/api-test/config"
)

// ConfigGenerator is a configuration generator used to generate shared configurations
type ConfigGenerator struct {
	*generator.Generator
}

// NewConfigGenerator creates a new configuration generator
func NewConfigGenerator(outputDir string) *ConfigGenerator {
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	return &ConfigGenerator{
		Generator: generator.NewGenerator(outputDir),
	}
}

// GenerateConfig generates configuration
func (g *ConfigGenerator) GenerateConfig() error {
	// Create configuration structure
	cfg := config.Config{
		ConfigValues: []config.ConfigValue{
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
