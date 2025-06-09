package workflow

import (
	"fmt"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
)

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
		case CheckTypeEquals:
			workflow.ExpectResponse.Body.Equals = append(workflow.ExpectResponse.Body.Equals, config.EqualsCheck{
				Field: check.Field,
				Value: check.Value,
				Type:  check.Type,
			})
		case CheckTypeMatches:
			workflow.ExpectResponse.Body.Matches = append(workflow.ExpectResponse.Body.Matches, config.MatchesCheck{
				Field: check.Field,
				Regex: check.Regex,
			})
		case CheckTypePresent:
			workflow.ExpectResponse.Body.Presents = append(workflow.ExpectResponse.Body.Presents, config.PresentCheck{
				Field: check.Field,
			})
		case CheckTypeNotPresent:
			workflow.ExpectResponse.Body.NotPresents = append(workflow.ExpectResponse.Body.NotPresents, config.NotPresentCheck{
				Field: check.Field,
			})
		case CheckTypeGreaterThan:
			if val, ok := check.Value.(int); ok {
				workflow.ExpectResponse.Body.GreaterThans = append(workflow.ExpectResponse.Body.GreaterThans, config.GreaterThanCheck{
					Field: check.Field,
					Value: val,
				})
			}
		case CheckTypeLessThan:
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

// GenerateAllWorkflows generates all workflow test cases
func (g *WorkflowGenerator) GenerateAllWorkflows() error {
	workflows := []struct {
		name     string
		generate func() []TestCase
		opts     []GenerateOption
	}{
		{
			name:     "simplified_config_workflow",
			generate: g.GenerateGeneralWorkflow,
			opts: []GenerateOption{
				WithHost("https://api.example.com"),
				WithNetworkEnable(true),
				WithSkip(false),
			},
		},
	}

	for _, wf := range workflows {
		fmt.Printf("Generating %s...\n", wf.name)
		testCases := wf.generate()
		if err := g.GenerateTestSuite(wf.name, testCases, wf.opts...); err != nil {
			return fmt.Errorf("failed to generate %s: %v", wf.name, err)
		}
	}

	return nil
}

// GenerateGeneralWorkflow generates general workflow test cases
func (g *WorkflowGenerator) GenerateGeneralWorkflow() []TestCase {
	return []TestCase{
		{
			Name:        "simplified_config_workflow",
			Description: "Simplified 4-step workflow demonstrating config variable usage with Faker API",
			Steps: []TestStep{
				{
					Name:   "get_user_data",
					Method: "GET",
					URI:    "https://fakerapi.it/api/v2/persons?_quantity=1&_locale=en_US",
					Headers: map[string]string{
						"Content-Type":  "application/json",
						"X-Merchant-ID": "{{.mid}}",
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							CheckType: "present",
							Field:     "data",
						},
						{
							CheckType: "equals",
							Field:     "total",
							Value:     1,
						},
						{
							CheckType: "present",
							Field:     "data[0].email",
						},
					},
					FromResponses: []FromResponse{
						{
							Step:      "get_user_data",
							Name:      "user_email",
							FromField: "data[0].email",
						},
						{
							Step:      "get_user_data",
							Name:      "user_name",
							FromField: "data[0].firstname",
						},
					},
				},
				{
					Name:   "create_user_account",
					Method: "POST",
					URI:    "/api/merchants/{{.mid}}/users",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"email":     "#{{.user_email}}",
						"name":      "#{{.user_name}}",
						"legacy_id": "#{{.legacy_id}}",
						"metadata": map[string]interface{}{
							"created_by":  "config_test",
							"merchant_id": "{{.mid}}",
						},
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							CheckType: "present",
							Field:     "data.id",
						},
						{
							CheckType: "equals",
							Field:     "data.email",
							Value:     "{{.user_email}}",
						},
					},
					FromResponses: []FromResponse{
						{
							Step:      "create_user_account",
							Name:      "user_id",
							FromField: "data.id",
						},
					},
				},
				{
					Name:   "create_order",
					Method: "POST",
					URI:    "/api/merchants/{{.mid}}/orders",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"user_id": "#{{.user_id}}",
						"items": []map[string]interface{}{
							{
								"legacy_id": "#{{.config_legacy_id}}",
								"quantity":  1,
								"price": map[string]interface{}{
									"cents":        1000,
									"currency_iso": "TWD",
								},
							},
						},
						"start_at": "#{{.next_day}}",
						"end_at":   "#{{.next_month}}",
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							CheckType: "present",
							Field:     "data.id",
						},
						{
							CheckType: "equals",
							Field:     "data.user_id",
							Value:     "{{.user_id}}",
						},
						{
							CheckType: "present",
							Field:     "data.items",
						},
					},
					FromResponses: []FromResponse{
						{
							Step:      "create_order",
							Name:      "order_id",
							FromField: "data.id",
						},
					},
				},
				{
					Name:   "get_order_status",
					Method: "GET",
					URI:    "/api/merchants/{{.mid}}/orders/{{.order_id}}",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							CheckType: "present",
							Field:     "data.id",
						},
						{
							CheckType: "equals",
							Field:     "data.id",
							Value:     "{{.order_id}}",
						},
						{
							CheckType: "present",
							Field:     "data.status",
						},
						{
							CheckType: "greater_than",
							Field:     "data.items.length",
							Value:     0,
						},
					},
				},
			},
		},
	}
}
