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
		body     interface{}
		expected string
	}{
		{
			name:     "nil body",
			body:     nil,
			expected: "",
		},
		{
			name:     "empty string body",
			body:     "",
			expected: `""`,
		},
		{
			name:     "simple string body",
			body:     "hello world",
			expected: `"hello world"`,
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
