package main

import (
	"log"
	"os"
	"strings"

	"github.com/subosito/gotenv"
)

const (
	APP_NAME         = `E2E_WEB_API`
	APP_VERSION      = `v.0.1.10.build.1209`
	CONNECT_TIMEOUT  = 5
	E2E_GREEN_LEVEL  = `operational`
	E2E_YELLOW_LEVEL = `operationalWithFailures`
	E2E_RED_LEVEL1   = `inoperational`
	E2E_RED_LEVEL2   = `unknown`
	E2E_DBTIMEOUT    = `E2E database connect timeout.`
	E2E_NOTUPDATE    = `The status is too old.`
)

var ENV struct {
	CMP_GRPC_URL            string
	CMP_NOTIFY_URL          string
	E2E_SERVER_PORT         string
	E2E_QUERY_TIMER_SEC     int
	E2E_QUERY_TIMER_MIN     int
	E2E_NOTUPDATE_TIMER_MIN int
	E2E_FAKEDATA_TIME_SEC   int
	SQL_HOSTNAME            string
	SQL_USER                string
	SQL_PWD                 string
	SQL_PORT                int
	SENARIOID_LIST          []string
	SENARIO_RULES           string
	EXECUTION_RULES         string
	SUMMARY_LIST            []string
	FAKE_NEED_LIST          []string
}

func ENV_INIT() {
	/*-----------錯誤回收與安全退出---------------*/
	var err error
	defer func() {
		if err := recover(); err != nil { //這裡的err是新的，和外面的沒關係
			log.Println("[Error] ", err)
		}
	}()

	/*---------函式執行與確認是否存在.env---------*/
	log.Println(`ENV_INIT()`)
	err = gotenv.Load()
	if err != nil {
		log.Fatal("[Error] loading .env file fail, ", err)
	}
	/*---------------取得環境參數----------------*/

	//從.env取得環境變數E2E_SERVER_PORT
	tmp := os.Getenv(`E2E_SERVER_PORT`)
	if tmp != "" {
		ENV.E2E_SERVER_PORT = tmp
		//log.Println("環境變數 E2E_SERVER_PORT:", ENV.E2E_SERVER_PORT)
	} else {
		ENV.E2E_SERVER_PORT = `8443`
		log.Println("[WARN] 讀不到環境變數 E2E_SERVER_PORT，使用程式預設:",
			ENV.E2E_SERVER_PORT)
	}
	//從.env取得環境變數SENARIOID_LIST
	tmp = os.Getenv(`SENARIOID_LIST`)
	if tmp != "" {
		ENV.SENARIOID_LIST = make([]string, 0)
		slist := strings.Split(tmp, ",")
		for _, value := range slist {
			ENV.SENARIOID_LIST = append(ENV.SENARIOID_LIST, value)
		}
		//log.Println("環境變數 SENARIOID_LIST:", ENV.SENARIOID_LIST)
	} else {
		log.Println("[ERROR] 讀不到環境變數 SENARIOID_LIST")
		panic("SENARIOID_LIST can not query from ENV!")
	}

	//從.env取得環境變數SENARIO_RULES
	tmp = os.Getenv(`SENARIO_RULES`)
	if tmp != "" {
		ENV.SENARIO_RULES = tmp
		//log.Println("環境變數 SENARIO_RULES:", ENV.SENARIO_RULES)
	} else {
		log.Println("[ERROR] 讀不到環境變數 SENARIO_RULES")
		panic("SENARIO_RULES can not query from ENV!")
	}
}
