package controllers

import (
	"WebApp/logic"
	"WebApp/modules"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 1.参数校验
	p := new(modules.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		// 非validator.ValidationErrors类型错误直接返回
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//fmt.Println(&p)
	// 2.业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("SignUp Failed!!!", zap.Error(err))
		ResponseError(c, CodeSignUpError)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

func SignInHandler(c *gin.Context) {
	// 1.参数校验
	p := new(modules.ParamSignIn)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，返回响应
		zap.L().Error("SignIn with invalid param", zap.String("uername", p.Username), zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		// 非validator.ValidationErrors类型错误直接返回
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2.业务处理
	if err := logic.SignIn(p); err != nil {
		zap.L().Error("SignIn Failed!!!", zap.String("username", p.Username), zap.Error(err))
		ResponseError(c, CodePasswordError)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}
