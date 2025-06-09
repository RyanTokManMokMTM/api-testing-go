package workflow

import (
	"encoding/json"
	"regexp"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
)

// convertHeadersToString converts a map of HTTP headers to a JSON string representation.
// This function is used to serialize headers for YAML configuration files.
//
// Parameters:
//   - headers: A map where keys are header names and values are header values
//
// Returns:
//   - Empty string if headers map is empty or nil
//   - JSON string representation of the headers map
//
// Example:
//
//	headers := map[string]string{"Content-Type": "application/json", "Authorization": "Bearer token"}
//	result := convertHeadersToString(headers)
//	// result: {"Content-Type":"application/json","Authorization":"Bearer token"}
func convertHeadersToString(headers map[string]string) string {
	if len(headers) == 0 {
		return ""
	}
	data, _ := json.Marshal(headers)
	return string(data)
}

// convertBodyToString converts a request body interface to a JSON string representation.
// This function handles template variables by removing quotes around them to maintain
// proper YAML template syntax.
//
// Parameters:
//   - body: The request body data (can be any type that can be marshaled to JSON)
//
// Returns:
//   - Empty string if body is nil or marshaling fails
//   - JSON string with template variables properly formatted
//
// Template Variable Handling:
//   - Template variables in the format {{.variable}} are preserved without quotes
//   - This allows the YAML parser to properly interpret template variables
//   - Variables marked with # prefix (e.g., "#{{.variable}}") will have quotes removed
//
// Example:
//
//	body := map[string]interface{}{
//	  "user_id": "#{{.user_id}}",  // # prefix indicates quote removal
//	  "amount": 100,
//	  "name": "{{.name}}"          // No # prefix, quotes will be preserved
//	}
//	result := convertBodyToString(body)
//	// result: {"user_id":{{.user_id}},"amount":100,"name":"{{.name}}"}
func convertBodyToString(body interface{}) string {
	if body == nil {
		return ""
	}

	// Serialize to JSON
	data, err := json.Marshal(body)
	if err != nil {
		return ""
	}

	// Process the serialized string to remove quotes from template variables
	// Variables with # prefix (e.g., "#{{.variable}}") will have quotes removed
	// This ensures template variables are properly formatted for YAML parsing
	result := string(data)

	// Remove quotes from variables marked with # prefix
	// Pattern: "#{{.variable}}" -> {{.variable}}
	re := regexp.MustCompile(`"#(\{\{\.\w+\}\})"`)
	result = re.ReplaceAllString(result, "$1")

	return result
}

// convertVariables converts a slice of variable names to a slice of config.Var structs.
// This function is used to transform simple variable names into the structured format
// required by the configuration system.
//
// Parameters:
//   - vars: A slice of variable names as strings
//
// Returns:
//   - nil if the input slice is empty or nil
//   - A slice of config.Var structs with the Name field populated
//
// Example:
//
//	vars := []string{"user_id", "merchant_id", "api_key"}
//	result := convertVariables(vars)
//	// result: []config.Var{{Name: "user_id"}, {Name: "merchant_id"}, {Name: "api_key"}}
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

// convertQueryToString converts a map of query parameters to a JSON string representation.
// This function is used to serialize query parameters for YAML configuration files.
//
// Parameters:
//   - query: A map where keys are parameter names and values are parameter values
//
// Returns:
//   - Empty string if query map is empty or nil
//   - JSON string representation of the query parameters map
//
// Example:
//
//	query := map[string]string{"page": "1", "limit": "10", "sort": "desc"}
//	result := convertQueryToString(query)
//	// result: {"page":"1","limit":"10","sort":"desc"}
func convertQueryToString(query map[string]string) string {
	if len(query) == 0 {
		return ""
	}
	data, _ := json.Marshal(query)
	return string(data)
}
