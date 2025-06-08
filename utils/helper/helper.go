package apihelper

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/RyanTokManMokMTM/api-testing-go/config/types"
	apierr "github.com/RyanTokManMokMTM/api-testing-go/utils/error"
)

func LoadYamlData[T any](path string, data *T) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if umarshalErr := yaml.Unmarshal(bytes, data); umarshalErr != nil {
		return umarshalErr
	}

	return nil
}

func GetRespFieldData(fromField string, resp any) (any, error) {
	// MAKR: split the field into layers by .
	// Example: data.data[0].id ->[data,data[0],id]
	layers := strings.Split(fromField, ".")
	layersLen := len(layers)

	currentData := resp
	for _, layer := range layers[:layersLen-1] {
		// convert current data to a key-value map
		// It shoule be a key-value pattern.
		data, ok := currentData.(map[string]any)
		if !ok {
			return nil, apierr.ErrDataTypeMustBeMap
		}

		// get data from either a map or a list
		result, err := fromMapOrList(layer, data)
		if err != nil {
			return nil, err
		}
		currentData = result
	}

	// after processing all layers except the last layer
	// example: data.data[0].id , the 'id' will be the result it want but haven't process yet.

	// Target field: data.data[0].id : 'id'
	targetField := layers[layersLen-1]

	// to check the field result
	// but current data can be map or a list
	// example: data.data[0], the target field is data[0]
	// example: data.data[0].id, the target field is id,which is a key of map

	// conver the result to a map
	data, ok := currentData.(map[string]any)
	if !ok {
		return nil, apierr.ErrDataTypeMustBeMap
	}

	// get data from either a map or a list
	result, err := fromMapOrList(targetField, data)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//nolint:nestif // if condition is short
func fromMapOrList(
	layer string,
	data map[string]any,
) (any, error) {
	var currentData any
	if strings.Contains(layer, "[") {
		// get the field name before '['
		field := strings.Split(layer, "[")[0]
		currentData = data[field] // should a list of map or list an any type.

		// To get index by regex
		re := regexp.MustCompile(`\[(\d+)\]`)
		matches := re.FindStringSubmatch(layer)

		//nolint:mnd // magic number
		if len(matches) < 2 {
			return nil, apierr.ErrIndexNotFound
		}

		index, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, apierr.ErrIndexInvaild
		}

		currentData, err = checkTypeAndGetValueByIndex(currentData, index)
		if err != nil {
			return nil, err
		}
	} else {
		// if key not found return error
		v, ok := data[layer]
		if !ok {
			return nil, fmt.Errorf("data field '%s' not found: %w", layer, apierr.ErrDataFieldNotFound)
		}
		currentData = v // is a map, get value by layer name
	}
	return currentData, nil
}

func checkTypeAndGetValueByIndex(
	data any, index int,
) (any, error) {
	// Check current type.
	var err error
	currentData := data
	switch list := currentData.(type) {
	case []map[string]any:
		currentData, err = getValue(index, list)
	case []int:
		currentData, err = getValue(index, list)
	case []int8:
		currentData, err = getValue(index, list)
	case []int16:
		currentData, err = getValue(index, list)
	case []int32:
		currentData, err = getValue(index, list)
	case []int64:
		currentData, err = getValue(index, list)
	case []uint:
		currentData, err = getValue(index, list)
	case []uint8:
		currentData, err = getValue(index, list)
	case []uint16:
		currentData, err = getValue(index, list)
	case []uint32:
		currentData, err = getValue(index, list)
	case []uint64:
		currentData, err = getValue(index, list)
	case []string:
		currentData, err = getValue(index, list)
	case []bool:
		currentData, err = getValue(index, list)
	case []float32:
		currentData, err = getValue(index, list)
	case []float64:
		currentData, err = getValue(index, list)
	case []complex64:
		currentData, err = getValue(index, list)
	case []complex128:
		currentData, err = getValue(index, list)
	case []any:
		currentData, err = getValue(index, list)
	default:
		err = apierr.ErrTypeNotSupported
	}
	if err != nil {
		return nil, err
	}
	return currentData, nil
}

func getValue[T any](index int, list []T) (T, error) {
	// list items is a list of any
	if index < 0 || index >= len(list) {
		return *new(T), apierr.ErrIndexOutOfRange
	}
	return list[index], nil
}

func ConvertStrToType(str string, convertTo types.FieldType) (any, error) {
	switch convertTo {
	case types.NumberType: // convert to int64
		return strconv.ParseInt(str, 10, 64)
	case types.StringType:
		return str, nil
	case types.BooleanType:
		return strconv.ParseBool(str)
	case types.ArrayType:
		// For array type, we'll return the string as is and let the caller handle the parsing
		return str, nil
	case types.ObjectType:
		// For object type, we'll return the string as is and let the caller handle the parsing
		return str, nil
	default:
		return str, nil
	}
}
