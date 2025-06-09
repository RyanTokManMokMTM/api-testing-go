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

var _ = Describe("[ API TEST WORK FLOW ]", func() {
	var (
		t             GinkgoTInterface
		respDataCahce map[string]interface{}
		runner        *apitestrunner.APITestRunner
	)

	BeforeEach(func() {
		t = GinkgoT()
		respDataCahce = make(map[string]interface{})
	})

	AfterEach(func() {
		for k := range respDataCahce {
			delete(respDataCahce, k)
		}
	})

	for _, apiTestCfg := range loadAllTestAnd() {
		// def gloabal hook
		Describe("", func() {
			BeforeEach(func() {
				if !apiTestCfg.APITest.Skip {
					runner = apitestrunner.NewAPITestRuuner(&apiTestCfg.APITest)
					runBeforeHook(t, runner, &apiTestCfg.APITest, apiTestCfg.APITest.GlobalHook, &respDataCahce)
				}
			})

			AfterEach(func() {
				if !apiTestCfg.APITest.Skip {
					// run after whole DescribeTable done
					runAfterHook(t, runner, &apiTestCfg.APITest, apiTestCfg.APITest.GlobalHook, &respDataCahce)

					// clear the cache after each scenario / test
					for k := range respDataCahce {
						delete(respDataCahce, k)
					}
				}
			})

			DescribeTableSubtree(fmt.Sprintf("< API Test : %s >", apiTestCfg.APITest.Name),
				func(
					cfg *config.APITest,
					scenarioIndex int,
				) {
					BeforeEach(func() {
						if !apiTestCfg.APITest.Skip && !cfg.Scenarios[scenarioIndex].Skip {
							// run before each scenario
							runBeforeHook(
								t,
								runner,
								cfg,
								cfg.Scenarios[scenarioIndex].Hook,
								&respDataCahce)
						}
					})

					AfterEach(func() {
						// run after each scenario
						if !apiTestCfg.APITest.Skip && !cfg.Scenarios[scenarioIndex].Skip {
							runAfterHook(
								t,
								runner,
								cfg,
								cfg.Scenarios[scenarioIndex].Hook,
								&respDataCahce)
						}
					})

					for _, scenario := range cfg.Scenarios {
						if !apiTestCfg.APITest.Skip && !scenario.Skip {
							It(fmt.Sprintf("< Scenario : %s >", scenario.Name), func() {
								runTest(
									t,
									runner,
									cfg,
									scenario.Workflows,
									&respDataCahce)
							})
						}
					}
				},
				createTestEntry(apiTestCfg),
			)
		})
	}
})

func loadAllTestAnd() []*config.APITesting {
	root := "../config/etc/api-test"
	configPath := path.Join(root, "config", "config.yaml")
	cfgs, configLoadErr := loadConfig[config.Config](configPath)
	if configLoadErr != nil {
		panic(configLoadErr)
	}
	configMap := map[string]interface{}{}
	for _, cfg := range cfgs.ConfigValues {
		if cfg.Command != "" {
			v, execErr := exec.Command("bash", "-c", cfg.Command).Output()
			if execErr != nil {
				panic(execErr)
			}
			converted, converErr := apihelper.ConvertStrToType(string(v[:len(v)-1]), cfg.Type)
			if converErr != nil {
				panic(converErr)
			}
			configMap[cfg.Name] = converted
		} else {
			configMap[cfg.Name] = cfg.Value
		}
	}

	workflowsPath := path.Join(root, "workflow")
	files, err := os.ReadDir(workflowsPath)
	if err != nil {
		panic(err)
	}
	result := make([]*config.APITesting, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}

		path := path.Join(root, f.Name())
		apiCfg, loadYamlErr := loadConfig[config.APITesting](path)
		if loadYamlErr != nil {
			panic(loadYamlErr)
		}
		apiCfg.APITest.Config = configMap
		result = append(result, apiCfg)
	}
	return result
}

func loadConfig[T any](path string) (*T, error) {
	// Load testing case
	var config T
	loadConfigErr := apihelper.LoadYamlData(path, &config)
	if loadConfigErr != nil {
		return nil, loadConfigErr
	}
	return &config, nil
}

func runTest(
	t GinkgoTInterface,
	starter *apitestrunner.APITestRunner,
	cfg *config.APITest,
	workFlows []config.Workflow,
	cache *map[string]interface{},
) {
	var respData any
	for _, w := range workFlows {
		apiResp, err := starter.Start(t, cfg, w, respData, cache)
		Expect(err).To(BeNil())

		// Store the response for next step
		apiResp.JSON(&respData)
		// Store resp to cache
		(*cache)[w.Step] = respData
	}
}

func createTestEntry(
	apiTestCfg *config.APITesting,
) []TableEntry {
	var tableEntries []TableEntry
	for i := range apiTestCfg.APITest.Scenarios {
		tableEntries = append(
			tableEntries,
			Entry(
				"",
				&apiTestCfg.APITest,
				i,
			))
	}
	return tableEntries
}

func runBeforeHook(
	t GinkgoTInterface,
	starter *apitestrunner.APITestRunner,
	cfg *config.APITest,
	hook config.Hook,
	respDataCahce *map[string]interface{},
) {
	runTest(t, starter, cfg, hook.Before.Workflows, respDataCahce)
	for _, regVar := range hook.Before.InitVars {
		resp, ok := (*respDataCahce)[regVar.FromStep]
		Expect(ok).To(BeTrue())
		Expect(resp).ToNot(BeNil())

		data, err := apihelper.GetRespFieldData(regVar.Field, resp)
		Expect(err).To(BeNil())
		Expect(data).ToNot(BeNil())

		cfg.Config[regVar.Name] = data
	}
}

func runAfterHook(
	t GinkgoTInterface,
	starter *apitestrunner.APITestRunner,
	cfg *config.APITest,
	hook config.Hook,
	respDataCahce *map[string]interface{},
) {
	runTest(t, starter, cfg, hook.After.Workflows, respDataCahce)
}
