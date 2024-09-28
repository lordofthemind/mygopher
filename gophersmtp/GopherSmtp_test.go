package gophersmtp

import (
	"testing"
	"time"
)

func TestSendTextEmail(t *testing.T) {
	emailService := NewEmailService("smtp.example.com", "587", "test@example.com", "password")
	err := emailService.SendTextEmail([]string{"recipient@example.com"}, "Test Subject", "Test Body")

	if err != nil {
		t.Errorf("SendTextEmail failed: %v", err)
	}
}

func TestSendHTMLEmail(t *testing.T) {
	emailService := NewEmailService("smtp.example.com", "587", "test@example.com", "password")
	err := emailService.SendHTMLEmail([]string{"recipient@example.com"}, "Test Subject", "<h1>Test Body</h1>")

	if err != nil {
		t.Errorf("SendHTMLEmail failed: %v", err)
	}
}

func TestSendEmailWithAttachment(t *testing.T) {
	emailService := NewEmailService("smtp.example.com", "587", "test@example.com", "password")
	err := emailService.SendEmailWithAttachment([]string{"recipient@example.com"}, "Test Subject", "Test Body", "path/to/attachment.txt")

	if err != nil {
		t.Errorf("SendEmailWithAttachment failed: %v", err)
	}
}

func TestSendEmailWithCCAndBCC(t *testing.T) {
	emailService := NewEmailService("smtp.example.com", "587", "test@example.com", "password")
	err := emailService.SendEmailWithCCAndBCC([]string{"recipient@example.com"}, []string{"cc@example.com"}, []string{"bcc@example.com"}, "Test Subject", "Test Body")

	if err != nil {
		t.Errorf("SendEmailWithCCAndBCC failed: %v", err)
	}
}

func TestScheduleEmail(t *testing.T) {
	emailService := NewEmailService("smtp.example.com", "587", "test@example.com", "password")
	sendAt := time.Now().Add(10 * time.Second)
	err := emailService.ScheduleEmail([]string{"recipient@example.com"}, "Test Subject", "Test Body", sendAt)

	if err != nil {
		t.Errorf("ScheduleEmail failed: %v", err)
	}
}

// Add more tests for other functions following the same structure
