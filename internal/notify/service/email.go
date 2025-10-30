package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/smtp"
)

type EmailService struct {
	smtpHost     string
	smtpPort     string
	fromEmail    string
	fromPassword string
	log          *zap.SugaredLogger
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	FromEmail    string
	FromPassword string
}

func NewEmailService(cfg EmailConfig, log *zap.SugaredLogger) *EmailService {
	return &EmailService{
		smtpHost:     cfg.SMTPHost,
		smtpPort:     cfg.SMTPPort,
		fromEmail:    cfg.FromEmail,
		fromPassword: cfg.FromPassword,
		log:          log,
	}
}

func (s *EmailService) SendVerificationEmail(ctx context.Context, email, username, token string) error {
	verificationURL := fmt.Sprintf("http://localhost:8081/api/v1/auth/verify/%s", token)

	subject := "Verify Your Email Address"
	body := s.buildVerificationEmail(username, verificationURL)

	return s.sendEmail(email, subject, body)
}

func (s *EmailService) buildVerificationEmail(username, verificationURL string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .button { background-color: #007bff; color: white; padding: 12px 24px; 
                  text-decoration: none; border-radius: 4px; display: inline-block; }
        .footer { margin-top: 20px; font-size: 12px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h2>Verify Your Email Address</h2>
        <p>Hello %s,</p>
        <p>Thank you for registering! Please click the button below to verify your email address:</p>
        <p>
            <a href="%s" class="button">Verify Email</a>
        </p>
        <p>Or copy and paste this link in your browser:</p>
        <p>%s</p>
        <div class="footer">
            <p>If you didn't create an account, please ignore this email.</p>
        </div>
    </div>
</body>
</html>`, username, verificationURL, verificationURL)
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.fromEmail, s.fromPassword, s.smtpHost)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	err := smtp.SendMail(
		s.smtpHost+":"+s.smtpPort,
		auth,
		s.fromEmail,
		[]string{to},
		msg,
	)

	if err != nil {
		s.log.Errorf("Failed to send email to %s: %v", to, err)
		return err
	}

	s.log.Infof("Verification email sent to %s", to)
	return nil
}
