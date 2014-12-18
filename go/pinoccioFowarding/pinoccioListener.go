package main

import (
   	"fmt"
//  	"strings"
    	"encoding/json"
    	"github.com/andelf/go-curl"
	"os"
//	"time"
  //  	"os/exec"
"github.com/ziutek/mymysql/mysql"
    _ "github.com/ziutek/mymysql/native" // Native engine
    // _ "github.com/ziutek/mymysql/thrsafe" // Thread safe engine
)

type Message struct{
	Data *jData
	DB   mysql.Conn
}

type jData struct {
	Type string
	Troop string
	Scout string
	Value *jValue
}

type jValue struct {
	Status string
	Type string 
	Name string
	Custom []string	
}

var redTroop, redId, greenTroop, greenId string
func startListen() {
    var redTroop, greenTroop, redId, greenId, readToken string
    checkArguments(os.Args)
    redTroop = os.Args[1]
    redId = os.Args[2]
    greenTroop = os.Args[3]
    greenId = os.Args[4]
    readToken = "98d30e31f54617dfe570a523463067db"
    enableProgram("6c5da8cf710ba4d745b8b00b523a9f6a")
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
    fooTest := func (buf []byte, userdata interface{}) bool {
	err := json.Unmarshal([]byte(buf), &m)
	if(err == nil && m.Data.Value.Type == "custom" && m.Data.Value.Name != "ping") {

		t = m.Data.Value.Name
		fmt.Println("Type: "+t)
		fmt.Print("\nTroop: "+m.Data.Troop)
		jd := new(jsonData)
		jd.cmdType, jd.cmdData = m.Data.Value.Name, m.Data.Value.Custom
		if(m.Data.Troop == greenTroop){
			jd.filter("red", redId)
		}else if(m.Data.Troop == redTroop){
			jd.filter("green", greenId)
		}
		
	}

	return true
    }

    easy.Setopt(curl.OPT_WRITEFUNCTION, fooTest)

    if err := easy.Perform(); err != nil {
        fmt.Printf("ERROR: %v\n", err)
    }
}

func enableProgram(write string){
        c := "https://api.pinocc.io/v1/" + greenTroop + "/"+greenId+"/command/en?token="+write
        runCmd(c)
        c = "https://api.pinocc.io/v1/" + redTroop + "/"+redId+"/command/en?token="+write
        runCmd(c)
}

func checkArguments(args []string){
	if(len(args) < 5){
		fmt.Println("Too few arguments, 4 needed\npinoccioFowarding [redTroopId] [redId] [greenTroopId] [greenId]")
		os.Exit(0)
	}else if(len(args) > 5){
		fmt.Println("Too many arguments, only 4 needed\n pinoccioFowarding [redTroopId] [redId] [greenTroopId] [greenId]")
		os.Exit(0)
	}
}


func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

