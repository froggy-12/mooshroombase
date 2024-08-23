package smtp_configs

import (
	"bytes"
	"html/template"
	"net/smtp"
	"strconv"

	"github.com/froggy-12/mooshroombase/config"
)

type VerificationEmailData struct {
	Code             string
	Applink          string
	VerificationLink string
}

func SendVerificationEmail(emailTo string, code int) error {
	smtpServer := config.Configs.SMTPServerAdress
	smtpPort := config.Configs.SMTPServerPort

	from := config.Configs.SMTPEmailFrom
	password := config.Configs.SMTPPassword
	subject := "Email Verification"

	data := VerificationEmailData{
		Code:             strconv.Itoa(code),
		Applink:          config.Configs.FrontEndUrl[0] + config.Configs.VerifyEmailRouteClient,
		VerificationLink: config.Configs.FrontEndUrl[0] + config.Configs.EmailVerificationRoute,
	}

	tmpl := template.Must(template.New("email").Parse(`
		<!DOCTYPE html>
		<html lang="en">
		
		<head>
		
		  <meta charset="UTF-8">
		  <meta http-equiv="X-UA-Compatible" content="IE=edge">
		  <meta name="viewport" content="width=device-width, initial-scale=1.0">
		  <title>Title</title>
		  <style>
			body {
			  font-family: Arial, sans-serif;
			  background-color: #f0f0f0;
			}
			.container {
			  max-width: 600px;
			  margin: 40px auto;
			  padding: 20px;
			  background-color: #fff;
			  border: 1px solid #ddd;
			  border-radius: 10px;
			  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}
			h1 {
			  color: #333;
			  font-size: 24px;
			  margin-bottom: 10px;
			}
			a {
			  text-decoration: none;
			  color: #337ab7;
			}
			a:hover {
			  color: #23527c;
			}
			button {
			  background-color: #337ab7;
			  color: #fff;
			  border: none;
			  padding: 10px 20px;
			  font-size: 16px;
			  cursor: pointer;
			}
			button:hover {
			  background-color: #23527c;
			}
		  </style>
		</head>
		
		<body>
		  <div class="container">
			<h1>Lets Verify your account</h1>
			<h1>Your Code is <span>{{ .Code }}</span></h1>
			<a href="{{ .Applink }}">
			  <button>Open App</button>
			</a>
			<p>Not able to paste this code click down here:</p>
			<a href="{{ .VerificationLink }}">
			  <button>Verify</button>
			  <a href="{{ .VerificationLink }}">link</a>
			</a>
		  </div>
		</body>
		
		</html>
	`))

	var html bytes.Buffer

	err := tmpl.Execute(&html, data)
	if err != nil {
		return err
	}
	msg := "To: " + emailTo + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		html.String()

	auth := smtp.PlainAuth(from, from, password, smtpServer)

	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, from, []string{emailTo}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}
