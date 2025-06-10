package workflow

import "github.com/RyanTokManMokMTM/api-testing-go/config/types"

// WorkflowType defines the type of workflow
type WorkflowType string
type CheckType string

const (
	// TypeCoupon represents a coupon workflow
	TypeCoupon WorkflowType = "coupon_workflow"
	// TypeSubscription represents a subscription workflow
	TypeSubscription WorkflowType = "subscription_workflow"
	// TypeOrder represents an order workflow
	TypeOrder WorkflowType = "order_workflow"
	// TypeGeneral represents a general workflow
	TypeGeneral WorkflowType = "general_workflow"
)

const (
	// CheckTypeEquals represents equality check
	CheckTypeEquals CheckType = "equals"

	// CheckTypeMatches represents regex match check
	CheckTypeMatches CheckType = "matches"

	// CheckTypePresent represents field presence check
	CheckTypePresent CheckType = "present"

	// CheckTypeNotPresent represents field absence check
	CheckTypeNotPresent CheckType = "not_present"

	// CheckTypeGreaterThan represents greater than comparison
	CheckTypeGreaterThan CheckType = "greater_than"

	// CheckTypeLessThan represents less than comparison
	CheckTypeLessThan CheckType = "less_than"
)

// GetWorkflowType returns the workflow type based on the name
func GetWorkflowType(name string) WorkflowType {
	switch name {
	case string(TypeCoupon):
		return TypeCoupon
	case string(TypeSubscription):
		return TypeSubscription
	case string(TypeOrder):
		return TypeOrder
	case string(TypeGeneral):
		return TypeGeneral
	default:
		return TypeGeneral
	}
}

// TestCase represents a test case
type TestCase struct {
	Name        string
	Description string
	Steps       []TestStep
}

// TestStep represents a test step
type TestStep struct {
	Name           string
	Method         string
	URI            string
	Headers        map[string]string
	Body           map[string]interface{}
	Query          map[string]string
	Variables      []string
	FromResponses  []FromResponse
	ExpectedCode   string
	ExpectedStatus int
	ResponseChecks []ResponseCheck
}

// FromResponse represents a response field reference
type FromResponse struct {
	Step      string
	Name      string
	FromField string
}

// ResponseCheck represents a response check
type ResponseCheck struct {
	CheckType CheckType
	Field     string
	Type      FieldType
	Value     interface{}
	Regex     string
}

// Hook represents a hook configuration
type Hook struct {
	Before HookActions
	After  HookActions
}

// HookActions represents actions to be performed before or after a test
type HookActions struct {
	InitVars  []InitVar
	Workflows []TestStep
}

// InitVar represents a variable initialization
type InitVar struct {
	Name     string
	FromStep string
	Field    string
}

type FieldType string

const (
	// StringType represents a string field
	StringType FieldType = "string"
	// NumberType represents a numeric field
	NumberType FieldType = "number"
	// BooleanType represents a boolean field
	BooleanType FieldType = "boolean"
)

func (ft FieldType) ToConfigTypeField() types.FieldType {
	return types.FieldType(ft)
}
