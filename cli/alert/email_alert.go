package alert

import (
	"bytes"
	"net/smtp"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Kshitiz-Mhto/xnap/pkg/config"
	"github.com/Kshitiz-Mhto/xnap/utility"
)

func sendHtmlEmailWithRetry(to []string, subject string, htmlBody string, maxRetries int, retryInterval time.Duration) error {
	auth := smtp.PlainAuth(
		"",
		config.Envs.FromEmail,
		config.Envs.FromEmailPassword,
		config.Envs.FromEmailSMTP,
	)

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	message := "Subject: " + subject + "\n" + headers + "\n\n" + htmlBody

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		lastError = smtp.SendMail(
			config.Envs.SMTPAddress,
			auth,
			config.Envs.FromEmail,
			to,
			[]byte(message),
		)
		if lastError == nil {
			return nil
		}

		// If failed, wait before retrying
		utility.Warning("Attempt %d failed: %s. Retrying in %v...", attempt, lastError.Error(), retryInterval)
		time.Sleep(retryInterval)

		// Exponentially increase the retry interval for next attempt
		retryInterval = retryInterval * 2
	}

	// Return the last error after all attempts
	utility.Error("Failed to send email after %d attempts: %s", maxRetries, lastError.Error())
	return lastError
}

func HTMLTemplateEmailHandler(addr string, vars map[string]interface{}) bool {
	basePathForEmailHtml := "./static/"
	emailSubject := "⚠️⚠️⚠️⚠️  ALERT  ⚠️⚠️⚠️⚠️"

	// Convert Param3 (comma-separated string) to a slice of strings
	to := strings.Split(addr, ",")

	// Parse the HTML template
	templatePath := filepath.Join(basePathForEmailHtml, "alert.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		utility.Error("failed to parse template: %v", err)
		return false
	}

	// Render the template with the map data
	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, vars); err != nil {
		utility.Error("failed to render template: %v", err)
		return false
	}

	// Define max retries and initial retry interval
	maxRetries := 3
	initialRetryInterval := 2 * time.Second

	// Attempt to send the email with retry logic
	err = sendHtmlEmailWithRetry(to, emailSubject, rendered.String(), maxRetries, initialRetryInterval)
	if err != nil {
		utility.Error("%s", err.Error())
		return false
	}

	return true
}
