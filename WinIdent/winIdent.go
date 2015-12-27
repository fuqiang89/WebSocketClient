package winIdent

import (
	"fmt"
	"os/exec"
)

func addStrCmdkey(ip, username, password string) string {
	return fmt.Sprintf("cmdkey /generic:TERMSRV/%s /user:%s /pass:%s", ip, username, password)
}

func delStrCmdkey(ip string) string {
	return fmt.Sprintf("cmdkey /delete TERMSRV/%s", ip)
}

func CreateRDPSession(ip, username, password string) error {
	strCmdkey := addStrCmdkey(ip, username, password)
	//	Cmd := exec.Command("PowerShell", "-Command", strCmdkey)
	Cmd := exec.Command("cmd", "/c", strCmdkey)
	_, err := Cmd.Output()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DelRDPSession(ip string) error {
	strCmdkey := delStrCmdkey(ip)
	//	Cmd := exec.Command("PowerShell", "-Command", strCmdkey)
	Cmd := exec.Command("cmd", "/c", strCmdkey)
	_, err := Cmd.Output()
	fmt.Println("del key")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
