package emails

import (
	"fmt"
	"net/smtp"
	"os"
)

func ConformationMail(data EmailDetails) (string, error) {
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpUsername == "" || smtpPassword == "" {
		fmt.Println("SMTP credentials not found in environment variables.")
		return "SMTP credentials not set", nil
	}

	// SMTP server details
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	//Email Details
	from := smtpUsername
	to := data.Mail
	headers := fmt.Sprintf("From: %s\r\nTo: %s\r\n", from, to)
	subject := "Subject: Your order Processing\r\n"
	body := fmt.Sprintf(`Dear %s 
Your order is currently being processed. Thank you for choosing Holoware Computers!

We appreciate your trust in Holoware. If you have any questions, feel free to reach out to us.


Best regards,
Holoware Computers Team
Support@holoware.co
www.holoware.co`, data.Name)

	message := []byte(headers + subject + "\r\n" + body)

	// Authentication
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	//adding bcc
	recipients := []string{to}

	// Send the email
	err2 := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, recipients, message)
	if err2 != nil {
		fmt.Println("Error sending email:", err2)
		return "", err2
	}
	return "Conformation Email Successfully sent", nil
}
