package main

import (
	"bytes"
	"fmt"
	"log"
	"net/smtp"
)

const bodyTmpl = "Dear %s:\r\n\r\n\tThe account information you have just modify is as follows:\r\n\tAccount: %s\r\n\tPassword: %s\r\n\r\n\t(This message is send by system, do not reply)\r\n\r\nEspressif Systems · Production System\r\n"
const mailTmpl = "To: <%s>\r\nFrom: Factory System <factory@espressif.com>\r\nSubject: %s\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:45.0) Gecko/20100101 Thunderbird/45.2.0\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8; format=flowed\r\nContent-Transfer-Encoding: 7bit\r\n\r\n%s"

func SendToMail_1(email, subject, body string) error {
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	c.Mail("factory@espressif.com")
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

func SendToMail(to, subject, body string) error {
	fmt.Println(to, subject)
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	c.Mail("factory@espressif.com")
	c.Rcpt(to)
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return err
	}
	defer wc.Close()
	buf := bytes.NewBufferString(fmt.Sprintf(mailTmpl, to, subject, body))
	if _, err = buf.WriteTo(wc); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func main() {
	to := "hongmj0815@163.com"
	subject := "golang mail test"
	body := "just for test"
	fmt.Println("send email")
	// err := SendToMail(to, subject, "html")
	body = fmt.Sprintf(bodyTmpl, "hmj", "test", "test")
	err := SendToMail(to, subject, body)
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}

}
