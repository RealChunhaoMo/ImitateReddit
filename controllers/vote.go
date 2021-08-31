package controllers

import (
	"WebApp/logic"
	"WebApp/modules"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func PostVoteHandler(c *gin.Context) {
	//1.参数校验
	p := new(modules.VoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) // 做一个类型断言，有可能你的参数错误还没到validator这一层
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errdata := removeTopStruct(errs.Translate(trans)) // 翻译错误并去除错误中的结构体
		ResponseErrorWithMsg(c, CodeInvalidParam, errdata)
		return
	}

	//2.业务处理
	UserID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	logic.PostVote(UserID, p)
	//3.返回响应
	ResponseSuccess(c, nil)
}
