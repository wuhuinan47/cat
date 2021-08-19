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
	"os/exec"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/wuhuinan47/cat/runner"
)

var Pool *sql.DB

// type CatData struct {
// 	GoldMineHelpList []GoldMineHelpList `json:"helpList"`
// }

// type GoldMineHelpList struct {
// 	Uid     int64 `json:"uid"`
// 	Quality int64 `json:"quality"`
// }

func main() {

	db, err := sql.Open("mysql", "cat:cyydmkj123@tcp(localhost:3306)/data_cat?charset=utf8&parseTime=true&collation=utf8mb4_unicode_ci")
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

	http.HandleFunc("/getServerURL", GetServerURLH)
	http.HandleFunc("/getZoneToken", GetZoneTokenH)

	//
	http.HandleFunc("/checkToken", CheckTokenH)

	//

	http.HandleFunc("/index.html", IndexH)
	http.HandleFunc("/", IndexH)

	sendMsg("cat is start")

	running := true
	{
		// go admindb.NamedRunnerWithHMS("ProcPlayerCashiers", 10, 00, 0, &running, admindb.ProcPlayerCashiers)
		go runner.NamedRunnerWithSeconds("ProcPlayerCashiers", 3600, &running, RunnerPullAnimal)

		// go admindb.NamedRunnerWithSeconds("ProcPlayerCashiers", 3600, &running, admindb.ProcPlayerCashiers)
	}

	log.Println("start port 33333 sucesss")
	http.ListenAndServe(":33333", nil)

}

func IndexH(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("ctrl.html")
	t.Execute(w, nil)
}

func SendQrcodeH(w http.ResponseWriter, req *http.Request) {
	qrcode := req.URL.Query().Get("qrcode")
	sendMsg(qrcode)
	log.Println("qrcode:", qrcode)
	io.WriteString(w, qrcode)
	return
}

func LoginByQrcodeH(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	str := "/mnt/app/cat/linuxWechat.py"

	if id == "2" {
		str = "/mnt/app/cat/linuxQQ.py"
	}

	go func() {
		cmd := exec.Command("/bin/python3", str)
		err := cmd.Run()
		if err != nil {
			return
		}
	}()

	io.WriteString(w, "SUCCESS")
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
		SQL := "select id, token from tokens"

		rows, err := Pool.Query(SQL)

		if err != nil {
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id, token string
			rows.Scan(&id, &token)

			serverURL := getServerURL()
			zoneToken := getZoneToken(serverURL, token)

			if zoneToken == "" {
				Pool.Exec("update tokens set token = '' where id = ?", id)
				sendMsg(id)
			}
			time.Sleep(time.Second * 2)

		}
		return
	}()
	io.WriteString(w, "SUCCESS")
	return

}

func DiamondH(w http.ResponseWriter, req *http.Request) {

	id := req.URL.Query().Get("id")
	quality := req.URL.Query().Get("quality")

	v2, _ := strconv.ParseFloat(quality, 64)

	log.Println("id:", id)
	log.Println("quality:", quality)
	log.Println("v2:", v2)
	getBoxPrizeGo(id, v2)

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

	log.Println("url:", url)
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

			time.Sleep(time.Second * 2)
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
	// uids := `301807377, 302691822, 309392050, 309433834, 374289806, 375912362, 382292124,
	// 	385498006, 403573789, 406961861, 408385382, 410572648, 425190502, 439943689, 440204933,
	// 	441013912, 444729748, 446399085, 693419844, 694981971`

	SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var uid, token string
		rows.Scan(&uid, &token)
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}

		foods := enterFamilyRob(serverURL, zoneToken)

		for _, v := range foods {
			myTeam := v["myTeam"].(int)
			if myTeam != 4 {
				robFamilyFood(serverURL, zoneToken, v["id"].(string))
				break
			}

		}
		// time.Sleep(time.Second * 1)
		log.Printf("serverURL:%v, zoneToken:%v\n", serverURL, zoneToken)
	}
	return
}

// 一键公会签到/签到/免费福利/免费转盘/免费夹子/赠送能量
func familySignGo() {
	log.Println("familySignGo start...")

	SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var uid, token string
		rows.Scan(&uid, &token)
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}

		log.Printf("[uid:%v] start familySign", uid)
		familySign(serverURL, zoneToken)
		log.Printf("[uid:%v] end familySign", uid)
		time.Sleep(time.Second * 1)
		log.Printf("[uid:%v] start getSignPrize", uid)
		getSignPrize(serverURL, zoneToken)
		log.Printf("[uid:%v] end getSignPrize", uid)
		time.Sleep(time.Second * 1)
		log.Printf("[uid:%v] start getFreeDailyGiftBox", uid)
		getFreeDailyGiftBox(serverURL, zoneToken)
		log.Printf("[uid:%v] end getFreeDailyGiftBox", uid)
		time.Sleep(time.Second * 1)

		log.Printf("[uid:%v] start playLuckyWheel", uid)
		shareAPI(serverURL, zoneToken)
		playLuckyWheel(serverURL, zoneToken)
		log.Printf("[uid:%v] end playLuckyWheel", uid)
		time.Sleep(time.Second * 1)

		log.Printf("[uid:%v] start getFreeClamp", uid)
		shareAPI(serverURL, zoneToken)
		getFreeClamp(serverURL, zoneToken)
		log.Printf("[uid:%v] end getFreeClamp", uid)

		time.Sleep(time.Second * 1)
		log.Printf("[uid:%v] start autoFriendEnergy", uid)
		autoFriendEnergy(serverURL, zoneToken)
		log.Printf("[uid:%v] end autoFriendEnergy", uid)

		gameList := []string{"535", "525", "157", "452", "411"}

		time.Sleep(time.Second * 1)
		for _, v := range gameList {
			getAward(token, v)
		}
	}

	time.Sleep(time.Second * 1)
	log.Println("start getFamilySignPrizeGo..")
	getFamilySignPrizeGo()
	log.Println("end getFamilySignPrizeGo..")

	getAwardForCowBoy()

	log.Println("familySignGo finish")
	return
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
	SQL := "select id, token from tokens where find_in_set(id, (select conf_value from config where conf_key = 'animalUids'))"

	rows, err := Pool.Query(SQL)

	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var uid, token string
		rows.Scan(&uid, &token)
		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}

		for i := 1; i <= 3; i++ {
			getFamilySignPrize(serverURL, zoneToken, i)
		}

	}
	return
}

// 抓特定箱子 quality-> 1=普通 2=稀有 3=传奇
func getBoxPrizeGo(uid string, quality float64) {

	SQL := "select token from tokens where id = ?"

	var token string
	Pool.QueryRow(SQL, uid).Scan(&token)

	serverURL := getServerURL()
	zoneToken := getZoneToken(serverURL, token)
	helpList := getGoldMineHelpList(serverURL, zoneToken, quality)

	ids := []int64{21, 22}

	for _, v := range helpList {

		for _, v2 := range ids {
			time.Sleep(time.Second * 1)
			goldMineFish(serverURL, zoneToken, v["uid"], v2)
		}
	}
}

// interface functions
func getServerURL() (serverURL string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := "https://api.11h5.com/conf?cmd=getGameInfo&gameid=147&" + now
	formData := httpGetReturnJson(url)
	serverURL, ok := formData["ext"].(map[string]interface{})["serverURL"].(string)
	if !ok {
		log.Println("get serverURL err")
		return
	}
	return

}

func getZoneToken(serverURL, token string) (zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := serverURL + "/zone?cmd=enter&token=" + token + "&yyb=0&inviteId=null&share_from=null&cp_shareId=null&now=" + now
	formData := httpGetReturnJson(url)
	zoneToken, ok := formData["zoneToken"].(string)
	if !ok {
		log.Println("get zoneToken err")
		return
	}
	return
}

func enterFamilyRob(serverURL, zoneToken string) (foods []map[string]interface{}) {

	// https://s147.11h5.com:3148/123_206_182_64/3148//game?cmd=enterFamilyRob&token=ild70NWsJczwFddT65UMNWv66oNaoFZMH_h&now=1629251946267

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
func goldMineFish(serverURL, zoneToken string, fuid, id interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)

	url := fmt.Sprintf("%v/game?cmd=goldMineFish&token=%v&fuid=%v&id=%v&now=%v", serverURL, zoneToken, fuid, id, now)
	httpGetReturnJson(url)
	return
}

// 获取帮助列表
func getBossHelpList(serverURL, zoneToken string) (bossList []map[string]interface{}) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=getBossHelpList&token=%v&now=%v", serverURL, zoneToken, now)
	formData := httpGetReturnJson(url)

	bossHelpList := formData["bossHelpList"].([]interface{})

	for _, v := range bossHelpList {
		boss := v.(map[string]interface{})["boss"].(map[string]interface{})
		leftHp := boss["leftHp"].(int)

		if leftHp > 0 {
			bossList = append(bossList, boss)
		}
	}
	log.Println("bossList:", bossList)
	return

}

// 小号获取BOSS状态
func enterBoss(serverURL, zoneToken, bossID string, selfUID int) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=enterBoss&token=%v&bossID=%v&now=%v", serverURL, zoneToken, bossID, now)
	formData := httpGetReturnJson(url)
	leftHp := formData["boss"].(map[string]interface{})["leftHp"].(int)
	if leftHp > 800 {
		bossAttackList := formData["bossAttackList"].([]interface{})
		for _, v := range bossAttackList {
			uid := v.(map[string]interface{})["uid"].(int)
			damage := v.(map[string]interface{})["damage"].(int)
			if uid == selfUID {
				if damage < 600 {
					return true
				} else {
					return false
				}
			}
		}
	}
	return true

}

// 开启Boss
func summonBoss(serverURL, zoneToken string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=summonBoss&token=%v&now=%v", serverURL, zoneToken, now)
	httpGetReturnJson(url)
	return
}

// 邀请BOSS
func inviteBoss(serverURL, zoneToken, bossID string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=inviteBoss&token=%v&bossID=%v&fuidList=[301807377,302691822,309433834,326941142,374289806,381909995,406378614,690708340,693419844,694068717,694981971]&now=%v", serverURL, zoneToken, bossID, now)
	httpGetReturnJson(url)
	return
}

// 小号打Boss
func attackBoss(serverURL, zoneToken, bossID string) {
	for i := 1; i <= 3; i++ {
		now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		damage := RangeRand(95, 100)
		url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=0&isDouble=0&now=%v", serverURL, zoneToken, bossID, damage, now)
		formData := httpGetReturnJson(url)
		log.Println("leftHp:", formData["boss"].(map[string]interface{})["leftHp"])
	}

	for i := 1; i <= 2; i++ {
		now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		damage := RangeRand(195, 200)
		url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=0&isDouble=0&now=%v", serverURL, zoneToken, bossID, damage, now)
		formData := httpGetReturnJson(url)
		log.Println("leftHp:", formData["boss"].(map[string]interface{})["leftHp"])
	}
}

// 大号打Boss
func attackBossByAdmin(serverURL, zoneToken, bossID string) {
	for i := 1; i <= 3; i++ {
		now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		damage := RangeRand(197, 200)
		url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=0&isDouble=1&now=%v", serverURL, zoneToken, bossID, damage, now)
		formData := httpGetReturnJson(url)
		log.Println("leftHp:", formData["boss"].(map[string]interface{})["leftHp"])
	}

	for i := 1; i <= 1; i++ {
		now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		damage := RangeRand(390, 400)
		url := fmt.Sprintf("%v/game?cmd=attackBoss&token=%v&bossID=%v&damage=%v&isPerfect=1&isDouble=1&now=%v", serverURL, zoneToken, bossID, damage, now)
		formData := httpGetReturnJson(url)
		log.Println("leftHp:", formData["boss"].(map[string]interface{})["leftHp"])
	}
}

// 攻击小岛
func attackIsland(serverURL, zoneToken, targetUid, building string) {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	url := fmt.Sprintf("%v/game?cmd=attack&token=%v&type=1&targetUid=%v&building=%v&now=%v", serverURL, zoneToken, targetUid, building, now)
	httpGetReturnJson(url)
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

// https://libao.11h5.com/wall?cmd=getAward&token=43b1d99e84da82a99879c289b56ca2ca&gameid=157&type=1&channel=147&now=1629364564823

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

func RunnerPullAnimal() (err error) {

	// 7 10 13 16 19 22 点执行

	hour := time.Now().Hour()

	if hour == 7 || hour == 10 || hour == 13 || hour == 16 || hour == 19 || hour == 22 {
		log.Println("after 30s to pull animal")
		time.Sleep(time.Second * 30)
		log.Println("now start pull animal")

		SQL := "select id, token from tokens where id = (select conf_value from config where conf_key = 'cowBoy')"
		var uid, token string

		err = Pool.QueryRow(SQL).Scan(&uid, &token)

		if err != nil {
			return
		}

		serverURL := getServerURL()
		zoneToken := getZoneToken(serverURL, token)

		if zoneToken == "" {
			sendMsg(uid)
			log.Printf("[uid: %v] token is invalid\n", uid)
		}

		foods := enterFamilyRob(serverURL, zoneToken)

		for _, v := range foods {
			myTeam := v["myTeam"].(int)
			if myTeam != 4 {
				robFamilyFood(serverURL, zoneToken, v["id"].(string))
				break
			}

		}
		// time.Sleep(time.Second * 1)
		log.Printf("cowboy serverURL:%v, zoneToken:%v\n", serverURL, zoneToken)

		pullAnimalGo()
		return
	}

	log.Println("RunnerPullAnimal do not run")

	return
}
