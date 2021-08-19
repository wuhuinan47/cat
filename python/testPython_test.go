package python

import (
	"os/exec"
	"testing"
)

func TestExec(t *testing.T) {
	err := CmdPythonSaveImageDpi()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("转换成功")
}

func CmdPythonSaveImageDpi() (err error) {
	cmd := exec.Command("/usr/local/bin/python3", "/Users/wuhuinan/go/src/github.com/wuhuinan47/cat/tools/wechat.py")
	err = cmd.Run()
	if err != nil {
		return
	}
	return
}
