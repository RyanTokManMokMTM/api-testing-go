package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/RyanTokManMokMTM/api-testing-go/config"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultOutputDir is the default output directory
	DefaultOutputDir = "config/etc/api-test/workflows"
)

// APITestOptions defines API test configuration options
type APITestOptions struct {
	// Basic configuration
	Host          string
	NetworkEnable bool
	Skip          bool
	// Global hook configuration
	GlobalHook *config.Hook
}

// DefaultAPITestOptions returns the default API test configuration
func DefaultAPITestOptions() *APITestOptions {
	return &APITestOptions{
		Host:          "localhost",
		NetworkEnable: false,
		Skip:          false,
		GlobalHook:    nil,
	}
}

// Generator is a base generator that provides common generation functionality
type Generator struct {
	outputDir string
	Options   *APITestOptions
}

// NewGenerator creates a new base generator
func NewGenerator(outputDir string, options *APITestOptions) *Generator {
	if options == nil {
		options = DefaultAPITestOptions()
	}
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	return &Generator{
		outputDir: outputDir,
		Options:   options,
	}
}

// GenerateTestSuite generates a test suite
func (g *Generator) GenerateTestSuite(name string, testCases []TestCase, opts ...GenerateOption) error {
	// Create new APITest configuration using Generator's options as default values
	apiTest := config.APITest{
		Name:          name,
		Host:          g.Options.Host,
		NetworkEnable: g.Options.NetworkEnable,
		Skip:          g.Options.Skip,
		// GlobalHook defaults to empty
	}

	// Apply all options
	for _, opt := range opts {
		opt(&apiTest)
	}

	// Build complete APITesting structure
	suite := &config.APITesting{
		APITest: apiTest,
	}

	// Generate test scenarios
	for _, tc := range testCases {
		scenario := g.generateScenario(tc)
		suite.APITest.Scenarios = append(suite.APITest.Scenarios, scenario)
	}

	return g.writeYAML(name, suite)
}

// GenerateOption defines a function type for generation options
type GenerateOption func(*config.APITest)

// WithHost sets the Host option
func WithHost(host string) GenerateOption {
	return func(apiTest *config.APITest) {
		apiTest.Host = host
	}
}

// WithNetworkEnable sets the NetworkEnable option
func WithNetworkEnable(enable bool) GenerateOption {
	return func(apiTest *config.APITest) {
		apiTest.NetworkEnable = enable
	}
}

// WithSkip sets the Skip option
func WithSkip(skip bool) GenerateOption {
	return func(apiTest *config.APITest) {
		apiTest.Skip = skip
	}
}

// WithGlobalHook sets the global hook
func WithGlobalHook(hook config.Hook) GenerateOption {
	return func(test *config.APITest) {
		test.GlobalHook = hook
	}
}

// generateScenario generates a test scenario
func (g *Generator) generateScenario(tc TestCase) config.Scenario {
	scenario := config.Scenario{
		Name:      tc.Name,
		Skip:      false,
		Workflows: make([]config.Workflow, len(tc.Steps)),
	}

	for i, step := range tc.Steps {
		scenario.Workflows[i] = g.generateWorkflow(step)
	}

	return scenario
}

// generateWorkflow generates a workflow
func (g *Generator) generateWorkflow(step TestStep) config.Workflow {
	workflow := config.Workflow{
		Step: step.Name,
		Request: config.Request{
			Method:  step.Method,
			URI:     step.URI,
			Headers: convertHeadersToString(step.Headers),
			Body:    convertBodyToString(step.Body),
			Vars:    convertVariables(step.Variables),
		},
		ExpectResponse: config.Response{
			Code:       step.ExpectedCode,
			StatusCode: step.ExpectedStatus,
			Body:       config.BodyCheck{},
		},
	}

	// Add FromResponse if any
	if len(step.FromResponses) > 0 {
		workflow.Request.FromResponse = make([]config.FromResponse, len(step.FromResponses))
		for i, fr := range step.FromResponses {
			workflow.Request.FromResponse[i] = config.FromResponse{
				Step:      fr.Step,
				Name:      fr.Name,
				FromField: fr.FromField,
			}
		}
	}

	// Add response checks
	for _, check := range step.ResponseChecks {
		switch check.Type {
		case "equals":
			workflow.ExpectResponse.Body.Equals = append(workflow.ExpectResponse.Body.Equals, config.EqualsCheck{
				Field: check.Field,
				Value: check.Value,
			})
		case "matches":
			workflow.ExpectResponse.Body.Matches = append(workflow.ExpectResponse.Body.Matches, config.MatchesCheck{
				Field: check.Field,
				Regex: check.Regex,
			})
		case "present":
			workflow.ExpectResponse.Body.Presents = append(workflow.ExpectResponse.Body.Presents, config.PresentCheck{
				Field: check.Field,
			})
		case "not_present":
			workflow.ExpectResponse.Body.NotPresents = append(workflow.ExpectResponse.Body.NotPresents, config.NotPresentCheck{
				Field: check.Field,
			})
		case "greater_than":
			if val, ok := check.Value.(int); ok {
				workflow.ExpectResponse.Body.GreaterThans = append(workflow.ExpectResponse.Body.GreaterThans, config.GreaterThanCheck{
					Field: check.Field,
					Value: val,
				})
			}
		case "less_than":
			if val, ok := check.Value.(int); ok {
				workflow.ExpectResponse.Body.LessThans = append(workflow.ExpectResponse.Body.LessThans, config.LessThanCheck{
					Field: check.Field,
					Value: val,
				})
			}
		}
	}

	return workflow
}

// convertHeadersToString converts headers map to JSON string
func convertHeadersToString(headers map[string]string) string {
	if len(headers) == 0 {
		return ""
	}
	data, _ := json.Marshal(headers)
	return string(data)
}

// convertBodyToString converts body to JSON string
func convertBodyToString(body interface{}) string {
	if body == nil {
		return ""
	}
	data, _ := json.Marshal(body)
	return string(data)
}

// convertVariables converts variable list to Var struct list
func convertVariables(vars []string) []config.Var {
	if len(vars) == 0 {
		return nil
	}
	result := make([]config.Var, len(vars))
	for i, v := range vars {
		result[i] = config.Var{Name: v}
	}
	return result
}

// writeYAML writes test suite to YAML file
func (g *Generator) writeYAML(name string, suite *config.APITesting) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename, convert to lowercase and replace spaces with underscores
	filename := filepath.Join(g.outputDir, fmt.Sprintf("%s.yaml", strings.ToLower(strings.ReplaceAll(name, " ", "_"))))

	// Convert struct to YAML
	data, err := yaml.Marshal(suite)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// GetOutputDir gets the output directory
func (g *Generator) GetOutputDir() string {
	return g.outputDir
}

// WriteConfigYAML writes configuration to YAML file
func (g *Generator) WriteConfigYAML(name string, cfg *config.Config) error {
	// Ensure output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Build output file path
	outputPath := filepath.Join(g.outputDir, fmt.Sprintf("%s.yaml", name))

	// Convert configuration to YAML
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %v", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}
