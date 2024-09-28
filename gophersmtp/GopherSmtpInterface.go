package gophersmtp

import (
	"time"
)

type GopherSmtpInterface interface {
	// SendEmail sends an email to the recipients. The isHtml flag determines whether it's text or HTML.
	SendEmail(to []string, subject, body string, isHtml bool) error

	// SendEmailWithAttachments sends an email with attachments. The isHtml flag determines text or HTML format.
	SendEmailWithAttachments(to []string, subject, body string, attachmentPaths []string, isHtml bool) error

	// SendEmailWithHeaders sends an email with custom headers. The isHtml flag determines text or HTML format.
	SendEmailWithHeaders(to []string, subject, body string, headers map[string]string, isHtml bool) error

	// ScheduleEmail schedules an email to be sent at a specific time. The isHtml flag determines text or HTML format.
	ScheduleEmail(to []string, subject, body string, sendAt time.Time, isHtml bool) error

	// SendEmailWithCCAndBCC sends an email with CC and BCC recipients. The isHtml flag determines text or HTML format.
	SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string, isHtml bool) error

	// SendBulkEmail sends bulk emails. The isHtml flag determines text or HTML format.
	SendBulkEmail(to []string, subject, body string, isHtml bool) error

	// SendEmailWithInLineImages sends an email with inline images.
	// Only applicable for HTML emails.
	SendEmailWithInLineImages(to []string, subject, body string, imagePaths []string) error

	// SendEmailWithAttachmentsAndInLineImages sends an email with both attachments and inline images.
	// Only applicable for HTML emails.
	SendEmailWithAttachmentsAndInLineImages(to []string, subject, body string, attachmentPaths, imagePaths []string) error
}
