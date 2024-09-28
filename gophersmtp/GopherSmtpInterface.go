// Package gophersmtp provides functionalities for sending emails via an SMTP server.
package gophersmtp

import "time"

// EmailServiceInterface defines all the necessary methods for sending various types of emails.
type EmailServiceInterface interface {
	// SendTextEmail sends a plain text email to the recipients.
	SendTextEmail(to []string, subject, body string) error

	// SendHTMLEmail sends an HTML email to the recipients.
	SendHTMLEmail(to []string, subject, body string) error

	// SendEmailWithAttachment sends an email with a single attachment.
	SendEmailWithAttachment(to []string, subject, body, attachmentPath string) error

	// SendEmailWithCCAndBCC sends an email with CC and BCC recipients.
	SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string) error

	// SendEmailWithMultipleAttachments sends an email with multiple attachments.
	SendEmailWithMultipleAttachments(to []string, subject, body string, attachmentPaths []string) error

	// SendEmailWithInlineImages sends an email with inline images (e.g., for HTML emails).
	SendEmailWithInlineImages(to []string, subject, body string, imagePaths []string) error

	// SendEmailWithHeaders allows sending emails with custom headers.
	SendEmailWithHeaders(to []string, subject, body string, headers map[string]string) error

	// SendPriorityEmail allows sending emails with a specific priority level.
	SendPriorityEmail(to []string, subject, body string, priority string) error

	// ScheduleEmail schedules an email to be sent at a specific time.
	ScheduleEmail(to []string, subject, body string, sendAt time.Time) error

	// SendEmailWithReplyTo sends an email with a "Reply-To" address.
	SendEmailWithReplyTo(to []string, subject, body, replyTo string) error

	// SendBatchEmail sends an email to a batch of recipients.
	SendBatchEmail(to []string, subject, body string) error

	// SendEmailWithTracking sends an email with tracking features enabled (e.g., for marketing purposes).
	SendEmailWithTracking(to []string, subject, body string, trackingID string) error

	// SendEmailWithAttachmentsAndInlineImages sends an email with both attachments and inline images.
	SendEmailWithAttachmentsAndInlineImages(to []string, subject, body string, attachmentPaths, imagePaths []string) error
}
