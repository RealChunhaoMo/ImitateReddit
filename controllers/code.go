package controllers

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodePasswordError
	CodeServerBusy
	CodeNeedLogin
	CodeSignUpError

	CodeEmptyAuth
	CodeErrorAuthFormat
	CodeNoworkToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "successful",
	CodeInvalidParam:    "请求参数错误",
	CodeUserExist:       "用户已存在",
	CodeUserNotExist:    "该用户不存在",
	CodeNeedLogin:       "你需要登录",
	CodePasswordError:   "密码错误",
	CodeServerBusy:      "服务器繁忙",
	CodeSignUpError:     "注册失败",
	CodeEmptyAuth:       "请求头auth为空",
	CodeErrorAuthFormat: "请求头auth格式有误",
	CodeNoworkToken:     "无效的token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
