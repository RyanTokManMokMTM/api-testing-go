// Package workflow provides workflow generators and utilities for creating API test workflows.
package workflow

import (
	"fmt"
	"log"
	"os"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
)

// GenerateOption defines a function type for generation options
type GenerateOption func(*config.APITest)

// ScenarioOption defines a function type for modifying a Scenario during generation.
type ScenarioOption func(*config.Scenario)

// TestCaseWithOptions represents a test case with its options
type TestCaseWithOptions struct {
	TestCase TestCase
	Options  []ScenarioOption
}

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
func WithGlobalHook(hook Hook) GenerateOption {
	return func(test *config.APITest) {
		test.GlobalHook = convertHook(hook)
	}
}

// WithSkipScenario sets the Skip Scenario option
func WithSkipScenario(skip bool) ScenarioOption {
	return func(scenario *config.Scenario) {
		scenario.Skip = skip
	}
}

// WithScenarioHook sets the scenario with hook opts
func WithScenarioHook(hook Hook) ScenarioOption {
	return func(scenario *config.Scenario) {
		scenario.Hook = convertHook(hook)
	}
}

// Generator is a workflow generator used to generate workflow-related test cases.
type Generator struct {
	*generator.Generator
	logger *log.Logger
}

// generateWorkflow generates a workflow
func (g *Generator) generateWorkflow(step TestStep) config.Workflow {
	workflow := config.Workflow{
		Step: step.Name,
		Request: config.Request{
			Method:       step.Method,
			URI:          step.URI,
			Headers:      convertHeadersToString(step.Headers),
			Body:         convertBodyToString(step.Body),
			Query:        convertQueryToString(step.Query),
			Vars:         convertVariables(step.Variables),
			FromResponse: convertFromResponses(step.FromResponses),
		},
		ExpectResponse: config.Response{
			StatusCode: step.ExpectedStatus,
			Body:       config.BodyCheck{},
		},
	}

	// Add response checks
	for _, check := range step.ResponseChecks {

		switch check.CheckType {
		case CheckTypeEquals:
			workflow.ExpectResponse.Body.Equals = append(workflow.ExpectResponse.Body.Equals, config.EqualsCheck{
				Field: check.Field,
				Value: check.Value,
				Type:  check.Type.ToConfigTypeField(),
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

// NewGenerator creates a new workflow generator with the specified output directory.
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		Generator: generator.NewGenerator(outputDir),
		logger:    log.New(os.Stdout, "[WorkflowGenerator] ", log.LstdFlags),
	}
}

// generateScenario generates a test scenario from a test case and applies options
func (g *Generator) generateScenario(tc TestCase, opts ...ScenarioOption) config.Scenario {
	scenario := config.Scenario{
		Name:      tc.Name,
		Skip:      false,
		Workflows: make([]config.Workflow, len(tc.Steps)),
	}

	// Apply all options
	for _, opt := range opts {
		opt(&scenario)
	}

	for i, step := range tc.Steps {
		scenario.Workflows[i] = g.generateWorkflow(step)
	}

	return scenario
}

// GenerateTestSuite generates a test suite from a list of test cases and options.
func (g *Generator) GenerateTestSuite(name string, testCases []TestCaseWithOptions, opts ...GenerateOption) error {
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

	suite := &config.APITesting{
		APITest: apiTest,
	}

	// Generate test scenarios
	for _, tc := range testCases {
		scenario := g.generateScenario(tc.TestCase, tc.Options...)
		suite.APITest.Scenarios = append(suite.APITest.Scenarios, scenario)
	}

	return g.WriteYAML(name, suite)
}

// GenerateAllWorkflows generates all workflow test cases.
func (g *Generator) GenerateAllWorkflows() error {
	workflows := []struct {
		name     string
		generate func() []TestCaseWithOptions
		opts     []GenerateOption
	}{
		{
			name:     "spanish_data_workflow",
			generate: g.GenerateExampleWorkflow,
			opts: []GenerateOption{
				WithSkip(false),
				WithNetworkEnable(true),
				WithHost("api.generadordni.es"),
			},
		},
	}

	for _, wf := range workflows {
		g.logger.Printf("Generating %s...", wf.name)
		testCases := wf.generate()
		if err := g.GenerateTestSuite(wf.name, testCases, wf.opts...); err != nil {
			return fmt.Errorf("failed to generate %s: %w", wf.name, err)
		}
		g.logger.Printf("Successfully generated %s", wf.name)
	}

	return nil
}

// GenerateExampleWorkflow generates Spanish data generation API workflow test cases
func (g *Generator) GenerateExampleWorkflow() []TestCaseWithOptions {
	return []TestCaseWithOptions{
		{
			TestCase: TestCase{
				Name:        "Spanish Person Profile Generation",
				Description: "Generate Spanish person profiles using generadordni.es API",
				Steps: []TestStep{
					{
						Name:   "Generate Person Profile",
						Method: "GET",
						URI:    "/v2/profiles/person",
						Headers: map[string]string{
							"Accept": "application/json",
						},
						Query: map[string]string{
							"results": "1",
						},
						ExpectedStatus: 200,
					},
				},
			},
		},
	}
}
