package crawler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"github.com/wuhuinan47/cat/catdb"
)

// func CatRun() {
// 	seleniumPath := "/usr/bin/chromedriver"
// 	port := 49515
// 	ops := []selenium.ServiceOption{}
// 	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
// 	if err != nil {
// 		fmt.Printf("Error starting the ChromeDriver server: %v", err)
// 		return
// 	}
// 	imagCaps := map[string]interface{}{

// 		"profile.managed_default_content_settings.images": 2,
// 	}
// 	chromeCaps := chrome.Capabilities{

// 		Prefs: imagCaps,

// 		Path: "",

// 		Args: []string{

// 			"--headless", // 设置Chrome无头模式

// 			"--no-sandbox",

// 			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36", // 模拟user-agent，防反爬

// 		},
// 	}

// 	//延迟关闭服务
// 	defer service.Stop()

// 	caps := selenium.Capabilities{
// 		"browserName": "chrome",
// 	}
// 	caps.AddChrome(chromeCaps)
// 	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
// 	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:49515/wd/hub")
// 	if err != nil {
// 		log.Println("CatRun err :", err)
// 		return
// 		// panic(err)
// 	}
// 	//延迟退出chrome
// 	defer wd.Quit()

// 	var globalURL = "https://play.h5avu.com/game/?gameid=147&token="
// 	var mtoken string

// 	catdb.Pool.QueryRow("select token from tokens where id = 302691822").Scan(&mtoken)

// 	URL := globalURL + mtoken

// 	if err = wd.Get(URL); err != nil {
// 		log.Println("get url:", err)
// 		return
// 	}

// 	var updateToken string

// 	for i := 0; i < 100; i++ {
// 		ygToken, _ = wd.ExecuteScript("return localStorage.getItem('yg_token')", nil)
// 		userID, _ = wd.ExecuteScript("return localStorage.getItem('__TD_userID')", nil)
// 		if ygToken != nil && userID != nil {
// 			log.Println("获取到token is", ygToken, " uid is", userID)
// 			break
// 		}
// 		log.Println("qq Scan正在获取token...")

// 		if i == 30 {
// 			imgBytes, _ := wd.Screenshot()
// 			ioutil.WriteFile("./static/qqQrCode.png", imgBytes, 0644)
// 		}
// 		time.Sleep(time.Second * 1)
// 	}

// 	log.Printf("login success")

// 	for {

// 		var crawlerStatus string
// 		catdb.Pool.QueryRow("select conf_value from config where conf_key = 'crawlerStatus'").Scan(&crawlerStatus)

// 		if crawlerStatus != "1" {
// 			log.Println("crawlerStatus stop")
// 			return
// 		}

// 		catdb.Pool.QueryRow("select token from tokens where id = 302691822")
// 	}

// }

func QQScan() (ygToken, userID interface{}, err error) {
	seleniumPath := "/usr/bin/chromedriver"
	port := 9515

	//1.开启selenium服务
	//设置selenium服务的选项,设置为空。根据需要设置。
	ops := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	imagCaps := map[string]interface{}{

		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{

		Prefs: imagCaps,

		Path: "",

		Args: []string{

			"--headless", // 设置Chrome无头模式

			"--no-sandbox",

			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36", // 模拟user-agent，防反爬

		},
	}

	//延迟关闭服务
	defer service.Stop()

	//2.调用浏览器实例
	//设置浏览器兼容性，我们设置浏览器名称为chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	caps.AddChrome(chromeCaps)
	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		log.Println("qq scan new remote err")
		return
		// panic(err)
	}
	//延迟退出chrome
	defer wd.Quit()

	// 3单选radio，多选checkbox，select框操作(功能待完善，https://github.com/tebeka/selenium/issues/141)
	// if err := wd.Get("https://open.weixin.qq.com/connect/qrconnect?appid=wx22f69b39568e9cb3&redirect_uri=http%3A%2F%2Flogin.11h5.com%2Faccount%2Fapi.php%3Fc%3Dwxlogin%26d%3DwxQrcodeAuth%26pf%3Dwxqrcode%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147%2526fuid%253D302691822%2526statid%253D1785%2526share_from%253Dmsg%2526cp_from%253Dmsg%2526cp_shareId%253D55&response_type=code&scope=snsapi_login&state=#wechat_redirect"); err != nil {
	// 	panic(err)
	// }

	URL := "https://graph.qq.com/oauth2.0/show?which=Login&display=pc&response_type=token&client_id=101206450&state=&redirect_uri=http%3A%2F%2Flogin.vutimes.com%2Faccount%2Fpage%2FqqAuthCallback.html%3FswitchVersion%3D1%26pf%3Dqq%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147%2526fuid%253D302691822%2526statid%253D1785%2526share_from%253Dmsg%2526cp_from%253Dmsg%2526cp_shareId%253D55"

	if err = wd.Get(URL); err != nil {
		log.Println("get url:", err)
		return
	}

	imgBytes, err := wd.Screenshot()
	if err != nil {
		log.Println("Screenshot err:", err)
		return
	}

	if err = ioutil.WriteFile("./static/qqQrCode.png", imgBytes, 0644); err != nil {
		log.Println("WriteFile:", err)
		return
	}

	if err = ioutil.WriteFile("./www/cat_demo/qrcode/qqQrCode.png", imgBytes, 0644); err != nil {
		log.Println("WriteFile:", err)
		// return
	}

	for i := 0; i < 100; i++ {
		ygToken, _ = wd.ExecuteScript("return localStorage.getItem('yg_token')", nil)
		userID, _ = wd.ExecuteScript("return localStorage.getItem('__TD_userID')", nil)
		if ygToken != nil && userID != nil {
			log.Println("获取到token is", ygToken, " uid is", userID)
			break
		}
		log.Println("qq Scan正在获取token...")

		if i == 30 {
			imgBytes, _ := wd.Screenshot()
			ioutil.WriteFile("./static/qqQrCode.png", imgBytes, 0644)
		}
		time.Sleep(time.Second * 1)
	}

	// http.Get(fmt.Sprintf("https://mcps.51yizhuan.com:13010/update?id=%v&token=%v", userID, ygToken))

	time.Sleep(time.Second * 1)

	return
}

func WechatScan() (ygToken, userID interface{}, err error) {
	seleniumPath := "/usr/bin/chromedriver"
	port := 39515

	//1.开启selenium服务
	//设置selenium服务的选项,设置为空。根据需要设置。
	ops := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}

	imagCaps := map[string]interface{}{

		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{

		Prefs: imagCaps,

		Path: "",

		Args: []string{

			"--headless", // 设置Chrome无头模式

			"--no-sandbox",

			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36", // 模拟user-agent，防反爬

		},
	}

	//延迟关闭服务
	defer service.Stop()

	//2.调用浏览器实例
	//设置浏览器兼容性，我们设置浏览器名称为chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	caps.AddChrome(chromeCaps)
	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:39515/wd/hub")
	if err != nil {
		log.Println("qq scan new remote err")
		return
		// panic(err)
	}
	//延迟退出chrome
	defer wd.Quit()

	// 3单选radio，多选checkbox，select框操作(功能待完善，https://github.com/tebeka/selenium/issues/141)
	// if err := wd.Get("https://open.weixin.qq.com/connect/qrconnect?appid=wx22f69b39568e9cb3&redirect_uri=http%3A%2F%2Flogin.11h5.com%2Faccount%2Fapi.php%3Fc%3Dwxlogin%26d%3DwxQrcodeAuth%26pf%3Dwxqrcode%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147%2526fuid%253D302691822%2526statid%253D1785%2526share_from%253Dmsg%2526cp_from%253Dmsg%2526cp_shareId%253D55&response_type=code&scope=snsapi_login&state=#wechat_redirect"); err != nil {
	// 	panic(err)
	// }

	URL := "https://open.weixin.qq.com/connect/qrconnect?appid=wx22f69b39568e9cb3&redirect_uri=http%3A%2F%2Flogin.11h5.com%2Faccount%2Fapi.php%3Fc%3Dwxlogin%26d%3DwxQrcodeAuth%26pf%3Dwxqrcode%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147%2526fuid%253D302691822%2526statid%253D1785%2526share_from%253Dmsg%2526cp_from%253Dmsg%2526cp_shareId%253D55&response_type=code&scope=snsapi_login&state=#wechat_redirect"

	if err = wd.Get(URL); err != nil {
		log.Println("get url:", err)
		return
	}
	log.Println("wd get ", err)

	webElement, err := wd.FindElement("xpath", "/html/body/div[1]/div/div/div[2]/div[1]/img")
	log.Println("FindElement err:", err)

	if err != nil {
		log.Println("FindElement err:", err)
		return
	}

	qrcode, err := webElement.GetAttribute("src")
	log.Println("GetAttribute qrcode", qrcode, err)

	if err != nil {
		log.Println("GetAttribute err:", err)
		return
	}

	catdb.Pool.Exec("update config set conf_value = ? where conf_key = 'wechatLoginQrcode'", qrcode)

	// qrcode=self.browser.find_element_by_xpath("/html/body/div[1]/div/div/div[2]/div[1]/img").get_attribute('src')

	// imgBytes, err := wd.Screenshot()
	// if err != nil {
	// 	log.Println("Screenshot err:", err)
	// 	return
	// }

	// if err = ioutil.WriteFile("./www/cat_demo/qrcode/wechatQrCode.png", imgBytes, 0644); err != nil {
	// 	log.Println("WriteFile:", err)
	// 	return
	// }

	for i := 0; i < 100; i++ {
		ygToken, _ = wd.ExecuteScript("return localStorage.getItem('yg_token')", nil)
		userID, _ = wd.ExecuteScript("return localStorage.getItem('__TD_userID')", nil)
		if ygToken != nil && userID != nil {
			log.Println("获取到token is", ygToken, " uid is", userID)
			break
		}
		log.Println("wechat Scan正在获取token...")

		if i == 30 {
			imgBytes, _ := wd.Screenshot()
			ioutil.WriteFile("./www/cat_demo/qrcode/wechatQrCode.png", imgBytes, 0644)
		}
		time.Sleep(time.Second * 1)
	}

	// http.Get(fmt.Sprintf("https://mcps.51yizhuan.com:13010/update?id=%v&token=%v", userID, ygToken))

	time.Sleep(time.Second * 1)

	return
}

func checkURL() chromedp.ActionFunc {
	return func(ctx context.Context) (err error) {

		for {
			time.Sleep(time.Second * 3)
			if globalCheckURL != "" {
				// log.Println("globalCheckURL:", globalCheckURL)
				globalFailCount = 0
				if catdb.CheckStat(globalCheckURL) {
					if !globalCheckBool {
						// log.Println("globalCheckURL:", globalCheckURL)

						sL := strings.Split(globalCheckURL, "?")
						// log.Println("sL:", sL)
						var zoneToken string
						sL1 := strings.Split(sL[1], "&")
						// log.Println("sL1:", sL1)

						for _, v2 := range sL1 {
							sL2 := strings.Split(v2, "=")
							// log.Println("sL2:", sL2)

							if sL2[0] == "token" {
								zoneToken = sL2[1]
								break
							}
						}
						log.Println("zoneToken:", zoneToken)

						catdb.Pool.Exec("update tokens set zoneToken = ? where id = 301807377", zoneToken)
					}
					globalCheckBool = true
				}

			} else {
				globalFailCount++
			}
		}
	}

}

func DemoChromedp(URL string) {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		// chromedp.Flag("window-size", "50,400"),
		chromedp.UserDataDir(dir),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	taskCtx, cancel = context.WithCancel(taskCtx)
	defer cancel()

	// ensure that the browser process is started
	if err := chromedp.Run(taskCtx); err != nil {
		panic(err)
	}

	// listen network event
	listenForNetworkEvent(taskCtx)

	for {
		globalCheckURL = ""
		chromedp.Run(taskCtx,
			network.Enable(),
			chromedp.Navigate(URL),
			// saveCookies(),
			// chromedp.MouseClickXY(150, 250),
			checkURL(),
			// chromedp.WaitVisible(`body`, chromedp.BySearch),
		)

		if globalFailCount > 10 {
			return
		}
	}
}

var globalCheckURL = ""

var globalFailCount = 0

var globalCheckBool = false

func listenForNetworkEvent(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch ev := ev.(type) {

		case *network.EventResponseReceived:
			resp := ev.Response
			if len(resp.Headers) != 0 {

				if strings.Contains(resp.URL, "s147.11h5.com//game?") {
					globalCheckURL = resp.URL
					log.Printf("received URL: %v", resp.URL)

				}

			}

		case *network.EventRequestWillBeSent:
			if strings.Contains(ev.Request.URL, "s147.11h5.com//game?") {
				log.Printf("EventRequestWillBeSent: %v", ev.Request.URL)
			}

		// case *network.EventWebSocketFrameReceived:
		// 	log.Printf("EventWebSocketFrameReceived:%v", ev.Response)
		// case *network.EventWebSocketFrameSent:
		// 	log.Printf("EventWebSocketFrameSent:%v", ev.Response)
		// case *network.EventWebSocketHandshakeResponseReceived:
		// 	log.Printf("EventWebSocketHandshakeResponseReceived:%v", ev.Response)
		// case *network.EventWebSocketWillSendHandshakeRequest:
		// 	log.Printf("EventWebSocketWillSendHandshakeRequest:%v", ev.Request)
		// case *network.EventSubresourceWebBundleInnerResponseParsed:
		// 	log.Printf("EventSubresourceWebBundleInnerResponseParsed:%v", ev.InnerRequestURL)
		default:
			// log.Printf("received ev: %v", ev)

		}

		// other needed network Event
	})
}

// // 保存Cookies
// func saveCookies() chromedp.ActionFunc {
// 	return func(ctx context.Context) (err error) {

// 		var cookiesData []byte
// 		for {
// 			cookiesData = getCookies(ctx)
// 			if string(cookiesData) != "{}" {
// 				log.Println("cookiesData:", cookiesData)
// 				return
// 			}
// 		}
// 	}
// }

// func getCookies(ctx context.Context) (cookiesData []byte) {
// 	cookies, err := network.GetAllCookies().Do(ctx)
// 	if err != nil {
// 		return
// 	}
// 	// 2. 序列化
// 	cookiesData, err = network.GetAllCookiesReturns{Cookies: cookies}.MarshalJSON()
// 	if err != nil {
// 		return
// 	}
// 	return
// }
