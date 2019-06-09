package util

import (
	"encoding/json"
	"fmt"
)

// RespMsg : http响应数据的通用结构
type RespMsg struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewRespMsg : 生成response对象
func NewRespMsg(code int32, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// JSONBytes : 对象转json格式的二进制数组
func (resp *RespMsg) JSONByte() []byte {
	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("json parse err:", err.Error())
		return nil
	}
	return data
}

// JSONString : 对象转json格式的string
func (resp *RespMsg) JSONString() string {
	data := resp.JSONByte()
	return string(data)
}

// GenSimpleRespStream : 只包含code和message的响应体([]byte)
func GenSimpleRespStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))
}

// GenSimpleRespString : 只包含code和message的响应体(string)
func GenSimpleRespString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}
