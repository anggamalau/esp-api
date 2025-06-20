package services

import (
	"fmt"
	"log"

	"backend/config"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService struct {
	client *sendgrid.Client
}

func NewEmailService() *EmailService {
	client := sendgrid.NewSendClient(config.AppConfig.SendGridAPIKey)
	return &EmailService{
		client: client,
	}
}

// SendPasswordResetEmail sends a new password to the user's email
func (s *EmailService) SendPasswordResetEmail(userEmail, userName, newPassword string) error {
	from := mail.NewEmail(config.AppConfig.SendGridFromName, config.AppConfig.SendGridFromEmail)
	to := mail.NewEmail(userName, userEmail)

	subject := config.AppConfig.ResetPasswordSubject

	// Create email content
	plainTextContent := s.formatPlainTextEmail(userName, newPassword)
	htmlContent := s.formatHTMLEmail(userName, newPassword)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	response, err := s.client.Send(message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	// Check response status
	if response.StatusCode >= 400 {
		log.Printf("SendGrid returned error status: %d, Response: %s", response.StatusCode, response.Body)
		return fmt.Errorf("email service returned status %d", response.StatusCode)
	}

	log.Printf("Password reset email sent successfully to %s", userEmail)
	return nil
}

// formatPlainTextEmail creates the plain text version of the password reset email
func (s *EmailService) formatPlainTextEmail(userName, newPassword string) string {
	return fmt.Sprintf(`Hello %s,

Your password has been reset as requested.

Your new temporary password is: %s

For security reasons, please log in with this password and change it immediately to something you can remember.

Steps to change your password:
1. Log in with the temporary password above
2. Go to your profile settings
3. Update your password

If you did not request this password reset, please contact our support team immediately.

Best regards,
The Support Team`, userName, newPassword)
}

// formatHTMLEmail creates the HTML version of the password reset email
func (s *EmailService) formatHTMLEmail(userName, newPassword string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Password Reset</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #f8f9fa; padding: 20px; text-align: center; border-radius: 5px; }
        .content { padding: 20px 0; }
        .password-box { 
            background-color: #e9ecef; 
            padding: 15px; 
            border-radius: 5px; 
            text-align: center;
            font-family: monospace;
            font-size: 18px;
            font-weight: bold;
            margin: 20px 0;
        }
        .warning { 
            background-color: #fff3cd; 
            border: 1px solid #ffeaa7; 
            padding: 15px; 
            border-radius: 5px; 
            margin: 20px 0;
        }
        .footer { 
            margin-top: 30px; 
            padding-top: 20px; 
            border-top: 1px solid #dee2e6; 
            font-size: 14px; 
            color: #6c757d; 
        }
        ol { padding-left: 20px; }
        li { margin: 10px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset</h1>
        </div>
        
        <div class="content">
            <p>Hello <strong>%s</strong>,</p>
            
            <p>Your password has been reset as requested.</p>
            
            <p>Your new temporary password is:</p>
            <div class="password-box">%s</div>
            
            <div class="warning">
                <strong>⚠️ Important:</strong> For security reasons, please log in with this password and change it immediately to something you can remember.
            </div>
            
            <p><strong>Steps to change your password:</strong></p>
            <ol>
                <li>Log in with the temporary password above</li>
                <li>Go to your profile settings</li>
                <li>Update your password</li>
            </ol>
            
            <p>If you did not request this password reset, please contact our support team immediately.</p>
        </div>
        
        <div class="footer">
            <p>Best regards,<br>The Support Team</p>
            <p><em>This is an automated message. Please do not reply to this email.</em></p>
        </div>
    </div>
</body>
</html>`, userName, newPassword)
}
