package workflow

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/config/types"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
)

// TestCase represents a test case
type TestCase struct {
	Name        string
	Description string
	Steps       []TestStep
}

// TestStep represents a test step
type TestStep struct {
	Name           string
	Method         string
	URI            string
	Headers        map[string]string
	Body           interface{}
	Query          map[string]string
	Variables      []string
	FromResponses  []FromResponse
	ExpectedCode   string
	ExpectedStatus int
	ResponseChecks []ResponseCheck
}

// FromResponse represents a response field reference
type FromResponse struct {
	Step      string
	Name      string
	FromField string
}

// ResponseCheck represents a response check
type ResponseCheck struct {
	CheckType string
	Field     string
	Type      types.FieldType
	Value     interface{}
	Regex     string
}

// Generator is a workflow generator used to generate workflow-related test cases
type WorkflowGenerator struct {
	*generator.Generator
}

// NewGenerator creates a new workflow generator
func NewWorkflowGenerator(outputDir string) *WorkflowGenerator {
	return &WorkflowGenerator{
		Generator: generator.NewGenerator(outputDir),
	}
}

// GenerateAllWorkflows generates all workflow test cases
func (g *WorkflowGenerator) GenerateAllWorkflows() error {
	workflows := []struct {
		name     string
		generate func() []TestCase
		opts     []GenerateOption
	}{
		{
			name:     "example_workflow",
			generate: g.GenerateGeneralWorkflow,
			opts: []GenerateOption{
				WithNetworkEnable(false),
			},
		},
	}

	for _, wf := range workflows {
		fmt.Printf("Generating %s...\n", wf.name)
		testCases := wf.generate()
		if err := g.GenerateTestSuite(wf.name, testCases, wf.opts...); err != nil {
			return fmt.Errorf("failed to generate %s: %v", wf.name, err)
		}
		fmt.Printf("Successfully generated %s\n", wf.name)
	}

	return nil
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

// GenerateTestSuite generates a test suite
func (g *WorkflowGenerator) GenerateTestSuite(name string, testCases []TestCase, opts ...GenerateOption) error {
	// Create new APITest configuration
	apiTest := config.APITest{
		Name:          name,
		Host:          "localhost",
		NetworkEnable: false,
		Skip:          false,
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

	return g.Generator.WriteYAML(name, suite)
}

// generateScenario generates a test scenario
func (g *WorkflowGenerator) generateScenario(tc TestCase) config.Scenario {
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
func (g *WorkflowGenerator) generateWorkflow(step TestStep) config.Workflow {
	workflow := config.Workflow{
		Step: step.Name,
		Request: config.Request{
			Method:  step.Method,
			URI:     step.URI,
			Headers: convertHeadersToString(step.Headers),
			Body:    convertBodyToString(step.Body),
			Vars:    convertVariables(step.Variables),
			Query:   convertQueryToString(step.Query),
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
				Type:  check.Type,
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

	// 序列化為 JSON
	data, err := json.Marshal(body)
	if err != nil {
		return ""
	}

	// 處理序列化後的字符串，去掉模板變量的引號
	result := string(data)
	re := regexp.MustCompile(`"@(\{\{\.\w+\}\})"`)
	result = re.ReplaceAllString(result, "$1")

	return result
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

// convertQueryToString converts query map to string
func convertQueryToString(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}
	data, _ := json.Marshal(query)
	return string(data)
}

// GenerateGeneralWorkflow generates general workflow test cases
func (g *WorkflowGenerator) GenerateGeneralWorkflow() []TestCase {
	return []TestCase{
		{
			Name:        "general_workflow",
			Description: "General API workflow test",
			Steps: []TestStep{
				{
					Name:   "health_check",
					Method: "GET",
					URI:    "/health_check",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							CheckType: "present",
							Field:     "status",
						},
					},
				},
			},
		},
	}
}
