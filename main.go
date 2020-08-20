package main

import (
	//"fmt"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	debug      *bool = flag.Bool("debug", false, "enable debugging")
	live_mode  *bool = flag.Bool("live", false, "enable live mode")       //切換DB
	fake_data  *bool = flag.Bool("fake", false, "enable insert fakedata") //注意: 以後此指令只用於決並是否插入和使用假資料(至於項目要不要用假資料用ENV.FAKE_NEED_LIST判斷)
	chan_t1    *time.Ticker
	chan_t2    *time.Ticker
	chan_flag1 chan int
)

func main() {
	//讀取外部CLI參數
	flag.Parse()

	//共用變數初始化(注意INIT的呼叫順序不能錯)
	ENV_INIT()
	VAR_INIT()

	//程式預設訊息
	var err error
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("[START]", APP_NAME, APP_VERSION)
	log.Println("Go version:", runtime.Version())
	log.Println("GIN version:", gin.Version)

	// 创建带有默认中间件的路由:
	router := gin.Default()
	//创建不带中间件的路由：
	//router := gin.New()

	router.GET("/", GetClientIP)

	router.GET("/favicon.ico", func(c *gin.Context) {
		c.JSON(http.StatusOK, "")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/time", func(c *gin.Context) {
		c.String(http.StatusOK, time.Now().Format(time.RFC3339))
	})
	/*---------方法一，無任何認證啟用URL-------------*/

	router.GET("/ARGUSView/serviceManagement/healthState/:senario_id", func(c *gin.Context) {
		senario_id := c.Param("senario_id")
		c.JSON(http.StatusOK, gin.H{
			"message": senario_id,
		})
	})

	router.GET("/ARGUSView.:senario_id/serviceManagement/executionState/rest", func(c *gin.Context) {
		//senario_id := c.Param("senario_id")
		//(DB有問題 先註解)ApiArgusGetServiceExecutionState(c, senario_id)
		c.String(http.StatusOK, "Hello")
	})

	router.POST("/ARGUSView.:senario_id/serviceManagement/serviceTest/rest", func(c *gin.Context) {
		var json TestJSON
		err := c.BindJSON(&json)
		if err == nil {
			for i := 0; i < len(json.TestList); i++ {
				fmt.Printf("%+v\n", json.TestList[i])
			}
			c.String(http.StatusOK, "Hello")
		} else {
			log.Println("[WARN] POST query error:", err)
			c.String(http.StatusBadRequest, "")
		}
	})

	/*---------方法二，認證JSON Token啟用URL--------*/

	/*---------第一種，開server方法(默認)-----------*/
	//router.Run() // listen and serve on 0.0.0.0:8080

	/*---------第二種，開server方法(HTTP)-----------*/

	server := http.Server{
		Addr:           fmt.Sprintf(":%s", ENV.E2E_SERVER_PORT),
		Handler:        router,
		ReadTimeout:    CONNECT_TIMEOUT * time.Second,
		WriteTimeout:   CONNECT_TIMEOUT * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("[INFO] Server is Ready....")
	err = server.ListenAndServe()

	/*---------第三種，開server方法(HTTPS)----------*/
	/*
		tlsConfig := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		server := http.Server{
			Addr:           fmt.Sprintf(":%s", ENV.E2E_SERVER_PORT),
			Handler:        router,
			ReadTimeout:    CONNECT_TIMEOUT * time.Second,
			WriteTimeout:   CONNECT_TIMEOUT * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      tlsConfig,
			TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
		//err = server.ListenAndServeTLS("CertB64.cer", "server.key")
		log.Println("[INFO] Server is Ready....")
		log.Println("[INFO] Start DB Process....")
		if *fake_data == true {
			go DB_Fakedata_Timer(chan_t2) //這裡是為了讓IoT大平台取得狀態好看，持續塞入假資料假裝有在更新
		}
		go Db_Query_Timer(chan_t1, chan_flag1) //這裡設計是只要DB Paser出錯，就直接退出程式
		err = server.ListenAndServeTLS("server.crt", "server.key")
	*/
	/*----------------Listen & Shutdown Server-----------------*/
	if err != nil && err != http.ErrServerClosed {
		log.Println("[ERROR]", err)
	} else {
		//do nothing.....
	}

	//Release connect & memory
	ctx, cancel := context.WithTimeout(context.Background(), CONNECT_TIMEOUT*time.Second)
	defer func() {
		//Notice! This can close database, redis, truncate message queues, etc.
		cancel()
	}()

	//use "graceful shutdown" for go http server
	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("[ERROR] Server Shutdown Failed: %+v", err)
	}

	log.Println("Close program after 3sec......")
	time.Sleep(3 * time.Second)
}
