package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// EmailService handles sending emails via SMTP or SendGrid Web API
type EmailService struct {
	// SMTP configuration (for MailHog/local dev)
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPassword string
	fromEmail    string
	
	// SendGrid API configuration (for production)
	sendgridAPIKey string
	useSendGridAPI bool
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	sendgridAPIKey := getEnv("SENDGRID_API_KEY", "")
	useAPI := sendgridAPIKey != ""
	
	return &EmailService{
		// SMTP config (for MailHog/local)
		smtpHost:     getEnv("SMTP_HOST", ""),
		smtpPort:     getEnv("SMTP_PORT", "587"),
		smtpUser:     getEnv("SMTP_USER", ""),
		smtpPassword: getEnv("SMTP_PASS", ""),
		fromEmail:    getEnv("SMTP_FROM", getEnv("SMTP_USER", "noreply@website-dummy.com")),
		
		// SendGrid API config
		sendgridAPIKey: sendgridAPIKey,
		useSendGridAPI: useAPI,
	}
}

// IsConfigured checks if email service is properly configured
// Either SMTP (for MailHog) or SendGrid API key must be set
func (es *EmailService) IsConfigured() bool {
	return es.useSendGridAPI || es.smtpHost != ""
}

// GetSMTPHost returns the SMTP host (for logging/debugging)
func (es *EmailService) GetSMTPHost() string {
	if es.useSendGridAPI {
		return "SendGrid API"
	}
	return es.smtpHost
}

// GetSMTPPort returns the SMTP port (for logging/debugging)
func (es *EmailService) GetSMTPPort() string {
	if es.useSendGridAPI {
		return "N/A (API)"
	}
	return es.smtpPort
}

// SendNotificationEmail sends a notification email to the configured recipients
func (es *EmailService) SendNotificationEmail(subject, body string) error {
	if !es.IsConfigured() {
		return fmt.Errorf("email service not configured (neither SENDGRID_API_KEY nor SMTP_HOST is set)")
	}

	notificationEmails := getEnv("NOTIFICATION_EMAILS", "")
	if notificationEmails == "" {
		return fmt.Errorf("NOTIFICATION_EMAILS not configured")
	}

	// Split comma-separated emails
	recipients := strings.Split(notificationEmails, ",")
	for i, email := range recipients {
		recipients[i] = strings.TrimSpace(email)
	}

	// Use SendGrid API if configured, otherwise fall back to SMTP
	if es.useSendGridAPI {
		return es.sendViaSendGridAPI(subject, body, recipients)
	}

	return es.sendViaSMTP(subject, body, recipients)
}

// sendViaSendGridAPI sends email using SendGrid Web API v3
func (es *EmailService) sendViaSendGridAPI(subject, body string, recipients []string) error {
	log.Printf("📧 Sending email via SendGrid Web API to %v", recipients)
	log.Printf("📧 From: %s", es.fromEmail)

	// Build SendGrid API request payload
	// See: https://docs.sendgrid.com/api-reference/mail-send/mail-send
	personalizations := []map[string]interface{}{
		{
			"to": buildEmailObjects(recipients),
		},
	}

	payload := map[string]interface{}{
		"personalizations": personalizations,
		"from": map[string]string{
			"email": es.fromEmail,
		},
		"subject": subject,
		"content": []map[string]string{
			{
				"type":  "text/html",
				"value": body,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Failed to marshal SendGrid payload: %v", err)
		return fmt.Errorf("failed to marshal SendGrid payload: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("❌ Failed to create SendGrid request: %v", err)
		return fmt.Errorf("failed to create SendGrid request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+es.sendgridAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// Send request with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	log.Printf("🌐 Sending request to SendGrid API...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ SendGrid API request failed: %v", err)
		return fmt.Errorf("SendGrid API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body for error details
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("✅ Email sent successfully via SendGrid API to %v", recipients)
		return nil
	}

	// Parse error response
	log.Printf("❌ SendGrid API error (status %d): %s", resp.StatusCode, bodyStr)
	
	// Check for common SendGrid errors
	if strings.Contains(bodyStr, "verified") || strings.Contains(bodyStr, "sender identity") {
		log.Printf("⚠️  SendGrid Error: The sender email address (%s) must be verified in your SendGrid account.", es.fromEmail)
		log.Printf("   Visit SendGrid Dashboard > Settings > Sender Authentication to verify your sender.")
	}
	
	if strings.Contains(bodyStr, "unauthorized") || strings.Contains(bodyStr, "Invalid API key") {
		log.Printf("⚠️  SendGrid Error: Invalid API key. Check your SENDGRID_API_KEY environment variable.")
	}

	return fmt.Errorf("SendGrid API error (status %d): %s", resp.StatusCode, bodyStr)
}

// sendViaSMTP sends email using SMTP (for MailHog/local development)
func (es *EmailService) sendViaSMTP(subject, body string, recipients []string) error {
	log.Printf("📧 Attempting to send email via SMTP %s:%s to %v", es.smtpHost, es.smtpPort, recipients)
	log.Printf("📧 From email: %s", es.fromEmail)

	// Build email message
	message := fmt.Sprintf("From: %s\r\n", es.fromEmail)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(recipients, ", "))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// For MailHog (local dev), we need to handle unencrypted connections
	if es.smtpPort == "1025" || es.smtpHost == "mailhog" {
		log.Printf("📮 Using MailHog mode (no TLS)")
		addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)
		
		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("failed to connect to SMTP server: %w", err)
		}
		defer client.Close()

		if err := client.Mail(es.fromEmail); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		for _, recipient := range recipients {
			if err := client.Rcpt(recipient); err != nil {
				return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
			}
		}

		writer, err := client.Data()
		if err != nil {
			return fmt.Errorf("failed to open data writer: %w", err)
		}

		_, err = writer.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write message: %w", err)
		}

		err = writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		log.Printf("✅ Email sent successfully via MailHog to %v", recipients)
		return client.Quit()
	}

	// For production SMTP servers (not recommended, use SendGrid API instead)
	log.Printf("🌐 Using production SMTP mode with TLS")
	log.Printf("⚠️  Consider using SendGrid API (SENDGRID_API_KEY) instead of SMTP for better reliability")
	
	// SMTP authentication
	var auth smtp.Auth
	if es.smtpUser != "" && es.smtpPassword != "" {
		auth = smtp.PlainAuth("", es.smtpUser, es.smtpPassword, es.smtpHost)
	}

	addr := fmt.Sprintf("%s:%s", es.smtpHost, es.smtpPort)
	
	// Connect with timeout
	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		log.Printf("❌ SMTP connection error: %v", err)
		return fmt.Errorf("failed to connect to SMTP server %s: %w", addr, err)
	}
	
	client, err := smtp.NewClient(conn, es.smtpHost)
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Start TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         es.smtpHost,
	}
	
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	// Send email
	if err := client.Mail(es.fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	for _, recipient := range recipients {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	log.Printf("✅ Email sent successfully via SMTP to %v", recipients)
	return client.Quit()
}

// buildEmailObjects converts email strings to SendGrid email objects
func buildEmailObjects(emails []string) []map[string]string {
	result := make([]map[string]string, len(emails))
	for i, email := range emails {
		result[i] = map[string]string{
			"email": strings.TrimSpace(email),
		}
	}
	return result
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

