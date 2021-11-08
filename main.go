package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/codingeasygo/util/xprop"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wuhuinan47/cat/crawler"
	"github.com/wuhuinan47/cat/runner"
	wechatapi "github.com/wuhuinan47/cat/wechatAPI"
)

var Pool *sql.DB

var run_dir, api_url string

// type CatData struct {
// 	GoldMineHelpList []GoldMineHelpList `json:"helpList"`
// }

// type GoldMineHelpList struct {
// 	Uid     int64 `json:"uid"`
// 	Quality int64 `json:"quality"`
// }

type User struct {
	Uid           string
	Token         string
	Name          string
	ServerURL     string
	ZoneToken     string
	Public        map[string]interface{}
	FamilyDayTask interface{}
}

func main() {

	confPath := "/etc/cat-srv/cat.properties"
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}
	var err error
	cfg := xprop.NewConfig()
	cfg.Load(confPath)
	cfg.Print()

	// db, err := sql.Open("mysql", "cat:cyydmkj123@tcp(localhost:3306)/data_cat?charset=utf8&parseTime=true&collation=utf8mb4_unicode_ci")
	db, err := sql.Open("mysql", cfg.Str("mysql_con"))
	if err != nil {
		panic(err)
	}

	run_dir = "/mnt/app/cat/"
	api_url = "https://mcps.51yizhuan.com:13010/"

	run_dir = cfg.Str("run_dir")
	api_url = cfg.Str("api_url")

	db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(100)                  //设置最大连接数
	db.SetMaxIdleConns(16)                   //设置闲置连接数

	Pool = db

	http.HandleFunc("/login", LoginH)
	http.HandleFunc("/update", UpdateH)
	http.HandleFunc("/test", TestH)
	http.HandleFunc("/pullAnimal", PullAnimalH)
	http.HandleFunc("/diamond", DiamondH)
	http.HandleFunc("/familySign", familySignH)
	http.HandleFunc("/sendQrcode", SendQrcodeH)
	http.HandleFunc("/sendQQQrcode", SendQQQrcodeH)

	http.HandleFunc("/loginByQrcode", LoginByQrcodeH)
	http.HandleFunc("/setConfig", SetConfigH)
	http.HandleFunc("/getAllTokens", GetAllTokensH)

	http.HandleFunc("/hitCandy", HitCandyH)
	http.HandleFunc("/getFreeBossCannon", GetFreeBossCannonH)

	http.HandleFunc("/attackBoss", AttackBossH)
	http.HandleFunc("/sonAttackBoss", SonAttackBossH)
	http.HandleFunc("/oneSonAttackBoss", OneSonAttackBossH)
	http.HandleFunc("/giftPiece", GiftPieceH)
	http.HandleFunc("/draw", DrawH)
	http.HandleFunc("/throwDice", ThrowDiceH)
	http.HandleFunc("/getMailPrize", GetMailPrizeH)
	http.HandleFunc("/getBossPrize", GetBossPrizeH)
	http.HandleFunc("/singleBossAttack", SingleBossAttackH)
	http.HandleFunc("/setPiece", SetPieceH)
	http.HandleFunc("/sixEnergy", SixEnergyH)
	http.HandleFunc("/attackMyBoss", AttackMyBossH)
	http.HandleFunc("/beachHelp", BeachHelpH)
	http.HandleFunc("/useShovel", UseShovelH)
	http.HandleFunc("/getBeachLineRewards", GetBeachLineRewardsH)
	http.HandleFunc("/setRunner", SetRunnerH)
	http.HandleFunc("/buildUp", BuildUpH)
	http.HandleFunc("/openSteamBox", OpenSteamBoxH)
	http.HandleFunc("/setPullRows", SetPullRowsH)
	http.HandleFunc("/addFirewood", AddFirewoodH)
	http.HandleFunc("/exchangeRiceCake", ExchangeRiceCakeH)
	http.HandleFunc("/unlockWorker", UnlockWorkerH)
	http.HandleFunc("/searchFamily", SearchFamilyH)

	http.HandleFunc("/getServerURL", GetServerURLH)
	http.HandleFunc("/getZoneToken", GetZoneTokenH)

	//
	http.HandleFunc("/checkToken", CheckTokenH)

	//

	// wechatAPI
	http.HandleFunc("/wechatAPI/sendMsg", wechatapi.SendMsgH)

	http.HandleFunc("/ctrl.html", IndexH)
	http.HandleFunc("/cat_demo.html", CatDemoH)
	// http.HandleFunc("/qqQrCode.png", QQQrCodeH)

	http.HandleFunc("/", StaticServer)

	http.HandleFunc("/getServerLogs", GetServerLogsH)
	http.HandleFunc("/delServerLogs", DelServerLogsH)
	//

	//
	// http.HandleFunc("/", IndexH)

	sendMsg("cat is start")

	running := true
	{
		go runner.NamedRunnerWithSeconds("RunnerPullAnimal", 1800, &running, RunnerPullAnimal)
		go runner.NamedRunnerWithSeconds("RunnerDraw", 3700, &running, RunnerDraw)
		go runner.NamedRunnerWithSeconds("RunnerSteamBox", 1900, &running, RunnerSteamBox)
	}

	defer func() {
		sendMsg("cat is stop")
		log.Println("All done, good bye")
	}()

	go http.ListenAndServe(cfg.Str("listen"), nil)
	log.Printf("listen port %s sucesss", cfg.Str("listen"))

	// if err = http.ListenAndServe(":33333", nil); err != nil {
	// 	log.Fatal(err)
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

}

func StaticServer(w http.ResponseWriter, r *http.Request) {

	http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))).ServeHTTP(w, r)
}

func IndexH(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("ctrl.html")
	t.Execute(w, nil)
}

func CatDemoH(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("cat.html")
	t.Execute(w, nil)
}

func GetServerLogsH(w http.ResponseWriter, req *http.Request) {
	_, err := os.Stat("screenlog.0")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	t, _ := template.ParseFiles("screenlog.0")
	t.Execute(w, nil)
}

func DelServerLogsH(w http.ResponseWriter, req *http.Request) {
	cmd := fmt.Sprintf("rm -rf %vscreenlog.0", run_dir)
	c := exec.Command("bash", "-c", cmd)
	_, err := c.CombinedOutput()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "SUCCESS")
}

func GetAllTokensH(w http.ResponseWriter, req *http.Request) {

	SQL := "select id, name from tokens"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()

	var list []map[string]interface{}

	for rows.Next() {

		var id, name string
		rows.Scan(&id, &name)
		list = append(list, map[string]interface{}{"id": id, "name": name})
	}

	bytes, err := json.Marshal(list)

	if err != nil {
		return
	}

	w.Write(bytes)
}

//

func SendQrcodeH(w http.ResponseWriter, req *http.Request) {
	qrcode := req.URL.Query().Get("qrcode")
	Pool.Exec("update config set conf_value = ? where conf_key = 'wechatLoginQrcode'", qrcode)

	sendMsg(qrcode)
	log.Println("qrcode:", qrcode)
	io.WriteString(w, qrcode)
}

func SendQQQrcodeH(w http.ResponseWriter, req *http.Request) {
	// qrcode := req.URL.Query().Get("qrcode")

	imgBytes, _ := ioutil.ReadAll(req.Body)

	cmd := fmt.Sprintf("rm -rf %vstatic/qqQrCode.png", run_dir)
	exec.Command("bash", "-c", cmd)

	if err := ioutil.WriteFile("./static/qqQrCode.png", imgBytes, 0644); err != nil {
		log.Println(err)
	}

	Pool.Exec("update config set conf_value = ? where conf_key = 'qqLoginQrcode'", fmt.Sprintf("%vstatic/qqQrCode.png", api_url))
	sendMsg(fmt.Sprintf("%vstatic/qqQrCode.png", api_url))

	io.WriteString(w, fmt.Sprintf("%vstatic/qqQrCode.png", api_url))
}

func SetConfigH(w http.ResponseWriter, req *http.Request) {
	confKey := req.URL.Query().Get("confKey")
	confValue := req.URL.Query().Get("confValue")

	_, err := Pool.Exec("update config set conf_value = ? where conf_key = ?", confValue, confKey)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, "SUCCESS")
}

func SingleBossAttackH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='mmBoss1'))"
	if id == "2" {
		SQL = "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='cowBoss1'))"

	} else if id == "3" {
		SQL = "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='boss3'))"
	} else if id == "4" {
		SQL = "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='newBoss1'))"
	} else if id == "5" {
		SQL = "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='newBoss2'))"
	} else if id == "0" {
		SQL = "select id, token from tokens where id <> 302691822 or id <> 309392050 order by id desc"
	}

	log.Println("SingleBossAttackH SQL :", SQL)
	rows, err := Pool.Query(SQL)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var uid, token string
		rows.Scan(&uid, &token)
		log.Println("SingleBossAttackH uid :", uid)

		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		if zoneToken != "" {
			bossList := getBossHelpList(serverURL, zoneToken)
			for _, v := range bossList {
				leftHp, ok := v["leftHp"].(float64)
				if ok {
					if leftHp >= 1150 {
						attackBoss(serverURL, zoneToken, v["id"].(string))
					}
				}

			}
		}

	}

	io.WriteString(w, "SUCCESS")
}

func AttackBossH(w http.ResponseWriter, req *http.Request) {
	// https://s147.11h5.com:3147/118_89_154_132/3150//game?cmd=getBossHelpList&token=ildvxZcIKv6sMNVWloJb72r4iiiurzy7376&now=1629467019423

	id := req.URL.Query().Get("id")

	log.Println("AttackBossH id :", id)

	var SQL, name, token string
	if id == "" {
		SQL = "select name, token from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
	} else {
		SQL = fmt.Sprintf("select name, token from tokens where id = %v", id)
	}

	Pool.QueryRow(SQL).Scan(&name, &token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	bossList := getBossHelpList(serverURL, zoneToken)
	log.Printf("---------------------------[%v]开始打龙---------------------------", name)

	for _, v := range bossList {
		leftHp := v["leftHp"].(float64)
		log.Println("leftHp:", leftHp)
		if leftHp <= 600 && leftHp >= 500 {
			attackBossAPI(serverURL, zoneToken, v["id"].(string), 3, 0, 1, 200, 200)
		} else if leftHp <= 1000 && leftHp >= 500 {
			attackBossAPI(serverURL, zoneToken, v["id"].(string), 3, 0, 1, 200, 200)
			attackBossAPI(serverURL, zoneToken, v["id"].(string), 1, 1, 1, 400, 400)
			log.Println("finish")
		} else {
			log.Println("ignore")
		}

		// if leftHp < 1000 && leftHp >= 500 {
		// 	attackBossByAdmin(serverURL, zoneToken, v["id"].(string))
		// }
	}

	log.Printf("---------------------------[%v]结束打龙---------------------------", name)

	io.WriteString(w, "SUCCESS")

}

//
func SixEnergyH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	for i := 0; i < 6; i++ {
		getSixEnergy(serverURL, zoneToken)
		time.Sleep(time.Millisecond * 600)
	}
	io.WriteString(w, "SUCCESS")
}
func SetPieceH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	ids := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, v := range ids {
		log.Printf("设置拼图%v", v)

		setPiece(serverURL, zoneToken, v)
		// time.Sleep(time.Millisecond * 300)
	}
	flag := getPiecePrize(serverURL, zoneToken)

	log.Printf("拼图奖励领取状态:%v", flag)

	if !flag {
		io.WriteString(w, "https://play.h5avu.com/game/?gameid=147&token="+token)
		return
	}

	io.WriteString(w, "SUCCESS")
}

//
func GetBossPrizeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	SQL := "select token, name from tokens where id != (select conf_value from config where conf_key = 'cowBoy')"

	if id != "" {
		SQL = fmt.Sprintf("select token, name from tokens where id = %v", id)
	}
	rows, err := Pool.Query(SQL)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var token, name string
		rows.Scan(&token, &name)

		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		bossIDList := getBossPrizeList(serverURL, zoneToken)

		if len(bossIDList) != 0 {
			log.Printf("[%v]开始领取龙奖励", name)
			for _, bossID := range bossIDList {
				getBossPrize(serverURL, zoneToken, bossID)
				log.Printf("[%v]领取了[%v]奖励", name, bossID)

				time.Sleep(time.Second * 1)
			}
			log.Printf("[%v]领取龙奖励完毕", name)
		} else {
			log.Printf("[%v]没有可领龙奖励", name)
		}

	}

	io.WriteString(w, "SUCCESS")

}

//
func GetMailPrizeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	title := req.URL.Query().Get("title")
	SQL := fmt.Sprintf("select name, token from tokens where id = %v", id)

	if id == "" {
		SQL = "select name, token from tokens"
	}

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name, token string

		rows.Scan(&name, &token)

		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		mailids := getMailList(serverURL, zoneToken, title)

		log.Printf("[%v] 开始领取邮件奖励[title:%v]", name, title)
		for _, mailid := range mailids {
			getMailAttachments(serverURL, zoneToken, mailid)
		}

		log.Printf("[%v] 领取邮件完毕", name)

		go func() {
			log.Printf("[%v] 开始删除邮件", name)
			for _, mailid := range mailids {
				log.Printf("deleteMail mailid:%v", mailid)
				readMail(serverURL, zoneToken, mailid)
				deleteMail(serverURL, zoneToken, mailid)
			}
			log.Printf("[%v] 删除邮件完毕", name)
		}()

		time.Sleep(time.Millisecond * 200)

	}

	io.WriteString(w, "SUCCESS")

}

func ThrowDiceH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	amount := req.URL.Query().Get("amount")
	intAmount, _ := strconv.Atoi(amount)
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	go func() {
		log.Println("start throwDice")

		for i := 1; i <= intAmount; i++ {
			throwDice(serverURL, zoneToken)

			if i == intAmount {
				break
			}
			time.Sleep(time.Second * 1)
		}
		log.Println("end throwDice")

	}()

	io.WriteString(w, "SUCCESS")
}

func UseShovelH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	toid := req.URL.Query().Get("toid")

	if id == toid {
		toid = ""
	}

	SQL := "select id, name, token from tokens where id = ?"

	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	if toid == "all" {
		SQL = "select id, name from tokens where find_in_set(id, (select conf_value from config where conf_key = 'beachUidList'))"
		rows, err := Pool.Query(SQL)
		if err != nil {
			return
		}
		defer rows.Close()

		for rows.Next() {
			var uuid, uname string
			rows.Scan(&uuid, &uname)
			for i := 1; i <= 5; i++ {
				log.Printf("[%v]为[%v]使用铲子x%v", name, uname, i)
				useShovel(serverURL, zoneToken, uuid)
				time.Sleep(time.Millisecond * 100)
			}
		}
		io.WriteString(w, "SUCCESS")
		return

	}

	if toid == "" {
		for i := 1; i <= 20; i++ {
			log.Printf("[%v]使用铲子x%v", name, i)
			useShovel(serverURL, zoneToken, "")
			time.Sleep(time.Millisecond * 100)
		}
	} else {

		var toname string
		Pool.QueryRow("select name from tokens where id = ?", toid).Scan(&toname)

		for i := 1; i <= 5; i++ {
			log.Printf("[%v]为[%v]使用铲子x%v", name, toname, i)
			useShovel(serverURL, zoneToken, toid)
			time.Sleep(time.Millisecond * 100)
		}
	}
	io.WriteString(w, "SUCCESS")
}

func GetBeachLineRewardsH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := fmt.Sprintf("select id, name, token from tokens where id = %s", id)

	if id == "" {
		SQL = "select id, name, token from tokens"
	}

	rows, err := Pool.Query(SQL)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var uid, name, token string
		rows.Scan(&uid, &name, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		for i := 0; i <= 4; i++ {
			log.Printf("[%v]领取海滩奖励x%v", name, i)
			getBeachLineRewards(serverURL, zoneToken, i, 0)
			time.Sleep(time.Millisecond * 100)
		}

		for i := 0; i <= 4; i++ {
			log.Printf("[%v]领取海滩奖励x%v", name, i)
			getBeachLineRewards(serverURL, zoneToken, i, 1)
			time.Sleep(time.Millisecond * 100)
		}
	}

	io.WriteString(w, "SUCCESS")
}

func OpenSteamBoxH(w http.ResponseWriter, req *http.Request) {
	openSteamBoxGo()
	io.WriteString(w, "SUCCESS")
}

func SetPullRowsH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	input := req.URL.Query().Get("input")
	if id == "" {
		Pool.Exec("update tokens set pull_rows = ?", input)
	} else {
		Pool.Exec("update tokens set pull_rows = ? where id = ?", input, id)
	}
	io.WriteString(w, "SUCCESS")
}

func ExchangeRiceCakeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := fmt.Sprintf("select id, name, token from tokens where id = %s", id)
	if id == "" {
		SQL = "select id, name, token from tokens"
	}
	rows, err := Pool.Query(SQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid, name, token string
		rows.Scan(&uid, &name, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		ids := []int64{5, 4, 3, 2, 1}
		for _, v := range ids {
			for {
				flag := exchangeRiceCake(serverURL, zoneToken, v)
				if !flag {
					break
				}
				log.Printf("[%s]领取汤圆奖励[%v]成功", name, v)
			}
			for {
				flag := exchangeXmas(serverURL, zoneToken, v)
				if !flag {
					break
				}
				log.Printf("[%s]领取骰子奖励[%v]成功", name, v)
			}

			for {
				flag := exchangeBeachReward(serverURL, zoneToken, v)
				if !flag {
					break
				}
				log.Printf("[%s]领取海滩奖励[%v]成功", name, v)

			}
		}
	}
	io.WriteString(w, "SUCCESS")
}

func AddFirewoodH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	SQL := "select id, name, token from tokens where id = ?"

	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	uids := getSteamBoxHelpList(serverURL, zoneToken)

	for _, v := range uids {
		fuid := fmt.Sprintf("%v", v)
		addFirewood(serverURL, zoneToken, fuid)
		log.Printf("[%v]给[%v]添加柴火", name, fuid)
	}

	io.WriteString(w, "SUCCESS")
}

func BuildUpH(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")

	SQL := "select id, name, token from tokens where id = ?"

	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	var ids = []int64{1, 2, 3, 4, 5}

	var islandid float64
	for _, id := range ids {
		if id == 1 {
			for i := 1; i <= 5; i++ {
				log.Printf("开始过岛 建筑%v->等级->%v", id, i)
				buildUp(serverURL, zoneToken, id)
			}
		} else {
			for i := 1; i <= 5; i++ {
				log.Printf("开始过岛 建筑%v->等级->%v", id, i)
				islandid = buildUp(serverURL, zoneToken, id)
			}
		}
	}

	if islandid != 0 {
		log.Printf("领取过岛能量")
		getIslandPrize(serverURL, zoneToken, islandid)
		log.Printf("领取过岛分享能量")
		getIslandEnergy(serverURL, zoneToken)
	}

	io.WriteString(w, "SUCCESS")
}

func SearchFamilyH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := "select id, name, token from tokens where id = ?"
	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)
	log.Printf("[name:%v][id:%v]", name, id)
	serverURL, zoneToken := getSeverURLAndZoneToken(token)
	familyId, timeFlushList := getFamilyId(serverURL, zoneToken)
	familyName := searchFamily(serverURL, zoneToken, familyId)

	familyName, _ = url.QueryUnescape(familyName)

	mapList := map[string]interface{}{"familyName": familyName, "familyId": familyId, "timeFlushList": timeFlushList}
	jsonBytes, err := json.Marshal(mapList)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, string(jsonBytes))
}

func UnlockWorkerH(w http.ResponseWriter, req *http.Request) {
	// mineList
	id := req.URL.Query().Get("id")
	SQL := "select id, name, token from tokens where id = ?"
	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)
	serverURL := getServerURL()

	zoneToken, mineList := getEnterInfo(uid, name, serverURL, token, "mineList")
	l1, ok := mineList.(map[string]interface{})
	if ok {
		for k, v := range l1 {
			v1, ok := v.(map[string]interface{})
			if ok {
				unlockNum, ok := v1["unlockNum"].(float64)
				if ok {
					if unlockNum != 5 {
						var i float64
						for i = 1; i <= 5-unlockNum; i++ {
							log.Printf("---------------------------[%v]解锁[%v]---------------------------", name, k)
							time.Sleep(time.Millisecond * 200)
							unlockWorker(serverURL, zoneToken, k)
						}
					}
				}
			}
		}
	}
	io.WriteString(w, "SUCCESS")
}

func SetRunnerH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := fmt.Sprintf("select id, name, token from tokens where id = %s", id)

	if id == "" {
		SQL = "select id, name, token from tokens"
	}

	Pool.QueryRow(SQL, id)

	rows, err := Pool.Query(SQL)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var uid, name, token string
		rows.Scan(&uid, &name, &token)
		serverURL := getServerURL()

		zoneToken, beach := getEnterInfo(uid, name, serverURL, token, "beach")

		if beach != nil {
			beachMap, ok := beach.(map[string]interface{})

			if ok {
				endTime, ok := beachMap["endTime"].(float64)
				if ok {
					now := time.Now().UnixNano() / 1e6

					seconds := int64(endTime) - now

					grids, ok := beachMap["grids"].([]interface{})
					var gridsnum = 0
					if ok {
						for _, v := range grids {
							vv, ok := v.([]interface{})
							if ok {
								for _, v2 := range vv {
									vvv, ok := v2.(map[string]interface{})
									if ok {
										gridsuid, ok := vvv["uid"].(string)
										if ok {
											if gridsuid != "" {
												gridsnum += 1
											}
										}
									}
								}
							}
						}
					}
					log.Printf("[%v] gridsnum is :%v", name, gridsnum)

					if gridsnum == 25 {
						if seconds > 0 {
							seconds += 2000
							fmt.Println("set runner after ", seconds)

							go func() {
								time.Sleep(time.Millisecond * time.Duration(seconds))
								// serverURL = getServerURL()
								// zoneToken := getZoneToken(serverURL, token)

								serverURL, zoneToken = getSeverURLAndZoneToken(token)

								refreshBeach(serverURL, zoneToken)
								fmt.Printf("[%s] refreshBeach finish ", name)
								for i := 1; i <= 20; i++ {
									log.Printf("[%v]使用铲子x%v", name, i)
									useShovel(serverURL, zoneToken, "")
									time.Sleep(time.Millisecond * 100)
								}
							}()
						} else {
							refreshBeach(serverURL, zoneToken)

							for i := 1; i <= 20; i++ {
								log.Printf("[%v]使用铲子x%v", name, i)
								useShovel(serverURL, zoneToken, "")
								time.Sleep(time.Millisecond * 100)
							}
						}
					}
				}
			}
		}
	}

	io.WriteString(w, "SUCCESS")
}

func DrawH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	drawMulti := req.URL.Query().Get("drawMulti")
	amount := req.URL.Query().Get("amount")
	intDrawMulti, _ := strconv.Atoi(drawMulti)
	intAmount, _ := strconv.Atoi(amount)

	SQL := "select id, name, token from tokens where id = ?"

	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)

	serverURL := getServerURL()

	zoneToken, familyDayTask := getEnterInfo(uid, name, serverURL, token, "familyDayTask")

	// zoneToken := getZoneToken(serverURL, token)

	go func() {
		log.Printf("---------------------------[%v]开始转盘---------------------------", name)
		followCompanion_1(serverURL, zoneToken, 2)

		energy := draw(uid, name, serverURL, zoneToken, drawMulti)

		targetEnergy := intDrawMulti * intAmount

		if targetEnergy > int(energy) {
			intAmount = int(energy) / intDrawMulti
		}
		time.Sleep(time.Millisecond * 2100)

		for i := 1; i <= intAmount; i++ {
			draw(uid, name, serverURL, zoneToken, drawMulti)
			log.Println("剩余转盘次数:", intAmount-i)
			time.Sleep(time.Millisecond * 2100)
		}

		taskIDs := getDayTasksInfo(serverURL, zoneToken)
		log.Printf("[%v]领取日常任务奖励:%v", name, taskIDs)
		for _, taskID := range taskIDs {
			// time.Sleep(time.Millisecond * 100)
			getDayTaskAward(serverURL, zoneToken, taskID)
		}

		log.Printf("[%v]领取超值返利", name)
		getElevenEnergyPrize(serverURL, zoneToken, 1)
		getElevenEnergyPrize(serverURL, zoneToken, 2)
		getElevenEnergyPrize(serverURL, zoneToken, 3)
		getElevenEnergyPrize(serverURL, zoneToken, 4)

		log.Printf("[%v]collectMineGold", name)
		collectMineGold(serverURL, zoneToken)

		if id == "302691822" || id == "309392050" {
			followCompanion_1(serverURL, zoneToken, 3)
		} else {
			followCompanion_1(serverURL, zoneToken, 1)
		}
		log.Printf("---------------------------[%v]结束转盘---------------------------", name)

		// serverURL := getServerURL()

		// zoneToken, familyDayTask := getEnterInfo(uid, name, serverURL, token, "familyDayTask")

		familyDayTask, ok := familyDayTask.(map[string]interface{})

		if ok {
			for k := range familyDayTask {
				getFamilyDayTaskPrize(serverURL, zoneToken, k)
				log.Printf("[%v]领取公会任务奖励[%v]", name, k)
			}
		}

	}()

	io.WriteString(w, "SUCCESS")

}

func GiftPieceH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	fromUid := req.URL.Query().Get("fromUid")
	toUid := req.URL.Query().Get("toUid")
	amount := req.URL.Query().Get("amount")

	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, fromUid).Scan(&token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	intAmount, _ := strconv.Atoi(amount)

	log.Println("amount:", amount)
	log.Println("intAmount:", intAmount)
	for i := 1; i <= intAmount; i++ {
		giftPiece(serverURL, zoneToken, id, toUid)
	}

	io.WriteString(w, "SUCCESS")
}

func BeachHelpH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		io.WriteString(w, "id is null")
		return
	}
	go func() {
		SQL := "select id, token, name from tokens where id != ?"

		rows, err := Pool.Query(SQL, id)

		if err != nil {
			return
		}

		defer rows.Close()

		for rows.Next() {
			var uid, token, name string
			rows.Scan(&uid, &token, &name)

			serverURL, zoneToken := getSeverURLAndZoneToken(token)

			if zoneToken != "" {
				// log.Printf("[%v] 海浪助力", name)
				// beachHelp(serverURL, zoneToken, id, 42)
				time.Sleep(time.Second * 1)
				log.Printf("[%v] 铲子助力", name)
				beachHelp(serverURL, zoneToken, id, 43)
			} else {
				sendMsg(uid + ":" + name)
				log.Printf("[uid: %v] token is invalid\n", uid)
			}
		}

	}()
	io.WriteString(w, "SUCCESS")

}

func AttackMyBossH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	mode := req.URL.Query().Get("mode")
	SQL := fmt.Sprintf("select name, token from tokens where id = %s", id)

	if id == "" {
		SQL = "select name, token from tokens where id <> 302691822 and id <> 309392050 "
	}

	rows, err := Pool.Query(SQL)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name, token string
		rows.Scan(&name, &token)
		go func() {

			serverURL, zoneToken := getSeverURLAndZoneToken(token)
			if zoneToken != "" {
				bossID := summonBoss(serverURL, zoneToken)
				if bossID == "" {
					log.Printf("---------------------------[%v]无法打龙---------------------------", name)
					return
				}
				inviteBoss(serverURL, zoneToken, bossID)
				time.Sleep(time.Second * 1)
				shareAPI(serverURL, zoneToken)
				getFreeBossCannon(serverURL, zoneToken)
				log.Printf("---------------------------[%v]开始打龙---------------------------", name)
				attackMyBoss(serverURL, zoneToken, bossID, mode)
				log.Printf("---------------------------[%v]结束打龙---------------------------", name)
			}
		}()
		time.Sleep(time.Millisecond * 1300)
	}

	io.WriteString(w, "SUCCESS")

}

func OneSonAttackBossH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	var confKey = "newBoss1"

	if id == "2" {
		confKey = "newBoss2"
	}
	if id == "3" {
		confKey = "newBoss3"
	}

	if id == "4" {
		confKey = "cowBoss1"
	}

	if id == "5" {
		confKey = "mmBoss1"
	}

	if id == "6" {
		confKey = "boss3"
	}

	SQL := fmt.Sprintf("select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='%s'))", confKey)

	rows, err := Pool.Query(SQL)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer rows.Close()
	var mmList1 []map[string]interface{}
	var mmList2 []map[string]interface{}
	var mmBossList []string
	var j = 1
	for rows.Next() {
		var token string
		var uid float64
		rows.Scan(&uid, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		if zoneToken != "" {
			bossID := summonBoss(serverURL, zoneToken)
			inviteBoss(serverURL, zoneToken, bossID)
			time.Sleep(time.Second * 1)
			shareAPI(serverURL, zoneToken)
			getFreeBossCannon(serverURL, zoneToken)
			if j >= 5 {
				mmList2 = append(mmList2, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
			} else {
				mmList1 = append(mmList1, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
			}
			j++
			mmBossList = append(mmBossList, bossID)
			time.Sleep(time.Second * 1)
		}
	}

	for _, v2 := range mmBossList {
		for _, v := range mmList1 {
			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag {
				go attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}
		}
		time.Sleep(time.Second * 18)
	}

	for _, v2 := range mmBossList {
		for _, v := range mmList2 {
			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag {
				attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}
		}
		time.Sleep(time.Second * 18)
	}

	// for _, v := range mmList {
	// 	for _, v2 := range mmBossList {
	// 		flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
	// 		if flag == true {
	// 			attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
	// 		}
	// 	}
	// }

	io.WriteString(w, "SUCCESS")
}

func SonAttackBossH(w http.ResponseWriter, req *http.Request) {
	SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='mmBoss1'))"
	rows, err := Pool.Query(SQL)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer rows.Close()
	var mmList []map[string]interface{}
	var mmBossList []string
	for rows.Next() {
		var token string
		var uid float64
		rows.Scan(&uid, &token)

		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		if zoneToken != "" {
			bossID := summonBoss(serverURL, zoneToken)
			inviteBoss(serverURL, zoneToken, bossID)
			mmList = append(mmList, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
			mmBossList = append(mmBossList, bossID)
		}
	}
	for _, v := range mmList {
		for _, v2 := range mmBossList {
			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag {
				attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}
		}
	}
	SQL = "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='cowBoss1'))"
	rows, err = Pool.Query(SQL)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	var nnList []map[string]interface{}
	var nnBossList []string
	for rows.Next() {
		var token string
		var uid float64
		rows.Scan(&uid, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		if zoneToken != "" {
			bossID := summonBoss(serverURL, zoneToken)
			inviteBoss(serverURL, zoneToken, bossID)
			nnList = append(nnList, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
			nnBossList = append(nnBossList, bossID)
		}

	}
	for _, v := range nnList {
		for _, v2 := range nnBossList {

			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag {
				attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}
		}
	}

	SQL = "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key='boss3'))"
	rows, err = Pool.Query(SQL)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	var boss3List []map[string]interface{}
	var boss3BossList []string
	for rows.Next() {
		var token string
		var uid float64
		rows.Scan(&uid, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		if zoneToken != "" {
			bossID := summonBoss(serverURL, zoneToken)
			inviteBoss(serverURL, zoneToken, bossID)
			boss3List = append(boss3List, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
			boss3BossList = append(boss3BossList, bossID)
		}
	}
	for _, v := range boss3List {
		for _, v2 := range boss3BossList {

			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag {
				attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}

		}
	}

	io.WriteString(w, "SUCCESS")
}

//

func LoginByQrcodeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	if id == "1" {
		Pool.Exec("update config set conf_value = '' where conf_key = 'wechatLoginQrcode'")
		str := fmt.Sprintf("%vlinuxWechat.py", run_dir)
		SQL := "select conf_value from config where conf_key = 'wechatLoginQrcode'"
		go func() {
			log.Println("start python")
			cmd := exec.Command("/bin/python3", str)
			err := cmd.Run()
			if err != nil {
				return
			}
			Pool.Exec("update config set conf_value = '' where conf_key = 'wechatLoginQrcode'")
			log.Println("end python")
		}()

		var qrcode string
		for i := 0; i < 4; i++ {
			Pool.QueryRow(SQL).Scan(&qrcode)
			if qrcode != "" {
				break
			}
			time.Sleep(time.Second * 1)
		}
		io.WriteString(w, qrcode)

		return
	}

	if id == "3" {
		Pool.Exec("update config set conf_value = '' where conf_key = 'wechatLoginQrcode'")
		str := fmt.Sprintf("%vlinuxWechat.py", run_dir)
		SQL := "select conf_value from config where conf_key = 'wechatLoginQrcode'"
		go func() {
			log.Println("start python")
			cmd := exec.Command("/bin/python3", str)
			err := cmd.Run()
			if err != nil {
				return
			}
			Pool.Exec("update config set conf_value = '' where conf_key = 'wechatLoginQrcode'")
			log.Println("end python")
		}()

		var qrcode string
		for i := 0; i < 4; i++ {
			Pool.QueryRow(SQL).Scan(&qrcode)
			if qrcode != "" {
				break
			}
			time.Sleep(time.Second * 1)
		}

		http.Redirect(w, req, qrcode, http.StatusTemporaryRedirect)
		return
	}

	if id == "2" {
		cmd := fmt.Sprintf("rm -rf %vstatic/qqQrCode.png", run_dir)
		c := exec.Command("bash", "-c", cmd)
		c.CombinedOutput()
		go func() {
			log.Println("start QQScan")
			token, uid, err := crawler.QQScan()
			if err != nil {
				log.Println("err QQScan", err)
				return
			}
			log.Println("end QQScan")

			var uuid string
			Pool.QueryRow("select id from tokens where id = ?", uid).Scan(&uuid)

			if uuid == "" {
				serverURL := getServerURL()
				_, name := getNickName(serverURL, fmt.Sprintf("%v", token))
				_, err := Pool.Exec("INSERT INTO tokens (id, name, token) values (?, ?, ?)", uid, name, token)
				log.Println("INSERT:", err)
			} else {
				_, err := Pool.Exec("update tokens set token = ? where id = ?", token, uid)
				log.Println("update:", err)

			}

		}()

		for i := 0; i < 6; i++ {
			_, err := os.Stat(fmt.Sprintf("%vstatic/qqQrCode.png", run_dir))
			if err == nil {
				http.Redirect(w, req, fmt.Sprintf("%vstatic/qqQrCode.png", api_url), http.StatusTemporaryRedirect)
				return
			}
			time.Sleep(time.Second * 1)
		}
		io.WriteString(w, "FAIL GET QR CODE PNG")

		return

	}

	if id == "4" {
		cmd := fmt.Sprintf("rm -rf %vstatic/qqQrCode.png", run_dir)
		c := exec.Command("bash", "-c", cmd)
		c.CombinedOutput()
		go func() {
			log.Println("start QQScan")
			token, uid, err := crawler.QQScan()
			if err != nil {
				log.Println("err QQScan", err)

				return
			}
			log.Println("end QQScan")
			var uuid string
			Pool.QueryRow("select id from tokens where id = ?", uid).Scan(&uuid)

			if uuid == "" {
				serverURL := getServerURL()
				_, name := getNickName(serverURL, fmt.Sprintf("%v", token))
				_, err := Pool.Exec("INSERT INTO tokens (id, name, token) values (?, ?, ?)", uid, name, token)
				log.Println("INSERT:", err)
			} else {
				_, err := Pool.Exec("update tokens set token = ? where id = ?", token, uid)
				log.Println("update:", err)

			}
		}()

		for i := 0; i < 6; i++ {
			_, err := os.Stat(fmt.Sprintf("%vstatic/qqQrCode.png", run_dir))
			if err == nil {
				io.WriteString(w, fmt.Sprintf("%vstatic/qqQrCode.png", api_url))
				return
			}
			time.Sleep(time.Second * 1)
		}
		io.WriteString(w, "FAIL GET QR CODE PNG")

		return

	}

	io.WriteString(w, "qrcode")

}

func TestH(w http.ResponseWriter, req *http.Request) {

	// serverURL := getServerURL()
	// SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"
	// rows, err := Pool.Query(SQL)
	// if err != nil {
	// 	return
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var uid, token string
	// 	rows.Scan(&uid, &token)
	// 	gameList := []string{"535", "525", "157", "452", "411"}
	// 	for _, v := range gameList {
	// 		getAward(token, v)
	// 	}
	// }
	qrcode := req.URL.Query().Get("qrcode")

	sendMsg(qrcode)
	log.Println("qrcode:", qrcode)
}

func GetServerURLH(w http.ResponseWriter, req *http.Request) {
	serverURL := getServerURL()
	io.WriteString(w, serverURL)
}

func GetZoneTokenH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	serverURL := req.URL.Query().Get("serverURL")
	var token string
	Pool.QueryRow("select token from tokens where id = ?", id).Scan(&token)

	log.Println("serverURL:", serverURL)
	log.Println("token:", token)
	zoneToken := getZoneToken(serverURL, token)
	io.WriteString(w, zoneToken)
}

func CheckTokenH(w http.ResponseWriter, req *http.Request) {
	// go func() {
	SQL := "select id, token, name from tokens"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()

	var groupconcat1, groupconcat2 string

	for rows.Next() {
		var id, token, name string
		rows.Scan(&id, &token, &name)

		serverURL := getServerURL()
		zoneToken, flag := getZoneToken_1(serverURL, token)

		if zoneToken == "" {
			Pool.Exec("update tokens set token = '' where id = ?", id)
			groupconcat1 += "["
			groupconcat1 += name
			groupconcat1 += "]"
			groupconcat1 += ":"
			groupconcat1 += "失效"
			groupconcat1 += "/"

			sendMsg(id + ":" + name)
		} else {
			if flag {
				groupconcat2 += name
				groupconcat2 += ":"
				groupconcat2 += fmt.Sprintf("%v", flag)
				groupconcat2 += "/"
			}

		}

		// time.Sleep(time.Millisecond * 100)

	}

	// SQL = "select group_concat(name) as groupname from tokens where token = ''"
	// var reply string
	// Pool.QueryRow(SQL).Scan(&reply)

	io.WriteString(w, groupconcat1+"||||"+groupconcat2)

}

func GetFreeBossCannonH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	var token string
	Pool.QueryRow("select token from tokens where id = ?", id).Scan(&token)
	serverURL, zoneToken := getSeverURLAndZoneToken(token)
	shareAPI(serverURL, zoneToken)
	getFreeBossCannon(serverURL, zoneToken)
	io.WriteString(w, "SUCCESS")
}

func HitCandyH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	quality := req.URL.Query().Get("quality")
	amount := req.URL.Query().Get("amount")

	log.Printf("id is %v, quality is %v, amount is %v", id, quality, amount)

	v2, _ := strconv.ParseFloat(quality, 64)
	v3, _ := strconv.ParseFloat(amount, 64)

	var token string
	Pool.QueryRow("select token from tokens where id = ?", id).Scan(&token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)

	uids := getFriendsCandyTreeInfo(serverURL, zoneToken, v2)
	time.Sleep(time.Second * 1)
	log.Println("getFriendsCandyTreeInfo uids:", uids)

	var targetAmount float64

	for _, v := range uids {
		log.Println("getCandyTreeInfo uid:", v)
		// posList
		getCandyTreeInfo(serverURL, zoneToken, v)
		// log.Println("posList:", posList)
		time.Sleep(time.Second * 1)
		log.Println("start uid:", v)

		testList := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		for _, v4 := range testList {

			if v3 == targetAmount {
				return
			}

			time.Sleep(time.Second * 2)

			flag := hitCandyTree(serverURL, zoneToken, v, v4)

			log.Println("hit flag:", flag)

			if flag == "err" {
				break
			}

			if flag == "true" {
				targetAmount += 1
				break
			}

		}

		log.Println("end uid:", v)

	}

	io.WriteString(w, "SUCCESS")

}

func DiamondH(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")
	quality := req.URL.Query().Get("quality")
	amount := req.URL.Query().Get("amount")

	v2, _ := strconv.ParseFloat(quality, 64)
	v3, _ := strconv.ParseFloat(amount, 64)

	total, all := getBoxPrizeGo(id, v2, v3)

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (total/all)*100), 64)
	str := fmt.Sprintf("抓取次数:%v, 抓中次数:%v, 命中率:%v", all, total, fmt.Sprintf("%v", rate)+"%")

	io.WriteString(w, str)
}

func familySignH(w http.ResponseWriter, req *http.Request) {

	uid := req.URL.Query().Get("id")

	if uid != "" {
		SQL := "select token, name from tokens where id = ?"

		var token, name string
		Pool.QueryRow(SQL, uid).Scan(&token, &name)

		log.Printf("[%v] getSeverURLAndZoneToken", name)

		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		log.Printf("[%v] zoneToken", zoneToken)
		log.Printf("[%v] token", token)

		if zoneToken == "" {
			io.WriteString(w, "FAIL")
			return
		}

		gameList := []string{"535", "525", "157", "452", "411"}

		for _, v := range gameList {
			getAward(token, v)
		}

		log.Printf("[%v] start familySign", name)
		familySign(serverURL, zoneToken)
		log.Printf("[%v] end familySign", name)
		time.Sleep(time.Second * 1)
		log.Printf("[%v] start getSignPrize", name)
		getSignPrize(serverURL, zoneToken)
		log.Printf("[%v] end getSignPrize", name)
		time.Sleep(time.Second * 1)
		log.Printf("[%v] start getFreeDailyGiftBox", name)
		getFreeDailyGiftBox(serverURL, zoneToken)
		log.Printf("[%v] end getFreeDailyGiftBox", name)
		time.Sleep(time.Second * 1)

		log.Printf("[%v] start playLuckyWheel", name)
		shareAPI(serverURL, zoneToken)
		playLuckyWheel(serverURL, zoneToken)
		log.Printf("[%v] end playLuckyWheel", name)
		time.Sleep(time.Second * 1)

		log.Printf("[%v] start getFreeClamp", name)
		shareAPI(serverURL, zoneToken)
		getFreeClamp(serverURL, zoneToken)
		log.Printf("[%v] end getFreeClamp", name)
		time.Sleep(time.Second * 1)

		log.Printf("[%v] start getInviteSnow", name)
		shareAPI(serverURL, zoneToken)
		getInviteSnow(serverURL, zoneToken)
		log.Printf("[%v] end getInviteSnow", name)

		time.Sleep(time.Second * 1)
		log.Printf("[%v] start autoFriendEnergy", name)
		autoFriendEnergy(serverURL, zoneToken)
		log.Printf("[%v] end autoFriendEnergy", name)

		log.Printf("[%v] start getSixEnergy", name)
		for i := 0; i < 6; i++ {
			getSixEnergy(serverURL, zoneToken)
		}
		log.Printf("[%v] end getSixEnergy", name)

		log.Printf("[%v] 公会聊天", name)
		familyChat(serverURL, zoneToken)

		log.Printf("[%v] collectMineGold", name)
		collectMineGold(serverURL, zoneToken)
		io.WriteString(w, "SUCCESS")
		return
	}

	go familySignGo()
	io.WriteString(w, "SUCCESS")
}

func PullAnimalH(w http.ResponseWriter, req *http.Request) {
	go pullAnimalGo()
	io.WriteString(w, "SUCCESS")

}

func UpdateH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	name := req.URL.Query().Get("name")
	token := req.URL.Query().Get("token")
	var serverURL, zoneToken string
	if name == "" {
		serverURL = getServerURL()
		zoneToken, name = getNickName(serverURL, token)
	}

	log.Println("update token id is ", id, "name is ", name, "token is", token)

	SQL := "replace into tokens (id, name, token, serverURL, zoneToken) values (?, ?, ?, ?, ?)"
	_, err := Pool.Exec(SQL, id, name, token, serverURL, zoneToken)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, "SUCCESS")

}

func LoginH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	var token string
	Pool.QueryRow("select token from tokens where id = ?", id).Scan(&token)

	if token == "" {
		io.WriteString(w, "cannot find token")
		return
	}
	url := "https://play.h5avu.com/game/?gameid=147&token="
	url += token

	http.Redirect(w, req, url, http.StatusTemporaryRedirect)

}

// goroutine

// 一键拉动物
func pullAnimalGo() {
	log.Println("start pullAnimalGo")
	SQL := "select id, token, name, pull_rows from tokens where id != 302691822"
	rows, err := Pool.Query(SQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid, token, name, pullRows string
		rows.Scan(&uid, &token, &name, &pullRows)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		if zoneToken == "" {
			sendMsg(uid + ":" + name)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}

		log.Printf("[%v]开始拉动物", name)
		foods := enterFamilyRob(serverURL, zoneToken)
		for _, v := range foods {
			// myTeam := v["myTeam"].(int)
			// if myTeam != 4 {
			if strings.Contains(pullRows, fmt.Sprintf("%v", v["row"])) {
				if robFamilyFood(serverURL, zoneToken, v["id"].(string)) {
					break
				}
			}
			// }

		}
		log.Printf("[%v]拉动物完成", name)
		// time.Sleep(time.Second * 1)
		// log.Printf("serverURL:%v, zoneToken:%v\n", serverURL, zoneToken)
	}
}

// 一键公会签到/签到/免费福利/免费转盘/免费夹子/赠送能量
func familySignGo() {
	log.Println("familySignGo start...")

	SQL := "select id, token, name from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()

	var tokenList []map[string]string

	var helpInfo = make(map[string]int64)

	for rows.Next() {
		var uid, token, name string
		rows.Scan(&uid, &token, &name)

		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		if zoneToken != "" {
			tokenList = append(tokenList, map[string]string{"uid": uid, "name": name, "serverURL": serverURL, "zoneToken": zoneToken, "token": token})
			helpInfo["uid"] = 0
		} else {
			sendMsg(uid + ":" + name)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}
	}

	var j = 1

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start enterSteamBox", user["name"])
		enterSteamBox(user["serverURL"], user["zoneToken"], user["uid"])
		log.Printf("[%v] end enterSteamBox", user["name"])
		j++
	}

	for _, user := range tokenList {

		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}

		log.Printf("[%v] start familySign", user["name"])
		familySign(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end familySign", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start getSignPrize", user["name"])
		getSignPrize(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getSignPrize", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start getFreeDailyGiftBox", user["name"])
		getFreeDailyGiftBox(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getFreeDailyGiftBox", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start playLuckyWheel", user["name"])
		shareAPI(user["serverURL"], user["zoneToken"])
		playLuckyWheel(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end playLuckyWheel", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start getFreeClamp", user["name"])
		shareAPI(user["serverURL"], user["zoneToken"])
		getFreeClamp(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getFreeClamp", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start getInviteSnow", user["name"])
		shareAPI(user["serverURL"], user["zoneToken"])
		getInviteSnow(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getInviteSnow", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start autoFriendEnergy", user["name"])
		autoFriendEnergy(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end autoFriendEnergy", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] start getSixEnergy", user["name"])
		for i := 0; i < 6; i++ {
			time.Sleep(time.Millisecond * 200)
			getSixEnergy(user["serverURL"], user["zoneToken"])
		}
		log.Printf("[%v] end getSixEnergy", user["name"])
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		gameList := []string{"535", "525", "157", "452", "411"}

		for _, v := range gameList {
			getAward(user["token"], v)
		}
		j++
	}

	for _, user := range tokenList {
		if j >= 2 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		for i := 1; i <= 3; i++ {
			getFamilySignPrize(user["serverURL"], user["zoneToken"], i)
		}
		j++
	}

	// 为大佬海浪、铲子助力

	var times int64 = 1
	for _, user := range tokenList {
		if times > 6 {
			break
		}
		if times > 3 {
			log.Printf("[%v] 为蜜蜜海浪助力", user["name"])
			if user["uid"] != "309392050" {
				beachHelp(user["serverURL"], user["zoneToken"], "309392050", 42)
				times++
			}
		} else {
			log.Printf("[%v] 为牛海浪助力", user["name"])
			if user["uid"] != "302691822" {
				beachHelp(user["serverURL"], user["zoneToken"], "302691822", 42)
				times++
			}
		}

	}

	// 互相助力逻辑
	for _, user := range tokenList {
		var chanziTimes = 1

		var nn = map[string]string{"uid": "302691822", "name": "大佬"}
		helpInfo["302691822"] = 0
		var tokenList1 = tokenList
		tokenList1 = append(tokenList1, nn)
		for _, v := range tokenList1 {
			if chanziTimes < 6 {
				if v["uid"] != user["uid"] {
					helpTimes := helpInfo[v["uid"]]
					if helpTimes != 3 {
						// log.Printf("[%v]为[%v]海浪助力", user["name"], v["name"])
						// beachHelp(user["serverURL"], user["zoneToken"], v["uid"], 42)
						// time.Sleep(time.Second * 2)
						log.Printf("[%v]为[%v]铲子助力", user["name"], v["name"])
						beachHelp(user["serverURL"], user["zoneToken"], v["uid"], 43)
						helpInfo[v["uid"]] = helpTimes + 1
						chanziTimes++
						time.Sleep(time.Second * 1)
					}
				}
			}

		}
	}

	// for _, user := range tokenList {
	// 	var hailangTimes = 1
	// 	for _, v := range tokenList {
	// 		if hailangTimes < 4 {
	// 			if v["uid"] != user["uid"] {
	// 				helpTimes := helpInfo[v["uid"]]
	// 				if helpTimes != 6 {
	// 					log.Printf("[%v]为[%v]海浪助力", user["name"], v["name"])
	// 					beachHelp(user["serverURL"], user["zoneToken"], v["uid"], 42)
	// 					time.Sleep(time.Second * 2)
	// 					helpInfo[v["uid"]] = helpTimes + 1
	// 					hailangTimes++
	// 					time.Sleep(time.Second * 2)
	// 				}
	// 			}
	// 		}

	// 	}
	// }
	// time.Sleep(time.Second * 1)

	// 领取海浪、铲子

	for _, user := range tokenList {
		if j >= 2 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] 领取海浪", user["name"])

		getHelpItem(user["serverURL"], user["zoneToken"], 0, 0)

		time.Sleep(time.Second * 1)

		for i := 0; i <= 2; i++ {
			log.Printf("[%v] 领取铲子", user["name"])

			getHelpItem(user["serverURL"], user["zoneToken"], 1, i)
		}
		j++
	}

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] 公会聊天", user["name"])
		familyChat(user["serverURL"], user["zoneToken"])
		j++
	}

	for _, user := range tokenList {
		if j >= 8 {
			time.Sleep(time.Second * 1)
			j = 1
		}
		log.Printf("[%v] collectMineGold", user["name"])
		collectMineGold(user["serverURL"], user["zoneToken"])
	}

	getAwardForCowBoy()

	log.Println("familySignGo finish")

	othersSign()
}

func othersSign() {
	SQL := "select id, token, name from tokens where !find_in_set(id, (select conf_value from config where conf_key = 'animalUids')) and id != (select conf_value from config where conf_key = 'cowBoy')"
	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()

	var tokenList []map[string]string

	var helpInfo = make(map[string]int64)

	for rows.Next() {
		var uid, token, name string
		rows.Scan(&uid, &token, &name)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		if zoneToken == "" {
			sendMsg(uid + ":" + name)
			log.Printf("[ %v] token is invalid\n", uid)
		}
		tokenList = append(tokenList, map[string]string{"uid": uid, "name": name, "serverURL": serverURL, "zoneToken": zoneToken, "token": token})
		helpInfo["uid"] = 0

		log.Printf("[%v] start familySign", uid)
		familySign(serverURL, zoneToken)
		log.Printf("[%v] end familySign", uid)
		time.Sleep(time.Second * 1)
		log.Printf("[%v] start getSignPrize", name)
		getSignPrize(serverURL, zoneToken)
		log.Printf("[%v] end getSignPrize", name)
		time.Sleep(time.Second * 1)
		log.Printf("[%v] start getFreeDailyGiftBox", name)
		getFreeDailyGiftBox(serverURL, zoneToken)
		log.Printf("[%v] end getFreeDailyGiftBox", name)
		time.Sleep(time.Second * 1)

		log.Printf("[%v] start playLuckyWheel", name)
		shareAPI(serverURL, zoneToken)
		playLuckyWheel(serverURL, zoneToken)
		log.Printf("[%v] end playLuckyWheel", name)
		time.Sleep(time.Second * 1)

		log.Printf("[%v] start getFreeClamp", name)
		shareAPI(serverURL, zoneToken)
		getFreeClamp(serverURL, zoneToken)
		log.Printf("[%v] end getFreeClamp", name)
		time.Sleep(time.Second * 1)

		log.Printf("[%v] start getInviteSnow", name)
		shareAPI(serverURL, zoneToken)
		getInviteSnow(serverURL, zoneToken)
		log.Printf("[%v] end getInviteSnow", name)

		time.Sleep(time.Second * 1)
		log.Printf("[%v] start autoFriendEnergy", name)
		autoFriendEnergy(serverURL, zoneToken)
		log.Printf("[%v] end autoFriendEnergy", name)

		log.Printf("[%v] start getSixEnergy", name)
		for i := 0; i < 6; i++ {
			time.Sleep(time.Second * 1)
			getSixEnergy(serverURL, zoneToken)
		}
		log.Printf("[%v] end getSixEnergy", name)

		gameList := []string{"535", "525", "157", "452", "411"}

		time.Sleep(time.Second * 1)
		for _, v := range gameList {
			getAward(token, v)
		}
		log.Printf("[%v] 公会聊天", name)
		familyChat(serverURL, zoneToken)

		log.Printf("[%v] collectMineGold", name)
		collectMineGold(serverURL, zoneToken)

	}

	// 互相助力逻辑
	for _, user := range tokenList {
		var chanziTimes = 1
		for _, v := range tokenList {
			if chanziTimes < 6 {
				if v["uid"] != user["uid"] {
					helpTimes := helpInfo[v["uid"]]
					if helpTimes != 3 {
						// log.Printf("[%v]为[%v]海浪助力", user["name"], v["name"])
						// beachHelp(user["serverURL"], user["zoneToken"], v["uid"], 42)
						// time.Sleep(time.Second * 2)
						log.Printf("[%v]为[%v]铲子助力", user["name"], v["name"])
						beachHelp(user["serverURL"], user["zoneToken"], v["uid"], 43)
						helpInfo[v["uid"]] = helpTimes + 1
						chanziTimes++
						time.Sleep(time.Second * 1)
					}
				}
			}

		}
	}

	for _, user := range tokenList {
		for i := 1; i <= 3; i++ {
			getFamilySignPrize(user["serverURL"], user["zoneToken"], i)
		}
	}

	// time.Sleep(time.Second * 1)

	// for _, user := range tokenList {
	// 	var hailangTimes = 1
	// 	for _, v := range tokenList {
	// 		if hailangTimes < 4 {
	// 			if v["uid"] != user["uid"] {
	// 				helpTimes := helpInfo[v["uid"]]
	// 				if helpTimes != 6 {
	// 					log.Printf("[%v]为[%v]海浪助力", user["name"], v["name"])
	// 					beachHelp(user["serverURL"], user["zoneToken"], v["uid"], 42)
	// 					time.Sleep(time.Second * 2)
	// 					helpInfo[v["uid"]] = helpTimes + 1
	// 					hailangTimes++
	// 					time.Sleep(time.Second * 2)
	// 				}
	// 			}
	// 		}

	// 	}s
	// }
	// time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		log.Printf("[%v] 领取海浪", user["name"])

		getHelpItem(user["serverURL"], user["zoneToken"], 0, 0)

		time.Sleep(time.Second * 1)

		for i := 0; i <= 2; i++ {
			log.Printf("[%v] 领取铲子", user["name"])

			getHelpItem(user["serverURL"], user["zoneToken"], 1, i)
		}
	}

}

func getAwardForCowBoy() {
	SQL := "select token from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
	var token string
	Pool.QueryRow(SQL).Scan(&token)
	gameList := []string{"535", "525", "157", "452", "411"}
	for _, v := range gameList {
		time.Sleep(time.Second * 1)
		getAward(token, v)
	}
}

func openSteamBoxGo() {
	SQL := "select id, name, token from tokens"
	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var uid, name, token string
		rows.Scan(&uid, &name, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		openSteamBox(serverURL, zoneToken, uid)
		log.Printf("[%v]领取汤圆成功", name)
	}
}

// 抓特定箱子 quality-> 1=普通 2=稀有 3=传奇
func getBoxPrizeGo(uid string, quality, amount float64) (total, all float64) {

	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, uid).Scan(&token)

	serverURL, zoneToken := getSeverURLAndZoneToken(token)
	helpList := getGoldMineHelpList(serverURL, zoneToken, quality)

	ids := []int64{21, 22}

	for _, v := range helpList {

		for _, v2 := range ids {
			if total == amount {
				log.Println("抓取完毕")
				rate, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (total/all)*100), 64)
				log.Println("命中率:", fmt.Sprintf("%v", rate)+"%")
				return
			}
			time.Sleep(time.Second * 1)
			getFlag := goldMineFish(serverURL, zoneToken, v["uid"], v2)
			if getFlag != -1 {
				all++
			}
			if getFlag == 1 {
				total += 1
				log.Printf("当前数量:%v 目标数量:%v", total, amount)
				break
			}
		}
	}

	rate, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", (total/all)*100), 64)
	log.Println("命中率:", fmt.Sprintf("%v", rate)+"%")

	return
}

// interface functions
func getServerURL() (serverURL string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := "https://api.11h5.com/conf?cmd=getGameInfo&gameid=147&" + now
	formData := httpGetReturnJson(url)

	ext, ok := formData["ext"].(map[string]interface{})

	if !ok {
		log.Println("get serverURL err")
		return
	}

	serverURL, ok = ext["serverURL"].(string)
	if !ok {
		log.Println("get serverURL err")
		return
	}
	return

}

func getEnterInfo(uid, name, serverURL, token, key string) (zoneToken string, info interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := serverURL + "/zone?cmd=enter&token=" + token + "&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now=" + now
	formData := httpGetReturnJson(URL)
	_, ok := formData["zoneToken"].(string)
	if !ok {
		log.Printf("[%v] token is invaild", name)
		sendMsg(uid + ":" + name)
		return
	}

	zoneToken, ok = formData["zoneToken"].(string)
	if !ok {
		log.Printf("token:%v get zoneToken err", token)
		return
	}

	info, ok = formData[key]
	if !ok {
		return
	}

	Pool.Exec("update tokens set serverURL = ?, zoneToken = ? where token = ?", serverURL, zoneToken, token)

	return
}

func checkZoneToken(serverURL, zoneToken string) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=hasNewMail&token=%v&now=%v", serverURL, zoneToken, now)
	// log.Println("checkZoneToken URL is :", URL)
	formData := httpGetReturnJson(URL)
	_, ok := formData["newMailCnt"].(float64)
	if ok {
		return true
	}
	log.Println("zoneToken is err :", formData["error"])
	return false
}

func getSeverURLAndZoneToken(token string) (serverURL, zoneToken string) {

	SQL := "select serverURL, zoneToken from tokens where token = ?"
	Pool.QueryRow(SQL, token).Scan(&serverURL, &zoneToken)
	if checkZoneToken(serverURL, zoneToken) {
		return
	}
	serverURL = getServerURL()
	zoneToken = getZoneToken(serverURL, token)
	if zoneToken != "" {
		SQL = "update tokens set serverURL = ?, zoneToken = ? where token = ?"
		Pool.Exec(SQL, serverURL, zoneToken, token)
		return
	}
	return
}

func getZoneToken(serverURL, token string) (zoneToken string) {

	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := serverURL + "/zone?cmd=enter&token=" + token + "&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now=" + now
	formData := httpGetReturnJson(URL)
	zoneToken, ok := formData["zoneToken"].(string)
	if !ok {
		log.Printf("token:%v get zoneToken err", token)
		Pool.Exec("update tokens set serverURL = '', zoneToken = '', token = '' where token = ?", token)
		return
	}

	buildPrice, ok := formData["buildPrice"].(map[string]interface{})

	if ok {
		ids := []string{"1", "2", "3", "4", "5"}
		var sum float64 = 0
		for _, id := range ids {
			price, ok := buildPrice[id].(map[string]interface{})
			if ok {
				for _, v := range ids {
					money, ok := price[v].(float64)
					if ok {
						sum += money
					}
				}
			}
		}
		sum, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", sum/1000000), 64)

		gold, _ := formData["gold"].(float64)
		gold, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", gold/1000000), 64)

		nickname, _ := formData["nickname"].(string)
		nickname, _ = url.QueryUnescape(nickname)

		if ok {
			log.Printf("%v 当前金币%vM 过岛费用%vM %v\n", nickname, gold, sum, gold >= sum)
			pieceList, ok := formData["pieceList"].(map[string]interface{})
			if ok {
				ids = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

				for _, v := range ids {
					piece, ok := pieceList[v].(map[string]interface{})
					if ok {
						log.Printf("拼图【%v】 数量：%v", v, piece["count"])
					}
				}

			}

		}

	}
	return
}

func getZoneToken_1(serverURL, token string) (zoneToken string, flag bool) {

	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := serverURL + "/zone?cmd=enter&token=" + token + "&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now=" + now
	formData := httpGetReturnJson(URL)
	zoneToken, ok := formData["zoneToken"].(string)
	if !ok {
		log.Printf("token:%v get zoneToken err", token)
		Pool.Exec("update tokens set serverURL = '', zoneToken = '', token = '' where token = ?", token)
		return
	}

	buildPrice, ok := formData["buildPrice"].(map[string]interface{})

	if ok {
		ids := []string{"1", "2", "3", "4", "5"}
		var sum float64 = 0
		for _, id := range ids {
			price, ok := buildPrice[id].(map[string]interface{})
			if ok {
				for _, v := range ids {
					money, ok := price[v].(float64)
					if ok {
						sum += money
					}
				}
			}
		}
		sum, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", sum/1000000), 64)

		gold, _ := formData["gold"].(float64)
		gold, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", gold/1000000), 64)

		nickname, _ := formData["nickname"].(string)
		nickname, _ = url.QueryUnescape(nickname)

		if ok {
			flag = gold >= sum
			log.Printf("%v 当前金币%vM 过岛费用%vM %v\n", nickname, gold, sum, flag)
			pieceList, ok := formData["pieceList"].(map[string]interface{})
			if ok {
				ids = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}

				for _, v := range ids {
					piece, ok := pieceList[v].(map[string]interface{})
					if ok {
						log.Printf("拼图【%v】 数量：%v", v, piece["count"])
					}
				}

			}

		}

	}
	return
}

func getNickName(serverURL, token string) (zoneToken, nickname string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := serverURL + "/zone?cmd=enter&token=" + token + "&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now=" + now
	formData := httpGetReturnJson(URL)
	zoneToken, ok := formData["zoneToken"].(string)
	if !ok {
		log.Printf("token:%v get zoneToken err", token)
		return
	}
	nickname, _ = formData["nickname"].(string)
	nickname, _ = url.QueryUnescape(nickname)
	return
}

func enterFamilyRob(serverURL, zoneToken string) (foods []map[string]interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := serverURL + "/game?cmd=enterFamilyRob&token=" + zoneToken + "&now=" + now
	formData := httpGetReturnJson(url)
	if formData == nil {
		return
	}
	familyRob, ok := formData["familyRob"].(map[string]interface{})
	if !ok {
		return
	}
	worker, ok := familyRob["worker"].(string)
	if !ok {
		return
	}
	if worker != "" {
		log.Println("worker:", worker)
		return
	}
	foodList, ok := familyRob["foodList"].([]interface{})
	if !ok {
		return
	}
	for _, v := range foodList {
		vv, ok := v.(map[string]interface{})
		if !ok {
			break
		}

		familyId, ok := vv["familyId"].(float64)
		if ok {
			if familyId != 1945 {
				searchFamily(serverURL, zoneToken, familyId)
			}
		}

		food := make(map[string]interface{})
		teamLen := len(vv["myTeam"].(map[string]interface{})["robList"].([]interface{}))
		food["id"] = vv["id"].(string)
		food["myTeam"] = teamLen
		food["row"] = vv["row"].(float64)
		foods = append(foods, food)
	}

	return
}

func getFamilyId(serverURL, zoneToken string) (familyId float64, timeFlushList []string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := serverURL + "/game?cmd=enterFamilyRob&token=" + zoneToken + "&now=" + now
	formData := httpGetReturnJson(url)
	if formData == nil {
		return
	}
	familyRob, ok := formData["familyRob"].(map[string]interface{})
	if !ok {
		return
	}
	foodList, ok := familyRob["foodList"].([]interface{})
	if !ok {
		return
	}

	log.Println(foodList)

	myFlushTimeList, ok := familyRob["myFlushTimeList"].([]interface{})
	log.Println(myFlushTimeList)

	if !ok {
		return
	}

	for _, v := range myFlushTimeList {
		log.Println(v)

		vv, ok := v.(string)
		if !ok {
			return
		}
		timeInt, err := strconv.Atoi(vv)
		log.Println(timeInt)

		if err != nil {
			return
		}

		date := time.Unix(int64(timeInt/1000), 0).Format("15:04")
		timeFlushList = append(timeFlushList, date)
	}

	for _, v := range foodList {
		vv, ok := v.(map[string]interface{})
		if !ok {
			break
		}

		familyId, ok = vv["familyId"].(float64)
		log.Println(familyId)
		if ok {
			if familyId != 1945 {
				return
			}
		}

	}

	return
}

func robFamilyFood(serverURL, zoneToken, foodId string) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := serverURL + "/game?cmd=robFamilyFood&token=" + zoneToken + "&foodId=" + foodId + "&now=" + now
	formData := httpGetReturnJson(url)
	if _, ok := formData["error"]; ok {
		return false
	}
	return true

}

// 查看宝箱帮助列表

func getGoldMineHelpList(serverURL, zoneToken string, quality float64) (data []map[string]interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getGoldMineHelpList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)

	helpList, ok := formData["helpList"].([]interface{})

	if !ok {
		return
	}

	for _, v := range helpList {
		vv, ok := v.(map[string]interface{})
		if !ok {
			break
		}
		if vv["quality"].(float64) == quality {
			data = append(data, vv)
		}
	}

	return
}

// 进入抓宝箱
// func enterGoldMine(serverURL, zoneToken string, fuid interface{}) {
// 	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)

// 	url := fmt.Sprintf("%v/game?cmd=enterGoldMine&token=%v&fuid=%v&type=0&now=%v", serverURL, zoneToken, fuid, now)
// 	// url := serverURL + "/game?cmd=enterGoldMine&token=" + zoneToken + "&fuid=" + fuid + "&type=0&now=" + now
// 	formData := httpGetReturnJson(url)
// 	log.Println("formData:", formData)
// }

// 抓宝箱
func goldMineFish(serverURL, zoneToken string, fuid, id interface{}) int64 {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)

	url := fmt.Sprintf("%v/game?cmd=goldMineFish&token=%v&fuid=%v&id=%v&now=%v", serverURL, zoneToken, fuid, id, now)
	formData := httpGetReturnJson(url)
	getItem, ok := formData["getItem"].(map[string]interface{})

	if !ok {
		return -1
	}

	_, ok = getItem["25"]
	if ok {
		log.Println("抓到黄水晶")
		return 1
	}

	_, ok = getItem["26"]
	if ok {
		log.Println("抓到紫水晶")
		return 1
	}

	_, ok = getItem["27"]
	if ok {
		log.Println("抓到黑水晶")
		return 1
	}

	_, ok = getItem["28"]
	if ok {
		log.Println("抓到绿宝石")
		return 1
	}

	_, ok = getItem["29"]
	if ok {
		log.Println("抓到红宝石")
		return 1
	}

	_, ok = getItem["30"]
	if ok {
		log.Println("抓到钻石")
		return 1
	}
	log.Println("没抓到宝石")
	return 0

}

// 获取帮助列表
func getBossHelpList(serverURL, zoneToken string) (bossList []map[string]interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getBossHelpList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)

	bossHelpList, ok := formData["bossHelpList"].([]interface{})
	if !ok {
		return
	}

	for _, v := range bossHelpList {
		boss := v.(map[string]interface{})["boss"].(map[string]interface{})
		leftHp := boss["leftHp"].(float64)

		if leftHp > 0 {
			bossList = append(bossList, boss)
		}
	}
	log.Println("bossList:", bossList)
	return

}

//
// 小号获取BOSS状态
func enterBoss(serverURL, zoneToken, bossID string, uid float64) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=enterBoss&token=%v&bossID=%v&now=%v", serverURL, zoneToken, bossID, now)
	formData := httpGetReturnJson(url)

	boss, ok := formData["boss"].(map[string]interface{})
	if !ok {
		return false
	}
	leftHp, ok := boss["leftHp"].(float64)

	if !ok {
		return false
	}
	if leftHp > 1000 {
		bossAttackList := formData["bossAttackList"].([]interface{})
		for _, v := range bossAttackList {
			vv := v.(map[string]interface{})
			if vv["uid"].(float64) == uid && vv["damage"].(float64) > 0 {
				return false
			}
		}
		return true
	}
	return false

}

// 开启Boss
func summonBoss(serverURL, zoneToken string) string {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=summonBoss&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	boss, ok := formData["boss"].(map[string]interface{})
	if !ok {
		now = fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		url = fmt.Sprintf("%v/game?cmd=getMyBoss&token=%v&now=%v", serverURL, zoneToken, now)
		formData = httpGetReturnJson(url)
		myBoss, ok := formData["boss"].(map[string]interface{})
		if !ok {
			return ""
		}
		bossID, ok := myBoss["id"].(string)
		if !ok {
			return ""
		}
		return bossID
	}

	bossID, ok := boss["id"].(string)
	if !ok {
		return ""
	}
	return bossID
}

// 邀请BOSS
func inviteBoss(serverURL, zoneToken, bossID string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=inviteBoss&token=%v&bossID=%v&fuidList=[301807377,302691822,309433834,326941142,374289806,381909995,406378614,690708340,693419844,694068717,694981971]&now=%v", serverURL, zoneToken, bossID, now)
	httpGetReturnJson(url)
}

// 打自己的BOSS 到600血

func attackMyBoss(serverURL, zoneToken, bossID, mode string) {

	// Pool.QueryRow("select conf_value from config where conf_key = 'attackBossMode'").Scan()

	if mode == "4500" {
		// 300

		attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 100, 100)

		// 1500

		attackBossAPI(serverURL, zoneToken, bossID, 6, 0, 1, 200, 200)

		// 1800

		attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 100, 100)

		// 3000
		attackBossAPI(serverURL, zoneToken, bossID, 6, 0, 1, 200, 200)

		// 3300
		attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 100, 100)

		// 4100

		attackBossAPI(serverURL, zoneToken, bossID, 4, 0, 1, 200, 200)

		// 4500
		attackBossAPI(serverURL, zoneToken, bossID, 1, 1, 1, 400, 400)

		return
	}

	// 300

	attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 95, 100)

	// 1500

	attackBossAPI(serverURL, zoneToken, bossID, 6, 0, 1, 195, 200)

	// 1800

	attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 95, 100)

	// 3000
	attackBossAPI(serverURL, zoneToken, bossID, 6, 0, 1, 195, 200)

	// 3300
	attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 95, 100)

	// 4100

	attackBossAPI(serverURL, zoneToken, bossID, 5, 0, 1, 195, 200)

	// 4500
	// attackBossAPI(serverURL, zoneToken, bossID, 1, 400, 1, 1)

}

// RangeRand(390, 400)

func attackBossAPI(serverURL, zoneToken, bossID string, amount, isPerfect, isDouble int, min, max int64) {
	for i := 1; i <= amount; i++ {
		now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		damage := RangeRand(min, max)
		url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=%v&isDouble=%v&now=%v", serverURL, zoneToken, bossID, damage, isPerfect, isDouble, now)
		formData := httpGetReturnJson(url)
		boss, ok := formData["boss"].(map[string]interface{})
		if !ok {
			time.Sleep(time.Second * 3)
			httpGetReturnJson(url)
			time.Sleep(time.Second * 3)
		} else {
			leftHp, ok := boss["leftHp"]
			if !ok {
				return
			}
			log.Println("leftHp:", leftHp)
			time.Sleep(time.Second * 3)
		}
	}
}

// 小号打Boss
func attackBoss(serverURL, zoneToken, bossID string) {
	// 1次 50
	attackBossAPI(serverURL, zoneToken, bossID, 1, 0, 0, 50, 50)
	// 3 次 100
	attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 100, 100)
	// 1 次 200
	attackBossAPI(serverURL, zoneToken, bossID, 1, 0, 1, 200, 200)
}

//
func getAttackEnemyList(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getAttackEnemyList&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 攻击小岛
func attackIsland(serverURL, zoneToken string, flag, targetUid, building interface{}) (attackUid int, buildings string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=attack&token=%v&type=%v&targetUid=%v&building=%v&now=%v", serverURL, zoneToken, flag, targetUid, building, now)
	formData := httpGetReturnJson(url)
	log.Printf("攻击小岛 目标uid:%v, 建筑:%v, 增加金币:%v, 额外金币:%v", targetUid, building, formData["addGold"], formData["companionGold"])

	attackData, ok := formData["attackData"].(map[string]interface{})

	if ok {

		island, _ := attackData["island"].(map[string]interface{})

		if island["1"].(map[string]interface{})["lv"].(float64) != 0 {
			buildings = "1"
		} else if island["2"].(map[string]interface{})["lv"].(float64) != 0 {
			buildings = "2"
		} else if island["3"].(map[string]interface{})["lv"].(float64) != 0 {
			buildings = "3"
		} else if island["4"].(map[string]interface{})["lv"].(float64) != 0 {
			buildings = "4"
		} else {
			buildings = "5"
		}
		attackUid = int(attackData["uid"].(float64))
		return
	}
	return
}

// 公会签到

func familySign(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=familySign&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 领取公会签到奖励
func getFamilySignPrize(serverURL, zoneToken string, id int) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getFamilySignPrize&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

// 领取签到奖励
func getSignPrize(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=sign&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 领取免费福利礼包
func getFreeDailyGiftBox(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getFreeDailyGiftBox&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 分享两部曲
func shareAPI(serverURL, zoneToken string) {
	addDayTaskShareCnt(serverURL, zoneToken)
	getSharePrize(serverURL, zoneToken)

}

// 获取6次能量
func getSixEnergy(serverURL, zoneToken string) {
	addDayTaskShareCnt(serverURL, zoneToken)
	getSharePrize0(serverURL, zoneToken)
}

// 分享调用接口1
func addDayTaskShareCnt(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=addDayTaskShareCnt&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 分享调用接口2
func getSharePrize(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getSharePrize&token=%v&noEnergy=1&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

func getSharePrize0(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getSharePrize&token=%v&noEnergy=0&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 免费转盘
func playLuckyWheel(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=playLuckyWheel&token=%v&ad=1&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 领取免费夹子
func getFreeClamp(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getFreeClamp&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 领取免费糖果炮弹
func getInviteSnow(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getInviteSnow&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

func searchFamily(serverURL, zoneToken string, id float64) (name string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=searchFamily&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	formData := httpGetReturnJson(URL)
	recommendFamilyInfo, ok := formData["recommendFamilyInfo"].(map[string]interface{})
	if ok {
		name, ok = recommendFamilyInfo["name"].(string)
		if ok {
			log.Printf("对方公会ID:%v, 公会名称:%v", id, name)
			return
		}
	}
	return
}

// 一键赠送和领取好友能量
func autoFriendEnergy(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=autoFriendEnergy&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 领取其他链接的奖励
func getAward(token, gameID string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("https://libao.11h5.com/wall?cmd=getAward&token=%v&gameid=%v&type=1&channel=147&now=%v", token, gameID, now)
	httpGetReturnJson(url)
}

// 查看好友糖果树
func getFriendsCandyTreeInfo(serverURL, zoneToken string, targetType float64) (uids []int) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getFriendsCandyTreeInfo&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	arr, ok := formData["arr"].([]interface{})
	if !ok {
		return
	}

	for _, v := range arr {
		vv, ok := v.(map[string]interface{})
		if !ok {
			break
		}
		if vv["quality"].(float64) == targetType {
			uids = append(uids, int(vv["uid"].(float64)))
		}
	}
	return
}

// 进入糖果树
func getCandyTreeInfo(serverURL, zoneToken string, opUid interface{}) (posList []float64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)

	url := fmt.Sprintf("%v/game?cmd=getCandyTreeInfo&token=%v&opUid=%v&now=%v", serverURL, zoneToken, opUid, now)
	formData := httpGetReturnJson(url)

	candyTree, ok := formData["candyTree"].(map[string]interface{})

	if !ok {
		return
	}

	candyBoxes, ok := candyTree["candyBoxes"].([]interface{})

	if !ok {
		return
	}

	for _, v := range candyBoxes {
		vv, ok := v.(map[string]interface{})

		if !ok {
			break
		}
		posList = append(posList, vv["pos"].(float64))
	}

	return
}

// 打糖果
func hitCandyTree(serverURL, zoneToken string, opUid, pos interface{}) string {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=hitCandyTree&token=%v&opUid=%v&pos=%v&now=%v", serverURL, zoneToken, opUid, pos, now)
	formData := httpGetReturnJson(url)

	getItem, ok := formData["getItem"]
	if !ok {
		log.Printf("hitCandyTree opUid:%v, pos:%v, formData:%v\n", opUid, pos, formData)
		return "err"
	}

	v, ok := getItem.(map[string]interface{})

	if !ok {
		log.Printf("hitCandyTree opUid:%v, pos:%v, formData:%v\n", opUid, pos, formData)
		return "err"
	}

	log.Println("hitCandyTree getItem:", v)
	_, ok = v["22"]
	if !ok {
		_, ok = v["21"]
		if !ok {
			_, ok = v["20"]
			if !ok {
				_, ok = v["19"]
				if !ok {
					_, ok = v["18"]
					if !ok {
						_, ok = v["17"]
						if !ok {
							return "false"
						}
					}
				}

			}
		}
	}

	return "true"
}

// 获取免费龙炮弹
func getFreeBossCannon(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getFreeBossCannon&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 赠送拼图
func giftPiece(serverURL, zoneToken, id, targetUid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=giftPiece&token=%v&id=%v&targetUid=%v&now=%v", serverURL, zoneToken, id, targetUid, now)
	httpGetReturnJson(url)
}

// 扔骰子
func throwDice(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=throwDice&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	getItem, ok := formData["getItem"]
	if ok {
		log.Println(getItem)
	}

}

func exchangeXmas(serverURL, zoneToken string, id int64) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=exchangeXmas&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	formData := httpGetReturnJson(URL)
	if _, ok := formData["error"]; ok {
		return false
	}
	return true
}

func exchangeBeachReward(serverURL, zoneToken string, id int64) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=exchangeBeachReward&token=%v&index=%v&now=%v", serverURL, zoneToken, id, now)
	formData := httpGetReturnJson(URL)
	if _, ok := formData["error"]; ok {
		return false
	}
	return true
}

// https://s147.11h5.com:3147/118_89_198_11/3147/

// 获取邮件列表
func getMailList(serverURL, zoneToken, title string) (mailids []string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getMailList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	mailList, ok := formData["mailList"].([]interface{})

	if !ok {
		return
	}

	for _, v := range mailList {
		if title != "" {
			if v.(map[string]interface{})["title"].(string) == title {
				mailid := v.(map[string]interface{})["mailid"].(string)
				mailids = append(mailids, mailid)
			}
		} else {
			mailid := v.(map[string]interface{})["mailid"].(string)
			mailids = append(mailids, mailid)
		}
	}
	return
}

// 领取邮件奖励
func getMailAttachments(serverURL, zoneToken, mailid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getMailAttachments&token=%v&mailid=%v&now=%v", serverURL, zoneToken, mailid, now)
	httpGetReturnJson(url)
}

// 删除邮件
func deleteMail(serverURL, zoneToken, mailid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=deleteMail&token=%v&mailid=%v&now=%v", serverURL, zoneToken, mailid, now)
	formData := httpGetReturnJson(URL)
	log.Println("deleteMail result:", formData)
}

func readMail(serverURL, zoneToken, mailid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=readMail&token=%v&mailid=%v&now=%v", serverURL, zoneToken, mailid, now)
	httpGetReturnJson(URL)
}

// https://s147.11h5.com:3148/111_231_17_85/3149/

func getBossPrizeList(serverURL, zoneToken string) (bossIDList []string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getBossPrizeList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)

	bossPrizeList, ok := formData["bossPrizeList"].([]interface{})
	if !ok {
		return
	}

	for _, v := range bossPrizeList {
		bossID := v.(map[string]interface{})["bossID"].(string)
		bossIDList = append(bossIDList, bossID)
	}
	return
}

func getBossPrize(serverURL, zoneToken, bossID string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getBossPrize&token=%v&bossID=%v&now=%v", serverURL, zoneToken, bossID, now)
	httpGetReturnJson(url)
}

// 发起拼图
func setPiece(serverURL, zoneToken string, id int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=setPiece&token=%v&id=%v&type=0&common=0&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

// 领取拼图奖励
func getPiecePrize(serverURL, zoneToken string) (ok bool) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getPiecePrize&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(URL)

	_, ok = formData["getItem"]
	return
}

func getFreeEnergy(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getFreeEnergy&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

func draw(uid, userName, serverURL, zoneToken string, drawMulti interface{}) float64 {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=draw&token=%v&drawMulti=%v&now=%v", serverURL, zoneToken, drawMulti, now)
	formData := httpGetReturnJson(url)

	idx, ok := formData["id"]

	if !ok {
		log.Printf("[%v]没有能量可摇！:%v", userName, formData)

		shareAPI(serverURL, zoneToken)
		getShareDrawDice(serverURL, zoneToken, 75)
		getShareDrawDice(serverURL, zoneToken, 97)
		time.Sleep(time.Second * 1)
		shareAPI(serverURL, zoneToken)
		getDayShareGold(serverURL, zoneToken)
		return 0
	}

	id := fmt.Sprintf("%v", idx)

	if id == "10" {

		stealData, ok := formData["stealData"].([]interface{})

		// var idx = 1
		idx := 1

		if ok {
			for i, v := range stealData {
				vv, ok := v.(map[string]interface{})
				if ok {
					rich, ok := vv["rich"].(float64)
					if ok {
						if rich == 1 {
							log.Printf("[%v]【摇一摇】rich is :%v", userName, i)
							idx = i
							break
						}
					}
				}
			}
		}
		// log.Println("stealData:", stealData)
		time.Sleep(time.Second * 1)
		stealResult, idx1 := steal(serverURL, zoneToken, idx)
		time.Sleep(time.Millisecond * 100)
		log.Printf("[%v]【摇一摇】偷取 结果:%v", userName, stealResult)

		for i := 0; i < 5; i++ {
			log.Println("start check steal idx:", idx1)
			if idx1 == -1 {
				log.Printf("[%v]【摇一摇】停止偷取", userName)
				break
			} else {
				time.Sleep(time.Millisecond * 200)
				stealResult1, idx2 := steal(serverURL, zoneToken, idx1)
				log.Printf("[%v]【摇一摇】偷取 结果:%v 下一个目标:%v", userName, stealResult1, idx2)
				idx1 = idx2
			}
		}

	} else if id == "3" {
		followCompanion(serverURL, zoneToken, 4)
		if uid == "302691822" || uid == "309392050" {
			followCompanion_1(serverURL, zoneToken, 4)
			time.Sleep(time.Second * 1)

			confKey := "attackIslandUid2"
			if uid == "302691822" {
				confKey = "attackIslandUid"
			}

			var attackIslandUid int64
			var attackIslandName string
			SQL := fmt.Sprintf("select conf_value, (select name from tokens where id = conf_value) as name from config where conf_key = '%s'", confKey)
			Pool.QueryRow(SQL).Scan(&attackIslandUid, &attackIslandName)
			log.Printf("[%v]【摇一摇】攻击好友【%v】", userName, attackIslandName)
			attuid, builds := attackIsland(serverURL, zoneToken, 1, attackIslandUid, 1)
			rebuild(attackIslandUid, 1)

			for {
				if attuid == 0 || builds == "" {
					log.Printf("[%v]【摇一摇】停止攻击", userName)
					break
				} else {
					time.Sleep(time.Millisecond * 200)
					attuid, builds = attackIsland(serverURL, zoneToken, 1, attackIslandUid, 1)
					rebuild(attackIslandUid, 1)
				}
			}
			followCompanion_1(serverURL, zoneToken, 2)

		} else {
			attackData, _ := formData["attackData"].(map[string]interface{})

			island, _ := attackData["island"].(map[string]interface{})

			var building string

			if island["1"].(map[string]interface{})["lv"].(float64) != 0 {
				building = "1"
			} else if island["2"].(map[string]interface{})["lv"].(float64) != 0 {
				building = "2"
			} else if island["3"].(map[string]interface{})["lv"].(float64) != 0 {
				building = "3"
			} else if island["4"].(map[string]interface{})["lv"].(float64) != 0 {
				building = "4"
			} else {
				building = "5"
			}
			attackUid := int(attackData["uid"].(float64))
			getAttackEnemyList(serverURL, zoneToken)
			time.Sleep(time.Second * 1)

			log.Printf("[%v]【摇一摇】攻击", userName)
			attackUid1, building1 := attackIsland(serverURL, zoneToken, 0, attackUid, building)

			for {
				if attackUid1 == 0 || building1 == "" {
					log.Printf("[%v]【摇一摇】停止攻击", userName)
					break
				} else {
					time.Sleep(time.Millisecond * 200)
					log.Printf("[%v]【摇一摇】继续攻击:%v", userName, attackUid1)
					attackUid1, building1 = attackIsland(serverURL, zoneToken, 0, attackUid1, building1)
				}
			}

		}
		followCompanion(serverURL, zoneToken, 2)

	}

	if formData["getRichmanDice"].(float64) == 1 {
		time.Sleep(time.Second * 1)
		shareAPI(serverURL, zoneToken)
		getShareDrawDice(serverURL, zoneToken, 75)
		log.Printf("[%v]【摇一摇】分享获取乐园骰子", userName)
	}

	if formData["getSnowball"].(float64) == 1 {
		log.Printf("[%v]【摇一摇】获取糖果炮弹 当前数量:%v", userName, formData["snowball"])
	}

	if formData["getClamp"].(float64) == 1 {
		log.Printf("[%v]【摇一摇】获取夹子 当前数量:%v", userName, formData["clamp"])
	}

	if formData["getShovel"].(float64) == 1 {
		log.Printf("[%v]【摇一摇】获得沙滩铲", userName)
	}

	if formData["getFirewood"].(float64) == 1 {
		time.Sleep(time.Second * 1)
		shareAPI(serverURL, zoneToken)
		getShareDrawDice(serverURL, zoneToken, 97)
		log.Printf("[%v]【摇一摇】分享获取柴火", userName)
	}

	if id == "5" {
		log.Printf("[%v]【摇一摇】获得盾", userName)
	}

	if formData["shareMulti"].(float64) != 0 {
		time.Sleep(time.Second * 1)
		shareAPI(serverURL, zoneToken)
		getDayShareGold(serverURL, zoneToken)
		log.Printf("[%v]【摇一摇】分享获取金币 倍数:%v", userName, formData["shareMulti"])
	}

	if getMiningItemId, ok := formData["getMiningItemId"].(float64); ok && getMiningItemId != 0 {
		if getMiningItemId == 184 {
			log.Printf("[%v]【摇一摇】获得鱼叉", userName)
		}
		if getMiningItemId == 185 {
			log.Printf("[%v]【摇一摇】获得鱼雷", userName)
		}
		if getMiningItemId == 186 {
			log.Printf("[%v]【摇一摇】获得水雷", userName)
		}
	}

	// getLabaCount: 2
	// getLabaId: 161 樱桃 159 西瓜

	getLabaId, ok := formData["getLabaId"].(float64)

	if ok {
		if getLabaId != 0 {

			shareAPI(serverURL, zoneToken)
			getShareDrawDice(serverURL, zoneToken, int64(getLabaId))

			var labaname string
			if getLabaId == 161 {
				labaname = "樱桃"
			}

			if getLabaId == 159 {
				labaname = "西瓜"
			}

			if getLabaId == 157 {
				labaname = "苹果"
			}

			if getLabaId == 160 {
				labaname = "草莓"
			}

			if getLabaId == 158 {
				labaname = "香蕉"
			}

			log.Printf("[%v]【摇一摇】获得[%s] 数量:[%v]", userName, labaname, formData["getLabaCount"])
		}
	}

	gold, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", formData["gold"].(float64)/1000000), 64)

	energy := formData["energy"].(float64)
	count, ok := formData["count"].(float64)

	if ok {
		count, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", count/1000000), 64)
	}

	if energy == 5 {
		log.Printf("[%v]【摇一摇】获取免费20能量", userName)
		getFreeEnergy(serverURL, zoneToken)
		log.Printf("[%v]【摇一摇】领取好友能量", userName)
		autoFriendEnergy(serverURL, zoneToken)

		taskIDs := getDayTasksInfo(serverURL, zoneToken)
		log.Printf("[%v]领取日常任务奖励:%v", userName, taskIDs)
		for _, taskID := range taskIDs {
			// time.Sleep(time.Millisecond * 100)
			getDayTaskAward(serverURL, zoneToken, taskID)
		}

		log.Printf("[%v]领取超值返利", userName)
		getElevenEnergyPrize(serverURL, zoneToken, 1)
		getElevenEnergyPrize(serverURL, zoneToken, 2)
		getElevenEnergyPrize(serverURL, zoneToken, 3)
		getElevenEnergyPrize(serverURL, zoneToken, 4)
	}

	log.Printf("[%v]转盘行为:%v, 剩余能量:%v, 当前金币:%vM, 增加金币:%vM", userName, id, energy, gold, count)
	return energy
}

func rebuild(uid, building int64) {
	SQL := "select token from tokens where id = ?"
	var token string
	Pool.QueryRow(SQL, uid).Scan(&token)
	serverURL, zoneToken := getSeverURLAndZoneToken(token)
	buildFix(serverURL, zoneToken, building)
	// buildUp(serverURL, zoneToken, building)
}

func buildUp(serverURL, zoneToken string, id int64) (islandid float64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=buildUp&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	formData := httpGetReturnJson(url)
	island, ok := formData["island"].(map[string]interface{})
	if ok {
		islandid = island["id"].(float64)
		return
	}
	return
}

// 领取过岛奖励
func getIslandPrize(serverURL, zoneToken string, id float64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getIslandPrize&token=%vid=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(URL)
}

// 过岛领取分享能量
func getIslandEnergy(serverURL, zoneToken string) {
	addDayTaskShareCnt(serverURL, zoneToken)
	getSharePrize(serverURL, zoneToken)

	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getIslandEnergy&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(URL)
}

// 修复岛
func buildFix(serverURL, zoneToken string, id int64) {

	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=buildFix&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

// 1=大队长 2=偷钱 4=打建筑 3=防偷钱
func followCompanion(serverURL, zoneToken string, id int64) {
	// now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	// url := fmt.Sprintf("%v/game?cmd=followCompanion&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	// httpGetReturnJson(url)
}

// 1=大队长 2=偷钱 4=打建筑 3=防偷钱
func followCompanion_1(serverURL, zoneToken string, id int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=followCompanion&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

// 【摇一摇】--分享获取金币
func getDayShareGold(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getDayShareGold&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 【摇一摇】--偷取
func steal(serverURL, zoneToken string, idx interface{}) (result bool, idx1 int) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=steal&token=%v&idx=%v&now=%v", serverURL, zoneToken, idx, now)
	formData := httpGetReturnJson(url)
	addGold, ok := formData["addGold"].(float64)
	if !ok {
		result = false
		idx1 = -1
		return
	}
	stealResult := formData["stealResult"].([]interface{})
	stealData, ok := formData["stealData"].([]interface{})

	// var idx = 1
	// idx1 = 0

	if ok {
		idx1 = 1
		for i, v := range stealData {
			vv, ok := v.(map[string]interface{})
			if ok {
				rich, ok := vv["rich"].(float64)
				if ok {
					if rich == 1 {
						log.Printf("【摇一摇】rich is :%v", i)
						idx1 = i
						break
					}
				}
			}
		}
	} else {
		idx1 = -1
	}
	result = false

	for _, v := range stealResult {
		vv := v.(map[string]interface{})
		if vv["gold"].(float64) == addGold {
			if vv["isRich"].(float64) == 1 {
				result = true
				return
			} else {
				result = false
			}
		}
	}

	return
}

// 【摇一摇】--分享获取乐园骰子
func getShareDrawDice(serverURL, zoneToken string, itemId int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getShareDrawDice&token=%v&share=1&itemId=%v&now=%v", serverURL, zoneToken, itemId, now)
	formData := httpGetReturnJson(url)
	log.Println("当前骰子数量:", formData["richmanDice"])
}

//shareType -> 43=铲子助力 42=海浪助力
func beachHelp(serverURL, zoneToken, uid string, shareType int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=beachHelp&token=%v&fuid=%v&shareType=%v&now=%v", serverURL, zoneToken, uid, shareType, now)
	httpGetReturnJson(URL)
}

// type=0 领取海浪 type=1 领取铲子 index=0,1,2
func getHelpItem(serverURL, zoneToken string, _type, index int) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getHelpItem&token=%v&type=1&index=%v&now=%v", serverURL, zoneToken, index, now)
	if _type == 0 {
		URL = fmt.Sprintf("%v/game?cmd=getHelpItem&token=%v&type=0&now=%v", serverURL, zoneToken, now)
	}

	httpGetReturnJson(URL)
}

func useShovel(serverURL, zoneToken, toid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=useShovel&token=%v&now=%v", serverURL, zoneToken, now)
	if toid != "" {
		URL = fmt.Sprintf("%v/game?cmd=useShovel&token=%v&targetUid=%v&now=%v", serverURL, zoneToken, toid, now)
	}
	httpGetReturnJson(URL)
}

// index->0...4 direction=0..1
func getBeachLineRewards(serverURL, zoneToken string, index, direction int) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getBeachLineRewards&token=%v&index=%v&direction=%v&now=%v", serverURL, zoneToken, index, direction, now)
	httpGetReturnJson(URL)
}

// 刷新海滩
func refreshBeach(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=refreshBeach&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(URL)
}

// 获取超值返利
func getElevenEnergyPrize(serverURL, zoneToken string, id int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getElevenEnergyPrize&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(URL)
}

// 获取日常任务奖励
func getDayTaskAward(serverURL, zoneToken, key string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getDayTaskAward&token=%v&type=%v&now=%v", serverURL, zoneToken, key, now)
	httpGetReturnJson(URL)

}

// 查看日常任务列表
func getDayTasksInfo(serverURL, zoneToken string) (keys []string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getDayTasksInfo&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(URL)

	tasks, ok := formData["tasks"].(map[string]interface{})

	if ok {
		for k := range tasks {
			keys = append(keys, k)
		}
	}

	return
}

func exchangeRiceCake(serverURL, zoneToken string, id int64) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=exchangeRiceCake&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	formData := httpGetReturnJson(URL)
	if _, ok := formData["error"]; ok {
		return false
	}
	return true
}

// 兑换金条
// func exchangeGoldChunk(serverURL, zoneToken string) {
// 	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
// 	URL := fmt.Sprintf("%v/game?cmd=exchangeGoldChunk&token=%v&now=%v", serverURL, zoneToken, now)
// 	httpGetReturnJson(URL)
// }

// func getFamilyRobTaskPrize(serverURL, zoneToken string) {

// }

func unlockWorker(serverURL, zoneToken, id string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=unlockWorker&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(URL)
}

// 公会赛季积分宝箱领取
// func getBoatRaceScorePrize(serverURL, zoenToken string, id int64) {
// 	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
// 	URL := fmt.Sprintf("%v/game?cmd=getBoatRaceScorePrize&token=%v&id=%v&now=%v", serverURL, zoenToken, id, now)
// 	httpGetReturnJson(URL)
// }

// 首次开启汤圆
func enterSteamBox(serverURL, zoneToken, fuid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=enterSteamBox&token=%v&fuid=%v&now=%v", serverURL, zoneToken, fuid, now)
	httpGetReturnJson(URL)
}

func openSteamBox(serverURL, zoneToken, fuid string) (startTime float64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=openSteamBox&token=%v&fuid=%v&now=%v", serverURL, zoneToken, fuid, now)
	formData := httpGetReturnJson(URL)

	steamBox, ok := formData["steamBox"].(map[string]interface{})
	if ok {
		startTime = steamBox["startTime"].(float64)
		return
	}
	return
}

func collectMineGold(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=collectMineGold&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(URL)
}

func getSteamBoxHelpList(serverURL, zoneToken string) (uids []float64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getSteamBoxHelpList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(URL)
	helpList, ok := formData["helpList"].([]interface{})
	if ok {
		for _, v := range helpList {
			vv, ok := v.(map[string]interface{})
			if ok {
				uid, ok := vv["uid"].(float64)
				if ok {
					uids = append(uids, uid)
				}
			}
		}
	}
	return
}

func addFirewood(serverURL, zoneToken, fuid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=addFirewood&token=%v&fuid=%v&now=%v", serverURL, zoneToken, fuid, now)
	httpGetReturnJson(URL)
}

func familyChat(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=familyChat&token=%v&content=1&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(URL)
}

// []string{"9","15","3"}
func getFamilyDayTaskPrize(serverURL, zoneToken, id string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=getFamilyDayTaskPrize&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(URL)
}

// // 统一下单接口
// func gameAPI(method, url, params string) (formData map[string]interface{}) {
// 	url = url + "/game?" + params
// 	if method == "GET" {
// 		httpGetReturnJson(url)
// 		formData = httpGetReturnJson(url)
// 		return
// 	}
// 	return
// }

func httpGetReturnJson(url string) (formData map[string]interface{}) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("httpGet err is %v, url is %v", err, url)
		return
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
	response, err := client.Do(request)

	if err != nil {
		log.Printf("httpGet err is %v, url is %v", err, url)
		return
	}
	defer response.Body.Close()
	formData = make(map[string]interface{})
	json.NewDecoder(response.Body).Decode(&formData)
	return
}

// 生成区间[-m, n]的安全随机数
func RangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}

	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))

		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

func sendMsg(msg string) {
	url := "https://rocket.chat.rosettawe.com/api/v1/login"

	client := &http.Client{}

	sendMsg, err := json.Marshal(map[string]interface{}{"user": "whn", "password": "Aa112211"})

	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(sendMsg))
	if err != nil {
		log.Printf("httpGet err is %v, url is %v", err, url)
		return
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36")
	request.Header.Add("Content-type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		log.Printf("httpGet err is %v, url is %v", err, url)
		return
	}
	defer response.Body.Close()
	formData := make(map[string]interface{})
	json.NewDecoder(response.Body).Decode(&formData)

	data, ok := formData["data"].(map[string]interface{})

	if !ok {
		return
	}

	authToken := data["authToken"].(string)

	me, ok := data["me"].(map[string]interface{})

	if !ok {
		return
	}

	id := me["_id"].(string)

	url = "https://rocket.chat.rosettawe.com/api/v1/chat.sendMessage"

	var message = map[string]interface{}{"message": map[string]interface{}{"rid": "48AM8JoiSdRYgCB9W", "msg": msg}}

	sendMsg, err = json.Marshal(message)

	if err != nil {
		return
	}

	client = &http.Client{}

	request, err = http.NewRequest("POST", url, bytes.NewBuffer(sendMsg))
	if err != nil {
		return
	}

	request.Header.Add("X-Auth-Token", authToken)
	request.Header.Add("X-User-Id", id)
	request.Header.Add("Content-type", "application/json")
	client.Do(request)

}

// runner

func RunnerSteamBox() (err error) {
	if runnerStatus("steamBoxStatus") == "0" {
		return
	}
	openSteamBoxGo()

	SQL := "select id, name, token from tokens where id <> 302691822 and id <> 309392050"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var uid, name, token string
		rows.Scan(&uid, &name, &token)
		serverURL, zoneToken := getSeverURLAndZoneToken(token)
		addFirewood(serverURL, zoneToken, "302691822")
		addFirewood(serverURL, zoneToken, "309392050")
	}

	return
}

func RunnerDraw() (err error) {
	// RunnerSteamBox()
	drawStatus := runnerStatus("drawStatus")
	if drawStatus == "0" {
		return
	}
	hour := time.Now().Hour()

	if hour == 4 {
		SQL := "select id, token, name from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
		var cowBoyUid, cowBoyToken, cowBoyName string
		Pool.QueryRow(SQL).Scan(&cowBoyUid, &cowBoyToken, &cowBoyName)
		serverURL, zoneToken := getSeverURLAndZoneToken(cowBoyToken)

		go func() {
			log.Printf("---------------------------[%v]开始转盘---------------------------", cowBoyName)
			followCompanion_1(serverURL, zoneToken, 2)
			for i := 0; i < 30; i++ {
				draw(cowBoyUid, cowBoyName, serverURL, zoneToken, 1)
				time.Sleep(time.Millisecond * 2500)
			}
			followCompanion_1(serverURL, zoneToken, 3)
			log.Printf("---------------------------[%v]结束转盘---------------------------", cowBoyName)

		}()

	}

	if drawStatus != "1" {
		SQL := "select id, token, name from tokens where id = ?"
		var cowBoyUid, cowBoyToken, cowBoyName string
		Pool.QueryRow(SQL, drawStatus).Scan(&cowBoyUid, &cowBoyToken, &cowBoyName)
		serverURL, zoneToken := getSeverURLAndZoneToken(cowBoyToken)

		go func() {
			log.Printf("---------------------------[%v]开始转盘---------------------------", cowBoyName)
			followCompanion_1(serverURL, zoneToken, 2)
			for i := 0; i < 30; i++ {
				draw(cowBoyUid, cowBoyName, serverURL, zoneToken, 1)
				time.Sleep(time.Millisecond * 2500)
			}
			followCompanion_1(serverURL, zoneToken, 3)

			taskIDs := getDayTasksInfo(serverURL, zoneToken)
			log.Printf("[%v]领取日常任务奖励:%v", cowBoyName, taskIDs)
			for _, taskID := range taskIDs {
				// time.Sleep(time.Millisecond * 100)
				getDayTaskAward(serverURL, zoneToken, taskID)
			}

			log.Printf("[%v]领取超值返利", cowBoyName)
			getElevenEnergyPrize(serverURL, zoneToken, 1)
			getElevenEnergyPrize(serverURL, zoneToken, 2)
			getElevenEnergyPrize(serverURL, zoneToken, 3)
			getElevenEnergyPrize(serverURL, zoneToken, 4)
			log.Printf("[%v]collectMineGold", cowBoyName)
			collectMineGold(serverURL, zoneToken)
			log.Printf("---------------------------[%v]结束转盘---------------------------", cowBoyName)

		}()
		return

	}

	if hour == 2 || hour == 5 || hour == 8 || hour == 11 || hour == 14 || hour == 17 || hour == 20 || hour == 23 {
		log.Println("start draw")
		SQL := "select id, name, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'drawIds'))"

		rows, err := Pool.Query(SQL)

		if err != nil {
			return err
		}
		defer rows.Close()

		var users []User

		for rows.Next() {
			var user User
			err = rows.Scan(&user.Uid, &user.Name, &user.Token)
			if err != nil {
				break
			}
			user.ServerURL = getServerURL()

			user.ZoneToken, user.FamilyDayTask = getEnterInfo(user.Uid, user.Name, user.ServerURL, user.Token, "familyDayTask")
			// user.ZoneToken = getZoneToken(user.ServerURL, user.Token)

			users = append(users, user)

		}
		for _, u := range users {

			goName := u.Name
			goUid := u.Uid
			goServerURL := u.ServerURL
			goZoneToken := u.ZoneToken
			goFamilyDayTask := u.FamilyDayTask

			if goUid == "302691822" {
				log.Printf("---------------------------[%v]开始转盘---------------------------", goName)
				followCompanion_1(goServerURL, goZoneToken, 2)
				amount := draw(goUid, goName, goServerURL, goZoneToken, 1)
				time.Sleep(time.Millisecond * 2100)
				for i := 0; i <= int(amount); i++ {
					draw(goUid, goName, goServerURL, goZoneToken, 1)
					time.Sleep(time.Millisecond * 2100)
				}
				if goUid == "302691822" || goUid == "309392050" {
					followCompanion_1(goServerURL, goZoneToken, 3)
				} else {
					followCompanion_1(goServerURL, goZoneToken, 1)
				}
				log.Printf("---------------------------[%v]结束转盘---------------------------", goName)

				taskIDs := getDayTasksInfo(goServerURL, goZoneToken)
				log.Printf("[%v]领取日常任务奖励:%v", goName, taskIDs)
				for _, taskID := range taskIDs {
					// time.Sleep(time.Millisecond * 100)
					getDayTaskAward(goServerURL, goZoneToken, taskID)
				}

				log.Printf("[%v]领取超值返利", goName)
				getElevenEnergyPrize(goServerURL, goZoneToken, 1)
				getElevenEnergyPrize(goServerURL, goZoneToken, 2)
				getElevenEnergyPrize(goServerURL, goZoneToken, 3)
				getElevenEnergyPrize(goServerURL, goZoneToken, 4)

				log.Printf("[%v]collectMineGold", goName)
				collectMineGold(goServerURL, goZoneToken)

				if goFamilyDayTask == nil {
					return err
				}

				for k := range goFamilyDayTask.(map[string]interface{}) {
					getFamilyDayTaskPrize(goServerURL, goZoneToken, k)
					log.Printf("[%v]领取公会任务奖励[%v]", goName, k)
				}
			} else {
				go func() {
					log.Printf("---------------------------[%v]开始转盘---------------------------", goName)
					followCompanion_1(goServerURL, goZoneToken, 2)
					amount := draw(goUid, goName, goServerURL, goZoneToken, 1)
					time.Sleep(time.Millisecond * 2100)
					for i := 0; i <= int(amount); i++ {
						draw(goUid, goName, goServerURL, goZoneToken, 1)
						time.Sleep(time.Millisecond * 2100)
					}
					if goUid == "302691822" || goUid == "309392050" {
						followCompanion_1(goServerURL, goZoneToken, 3)
					} else {
						followCompanion_1(goServerURL, goZoneToken, 1)
					}
					log.Printf("---------------------------[%v]结束转盘---------------------------", goName)

					taskIDs := getDayTasksInfo(goServerURL, goZoneToken)
					log.Printf("[%v]领取日常任务奖励:%v", goName, taskIDs)
					for _, taskID := range taskIDs {
						// time.Sleep(time.Millisecond * 100)
						getDayTaskAward(goServerURL, goZoneToken, taskID)
					}

					log.Printf("[%v]领取超值返利", goName)
					getElevenEnergyPrize(goServerURL, goZoneToken, 1)
					getElevenEnergyPrize(goServerURL, goZoneToken, 2)
					getElevenEnergyPrize(goServerURL, goZoneToken, 3)
					getElevenEnergyPrize(goServerURL, goZoneToken, 4)

					log.Printf("[%v]collectMineGold", goName)
					collectMineGold(goServerURL, goZoneToken)

					if goFamilyDayTask == nil {
						return
					}

					for k := range goFamilyDayTask.(map[string]interface{}) {
						getFamilyDayTaskPrize(goServerURL, goZoneToken, k)
						log.Printf("[%v]领取公会任务奖励[%v]", goName, k)
					}

				}()
			}

		}
		return err

	}

	log.Println("no draw")
	return
}

func RunnerPullAnimal() (err error) {

	// 7 10 13 16 19 22 点执行

	if runnerStatus("pullAnimalGoStatus") == "0" {
		return
	}

	hour := time.Now().Hour()

	if hour >= 7 && hour <= 22 {
		log.Println("5秒后开始拉动物...")

		for i := 1; i <= 5; i++ {
			time.Sleep(time.Second * 1)
			log.Printf("%v秒后开始拉动物...", 5-i)
		}

		log.Println("现在开始拉动物")

		SQL := "select id, token, name, pull_rows from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
		var uid, token, name, pullRows string

		err = Pool.QueryRow(SQL).Scan(&uid, &token, &name, &pullRows)

		if err != nil {
			return
		}

		serverURL, zoneToken := getSeverURLAndZoneToken(token)

		if zoneToken == "" {
			sendMsg(uid + ":" + name)
			log.Printf("[name: %v] token is invalid\n", name)
		}

		foods := enterFamilyRob(serverURL, zoneToken)

		var nextFlag bool
		for _, v := range foods {
			// myTeam := v["myTeam"].(int)
			// if myTeam != 4 {
			if strings.Contains(pullRows, fmt.Sprintf("%v", v["row"])) {
				robFamilyFood(serverURL, zoneToken, v["id"].(string))
				nextFlag = true
				break
			}
			// }

		}
		log.Printf("[%v]拉动物完成", name)
		// time.Sleep(time.Second * 1)
		// log.Printf("cowboy serverURL:%v, zoneToken:%v\n", serverURL, zoneToken)

		pullAnimalGo()

		time.Sleep(900 * time.Second)

		if nextFlag {
			serverURL, zoneToken = getSeverURLAndZoneToken(token)

			if zoneToken == "" {
				sendMsg(uid + ":" + name)
				log.Printf("[name: %v] token is invalid\n", name)
			}

			foods = enterFamilyRob(serverURL, zoneToken)

			for _, v := range foods {
				// myTeam := v["myTeam"].(int)
				// if myTeam != 4 {
				if strings.Contains(pullRows, fmt.Sprintf("%v", v["row"])) {
					robFamilyFood(serverURL, zoneToken, v["id"].(string))
					break
				}
				// }

			}
			log.Printf("[%v]拉动物完成", name)
		}

		return
	}

	log.Println("RunnerPullAnimal do not run")

	return
}

func runnerStatus(confKey string) (confValue string) {
	Pool.QueryRow("select conf_value from config where conf_key = ?", confKey).Scan(&confValue)
	return
}
