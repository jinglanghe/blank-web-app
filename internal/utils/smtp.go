package utils

import (
	"crypto/tls"
	"errors"
	"net"
	"net/smtp"
	"strings"
)

// Dial Here is the key, you need to call tls.Dial instead of smtp.Dial
// for smtp servers running on 465 that require an ssl connection
// from the very beginning (no starttls)
func Dial(addr string, tlsConfig *tls.Config) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	_, port, _ := net.SplitHostPort(addr)
	if port == "465" {
		return SendMailWithoutTLS(addr, a, from, to, msg)
	} else {
		return smtp.SendMail(addr, a, from, to, msg)
	}
}

func SendMailWithoutTLS(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}
	host, _, _ := net.SplitHostPort(addr)
	tlsConfig := &tls.Config{
		ServerName: host,
	}
	c, err := Dial(addr, tlsConfig)
	if err != nil {
		return err
	}
	defer c.Close()
	if a != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err = c.Auth(a); err != nil {
			return err
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

// validateLine checks to see if a line has CR or LF as per RFC 5321
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}
