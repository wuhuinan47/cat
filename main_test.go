package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/codingeasygo/util/converter"
	"github.com/codingeasygo/util/xhttp"
	"github.com/codingeasygo/util/xmap"
	"github.com/wuhuinan47/cat/catdb"
)

func init() {
	db, err := sql.Open("mysql", "root:123@tcp(localhost:3306)/data_cat")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(100)                  //设置最大连接数
	db.SetMaxIdleConns(16)                   //设置闲置连接数

	Pool = db
	catdb.Pool = db
}

func TestGetFriendsUserIDs(t *testing.T) {
	userIDs := getFriendsUserIDs("302691822", "https://s147.11h5.com/", "ildSnyoMAOQOr31T82un39YP9SqiulbUOpb")
	fmt.Println(converter.JSON(userIDs))
}

func TestGame(t *testing.T) {
	fmt.Println(converter.JSON(game("ildv-8OmqkUICqn4NausGpcDGpQBb4An_qn")))
}

func TestGameWabao(t *testing.T) {
	wabao("695923850", "大南是逗比", "https://s147.11h5.com/", "dd84fab9b6caa625e0efd08d93e8e945")
}

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

func TestGetCanUseRange(t *testing.T) {
	// grids := [][]float64{{3, 0, 1, 0, -1, 0, 0}, {0, 0, 0, -1, -1, -1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 3, -1, -1, -1, 1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, 0}}
	// grids := [][]float64{{0, 0, 0, -1, -1, -1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 3, -1, -1, -1, 1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 1, 0}}
	// grids := [][]float64{{-1, -1, -1, -1, -1, -1, -1}, {1, 0, 3, -1, -1, -1, 1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {0, 0, 3, 1, 0, 4, 3}}
	// grids := [][]float64{{1, 0, 3, -1, -1, -1, 1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {0, 0, -1, 1, 0, 4, 3}, {0, 0, 1, 6, 0, 3, 0}}
	// grids := [][]float64{{1, 0, 3, -1, -1, -1, 1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, 6, 0, 3, 0}}
	// grids := [][]float64{{-1, -1, -1, -1, -1, -1, -1}, {0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, 0, 0, 0, 0, 3, 0}}
	// grids := [][]float64{{0, 0, 0, 0, 0, 0, -1}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, 0, 0, -1, 0, 3, 0}, {0, 5, 1, 0, 0, 0, 4}}
	// grids := [][]float64{{-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, 0, 0, -1, 0, 3, 0}, {0, 5, 1, -1, 0, 0, 4}, {1, 3, 4, 5, 0, 0, 0}}
	// grids := [][]float64{{-1, -1, -1, -1, -1, -1, -1}, {0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, 0, 0, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 3, 4, 5, 0, 0, 0}}
	// grids := [][]float64{{0, 0, -1, 0, 0, 0, 0}, {0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {0, 4, 3, 1, 5, 1, 0}}
	// grids := [][]float64{{0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {0, -1, 3, 1, 5, 1, 0}, {1, 0, 0, 3, 0, 0, 0}}
	// grids := [][]float64{{0, 0, -1, 0, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 0, 3, 0, 0, 0}}
	// grids := [][]float64{{-1, -1, -1, -1, -1, -1, -1}, {0, 0, 1, -1, 0, 3, 0}, {3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 0, -1, 0, 0, 0}, {0, 0, 0, 0, 0, 1, 0}}
	// grids := [][]float64{{0, 0, 1, -1, 0, 3, 0}, {3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 0, -1, 0, 0, 0}, {0, 0, 0, -1, 0, 1, 0}, {0, 0, 0, 0, 0, 0, 5}}
	// grids := [][]float64{{3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 0, -1, 0, 0, 0}, {0, 0, 0, -1, 0, 1, 0}, {0, 0, 0, -1, 0, 0, 5}, {0, 3, 0, 0, 4, 0, 3}}
	// grids := [][]float64{{3, -1, -1, -1, 0, 3, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, -1, -1, -1, 0, 0, 0}, {-1, -1, -1, -1, -1, -1, -1}, {1, 0, 0, -1, 0, 0, 0}, {0, 0, 0, -1, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, 3, 0, 0, 4, 0, 3}}
	// grids := [][]float64{{4, 0, 0, -1, 0, 3, 0}, {0, 3, 0, -1, 3, 0, 0}, {0, 0, 0, -1, 0, 0, 0}, {0, 0, 1, -1, 0, 1, 0}, {-1, -1, -1, -1, -1, -1, -1}, {0, -1, -1, -1, 0, 0, 0}, {3, -1, -1, -1, 0, 4, 0}, {0, 3, 0, 0, 0, 0, 0}}
	// grids := [][]float64{{
	// 	-1,
	// 	-1,
	// 	-1,
	// 	-1,
	// 	-1,
	// 	-1,
	// 	-1,
	// },
	// 	{
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 	},
	// 	{
	// 		0,
	// 		0,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		0,
	// 		0,
	// 	},
	// 	{
	// 		0,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		0,
	// 		3,
	// 	},
	// 	{
	// 		0,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		3,
	// 		0,
	// 	},
	// 	{
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 		-1,
	// 	},
	// 	{
	// 		0,
	// 		7,
	// 		3,
	// 		0,
	// 		-1,
	// 		0,
	// 		3,
	// 	},
	// 	{
	// 		0,
	// 		5,
	// 		3,
	// 		0,
	// 		0,
	// 		0,
	// 		0,
	// 	}}
	grids := [][]float64{{-1, -1, -1, -1, -1, -1, -1}, {0, -1, 0, 0, 0, 0, 0}, {-1, -1, -1, 0, 0, 4, 0}, {-1, -1, -1, 0, 0, 2, 4}, {-1, -1, -1, 0, 0, 0, 0}, {2, -1, 0, 0, 4, 1, 0}, {0, -1, 4, 4, 0, 3, 0}, {3, 1, 0, 0, 5, 0, 0}}

	canUse, rewards := getCanUseRange(grids)
	fmt.Println(converter.JSON(canUse))
	fmt.Println(converter.JSON(rewards))
	// 统计reward奖励
	item184Key, item185Key, item186Key, finalKey, win := doneRewards(rewards)
	fmt.Println(item184Key, item185Key, item186Key, finalKey)

	fmt.Println(win[item184Key])
	fmt.Println(win[item185Key])
	fmt.Println(win[item186Key])
	fmt.Println(win[finalKey])

}

func TestHttp(t *testing.T) {
	for i := 0; i < 1; i++ {
		data, err := xhttp.GetMap("https://s147.11h5.com//game?cmd=useCompanionWine&token=ild-1RUD5bPxZf4OkMEL0FYFbz4MMMEYJRS&count=1&now=1688105725249")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(converter.JSON(data))
	}
}
