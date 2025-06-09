# API Testing Tool

A powerful and flexible API testing tool written in Go, designed to help you create, manage, and run API test cases efficiently.

## Features

- **YAML-based Test Configuration**: Define your API tests in YAML format for better readability and maintainability
- **Workflow Support**: Create complex test workflows with multiple steps and dependencies
- **Variable Support**: Use variables to make your tests more dynamic and reusable
- **Response Validation**: Multiple ways to validate API responses:
  - Equals check
  - Regex matching
  - Field presence/absence check
  - Greater than/Less than comparisons
- **Global Hooks**: Define before/after hooks for test setup and cleanup
- **Command-based Variables**: Support for dynamic variable values using shell commands

## Project Structure

```
.
├── config/
│   └── etc/
│       └── api-test/
│           ├── config/         # Configuration files
│           └── workflows/      # Test workflow files
├── utils/
│   └── tool/
│       ├── cmd/               # Command-line tools
│       ├── config/            # Configuration generation
│       ├── generator/         # Base YAML generation
│       └── workflow/          # Workflow generation
└── api_test_runner/          # Test runner implementation
```

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/RyanTokManMokMTM/api-testing-go.git
   cd api-testing-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

### Generating Test Files

1. Generate configuration:
   ```bash
   cd utils/tool/cmd
   go run main.go
   ```

This will generate:
- Configuration file at `config/etc/api-test/config/config.yaml`
- Workflow files at `config/etc/api-test/workflows/`

## Test Configuration

### Basic Structure

```yaml
api_test:
  host: '{{.API_HOST}}'
  network_enable: true
  name: your_test_name
  scenarios:
    - name: scenario_name
      workflows:
        - step: step_name
          request:
            method: POST
            uri: /api/endpoint
            headers: '{"Content-Type":"application/json"}'
            body: '{"key":"value"}'
          expect_response:
            code: SUCCESS
            status_code: 200
            body:
              presents:
                - field: data.id
```

### Available Response Checks

1. **Equals Check**:
   ```yaml
   equals:
     - field: data.name
       value: "expected_value"
   ```

2. **Regex Match**:
   ```yaml
   matches:
     - field: data.id
       regex: "^[0-9a-f]{24}$"
   ```

3. **Field Presence**:
   ```yaml
   presents:
     - field: data.id
   ```

4. **Field Absence**:
   ```yaml
   not_presents:
     - field: error
   ```

5. **Numeric Comparison**:
   ```yaml
   greater_thans:
     - field: data.count
       value: 10
   ```

### Using Variables

1. **Static Variables**:
   ```yaml
   config:
     - name: api_key
       value: "your-api-key"
   ```

2. **Command-based Variables**:
   ```yaml
   config:
     - name: next_day
       command: echo $(date -d "tomorrow" +%s%3N)
       type: number
   ```

3. **Response-based Variables**:
   ```yaml
   from_response:
     - step: previous_step
       name: variable_name
       from_field: data.id
   ```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Thanks to all contributors who have helped shape this project
- Inspired by various API testing frameworks and tools 