package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/codingeasygo/util/xmap"
)

func TestLogin(t *testing.T) {
	// 	l := xmap.M{
	// 		"438731474": "qweasd888",
	// 	}

	// 	fmt.Println(loginByPassword("test", "test"))

	//	for k, v := range l {
	//		fmt.Printf("update tokens set token='%v' where id = %v;\n", loginByPassword(k, fmt.Sprintf("%v", v)), k)
	//	}
}

func TestChina(t *testing.T) {
	挖矿配置()
}

func TestJson(t *testing.T) {
	jsonFile := `/Users/wfunc/Downloads/play.h5avu.com.har`
	// 读取json文件
	data, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return
	}
	// 解析json
	var result xmap.M
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return
	}
	log := result.Map("log")
	entries := log.ArrayMapDef([]xmap.M{}, "entries")
	var conf = [][]string{}
	for _, m := range entries {
		if strings.Contains(m.Map("request").Str("url"), "game?cmd=useMiningItem") {
			fmt.Println(m.Map("request").Str("url"))
			a := strings.Split(m.Map("request").Str("url"), "?")
			b := strings.Split(a[1], "&")
			var itemId, row, column, now string
			for _, v := range b {
				c := strings.Split(v, "=")
				if c[0] == "itemId" {
					itemId = c[1]
				}
				if c[0] == "row" {
					row = c[1]
				}
				if c[0] == "column" {
					column = c[1]
				}
				if c[0] == "now" {
					now = c[1]
				}
			}
			conf = append(conf, []string{itemId, row, column})
			fmt.Println(itemId, row, column, now)
		}
	}
	// 把conf写入文件
	data, err = json.Marshal(conf)
	if err != nil {
		fmt.Println("Marshal: ", err.Error())
		return
	}
	err = ioutil.WriteFile("conf.json", data, 0644)
	if err != nil {
		fmt.Println("WriteFile: ", err.Error())
		return
	}
}

func TestReadConf(t *testing.T) {
	data, err := ioutil.ReadFile("useMiningItem.json")
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return
	}
	// 解析json
	var result [][]string
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return
	}
	a, b, c := 0, 0, 0
	for i, v := range result {
		if v[0] == "184" {
			a++
		}
		if v[0] == "185" {
			b++
		}
		if v[0] == "186" {
			c++
		}
		fmt.Println(i, v)
	}
	fmt.Println(a, b, c)
}
