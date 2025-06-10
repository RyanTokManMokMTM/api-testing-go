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
//   - body: The request body data as map[string]interface{}
//
// Returns:
//   - "null" if body is nil
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
func convertBodyToString(body map[string]interface{}) string {
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

// convertInitVars converts our InitVar type to config.InitVar
// This function is used to transform the initialization variables from our internal format
// to the configuration format used by the test runner.
// Parameters:
//   - vars: A slice of InitVar containing variable initialization information
//
// Returns:
//   - A slice of config.InitVar with the same information in the config format
//   - nil if the input slice is empty
func convertInitVars(vars []InitVar) []config.InitVar {
	if len(vars) == 0 {
		return nil
	}
	result := make([]config.InitVar, len(vars))
	for i, v := range vars {
		result[i] = config.InitVar{
			Name:     v.Name,     // The name of the variable to be initialized
			FromStep: v.FromStep, // The step from which to get the value
			Field:    v.Field,    // The field path to extract the value from
		}
	}
	return result
}

// convertHook converts our Hook type to config.Hook
// This function transforms the hook configuration from our internal format
// to the configuration format used by the test runner.
// A hook consists of before and after actions that are executed
// before and after the main test workflow.
// Parameters:
//   - hook: A Hook struct containing before and after actions
//
// Returns:
//   - A config.Hook with the same information in the config format
func convertHook(hook Hook) config.Hook {
	return config.Hook{
		Before: convertHookActions(hook.Before), // Convert before hook actions
		After:  convertHookActions(hook.After),  // Convert after hook actions
	}
}

// convertHookActions converts our HookActions type to config.HookActions
// This function transforms the hook actions from our internal format
// to the configuration format used by the test runner.
// Hook actions include initialization variables and workflows to be executed.
// Parameters:
//   - actions: A HookActions struct containing init vars and workflows
//
// Returns:
//   - A config.HookActions with the same information in the config format
func convertHookActions(actions HookActions) config.HookActions {
	return config.HookActions{
		// Convert initialization variables
		InitVars: convertInitVars(actions.InitVars),
		// Convert workflows using an anonymous function to handle nil case
		Workflows: func() []config.Workflow {
			if len(actions.Workflows) == 0 {
				return nil
			}
			// Create a new slice to store converted workflows
			result := make([]config.Workflow, len(actions.Workflows))
			for i, wf := range actions.Workflows {
				// Convert each workflow step
				result[i] = config.Workflow{
					Step: wf.Name,
					// Convert request details
					Request: config.Request{
						Method:       wf.Method,                              // HTTP method (GET, POST, etc.)
						URI:          wf.URI,                                 // API endpoint
						Headers:      convertHeadersToString(wf.Headers),     // Convert headers to string format
						Body:         convertBodyToString(wf.Body),           // Convert body to string format
						Query:        convertQueryToString(wf.Query),         // Convert query parameters to string format
						Vars:         convertVariables(wf.Variables),         // Convert variables
						FromResponse: convertFromResponses(wf.FromResponses), // Convert from_response
					},
					// Convert expected response details
					ExpectResponse: config.Response{
						StatusCode: wf.ExpectedStatus, // Expected HTTP status code
					},
				}
			}
			return result
		}(),
	}
}

// convertFromResponses converts FromResponse slice to config.FromResponse slice
func convertFromResponses(fromResponses []FromResponse) []config.FromResponse {
	if fromResponses == nil {
		return nil
	}
	result := make([]config.FromResponse, len(fromResponses))
	for i, fr := range fromResponses {
		result[i] = config.FromResponse{
			Step:      fr.Step,
			Name:      fr.Name,
			FromField: fr.FromField,
		}
	}
	return result
}
