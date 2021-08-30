package crawler

import (
	"fmt"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {

	prefix := "K86TF"

	str := "K86TFtestobsx1001"
	fmt.Println(strings.HasPrefix(str, prefix))

	fmt.Println(str[len(prefix) : len(str)-0])

	// seleniumPath := "/usr/local/bin/chromeDriver"
	// port := 9515

	// //1.开启selenium服务
	// //设置selenium服务的选项,设置为空。根据需要设置。
	// ops := []selenium.ServiceOption{}
	// service, err := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	// if err != nil {
	// 	fmt.Printf("Error starting the ChromeDriver server: %v", err)
	// }
	// //延迟关闭服务
	// defer service.Stop()

	// //2.调用浏览器实例
	// //设置浏览器兼容性，我们设置浏览器名称为chrome
	// caps := selenium.Capabilities{
	// 	"browserName": "chrome",
	// }
	// //调用浏览器urlPrefix: 测试参考：DefaultURLPrefix = "http://127.0.0.1:4444/wd/hub"
	// wd, err := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	// if err != nil {
	// 	panic(err)
	// }
	// //延迟退出chrome
	// defer wd.Quit()

	// // 3单选radio，多选checkbox，select框操作(功能待完善，https://github.com/tebeka/selenium/issues/141)
	// if err := wd.Get("http://cdn1.python3.vip/files/selenium/test2.html"); err != nil {
	// 	panic(err)
	// }
	// //3.1操作单选radio
	// we, err := wd.FindElement(selenium.ByCSSSelector, `#s_radio > input[type=radio]:nth-child(3)`)
	// if err != nil {
	// 	panic(err)
	// }
	// we.Click()

	// //3.2操作多选checkbox
	// //删除默认checkbox
	// we, err = wd.FindElement(selenium.ByCSSSelector, `#s_checkbox > input[type=checkbox]:nth-child(5)`)
	// if err != nil {
	// 	panic(err)
	// }
	// we.Click()
	// //选择选项
	// we, err = wd.FindElement(selenium.ByCSSSelector, `#s_checkbox > input[type=checkbox]:nth-child(1)`)
	// if err != nil {
	// 	panic(err)
	// }
	// we.Click()
	// we, err = wd.FindElement(selenium.ByCSSSelector, `#s_checkbox > input[type=checkbox]:nth-child(3)`)
	// if err != nil {
	// 	panic(err)
	// }
	// we.Click()

	// //3.3 select多选
	// //删除默认选项

	// //选择默认项
	// we, err = wd.FindElement(selenium.ByCSSSelector, `#ss_multi > option:nth-child(3)`)
	// if err != nil {
	// 	panic(err)
	// }
	// we.Click()

	// we, err = wd.FindElement(selenium.ByCSSSelector, `#ss_multi > option:nth-child(2)`)
	// if err != nil {
	// 	panic(err)
	// }
	// we.Click()

	// //睡眠20秒后退出
	// time.Sleep(20 * time.Second)
}
