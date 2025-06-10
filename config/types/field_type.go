// Package types provides field type definitions and validation for API test configurations.
package types

// FieldType represents the type of a field in API responses
type FieldType string

const (
	// StringType represents a string field
	StringType FieldType = "string"
	// NumberType represents a numeric field
	NumberType FieldType = "number"
	// BooleanType represents a boolean field
	BooleanType FieldType = "boolean"
	// ArrayType represents an array field
	ArrayType FieldType = "array"
	// ObjectType represents an object field
	ObjectType FieldType = "object"
)

// IsValid checks if the field type is valid
func (ft FieldType) IsValid() bool {
	switch ft {
	case StringType, NumberType, BooleanType, ArrayType, ObjectType:
		return true
	default:
		return false
	}
}

// String returns the string representation of the field type
func (ft FieldType) String() string {
	return string(ft)
}
