package rdpConn

import (
	"fmt"
	"os/exec"
)

func StrPuttyConn(ip string, username string, password string, port string, path string) string {
	return fmt.Sprintf("%s -ssh %s -P %s -l %s -pw %s   ", path, ip, port, username, password)
}

func PuttyCmd(ip string, username string, password string, port string, path string) error {
	connStr := StrPuttyConn(ip, username, password, port, path)
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
