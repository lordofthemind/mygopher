package gophersmtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var EmailResultsChan = make(chan EmailResult)

type EmailResult struct {
	Recipient string
	Error     error
}

// EmailRoutineService introduces Go routines to enhance email sending efficiency.
type EmailRoutineService struct {
	smtpHost string
	smtpPort string
	username string
	password string
}

func NewEmailRoutineService(smtpHost, smtpPort, username, password string) GopherSmtpInterface {
	service := &EmailRoutineService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}

	// Start a goroutine to handle results
	go service.processEmailResults()

	return service
}

// SendEmail sends an email to the recipients using a Go routine and reports results via channel.
//
// This function sends a basic email to the specified recipients. It can handle either text or HTML content.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailRoutineService) SendEmail(to []string, subject, body string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	msg := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\";\r\n\r\n%s", subject, mime, body)

	// Go routine to send email asynchronously
	go func() {
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, []byte(msg))
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(to, ", "),
			Error:     err,
		}
	}()

	return nil
}

// SendEmailWithAttachments sends an email with attachments using a Go routine and reports results via channel.
//
// This function sends an email with one or more file attachments to the specified recipients.
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
func (e *EmailRoutineService) SendEmailWithAttachments(to []string, subject, body string, attachmentPaths []string, isHtml bool) error {
	mime := "text/plain"
	if isHtml {
		mime = "text/html"
	}

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/mixed; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	// Create body part of the email
	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Attach files to the email
	for _, path := range attachmentPaths {
		err := e.attachFile(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Go routine to send email asynchronously
	go func() {
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(to, ", "),
			Error:     err,
		}
	}()

	return nil
}

// SendEmailWithHeaders sends an email with custom headers using a Go routine and reports results via channel.
//
// This function sends an email with custom headers. It allows for additional headers like 'Reply-To' or 'From'.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - headers: A map of custom headers to include in the email.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
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
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, []byte(msg))
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(to, ", "),
			Error:     err,
		}
	}()

	return nil
}

// ScheduleEmail schedules an email to be sent at a specific time using a Go routine.
//
// This function schedules an email to be sent at a future time.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - sendAt: The time at which the email should be sent.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailRoutineService) ScheduleEmail(to []string, subject, body string, sendAt time.Time, isHtml bool) error {
	delay := time.Until(sendAt)
	if delay <= 0 {
		return fmt.Errorf("scheduled time is in the past")
	}

	// Schedule the email using a Go routine
	go func() {
		time.Sleep(delay)
		err := e.SendEmail(to, subject, body, isHtml)
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(to, ", "),
			Error:     err,
		}
	}()

	return nil
}

// SendEmailWithCCAndBCC sends an email with CC and BCC recipients using a Go routine.
//
// This function sends an email to the specified recipients, including CC and BCC recipients.
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
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, allRecipients, []byte(headers))
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(allRecipients, ", "),
			Error:     err,
		}
	}()

	return nil
}

// SendBulkEmail sends bulk emails using Go routines for each email.
//
// This function sends multiple emails to the specified list of recipients by spinning up a Go routine for each.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - isHtml: A flag indicating whether the email should be sent in HTML format.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailRoutineService) SendBulkEmail(to []string, subject, body string, isHtml bool) error {
	for _, recipient := range to {
		// Send each email in a Go routine
		go func(recipient string) {
			err := e.SendEmail([]string{recipient}, subject, body, isHtml)
			EmailResultsChan <- EmailResult{
				Recipient: recipient,
				Error:     err,
			}
		}(recipient)
	}
	return nil
}

// SendEmailWithInLineImages sends an email with inline images using a Go routine.
//
// This function sends an email that contains inline images to the specified recipients.
//
// Params:
//   - to: A list of recipient email addresses.
//   - subject: The subject of the email.
//   - body: The content of the email.
//   - imagePaths: A list of file paths for the inline images.
//
// Returns:
//   - error: An error message if the email fails to send.
func (e *EmailRoutineService) SendEmailWithInLineImages(to []string, subject, body string, imagePaths []string) error {
	mime := "text/html"

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	headers := fmt.Sprintf("Subject: %s\r\nMIME-version: 1.0;\r\nContent-Type: multipart/related; boundary=%s\r\n", subject, writer.Boundary())
	buffer.Write([]byte(headers))

	// Create body part of the email
	bodyPart, err := writer.CreatePart(map[string][]string{
		"Content-Type": {mime + "; charset=\"UTF-8\""},
	})
	if err != nil {
		return err
	}
	bodyPart.Write([]byte(body))

	// Attach inline images
	for _, path := range imagePaths {
		err := e.attachInlineImage(writer, path)
		if err != nil {
			return err
		}
	}
	writer.Close()

	// Go routine to send email asynchronously
	go func() {
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(to, ", "),
			Error:     err,
		}
	}()

	return nil
}

// SendEmailWithCCAndBCCAndAttachments sends an email with CC, BCC recipients, and attachments using a Go routine.
// The isHtml flag determines whether it's text or HTML, and the result is reported via a channel.
func (e *EmailRoutineService) SendEmailWithCCAndBCCAndAttachments(to, cc, bcc []string, subject, body string, attachmentPaths []string, isHtml bool) error {
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

	// Go routine to send email asynchronously
	go func() {
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, allRecipients, buffer.Bytes())
		// Send the result to the channel
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(allRecipients, ", "),
			Error:     err,
		}
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
		err := smtp.SendMail(e.smtpHost+":"+e.smtpPort, smtp.PlainAuth("", e.username, e.password, e.smtpHost), e.username, to, buffer.Bytes())
		// Send the result to the channel
		EmailResultsChan <- EmailResult{
			Recipient: strings.Join(to, ", "),
			Error:     err,
		}
	}()

	return nil
}

// Attach file helper function
func (e *EmailRoutineService) attachFile(writer *multipart.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	filePart, err := writer.CreateFormFile("attachment", filepath.Base(path))
	if err != nil {
		return err
	}

	_, err = io.Copy(filePart, file)
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

func (e *EmailRoutineService) processEmailResults() {
	for result := range EmailResultsChan {
		if result.Error != nil {
			log.Printf("Failed to send email to %s: %v\n", result.Recipient, result.Error)
		} else {
			log.Printf("Email sent successfully to %s!\n", result.Recipient)
		}
	}
}
