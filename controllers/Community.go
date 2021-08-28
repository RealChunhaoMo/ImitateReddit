package controllers

import (
	"WebApp/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//社区相关的
func CommunityHandler(c *gin.Context) {
	//查询所有社区以列表形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不要把服务器端错误暴露给外面
		return
	}
	ResponseSuccess(c, data)
}

//社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	//获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("GetCommunityDetail Failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}
