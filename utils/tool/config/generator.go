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
			{
				Name:  "legacy_id",
				Value: "SL101PRO6828965740321448017_SKU6828965740942204962",
			},
			{
				Name:  "config_legacy_id",
				Value: "SL101PRO4409746819723214928_SL101SKU4409746892519555080",
			},
			{
				Name:  "test_merchant_id",
				Value: "test_merchant_12345",
			},
			{
				Name:  "test_user_id",
				Value: "test_user_67890",
			},

			// Time-based values
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
			{
				Name:    "current_timestamp",
				Command: "echo $(date +%s%3N)",
				Type:    "number",
			},
			{
				Name:    "yesterday_timestamp",
				Command: "echo $(date -d \"yesterday\" +%s%3N)",
				Type:    "number",
			},

			// Test data
			{
				Name:  "test_email",
				Value: "test@example.com",
			},
			{
				Name:  "test_phone",
				Value: "+886912345678",
			},
			{
				Name:  "test_name",
				Value: "Test User",
			},
			{
				Name:  "test_address",
				Value: "123 Test Street, Taipei, Taiwan",
			},

			// API endpoints
			{
				Name:  "api_base_url",
				Value: "https://api.example.com",
			},
			{
				Name:  "faker_api_url",
				Value: "https://fakerapi.it/api/v2",
			},

			// Currency and pricing
			{
				Name:  "default_currency",
				Value: "TWD",
			},
			{
				Name:  "test_price_cents",
				Value: 1000,
			},
			{
				Name:  "test_price_usd_cents",
				Value: 50,
			},

			// Authentication
			{
				Name:  "test_api_key",
				Value: "test_api_key_12345",
			},
			{
				Name:  "test_bearer_token",
				Value: "Bearer test_token_67890",
			},

			// Product and SKU related
			{
				Name:  "test_product_id",
				Value: "PROD_001",
			},
			{
				Name:  "test_sku_id",
				Value: "SKU_001",
			},
			{
				Name:  "test_category_id",
				Value: "CAT_001",
			},

			// Order and subscription related
			{
				Name:  "test_order_id",
				Value: "ORDER_001",
			},
			{
				Name:  "test_subscription_id",
				Value: "SUB_001",
			},
			{
				Name:  "test_payment_id",
				Value: "PAY_001",
			},

			// Status values
			{
				Name:  "status_active",
				Value: "active",
			},
			{
				Name:  "status_pending",
				Value: "pending",
			},
			{
				Name:  "status_completed",
				Value: "completed",
			},
			{
				Name:  "status_cancelled",
				Value: "cancelled",
			},

			// Locale and language
			{
				Name:  "locale_en",
				Value: "en_US",
			},
			{
				Name:  "locale_zh",
				Value: "zh_TW",
			},
			{
				Name:  "locale_ja",
				Value: "ja_JP",
			},

			// HTTP status codes
			{
				Name:  "status_ok",
				Value: 200,
			},
			{
				Name:  "status_created",
				Value: 201,
			},
			{
				Name:  "status_bad_request",
				Value: 400,
			},
			{
				Name:  "status_unauthorized",
				Value: 401,
			},
			{
				Name:  "status_not_found",
				Value: 404,
			},
			{
				Name:  "status_internal_error",
				Value: 500,
			},

			// Test metadata
			{
				Name:  "test_metadata_key",
				Value: "test_key",
			},
			{
				Name:  "test_metadata_value",
				Value: "test_value",
			},
			{
				Name:  "test_description",
				Value: "This is a test description for API testing",
			},
		},
	}

	// Write configuration file
	return g.Generator.WriteYAML("config", &cfg)
}
