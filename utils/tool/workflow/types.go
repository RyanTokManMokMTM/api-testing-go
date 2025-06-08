package workflow

// WorkflowType 定義工作流類型
type WorkflowType string

const (
	// TypeCoupon 優惠券工作流
	TypeCoupon WorkflowType = "coupon_workflow"
	// TypeSubscription 訂閱工作流
	TypeSubscription WorkflowType = "subscription_workflow"
	// TypeOrder 訂單工作流
	TypeOrder WorkflowType = "order_workflow"
	// TypeGeneral 通用工作流
	TypeGeneral WorkflowType = "general_workflow"
)

// GetWorkflowType 根據名稱獲取工作流類型
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
