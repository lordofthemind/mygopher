// Package gophersmtp provides a versatile SMTP email service with support for plain text, HTML, attachments,
// CC/BCC, scheduling, and more. It simplifies sending various types of emails through SMTP.
package gophersmtp

import (
	"time"
)

// GopherSmtpInterface defines the necessary methods for sending various types of emails.
//
// These functions allow for sending plain text and HTML emails with additional options
// like attachments, headers, CC/BCC, scheduling, and bulk email support.
type GopherSmtpInterface interface {
	// SendEmail sends an email to the recipients. The isHtml flag determines whether it's text or HTML.
	//
	// This method sends a basic email either in plain text or HTML format depending on the value
	// of the `isHtml` flag. It supports sending emails to multiple recipients.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email.
	//	isHtml - A flag indicating whether the email should be sent in HTML format.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendEmail(to []string, subject, body string, isHtml bool) error

	// SendEmailWithAttachments sends an email with attachments. The isHtml flag determines text or HTML format.
	//
	// This method sends an email with one or more file attachments, with support for either
	// plain text or HTML content. The attachments are specified by their file paths.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email.
	//	attachmentPaths - A list of file paths for the attachments.
	//	isHtml - A flag indicating whether the email should be sent in HTML format.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendEmailWithAttachments(to []string, subject, body string, attachmentPaths []string, isHtml bool) error

	// SendEmailWithHeaders sends an email with custom headers. The isHtml flag determines text or HTML format.
	//
	// This method allows for sending an email with custom headers, such as tracking or
	// additional metadata. The headers are passed as a map of key-value pairs.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email.
	//	headers - A map of custom headers.
	//	isHtml - A flag indicating whether the email should be sent in HTML format.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendEmailWithHeaders(to []string, subject, body string, headers map[string]string, isHtml bool) error

	// ScheduleEmail schedules an email to be sent at a specific time. The isHtml flag determines text or HTML format.
	//
	// This method schedules an email to be sent at a future time, as specified by the
	// `sendAt` parameter. The email can be sent either as plain text or HTML.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email.
	//	sendAt - The time at which the email should be sent.
	//	isHtml - A flag indicating whether the email should be sent in HTML format.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	ScheduleEmail(to []string, subject, body string, sendAt time.Time, isHtml bool) error

	// SendEmailWithCCAndBCC sends an email with CC and BCC recipients. The isHtml flag determines text or HTML format.
	//
	// This method sends an email with additional recipients specified in the CC and BCC
	// fields. It can send both plain text and HTML emails depending on the value of `isHtml`.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	cc - A list of CC recipient email addresses.
	//	bcc - A list of BCC recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email.
	//	isHtml - A flag indicating whether the email should be sent in HTML format.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendEmailWithCCAndBCC(to []string, cc []string, bcc []string, subject, body string, isHtml bool) error

	// SendBulkEmail sends bulk emails. The isHtml flag determines text or HTML format.
	//
	// This method allows for sending bulk emails to multiple recipients in one go. It supports
	// both plain text and HTML formats.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email.
	//	isHtml - A flag indicating whether the email should be sent in HTML format.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendBulkEmail(to []string, subject, body string, isHtml bool) error

	// SendEmailWithAttachmentsAndInLineImages sends an email with both attachments and inline images.
	// Only applicable for HTML emails.
	//
	// This method sends an HTML email with inline images (which are embedded in the email body)
	// as well as file attachments. Inline images are usually specified by their file paths
	// and referenced in the email body using the `cid` attribute.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email (in HTML format).
	//	attachmentPaths - A list of file paths for the attachments.
	//	imagePaths - A list of file paths for the inline images.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendEmailWithAttachmentsAndInLineImages(to []string, subject, body string, attachmentPaths, imagePaths []string) error

	// SendEmailWithInLineImages sends an email with inline images.
	// Only applicable for HTML emails.
	//
	// This method sends an HTML email with images embedded within the email body, specified
	// by their file paths. The images are referenced in the email body using the `cid` attribute.
	//
	// Params:
	//	to - A list of recipient email addresses.
	//	subject - The subject of the email.
	//	body - The body content of the email (in HTML format).
	//	imagePaths - A list of file paths for the inline images.
	//
	// Returns:
	//	error - An error message if the email fails to send.
	SendEmailWithInLineImages(to []string, subject, body string, imagePaths []string) error
}
