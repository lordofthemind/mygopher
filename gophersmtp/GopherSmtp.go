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

// EmailService is responsible for handling email sending with various functionalities
// such as sending plain text, HTML, attachments, and more.
type EmailService struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

// NewEmailService creates a new instance of EmailService with the given SMTP configurations.
// Parameters:
// - smtpHost: The host of the SMTP server.
// - smtpPort: The port of the SMTP server.
// - username: The sender's email address.
// - password: The sender's email account password (used for authentication).
func NewEmailService(smtpHost, smtpPort, username, password string) GopherSmtpInterface {
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

// SendEmail sends an email to the recipients. The isHtml flag determines whether it's text or HTML.
//
// This function composes and sends a basic email to the specified recipients. It can send both plain
// text and HTML emails based on the `isHtml` flag.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmail(to []string, subject, body string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}
	msg := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", subject, mime, body)

	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, []byte(msg))
}

// SendEmailWithAttachments sends an email with attachments. The isHtml flag determines text or HTML format.
//
// This function attaches one or more files to the email and sends it to the recipients. The email can be
// either plain text or HTML based on the `isHtml` flag.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - attachmentPaths: A list of file paths for the attachments.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmailWithAttachments(to []string, subject, body string, attachmentPaths []string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Set headers
	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	// Add body part
	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Attach files
	for _, path := range attachmentPaths {
		err := e.attachFile(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Send the email
	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
}

// SendEmailWithInLineImages sends an email with inline images only.
//
// This function allows embedding images directly into the email content. The email can either be
// plain text or HTML based on the `isHtml` flag.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - inlineImagePaths: A list of file paths for the inline images.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmailWithInLineImages(to []string, subject, body string, inlineImagePaths []string) error {
	mime := "text/html" // If you want to send HTML, else set to "text/plain"

	// Create email body
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Set headers
	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/related; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	// Add body part
	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Attach inline images
	for _, path := range inlineImagePaths {
		err := e.attachInlineImage(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Send the email
	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
}

// SendEmailWithHeaders sends an email with custom headers. The isHtml flag determines text or HTML format.
//
// This function allows setting custom headers such as priority, tracking, and metadata. The email can
// either be plain text or HTML based on the `isHtml` flag.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - headers: A map of custom headers.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmailWithHeaders(to []string, subject, body string, headers map[string]string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	// Compose custom headers
	headerText := ""
	for key, value := range headers {
		headerText += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Complete message
	msg := fmt.Sprintf("%sSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", headerText, subject, mime, body)

	// Send email
	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, []byte(msg))
}

// ScheduleEmail schedules an email to be sent at a specific time. The isHtml flag determines text or HTML format.
//
// This function schedules the email to be sent at a specific time using a goroutine and timer to delay
// execution.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - sendAt: The time when the email should be sent.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the scheduling fails.
func (e *EmailService) ScheduleEmail(to []string, subject, body string, sendAt time.Time, isHtml bool) error {
	delay := time.Until(sendAt)
	if delay <= 0 {
		return fmt.Errorf("scheduled time is in the past")
	}

	go func() {
		time.Sleep(delay)
		e.SendEmail(to, subject, body, isHtml)
	}()

	return nil
}

// SendEmailWithCCAndBCC sends an email with CC and BCC recipients. The isHtml flag determines text or HTML format.
//
// This function sends an email with additional CC and BCC recipients. Both CC and BCC lists are supported
// to allow visibility and hidden recipients.
//
// Params:
//   - to: A list of recipient email addresses.
//   - cc: A list of CC recipient email addresses.
//   - bcc: A list of BCC recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmailWithCCAndBCC(to, cc, bcc []string, subject, body string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	// Merge recipients
	allRecipients := append(to, cc...)
	allRecipients = append(allRecipients, bcc...)

	// Construct headers
	ccHeader := strings.Join(cc, ",")
	bccHeader := strings.Join(bcc, ",")
	headers := fmt.Sprintf("Subject: %s\r\nCC: %s\r\nBCC: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", subject, ccHeader, bccHeader, mime, body)

	// Send email
	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, allRecipients, []byte(headers))
}

// SendBulkEmail sends bulk emails. The isHtml flag determines text or HTML format.
//
// This function is designed for sending the same email to multiple recipients in bulk.
// It can handle plain text and HTML emails based on the `isHtml` flag.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the bulk email fails to send.
func (e *EmailService) SendBulkEmail(to []string, subject, body string, isHtml bool) error {
	for _, recipient := range to {
		if err := e.SendEmail([]string{recipient}, subject, body, isHtml); err != nil {
			return err
		}
	}
	return nil
}

// SendEmailWithCCAndBCCAndAttachments sends an email with CC, BCC recipients, and attachments.
// The isHtml flag determines whether it's text or HTML.
//
// This function sends an email to the specified recipients, including CC, BCC recipients,
// and attaches one or more files to the email.
//
// Params:
//   - to: A list of recipient email addresses.
//   - cc: A list of CC recipient email addresses.
//   - bcc: A list of BCC recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - attachmentPaths: A list of file paths for the attachments.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmailWithCCAndBCCAndAttachments(to, cc, bcc []string, subject, body string, attachmentPaths []string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Set headers
	ccHeader := strings.Join(cc, ",")
	bccHeader := strings.Join(bcc, ",")
	headers := fmt.Sprintf("Subject: %s\r\nCC: %s\r\nBCC: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=%s\r\n", subject, ccHeader, bccHeader, writer.Boundary())
	buffer.Write([]byte(headers))

	// Add body part
	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Attach files
	for _, path := range attachmentPaths {
		err := e.attachFile(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Merge recipients
	allRecipients := append(to, cc...)
	allRecipients = append(allRecipients, bcc...)

	// Send the email
	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, allRecipients, buffer.Bytes())
}

// SendEmailWithAttachmentsAndInLineImages sends an email with both attachments and inline images.
//
// This function combines attachments and inline images into a single email. It supports sending
// both plain text and HTML content, and allows the inclusion of image references in the email body.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - attachmentPaths: A list of file paths for the attachments.
//   - inlineImagePaths: A list of file paths for the inline images.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailService) SendEmailWithAttachmentsAndInLineImages(to []string, subject, body string, attachmentPaths []string, inlineImagePaths []string) error {
	mime := "text/html"

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Set headers
	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	// Add body part
	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Attach inline images
	for _, path := range inlineImagePaths {
		err := e.attachInlineImage(writer, path)
		if err != nil {
			return err
		}
	}

	// Attach other files
	for _, path := range attachmentPaths {
		err := e.attachFile(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Send the email
	return smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
}

// Helper function to attach a file to the email.
func (e *EmailService) attachFile(writer *multipart.Writer, filePath string) error {
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
func (e *EmailService) attachInlineImage(writer *multipart.Writer, imagePath string) error {
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
