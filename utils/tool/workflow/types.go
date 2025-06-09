package workflow

// WorkflowType defines the type of workflow
type WorkflowType string

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
