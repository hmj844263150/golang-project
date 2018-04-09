package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

var semiRe = regexp.MustCompile("^([a-zA-Z0-9_-]+)[,:](.+)")
var vdd33Re = regexp.MustCompile("^vdd33=([0-9]+),")
var txVdd33Re = regexp.MustCompile("^TX_VDD33=([0-9]+),[ ]*([0-9]+)")
var failureModeRe = regexp.MustCompile("^Part failure in ([a-zA-Z0-9_]+)")
var keys = map[string]interface{}{"DUT_RXRSSI": struct{}{}, "FB_RXRSSI": struct{}{}, "FREQ_OFFSET": struct{}{}, "MAC": struct{}{}, "RXDC": struct{}{}, "RXIQ": struct{}{}, "RX_NOISEFLOOR": struct{}{}, "TXDC": struct{}{}, "TXIQ": struct{}{}, "TX_POWER_BACKOFF": struct{}{}, "TXP_RESULT": struct{}{}, "TX_VDD": struct{}{}, "TX_VDD33": struct{}{}, "VDD": struct{}{}, "VDD33": struct{}{}, "TXCAP_TMX2G_CCT_LOAD": struct{}{}, "TXCAP_PA2G_CCT_STG1": struct{}{}, "TXCAP_PA2G_CCT_STG2": struct{}{}, "TX_PWRCTRL_ATTEN": struct{}{}, "BT_PA_GAIN": struct{}{}, "BT_DIG_ATTEN": struct{}{}, "BT_TX_BB": struct{}{}, "BT_TXIQ": struct{}{}, "BT_TXDC": struct{}{}}

func put(m map[string][]string, key string, values string) {
	vs := []string{}
	for _, v := range strings.Split(values, " ") {
		v = strings.Trim(v, " \r\n,;")
		v = strings.Replace(v, ",", "", -1)
		v = strings.Replace(v, ";", "", -1)
		if v == "" {
			continue
		}
		vs = append(vs, v)
	}
	m[key] = vs
}

func sortAndPrint(m map[string][]string) {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		log.Println(k, m[k])
	}
}

func sortAndWriteXlsx(m map[string][]string) {
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	row = sheet.AddRow()
	for _, key := range keys {
		values, ok := m[key]
		if !ok {
			panic(fmt.Sprintf("key '%s' not found", key))
		}
		if len(values) == 1 {
			cell = row.AddCell()
			cell.Value = key
			continue
		}
		for i, _ := range values {
			cell = row.AddCell()
			cell.Value = fmt.Sprintf("%s_%d", key, i)
		}
	}

	row = sheet.AddRow()
	for _, key := range keys {
		values := m[key]
		for _, v := range values {
			cell = row.AddCell()
			cell.Value = v
		}
	}

	err = file.Save("failure.xlsx")
	if err != nil {
		panic(err)
	}
}
func main() {
	f, err := os.Open("/home/cloudzhou/fail_240AC4041740_2017-03-16_155856.txt")
	//f, err := os.Open("/home/cloudzhou/Downloads/fail_5CCF7F360616_2017-02-07_160731.txt")
	//f, err := os.Open("/home/cloudzhou/Downloads/fail_5CCF7F2AE5D9_2016-12-28_181210.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	m := map[string][]string{}
	failureModes := []string{}
	for {
		bs, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		line := bytes.Trim(bs, "\r\n ")
		if semiRe.Match(line) {
			subs := semiRe.FindSubmatch(line)
			key := string(bytes.ToUpper(subs[1]))
			if _, ok := keys[key]; !ok {
				continue
			}
			values := string(bytes.ToUpper(subs[2]))
			put(m, key, values)
			continue
		}
		if vdd33Re.Match(line) {
			subs := vdd33Re.FindSubmatch(line)
			put(m, "VDD33", string(bytes.ToUpper(subs[1])))
		}
		if txVdd33Re.Match(line) {
			subs := txVdd33Re.FindSubmatch(line)
			put(m, "TX_VDD33_1", string(bytes.ToUpper(subs[1])))
			put(m, "TX_VDD33_2", string(bytes.ToUpper(subs[2])))
		}
		if failureModeRe.Match(line) {
			subs := failureModeRe.FindSubmatch(line)
			if len(subs) < 2 {
				continue
			}
			key := string(bytes.ToUpper(subs[1]))
			if _, ok := keys[key]; !ok {
				continue
			}
			found := false
			for _, failure := range failureModes {
				if key == failure {
					found = true
				}
			}
			if found {
				continue
			}
			failureModes = append(failureModes, key)
		}
		//log.Println(line)
	}
	m["FAILURE_MODES"] = failureModes
	raw, err := json.Marshal(&m)
	log.Println(string(raw))
	//log.Println(m)
	//sortAndPrint(m)
	//sortAndWriteXlsx(m)
}
