package apihelper

import (
	"testing"

	apierror "github.com/RyanTokManMokMTM/api-testing-go/utils/error"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromMapOrList(t *testing.T) {
	tests := []struct {
		name     string
		layerStr string
		data     map[string]any
		throwErr bool
		err      error
		result   any
	}{
		{
			name:     "Should return data from a map",
			layerStr: "data",
			data: map[string]any{
				"data": map[string]any{
					"id": 1,
				},
			},
			throwErr: false,
			err:      nil,
			result: map[string]any{
				"id": 1,
			},
		},
		{
			name:     "Should return data from a list",
			layerStr: "data[0]",
			data: map[string]any{
				"data": []any{
					map[string]any{
						"id": 1,
					},
				},
			},
			throwErr: false,
			err:      nil,
			result: map[string]any{
				"id": 1,
			},
		},
		{
			name:     "Should return error if field not found",
			layerStr: "data.item",
			data: map[string]any{
				"data": map[string]any{},
			},
			throwErr: true,
			err:      apierror.ErrDataFieldNotFound,
			result:   nil,
		},
		{
			name:     "Should return error if index out of range",
			layerStr: "data[1]",
			data: map[string]any{
				"data": []any{
					map[string]any{
						"id": 1,
					},
				},
			},
			throwErr: true,
			err:      apierror.ErrIndexOutOfRange,
			result:   nil,
		},
		{
			name:     "Should return error if index not found",
			layerStr: "data[",
			data: map[string]any{
				"data": []any{
					map[string]any{
						"id": 1,
					},
				},
			},
			throwErr: true,
			err:      apierror.ErrIndexNotFound,
			result:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fromMapOrList(tt.layerStr, tt.data)
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

func TestCheckTypeAndGetValueByIndex(t *testing.T) {
	testTable := []struct {
		name     string
		index    int
		data     any
		throwErr bool
		err      error
		result   any
	}{
		{
			"Should get data from a map list with index 0",
			0,
			[]map[string]any{
				{
					"id": 1,
				},
				{
					"id": 2,
				},
			},
			false,
			nil,
			map[string]any{
				"id": 1,
			},
		},
		{
			"Should get data from a int list with index 0",
			0,
			[]int{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			1,
		},
		{
			"Should get data from a int8 list with index 0",
			0,
			[]int8{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			int8(1),
		},
		{
			"Should get data from a int16 list with index 0",
			0,
			[]int16{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			int16(1),
		},
		{
			"Should get data from a int32 list with index 0",
			0,
			[]int32{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			int32(1),
		},
		{
			"Should get data from a int64 list with index 0",
			0,
			[]int64{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			int64(1),
		},

		{
			"Should get data from a uint list with index 0",
			0,
			[]uint{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			uint(1),
		},
		{
			"Should get data from a uint8 list with index 0",
			0,
			[]uint8{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			uint8(1),
		},
		{
			"Should get data from a uint16 list with index 0",
			0,
			[]uint16{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			uint16(1),
		},
		{
			"Should get data from a uint32 list with index 0",
			0,
			[]uint32{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			uint32(1),
		},
		{
			"Should get data from a uint32 list with index 0",
			0,
			[]uint32{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			uint32(1),
		},
		{
			"Should get data from a uint64 list with index 0",
			0,
			[]uint64{
				1, 2, 3, 4, 5,
			},
			false,
			nil,
			uint64(1),
		},
		{
			"Should get data from a string list with index 0",
			0,
			[]string{
				"hello,world",
			},
			false,
			nil,
			"hello,world",
		},
		{
			"Should get data from a bool list with index 0",
			0,
			[]bool{
				true, false,
			},
			false,
			nil,
			true,
		},
		{
			"Should get data from a float32 list with index 0",
			0,
			[]float32{
				3.21, 4.56,
			},
			false,
			nil,
			float32(3.21),
		},
		{
			"Should get data from a float64 list with index 0",
			0,
			[]float64{
				3.2112223, 4.5623,
			},
			false,
			nil,
			float64(3.2112223),
		},
		{
			"Should get data from a complex64 list with index 0",
			0,
			[]complex64{
				3.21, 4.56,
			},
			false,
			nil,
			complex64(3.21),
		},
		{
			"Should get data from a complex128 list with index 0",
			0,
			[]complex128{
				3.141592653589793, 4.56,
			},
			false,
			nil,
			complex128(3.141592653589793),
		},
		{
			"Should get data from a any list with index 0",
			0,
			[]any{
				"hello",
			},
			false,
			nil,
			"hello",
		},
		{
			"Should return error if type not supported",
			0,
			[]uintptr{
				1,
			},
			true,
			apierror.ErrTypeNotSupported,
			nil,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			result, err := checkTypeAndGetValueByIndex(tt.data, tt.index)
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

func TestGetValue(t *testing.T) {
	testTable := []struct {
		name     string
		index    int
		list     []int
		throwErr bool
		err      error
		result   any
	}{
		{
			name:     "Should get data from a list with index 0",
			index:    2,
			list:     []int{1, 2, 3, 4, 5},
			throwErr: false,
			err:      nil,
			result:   3,
		},
		{
			name:     "Should return error if index out of range",
			index:    100,
			list:     []int{1, 2, 3, 4, 5},
			throwErr: true,
			err:      apierror.ErrIndexOutOfRange,
			result:   0,
		},
		{
			name:     "Should return error if index not found",
			index:    -1,
			list:     []int{1, 2, 3, 4, 5},
			throwErr: true,
			err:      apierror.ErrIndexOutOfRange,
			result:   0,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getValue(tt.index, tt.list)
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
