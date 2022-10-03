package response

import (
	"github.com/JeffJiang620/anyu_gin/renders"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	statusCode int
	Code       string      `json:"code"`
	Data       interface{} `json:"data"`
	Msg        string      `json:"msg"`
}

func NewResponse(statusCode int, code string, data interface{}, msg string) *Response {
	return &Response{statusCode: statusCode, Code: code, Data: data, Msg: msg}
}

func (resp *Response) clone() *Response {
	c := *resp
	return &c
}

func (resp *Response) StatusCode() int {
	return resp.statusCode
}

func (resp *Response) Render(ctx *gin.Context) {
	ctx.Render(resp.statusCode, renders.JSON{Data: resp})
}

func (resp *Response) WithMsg(msg string) *Response {
	c := resp.clone()
	c.Msg = msg
	return c
}

func (resp *Response) WithData(data interface{}) *Response {
	c := resp.clone()
	c.Data = data
	return c
}

func (resp *Response) WithStatusCode(statusCode int) *Response {
	c := resp.clone()
	c.statusCode = statusCode
	return c
}

const (
	SuccessCode = "Success"
	SuccessMsg  = "成功"
)

var SuccessResponse = NewResponse(http.StatusOK, SuccessCode, nil, SuccessMsg)

func SuccessWithData(data interface{}) *Response {
	return SuccessResponse.WithData(data)
}

func SuccessWithDataAndMsg(data interface{}, msg string) *Response {
	return SuccessResponse.WithData(data).WithMsg(msg)
}
