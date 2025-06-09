package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
		expected  string
	}{
		{
			name:      "custom output directory",
			outputDir: "/custom/path",
			expected:  "/custom/path",
		},
		{
			name:      "empty output directory",
			outputDir: "",
			expected:  DefaultOutputDir,
		},
		{
			name:      "relative path",
			outputDir: "test/output",
			expected:  "test/output",
		},
		{
			name:      "absolute path with spaces",
			outputDir: "/path with spaces",
			expected:  "/path with spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerator(tt.outputDir)
			assert.NotNil(t, gen)
			assert.Equal(t, tt.expected, gen.outputDir)
		})
	}
}

func TestWriteYAML(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		outputDir    string
		filename     string
		data         interface{}
		expectError  bool
		expectedFile string
	}{
		{
			name:      "simple struct",
			outputDir: tempDir,
			filename:  "test",
			data: struct {
				Name string `yaml:"name"`
				Age  int    `yaml:"age"`
			}{
				Name: "John",
				Age:  30,
			},
			expectError:  false,
			expectedFile: "test.yaml",
		},
		{
			name:      "filename with spaces",
			outputDir: tempDir,
			filename:  "Test File Name",
			data: struct {
				Title string `yaml:"title"`
			}{
				Title: "Test Title",
			},
			expectError:  false,
			expectedFile: "test_file_name.yaml",
		},
		{
			name:      "filename with uppercase",
			outputDir: tempDir,
			filename:  "UPPERCASE",
			data: struct {
				Value string `yaml:"value"`
			}{
				Value: "test",
			},
			expectError:  false,
			expectedFile: "uppercase.yaml",
		},
		{
			name:      "complex nested struct",
			outputDir: tempDir,
			filename:  "complex",
			data: struct {
				User struct {
					Name  string `yaml:"name"`
					Email string `yaml:"email"`
				} `yaml:"user"`
				Settings map[string]interface{} `yaml:"settings"`
			}{
				User: struct {
					Name  string `yaml:"name"`
					Email string `yaml:"email"`
				}{
					Name:  "Alice",
					Email: "alice@example.com",
				},
				Settings: map[string]interface{}{
					"enabled": true,
					"count":   42,
				},
			},
			expectError:  false,
			expectedFile: "complex.yaml",
		},
		{
			name:         "nil data",
			outputDir:    tempDir,
			filename:     "nil_test",
			data:         nil,
			expectError:  false,
			expectedFile: "nil_test.yaml",
		},
		{
			name:         "empty string data",
			outputDir:    tempDir,
			filename:     "empty_string",
			data:         "",
			expectError:  false,
			expectedFile: "empty_string.yaml",
		},
		{
			name:      "map data",
			outputDir: tempDir,
			filename:  "map_test",
			data: map[string]interface{}{
				"key1": "value1",
				"key2": 123,
				"key3": []string{"a", "b", "c"},
			},
			expectError:  false,
			expectedFile: "map_test.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerator(tt.outputDir)

			err := gen.WriteYAML(tt.filename, tt.data)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify file was created
				expectedPath := filepath.Join(tt.outputDir, tt.expectedFile)
				assert.FileExists(t, expectedPath)

				// Verify file content can be read back
				content, err := os.ReadFile(expectedPath)
				require.NoError(t, err)
				assert.NotEmpty(t, content)

				// Verify YAML can be unmarshaled back
				var unmarshaled interface{}
				err = yaml.Unmarshal(content, &unmarshaled)
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteYAML_ErrorCases(t *testing.T) {
	tests := []struct {
		name        string
		outputDir   string
		filename    string
		data        interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name:        "invalid output directory",
			outputDir:   "/invalid/path/that/does/not/exist/and/has/very/long/path/name",
			filename:    "test",
			data:        "test data",
			expectError: true,
			errorMsg:    "failed to create output directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerator(tt.outputDir)

			err := gen.WriteYAML(tt.filename, tt.data)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.errorMsg)
		})
	}
}

func TestWriteYAML_UnmarshallableData(t *testing.T) {
	tempDir := t.TempDir()
	gen := NewGenerator(tempDir)

	// Test with channel (cannot be marshaled to YAML)
	// yaml.v2 panics on unmarshallable data, so we need to recover
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for unmarshallable data, but no panic occurred")
			} else {
				// Verify the panic message contains expected content
				panicMsg := fmt.Sprintf("%v", r)
				assert.Contains(t, panicMsg, "cannot marshal type")
			}
		}()

		// This should cause a panic
		gen.WriteYAML("test", make(chan int))
	}()
}

func TestWriteYAML_FileContent(t *testing.T) {
	tempDir := t.TempDir()
	gen := NewGenerator(tempDir)

	testData := struct {
		Name   string            `yaml:"name"`
		Age    int               `yaml:"age"`
		Active bool              `yaml:"active"`
		Tags   []string          `yaml:"tags"`
		Config map[string]string `yaml:"config"`
	}{
		Name:   "Test User",
		Age:    25,
		Active: true,
		Tags:   []string{"tag1", "tag2", "tag3"},
		Config: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	err := gen.WriteYAML("content_test", testData)
	require.NoError(t, err)

	// Read and verify file content
	filePath := filepath.Join(tempDir, "content_test.yaml")
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	// Verify YAML structure
	var result map[string]interface{}
	err = yaml.Unmarshal(content, &result)
	require.NoError(t, err)

	assert.Equal(t, "Test User", result["name"])
	assert.Equal(t, 25, result["age"])
	assert.Equal(t, true, result["active"])
	assert.Equal(t, []interface{}{"tag1", "tag2", "tag3"}, result["tags"])
	assert.Equal(t, map[interface{}]interface{}{"key1": "value1", "key2": "value2"}, result["config"])
}

func TestGetOutputDir(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
		expected  string
	}{
		{
			name:      "custom directory",
			outputDir: "/custom/path",
			expected:  "/custom/path",
		},
		{
			name:      "default directory",
			outputDir: "",
			expected:  DefaultOutputDir,
		},
		{
			name:      "relative path",
			outputDir: "relative/path",
			expected:  "relative/path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewGenerator(tt.outputDir)
			result := gen.GetOutputDir()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWriteYAML_MultipleFiles(t *testing.T) {
	tempDir := t.TempDir()
	gen := NewGenerator(tempDir)

	// Write multiple files
	files := []struct {
		name string
		data interface{}
	}{
		{
			name: "file1",
			data: map[string]string{"key1": "value1"},
		},
		{
			name: "file2",
			data: map[string]int{"count": 42},
		},
		{
			name: "file3",
			data: []string{"item1", "item2", "item3"},
		},
	}

	for _, file := range files {
		err := gen.WriteYAML(file.name, file.data)
		assert.NoError(t, err)
	}

	// Verify all files were created
	for _, file := range files {
		expectedPath := filepath.Join(tempDir, file.name+".yaml")
		assert.FileExists(t, expectedPath)
	}
}

func TestWriteYAML_OverwriteExisting(t *testing.T) {
	tempDir := t.TempDir()
	gen := NewGenerator(tempDir)

	filename := "overwrite_test"
	originalData := "original content"
	updatedData := "updated content"

	// Write original file
	err := gen.WriteYAML(filename, originalData)
	require.NoError(t, err)

	// Verify original content
	filePath := filepath.Join(tempDir, filename+".yaml")
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "original content")

	// Overwrite with new content
	err = gen.WriteYAML(filename, updatedData)
	require.NoError(t, err)

	// Verify updated content
	content, err = os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "updated content")
}

// Benchmark tests
func BenchmarkWriteYAML(b *testing.B) {
	tempDir := b.TempDir()
	gen := NewGenerator(tempDir)

	testData := map[string]interface{}{
		"name":    "Benchmark Test",
		"age":     30,
		"active":  true,
		"tags":    []string{"tag1", "tag2", "tag3"},
		"config":  map[string]string{"key1": "value1", "key2": "value2"},
		"numbers": []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filename := fmt.Sprintf("benchmark_%d", i)
		err := gen.WriteYAML(filename, testData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewGenerator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gen := NewGenerator("/test/path")
		if gen == nil {
			b.Fatal("generator is nil")
		}
	}
}
