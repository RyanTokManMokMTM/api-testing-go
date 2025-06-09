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
			{
				Name:  "mid",
				Value: "6552f99b99821c568c0115cc",
			},
			{
				Name:  "legacy_id",
				Value: "SL101PRO6828965740321448017_SKU6828965740942204962",
			},
			{
				Name:  "config_legacy_id",
				Value: "SL101PRO4409746819723214928_SL101SKU4409746892519555080",
			},
			{
				Name:    "next_day",
				Command: "echo $(date -d \"tomorrow\" +%s%3N)",
				Type:    "number",
			},
			{
				Name:    "next_month",
				Command: "echo $(date -d \"+1 month\" +%s%3N)",
				Type:    "number",
			},
		},
	}

	// Write configuration file
	return g.Generator.WriteYAML("config", &cfg)
}
