package catapi

import (
	"fmt"
	"testing"
)

func TestExec(t *testing.T) {
	init := new(GetRequest).Init()
	params := init.AddParam("cmd", "getSharePrize").AddParam("token", "zoneToken").BuildParams()
	fmt.Println("params:", params)
}
