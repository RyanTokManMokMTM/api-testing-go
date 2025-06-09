package api_test

//nolint:revive // must use this style
import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	apitestrunner "github.com/RyanTokManMokMTM/api-testing-go/api_test_runner"
	"github.com/RyanTokManMokMTM/api-testing-go/config"
	apihelper "github.com/RyanTokManMokMTM/api-testing-go/utils/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Constants for configuration paths
const (
	configRootPath = "../config/etc/api-test"
	configDir      = "config"
	workflowsDir   = "workflows"
	configFile     = "config.yaml"
	yamlExtension  = ".yaml"
)

// TestSuite represents the main API test suite
var _ = Describe("[ API TEST WORKFLOW ]", func() {
	var (
		t             GinkgoTInterface
		responseCache map[string]interface{}
		runner        *apitestrunner.APITestRunner
	)

	BeforeEach(func() {
		t = GinkgoT()
		responseCache = make(map[string]interface{})
	})

	AfterEach(func() {
		clearResponseCache(responseCache)
	})

	// Load and run all test configurations
	for _, apiTestConfig := range loadAllTestConfigs() {
		apiTestConfig := apiTestConfig // Capture loop variable
		Describe("", func() {
			BeforeEach(func() {
				if !apiTestConfig.APITest.Skip {
					runner = apitestrunner.NewAPITestRuuner(&apiTestConfig.APITest)
					runBeforeHook(t, runner, &apiTestConfig.APITest, apiTestConfig.APITest.GlobalHook, &responseCache)
				}
			})

			AfterEach(func() {
				if !apiTestConfig.APITest.Skip {
					// Run after whole DescribeTable done
					runAfterHook(t, runner, &apiTestConfig.APITest, apiTestConfig.APITest.GlobalHook, &responseCache)
					// Clear the cache after each scenario / test
					clearResponseCache(responseCache)
				}
			})

			DescribeTableSubtree(fmt.Sprintf("< API Test : %s >", apiTestConfig.APITest.Name),
				func(cfg *config.APITest, scenarioIndex int) {
					BeforeEach(func() {
						if !apiTestConfig.APITest.Skip && !cfg.Scenarios[scenarioIndex].Skip {
							// Run before each scenario
							runBeforeHook(t, runner, cfg, cfg.Scenarios[scenarioIndex].Hook, &responseCache)
						}
					})

					AfterEach(func() {
						// Run after each scenario
						if !apiTestConfig.APITest.Skip && !cfg.Scenarios[scenarioIndex].Skip {
							runAfterHook(t, runner, cfg, cfg.Scenarios[scenarioIndex].Hook, &responseCache)
						}
					})

					// Execute all scenarios for this test configuration
					for _, scenario := range cfg.Scenarios {
						scenario := scenario // Capture loop variable
						if !apiTestConfig.APITest.Skip && !scenario.Skip {
							It(fmt.Sprintf("< Scenario : %s >", scenario.Name), func() {
								runTest(t, runner, cfg, scenario.Workflows, &responseCache)
							})
						}
					}
				},
				createTestEntries(apiTestConfig),
			)
		})
	}
})

// loadAllTestConfigs loads all test configurations from YAML files
func loadAllTestConfigs() []*config.APITesting {
	configMap, err := loadConfigMap()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config map: %v", err))
	}

	apiTestConfigs, err := loadAPITestConfigs(configMap)
	if err != nil {
		panic(fmt.Sprintf("Failed to load API test configs: %v", err))
	}

	return apiTestConfigs
}

// loadConfigMap loads and processes the main configuration file
func loadConfigMap() (map[string]interface{}, error) {
	configPath := path.Join(configRootPath, configDir, configFile)
	configs, err := loadConfig[config.Config](configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	configMap := make(map[string]interface{})
	for _, cfg := range configs.ConfigValues {
		value, err := processConfigValue(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to process config value %s: %w", cfg.Name, err)
		}
		configMap[cfg.Name] = value
	}

	return configMap, nil
}

// processConfigValue processes a single config value, executing commands if needed
func processConfigValue(cfg config.ConfigValue) (interface{}, error) {
	if cfg.Command != "" {
		// Execute command and get output
		output, err := exec.Command("bash", "-c", cfg.Command).Output()
		if err != nil {
			return nil, fmt.Errorf("command execution failed: %w", err)
		}

		// Remove trailing newline and convert to appropriate type
		outputStr := strings.TrimSpace(string(output))
		converted, err := apihelper.ConvertStrToType(outputStr, cfg.Type)
		if err != nil {
			return nil, fmt.Errorf("type conversion failed: %w", err)
		}

		return converted, nil
	}

	return cfg.Value, nil
}

// loadAPITestConfigs loads all API test configuration files from the workflows directory
func loadAPITestConfigs(configMap map[string]interface{}) ([]*config.APITesting, error) {
	workflowsPath := path.Join(configRootPath, workflowsDir)
	files, err := os.ReadDir(workflowsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read workflows directory %s: %w", workflowsPath, err)
	}

	var apiTestConfigs []*config.APITesting
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), yamlExtension) {
			continue
		}

		filePath := path.Join(workflowsPath, file.Name())
		apiConfig, err := loadConfig[config.APITesting](filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load API test config from %s: %w", filePath, err)
		}

		apiConfig.APITest.Config = configMap
		apiTestConfigs = append(apiTestConfigs, apiConfig)
	}

	return apiTestConfigs, nil
}

// loadConfig loads a YAML configuration file into the specified type
func loadConfig[T any](filePath string) (*T, error) {
	var config T
	if err := apihelper.LoadYamlData(filePath, &config); err != nil {
		return nil, fmt.Errorf("failed to load YAML data from %s: %w", filePath, err)
	}
	return &config, nil
}

// runTest executes a series of workflows and stores responses in cache
func runTest(
	t GinkgoTInterface,
	runner *apitestrunner.APITestRunner,
	cfg *config.APITest,
	workflows []config.Workflow,
	cache *map[string]interface{},
) {
	var responseData interface{}
	for _, workflow := range workflows {
		apiResponse, err := runner.Start(t, cfg, workflow, responseData, cache)
		Expect(err).To(BeNil(), "API request failed for workflow %s", workflow.Step)

		// Parse response and store in cache
		apiResponse.JSON(&responseData)
		(*cache)[workflow.Step] = responseData
	}
}

// createTestEntries creates table entries for Ginkgo's DescribeTableSubtree
func createTestEntries(apiTestConfig *config.APITesting) []TableEntry {
	var tableEntries []TableEntry
	for i := range apiTestConfig.APITest.Scenarios {
		tableEntries = append(tableEntries, Entry("", &apiTestConfig.APITest, i))
	}
	return tableEntries
}

// runBeforeHook executes before hooks and initializes variables
func runBeforeHook(
	t GinkgoTInterface,
	runner *apitestrunner.APITestRunner,
	cfg *config.APITest,
	hook config.Hook,
	responseCache *map[string]interface{},
) {
	// Execute before workflows
	runTest(t, runner, cfg, hook.Before.Workflows, responseCache)

	// Initialize variables from response data
	for _, regVar := range hook.Before.InitVars {
		response, exists := (*responseCache)[regVar.FromStep]
		Expect(exists).To(BeTrue(), "Response data not found for step %s", regVar.FromStep)
		Expect(response).ToNot(BeNil(), "Response data is nil for step %s", regVar.FromStep)

		data, err := apihelper.GetRespFieldData(regVar.Field, response)
		Expect(err).To(BeNil(), "Failed to extract field %s from response", regVar.Field)
		Expect(data).ToNot(BeNil(), "Extracted data is nil for field %s", regVar.Field)

		cfg.Config[regVar.Name] = data
	}
}

// runAfterHook executes after hooks
func runAfterHook(
	t GinkgoTInterface,
	runner *apitestrunner.APITestRunner,
	cfg *config.APITest,
	hook config.Hook,
	responseCache *map[string]interface{},
) {
	runTest(t, runner, cfg, hook.After.Workflows, responseCache)
}

// clearResponseCache clears all entries from the response cache
func clearResponseCache(cache map[string]interface{}) {
	for key := range cache {
		delete(cache, key)
	}
}
