package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	addr = "120.76.204.21:6666"
)

//esp_mac 11:11:11:11:11:11
//cus_mac 112211221122
//BacthSid bf54c4adf0

type TestData struct {
	DeviceType   string `json:"device_type"`
	FwVer        string `json:"fw_ver"`
	EspMac       string `json:"esp_mac"`
	CusMac       string `json:"cus_mac"`
	FlashID      string `json:"flash_id"`
	TestResult   string `json:"test_result"`
	FactorySid   string `json:"factory_sid"`
	BacthSid     string `json:"batch_sid"`
	Efuse        string `json:"efuse"`
	ChkRepeatFlg bool   `json:"chk_repeat_flg"`
	PoType       int    `json:"po_type"`
}

type HttpJson struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Data   TestData `json:"testdata"`
}

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	fmt.Println("已连接服务器")
	defer conn.Close()
	Client(conn)
}

func Client(conn net.Conn) {
	batchs := [2]string{"178333bef2", "bf54c4adf0"}
	esp_mac_s := "22:22:22:43:22:22"
	cus_mac_s := "33:33:33:43:33:33"

	//batch,esp_mac,cus_mac,cus_mac(r),op_type
	/*
			examples := [...][5]int{
				{1, 1, 1, 1, 0},
				{0, 0, 0, 0, 0},
				{1, 0, 0, 0, 0},
				{0, 1, 0, 0, 0},
				{1, 1, 0, 0, 0},
				{0, 0, 1, 0, 0},
				{1, 0, 1, 0, 0},
				{0, 1, 1, 0, 0},
				{1, 1, 1, 0, 0},
				{0, 0, 1, 1, 0},
				{1, 0, 1, 1, 0},
				{0, 1, 1, 1, 0},
				{1, 1, 1, 1, 0},
				{0, 0, 0, 0, 1},
				{1, 0, 0, 0, 1},
				{0, 1, 0, 0, 1},
				{1, 1, 0, 0, 1},
				{0, 0, 1, 0, 1},
				{1, 0, 1, 0, 1},
				{0, 1, 1, 0, 1},
				{1, 1, 1, 0, 1},
				{0, 0, 1, 1, 1},
				{1, 0, 1, 1, 1},
				{0, 1, 1, 1, 1},
				{1, 1, 1, 1, 1},
			}

			normal := [...][5]int{
				{1, 1, 1, 1, 0},
				{1, 0, 0, 0, 0},
			}

		burnMac := [...][5]int{
			{1, 1, 1, 1, 0},
			{0, 0, 1, 0, 1},
			{0, 1, 1, 0, 1},
			{0, 0, 1, 1, 1},
			{0, 1, 1, 1, 1},
		}
	*/
	redo := [...][5]int{
		{1, 1, 1, 1, 0},
		{0, 1, 0, 0, 1},
		{0, 1, 1, 0, 1},
		{0, 1, 1, 1, 1},
		{0, 0, 1, 1, 1},
	}

	var sendData HttpJson
	sendData = HttpJson{
		Path:   "/testdata",
		Method: "POST",
		Data: TestData{
			DeviceType:   "ESP_WROOM02",
			FwVer:        "v1.0.0.0",
			EspMac:       "22:22:22:22:22:10",
			CusMac:       "33:33:33:33:33:10",
			FlashID:      "19191919",
			TestResult:   "success",
			FactorySid:   "esp-own-test-FID-dbd42d01",
			BacthSid:     "bf54c4adf0",
			Efuse:        "1122334455667788",
			ChkRepeatFlg: true,
			PoType:       0,
		},
	}
	var result [100]string
	ri := 0
	esp_mac_add := 0
	cus_mac_add := 0
	for times := 0; times < 1; times++ {
		for i, example := range redo {
			if i == 0 && times != 0 {
				continue
			}

			sendData.Data.BacthSid = batchs[example[0]]
			if example[1] == 0 {
				esp_mac_add += 1
				sendData.Data.EspMac = esp_mac_s[:15] + strconv.Itoa(10+esp_mac_add)
			} else {
				sendData.Data.EspMac = esp_mac_s[:15] + strconv.Itoa(10)
			}

			if example[2] == 0 {
				sendData.Data.CusMac = ""
			} else {
				if example[3] == 0 {
					cus_mac_add += 1
					sendData.Data.CusMac = cus_mac_s[:15] + strconv.Itoa(10+cus_mac_add)
				} else {
					sendData.Data.CusMac = cus_mac_s[:15] + strconv.Itoa(10)
				}
			}

			sendData.Data.PoType = example[4]

			fmt.Println(sendData)
			bytes, _ := json.Marshal(sendData)
			bytes = append(bytes, '\n')
			conn.Write([]byte(bytes))

			buf := make([]byte, 1024)
			c, err := conn.Read(buf)
			if err != nil {
				fmt.Println("读取服务器数据异常:", err.Error())
			}
			rst := string(buf[0:c])
			fmt.Println(rst)
			status_index := strings.Index(rst, "\"status\": ")
			status := rst[status_index+10 : status_index+13]

			if status == "500" {
				result[ri] = "E"
			} else {
				tpr_index := strings.Index(rst, "\"test_pass_record\": \"")
				tpr := rst[tpr_index+21 : tpr_index+22]
				if tpr == "0" {
					result[ri] = "Y"
				} else {
					result[ri] = "W"
				}
			}
			ri += 1
		}
	}
	fmt.Println(result)
}
