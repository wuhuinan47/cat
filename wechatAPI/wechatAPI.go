package wechatapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type WechatAccessToken struct {
	AccessToken string
	ExpiryTime  time.Time
}

var wechatAccessToken WechatAccessToken

func Init() {
	getAccessToken()
}

var openID = "o9od753rqL522SlZkuIEc3NVHBKA"
var openID2 = "o9od75ydpNatcX-e0RMnXtnP8NDM"

func SendMsgH(w http.ResponseWriter, req *http.Request) {

	qrcode := req.URL.Query().Get("qrcode")

	url := "https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=" + getAccessToken()
	var body = make(map[string]interface{})

	body["touser"] = openID
	body["msgtype"] = "news"
	body["news"] = map[string]interface{}{
		"articles": []interface{}{
			map[string]interface{}{
				"title":       "Happy Day",
				"description": "Is Really A Happy Day",
				"url":         qrcode,
				"picurl":      qrcode,
			},
		},
	}
	result, _ := httpPostJson(url, body)

	msg, _ := json.Marshal(result)

	log.Println("Msg:", string(msg))

	w.Write(msg)

}

func getAccessToken() (accessToken string) {

	now := time.Now()
	if wechatAccessToken.ExpiryTime.After(now) {
		accessToken = wechatAccessToken.AccessToken
		log.Println("1accessToken is ", accessToken)
		return
	}

	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx5a34bbd4909c1a37&secret=1cffa17bd60fba501628cc86a20fdcc3"

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
	formData := make(map[string]interface{})
	json.NewDecoder(response.Body).Decode(&formData)
	accessToken, ok := formData["access_token"].(string)
	log.Println("2accessToken is ", accessToken)

	if !ok {
		return
	}
	expiresIn, ok := formData["expires_in"].(float64)

	wechatAccessToken.AccessToken = accessToken
	wechatAccessToken.ExpiryTime = time.Now().Add(time.Second * time.Duration(expiresIn))

	return
}

func httpPostJson(url string, body map[string]interface{}) (result io.ReadCloser, err error) {
	client := &http.Client{}

	sendMsg, err := json.Marshal(body)

	if err != nil {
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(sendMsg))
	if err != nil {
		log.Printf("http post err is %v, url is %v", err, url)
		return
	}
	request.Header.Add("Content-type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		log.Printf("http post err is %v, url is %v", err, url)
		return
	}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(&result)
	return

}
