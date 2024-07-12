package smtp

import (
	"net/smtp"
	"strings"
)

type SMTPClient interface {
	SendMail(to []string, msg []byte) error
}

type SMTPConfig struct {
	Addr     string
	Username string
	Password string
	From     string
}

type smtpClient struct {
	auth   smtp.Auth
	config SMTPConfig
}

func NewSMTPClient(config SMTPConfig) SMTPClient {
	host := strings.Split(config.Addr, ":")[0]
	auth := smtp.PlainAuth("", config.Username, config.Password, host)
	return &smtpClient{
		auth:   auth,
		config: config,
	}
}

func (c *smtpClient) SendMail(to []string, msg []byte) error {
	return smtp.SendMail(c.config.From, c.auth, c.config.Addr, to, msg)
}
