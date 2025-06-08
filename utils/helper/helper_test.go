package apihelper_test

import (
	"testing"

	"github.com/RyanTokManMokMTM/api-testing-go/config"
	apierror "github.com/RyanTokManMokMTM/api-testing-go/utils/error"
	apihelper "github.com/RyanTokManMokMTM/api-testing-go/utils/helper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadYamlData(t *testing.T) {
	testTable := []struct {
		name     string
		path     string
		anyError bool
	}{
		{
			"Should load yaml data",
			"../config/etc/api-test/general_work_flow.yaml",
			false,
		},
		{
			"Should return error if file not found",
			"../config/etc/api-test/general_work.yaml",
			true,
		},
	}

	for _, v := range testTable {
		var cfg config.APITest
		err := apihelper.LoadYamlData(v.path, &cfg)
		if v.anyError {
			require.Error(t, err)
		}
	}
}

func TestGetRespFileData(t *testing.T) {
	resp := map[string]any{
		"data": map[string]any{
			"items": []map[string]any{
				{
					"id":     1,
					"key":    "My key",
					"status": "pending",
				},
				{
					"id":     2,
					"key":    "My key 2",
					"status": "completed",
				},
				{
					"id":     3,
					"key":    "My key 3",
					"status": "cancelled",
				},
			},
		},
	}

	testTable := []struct {
		name      string
		fromField string
		resp      any
		throwErr  bool
		err       error
		result    any
	}{
		{
			name:      "should return data from a map",
			fromField: "data.items[0].key",
			resp:      resp,
			throwErr:  false,
			err:       nil,
			result:    "My key",
		},
		{
			name:      "should return list of map",
			fromField: "data.items[0]",
			resp:      resp,
			throwErr:  false,
			err:       nil,
			result: map[string]any{
				"id":     1,
				"key":    "My key",
				"status": "pending",
			},
		},
		{
			name:      "should return error if field not found",
			fromField: "data.item",
			resp:      resp,
			throwErr:  true,
			err:       apierror.ErrDataFieldNotFound,
			result:    nil,
		},
		{
			name:      "should return error if index not found",
			fromField: "data.items[]",
			resp:      resp,
			throwErr:  true,
			err:       apierror.ErrIndexNotFound,
			result:    nil,
		},
		{
			name:      "should return error if index invaild",
			fromField: "data.items[x]",
			resp:      resp,
			throwErr:  true,
			err:       apierror.ErrIndexNotFound,
			result:    nil,
		},
		{
			name:      "should return error if resp type not map",
			fromField: "data.items[x]",
			resp: []string{
				"6666",
			},
			throwErr: true,
			err:      apierror.ErrDataTypeMustBeMap,
			result:   nil,
		},
		{
			name:      "should return error if target type is not a map",
			fromField: "data",
			resp: []string{
				"6666",
			},
			throwErr: true,
			err:      apierror.ErrDataTypeMustBeMap,
			result:   nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			result, err := apihelper.GetRespFieldData(tt.fromField, tt.resp)
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
