package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("err: usage print.exe factory_sid batch_sid 5C:CF:7F:D0:FA:A1 [true/false]")
		return
	}
	factorySid := os.Args[1]
	batchSid := os.Args[2]
	espMac := os.Args[3]
	dryrun := "false"
	if len(os.Args) >= 5 && os.Args[4] == "true" {
		dryrun = "true"
	}
	pass := "FAILED"
	printTimes := 0

	resp, err := dial(factorySid, batchSid, espMac, dryrun)
	if err != nil {
		fmt.Println(err)
		return
	}
	m := map[string]interface{}{}
	err = json.Unmarshal([]byte(resp), &m)
	if err != nil {
		fmt.Println("err: parse response failed")
		return
	}
	status := int(m["status"].(float64))
	if status == 200 {
		result := m["test_result"].(string)
		if result == "success" {
			pass = "PASS"
			printTimes = int(m["print_times"].(float64))
		}
	}
	fmt.Println(fmt.Sprintf("%s, %d", pass, printTimes))
}

func dial(factorySid, batchSid, espMac, dryrun string) (string, error) {
	i := 0
	var rerr error
	for i < 3 {
		i++
		conn, err := net.DialTimeout("tcp", "factory.espressif.cn:6666", 3*time.Second)
		if err != nil {
			rerr = err
			continue
		}
		printReq := fmt.Sprintf(`{"path":"/testdata/print","method":"POST","get":{"factory_sid":"%s","batch_sid":"%s","esp_mac":"%s","dryrun":"%s"}}%s`, factorySid, batchSid, espMac, dryrun, "\n")
		_, err = conn.Write([]byte(printReq))
		if err != nil {
			rerr = err
			continue
		}
		resp, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			rerr = err
			continue
		}
		return resp, nil
	}
	return "", rerr
}
