package controllers

import (
	"WebApp/logic"
	"WebApp/modules"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

//CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	//1.获取参数以及参数校验
	p := new(modules.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON failed", zap.Error(err))
		zap.L().Error("Create Post with Invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}

	//获取当前用户id
	userID, err := GetCurrentUserID(c)
	fmt.Println("usrID = ", userID)
	if err != nil {
		zap.L().Error("GetCurrentUserID(c) failed", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	//2.创建帖子
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, nil)
}

//GetPostDetailHandler 根据帖子id获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	//1.获取参数，从url中获取帖子的id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.通过id获取帖子数据
	data, err := logic.GetPostDetail(id)
	if err != nil {
		zap.L().Error("logic.GetPostDetail(id) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, data)
}

//GetPostListHandler 根据请求参数获取帖子列表
func GetPostListHandler(c *gin.Context) {
	page, size := GetPageInfo(c)
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList()", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//2.返回响应
	ResponseSuccess(c, data)
}

//GetPostListHandler2 获取帖子列表，帖子列表可根据创建时间和分数来排列
func GetPostListHandler2(c *gin.Context) {
	//从query string获取参数
	p := &modules.ParamPostList{
		Page:  modules.DefaultPage,
		Size:  modules.DefaultSize,
		Order: modules.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("ParamPostList with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.GetPostListUnion(p)
	if err != nil {
		zap.L().Error("logic.GetPostListNew(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

////GetCommunityPostListHandlers 根据社区分类，获取同属一个社区的帖子
//func GetCommunityPostListHandlers(c *gin.Context) {
//	p := &modules.ParamCommunityPostList{
//		ParamPostList: &modules.ParamPostList{
//			Page:  modules.DefaultPage,
//			Size:  modules.DefaultSize,
//			Order: modules.OrderTime,
//		},
//	}
//	//从query string获取参数
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("ParamCommunityPostList with invalid param", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	data, err := logic.GetCommunitPostList(p)
//	if err != nil {
//		zap.L().Error("logic.GetPostList()", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//	//2.返回响应
//	ResponseSuccess(c, data)
//}
