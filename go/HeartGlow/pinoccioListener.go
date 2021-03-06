package main

import (
	"fmt"
	//  	"strings"
	"encoding/json"
	"github.com/andelf/go-curl"
	"time"
	//    	"os/exec"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native" // Native engine
	// _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

type Message struct {
	Data *jData
	DB   mysql.Conn
}

type jData struct {
	Type  string
	Troop string
	Scout string
	Value *jValue
}

type jValue struct {
	Status string
	Type   string
	Name   string
	Custom []string
}

func startListen() {
	sendingTroop := "13"
	readToken := "d12c26b23f51808b9d872bbba7be40da"
	fmt.Println("Curl is listening to Pinoccio...")
	var m Message
	var t string
	easy := curl.EasyInit()
	defer easy.Cleanup()
	m.DB = openDB()
	m.Data = new(jData)
	m.Data.Value = new(jValue)
	easy.Setopt(curl.OPT_URL, "https://api.pinocc.io/v1/sync?token="+readToken)
	// make a callback function
	fooTest := func(buf []byte, userdata interface{}) bool {
		err := json.Unmarshal([]byte(buf), &m)
		if err == nil && m.Data.Value.Type == "custom" && m.Data.Value.Name != "ping" {

			t = m.Data.Value.Name
			if m.Data.Troop == sendingTroop {
				fmt.Println("Type: " + t)
				fmt.Println("Type: " + t)
				fmt.Print("\n")
				jd := new(jsonData)
				jd.cmdType, jd.cmdData = m.Data.Value.Name, m.Data.Value.Custom
				jd.filter()
			}
		}

		return true
	}

	easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

	if err := easy.Perform(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
