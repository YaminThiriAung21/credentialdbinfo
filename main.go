
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)
type ConfigArray struct{
	DbConfig []DbConfig
}
type DbConfig struct {
	Engine string
	Host string
	Port string
	Username string
	Password string
}

func main(){
	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened config.json")
	defer jsonFile.Close()

	jsonValue, _ := ioutil.ReadAll(jsonFile)
	var dbconfig ConfigArray
	json.Unmarshal(jsonValue, &dbconfig)

	credentialport := "9798"
	hostraw, err := exec.Command("sh", "-c", fmt.Sprintf("hostname -i")).Output()
  	hoststr := string(hostraw)
  	host := strings.TrimSuffix(hoststr, "\n")
	raw_connect(host,credentialport,dbconfig)
}

func raw_connect(host string, port string,config ConfigArray) {
		timeout := time.Second
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
		if err != nil {
			fmt.Println("Connecting error:", err)
		}
		if conn != nil {
      fmt.Println("Opened", net.JoinHostPort(host, port))
      insert_dbinfo(host,port,config)
			defer conn.Close() 
   }
}
func insert_dbinfo(host string,port string,config ConfigArray){
			
		url := "http://"+host+":"+port+"/api/credential/v6/dbinfos"
			fmt.Println("URL:>", url)
	for _, s := range config.DbConfig {

		var db = []byte(fmt.Sprintf(`[{"Engine":"%s","Host": "%s","Port":"%s","Username":"%s","Password":"%s"}]`, s.Engine, s.Host, s.Port, s.Username, s.Password))
		reqdb, err := http.NewRequest("POST", url, bytes.NewBuffer(db))
		reqdb.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		respdb, err := client.Do(reqdb)
		if err != nil {
			panic(err)
		}
		defer respdb.Body.Close()
		body, _ := ioutil.ReadAll(respdb.Body)
		fmt.Println("response Body:", string(body))
	}
}
