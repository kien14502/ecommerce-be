package sendto

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"

	"github.com/kien14502/ecommerce-be/global"
	"go.uber.org/zap"
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ","))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
func SendTextEmailOTP(to []string, from, otp string) error {
	smtpConfig := global.Config.Smtp
	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "Test"},
		To:      to,
		Subject: "OTP Verification",
		Body:    fmt.Sprintf("Your OTP is %s. Please enter it to verify your account.", otp),
	}
	messageEmail := BuildMessage(contentEmail)

	auth := smtp.PlainAuth("", smtpConfig.User, smtpConfig.Pass, smtpConfig.Host)

	err := smtp.SendMail(smtpConfig.Host+":25", auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Send email failed", zap.Error((err)))
		return err
	}
	return nil
}

func SendTemplateEmailOTP(to []string, from, nameTemplate string, dataTemplate map[string]interface{}) error {
	htmlBody, err := getTemplate(nameTemplate, dataTemplate)
	if err != nil {
		return err
	}

	return sendTo(to, from, htmlBody)
}

func getTemplate(templateName string, templateData map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(templateName).ParseFiles("pkg/template/" + templateName))
	err := t.Execute(htmlTemplate, templateData)
	if err != nil {
		fmt.Println("err send email", err)
		return "", err
	}
	return htmlTemplate.String(), nil
}

func sendTo(to []string, from, template string) error {
	smtpConfig := global.Config.Smtp
	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "Test"},
		To:      to,
		Subject: "OTP Verification",
		Body:    template,
	}
	messageEmail := BuildMessage(contentEmail)

	auth := smtp.PlainAuth("", smtpConfig.User, smtpConfig.Pass, smtpConfig.Host)

	err := smtp.SendMail(smtpConfig.Host+":25", auth, from, to, []byte(messageEmail))
	if err != nil {
		global.Logger.Error("Send email failed", zap.Error((err)))
		return err
	}
	return nil
}
