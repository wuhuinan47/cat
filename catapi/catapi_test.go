package catapi

import (
	"fmt"
	"strings"
	"testing"
)

func TestExec(t *testing.T) {
	// init := new(GetRequest).Init()
	// params := init.AddParam("cmd", "getSharePrize").AddParam("token", "zoneToken").BuildParams()
	// fmt.Println("params:", params)

	s := "1,2,3,4,5,6"

	var ids = []float64{1, 2, 3, 4, 5, 6, 7}
	for _, v := range ids {
		fmt.Println("ssss:", strings.Contains(s, fmt.Sprintf("%v", v)))
	}

}
