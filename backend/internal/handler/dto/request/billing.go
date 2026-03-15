package request

type CheckoutRequest struct {
	PlanID string `json:"plan_id"`
}

type DemoCheckoutRequest struct {
	PlanID     string `json:"plan_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}
