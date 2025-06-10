package config

import (
	"fmt"

	"github.com/RyanTokManMokMTM/api-testing-go/config/types"
)

// Config represents the global configuration
type Config struct {
	ConfigValues []ConfigValue `yaml:"config" validate:"required"`
}

// ConfigValue represents a single configuration value
type ConfigValue struct {
	Name    string          `yaml:"name" validate:"required"`
	Value   any             `yaml:"value,omitempty"`
	Command string          `yaml:"command,omitempty"`
	Type    types.FieldType `yaml:"type,omitempty"`
}

// APITesting represents a complete API test suite
type APITesting struct {
	APITest APITest `yaml:"api_test" validate:"required"`
}

// APITest represents a single API test configuration
type APITest struct {
	Host          string         `yaml:"host" validate:"required"`
	NetworkEnable bool           `yaml:"network_enable"`
	Name          string         `yaml:"name" validate:"required"`
	Config        map[string]any `yaml:"-"`
	GlobalHook    Hook           `yaml:"global_hook,omitempty"`
	Skip          bool           `yaml:"skip"`
	Scenarios     []Scenario     `yaml:"scenarios" validate:"required,dive"`
}

// Scenario represents a test scenario
type Scenario struct {
	Skip      bool       `yaml:"skip"`
	Name      string     `yaml:"name" validate:"required"`
	Hook      Hook       `yaml:"hook,omitempty"`
	Workflows []Workflow `yaml:"workflows" validate:"required,dive"`
}

// Hook represents before/after hooks
type Hook struct {
	Before HookActions `yaml:"before,omitempty"`
	After  HookActions `yaml:"after,omitempty"`
}

// HookActions represents actions to be taken in hooks
type HookActions struct {
	InitVars  []InitVar  `yaml:"init_vars,omitempty"`
	Workflows []Workflow `yaml:"workflows,omitempty"`
}

// InitVar represents a variable initialization
type InitVar struct {
	Name     string `yaml:"name" validate:"required"`
	FromStep string `yaml:"from_step" validate:"required"`
	Field    string `yaml:"field" validate:"required"`
}

// Workflow represents a single API test workflow
type Workflow struct {
	Step           string   `yaml:"step" validate:"required"`
	Request        Request  `yaml:"request" validate:"required"`
	ExpectResponse Response `yaml:"expect_response" validate:"required"`
}

// Request represents an API request
type Request struct {
	Vars         []Var          `yaml:"vars,omitempty"`
	FromResponse []FromResponse `yaml:"from_response,omitempty"`
	Method       string         `yaml:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	URI          string         `yaml:"uri" validate:"required"`
	Headers      string         `yaml:"headers,omitempty"`
	Query        string         `yaml:"query,omitempty"`
	Body         string         `yaml:"body,omitempty"`
}

// Response represents expected API response
type Response struct {
	StatusCode int       `yaml:"status_code" validate:"required"`
	Body       BodyCheck `yaml:"body,omitempty"`
}

// BodyCheck represents response body validation
type BodyCheck struct {
	Equals       []EqualsCheck      `yaml:"equals,omitempty"`
	GreaterThans []GreaterThanCheck `yaml:"greater_thans,omitempty"`
	LessThans    []LessThanCheck    `yaml:"less_thans,omitempty"`
	Matches      []MatchesCheck     `yaml:"matches,omitempty"`
	Presents     []PresentCheck     `yaml:"presents,omitempty"`
	NotPresents  []NotPresentCheck  `yaml:"not_presents,omitempty"`
}

// Var represents a variable definition
type Var struct {
	Name string `yaml:"name" validate:"required"`
}

// FromResponse represents a response field reference
type FromResponse struct {
	Step      string `yaml:"step,omitempty"`
	Name      string `yaml:"name" validate:"required"`
	FromField string `yaml:"from_field" validate:"required"`
}

// EqualsCheck represents an equality check
type EqualsCheck struct {
	Field string          `yaml:"field" validate:"required"`
	Type  types.FieldType `yaml:"type,omitempty"`
	Value any             `yaml:"value" validate:"required"`
}

// GreaterThanCheck represents a greater than check
type GreaterThanCheck struct {
	Field string `yaml:"field" validate:"required"`
	Value int    `yaml:"value" validate:"required"`
}

// LessThanCheck represents a less than check
type LessThanCheck struct {
	Field string `yaml:"field" validate:"required"`
	Value int    `yaml:"value" validate:"required"`
}

// MatchesCheck represents a regex match check
type MatchesCheck struct {
	Regex string `yaml:"regex" validate:"required"`
	Field string `yaml:"field" validate:"required"`
}

// PresentCheck represents a field presence check
type PresentCheck struct {
	Field string `yaml:"field" validate:"required"`
}

// NotPresentCheck represents a field absence check
type NotPresentCheck struct {
	Field string `yaml:"field" validate:"required"`
}

// Validate performs validation on the configuration
func (c *Config) Validate() error {
	for _, cv := range c.ConfigValues {
		if cv.Name == "" {
			return fmt.Errorf("config value name cannot be empty")
		}
		if cv.Command != "" && cv.Type == "" {
			return fmt.Errorf("config value with command must specify type")
		}
	}
	return nil
}

// GetValue retrieves a configuration value
func (c *Config) GetValue(name string) (any, error) {
	for _, cv := range c.ConfigValues {
		if cv.Name == name {
			return cv.Value, nil
		}
	}
	return nil, fmt.Errorf("config value %s not found", name)
}
