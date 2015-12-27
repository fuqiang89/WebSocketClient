package rdpConn

import (
	"fmt"
	"os/exec"

	"WinIdent"
)

func StrMstsConn(ip string, port string) string {
	//	return fmt.Sprintf("mstsc /v:'%s'", ip)
	return fmt.Sprintf("mstsc /v:%s:%s", ip, port)
}

func MstsCmd(ip string, username string, password string, port string) error {
	connStr := StrMstsConn(ip, port)
	err := winIdent.CreateRDPSession(ip, username, password)
	defer winIdent.DelRDPSession(ip)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//	Cmd := exec.Command("PowerShell", "-Command", connStr)
	fmt.Println(connStr)
	Cmd := exec.Command("cmd", "/c", connStr)
	_, errCmd := Cmd.Output()
	if errCmd != nil {
		fmt.Println(errCmd)
		return errCmd
	}
	return nil
}
