package user

import (
	"bytes"
	"crypto/md5"
	"errors"
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
	"fmt"
	"log"
	"net/smtp"
)

const bodyTmpl = "Dear %s:\r\n\r\n\tThe account information you have just modify is as follows:\r\n\tAccount: %s\r\n\tPassword: %s\r\n\r\n\t(This message is send by system, do not reply)\r\n\r\nEspressif Systems · Production System\r\n"
const mailTmpl = "To: <%s>\r\nFrom: hongmingjie <hongmingjie@espressif.com>\r\nSubject: %s\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:45.0) Gecko/20100101 Thunderbird/45.2.0\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8; format=flowed\r\nContent-Transfer-Encoding: 7bit\r\n\r\n%s"

func User(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Post:
		userPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func userPost(req *rpc.Request, resp *rpc.Response) {
	user, err := dal.UnmarshalUser(req.Ctx, req.GetValue("user", req.Body))
	if err != nil {
		resp.Err = err
		return
	}
	exist := dal.FindUserByAccount(req.Ctx, user.Account)
	if exist != nil {
		resp.Err = errors.New("Account Already Exist")
		return
	}
	originPassword := user.Password
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))
	err = user.Save()
	if err != nil {
		resp.Err = err
		return
	}
	resp.Body["user"] = user
	body := fmt.Sprintf(bodyTmpl, user.Name, user.Account, originPassword)
	err = SendToMail(user.Email, "factory account and password(auto send noreplay)", body)
}

func Login(req *rpc.Request, resp *rpc.Response) {
	account, _ := req.GetString("account", req.Get)
	password, _ := req.GetString("password", req.Get)
	user := dal.FindUserByAccount(req.Ctx, account)
	if user == nil {
		resp.Err = rpc.NotFound
		return
	}
	hash := md5.Sum([]byte(password))
	md5str := fmt.Sprintf("%x", hash) //[]byte16
	if md5str != user.Password {
		resp.Err = rpc.Conflict
		return
	}
	factory := dal.FindFactoryBySid(req.Ctx, user.FactorySid)
	if factory == nil {
		resp.Err = rpc.Conflict
		return
	}
	resp.Body["token"] = factory.Token
}

func Modify(req *rpc.Request, resp *rpc.Response) {
	switch req.Method {
	case rpc.Get:
		userModifyGet(req, resp)
	case rpc.Post:
		userModifyPost(req, resp)
	default:
		resp.Err = rpc.MethodNotAllowed
	}
}

func userModifyGet(req *rpc.Request, resp *rpc.Response) {
	account, _ := req.GetString("account", req.Get)
	user := dal.FindUserByAccount(req.Ctx, account)
	if user == nil {
		resp.Err = rpc.NotFound
		return
	}
	user.Password = ""
	resp.Body["user"] = user
}

func userModifyPost(req *rpc.Request, resp *rpc.Response) {
	user, err := dal.UnmarshalUser(req.Ctx, req.GetValue("user", req.Body))
	if err != nil {
		resp.Err = err
		return
	}
	originPassword := user.Password
	user.Password = fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))
	err = user.Update()
	if err != nil {
		resp.Err = err
		return
	}
	user.Password = ""
	resp.Body["user"] = user
	body := fmt.Sprintf(bodyTmpl, user.Name, user.Account, originPassword)
	err = SendToMail(user.Email, "factory account and password(auto send noreplay)", body)
}

func SendToMail(to, subject, body string) error {
	c, err := smtp.Dial("localhost:25")
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()
	c.Mail("noreply@factory.espressif.cn")
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
