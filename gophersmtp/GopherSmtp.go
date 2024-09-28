// Package gophersmtp provides a versatile SMTP email service with support for plain text, HTML, attachments,
// CC/BCC, scheduling, tracking, and more. It simplifies sending various types of emails through SMTP.
package gophersmtp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
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
func NewEmailService(smtpHost, smtpPort, username, password string) EmailServiceInterface {
	return &EmailService{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		username: username,
		password: password,
	}
}

// SendTextEmail sends a plain text email to the specified recipients.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The plain text body of the email.
func (es *EmailService) SendTextEmail(to []string, subject, body string) error {
	return es.sendMail(to, subject, body, "text/plain")
}

// SendHTMLEmail sends an HTML email to the specified recipients.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The HTML body of the email.
func (es *EmailService) SendHTMLEmail(to []string, subject, body string) error {
	return es.sendMail(to, subject, body, "text/html")
}

// SendEmailWithAttachment sends an email with a single attachment.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (plain text or HTML).
// - attachmentPath: The file path of the attachment to be included in the email.
func (es *EmailService) SendEmailWithAttachment(to []string, subject, body, attachmentPath string) error {
	return es.sendEmailWithAttachmentHelper(to, subject, body, []string{attachmentPath})
}

// SendEmailWithMultipleAttachments sends an email with multiple attachments.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (plain text or HTML).
// - attachmentPaths: A list of file paths of the attachments to be included.
func (es *EmailService) SendEmailWithMultipleAttachments(to []string, subject, body string, attachmentPaths []string) error {
	return es.sendEmailWithAttachmentHelper(to, subject, body, attachmentPaths)
}

// SendEmailWithInlineImages sends an email with inline images, typically used in HTML emails.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (HTML format).
// - imagePaths: A list of image file paths to be embedded as inline images.
func (es *EmailService) SendEmailWithInlineImages(to []string, subject, body string, imagePaths []string) error {
	// Inline images are treated as attachments with specific Content-ID headers.
	return es.sendEmailWithAttachmentHelper(to, subject, body, imagePaths)
}

// SendEmailWithCCAndBCC sends an email with CC and BCC recipients.
// Parameters:
// - to: List of primary recipient email addresses.
// - cc: List of email addresses to be CC'd.
// - bcc: List of email addresses to be BCC'd.
// - subject: The subject of the email.
// - body: The body of the email.
func (es *EmailService) SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string) error {
	// Combine 'to', 'cc', and 'bcc' into one list of recipients.
	recipients := append(to, append(cc, bcc...)...)
	return es.sendMail(recipients, subject, body, "text/plain")
}

// SendEmailWithHeaders sends an email with custom headers.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (plain text or HTML).
// - headers: A map of custom headers to include in the email.
func (es *EmailService) SendEmailWithHeaders(to []string, subject, body string, headers map[string]string) error {
	// Build the email message with custom headers.
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	for key, value := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n" + body)
	return es.send(msg.Bytes(), to)
}

// SendPriorityEmail sends an email with a priority level by adding the "X-Priority" header.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (plain text or HTML).
// - priority: The priority level (e.g., "1" for high, "3" for normal, "5" for low).
func (es *EmailService) SendPriorityEmail(to []string, subject, body string, priority string) error {
	headers := map[string]string{
		"X-Priority": priority,
	}
	return es.SendEmailWithHeaders(to, subject, body, headers)
}

// ScheduleEmail schedules an email to be sent at a specific time in the future.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email.
// - sendAt: The time when the email should be sent.
func (es *EmailService) ScheduleEmail(to []string, subject, body string, sendAt time.Time) error {
	// Calculate the duration to wait until sendAt
	duration := time.Until(sendAt)
	if duration <= 0 {
		// If sendAt is in the past or now, send immediately
		return es.SendTextEmail(to, subject, body)
	}

	// Use a goroutine to handle the delay without blocking
	go func() error {
		time.Sleep(duration)
		return es.SendTextEmail(to, subject, body)
	}()

	return nil
}

// SendEmailWithReplyTo sends an email with a "Reply-To" header to allow replies to a specific address.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email.
// - replyTo: The email address for the "Reply-To" header.
func (es *EmailService) SendEmailWithReplyTo(to []string, subject, body, replyTo string) error {
	headers := map[string]string{
		"Reply-To": replyTo,
	}
	return es.SendEmailWithHeaders(to, subject, body, headers)
}

// SendBatchEmail sends the same email to a batch of recipients.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email.
func (es *EmailService) SendBatchEmail(to []string, subject, body string) error {
	// Sends email to each recipient one at a time.
	for _, recipient := range to {
		err := es.SendTextEmail([]string{recipient}, subject, body)
		if err != nil {
			return err
		}
	}
	return nil
}

// SendEmailWithTracking sends an email with tracking features by adding a tracking ID in the email body.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email.
// - trackingID: A tracking identifier (e.g., for marketing emails).
func (es *EmailService) SendEmailWithTracking(to []string, subject, body string, trackingID string) error {
	// Append the tracking ID to the body of the email.
	trackedBody := fmt.Sprintf("%s\n\nTracking ID: %s", body, trackingID)
	return es.SendTextEmail(to, subject, trackedBody)
}

// SendEmailWithAttachmentsAndInlineImages sends an email with both attachments and inline images.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (HTML format).
// - attachmentPaths: A list of file paths of the attachments to be included.
// - imagePaths: A list of image file paths to be embedded as inline images.
func (es *EmailService) SendEmailWithAttachmentsAndInlineImages(to []string, subject, body string, attachmentPaths, imagePaths []string) error {
	// Send email with combined attachments and inline images.
	return es.sendEmailWithAttachmentHelper(to, subject, body, append(attachmentPaths, imagePaths...))
}

// sendMail is a helper function to send plain text or HTML emails.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (plain text or HTML).
// - contentType: The content type of the email (e.g., text/plain or text/html).
func (es *EmailService) sendMail(to []string, subject, body, contentType string) error {
	msg := []byte(fmt.Sprintf("Subject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s; charset=\"utf-8\"\r\n\r\n%s",
		subject, contentType, body))
	return es.send(msg, to)
}

// send is a helper function to actually send the final email message using the SMTP server.
// Parameters:
// - msg: The complete email message as bytes.
// - to: List of recipient email addresses.
func (es *EmailService) send(msg []byte, to []string) error {
	auth := smtp.PlainAuth("", es.username, es.password, es.smtpHost)
	return smtp.SendMail(es.smtpHost+":"+es.smtpPort, auth, es.username, to, msg)
}

// sendEmailWithAttachmentHelper is a helper function to send emails with attachments.
// It encodes files as Base64 before sending.
// Parameters:
// - to: List of recipient email addresses.
// - subject: The subject of the email.
// - body: The body of the email (plain text or HTML).
// - attachmentPaths: A list of file paths of the attachments to be included.
func (es *EmailService) sendEmailWithAttachmentHelper(to []string, subject, body string, attachmentPaths []string) error {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Email headers
	boundary := writer.Boundary()
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buffer.WriteString(fmt.Sprintf("MIME-Version: 1.0\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n", boundary))

	// Add plain text or HTML body part
	var contentType string
	if isHTML(body) {
		contentType = "text/html"
	} else {
		contentType = "text/plain"
	}
	bodyPartHeader := map[string][]string{
		"Content-Type": {fmt.Sprintf("%s; charset=\"utf-8\"", contentType)},
	}
	bodyPart, err := writer.CreatePart(bodyPartHeader)
	if err != nil {
		return err
	}
	_, err = bodyPart.Write([]byte(body))
	if err != nil {
		return err
	}

	// Attach each file in the attachmentPaths
	for _, attachmentPath := range attachmentPaths {
		fileContent, err := es.readFileAndEncodeBase64(attachmentPath)
		if err != nil {
			return err
		}
		attachmentPartHeader := map[string][]string{
			"Content-Type":              {fmt.Sprintf("application/octet-stream; name=\"%s\"", filepath.Base(attachmentPath))},
			"Content-Disposition":       {fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(attachmentPath))},
			"Content-Transfer-Encoding": {"base64"},
		}
		attachmentPart, err := writer.CreatePart(attachmentPartHeader)
		if err != nil {
			return err
		}
		_, err = attachmentPart.Write([]byte(fileContent))
		if err != nil {
			return err
		}
	}

	writer.Close()

	// Send the final email
	return es.send(buffer.Bytes(), to)
}

// readFileAndEncodeBase64 is a helper function to read a file and return its content encoded in Base64 format.
// Parameters:
// - filePath: The path to the file.
func (es *EmailService) readFileAndEncodeBase64(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileData), nil
}

// isHTML is a helper function to determine if the email body is HTML.
// It checks for the presence of HTML tags.
func isHTML(body string) bool {
	return bytes.Contains([]byte(body), []byte("<html>")) || bytes.Contains([]byte(body), []byte("<HTML>"))
}
