package nmcli

import (
	"errors"
	"os/exec"
	"strings"
)

// Checks whether nmcli is installed on the system
// Return an error if nmcli not installed, and a version string if it is installed
func CheckNmcliExist() (string, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		"nmcli --version",
	).CombinedOutput()
	if err != nil {
		return "", errors.New("nmcli not found on this system")
	}
	return string(res), nil
}

func InstallNmcli() (string, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		"apt install network-manager",
	).CombinedOutput()
	if err != nil {
		return "", errors.New("nmcli install completed")
	}
	return string(res), nil
}


// check nmcli is enable
func CheckNmcliEnable() (bool, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		"nmcli networking",
	).CombinedOutput()
	if err != nil {
		return false, errors.New("nmcli networking command error:"+err.Error())
	}
	status:=string(res)
	status,_=strings.CutSuffix(status,"\n")
	if status=="enabled"{
		return true,nil
	}
	return false, nil
}




// nmcli n on/off
func ChangeNmcliStatus(status string) (string, error) {
	res, err := exec.Command(
		"bash",
		"-c",
		"nmcli n "+status,
	).CombinedOutput()
	if err != nil {
		return "",errors.New("nmcli n command error:"+err.Error())
	}
	return string(res), nil
}


