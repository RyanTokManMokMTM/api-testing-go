package apitestrunner

import (
	"testing"
	"text/template"

	apierror "github.com/RyanTokManMokMTM/api-testing-go/utils/error"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeader(t *testing.T) {
	runner := NewAPITestRuuner(nil)
	testTables := []struct {
		name      string
		headerStr string
		variables map[string]any
		isAnyErr  bool
		result    map[string]string
	}{
		{
			"Should return result if parsing success",
			`{"content-type": "application/json","custom_value":"{{.custom_value}}"}`,
			map[string]any{
				"custom_value": "test",
			},
			false,
			map[string]string{
				"content-type": "application/json",
				"custom_value": "test",
			},
		},
		{
			"Should return empty map if no header",
			"",
			nil,
			false,
			map[string]string{},
		},
		{
			"Should return errir if parsing failed",
			`{"content-type": "application/json","custom-value":"{{.custom-value}}"}`,
			map[string]any{
				"custom-value": "test",
			},
			true,
			nil,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runner.parseHeader(tt.variables, tt.headerStr)
			if tt.isAnyErr {
				require.Error(t, err)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestParseQuery(t *testing.T) {
	runner := NewAPITestRuuner(nil)
	testTables := []struct {
		name      string
		queryStr  string
		variables map[string]any
		isAnyErr  bool
		result    map[string]string
	}{
		{
			"Should return result if parsing success",
			`{"type":"{{.type}}","ids":"{{.ids}}","limit":"{{.limit}}","offset":"{{.offset}}"}`,
			map[string]any{
				"type":   "plan",
				"ids":    "ID1,ID2",
				"limit":  15,
				"offset": 0,
			},
			false,
			map[string]string{
				"type":   "plan",
				"ids":    "ID1,ID2",
				"limit":  "15",
				"offset": "0",
			},
		},
		{
			"Should return empty map if no query string",
			"",
			nil,
			false,
			map[string]string{},
		},
		{
			"Should return error if parsing failed",
			`{"custom-type":"{{.custom-type}}"`,
			map[string]any{
				"custom-value": "planC",
			},
			true,
			nil,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runner.parseQuery(tt.variables, tt.queryStr)
			if tt.isAnyErr {
				require.Error(t, err)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}
func TestParseBody(t *testing.T) {
	runner := NewAPITestRuuner(nil)
	testTables := []struct {
		name      string
		queryStr  string
		variables map[string]any
		isAnyErr  bool
		result    map[string]string
	}{
		{
			"Should return result if parsing success",
			`{"type":"{{.type}}","ids":"{{.ids}}","limit":"{{.limit}}","offset":"{{.offset}}"}`,
			map[string]any{
				"type":   "plan",
				"ids":    "ID1,ID2",
				"limit":  15,
				"offset": 0,
			},
			false,
			map[string]string{
				"type":   "plan",
				"ids":    "ID1,ID2",
				"limit":  "15",
				"offset": "0",
			},
		},
		{
			"Should return empty map if no query string",
			"",
			nil,
			false,
			map[string]string{},
		},
		{
			"Should return error if parsing failed",
			`{"custom-type":"{{.custom-type}}"`,
			map[string]any{
				"custom-value": "planC",
			},
			true,
			nil,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runner.parseQuery(tt.variables, tt.queryStr)
			if tt.isAnyErr {
				require.Error(t, err)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestParseURI(t *testing.T) {
	runner := NewAPITestRuuner(nil)
	testTables := []struct {
		name      string
		uri       string
		variables map[string]any
		isAnyErr  bool
		result    string
	}{
		{
			"Should return result if parsing success",
			"/api/merchants/{{.mid}}/subscriptions",
			map[string]any{
				"mid": "SDFVGBCNHYIW3456YHGHJI90",
			},
			false,
			"/api/merchants/SDFVGBCNHYIW3456YHGHJI90/subscriptions",
		},
		{
			"Should return empty map if no query string",
			"",
			nil,
			false,
			"",
		},
		{
			"Should return error if parsing failed",
			`/api/merchants/{{.custom-mid}}/subscriptions`,
			map[string]any{
				"custom-mid": "SDFVGBCNHYIW3456YHGHJI90",
			},
			true,
			"",
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runner.parseURI(tt.variables, tt.uri)
			if tt.isAnyErr {
				require.Error(t, err)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestParseURIQuery(t *testing.T) {
	runner := NewAPITestRuuner(nil)
	testTables := []struct {
		name      string
		template  *template.Template
		input     string
		variables map[string]any
		isAnyErr  bool
		result    string
	}{
		{
			"should return string if parsing success",
			uriTemp,
			"/api/merchants/{{.mid}}/subscriptions",
			map[string]any{
				"mid": "SDFVGBCNHYIW3456YHGHJI90",
			},
			false,
			"/api/merchants/SDFVGBCNHYIW3456YHGHJI90/subscriptions",
		},
		{
			"should return err if parse string failed",
			uriTemp,
			"/api/merchants/{{.mid}/subscriptions",
			map[string]any{
				"mid": "SDFVGBCNHYIW3456YHGHJI90",
			},
			true,
			"",
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runner.parseTemplate(tt.template, tt.input, tt.variables)
			if tt.isAnyErr {
				require.Error(t, err)
			}
			assert.Equal(t, tt.result, result)
		})
	}
}

func TestGetValueFromResponse(t *testing.T) {
	runner := NewAPITestRuuner(nil)
	resp := map[string]any{
		"data": map[string]any{
			"items": []map[string]any{
				{
					"id": 1,
				},
			},
		},
	}

	respCache := map[string]any{
		"test_step": resp,
	}

	testTables := []struct {
		name      string
		step      string
		varName   string
		formField string
		throwErr  bool
		result    any
		err       error
	}{
		{
			"Should return value if success",
			"test_step",
			"id",
			"data.items[0].id",
			false,
			1,
			nil,
		},
		{
			"Should return error if step not found",
			"test_step_demo",
			"id",
			"data.items[0].id",
			true,
			nil,
			apierror.ErrStepRespNotFound,
		},
		{
			"Should return error if data is nil",
			"test_step",
			"id",
			"data.items[0].temp",
			true,
			nil,
			apierror.ErrDataFieldNotFound,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			result, err := runner.getValueFromResponse(
				&respCache,
				tt.varName,
				tt.step,
				tt.formField)

			assert.Equal(t, tt.result, result)
			if tt.throwErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.result, result)
			}
		})
	}
}

func TestIsExistVars(t *testing.T) {
	runner := NewAPITestRuuner(nil)

	testTables := []struct {
		name     string
		vars     map[string]any
		str      string
		isAnyErr bool
		err      error
	}{
		{
			"Should not return error if var exist",
			map[string]any{
				"mid": "SDFVGBCNHYIW3456YHGHJI90",
			},
			"{{.mid}} ids",
			false,
			nil,
		},
		{
			"Should not return error if var not exist",
			map[string]any{
				"mid": "SDFVGBCNHYIW3456YHGHJI90",
			},
			"{{.mid}} ids {{.custom-mid}}",
			true,
			apierror.ErrVarsNotFound,
		},
	}

	for _, tt := range testTables {
		t.Run(tt.name, func(t *testing.T) {
			err := runner.isExistVars(
				tt.str,
				tt.vars)
			if tt.isAnyErr {
				require.Error(t, err)
			}
		})
	}
}
