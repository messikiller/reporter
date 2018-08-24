package util

import (
    "crypto/tls"

    gomail "github.com/go-mail/gomail"
)

func Mail(to string, subject string, body string) {
    from := ""
    host := ""
    port := 465
    username := ""
    password := ""

    m := gomail.NewMessage()

    m.SetHeader("From", from)

    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    d := gomail.NewDialer(host, port, username, password)
    d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

    if err := d.DialAndSend(m); err != nil {
        panic(err)
    }
}
