package workflow

import (
	"fmt"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/RyanTokManMokMTM/api-testing-go/utils/tool/generator"
)

// Generator 工作流生成器，用於生成工作流相關的測試用例
type WorkflowGenerator struct {
	*generator.Generator
}

// NewGenerator 創建一個新的工作流生成器
func NewWorkflowGenerator(outputDir string) *WorkflowGenerator {
	return &WorkflowGenerator{
		Generator: generator.NewGenerator(outputDir, nil), // 使用默認配置
	}
}

// GenerateAllWorkflows 生成所有工作流的測試用例
func (g *WorkflowGenerator) GenerateAllWorkflows() error {
	workflows := []struct {
		name     string
		generate func() []generator.TestCase
		opts     []generator.GenerateOption
	}{
		{
			name:     "subscription_workflow",
			generate: g.GenerateSubscriptionWorkflow,
			opts: []generator.GenerateOption{
				generator.WithHost("{{.API_HOST}}"),
				generator.WithNetworkEnable(true),
				generator.WithGlobalHook(config.Hook{
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
			opts: []generator.GenerateOption{
				generator.WithHost("{{.API_HOST}}"),
				generator.WithNetworkEnable(true),
				generator.WithGlobalHook(config.Hook{
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
			opts: []generator.GenerateOption{
				generator.WithHost("localhost"),
				generator.WithNetworkEnable(false),
			},
		},
		{
			name:     "general_workflow",
			generate: g.GenerateGeneralWorkflow,
			opts: []generator.GenerateOption{
				generator.WithHost("localhost"),
				generator.WithNetworkEnable(false),
			},
		},
	}

	for _, wf := range workflows {
		fmt.Printf("Generating %s...\n", wf.name)
		testCases := wf.generate()
		if err := g.Generator.GenerateTestSuite(wf.name, testCases, wf.opts...); err != nil {
			return fmt.Errorf("failed to generate %s: %v", wf.name, err)
		}
		fmt.Printf("Successfully generated %s\n", wf.name)
	}

	return nil
}

// GenerateSubscriptionWorkflow 生成訂閱工作流測試用例
func (g *WorkflowGenerator) GenerateSubscriptionWorkflow() []generator.TestCase {
	return []generator.TestCase{
		{
			Name:        "subscription_workflow",
			Description: "訂閱 API 工作流測試",
			Steps: []generator.TestStep{
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
					ResponseChecks: []generator.ResponseCheck{
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
					FromResponses: []generator.FromResponse{
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
					ResponseChecks: []generator.ResponseCheck{
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

// GenerateOrderWorkflow 生成訂單工作流測試用例
func (g *WorkflowGenerator) GenerateOrderWorkflow() []generator.TestCase {
	return []generator.TestCase{
		{
			Name:        "order_workflow",
			Description: "訂單 API 工作流測試",
			Steps: []generator.TestStep{
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
					ResponseChecks: []generator.ResponseCheck{
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

// GenerateCouponWorkflow 生成優惠券工作流測試用例
func (g *WorkflowGenerator) GenerateCouponWorkflow() []generator.TestCase {
	return []generator.TestCase{
		{
			Name:        "coupon_workflow",
			Description: "優惠券 API 工作流測試",
			Steps: []generator.TestStep{
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
					ResponseChecks: []generator.ResponseCheck{
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

// GenerateGeneralWorkflow 生成通用工作流測試用例
func (g *WorkflowGenerator) GenerateGeneralWorkflow() []generator.TestCase {
	return []generator.TestCase{
		{
			Name:        "general_workflow",
			Description: "通用 API 工作流測試",
			Steps: []generator.TestStep{
				{
					Name:   "health_check",
					Method: "GET",
					URI:    "/health_check",
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					ExpectedCode:   "SUCCESS",
					ExpectedStatus: 200,
					ResponseChecks: []generator.ResponseCheck{
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
