# Workflow Generator Guide

A comprehensive guide on how to use the workflow generator tool to create automated API test workflows.

## 📋 Overview

The workflow generator is a powerful tool that allows you to programmatically generate API test workflows from Go code. It converts structured test definitions into YAML configuration files that can be executed by the API test runner.

## 🏗️ Architecture

The workflow generator consists of several key components:

```
utils/tool/
├── cmd/           # Command line interface
├── config/        # Configuration generation utilities
├── generator/     # Base generator functionality
└── workflow/      # Workflow-specific generation logic
```

## 🚀 Quick Start

### 1. Basic Workflow Generation

To generate workflows, run the generator command:

```bash
make gen-workflow
```

This will execute the main generator and create workflow files in `config/etc/api-test/workflows/`.

### 2. Understanding the Generation Process

The generator follows this flow:

1. **Define Test Cases**: Create test case structures in Go code
2. **Configure Generation Options**: Set up host, network settings, etc.
3. **Generate YAML**: Convert Go structures to YAML configuration
4. **Output Files**: Save generated workflows to the config directory

## 📝 Creating Custom Workflows

### Step 1: Define Test Steps

Create test steps using the `TestStep` struct:

```go
TestStep{
    Name:           "step_name",
    Method:         "GET",
    URI:            "/api/endpoint",
    Headers:        map[string]string{"Accept": "application/json"},
    Body:           map[string]interface{}{"key": "value"},
    Query:          map[string]string{"param": "value"},
    Variables:      []string{"var1", "var2"},
    FromResponses:  []FromResponse{},
    ExpectedStatus: 200,
    ResponseChecks: []ResponseCheck{
        {
            CheckType: CheckTypePresent,
            Field:     "data",
        },
        {
            CheckType: CheckTypeMatches,
            Field:     "data.id",
            Regex:     "^[0-9]+$",
        },
    },
}
```

### Step 2: Define Response Checks

Available response check types:

```go
// Check if field exists
ResponseCheck{
    CheckType: CheckTypePresent,
    Field:     "data.user.name",
}

// Check if field matches regex
ResponseCheck{
    CheckType: CheckTypeMatches,
    Field:     "data.email",
    Regex:     "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
}

// Check if field equals value
ResponseCheck{
    CheckType: CheckTypeEquals,
    Field:     "data.status",
    Value:     "active",
    Type:      StringType,
}

// Check if numeric field is greater than
ResponseCheck{
    CheckType: CheckTypeGreaterThan,
    Field:     "data.count",
    Value:     10,
}

// Check if field does not exist
ResponseCheck{
    CheckType: CheckTypeNotPresent,
    Field:     "data.error",
}
```

### Step 3: Create Test Cases

Combine test steps into test cases:

```go
TestCase{
    Name: "My API Test Workflow",
    Steps: []TestStep{
        step1,
        step2,
        step3,
    },
}
```

### Step 4: Add Generation Options

Configure how the workflow should be generated:

```go
TestCaseWithOptions{
    TestCase: testCase,
    Options: []ScenarioOption{
        WithSkipScenario(false),
        WithScenarioHook(hook),
    },
}
```

## 🔧 Generation Options

### Global Options

```go
// Set the host for all requests
WithHost("https://api.example.com")

// Enable network access for real API calls
WithNetworkEnable(true)

// Skip the entire test suite
WithSkip(false)

// Add global hooks
WithGlobalHook(hook)
```

### Scenario Options

```go
// Skip specific scenarios
WithSkipScenario(true)

// Add scenario-specific hooks
WithScenarioHook(hook)
```

## 🌟 Example: Complete Workflow Generator

Here's a complete example of how to create a workflow generator function:

```go
func (g *WorkflowGenerator) GenerateMyWorkflow() []TestCaseWithOptions {
    return []TestCaseWithOptions{
        {
            TestCase: TestCase{
                Name: "User Management API Test",
                Steps: []TestStep{
                    {
                        Name:   "create_user",
                        Method: "POST",
                        URI:    "/api/users",
                        Headers: map[string]string{
                            "Content-Type": "application/json",
                        },
                        Body: map[string]interface{}{
                            "name":  "John Doe",
                            "email": "john@example.com",
                        },
                        ExpectedStatus: 201,
                        ResponseChecks: []ResponseCheck{
                            {
                                CheckType: CheckTypePresent,
                                Field:     "data.id",
                            },
                            {
                                CheckType: CheckTypePresent,
                                Field:     "data.name",
                            },
                        },
                    },
                    {
                        Name:   "get_user",
                        Method: "GET",
                        URI:    "/api/users/{{.user_id}}",
                        Variables: []string{"user_id"},
                        FromResponses: []FromResponse{
                            {
                                Step:      "create_user",
                                Name:      "user_id",
                                FromField: "data.id",
                            },
                        },
                        ExpectedStatus: 200,
                        ResponseChecks: []ResponseCheck{
                            {
                                CheckType: CheckTypeEquals,
                                Field:     "data.name",
                                Value:     "John Doe",
                            },
                        },
                    },
                },
            },
        },
    }
}
```

## 🔄 Adding to Generator

### Step 1: Add to GenerateAllWorkflows

Add your new workflow to the main generation function:

```go
func (g *WorkflowGenerator) GenerateAllWorkflows() error {
    workflows := []struct {
        name     string
        generate func() []TestCaseWithOptions
        opts     []GenerateOption
    }{
        {
            name:     "my_workflow",
            generate: g.GenerateMyWorkflow,
            opts: []GenerateOption{
                WithSkip(false),
                WithNetworkEnable(true),
                WithHost("https://api.example.com"),
            },
        },
    }
    
    // ... rest of the function
}
```

### Step 2: Generate and Test

1. Run the generator:
```bash
make gen-workflow
```

2. Check the generated file:
```bash
cat config/etc/api-test/workflows/my_workflow.yaml
```

3. Run the tests:
```bash
make test/api-junit
```

## 📊 Generated YAML Structure

The generator creates YAML files with this structure:

```yaml
api_test:
  host: https://api.example.com
  network_enable: true
  name: my_workflow
  skip: false
  scenarios:
    - name: User Management API Test
      skip: false
      workflows:
        - step: create_user
          request:
            method: POST
            uri: /api/users
            headers: '{"Content-Type":"application/json"}'
            body: '{"name":"John Doe","email":"john@example.com"}'
          expect_response:
            status_code: 201
            body:
              presents:
                - field: data.id
                - field: data.name
```

## 🎯 Best Practices

### 1. Naming Conventions

- Use descriptive step names: `create_user`, `get_user_by_id`
- Use clear test case names: `User Management API Test`
- Use consistent naming patterns across workflows

### 2. Response Validation

- Always validate the presence of expected fields
- Use regex patterns for format validation
- Include both positive and negative test cases
- Validate status codes appropriately

### 3. Variable Management

- Use meaningful variable names
- Extract values from responses for subsequent requests
- Document the data flow between steps

### 4. Error Handling

- Include error scenarios in your workflows
- Test both valid and invalid inputs
- Validate error responses and status codes

## 🔍 Debugging

### Common Issues

1. **YAML Generation Errors**: Check struct field tags and types
2. **Missing Fields**: Ensure all required fields are populated
3. **Invalid Regex**: Test regex patterns separately
4. **Network Issues**: Verify host and network settings

### Debug Commands

```bash
# Generate with verbose output
go run utils/tool/cmd/main.go -v

# Check generated YAML syntax
yamllint config/etc/api-test/workflows/*.yaml

# Test specific workflow
make test/api-junit
```

## 📚 Advanced Features

### Custom Response Validators

You can extend the validation system by adding new check types:

```go
// Add to CheckType enum
const (
    CheckTypeCustom CheckType = "custom"
)

// Implement validation logic in generateWorkflow function
case CheckTypeCustom:
    // Custom validation logic
```

### Dynamic Data Generation

Integrate with external APIs for dynamic test data:

```go
// Use external data generation APIs
URI: "/v2/profiles/person",
Query: map[string]string{"results": "1"},
```

### Hook Integration

Add setup and cleanup hooks:

```go
Hook{
    Before: HookActions{
        InitVars: []InitVar{
            {
                Name:     "test_user_id",
                FromStep: "setup_user",
                Field:    "data.id",
            },
        },
    },
    After: HookActions{
        Workflows: []TestStep{
            {
                Name:   "cleanup_user",
                Method: "DELETE",
                URI:    "/api/users/{{.test_user_id}}",
            },
        },
    },
}
```

## 🎉 Conclusion

The workflow generator provides a powerful and flexible way to create automated API test workflows. By following this guide, you can create comprehensive test suites that are maintainable, readable, and effective.

For more information, refer to the main [README.md](../README.md) and the [apitest documentation](https://apitest.dev/). 