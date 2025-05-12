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

	cleanBody := strings.Replace(body, "Your verification code is: ", "", 1)

	htmlBody := fmt.Sprintf(`
	<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verification Code</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f4f4f4;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: #ffffff;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: #007bff;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .logo {
            max-width: 150px;
            height: auto;
            margin: 10px auto;
            display: block;
        }
        .content {
            padding: 30px;
        }
        .code {
            font-size: 24px;
            font-weight: bold;
            color: #007bff;
            text-align: center;
            padding: 20px;
            background: #f8f9fa;
            border-radius: 4px;
            margin: 20px 0;
        }
        .footer {
            text-align: center;
            padding: 20px;
            font-size: 12px;
            color: #666;
            border-top: 1px solid #eee;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Volatex</h1> <br/>
			<h3>Ваш код верификаци</h3>
        </div>
        <div class="content">
            <p>Здравствуйте,</p>
            <p>Пожалуйста, используйте следующий код для завершения регистрации:</p>
            <div class="code">%s</div>
            <p>Этот код действителен в течение 15 минут.</p>
            <p>Если вы не запрашивали этот код, пожалуйста, проигнорируйте это письмо.</p>
        </div>
    </div>
</body>
</html>`, cleanBody) // TODO: Вынести html в отдельный файл

	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", c.from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-version: 1.0;\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n")
	msg.WriteString(htmlBody)

	err := smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg.String()))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
