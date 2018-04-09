package testdata

import (
	"archive/zip"
	"bytes"
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func Dump(req *rpc.Request, resp *rpc.Response) {
	action, err := req.GetString("action", req.Get)
	if err == nil && action == "download" {
		filename, _ := req.GetString("filename", req.Get)
		filepath := "/opt/backup/testdata/" + filename
		file, err := os.Open(filepath)
		if err != nil {
			log.Println(err)
			resp.Err = err
			return
		}
		w := resp.W
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
		io.Copy(w, file)
		return
	}
	email, err := req.GetString("email", req.Get)
	if err == nil && email != "" && strings.HasSuffix(email, "@espressif.com") {
		dumpAndSendToEmail(req, resp, email)
		return
	}

	testdatas := listTestdata(req)
	resp.Body["testdatas"] = testdatas
}

func dumpAndSendToEmail(req *rpc.Request, resp *rpc.Response, email string) {
	factorySid, _ := req.GetString("factory_sid", req.Get)
	batchSid, _ := req.GetString("batch_sid", req.Get)
	if factorySid == "" || batchSid == "" {
		resp.Err = rpc.BadRequest
		return
	}
	ts := time.Now().Unix()
	randstr := dal.Randstr(8)
	filename := fmt.Sprintf("dump.%s.%s.%d.%s.csv.zip", factorySid, batchSid, ts, randstr)
	go dumpAsCsv(factorySid, batchSid, filename, email)
	resp.Message = "send csv dump file url to your email, wait and see"
	resp.Body["url"] = "https://factory.espressif.cn/testdata/dump?action=download&filename=" + filename
}

func dumpAsCsv(factorySid, batchSid, filename, email string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	filepath := "/opt/backup/testdata/" + dal.Randstr(10)
	file, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	archive := zip.NewWriter(file)
	defer archive.Close()
	writer, err := archive.Create(filename[0 : len(filename)-4])
	if err != nil {
		log.Println(err)
		return err
	}

	lineTmpl := `%d, "%s", "%s", %t, %d, "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", %d, %d, %d, %t` + "\r\n"
	_, err = writer.Write([]byte("Id, Created, Updated, Visibly, ModuleId, DeviceType, FwVer, EspMac, CusMac, FlashId, TestResult, TestMsg, FactorySid, BatchSid, Efuse, QueryTimes, PrintTimes, BatchIndex, Latest\r\n"))
	if err != nil {
		log.Println(err)
		return err
	}
	i := 0
	for i < 100000 {
		testdatas := dal.ListTestdataByBatch(nil, batchSid, i, 1000)
		if len(testdatas) == 0 {
			break
		}
		i = i + 1000
		for _, x := range testdatas {
			line := fmt.Sprintf(lineTmpl, x.Id, x.Created, x.Updated, x.Visibly, x.ModuleId, x.DeviceType, x.FwVer, x.EspMac, x.CusMac, x.FlashId, x.TestResult, x.TestMsg, x.FactorySid, x.BatchSid, x.Efuse, x.QueryTimes, x.PrintTimes, x.BatchIndex, x.Latest)
			_, err := writer.Write([]byte(line))
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	archive.Flush()
	os.Rename(filepath, "/opt/backup/testdata/"+filename)
	subject := fmt.Sprintf("csv file dump for factory: %s, batch: %s", factorySid, batchSid)
	body := "click https://factory.espressif.cn/testdata/dump?action=download&filename=" + filename
	sendDumpCsvFileUrl(filename, email, subject, body)
	return nil
}

const mailTmpl = "To: <%s>\r\nFrom: wuyunzhou <wuyunzhou@espressif.com>\r\nSubject: %s\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:45.0) Gecko/20100101 Thunderbird/45.2.0\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8; format=flowed\r\nContent-Transfer-Encoding: 7bit\r\n\r\n%s"

func sendDumpCsvFileUrl(filename, email, subject, body string) error {
	log.Println(filename, email, subject, body)
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

func listTestdata(req *rpc.Request) []*dal.Testdata {
	req.FillRange(true, false)
	_, err := req.GetInt("row_count", req.Get)
	if err != nil {
		req.RowCount = 100
	}
	req.Factory = dal.FindFactoryByToken(req.Ctx, req.Token)
	espMac, err := req.GetString("esp_mac", req.Get)
	if err == nil {
		if req.Factory.IsStaff {
			return dal.ListTestdataByEspMac(req.Ctx, espMac, req.Offset, req.RowCount)
		} else {
			return dal.ListTestdataByFactoryEspMac(req.Ctx, req.Factory.Sid, espMac, req.Offset, req.RowCount)
		}
	}
	batchSid, err := req.GetString("batch_sid", req.Get)
	if err == nil {
		return dal.ListTestdataByBatch(req.Ctx, batchSid, req.Offset, req.RowCount)
	}
	factorySid, err := req.GetString("factory_sid", req.Get)
	if err == nil {
		return dal.ListTestdataByFactory(req.Ctx, factorySid, req.Offset, req.RowCount)
	}

	if req.Factory != nil {
		if req.Factory.IsStaff {
			return dal.ListTestdataAll(req.Ctx, req.Offset, req.RowCount)
		}
		return dal.ListTestdataByFactory(req.Ctx, req.Factory.Sid, req.Offset, req.RowCount)
	}
	return nil
}
