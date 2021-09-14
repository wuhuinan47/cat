package crawler

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestExec(t *testing.T) {

	seleniumPath := "/usr/local/bin/chromeDriver"
	port := 9515

	//1.开启selenium服务
	//设置selenium服务的选项,设置为空。根据需要设置。
	ops := []selenium.ServiceOption{}

	service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	if err != nil {
		fmt.Printf("Error starting the ChromeDriver server: %v", err)
	}
	//延迟关闭服务
	defer service.Stop()

	//2.调用浏览器实例
	//设置浏览器兼容性，我们设置浏览器名称为chrome
	// options.add_experimental_option('excludeSwitches',
	//                                     ['enable-automation'])
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	//调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	if err != nil {
		panic(err)
	}
	//延迟退出chrome
	defer wd.Quit()

	// 3单选radio，多选checkbox，select框操作(功能待完善，https://github.com/tebeka/selenium/issues/141)
	// if err := wd.Get("https://open.weixin.qq.com/connect/qrconnect?appid=wx22f69b39568e9cb3&redirect_uri=http%3A%2F%2Flogin.11h5.com%2Faccount%2Fapi.php%3Fc%3Dwxlogin%26d%3DwxQrcodeAuth%26pf%3Dwxqrcode%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147%2526fuid%253D302691822%2526statid%253D1785%2526share_from%253Dmsg%2526cp_from%253Dmsg%2526cp_shareId%253D55&response_type=code&scope=snsapi_login&state=#wechat_redirect"); err != nil {
	// 	panic(err)
	// }

	URL := "https://graph.qq.com/oauth2.0/show?which=Login&display=pc&response_type=token&client_id=101206450&state=&redirect_uri=http%3A%2F%2Flogin.vutimes.com%2Faccount%2Fpage%2FqqAuthCallback.html%3FswitchVersion%3D1%26pf%3Dqq%26ssl%3D1%26back_url%3Dhttps%253A%252F%252Fplay.h5avu.com%252Fgame%252F%253Fgameid%253D147%2526fuid%253D302691822%2526statid%253D1785%2526share_from%253Dmsg%2526cp_from%253Dmsg%2526cp_shareId%253D55"

	fmt.Println("123")
	if err := wd.Get(URL); err != nil {
		fmt.Println("err1:", err)
		panic(err)
	}

	time.Sleep(time.Second * 2)

	// imgBytes, err := wd.Screenshot()

	// body := bytes.NewReader(imgBytes)
	// client := &http.Client{}
	// request, _ := http.NewRequest("POST", "https://mcps.51yizhuan.com:13010/sendQQQrcode", body)
	// _, err = client.Do(request)

	var ygToken, userID interface{}

	for i := 0; i < 100; i++ {
		ygToken, _ = wd.ExecuteScript("return localStorage.getItem('yg_token')", nil)

		fmt.Println(ygToken)

		userID, _ = wd.ExecuteScript("return localStorage.getItem('__TD_userID')", nil)

		fmt.Println(userID)

		if ygToken != nil && userID != nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	http.Get(fmt.Sprintf("https://mcps.51yizhuan.com:13010/update?id=%v&token=%v", userID, ygToken))

	time.Sleep(time.Second * 1)

}
