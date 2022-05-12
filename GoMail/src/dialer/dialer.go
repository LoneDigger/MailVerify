package dialer

import (
	"fmt"

	"github.com/go-mail/mail"
	"me.mail/src/config"
	"me.mail/src/mq"
)

type Dialer struct {
	d    *mail.Dialer
	host string
}

func NewMail(host string, m config.Mail) *Dialer {
	d := mail.NewDialer(m.Host, m.Port, m.Username, m.Password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	return &Dialer{
		d:    d,
		host: host,
	}
}

func (d *Dialer) Send(msg mq.Mail) error {
	link := fmt.Sprintf("<a href='http://%s/register/%s'>link</a>", d.host, msg.Token)

	message := mail.NewMessage()
	message.SetHeader("From", "bot.junx@gmail.com")
	message.SetHeader("To", msg.Mail)
	message.SetHeader("Subject", "信箱驗證郵件")
	message.SetAddressHeader("Bcc", "bot.junx@gmail.com", "bot")
	message.SetBody("text/html", link)

	return d.d.DialAndSend(message)
}
