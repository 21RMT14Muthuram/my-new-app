package controller

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
}

// Global email configuration
var EmailCfg = EmailConfig{}

// InitEmailConfig initializes email configuration from environment variables
func InitEmailConfig() {

	   err := godotenv.Load(".env")
    if err != nil {
        fmt.Println(" Warning: Could not load .env file:", err)
    } else {
        fmt.Println(".env file loaded successfully!")
    }
	EmailCfg = EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("FROM_EMAIL"),
	}
	
	// fmt.Print("Env List", EmailCfg)
}


// SendOTPEmail sends OTP to user's email
func SendOTPEmail(toEmail, otpCode string) error {
	// Email subject and body
	subject := "Your OTP Verification Code"
	body := fmt.Sprintf(`
Hello,

Your OTP verification code is: %s

This code will expire in 10 minutes. so "Veegama Sollu da"

If you didn't request this code, please ignore this email.

Best regards,
Un Nanban
	`, otpCode)

	// SMTP authentication
	auth := smtp.PlainAuth("", EmailCfg.SMTPUsername, EmailCfg.SMTPPassword, EmailCfg.SMTPHost)

	// Email message
	message := []byte(
		"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n")

	// Send email
	err := smtp.SendMail(
		EmailCfg.SMTPHost+":"+EmailCfg.SMTPPort,
		auth,
		EmailCfg.FromEmail,
		[]string{toEmail},
		message,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// IsEmailConfigured checks if email configuration is set
func IsEmailConfigured() bool {
	return EmailCfg.SMTPUsername != "" && EmailCfg.SMTPPassword != ""
}