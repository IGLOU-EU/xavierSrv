package tools

import (
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	mail "github.com/xhit/go-simple-mail"
)

func HttpReport(m string, u string, d string) {
	var res *http.Response
	var err error

	if m == "POST" {
		res, err = http.Post(u, "application/json", bytes.NewBufferString(d))
	} else {
		res, err = http.Get(u + url.QueryEscape(d))
	}

	if err != nil {
		PushToLog(3, err)
	}

	if res.StatusCode != 200 {
		PushToLog(3, errors.New("Request to `"+u+url.QueryEscape(d)+"` has failed: "+res.Status))
	}
}

func SmtpReport(host string, encrypt string, from string, user string, passwd string, recip []string, subject string, message string) {
	ht := strings.Split(host, ":")
	pt, _ := strconv.Atoi(ht[1])
	server := mail.NewSMTPClient()

	//SMTP Server
	server.Host = ht[0]
	server.Port = pt
	server.Username = user
	server.Password = passwd
	server.Authentication = mail.AuthPlain
	server.KeepAlive = true
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	smtpClient, err := server.Connect()

	switch encrypt {
	case "tls":
		server.Encryption = mail.EncryptionTLS
	case "ssl":
		server.Encryption = mail.EncryptionSSL
	default:
		server.Encryption = mail.EncryptionNone
	}

	if err != nil {
		PushToLog(3, err)
	}

	for _, to := range recip {
		email := mail.NewMSG()

		email.SetFrom(from).
			AddTo(to).
			SetSubject(subject)

		email.SetBody(mail.TextPlain, message)

		if err = email.Send(smtpClient); err != nil {
			PushToLog(3, err)
		}
	}
}
