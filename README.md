# API Test Generator

A comprehensive API testing framework written in Go, designed to generate and execute automated API test workflows with support for data generation APIs and complex test scenarios.

## 🚀 Features

- **Automated Test Generation**: Generate test workflows from configuration files
- **Data Generation API Support**: Built-in support for external data generation APIs (e.g., Spanish DNI generator)
- **Flexible Configuration**: YAML-based configuration for easy test setup
- **Response Validation**: Multiple validation types including presence checks, regex matching, and equality comparisons
- **Hook System**: Before and after hooks for test setup and cleanup
- **Variable Management**: Dynamic variable extraction from API responses
- **Network Support**: Configurable network access for real API testing
- **Comprehensive Reporting**: JSON and JUnit XML report formats
- **Built on apitest**: Leverages [apitest](https://apitest.dev/) for robust HTTP testing with JSON path assertions and sequence diagrams

## 📋 Prerequisites

- Go 1.23.0 or higher
- Git

## 🛠️ Installation

1. Clone the repository:
```bash
git clone https://github.com/RyanTokManMokMTM/api-testing-go.git
cd api-testing-go
```

2. Install dependencies:
```bash
go mod download
```

3. (Optional) Install Ginkgo for enhanced testing:
```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest
```

## 🏗️ Project Structure

```
api-testing-go/
├── api_test_runner/     # Core API test runner (built on apitest)
├── config/             # Configuration files and generated test suites
│   └── etc/api-test/
│       ├── config/     # Global configuration
│       └── workflows/  # Generated test workflows
├── test/               # Test files
├── utils/              # Utility packages
│   ├── error/          # Error definitions
│   ├── helper/         # Helper functions
│   ├── tool/           # Test generation tools
│   │   ├── cmd/        # Command line tools
│   │   ├── config/     # Configuration tools
│   │   ├── generator/  # Test case generator
│   │   └── workflow/   # Workflow generation
│   └── util/           # General utilities
├── go.mod              # Go module file
├── go.sum              # Go dependencies checksum
├── makefile            # Build and test commands
└── README.md           # This file
```

### Technology Stack

- **Core Testing**: [apitest](https://apitest.dev/) - Go HTTP testing library with JSON path assertions
- **Test Framework**: [Ginkgo](https://onsi.github.io/ginkgo/) - BDD testing framework for Go
- **Configuration**: YAML-based configuration system
- **Data Generation**: Integration with external APIs for test data generation

## 🚀 Quick Start

### 1. Generate Test Workflows

Generate test workflows from predefined templates:

```bash
make gen
```

This will generate test configurations in `config/etc/api-test/workflows/`.

### 2. Run Tests

Run all tests:

```bash
# Using go test (recommended)
go test ./test/...

# Using Ginkgo (if installed)
make test/api
```

### 3. Generate Reports

Generate test reports in various formats:

```bash
# JSON report
make test-json

# JUnit XML report
make test/api
```

## 📝 Configuration

### Test Workflow Configuration

Test workflows are defined in YAML format:

```yaml
api_test:
  host: https://api.example.com
  network_enable: true
  name: example_workflow
  skip: false
  scenarios:
    - name: Test Scenario
      skip: false
      workflows:
        - step: test_step_name
          request:
            method: GET
            uri: /api/endpoint
            headers: '{"Accept":"application/json"}'
            query: '{"param":"value"}'
          expect_response:
            status_code: 200
            body:
              presents:
                - field: data
              matches:
                - field: data.id
                  regex: "^[0-9]+$"
```

### Response Validation Types

- **Present**: Check if a field exists in the response
- **Not Present**: Check if a field does not exist
- **Equals**: Check if a field equals a specific value
- **Matches**: Check if a field matches a regex pattern
- **Greater Than**: Check if a numeric field is greater than a value
- **Less Than**: Check if a numeric field is less than a value

## 🔧 Available Commands

### Make Commands

```bash
# Generate test workflows
make gen-workflow

# Run tests with JUnit XML report
make test/api-junit

# Run tests with JSON report
make test/api-json
```

## 🌟 Example: Spanish Data Generation API

The framework includes built-in support for the [Spanish Data Generation API](https://api.generadordni.es/), which can generate:

- Personal profiles (DNI, names, addresses)
- Bank accounts (IBAN, BIC)
- Credit cards
- Company profiles (CIF, company names)
- Vehicle plate numbers

### Generated Workflow Example

```yaml
api_test:
  host: https://api.generadordni.es
  network_enable: true
  name: example_workflow
  skip: false
  scenarios:
    - name: Test Spanish Data Generation API Workflow
      workflows:
        - step: generate_person_profile
          request:
            method: GET
            uri: /v2/profiles/person
            headers: '{"Accept":"application/json"}'
            query: '{"results":"1"}'
          expect_response:
            status_code: 200
            body:
              presents:
                - field: data
                - field: data[0].name
                - field: data[0].surname
                - field: data[0].dni
```

## 🔄 Hook System

The framework supports before and after hooks for test setup and cleanup:

```yaml
global_hook:
  before:
    init_vars:
      - name: setup_var
        from_step: setup
        field: data.id
    workflows:
      - step: setup_test
        request:
          method: POST
          uri: /api/setup
  after:
    workflows:
      - step: cleanup_test
        request:
          method: DELETE
          uri: /api/cleanup
```

## 📊 Test Reports

The framework generates comprehensive test reports:

- **JUnit XML**: Compatible with CI/CD systems
- **JSON**: Machine-readable format for custom processing

Reports are saved in the `reports/` directory by default.

## 🛠️ Development

### Adding New Test Workflows

1. Create a new function in `utils/tool/workflow/generator.go`
2. Define your test steps with appropriate validation
3. Add the workflow to the `GenerateAllWorkflows` function
4. Run `make gen` to generate the configuration

### Custom Response Validations

Add new validation types by extending the `ResponseCheck` struct and implementing the validation logic in the workflow generator.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the existing issues in the repository
2. Create a new issue with detailed information
3. Include test configuration and error messages

## 🔗 Related Links

- [Spanish Data Generation API](https://api.generadordni.es/)
- [apitest - Go HTTP Testing Library](https://apitest.dev/)
- [Ginkgo Testing Framework](https://onsi.github.io/ginkgo/)
- [Go Testing Package](https://golang.org/pkg/testing/)
