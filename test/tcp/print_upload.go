package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"
)

const (
	addr = "192.168.16.186:6666"
)

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

type BodyData struct {
	DeviceType string `json:"device_type"`
	BacthSid   string `json:"batch_sid"`
}

type GetData struct {
	EspMac string `json:"esp_mac"`
	Dryrun bool   `json:"dryrun"`
}

type HttpJson struct {
	Path   string   `json:"path"`
	Method string   `json:"method"`
	Data   TestData `json:"testdata"`
	Body   BodyData `json:"body"`
	Get    GetData  `json:"get"`
}

//333334331210
func main() {
	esp_mac_1 := "22:22:22:22:24:10"

	channel1 := make(chan float64, 1)
	go Client(esp_mac_1, channel1)
	<-channel1
}

type Requestbody struct {
	req string
}

func (r *Requestbody) Json2map() (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(r.req), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func Client(esp_mac_s string, channel chan float64) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	fmt.Println("已连接服务器")
	defer conn.Close()

	var r Requestbody
	//cus_mac_s := "33:33:33:43:33:33"

	var sendData HttpJson
	sendData = HttpJson{
		Path:   "/testdata/print",
		Method: "POST",
		Data: TestData{
			DeviceType:   "ESP_WROOM02",
			FwVer:        "v1.0.0.0",
			EspMac:       "22:22:22:22:24:10",
			CusMac:       "",
			FlashID:      "19191919",
			TestResult:   "success",
			FactorySid:   "esp-fae-test-a95342f3",
			BacthSid:     "fd90f554b6",
			Efuse:        "1122334455667788",
			ChkRepeatFlg: true,
			PoType:       0,
		},
		Body: BodyData{
			DeviceType: "ESP-WROOM-02E",
			BacthSid:   "fd90f554b6",
		},
		Get: GetData{
			EspMac: "333335331311",
			Dryrun: false,
		},
	}

	var result [100]float64
	ri := 0
	// batch_index := -1
	esp_mac_add := 0
	// cus_mac_add := 0
	fmt.Println("start:", time.Now())
	for times := 0; times < 1; times++ {

		esp_mac_add += 1
		sendData.Data.EspMac = esp_mac_s[:15] + strconv.Itoa(10+esp_mac_add)

		// fmt.Println(sendData)
		bytes, _ := json.Marshal(sendData)
		bytes = append(bytes, '\n')
		conn.Write([]byte(bytes))

		buf := make([]byte, 1024)
		c, err := conn.Read(buf)
		if err != nil {
			fmt.Println("读取服务器数据异常:", err.Error())
		}
		r.req = string(buf[0:c])

		rst, err := r.Json2map()
		fmt.Println(rst["print_times"])

		ri += 1
	}
	fmt.Println(result)
	fmt.Println("end:", time.Now())
	channel <- 0
}
