package mailer

import (
	"fmt"
	"net/smtp"
	"strings"
)

type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`
}

type Mailer struct {
	cfg Config
}

func New(cfg Config) (*Mailer, error) {
	if cfg.Host == "" || cfg.Port == 0 || cfg.Username == "" || cfg.Password == "" || cfg.From == "" {
		return nil, fmt.Errorf("mailer: invalid smtp configuration")
	}

	return &Mailer{cfg: cfg}, nil
}

func (m *Mailer) SendResetPasswordEmail(to, resetLink string) error {
	subject := "Сброс пароля администратора"
	body := fmt.Sprintf(
		"Для восстановления пароля перейдите по ссылке:\n\n%s\n\nСсылка действительна 30 минут.",
		resetLink,
	)

	msg := buildMessage(m.cfg.From, []string{to}, subject, body)

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, []byte(msg))
}

func (m *Mailer) SendVerificationEmail(to, verificationLink string) error {
	subject := "Подтверждение email адреса"
	body := fmt.Sprintf(
		"Добро пожаловать!\n\nДля подтверждения вашего email адреса перейдите по ссылке:\n\n%s\n\nСсылка действительна 24 часа.\n\nЕсли вы не регистрировались на нашем сайте, просто проигнорируйте это письмо.",
		verificationLink,
	)

	msg := buildMessage(m.cfg.From, []string{to}, subject, body)

	addr := fmt.Sprintf("%s:%d", m.cfg.Host, m.cfg.Port)
	auth := smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.Host)

	return smtp.SendMail(addr, auth, m.cfg.From, []string{to}, []byte(msg))
}

func buildMessage(from string, to []string, subject, body string) string {
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", strings.Join(to, ", ")),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
		"",
	}

	return strings.Join(headers, "\r\n") + body
}
