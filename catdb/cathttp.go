package catdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func CheckStat(URL string) bool {
	formData := httpGetReturnJson(URL)
	_, ok := formData["error"]
	log.Println("CheckStat:", formData["error"])

	return !ok
}

func CheckZoneToken(serverURL, zoneToken string) bool {
	now := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	URL := fmt.Sprintf("%v/game?cmd=hasNewMail&token=%v&now=%v", serverURL, zoneToken, now)
	// log.Println("checkZoneToken URL is :", URL)
	formData := httpGetReturnJson(URL)
	_, ok := formData["newMailCnt"].(float64)
	return ok
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
