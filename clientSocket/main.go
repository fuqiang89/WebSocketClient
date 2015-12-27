package main

import (
	"flag"
	"net/url"
	"strconv"
	"time"
	//    "encoding/json"
	"fmt"
	"reflect"

	//    "github.com/donnie4w/go-logger/logger"
	"github.com/bitly/go-simplejson"
	"github.com/fatih/structs"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"

	"rdpConn"
)

type SendMes struct {
	Version  string
	Username string
}

type ReadMes struct {
	MesType string      `1: remote connect`
	MesData interface{} `Important data`
	Status  string      `0: ok 1: error 2: warn`
	Info    string      `someting  messages`
}

type RdpData struct {
	Ip       string
	Port     string
	UserName string
	Password string
	ConType  string `1: mestc 2:putty 3:xshell`
	Path     string `rdp clinet file path`
}

func main() {
	//    logger.SetConsole(true)
	//    logger.SetRollingDaily(".", "test.log")
	//    logger.SetLevel(logger.ALL)

	flag.Parse()
	var addr = flag.String("addr", "localhost:8080", "http service address")
	var dialer = websocket.Dialer{}

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/gows"}

	for {
		errCh := make(chan error)
		ws, _, err := dialer.Dial(u.String(), nil)
		defer ws.Close()
		if err == nil {
			fmt.Println("runing")
			go readMessages(ws, errCh)
			go writeMessages(ws)
			fmt.Println(<-errCh)
			close(errCh)
		}
		if ws != nil {
			ws.Close()
		}
		fmt.Println("reconning....")
		time.Sleep(5 * time.Second)
	}
}

func readMessages(ws *websocket.Conn, errCh chan error) {
	for {
		//        rmes := new(SendMes)
		_, message, err := ws.ReadMessage()
		//        err := ws.ReadJSON(rmes)
		if err != nil {
			fmt.Println("read: >>>>>>> ", err)
			errCh <- err
			return
		}
		jsMess, err := simplejson.NewJson(message)
		if err != nil {
			fmt.Println(err)
		}
		go wsMessage(jsMess.Interface())
	}
}

func writeMessages(ws *websocket.Conn) {
	wmes := new(SendMes)
	wmes.Version = "0.1"
	wmes.Username = "test"
	sm := structs.Map(wmes)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		//        err := ws.WriteMessage(websocket.TextMessage, []byte("ddd"))
		err := ws.WriteJSON(sm)
		if err != nil {
			fmt.Println("write: >>>>>>> ", err)
			return
		}
	}
}

func wsMessage(message interface{}) {

	var readMes ReadMes
	err := mapstructure.Decode(message, &readMes)
	if err != nil {
		fmt.Println(err)
		return
	}
	mesType, err := strconv.Atoi(readMes.MesType)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reflect.TypeOf(mesType))
	switch mesType {
	case 1:
		//		rdp  remote connect
		connClient(readMes)
	default:
		fmt.Println("unknow message type, check you update info!")
	}
}

func connClient(message ReadMes) {
	var mesdata RdpData
	err := mapstructure.Decode(message.MesData, &mesdata)
	if err != nil {
		fmt.Println(err)
		return
	}
	conType, errConv := strconv.Atoi(mesdata.ConType)
	if errConv != nil {
		fmt.Println(errConv)
		return
	}
	switch conType {
	case 1:
		//		mstsc windows
		errMstsc := rdpConn.MstsCmd(mesdata.Ip, mesdata.UserName, mesdata.Password, mesdata.Port)
		if errMstsc != nil {
			fmt.Println("error mstsc", errMstsc)
		}

	case 2:
		//	putty linux
		errPutty := rdpConn.PuttyCmd(mesdata.Ip, mesdata.UserName, mesdata.Password, mesdata.Port, mesdata.Path)
		if errPutty != nil {
			fmt.Println("error mstsc", errPutty)
		}

	default:
		fmt.Println("unknow rdp connect type, check you update info!")
	}
}
