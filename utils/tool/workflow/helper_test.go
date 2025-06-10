package workflow

import (
	"testing"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
	"github.com/stretchr/testify/assert"
)

func TestConvertHeadersToString(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected string
	}{
		{
			name:     "empty headers",
			headers:  map[string]string{},
			expected: "",
		},
		{
			name:     "nil headers",
			headers:  nil,
			expected: "",
		},
		{
			name: "single header",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			expected: `{"Content-Type":"application/json"}`,
		},
		{
			name: "multiple headers",
			headers: map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token123",
				"User-Agent":    "TestClient/1.0",
			},
			expected: `{"Authorization":"Bearer token123","Content-Type":"application/json","User-Agent":"TestClient/1.0"}`,
		},
		{
			name: "headers with special characters",
			headers: map[string]string{
				"X-Custom-Header": "value with spaces",
				"Accept":          "text/html,application/json",
			},
			expected: `{"Accept":"text/html,application/json","X-Custom-Header":"value with spaces"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertHeadersToString(tt.headers)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertBodyToString(t *testing.T) {
	tests := []struct {
		name     string
		body     map[string]interface{}
		expected string
	}{
		{
			name:     "nil body",
			body:     nil,
			expected: "",
		},
		{
			name: "simple map body",
			body: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			expected: `{"name":"test","value":123}`,
		},
		{
			name: "body with template variables (no prefix)",
			body: map[string]interface{}{
				"user_id": "{{.user_id}}",
				"amount":  100,
			},
			expected: `{"amount":100,"user_id":"{{.user_id}}"}`,
		},
		{
			name: "body with template variables (with # prefix)",
			body: map[string]interface{}{
				"user_id": "#{{.user_id}}",
				"amount":  100,
				"name":    "{{.name}}",
			},
			expected: `{"amount":100,"name":"{{.name}}","user_id":{{.user_id}}}`,
		},
		{
			name: "body with mixed template variables",
			body: map[string]interface{}{
				"user_id":     "#{{.user_id}}",
				"merchant_id": "#{{.merchant_id}}",
				"name":        "{{.name}}",
				"status":      "#{{.status}}",
				"amount":      100,
			},
			expected: `{"amount":100,"merchant_id":{{.merchant_id}},"name":"{{.name}}","status":{{.status}},"user_id":{{.user_id}}}`,
		},
		{
			name: "nested map with template variables",
			body: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   "#{{.user_id}}",
					"name": "{{.user_name}}",
				},
				"order": map[string]interface{}{
					"id":     "#{{.order_id}}",
					"amount": 100,
				},
			},
			expected: `{"order":{"amount":100,"id":{{.order_id}}},"user":{"id":{{.user_id}},"name":"{{.user_name}}"}}`,
		},
		{
			name: "array with template variables",
			body: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"id":   "#{{.item_id}}",
						"name": "{{.item_name}}",
					},
				},
			},
			expected: `{"items":[{"id":{{.item_id}},"name":"{{.item_name}}"}]}`,
		},
		{
			name: "complex nested structure",
			body: map[string]interface{}{
				"request_id": "#{{.request_id}}",
				"data": map[string]interface{}{
					"user": map[string]interface{}{
						"id":    "#{{.user_id}}",
						"email": "{{.user_email}}",
						"metadata": map[string]interface{}{
							"created_at": "#{{.created_at}}",
							"updated_at": "{{.updated_at}}",
						},
					},
				},
			},
			expected: `{"data":{"user":{"email":"{{.user_email}}","id":{{.user_id}},"metadata":{"created_at":{{.created_at}},"updated_at":"{{.updated_at}}"}}},"request_id":{{.request_id}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertBodyToString(tt.body)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertVariables(t *testing.T) {
	tests := []struct {
		name     string
		vars     []string
		expected []config.Var
	}{
		{
			name:     "empty variables",
			vars:     []string{},
			expected: nil,
		},
		{
			name:     "nil variables",
			vars:     nil,
			expected: nil,
		},
		{
			name: "single variable",
			vars: []string{"user_id"},
			expected: []config.Var{
				{Name: "user_id"},
			},
		},
		{
			name: "multiple variables",
			vars: []string{"user_id", "merchant_id", "api_key"},
			expected: []config.Var{
				{Name: "user_id"},
				{Name: "merchant_id"},
				{Name: "api_key"},
			},
		},
		{
			name: "variables with special characters",
			vars: []string{"user_id", "api_key_v2", "merchant_id_123"},
			expected: []config.Var{
				{Name: "user_id"},
				{Name: "api_key_v2"},
				{Name: "merchant_id_123"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertVariables(tt.vars)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertQueryToString(t *testing.T) {
	tests := []struct {
		name     string
		query    map[string]string
		expected string
	}{
		{
			name:     "empty query",
			query:    map[string]string{},
			expected: "",
		},
		{
			name:     "nil query",
			query:    nil,
			expected: "",
		},
		{
			name: "single query parameter",
			query: map[string]string{
				"page": "1",
			},
			expected: `{"page":"1"}`,
		},
		{
			name: "multiple query parameters",
			query: map[string]string{
				"page":  "1",
				"limit": "10",
				"sort":  "desc",
			},
			expected: `{"limit":"10","page":"1","sort":"desc"}`,
		},
		{
			name: "query parameters with special characters",
			query: map[string]string{
				"search":  "test query",
				"filter":  "status=active",
				"include": "user,orders",
			},
			expected: `{"filter":"status=active","include":"user,orders","search":"test query"}`,
		},
		{
			name: "query parameters with numbers",
			query: map[string]string{
				"id":     "12345",
				"amount": "100.50",
				"count":  "0",
			},
			expected: `{"amount":"100.50","count":"0","id":"12345"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertQueryToString(tt.query)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertInitVars(t *testing.T) {
	tests := []struct {
		name     string
		vars     []InitVar
		expected []config.InitVar
	}{
		{
			name:     "empty init vars",
			vars:     []InitVar{},
			expected: nil,
		},
		{
			name:     "nil init vars",
			vars:     nil,
			expected: nil,
		},
		{
			name: "single init var",
			vars: []InitVar{
				{
					Name:     "user_id",
					FromStep: "create_user",
					Field:    "data.id",
				},
			},
			expected: []config.InitVar{
				{
					Name:     "user_id",
					FromStep: "create_user",
					Field:    "data.id",
				},
			},
		},
		{
			name: "multiple init vars",
			vars: []InitVar{
				{
					Name:     "user_id",
					FromStep: "create_user",
					Field:    "data.id",
				},
				{
					Name:     "merchant_id",
					FromStep: "create_merchant",
					Field:    "merchant.id",
				},
				{
					Name:     "api_key",
					FromStep: "get_api_key",
					Field:    "key",
				},
			},
			expected: []config.InitVar{
				{
					Name:     "user_id",
					FromStep: "create_user",
					Field:    "data.id",
				},
				{
					Name:     "merchant_id",
					FromStep: "create_merchant",
					Field:    "merchant.id",
				},
				{
					Name:     "api_key",
					FromStep: "get_api_key",
					Field:    "key",
				},
			},
		},
		{
			name: "init vars with special characters",
			vars: []InitVar{
				{
					Name:     "user_id_v2",
					FromStep: "create_user_v2",
					Field:    "data.user.id",
				},
				{
					Name:     "api_key_123",
					FromStep: "get_api_key_v1",
					Field:    "keys.primary",
				},
			},
			expected: []config.InitVar{
				{
					Name:     "user_id_v2",
					FromStep: "create_user_v2",
					Field:    "data.user.id",
				},
				{
					Name:     "api_key_123",
					FromStep: "get_api_key_v1",
					Field:    "keys.primary",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertInitVars(tt.vars)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertHook(t *testing.T) {
	tests := []struct {
		name     string
		hook     Hook
		expected config.Hook
	}{
		{
			name: "empty hook",
			hook: Hook{
				Before: HookActions{},
				After:  HookActions{},
			},
			expected: config.Hook{
				Before: config.HookActions{},
				After:  config.HookActions{},
			},
		},
		{
			name: "hook with before actions only",
			hook: Hook{
				Before: HookActions{
					InitVars: []InitVar{
						{
							Name:     "setup_var",
							FromStep: "setup",
							Field:    "data.value",
						},
					},
					Workflows: []TestStep{
						{
							Name:   "setup_request",
							Method: "POST",
							URI:    "/setup",
						},
					},
				},
				After: HookActions{},
			},
			expected: config.Hook{
				Before: config.HookActions{
					InitVars: []config.InitVar{
						{
							Name:     "setup_var",
							FromStep: "setup",
							Field:    "data.value",
						},
					},
					Workflows: []config.Workflow{
						{
							Step: "setup_request",
							Request: config.Request{
								Method: "POST",
								URI:    "/setup",
							},
							ExpectResponse: config.Response{
								StatusCode: 0,
							},
						},
					},
				},
				After: config.HookActions{},
			},
		},
		{
			name: "hook with after actions only",
			hook: Hook{
				Before: HookActions{},
				After: HookActions{
					InitVars: []InitVar{
						{
							Name:     "cleanup_var",
							FromStep: "cleanup",
							Field:    "result.status",
						},
					},
					Workflows: []TestStep{
						{
							Name:           "cleanup_request",
							Method:         "DELETE",
							URI:            "/cleanup",
							ExpectedStatus: 204,
						},
					},
				},
			},
			expected: config.Hook{
				Before: config.HookActions{},
				After: config.HookActions{
					InitVars: []config.InitVar{
						{
							Name:     "cleanup_var",
							FromStep: "cleanup",
							Field:    "result.status",
						},
					},
					Workflows: []config.Workflow{
						{
							Step: "cleanup_request",
							Request: config.Request{
								Method: "DELETE",
								URI:    "/cleanup",
							},
							ExpectResponse: config.Response{
								StatusCode: 204,
							},
						},
					},
				},
			},
		},
		{
			name: "complete hook with both before and after",
			hook: Hook{
				Before: HookActions{
					InitVars: []InitVar{
						{
							Name:     "user_id",
							FromStep: "create_user",
							Field:    "data.id",
						},
					},
					Workflows: []TestStep{
						{
							Name:           "create_user",
							Method:         "POST",
							URI:            "/users",
							ExpectedStatus: 201,
						},
					},
				},
				After: HookActions{
					InitVars: []InitVar{
						{
							Name:     "cleanup_id",
							FromStep: "cleanup",
							Field:    "id",
						},
					},
					Workflows: []TestStep{
						{
							Name:           "delete_user",
							Method:         "DELETE",
							URI:            "/users/{{.user_id}}",
							ExpectedStatus: 204,
						},
					},
				},
			},
			expected: config.Hook{
				Before: config.HookActions{
					InitVars: []config.InitVar{
						{
							Name:     "user_id",
							FromStep: "create_user",
							Field:    "data.id",
						},
					},
					Workflows: []config.Workflow{
						{
							Step: "create_user",
							Request: config.Request{
								Method: "POST",
								URI:    "/users",
							},
							ExpectResponse: config.Response{
								StatusCode: 201,
							},
						},
					},
				},
				After: config.HookActions{
					InitVars: []config.InitVar{
						{
							Name:     "cleanup_id",
							FromStep: "cleanup",
							Field:    "id",
						},
					},
					Workflows: []config.Workflow{
						{
							Step: "delete_user",
							Request: config.Request{
								Method: "DELETE",
								URI:    "/users/{{.user_id}}",
							},
							ExpectResponse: config.Response{
								StatusCode: 204,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertHook(tt.hook)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertHookActions(t *testing.T) {
	tests := []struct {
		name     string
		actions  HookActions
		expected config.HookActions
	}{
		{
			name: "empty hook actions",
			actions: HookActions{
				InitVars:  []InitVar{},
				Workflows: []TestStep{},
			},
			expected: config.HookActions{
				InitVars:  nil,
				Workflows: nil,
			},
		},
		{
			name: "hook actions with init vars only",
			actions: HookActions{
				InitVars: []InitVar{
					{
						Name:     "test_var",
						FromStep: "test_step",
						Field:    "data.value",
					},
				},
				Workflows: []TestStep{},
			},
			expected: config.HookActions{
				InitVars: []config.InitVar{
					{
						Name:     "test_var",
						FromStep: "test_step",
						Field:    "data.value",
					},
				},
				Workflows: nil,
			},
		},
		{
			name: "hook actions with workflows only",
			actions: HookActions{
				InitVars: []InitVar{},
				Workflows: []TestStep{
					{
						Name:           "test_workflow",
						Method:         "GET",
						URI:            "/test",
						ExpectedStatus: 200,
					},
				},
			},
			expected: config.HookActions{
				InitVars: nil,
				Workflows: []config.Workflow{
					{
						Step: "test_workflow",
						Request: config.Request{
							Method: "GET",
							URI:    "/test",
						},
						ExpectResponse: config.Response{
							StatusCode: 200,
						},
					},
				},
			},
		},
		{
			name: "complete hook actions",
			actions: HookActions{
				InitVars: []InitVar{
					{
						Name:     "user_id",
						FromStep: "create_user",
						Field:    "data.id",
					},
					{
						Name:     "api_key",
						FromStep: "get_key",
						Field:    "key",
					},
				},
				Workflows: []TestStep{
					{
						Name:           "create_user",
						Method:         "POST",
						URI:            "/users",
						Headers:        map[string]string{"Content-Type": "application/json"},
						Body:           map[string]interface{}{"name": "test"},
						Query:          map[string]string{"type": "user"},
						Variables:      []string{"user_id"},
						ExpectedStatus: 201,
					},
					{
						Name:           "get_api_key",
						Method:         "GET",
						URI:            "/keys",
						ExpectedStatus: 200,
					},
				},
			},
			expected: config.HookActions{
				InitVars: []config.InitVar{
					{
						Name:     "user_id",
						FromStep: "create_user",
						Field:    "data.id",
					},
					{
						Name:     "api_key",
						FromStep: "get_key",
						Field:    "key",
					},
				},
				Workflows: []config.Workflow{
					{
						Step: "create_user",
						Request: config.Request{
							Method:  "POST",
							URI:     "/users",
							Headers: `{"Content-Type":"application/json"}`,
							Body:    `{"name":"test"}`,
							Query:   `{"type":"user"}`,
							Vars: []config.Var{
								{Name: "user_id"},
							},
						},
						ExpectResponse: config.Response{
							StatusCode: 201,
						},
					},
					{
						Step: "get_api_key",
						Request: config.Request{
							Method: "GET",
							URI:    "/keys",
						},
						ExpectResponse: config.Response{
							StatusCode: 200,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertHookActions(tt.actions)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertFromResponses(t *testing.T) {
	tests := []struct {
		name          string
		fromResponses []FromResponse
		expected      []config.FromResponse
	}{
		{
			name:          "empty from responses",
			fromResponses: []FromResponse{},
			expected:      []config.FromResponse{},
		},
		{
			name:          "nil from responses",
			fromResponses: nil,
			expected:      nil,
		},
		{
			name: "single from response",
			fromResponses: []FromResponse{
				{
					Step:      "create_user",
					Name:      "user_id",
					FromField: "data.id",
				},
			},
			expected: []config.FromResponse{
				{
					Step:      "create_user",
					Name:      "user_id",
					FromField: "data.id",
				},
			},
		},
		{
			name: "multiple from responses",
			fromResponses: []FromResponse{
				{
					Step:      "create_user",
					Name:      "user_id",
					FromField: "data.id",
				},
				{
					Step:      "create_merchant",
					Name:      "merchant_id",
					FromField: "merchant.id",
				},
				{
					Step:      "get_api_key",
					Name:      "api_key",
					FromField: "key",
				},
			},
			expected: []config.FromResponse{
				{
					Step:      "create_user",
					Name:      "user_id",
					FromField: "data.id",
				},
				{
					Step:      "create_merchant",
					Name:      "merchant_id",
					FromField: "merchant.id",
				},
				{
					Step:      "get_api_key",
					Name:      "api_key",
					FromField: "key",
				},
			},
		},
		{
			name: "from responses with special characters",
			fromResponses: []FromResponse{
				{
					Step:      "create_user_v2",
					Name:      "user_id_v2",
					FromField: "data.user.id",
				},
				{
					Step:      "get_api_key_v1",
					Name:      "api_key_123",
					FromField: "keys.primary",
				},
			},
			expected: []config.FromResponse{
				{
					Step:      "create_user_v2",
					Name:      "user_id_v2",
					FromField: "data.user.id",
				},
				{
					Step:      "get_api_key_v1",
					Name:      "api_key_123",
					FromField: "keys.primary",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertFromResponses(tt.fromResponses)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests for performance testing
func BenchmarkConvertHeadersToString(b *testing.B) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token123",
		"User-Agent":    "TestClient/1.0",
		"Accept":        "application/json",
		"X-Request-ID":  "req-12345",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		convertHeadersToString(headers)
	}
}

func BenchmarkConvertBodyToString(b *testing.B) {
	body := map[string]interface{}{
		"user_id":     "#{{.user_id}}",
		"merchant_id": "#{{.merchant_id}}",
		"name":        "{{.name}}",
		"status":      "#{{.status}}",
		"amount":      100,
		"metadata": map[string]interface{}{
			"created_at": "#{{.created_at}}",
			"updated_at": "{{.updated_at}}",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		convertBodyToString(body)
	}
}

func BenchmarkConvertVariables(b *testing.B) {
	vars := []string{"user_id", "merchant_id", "api_key", "request_id", "session_id"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		convertVariables(vars)
	}
}

func BenchmarkConvertQueryToString(b *testing.B) {
	query := map[string]string{
		"page":    "1",
		"limit":   "10",
		"sort":    "desc",
		"filter":  "status=active",
		"include": "user,orders",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		convertQueryToString(query)
	}
}
