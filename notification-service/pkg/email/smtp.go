package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPClient struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewSMTPClient(host string, port int, username, password, from string) *SMTPClient {
	return &SMTPClient{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (c *SMTPClient) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)

	auth := smtp.PlainAuth("", c.username, c.password, c.host)

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", c.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-version: 1.0;\r\nContent-Type: text/plain; charset=\"UTF-8\";\r\n\r\n")
	msg.WriteString(body)

	fmt.Println("Sending email to:", to)
	fmt.Println("Message:", msg.String())

	err := smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg.String()))
	if err != nil {
		fmt.Println("Error sending email:", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
