package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
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
	Type  string
	Field string
	Value interface{}
	Regex string
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
			name:     "subscription_workflow",
			generate: g.GenerateSubscriptionWorkflow,
			opts: []GenerateOption{
				WithHost("{{.API_HOST}}"),
				WithNetworkEnable(true),
				WithGlobalHook(config.Hook{
					After: config.HookActions{
						Workflows: []config.Workflow{
							{
								Step: "cleanup_subscription",
								Request: config.Request{
									Method: "DELETE",
									URI:    "/api/merchants/{{.mid}}/subscriptions/{{.subscription_id}}",
								},
								ExpectResponse: config.Response{
									Code:       "SUCCESS",
									StatusCode: 200,
								},
							},
						},
					},
				}),
			},
		},
		{
			name:     "order_workflow",
			generate: g.GenerateOrderWorkflow,
			opts: []GenerateOption{
				WithHost("{{.API_HOST}}"),
				WithNetworkEnable(true),
				WithGlobalHook(config.Hook{
					After: config.HookActions{
						Workflows: []config.Workflow{
							{
								Step: "cleanup_order",
								Request: config.Request{
									Method: "DELETE",
									URI:    "/api/merchants/{{.mid}}/orders/{{.order_id}}",
								},
								ExpectResponse: config.Response{
									Code:       "SUCCESS",
									StatusCode: 200,
								},
							},
						},
					},
				}),
			},
		},
		{
			name:     "coupon_workflow",
			generate: g.GenerateCouponWorkflow,
			opts: []GenerateOption{
				WithHost("localhost"),
				WithNetworkEnable(false),
			},
		},
		{
			name:     "general_workflow",
			generate: g.GenerateGeneralWorkflow,
			opts: []GenerateOption{
				WithHost("localhost"),
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

//

// GenerateSubscriptionWorkflow generates subscription workflow test cases
func (g *WorkflowGenerator) GenerateSubscriptionWorkflow() []TestCase {
	return []TestCase{
		{
			Name:        "subscription_workflow",
			Description: "Subscription API workflow test",
			Steps: []TestStep{
				{
					Name:   "create_subscribable_entity_plan",
					Method: "POST",
					URI:    "/api/subscribables_entities",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Variables: []string{
						"legacy_id",
						"config_legacy_id",
					},
					Body: map[string]interface{}{
						"legacy_id": "{{.legacy_id}}",
						"type":      "plan",
						"key":       "test_plan",
						"description_translations": map[string]string{
							"en": "test plan",
						},
						"configurations": []map[string]interface{}{
							{
								"legacy_id": "{{.config_legacy_id}}",
								"name_translations": map[string]string{
									"en": "test",
								},
								"key":           "test_month",
								"is_trial":      false,
								"recurring_day": 0,
								"duration":      1,
								"duration_unit": "month",
								"charges": []map[string]interface{}{
									{
										"price": map[string]interface{}{
											"cents":        10,
											"currency_iso": "HKD",
										},
										"one_off": true,
										"title":   "subscription_fee",
									},
								},
							},
						},
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							Type:  "present",
							Field: "data.id",
						},
					},
				},
				{
					Name:   "create_subscription",
					Method: "POST",
					URI:    "/api/merchants/{{.mid}}/subscriptions",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Variables: []string{
						"mid",
					},
					FromResponses: []FromResponse{
						{
							Step:      "create_subscribable_entity_plan",
							Name:      "subscribable_entity_id",
							FromField: "data.id",
						},
					},
					Body: map[string]interface{}{
						"subscribable_entity_id": "{{.subscribable_entity_id}}",
						"start_at":               "{{.next_day}}",
						"end_at":                 "{{.next_month}}",
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							Type:  "present",
							Field: "data.id",
						},
					},
				},
			},
		},
	}
}

// GenerateOrderWorkflow generates order workflow test cases
func (g *WorkflowGenerator) GenerateOrderWorkflow() []TestCase {
	return []TestCase{
		{
			Name:        "order_workflow",
			Description: "Order API workflow test",
			Steps: []TestStep{
				{
					Name:   "create_order",
					Method: "POST",
					URI:    "/api/merchants/{{.mid}}/orders",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Variables: []string{
						"mid",
					},
					Body: map[string]interface{}{
						"items": []map[string]interface{}{
							{
								"name":     "Test Item",
								"quantity": 1,
								"price": map[string]interface{}{
									"cents":        100,
									"currency_iso": "HKD",
								},
							},
						},
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							Type:  "present",
							Field: "data.id",
						},
					},
				},
			},
		},
	}
}

// GenerateCouponWorkflow generates coupon workflow test cases
func (g *WorkflowGenerator) GenerateCouponWorkflow() []TestCase {
	return []TestCase{
		{
			Name:        "coupon_workflow",
			Description: "Coupon API workflow test",
			Steps: []TestStep{
				{
					Name:   "create_coupon",
					Method: "POST",
					URI:    "/api/merchants/{{.mid}}/coupons",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Variables: []string{
						"mid",
					},
					Body: map[string]interface{}{
						"code":        "TEST_COUPON",
						"type":        "percentage",
						"value":       10,
						"start_at":    "{{.next_day}}",
						"end_at":      "{{.next_month}}",
						"usage_limit": 100,
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []ResponseCheck{
						{
							Type:  "present",
							Field: "data.id",
						},
					},
				},
			},
		},
	}
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
							Type:  "present",
							Field: "status",
						},
					},
				},
			},
		},
	}
}
