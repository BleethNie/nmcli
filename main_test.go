package nmcli

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_CheckNmcliExist(t *testing.T) {
	info, err := CheckNmcliExist()
	if err != nil {
		fmt.Println("nmcli check error:" + err.Error())
	}
	fmt.Println("nmcli check info:" + info)
}

func Test_InstallNmcli(t *testing.T) {
	info, err := InstallNmcli()
	if err != nil {
		fmt.Println("nmcli install error:" + err.Error())
	}
	fmt.Println("nmcli install info:" + info)
}

func Test_CheckNmcliEnable(t *testing.T) {
	enable, err := CheckNmcliEnable()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("nmcli enable status:" + strconv.FormatBool(enable))
}
