package notifications

import (
	"bytes"
	"fmt"
	"net/smtp"
	"slices"
	"strings"
	"text/template"

	"github.com/willtrojniak/TabAppBackend/models"
)

type MailDriver struct {
	from     string
	addr     string
	auth     smtp.Auth
	template *template.Template
}

func NewMailDriver(username, password, host, port string) *MailDriver {
	const templateName = "resources/templates/notifications/template.html"
	auth := smtp.PlainAuth("", username, password, host)
	t, err := template.ParseFiles(templateName)
	if err != nil {
		panic(err)
	}

	return &MailDriver{
		from:     username,
		addr:     fmt.Sprintf("%v:%v", host, port),
		auth:     auth,
		template: t,
	}
}

func (driver *MailDriver) Name() string {
	return "Mail"
}

func (driver *MailDriver) isDisabledFor(user *models.User) bool {
	return !user.EnableEmails
}

func (driver *MailDriver) NotifyShop(shop *models.Shop, n Notification) error {
	to := slices.DeleteFunc(shop.ConfirmedUsers(), func(u *models.User) bool {
		return n.IsDisabledFor(u, shop)
	})

	return driver.NotifyUsers(to, n)
}

func (driver *MailDriver) NotifyUsers(to []*models.User, n Notification) error {
	to = slices.DeleteFunc(to, driver.isDisabledFor)

	emails := make([]string, len(to))
	for i, u := range to {
		emails[i] = u.Email
	}

	html, err := driver.toHTML(emails, n)
	if err != nil {
		return err
	}

	return smtp.SendMail(
		driver.addr,
		driver.auth,
		driver.from,
		emails,
		html)
}

func (driver *MailDriver) toHTML(emails []string, n Notification) ([]byte, error) {
	var res bytes.Buffer
	err := driver.template.Execute(&res, n)
	if err != nil {
		return nil, err
	}

	return []byte(
		fmt.Sprintf("To: %v\n", strings.Join(emails, ",")) +
			fmt.Sprintf("Subject: %v\n", n.Heading()) +
			"MIME-version: 1.0;\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\n" +
			"\n" +
			fmt.Sprintf("%s\r\n", res.String())), nil

}
