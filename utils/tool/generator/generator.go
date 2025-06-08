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
	// DefaultOutputDir 是默認的輸出目錄
	DefaultOutputDir = "config/etc/api-test/workflows"
)

// APITestOptions 定義 API 測試的配置選項
type APITestOptions struct {
	// 基本配置
	Host          string
	NetworkEnable bool
	Skip          bool
	// 全局鉤子配置
	GlobalHook *config.Hook
}

// DefaultAPITestOptions 返回默認的 API 測試配置
func DefaultAPITestOptions() *APITestOptions {
	return &APITestOptions{
		Host:          "localhost",
		NetworkEnable: false,
		Skip:          false,
		GlobalHook:    nil,
	}
}

// Generator 基礎生成器，提供通用的生成功能
type Generator struct {
	outputDir string
	options   *APITestOptions
}

// NewGenerator 創建一個新的基礎生成器
func NewGenerator(outputDir string, options *APITestOptions) *Generator {
	if options == nil {
		options = DefaultAPITestOptions()
	}
	if outputDir == "" {
		outputDir = DefaultOutputDir
	}
	return &Generator{
		outputDir: outputDir,
		options:   options,
	}
}

// GenerateTestSuite 生成測試套件
func (g *Generator) GenerateTestSuite(name string, testCases []TestCase, opts ...GenerateOption) error {
	// 創建新的 APITest 配置，使用 Generator 的 options 作為默認值
	apiTest := config.APITest{
		Name:          name,
		Host:          g.options.Host,
		NetworkEnable: g.options.NetworkEnable,
		Skip:          g.options.Skip,
		// GlobalHook 默認為空
	}

	// 應用所有選項
	for _, opt := range opts {
		opt(&apiTest)
	}

	// 構建完整的 APITesting 結構
	suite := &config.APITesting{
		APITest: apiTest,
	}

	// 生成測試場景
	for _, tc := range testCases {
		scenario := g.generateScenario(tc)
		suite.APITest.Scenarios = append(suite.APITest.Scenarios, scenario)
	}

	return g.writeYAML(name, suite)
}

// GenerateOption 定義生成選項的函數類型
type GenerateOption func(*config.APITest)

// WithHost 設置 Host 選項
func WithHost(host string) GenerateOption {
	return func(apiTest *config.APITest) {
		apiTest.Host = host
	}
}

// WithNetworkEnable 設置 NetworkEnable 選項
func WithNetworkEnable(enable bool) GenerateOption {
	return func(apiTest *config.APITest) {
		apiTest.NetworkEnable = enable
	}
}

// WithSkip 設置 Skip 選項
func WithSkip(skip bool) GenerateOption {
	return func(apiTest *config.APITest) {
		apiTest.Skip = skip
	}
}

// WithGlobalHook 設置全局鉤子
func WithGlobalHook(hook config.Hook) GenerateOption {
	return func(test *config.APITest) {
		test.GlobalHook = hook
	}
}

// generateScenario 生成測試場景
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

// generateWorkflow 生成工作流程
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

// convertHeadersToString 將 headers map 轉換為 JSON 字符串
func convertHeadersToString(headers map[string]string) string {
	if len(headers) == 0 {
		return ""
	}
	data, _ := json.Marshal(headers)
	return string(data)
}

// convertBodyToString 將 body 轉換為 JSON 字符串
func convertBodyToString(body interface{}) string {
	if body == nil {
		return ""
	}
	data, _ := json.Marshal(body)
	return string(data)
}

// convertVariables 將變量列表轉換為 Var 結構體列表
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

// writeYAML 將測試套件寫入 YAML 文件
func (g *Generator) writeYAML(name string, suite *config.APITesting) error {
	// 確保輸出目錄存在
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// 生成文件名，使用小寫並將空格替換為下劃線
	filename := filepath.Join(g.outputDir, fmt.Sprintf("%s.yaml", strings.ToLower(strings.ReplaceAll(name, " ", "_"))))

	// 將結構體轉換為 YAML
	data, err := yaml.Marshal(suite)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// 寫入文件
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}

	return nil
}

// GetOutputDir 獲取輸出目錄
func (g *Generator) GetOutputDir() string {
	return g.outputDir
}
