package gophersmtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// EmailRoutineService introduces Go routines to enhance email sending efficiency.
type EmailRoutineService struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

// NewEmailRoutineService creates a new instance of EmailRoutineService with the given SMTP configurations.
func NewEmailRoutineService(smtpHost, smtpPort, username, password string) GopherSmtpInterface {
	return &EmailRoutineService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

// SendEmail sends an email to the recipients using a Go routine.
func (e *EmailRoutineService) SendEmail(to []string, subject, body string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	msg := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", subject, mime, body)

	// Go routine to send email asynchronously
	go func() {
		smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, []byte(msg))
	}()

	return nil
}

// SendEmailWithAttachments sends an email with attachments using a Go routine.
func (e *EmailRoutineService) SendEmailWithAttachments(to []string, subject, body string, attachmentPaths []string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	for _, path := range attachmentPaths {
		err := e.attachFile(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Go routine to send email asynchronously
	go func() {
		smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
	}()

	return nil
}

// SendEmailWithHeaders sends an email with custom headers using a Go routine.
func (e *EmailRoutineService) SendEmailWithHeaders(to []string, subject, body string, headers map[string]string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	headerText := ""
	for key, value := range headers {
		headerText += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	msg := fmt.Sprintf("%sSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", headerText, subject, mime, body)

	// Go routine to send email asynchronously
	go func() {
		smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, []byte(msg))
	}()

	return nil
}

// ScheduleEmail schedules an email to be sent at a specific time using a Go routine.
func (e *EmailRoutineService) ScheduleEmail(to []string, subject, body string, sendAt time.Time, isHtml bool) error {
	delay := time.Until(sendAt)
	if delay <= 0 {
		return fmt.Errorf("scheduled time is in the past")
	}

	// Schedule the email using a Go routine
	go func() {
		time.Sleep(delay)
		e.SendEmail(to, subject, body, isHtml)
	}()

	return nil
}

// SendEmailWithCCAndBCC sends an email with CC and BCC recipients using a Go routine.
func (e *EmailRoutineService) SendEmailWithCCAndBCC(to, cc, bcc []string, subject, body string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	allRecipients := append(to, cc...)
	allRecipients = append(allRecipients, bcc...)

	ccHeader := strings.Join(cc, ",")
	bccHeader := strings.Join(bcc, ",")
	headers := fmt.Sprintf("Subject: %s\r\nCC: %s\r\nBCC: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", subject, ccHeader, bccHeader, mime, body)

	// Go routine to send email asynchronously
	go func() {
		smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, allRecipients, []byte(headers))
	}()

	return nil
}

// SendBulkEmail sends bulk emails using Go routines for each email.
func (e *EmailRoutineService) SendBulkEmail(to []string, subject, body string, isHtml bool) error {
	for _, recipient := range to {
		// Send each email in a Go routine
		go e.SendEmail([]string{recipient}, subject, body, isHtml)
	}
	return nil
}

// SendEmailWithInLineImages sends an email with inline images using a Go routine.
func (e *EmailRoutineService) SendEmailWithInLineImages(to []string, subject, body string, imagePaths []string) error {
	mime := "text/html"

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/related; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	for _, path := range imagePaths {
		err := e.attachInlineImage(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Go routine to send email asynchronously
	go func() {
		smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
	}()

	return nil
}

// SendEmailWithAttachmentsAndInLineImages sends an email with both attachments and inline images using a Go routine.
func (e *EmailRoutineService) SendEmailWithAttachmentsAndInLineImages(to []string, subject, body string, attachmentPaths, imagePaths []string) error {
	mime := "text/html"

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	for _, path := range imagePaths {
		err := e.attachInlineImage(writer, path)
		if err != nil {
			return err
		}
	}

	for _, path := range attachmentPaths {
		err := e.attachFile(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Go routine to send email asynchronously
	go func() {
		smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
	}()

	return nil
}

// Helper function to attach a file to the email.
func (e *EmailRoutineService) attachFile(writer *multipart.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("attachment", filepath.Base(filePath))
	if err != nil {
		return err
	}

	_, err = part.Write([]byte(filePath))
	return err
}

// Helper function to attach an inline image to the email.
func (e *EmailRoutineService) attachInlineImage(writer *multipart.Writer, imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the file's MIME type
	mimeType := "image/" + strings.TrimPrefix(filepath.Ext(imagePath), ".")
	partHeader := make(map[string][]string)
	partHeader["Content-Type"] = []string{mimeType}
	partHeader["Content-Transfer-Encoding"] = []string{"base64"}
	partHeader["Content-Disposition"] = []string{`inline; filename="` + filepath.Base(imagePath) + `";`}
	partHeader["Content-ID"] = []string{`<` + filepath.Base(imagePath) + `>`}

	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return err
	}

	// Read the image and encode it in base64
	imageData := make([]byte, base64.StdEncoding.EncodedLen(len(imagePath)))
	base64.StdEncoding.Encode(imageData, []byte(imagePath))

	_, err = part.Write(imageData)
	return err
}
