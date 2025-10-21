package service

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"backend/internal/models"

	"gopkg.in/gomail.v2"
)

var (
	ErrIncompleteSMTPConfig = errors.New("SMTP configuration is incomplete")
	ErrEmailSendFailed      = errors.New("failed to send email")
)

// EmailService handles sending emails
type EmailService struct {
	encryptionService *EncryptionService
}

// NewEmailService creates a new EmailService
func NewEmailService(encryptionService *EncryptionService) *EmailService {
	return &EmailService{
		encryptionService: encryptionService,
	}
}

// SubmissionNotificationData represents data for submission notification email
type SubmissionNotificationData struct {
	UserName     string
	UserEmail    string
	SubmissionID string
	DropboxLink  string
	FormName     string
	SubmittedAt  string
}

// SendSubmissionNotification sends an email notification for a new submission
func (s *EmailService) SendSubmissionNotification(
	smtpConfig models.SMTPConfig,
	recipients []string,
	data SubmissionNotificationData,
) error {
	// Validate SMTP config
	if !smtpConfig.IsComplete() {
		return ErrIncompleteSMTPConfig
	}

	// Decrypt password
	decryptedPassword, err := s.encryptionService.Decrypt(smtpConfig.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt SMTP password: %w", err)
	}

	// Create email message
	subject := fmt.Sprintf("New Ethics Approval Submission from %s", data.UserName)
	htmlBody, err := s.generateSubmissionEmailHTML(data)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	// Create mailer
	m := gomail.NewMessage()
	m.SetHeader("From", smtpConfig.FromEmail)
	m.SetHeader("To", recipients...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	// Create dialer and send
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, decryptedPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	return nil
}

// generateSubmissionEmailHTML generates the HTML body for submission notification
func (s *EmailService) generateSubmissionEmailHTML(data SubmissionNotificationData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
        }
        .header {
            background-color: #0066cc;
            color: white;
            padding: 20px;
            text-align: center;
            border-radius: 5px 5px 0 0;
        }
        .content {
            background-color: #f9f9f9;
            padding: 30px;
            border: 1px solid #ddd;
            border-radius: 0 0 5px 5px;
        }
        .info-row {
            margin: 15px 0;
        }
        .label {
            font-weight: bold;
            color: #555;
        }
        .button {
            display: inline-block;
            padding: 12px 30px;
            background-color: #0066cc;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            margin-top: 20px;
        }
        .footer {
            margin-top: 30px;
            text-align: center;
            color: #777;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>New Ethics Approval Submission</h1>
        </div>
        <div class="content">
            <p>A new ethics approval has been submitted to the African HOPeR Registry.</p>
            
            <div class="info-row">
                <span class="label">Submitted by:</span> {{.UserName}} ({{.UserEmail}})
            </div>
            
            <div class="info-row">
                <span class="label">Form:</span> {{.FormName}}
            </div>
            
            <div class="info-row">
                <span class="label">Submission ID:</span> {{.SubmissionID}}
            </div>
            
            <div class="info-row">
                <span class="label">Submitted at:</span> {{.SubmittedAt}}
            </div>
            
            <div style="text-align: center;">
                <a href="{{.DropboxLink}}" class="button">View Documents in Dropbox</a>
            </div>
            
            <div class="footer">
                <p>This is an automated notification from the BLOODSA Doctor's Workspace.</p>
            </div>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// SendTestEmail sends a test email to verify SMTP configuration
func (s *EmailService) SendTestEmail(smtpConfig models.SMTPConfig, recipient string) error {
	// Validate SMTP config
	if !smtpConfig.IsComplete() {
		return ErrIncompleteSMTPConfig
	}

	// Decrypt password
	decryptedPassword, err := s.encryptionService.Decrypt(smtpConfig.Password)
	if err != nil {
		return fmt.Errorf("failed to decrypt SMTP password: %w", err)
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpConfig.FromEmail)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", "Test Email - BLOODSA Registry Configuration")
	m.SetBody("text/html", `
		<html>
		<body>
			<h2>Test Email</h2>
			<p>This is a test email from the BLOODSA Doctor's Workspace.</p>
			<p>If you received this email, your SMTP configuration is working correctly.</p>
		</body>
		</html>
	`)

	// Create dialer and send
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, decryptedPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	return nil
}
