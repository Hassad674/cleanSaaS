package notification

type Template string

const (
	TemplateWelcome       Template = "welcome"
	TemplateVerifyEmail   Template = "verify_email"
	TemplateResetPassword Template = "reset_password"
	TemplateInvoice       Template = "invoice"
	TemplateSubCanceled   Template = "subscription_canceled"
)
