package mailer

import (
	"bytes"
	"fmt"
	"net/smtp"
	"sync"
)

// Mailer defines the interface for sending emails
type Mailer interface {
	SendEmail(to, subject, body string) error
}

// SMTPMailer implements the Mailer interface using standard net/smtp
type SMTPMailer struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// NewSMTPMailer creates a new SMTPMailer
func NewSMTPMailer(host string, port int, user, pass, from string) *SMTPMailer {
	return &SMTPMailer{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
		From: from,
	}
}

// SendEmail sends an email using SMTP
func (m *SMTPMailer) SendEmail(to, subject, body string) error {
	// If no host is configured, we can just log or skip.
	// This acts as a safeguard if SMTP isn't configured.
	if m.Host == "" {
		fmt.Printf("Mock email sent to: %s, Subject: %s\n", to, subject)
		return nil
	}

	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)
	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	// Format headers and body
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", to))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buf.WriteString("MIME-version: 1.0;\r\n")
	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	buf.WriteString(body)

	err := smtp.SendMail(addr, auth, m.From, []string{to}, buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

// MockMailer is a mock implementation of Mailer for testing
type MockMailer struct {
	mu         sync.Mutex
	SentEmails []MockEmail
}

// MockEmail represents a captured email in the mock mailer
type MockEmail struct {
	To      string
	Subject string
	Body    string
}

// NewMockMailer creates a new MockMailer
func NewMockMailer() *MockMailer {
	return &MockMailer{
		SentEmails: []MockEmail{},
	}
}

// SendEmail records the email in the mock mailer
func (m *MockMailer) SendEmail(to, subject, body string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SentEmails = append(m.SentEmails, MockEmail{
		To:      to,
		Subject: subject,
		Body:    body,
	})
	return nil
}
