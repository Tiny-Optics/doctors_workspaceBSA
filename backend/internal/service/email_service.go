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

// SendStatusChangeNotification sends email to user when their submission status changes
func (s *EmailService) SendStatusChangeNotification(
	smtpConfig models.SMTPConfig,
	userEmail string,
	userName string,
	formName string,
	status string,
	reviewNotes string,
	submissionID string,
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

	// Prepare email subject and body based on status
	var subject string
	var emailBody string

	switch status {
	case "approved":
		subject = "Your Registry Submission Has Been Approved"
		emailBody = fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
					<h2 style="color: #059669;">Submission Approved</h2>
					<p>Dear %s,</p>
					<p>We are pleased to inform you that your submission for <strong>%s</strong> has been approved.</p>
					<div style="background-color: #f0fdf4; border-left: 4px solid #059669; padding: 15px; margin: 20px 0;">
						<p style="margin: 0;"><strong>Submission ID:</strong> %s</p>
						<p style="margin: 5px 0 0 0;"><strong>Status:</strong> <span style="color: #059669;">Approved</span></p>
					</div>
					<p>Thank you for your submission.</p>
					<hr style="border: none; border-top: 1px solid #e5e7eb; margin: 20px 0;">
					<p style="font-size: 12px; color: #6b7280;">
						This is an automated message from the BLOODSA Doctor's Workspace African HOPeR Registry system.
					</p>
				</div>
			</body>
			</html>
		`, userName, formName, submissionID)

	case "rejected":
		subject = "Your Registry Submission Has Been Rejected"
		emailBody = fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
					<h2 style="color: #dc2626;">Submission Rejected</h2>
					<p>Dear %s,</p>
					<p>We regret to inform you that your submission for <strong>%s</strong> has been rejected.</p>
					<div style="background-color: #fef2f2; border-left: 4px solid #dc2626; padding: 15px; margin: 20px 0;">
						<p style="margin: 0;"><strong>Submission ID:</strong> %s</p>
						<p style="margin: 5px 0;"><strong>Status:</strong> <span style="color: #dc2626;">Rejected</span></p>
						<p style="margin: 5px 0 0 0;"><strong>Reason:</strong></p>
						<p style="margin: 5px 0 0 0; padding: 10px; background-color: white; border-radius: 4px;">%s</p>
					</div>
					<p>If you have any questions or would like to discuss this further, please contact the administrator.</p>
					<hr style="border: none; border-top: 1px solid #e5e7eb; margin: 20px 0;">
					<p style="font-size: 12px; color: #6b7280;">
						This is an automated message from the BLOODSA Doctor's Workspace African HOPeR Registry system.
					</p>
				</div>
			</body>
			</html>
		`, userName, formName, submissionID, reviewNotes)

	case "pending":
		subject = "Your Registry Submission Status Has Been Updated"
		emailBody = fmt.Sprintf(`
			<html>
			<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
				<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
					<h2 style="color: #d97706;">Submission Status Update</h2>
					<p>Dear %s,</p>
					<p>The status of your submission for <strong>%s</strong> has been updated to pending review.</p>
					<div style="background-color: #fffbeb; border-left: 4px solid #d97706; padding: 15px; margin: 20px 0;">
						<p style="margin: 0;"><strong>Submission ID:</strong> %s</p>
						<p style="margin: 5px 0 0 0;"><strong>Status:</strong> <span style="color: #d97706;">Pending Review</span></p>
					</div>
					<p>We will notify you once your submission has been reviewed.</p>
					<hr style="border: none; border-top: 1px solid #e5e7eb; margin: 20px 0;">
					<p style="font-size: 12px; color: #6b7280;">
						This is an automated message from the BLOODSA Doctor's Workspace African HOPeR Registry system.
					</p>
				</div>
			</body>
			</html>
		`, userName, formName, submissionID)

	default:
		return fmt.Errorf("unknown status: %s", status)
	}

	// Create email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpConfig.FromEmail)
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	// Send email
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, decryptedPassword)
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("%w: %v", ErrEmailSendFailed, err)
	}

	return nil
}
