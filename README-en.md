# API Testing Framework 🚀

A powerful and flexible API testing framework that enables you to define, execute, and maintain API test workflows using YAML configurations. This framework is designed to make API testing more efficient, maintainable, and scalable.

## ✨ Features

- **YAML-based Configuration**: Define test workflows using simple YAML syntax
- **Dynamic Variable Support**: Use pre-defined and response-extracted variables
- **Comprehensive Validation**: Multiple assertion types for thorough API testing
- **Flexible Hooks System**: Global and scenario-specific hooks for setup and cleanup
- **Extensible Architecture**: Easy to add new routes and test scenarios
- **Go-based Implementation**: High performance and reliability

## 📋 Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
- [Core Concepts](#core-concepts)
  - [API Router Configuration](#api-router-configuration)
  - [Test Workflow Structure](#test-workflow-structure)
  - [Variables and Dynamic Values](#variables-and-dynamic-values)
  - [Hooks System](#hooks-system)
- [Writing Tests](#writing-tests)
  - [Basic Test Structure](#basic-test-structure)
  - [Request Configuration](#request-configuration)
  - [Response Validation](#response-validation)
  - [Advanced Features](#advanced-features)
- [Best Practices](#best-practices)
- [Contributing](#contributing)
- [Support](#support)

## 🚀 Getting Started

### Prerequisites

- Go 1.x or higher
- Basic understanding of REST APIs
- Familiarity with YAML syntax
- Git for version control

### Installation

1. Clone the repository:
```bash
git clone https://github.com/your-org/github.com/RyanTokManMokMTM/api-testing-go.git
cd github.com/RyanTokManMokMTM/api-testing-go
```

2. Install dependencies:
```bash
go mod download
```

### Quick Start

1. Navigate to the project directory
2. Run the test suite:
```bash
go test ./test/...
```

3. Check the example test workflows in `/config/etc/api-test/` to understand the structure

## 🎯 Core Concepts

### API Router Configuration

#### Available Routes
```yaml
- Order Route: /api/merchants/:mid/orders
- Subscription Route: /api/merchants/:mid/subscriptions
- Subscribable Entity Route: /api/subscribables_entities
```

#### Adding New Routes

1. Define the route in `/utils/util/var.go`:
```go
const (
    // Existing routes
    SubscriptionPrefix = "/api"
    SubscribablePrefix = "/api/subscribables_entities"
    OrderPrefix = "/api/merchants/:mid/orders"

    // Add your new route
    NewRoutePrefix = "/api/your/new/route"
)
```

2. Register the route in `/test/api_test_suite_test.go`:
```go
func initRoute() {
    // ... existing code ...
    
    // Register your new route
    newService := new.NewService(uow, eventService)
    newHandler := newhandler.NewHandler(newService)
    newHandler.Router(app.Group(apiutil.NewRoutePrefix))
}
```

## 📝 Writing Tests

### Basic Test Structure

Create your test workflow YAML file in `/config/etc/api-test/`:

```yaml
api_test:
  name: my_test_workflow
  description: "Description of what this test workflow does"
  scenarios:
    - name: test_scenario
      description: "Description of this specific scenario"
      workflow:
        - step: health_check
          request:
            method: GET
            headers: '{"Content-Type": "application/json"}'
            uri: /health_check
          expect_response:
            code: SUCCESS
            status_code: 200
```

### Request Configuration

```yaml
request:
  method: POST                    # HTTP method (GET, POST, PUT, PATCH, DELETE)
  uri: /api/endpoint             # API endpoint
  headers: '{"key": "value"}'    # Request headers
  query: '{"param": "value"}'    # Query parameters
  body: '{"data": "value"}'      # Request body
  timeout: 30                    # Request timeout in seconds (optional)
  retry:                         # Retry configuration (optional)
    attempts: 3
    delay: 1
```

### Response Validation

```yaml
expect_response:
  code: SUCCESS                  # Expected response code
  status_code: 200              # Expected HTTP status code
  timeout: 5                    # Response timeout in seconds
  equals:                       # Exact value matching
    - field: data.id
      value: "expected-id"
      message: "ID should match expected value"
  matches:                      # Regex pattern matching
    - field: data.id
      value: "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
      message: "ID should match UUID format"
  presents:                     # Field existence check
    - field: data.id
      message: "ID field should be present"
  not_presents:                 # Field non-existence check
    - field: data.error
      message: "Error field should not be present"
```

### Variables and Dynamic Values

#### Pre-defined Variables
Define variables in `/config/etc/api-test/config/config.yaml`:
```yaml
environment:
  development:
    mid: "6552f99b99821c568c0115cc"
    legacy_id: "SL101PRO6828965740321448017_SKU6828965740942204962"
  production:
    mid: "your-production-mid"
    legacy_id: "your-production-legacy-id"
```

#### Dynamic Variables from Responses
```yaml
request:
  vars:
    - name: mid
  from_response:
    - step: create_subscription
      name: subscription_id
      from_field: data.id
      default: "fallback-value"    # Optional default value
  uri: /api/merchants/{{.mid}}/subscriptions/{{.subscription_id}}
```

### Hooks System

#### Global Hooks
```yaml
api_test:
  name: workflow_with_hooks
  description: "Test workflow with global hooks"
  global_hook:
    before:
      workflows:
        - step: setup_environment
          description: "Prepare test environment"
    after:
      workflows:
        - step: cleanup_environment
          description: "Clean up test environment"
  scenarios:
    - name: test_scenario
      # ... scenario definition
```

#### Scenario-specific Hooks
```yaml
api_test:
  name: workflow_with_scenario_hooks
  scenarios:
    - name: test_scenario
      hook:
        before:
          workflows:
            - step: scenario_setup
              description: "Prepare scenario-specific data"
        after:
          workflows:
            - step: scenario_cleanup
              description: "Clean up scenario-specific data"
      workflows:
        - step: main_test_step
```

## 💡 Best Practices

1. **Organization**
   - Group related scenarios logically
   - Use clear, descriptive names for scenarios and steps
   - Include descriptions for workflows and scenarios
   - Maintain a consistent directory structure

2. **Variables Management**
   - Use environment-specific variables
   - Document all variable dependencies
   - Provide default values where appropriate
   - Use meaningful variable names

3. **Test Design**
   - Keep tests independent and isolated
   - Use hooks for setup and cleanup
   - Include both positive and negative test cases
   - Validate both success and error scenarios

4. **Validation Strategy**
   - Validate status codes and response bodies
   - Use appropriate assertion types
   - Include meaningful error messages
   - Test edge cases and error conditions

5. **Maintenance**
   - Regular review of test cases
   - Remove obsolete tests
   - Update tests when APIs change
   - Document test dependencies

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 📞 Support

- **Documentation**: [Link to detailed documentation]
- **Issues**: [GitHub Issues](https://github.com/your-org/github.com/RyanTokManMokMTM/api-testing-go/issues)
- **Contact**: development-team@your-org.com

---

Made with ❤️ by Your Organization