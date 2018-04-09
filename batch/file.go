package batch

import (
	"espressif.com/chip/factory/config"
	"espressif.com/chip/factory/dal"
	"espressif.com/chip/factory/rpc"
	"fmt"
	"io"
	"log"
	"os"
)

func File(req *rpc.Request, resp *rpc.Response) {
	action, _ := req.GetString("action", req.Get)
	batchSid, _ := req.GetString("batch_sid", req.Get)
	batch := dal.FindBatchBySid(req.Ctx, batchSid)
	if batch == nil {
		resp.Err = rpc.NotFound
		return
	}
	var err error
	switch action {
	case "upload":
		err = upload(req, resp, batch)
	case "download":
		err = download(req, resp, batch)
	}
	resp.Err = rpc.Skip
	if err != nil {
		resp.Err = err
	}
}

func upload(req *rpc.Request, resp *rpc.Response, batch *dal.Batch) error {
	err := req.R.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, xx := range req.R.MultipartForm.File {
		for _, x := range xx {
			file, err := x.Open()
			if err != nil {
				log.Println(err)
				return err
			}
			binfile, err := os.Create(fmt.Sprintf("%s/batch/%s.file", config.Cfg.DataDir, batch.Sid))
			if err != nil {
				log.Println(err)
				return err
			}
			defer binfile.Close()
			defer file.Close()
			io.CopyN(binfile, file, 10*1024*1024)
			break
		}
	}
	resp.Message = "upload batch file success"
	return nil
}

func download(req *rpc.Request, resp *rpc.Response, batch *dal.Batch) error {
	absPath := fmt.Sprintf("%s/batch/%s.file", config.Cfg.DataDir, batch.Sid)
	file, err := os.Open(absPath)
	if err != nil {
		log.Println(err)
		return err
	}
	resp.W.Header().Set("Content-Type", "application/octet-stream")
	resp.W.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s.file", batch.Sid))
	io.Copy(resp.W, file)
	return nil
}
