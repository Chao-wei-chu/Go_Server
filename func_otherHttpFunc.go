package main

import (
	//"fmt"

	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpReplyMsg struct {
	Func      string      `json:"func"`
	ReqId     string      `json:"req_id"`
	Time      time.Time   `json:"time"`
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Result    interface{} `json:"result"`
}

func GetFuncName(c *gin.Context) string {
	var funcNameStr string
	funcName, exists := c.Get("x-func")
	if exists {
		if str, ok := funcName.(string); ok {
			funcNameStr = str
		} else {
			/* not string */
		}
	}
	return funcNameStr
}

func NewBasicHttpReplyMsg(c *gin.Context) *HttpReplyMsg {
	now := time.Now()
	replyMsg := &HttpReplyMsg{
		GetFuncName(c), strconv.FormatInt(now.Unix(), 10), now, 0, "OK", nil,
	}
	return replyMsg
}

func GetClientIP(c *gin.Context) {
	reply := NewBasicHttpReplyMsg(c)
	reply.Result = fmt.Sprintf("Your IP is %s", c.ClientIP())
	c.JSON(http.StatusOK, reply)
}
