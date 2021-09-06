package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wuhuinan47/cat/runner"
	wechatapi "github.com/wuhuinan47/cat/wechatAPI"
)

var Pool *sql.DB

// type CatData struct {
// 	GoldMineHelpList []GoldMineHelpList `json:"helpList"`
// }

// type GoldMineHelpList struct {
// 	Uid     int64 `json:"uid"`
// 	Quality int64 `json:"quality"`
// }

type User struct {
	Uid       string
	Token     string
	Name      string
	ServerURL string
	ZoneToken string
	Public    map[string]interface{}
}

func main() {

	// db, err := sql.Open("mysql", "cat:cyydmkj123@tcp(localhost:3306)/data_cat?charset=utf8&parseTime=true&collation=utf8mb4_unicode_ci")
	db, err := sql.Open("mysql", "cat:cyydmkj123@tcp(localhost:3306)/data_cat")
	if err != nil {
		panic(err)
	}
	db.SetConnMaxLifetime(100 * time.Second) //最大连接周期，超过时间的连接就close
	db.SetMaxOpenConns(100)                  //设置最大连接数
	db.SetMaxIdleConns(16)                   //设置闲置连接数

	Pool = db

	http.HandleFunc("/autoDraw", autoDrawH)
	http.HandleFunc("/login", LoginH)
	http.HandleFunc("/update", UpdateH)
	http.HandleFunc("/test", TestH)
	http.HandleFunc("/pullAnimal", PullAnimalH)
	http.HandleFunc("/diamond", DiamondH)
	http.HandleFunc("/familySign", familySignH)
	http.HandleFunc("/sendQrcode", SendQrcodeH)
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
	http.HandleFunc("/testDraw", TestDrawH)

	//

	//

	//

	//OneSonAttackBossH

	//

	http.HandleFunc("/getServerURL", GetServerURLH)
	http.HandleFunc("/getZoneToken", GetZoneTokenH)

	//
	http.HandleFunc("/checkToken", CheckTokenH)

	//

	// wechatAPI
	http.HandleFunc("/wechatAPI/sendMsg", wechatapi.SendMsgH)

	http.HandleFunc("/ctrl.html", IndexH)
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
	}

	defer func() {
		sendMsg("cat is stop")
		log.Println("All done, good bye")
	}()

	go http.ListenAndServe(":33333", nil)
	log.Println("listen port 33333 sucesss")

	// if err = http.ListenAndServe(":33333", nil); err != nil {
	// 	log.Fatal(err)
	// }

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c

}

func IndexH(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("ctrl.html")
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
	cmd := "rm -rf /mnt/app/cat/screenlog.0"
	c := exec.Command("bash", "-c", cmd)
	_, err := c.CombinedOutput()
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "SUCCESS")
	return
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

	// io.WriteString(w, string(bytes))
	// w.Header().Set("content-type", "application/json;charset=utf-8")
	w.Write(bytes)
}

//

func SendQrcodeH(w http.ResponseWriter, req *http.Request) {
	qrcode := req.URL.Query().Get("qrcode")
	Pool.Exec("update config set conf_value = ? where conf_key = 'wechatLoginQrcode'", qrcode)

	sendMsg(qrcode)
	log.Println("qrcode:", qrcode)
	io.WriteString(w, qrcode)
	return
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
	return
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
	}

	rows, err := Pool.Query(SQL)

	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var uid, token string
		rows.Scan(&uid, &token)

		serverURL := getServerURL()

		zoneToken := getZoneToken(serverURL, token)

		bossList := getBossHelpList(serverURL, zoneToken)

		for _, v := range bossList {
			leftHp := v["leftHp"].(float64)
			if leftHp >= 1500 {
				attackBoss(serverURL, zoneToken, v["id"].(string))
			}
		}

		// go func() {

		// }()

	}

	io.WriteString(w, "SUCCESS")
	return
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

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	bossList := getBossHelpList(serverURL, zoneToken)
	log.Printf("---------------------------[%v]开始打龙---------------------------", name)

	for _, v := range bossList {
		leftHp := v["leftHp"].(float64)
		log.Println("leftHp:", leftHp)
		if leftHp <= 600 {
			attackBossAPI(serverURL, zoneToken, v["id"].(string), 3, 0, 1, 200, 200)
		} else if leftHp <= 1000 {
			attackBossAPI(serverURL, zoneToken, v["id"].(string), 3, 0, 1, 195, 200)
			attackBossAPI(serverURL, zoneToken, v["id"].(string), 1, 1, 1, 395, 400)
		}

		// if leftHp < 1000 && leftHp >= 500 {
		// 	attackBossByAdmin(serverURL, zoneToken, v["id"].(string))
		// }
	}

	log.Printf("---------------------------[%v]结束打龙---------------------------", name)

	io.WriteString(w, "SUCCESS")
	return

}

//
func SixEnergyH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	for i := 0; i < 6; i++ {
		getSixEnergy(serverURL, zoneToken)
		time.Sleep(time.Millisecond * 600)
	}
	io.WriteString(w, "SUCCESS")
	return
}
func SetPieceH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	ids := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, v := range ids {
		setPiece(serverURL, zoneToken, v)
		time.Sleep(time.Second * 1)
	}
	io.WriteString(w, "SUCCESS")
	return
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
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)
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
	return

}

//
func GetMailPrizeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	mailids := getMailList(serverURL, zoneToken)

	log.Printf("id:%v 开始领取邮件奖励", id)
	for _, mailid := range mailids {
		getMailAttachments(serverURL, zoneToken, mailid)
		time.Sleep(time.Second * 1)
	}

	log.Printf("id:%v 领取邮件完毕", id)

	go func() {
		log.Printf("id:%v 开始删除邮件", id)
		for _, mailid := range mailids {
			deleteMail(serverURL, zoneToken, mailid)
			time.Sleep(time.Second * 1)
		}
		log.Printf("id:%v 删除邮件完毕", id)
	}()

	io.WriteString(w, "SUCCESS")
	return

}

func ThrowDiceH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	amount := req.URL.Query().Get("amount")
	intAmount, _ := strconv.Atoi(amount)
	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, id).Scan(&token)

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	go func() {
		log.Println("start throwDice")

		for i := 1; i <= intAmount; i++ {
			throwDice(serverURL, zoneToken)

			if i == intAmount {
				break
			}
			time.Sleep(time.Second * 3)
		}
		log.Println("end throwDice")

	}()

	io.WriteString(w, "SUCCESS")
	return
}

func DrawH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	drawMulti := req.URL.Query().Get("drawMulti")
	amount := req.URL.Query().Get("amount")
	intAmount, _ := strconv.Atoi(amount)

	SQL := "select id, name, token from tokens where id = ?"

	var uid, name, token string
	Pool.QueryRow(SQL, id).Scan(&uid, &name, &token)

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	go func() {
		log.Printf("---------------------------[%v]开始转盘---------------------------", name)
		followCompanion(serverURL, zoneToken, 2)
		for i := 1; i <= intAmount; i++ {
			draw(uid, name, serverURL, zoneToken, drawMulti)
			log.Println("剩余转盘次数:", intAmount-i)
			time.Sleep(time.Millisecond * 2100)
		}
		if id == "302691822" {
			followCompanion(serverURL, zoneToken, 3)
		} else {
			followCompanion(serverURL, zoneToken, 1)
		}
		log.Printf("---------------------------[%v]结束转盘---------------------------", name)
		return
	}()

	io.WriteString(w, "SUCCESS")

	return

}

func GiftPieceH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	fromUid := req.URL.Query().Get("fromUid")
	toUid := req.URL.Query().Get("toUid")
	amount := req.URL.Query().Get("amount")

	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, fromUid).Scan(&token)

	serverURL := getServerURL()

	zoneToken := getZoneToken(serverURL, token)

	intAmount, _ := strconv.Atoi(amount)

	log.Println("amount:", amount)
	log.Println("intAmount:", intAmount)
	for i := 1; i <= intAmount; i++ {
		giftPiece(serverURL, zoneToken, id, toUid)
		time.Sleep(time.Second * 1)
	}

	io.WriteString(w, "SUCCESS")
	return
}

func TestDrawH(w http.ResponseWriter, req *http.Request) {
	log.Println("start draw")
	SQL := "select id, name, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'drawIds'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		io.WriteString(w, "FAIL")

		return
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
		user.ZoneToken = getZoneToken(user.ServerURL, user.Token)

		users = append(users, user)

	}
	for _, u := range users {

		goName := u.Name
		goUid := u.Uid
		goServerURL := u.ServerURL
		goZoneToken := u.ZoneToken
		go func() {
			log.Printf("---------------------------[%v]开始转盘---------------------------", goName)
			followCompanion(goServerURL, goZoneToken, 2)

			energy := draw(goUid, goName, goServerURL, goZoneToken, 1)

			time.Sleep(time.Millisecond * 2100)

			for i := 0; i <= int(energy); i++ {
				draw(goUid, goName, goServerURL, goZoneToken, 1)
				time.Sleep(time.Millisecond * 2100)
			}
			if goUid == "302691822" {
				followCompanion(goServerURL, goZoneToken, 3)
			} else {
				followCompanion(goServerURL, goZoneToken, 1)
			}
			log.Printf("---------------------------[%v]结束转盘---------------------------", goName)

		}()
	}
	io.WriteString(w, "SUCCESS")
	return
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
			serverURL := getServerURL()
			zoneToken := getZoneToken(serverURL, token)

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

	return
}

func AttackMyBossH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	mode := req.URL.Query().Get("mode")

	SQL := "select name, token from tokens where id = ?"

	var name, token string

	err := Pool.QueryRow(SQL, id).Scan(&name, &token)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	serverURL := getServerURL()
	zoneToken := getZoneToken(serverURL, token)
	bossID := summonBoss(serverURL, zoneToken)
	inviteBoss(serverURL, zoneToken, bossID)
	time.Sleep(time.Second * 1)
	shareAPI(serverURL, zoneToken)
	getFreeBossCannon(serverURL, zoneToken)
	log.Printf("---------------------------[%v]开始打龙---------------------------", name)
	attackMyBoss(serverURL, zoneToken, bossID, mode)
	log.Printf("---------------------------[%v]结束打龙---------------------------", name)

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
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)
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

	for _, v2 := range mmBossList {
		for _, v := range mmList1 {
			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag == true {
				go attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}
		}
		time.Sleep(time.Second * 18)
	}

	for _, v2 := range mmBossList {
		for _, v := range mmList2 {
			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag == true {
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
	return

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

		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)
		bossID := summonBoss(serverURL, zoneToken)
		inviteBoss(serverURL, zoneToken, bossID)
		mmList = append(mmList, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
		mmBossList = append(mmBossList, bossID)
	}
	for _, v := range mmList {
		for _, v2 := range mmBossList {
			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag == true {
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
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)
		bossID := summonBoss(serverURL, zoneToken)
		inviteBoss(serverURL, zoneToken, bossID)
		nnList = append(nnList, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
		nnBossList = append(nnBossList, bossID)
	}
	for _, v := range nnList {
		for _, v2 := range nnBossList {

			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag == true {
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
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)
		bossID := summonBoss(serverURL, zoneToken)
		inviteBoss(serverURL, zoneToken, bossID)
		boss3List = append(boss3List, map[string]interface{}{"uid": uid, "serverURL": serverURL, "zoneToken": zoneToken})
		boss3BossList = append(boss3BossList, bossID)
	}
	for _, v := range boss3List {
		for _, v2 := range boss3BossList {

			flag := enterBoss(v["serverURL"].(string), v["zoneToken"].(string), v2, v["uid"].(float64))
			if flag == true {
				attackBoss(v["serverURL"].(string), v["zoneToken"].(string), v2)
			}

		}
	}

	io.WriteString(w, "SUCCESS")
	return
}

//

func LoginByQrcodeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	str := "/mnt/app/cat/linuxWechat.py"

	SQL := "select conf_value from config where conf_key = 'wechatLoginQrcode'"
	if id == "2" {
		str = "/mnt/app/cat/linuxQQ.py"
		SQL = "select conf_value from config where conf_key = 'qqLoginQrcode'"
	}

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

	if id == "3" {
		http.Redirect(w, req, qrcode, http.StatusTemporaryRedirect)
	}
	io.WriteString(w, qrcode)

	return

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
	return
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
	go func() {
		SQL := "select id, token, name from tokens"

		rows, err := Pool.Query(SQL)

		if err != nil {
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id, token, name string
			rows.Scan(&id, &token, &name)

			serverURL := getServerURL()
			zoneToken := getZoneToken(serverURL, token)

			if zoneToken == "" {
				Pool.Exec("update tokens set token = '' where id = ?", id)
				sendMsg(id + ":" + name)
			}
			time.Sleep(time.Second * 1)

		}
		return
	}()
	io.WriteString(w, "SUCCESS")
	return

}

func GetFreeBossCannonH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	var token string
	Pool.QueryRow("select token from tokens where id = ?", id).Scan(&token)
	serverURL := getServerURL()
	zoneToken := getZoneToken(serverURL, token)
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

	serverURL := getServerURL()
	zoneToken := getZoneToken(serverURL, token)

	uids := getFriendsCandyTreeInfo(serverURL, zoneToken, v2)
	time.Sleep(time.Second * 1)
	log.Println("getFriendsCandyTreeInfo uids:", uids)

	var targetAmount float64

	for _, v := range uids {
		log.Println("getCandyTreeInfo uid:", v)

		posList := getCandyTreeInfo(serverURL, zoneToken, v)
		log.Println("posList:", posList)
		time.Sleep(time.Second * 1)
		log.Println("start uid:", v)

		for _, v4 := range posList {

			if v3 == targetAmount {
				return
			}

			time.Sleep(time.Second * 3)

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

	getBoxPrizeGo(id, v2, v3)

	io.WriteString(w, "SUCCESS")
}

func familySignH(w http.ResponseWriter, req *http.Request) {
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

	SQL := "replace into tokens (id, name, token) values (?, ?, ?)"
	_, err := Pool.Exec(SQL, id, name, token)

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

func autoDrawH(w http.ResponseWriter, req *http.Request) {
	url := "https://s147.11h5.com:3147/123_206_192_93/3148//game?cmd=draw&token=ilddcNcOr4NRTwKqvAL6v76tCgPz-_nJVai&drawMulti=1&now=1628759013206"

	resp, _ := http.Get(url)

	data := make(map[string]interface{})

	json.NewDecoder(resp.Body).Decode(&data)

	lastEnergyTime, _ := data["lastEnergyTime"].(float64)

	log.Println("lastEneygyTime:", lastEnergyTime)
	energy, _ := data["energy"].(float64)

	log.Println("energy:", energy)

	go func() {

		var lastTime float64
		for i := 0; i < int(energy); i++ {
			log.Println("now lastTime:", lastTime)

			url := fmt.Sprintf("https://s147.11h5.com:3147/123_206_192_93/3148//game?cmd=draw&token=ilddcNcOr4NRTwKqvAL6v76tCgPz-_nJVai&drawMulti=1&now=%v", lastTime)

			resp, _ := http.Get(url)

			data := make(map[string]interface{})

			json.NewDecoder(resp.Body).Decode(&data)

			lastEnergyTime, _ := data["lastEnergyTime"].(float64)
			lastTime = lastEnergyTime

			energy, ok := data["energy"].(float64)
			log.Println("now energy:", energy)
			if ok && energy <= 0 {
				break
			}

			time.Sleep(time.Second * 3)
		}

	}()

	io.WriteString(w, "success")

}

// goroutine

func attackBossGo() {
	SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'cowBoss1'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()

}

// 一键拉动物
func pullAnimalGo() {
	log.Println("start pullAnimalGo")
	SQL := "select id, token, name from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"
	rows, err := Pool.Query(SQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var uid, token, name string
		rows.Scan(&uid, &token, &name)
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid + ":" + name)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}

		log.Printf("[%v]开始拉动物", name)
		foods := enterFamilyRob(serverURL, zoneToken)

		for _, v := range foods {
			myTeam := v["myTeam"].(int)
			if myTeam != 4 {
				robFamilyFood(serverURL, zoneToken, v["id"].(string))
				break
			}

		}
		log.Printf("[%v]拉动物完成", name)
		// time.Sleep(time.Second * 1)
		// log.Printf("serverURL:%v, zoneToken:%v\n", serverURL, zoneToken)
	}
	return
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
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

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
			time.Sleep(time.Second * 2)
			j = 1
		}

		log.Printf("[%v] start familySign", user["name"])
		familySign(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end familySign", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start getSignPrize", user["name"])
		getSignPrize(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getSignPrize", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start getFreeDailyGiftBox", user["name"])
		getFreeDailyGiftBox(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getFreeDailyGiftBox", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start playLuckyWheel", user["name"])
		shareAPI(user["serverURL"], user["zoneToken"])
		playLuckyWheel(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end playLuckyWheel", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start getFreeClamp", user["name"])
		shareAPI(user["serverURL"], user["zoneToken"])
		getFreeClamp(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getFreeClamp", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start getInviteSnow", user["name"])
		shareAPI(user["serverURL"], user["zoneToken"])
		getInviteSnow(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end getInviteSnow", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start autoFriendEnergy", user["name"])
		autoFriendEnergy(user["serverURL"], user["zoneToken"])
		log.Printf("[%v] end autoFriendEnergy", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		log.Printf("[%v] start getSixEnergy", user["name"])
		for i := 0; i < 6; i++ {
			time.Sleep(time.Second * 1)
			getSixEnergy(user["serverURL"], user["zoneToken"])
		}
		log.Printf("[%v] end getSixEnergy", user["name"])
		j++
	}
	time.Sleep(time.Second * 1)

	for _, user := range tokenList {
		if j >= 6 {
			time.Sleep(time.Second * 2)
			j = 1
		}
		gameList := []string{"535", "525", "157", "452", "411"}

		for _, v := range gameList {
			getAward(user["token"], v)
		}
		j++
	}
	time.Sleep(time.Second * 1)

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
	time.Sleep(time.Second * 1)

	// 为大佬海浪、铲子助力

	var times int64 = 1
	for _, user := range tokenList {
		if times > 6 {
			break
		}
		if times > 3 {
			log.Printf("[%v] 为蜜蜜海浪助力", user["name"])
			beachHelp(user["serverURL"], user["zoneToken"], "309392050", 42)
		} else {
			log.Printf("[%v] 为牛海浪助力", user["name"])
			beachHelp(user["serverURL"], user["zoneToken"], "302691822", 42)
		}
		times++
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
						time.Sleep(time.Second * 2)
					}
				}
			}

		}
	}

	time.Sleep(time.Second * 1)

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

	getAwardForCowBoy()

	// for rows.Next() {
	// 	var uid, token, name string
	// 	rows.Scan(&uid, &token, &name)
	// 	serverURL := getServerURL()
	// 	zoneToken := getZoneToken(serverURL, token)

	// 	if zoneToken == "" {
	// 		sendMsg(uid + ":" + name)
	// 		log.Printf("[uid: %v] token is invalid\n", uid)
	// 	}

	// 	log.Printf("[%v] start familySign", name)
	// 	familySign(serverURL, zoneToken)
	// 	log.Printf("[%v] end familySign", name)
	// 	time.Sleep(time.Second * 1)

	// 	log.Printf("[%v] start getSignPrize", name)
	// 	getSignPrize(serverURL, zoneToken)
	// 	log.Printf("[%v] end getSignPrize", name)
	// 	time.Sleep(time.Second * 1)

	// 	log.Printf("[%v] start getFreeDailyGiftBox", name)
	// 	getFreeDailyGiftBox(serverURL, zoneToken)
	// 	log.Printf("[%v] end getFreeDailyGiftBox", name)
	// 	time.Sleep(time.Second * 1)

	// 	log.Printf("[%v] start playLuckyWheel", name)
	// 	shareAPI(serverURL, zoneToken)
	// 	playLuckyWheel(serverURL, zoneToken)
	// 	log.Printf("[%v] end playLuckyWheel", name)
	// 	time.Sleep(time.Second * 1)

	// 	log.Printf("[%v] start getFreeClamp", name)
	// 	shareAPI(serverURL, zoneToken)
	// 	getFreeClamp(serverURL, zoneToken)
	// 	log.Printf("[%v] end getFreeClamp", name)
	// 	time.Sleep(time.Second * 1)

	// 	log.Printf("[%v] start getInviteSnow", name)
	// 	shareAPI(serverURL, zoneToken)
	// 	getInviteSnow(serverURL, zoneToken)
	// 	log.Printf("[%v] end getInviteSnow", name)

	// 	time.Sleep(time.Second * 1)
	// 	log.Printf("[%v] start autoFriendEnergy", name)
	// 	autoFriendEnergy(serverURL, zoneToken)
	// 	log.Printf("[%v] end autoFriendEnergy", name)

	// 	log.Printf("[%v] start getSixEnergy", name)
	// 	for i := 0; i < 6; i++ {
	// 		time.Sleep(time.Second * 1)
	// 		getSixEnergy(serverURL, zoneToken)
	// 	}
	// 	log.Printf("[%v] end getSixEnergy", name)

	// 	gameList := []string{"535", "525", "157", "452", "411"}

	// 	time.Sleep(time.Second * 1)
	// 	for _, v := range gameList {
	// 		getAward(token, v)
	// 	}
	// }

	// time.Sleep(time.Second * 1)
	// log.Println("start getFamilySignPrizeGo..")
	// getFamilySignPrizeGo()
	// log.Println("end getFamilySignPrizeGo..")

	// getAwardForCowBoy()

	log.Println("familySignGo finish")

	othersSign()
	return
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
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

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
						time.Sleep(time.Second * 2)
					}
				}
			}

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

// 一键领取公会签到奖励
func getFamilySignPrizeGo() {
	SQL := "select id, token, name from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var uid, token, name string
		rows.Scan(&uid, &token, &name)
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid + ":" + name)
			log.Printf("[%v] token is invalid\n", name)
		}

		for i := 1; i <= 3; i++ {
			getFamilySignPrize(serverURL, zoneToken, i)
		}

	}
	return
}

// 抓特定箱子 quality-> 1=普通 2=稀有 3=传奇
func getBoxPrizeGo(uid string, quality, amount float64) {

	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, uid).Scan(&token)

	serverURL := getServerURL()
	zoneToken := getZoneToken(serverURL, token)
	helpList := getGoldMineHelpList(serverURL, zoneToken, quality)

	ids := []int64{21, 22}

	var total float64

	var all float64

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

func getZoneToken(serverURL, token string) (zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := serverURL + "/zone?cmd=enter&token=" + token + "&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now=" + now
	formData := httpGetReturnJson(URL)
	zoneToken, ok := formData["zoneToken"].(string)
	if !ok {
		log.Printf("token:%v get zoneToken err", token)
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
		}
	}
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
		food := make(map[string]interface{})
		teamLen := len(vv["myTeam"].(map[string]interface{})["robList"].([]interface{}))
		food["id"] = vv["id"].(string)
		food["myTeam"] = teamLen
		foods = append(foods, food)
	}

	return
}

func robFamilyFood(serverURL, zoneToken, foodId string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := serverURL + "/game?cmd=robFamilyFood&token=" + zoneToken + "&foodId=" + foodId + "&now=" + now
	httpGetReturnJson(url)
	return
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
func enterGoldMine(serverURL, zoneToken string, fuid interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)

	url := fmt.Sprintf("%v/game?cmd=enterGoldMine&token=%v&fuid=%v&type=0&now=%v", serverURL, zoneToken, fuid, now)
	// url := serverURL + "/game?cmd=enterGoldMine&token=" + zoneToken + "&fuid=" + fuid + "&type=0&now=" + now
	formData := httpGetReturnJson(url)
	log.Println("formData:", formData)
	return
}

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

	bossHelpList := formData["bossHelpList"].([]interface{})

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
func summonBoss(serverURL, zoneToken string) (bossID string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=summonBoss&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	boss, ok := formData["boss"].(map[string]interface{})
	if !ok {
		now = fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		url = fmt.Sprintf("%v/game?cmd=getMyBoss&token=%v&now=%v", serverURL, zoneToken, now)
		formData = httpGetReturnJson(url)
		bossID, ok = formData["boss"].(map[string]interface{})["id"].(string)
		return
	}

	bossID, ok = boss["id"].(string)
	if !ok {

		return
	}
	return
}

// 邀请BOSS
func inviteBoss(serverURL, zoneToken, bossID string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=inviteBoss&token=%v&bossID=%v&fuidList=[301807377,302691822,309433834,326941142,374289806,381909995,406378614,690708340,693419844,694068717,694981971]&now=%v", serverURL, zoneToken, bossID, now)
	httpGetReturnJson(url)
	return
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
	// 3 次 100
	attackBossAPI(serverURL, zoneToken, bossID, 3, 0, 0, 95, 100)
	// 2 次 200
	attackBossAPI(serverURL, zoneToken, bossID, 2, 0, 1, 195, 200)

}

// 大号打Boss
func attackBossByAdmin(serverURL, zoneToken, bossID string) {
	for i := 1; i <= 3; i++ {
		now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		damage := 200
		url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=0&isDouble=1&now=%v", serverURL, zoneToken, bossID, damage, now)
		formData := httpGetReturnJson(url)
		log.Println("leftHp:", formData["boss"].(map[string]interface{})["leftHp"])
		time.Sleep(time.Second * 3)
	}

	// for i := 1; i <= 1; i++ {
	// 	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	// 	damage := RangeRand(390, 400)
	// 	url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=1&isDouble=1&now=%v", serverURL, zoneToken, bossID, damage, now)
	// 	formData := httpGetReturnJson(url)
	// 	log.Println("leftHp:", formData["boss"].(map[string]interface{})["leftHp"])
	// 	time.Sleep(time.Second * 3)
	// }
}

// https://s147.11h5.com:3148/123_207_183_233/3147/

//
func getAttackEnemyList(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getAttackEnemyList&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
	return
}

// 攻击小岛
func attackIsland(serverURL, zoneToken string, flag, targetUid, building interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=attack&token=%v&type=%v&targetUid=%v&building=%v&now=%v", serverURL, zoneToken, flag, targetUid, building, now)
	formData := httpGetReturnJson(url)
	log.Printf("攻击小岛 目标uid:%v, 建筑:%v, 增加金币:%v, 额外金币:%v", targetUid, building, formData["addGold"], formData["companionGold"])
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

// https://s147.11h5.com:3148/123_206_182_64/3148//game?cmd=addDayTaskShareCnt&token=ildWvxWpQFhVT72lfWv0MKmpOmr1bzw_uj0&now=1630027796215
// https://s147.11h5.com:3148/123_206_182_64/3148//game?cmd=getSharePrize&token=ildWvxWpQFhVT72lfWv0MKmpOmr1bzw_uj0&noEnergy=0&now=1630027796215

// 分享调用接口1
func addDayTaskShareCnt(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	// https://s147.11h5.com:3148/111_231_17_85/3149//game?cmd=addDayTaskShareCnt&token=ildSoOn78OOpalaqj0NRT34KUDyoL0Of6tc&now=1630076872090
	url := fmt.Sprintf("%v/game?cmd=addDayTaskShareCnt&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 分享调用接口2
func getSharePrize(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	// https://s147.11h5.com:3148/111_231_17_85/3149//game?cmd=getSharePrize&token=ildSoOn78OOpalaqj0NRT34KUDyoL0Of6tc&noEnergy=0&now=1630076872090
	url := fmt.Sprintf("%v/game?cmd=getSharePrize&token=%v&noEnergy=1&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

func getSharePrize0(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getSharePrize&token=%v&noEnergy=0&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// https://s147.11h5.com:3148/111_231_17_85/3149//game?cmd=stat&token=ildSoOn78OOpalaqj0NRT34KUDyoL0Of6tc&item=%E5%88%86%E4%BA%AB%E8%AF%AD%E8%AE%B0%E5%BD%95&subitem=55&now=1630076872090

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

// https://s147.11h5.com:3147/115_159_98_146/3150/

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
	log.Println("getCandyTreeInfo url:", url)
	formData := httpGetReturnJson(url)

	log.Println("getCandyTreeInfo formData:", formData)
	candyTree, ok := formData["candyTree"].(map[string]interface{})
	log.Println("getCandyTreeInfo candyTree:", candyTree)

	if !ok {
		return
	}

	candyBoxes, ok := candyTree["candyBoxes"].([]interface{})
	log.Println("getCandyTreeInfo candyBoxes:", candyBoxes)

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
	httpGetReturnJson(url)

}

// 获取邮件列表
func getMailList(serverURL, zoneToken string) (mailids []string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getMailList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	mailList, ok := formData["mailList"].([]interface{})

	if !ok {
		return
	}

	for _, v := range mailList {
		mailid := v.(map[string]interface{})["mailid"].(string)
		mailids = append(mailids, mailid)
	}
	return
}

// 领取邮件奖励
func getMailAttachments(serverURL, zoneToken, mailid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getMailAttachments&token=%v&mailid=%v&now=%v", serverURL, zoneToken, mailid, now)
	httpGetReturnJson(url)
	return
}

// 删除邮件
func deleteMail(serverURL, zoneToken, mailid string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=deleteMail&token=%v&mailid=%v&now=%v", serverURL, zoneToken, mailid, now)
	httpGetReturnJson(url)
	return
}

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
	return
}

func setPiece(serverURL, zoneToken string, id int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=setPiece&token=%v&id=%v&type=0&common=0&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
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

	id, ok := formData["id"].(string)

	if !ok {
		log.Printf("[%v]没有能量可摇！:%v", userName, formData)

		shareAPI(serverURL, zoneToken)
		getShareDrawDice(serverURL, zoneToken)
		time.Sleep(time.Second * 1)
		shareAPI(serverURL, zoneToken)
		getDayShareGold(serverURL, zoneToken)
		return 0
	}

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
		stealResult := steal(serverURL, zoneToken, idx)
		log.Printf("[%v]【摇一摇】偷取 结果:%v", userName, stealResult)
	} else if id == "3" {
		followCompanion(serverURL, zoneToken, 4)
		if uid == "302691822" || uid == "309392050" {
			time.Sleep(time.Second * 1)

			confKey := "attackIslandUid2"
			if uid == "302691822" {
				confKey = "attackIslandUid"
			}

			var attackIslandUid int64
			var attackIslandName string
			SQL := fmt.Sprintf("select conf_value, (select name from tokens where id = conf_value) as name from config where conf_key = '%s'", confKey)
			Pool.QueryRow(SQL).Scan(&attackIslandUid, &attackIslandName)

			attackIsland(serverURL, zoneToken, 1, attackIslandUid, 1)
			log.Printf("[%v]【摇一摇】攻击好友【%v】", userName, attackIslandName)
			rebuild(attackIslandUid, 1)
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
			attackIsland(serverURL, zoneToken, 0, attackUid, building)
			log.Printf("[%v]【摇一摇】攻击", userName)
		}
		followCompanion(serverURL, zoneToken, 2)

	}

	if formData["getRichmanDice"].(float64) == 1 {
		time.Sleep(time.Second * 1)
		shareAPI(serverURL, zoneToken)
		getShareDrawDice(serverURL, zoneToken)
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

	gold, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", formData["gold"].(float64)/1000000), 64)

	energy := formData["energy"].(float64)

	if energy == 0 {
		log.Printf("[%v]【摇一摇】获取免费20能量", userName)
		getFreeEnergy(serverURL, zoneToken)
	}

	log.Printf("[%v]转盘行为:%v, 剩余能量:%v, 当前金币:%vM", userName, id, energy, gold)
	return energy
}

func rebuild(uid, building int64) {
	SQL := "select token from tokens where id = ?"
	var token string
	Pool.QueryRow(SQL, uid).Scan(&token)
	serverURL := getServerURL()
	zoneToken := getZoneToken(serverURL, token)
	buildFix(serverURL, zoneToken, building)
	// buildUp(serverURL, zoneToken, building)
}

func buildUp(serverURL, zoneToken string, id int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=buildUp&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

func buildFix(serverURL, zoneToken string, id int64) {

	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=buildFix&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

// 1=大队长 2=偷钱 4=打建筑 3=防偷钱
func followCompanion(serverURL, zoneToken string, id int64) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=followCompanion&token=%v&id=%v&now=%v", serverURL, zoneToken, id, now)
	httpGetReturnJson(url)
}

// https://s147.11h5.com:3148/115_159_203_142/3147/

// https://s147.11h5.com:3147/118_89_198_182/3149/

// https://s147.11h5.com:3147/118_89_198_182/3149/

// 【摇一摇】--分享获取金币
func getDayShareGold(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getDayShareGold&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
}

// 【摇一摇】--偷取
func steal(serverURL, zoneToken string, idx interface{}) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=steal&token=%v&idx=%v&now=%v", serverURL, zoneToken, idx, now)
	formData := httpGetReturnJson(url)
	addGold := formData["addGold"].(float64)
	stealResult := formData["stealResult"].([]interface{})

	for _, v := range stealResult {
		vv := v.(map[string]interface{})
		if vv["gold"].(float64) == addGold {
			if vv["isRich"].(float64) == 1 {
				return true
			}
			return false
		}
	}

	return false
}

// 【摇一摇】--分享获取乐园骰子
func getShareDrawDice(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getShareDrawDice&token=%v&share=1&itemId=75&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)
	log.Println("当前骰子数量:", formData["richmanDice"])
}

// https://s147.11h5.com:3148/123_207_183_233/3150//game?cmd=exchangeXmas&token=ilddVUBelY_j5UCiqf_gJav9f3-2v4C4Nau&id=5&now=1630768773820

// https://s147.11h5.com:3147/123_206_192_93/3148/

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

// 领取海浪

// https://s147.11h5.com:3147/118_89_154_91/3149/

// 统一下单接口
func gameAPI(method, url, params string) (formData map[string]interface{}) {
	url = url + "/game?" + params
	if method == "GET" {
		httpGetReturnJson(url)
		formData = httpGetReturnJson(url)
		return
	}
	return
}

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

	request.Header.Add("X-Auth-Token", authToken)
	request.Header.Add("X-User-Id", id)
	request.Header.Add("Content-type", "application/json")
	client.Do(request)

	return
}

// runner

func RunnerDraw() (err error) {

	if runnerStatus("drawStatus") == "0" {
		return
	}

	hour := time.Now().Hour()
	if hour == 1 || hour == 4 || hour == 7 || hour == 10 || hour == 13 || hour == 16 || hour == 19 || hour == 22 {
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
			user.ZoneToken = getZoneToken(user.ServerURL, user.Token)

			users = append(users, user)

		}
		for _, u := range users {

			goName := u.Name
			goUid := u.Uid
			goServerURL := u.ServerURL
			goZoneToken := u.ZoneToken
			go func() {
				log.Printf("---------------------------[%v]开始转盘---------------------------", goName)
				followCompanion(goServerURL, goZoneToken, 2)
				amount := draw(goUid, goName, goServerURL, goZoneToken, 1)
				time.Sleep(time.Millisecond * 2100)
				for i := 0; i <= int(amount); i++ {
					draw(goUid, goName, goServerURL, goZoneToken, 1)
					time.Sleep(time.Millisecond * 2100)
				}
				if goUid == "302691822" {
					followCompanion(goServerURL, goZoneToken, 3)
				} else {
					followCompanion(goServerURL, goZoneToken, 1)
				}
				log.Printf("---------------------------[%v]结束转盘---------------------------", goName)

			}()
		}

		if hour == 4 {
			SQL = "select id, token, name from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
			var cowBoyUid, cowBoyToken, cowBoyName string
			Pool.QueryRow(SQL).Scan(&cowBoyUid, &cowBoyToken, &cowBoyName)
			serverURL := getServerURL()
			zoneToken := getZoneToken(serverURL, cowBoyToken)

			go func() {
				log.Printf("---------------------------[%v]开始转盘---------------------------", cowBoyName)
				followCompanion(serverURL, zoneToken, 2)
				for i := 0; i < 20; i++ {
					draw(cowBoyUid, cowBoyName, serverURL, zoneToken, 1)
					time.Sleep(time.Millisecond * 2500)
				}
				followCompanion(serverURL, zoneToken, 3)
				log.Printf("---------------------------[%v]结束转盘---------------------------", cowBoyName)

			}()

		}

	} else {
		log.Println("no draw")
	}

	return
}

func RunnerPullAnimal() (err error) {

	// 7 10 13 16 19 22 点执行

	if runnerStatus("pullAnimalGoStatus") == "0" {
		return
	}

	hour := time.Now().Hour()

	if hour >= 7 && hour <= 22 {
		log.Println("30秒后开始拉动物...")
		time.Sleep(time.Second * 30)
		log.Println("现在开始拉动物")

		SQL := "select id, token, name from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
		var uid, token, name string

		err = Pool.QueryRow(SQL).Scan(&uid, &token, &name)

		if err != nil {
			return
		}

		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid + ":" + name)
			log.Printf("[name: %v] token is invalid\n", name)
		}

		foods := enterFamilyRob(serverURL, zoneToken)

		for _, v := range foods {
			myTeam := v["myTeam"].(int)
			if myTeam != 4 {
				robFamilyFood(serverURL, zoneToken, v["id"].(string))
				break
			}

		}
		log.Printf("[%v]拉动物完成", name)
		// time.Sleep(time.Second * 1)
		// log.Printf("cowboy serverURL:%v, zoneToken:%v\n", serverURL, zoneToken)

		pullAnimalGo()
		return
	}

	log.Println("RunnerPullAnimal do not run")

	return
}

func runnerStatus(confKey string) (confValue string) {
	Pool.QueryRow("select conf_value from config where conf_key = ?", confKey).Scan(&confValue)
	return
}
