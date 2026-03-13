package resend

import "fmt"

// renderTemplate generates subject and HTML body for known templates.
func renderTemplate(template string, data map[string]string) (subject, body string) {
	switch template {
	case "welcome":
		return WelcomeEmail(data["name"])
	case "verification":
		return VerificationEmail(data["name"], data["link"])
	case "password_reset":
		return PasswordResetEmail(data["name"], data["link"])
	default:
		return "CleanSaaS Notification", fmt.Sprintf("<p>%s</p>", data["message"])
	}
}

// WelcomeEmail returns subject and HTML body for welcome emails.
func WelcomeEmail(name string) (subject, body string) {
	subject = "Welcome to CleanSaaS!"
	body = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px; color: #1a1a1a;">
  <div style="text-align: center; margin-bottom: 32px;">
    <h1 style="color: #e11d48; font-size: 24px; margin: 0;">CleanSaaS</h1>
  </div>
  <h2 style="font-size: 20px; margin-bottom: 16px;">Welcome, %s!</h2>
  <p style="font-size: 16px; line-height: 1.6; color: #4a4a4a;">
    Thank you for joining CleanSaaS. We're excited to have you on board.
  </p>
  <p style="font-size: 16px; line-height: 1.6; color: #4a4a4a;">
    Get started by exploring your dashboard and setting up your profile.
  </p>
  <div style="text-align: center; margin: 32px 0;">
    <a href="#" style="background-color: #e11d48; color: white; padding: 12px 32px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Go to Dashboard</a>
  </div>
  <p style="font-size: 14px; color: #9a9a9a; margin-top: 40px; border-top: 1px solid #e5e5e5; padding-top: 20px;">
    This email was sent by CleanSaaS. If you didn't create an account, you can safely ignore this email.
  </p>
</body>
</html>`, name)
	return
}

// VerificationEmail returns subject and HTML body for email verification.
func VerificationEmail(name, link string) (subject, body string) {
	subject = "Verify your email address"
	body = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px; color: #1a1a1a;">
  <div style="text-align: center; margin-bottom: 32px;">
    <h1 style="color: #e11d48; font-size: 24px; margin: 0;">CleanSaaS</h1>
  </div>
  <h2 style="font-size: 20px; margin-bottom: 16px;">Verify your email</h2>
  <p style="font-size: 16px; line-height: 1.6; color: #4a4a4a;">
    Hi %s, please click the button below to verify your email address.
  </p>
  <div style="text-align: center; margin: 32px 0;">
    <a href="%s" style="background-color: #e11d48; color: white; padding: 12px 32px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Verify Email</a>
  </div>
  <p style="font-size: 14px; color: #4a4a4a;">
    Or copy this link into your browser:<br>
    <a href="%s" style="color: #e11d48; word-break: break-all;">%s</a>
  </p>
  <p style="font-size: 14px; color: #9a9a9a;">This link expires in 24 hours.</p>
  <p style="font-size: 14px; color: #9a9a9a; margin-top: 40px; border-top: 1px solid #e5e5e5; padding-top: 20px;">
    If you didn't create an account, you can safely ignore this email.
  </p>
</body>
</html>`, name, link, link, link)
	return
}

// PasswordResetEmail returns subject and HTML body for password reset.
func PasswordResetEmail(name, link string) (subject, body string) {
	subject = "Reset your password"
	body = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px; color: #1a1a1a;">
  <div style="text-align: center; margin-bottom: 32px;">
    <h1 style="color: #e11d48; font-size: 24px; margin: 0;">CleanSaaS</h1>
  </div>
  <h2 style="font-size: 20px; margin-bottom: 16px;">Reset your password</h2>
  <p style="font-size: 16px; line-height: 1.6; color: #4a4a4a;">
    Hi %s, we received a request to reset your password. Click the button below to choose a new one.
  </p>
  <div style="text-align: center; margin: 32px 0;">
    <a href="%s" style="background-color: #e11d48; color: white; padding: 12px 32px; border-radius: 8px; text-decoration: none; font-weight: 600; font-size: 16px;">Reset Password</a>
  </div>
  <p style="font-size: 14px; color: #4a4a4a;">
    Or copy this link into your browser:<br>
    <a href="%s" style="color: #e11d48; word-break: break-all;">%s</a>
  </p>
  <p style="font-size: 14px; color: #9a9a9a;">This link expires in 1 hour. If you didn't request a reset, ignore this email.</p>
  <p style="font-size: 14px; color: #9a9a9a; margin-top: 40px; border-top: 1px solid #e5e5e5; padding-top: 20px;">
    This email was sent by CleanSaaS.
  </p>
</body>
</html>`, name, link, link, link)
	return
}
