package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

func SendToMail(to, subject, mailtype string) error {
	body = `
			<html>
			<body>
			<h3>
			"Test send to email"
			</h3>
			</body>
			</html>
			`
	user = "hongmj0815@163.com"
	password = "hmj8....."
	host = "smtp.163.com:25"
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

const mailTmpl = "To: <%s>\r\nFrom: hongmingjie <844263150@qq.com>\r\nSubject: %s\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:45.0) Gecko/20100101 Thunderbird/45.2.0\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8; format=flowed\r\nContent-Transfer-Encoding: 7bit\r\n\r\n%s"

func SendToMail_1(email, subject, body string) error {
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	c.Mail("noreply@factory.espressif.cn")
	c.Rcpt(email)
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return err
	}
	defer wc.Close()
	buf := bytes.NewBufferString(fmt.Sprintf(mailTmpl, email, subject, body))
	if _, err = buf.WriteTo(wc); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func main() {
	to := "hongmingjie@espressif.com"
	subject := "golang mail test"
	body := "just for test"
	fmt.Println("send email")
	// err := SendToMail(to, subject, "html")
	err := SendToMail_1(to, subject, body)
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

}
