package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

var ErrUninitialized = errors.New("smtp 未初始化")

func MailTo(to string, subject string, body string) error {
	if smtpConf == nil {
		return ErrUninitialized
	}

	return Smtp(smtpConf, to, subject, body)
}

var smtpConf *SmtpConfig

func Init(config SmtpConfig) error {
	smtpConf = &config
	return nil
}

type SmtpConfig struct {
	Host     string `json:"host" toml:"host"`
	Port     int    `json:"port" toml:"port"`
	From     string `json:"from" toml:"from"`         // 控制台创建的发信地址
	Password string `json:"password" toml:"password"` // 控制台设置的SMTP密码
}

func Smtp(config *SmtpConfig, to string, subject string, body string) error {

	header := make(map[string]string)
	header["From"] = "P2Link " + "<" + config.From + ">"
	header["To"] = to
	header["Subject"] = subject
	//html格式邮件
	header["Content-Type"] = "text/html; charset=UTF-8"

	//纯文本格式邮件
	//header["Content-Type"] = "text/plain; charset=UTF-8"
	//body := "test body"
	var message bytes.Buffer
	for k, v := range header {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(body)
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	err := SendMailWithTLS(
		fmt.Sprintf("%s:%d", config.Host, config.Port),
		auth,
		config.From,
		[]string{to},
		message.Bytes(),
	)
	return err
}

// DialSmtp return a smtp client
func DialSmtp(addr string) (*smtp.Client, error) {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS10, // TLSv1.0
		MaxVersion: tls.VersionTLS12, // TLSv1.2
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_RC4_128_SHA,
			tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
			// tls.TLS_RSA_WITH_RC4_128_MD5,
		},
	}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		log.Println("tls.Dial Error:", err)
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

// SendMailWithTLS send email with tls
func SendMailWithTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := DialSmtp(addr)
	if err != nil {
		log.Println("Create smtp client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
