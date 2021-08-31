package logic

import (
	"WebApp/dao/mysql"
	"WebApp/modules"
	"WebApp/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *modules.Post) (err error) {
	//1.生成帖子ID
	p.ID = snowflake.GenID()
	//2.保存到数据库
	return mysql.CreatePost(p)
	//3.返回
}

func GetPostDetail(id int64) (data *modules.ApiPostDetail, err error) {
	post, err := mysql.GetPostDetailByID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailByID(id) failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID)", zap.Error(err))
		return
	}
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)", zap.Error(err))
		return
	}
	data = &modules.ApiPostDetail{
		AuthorName:      "",
		Post:            post,
		CommunityDetail: community,
	}
	data.AuthorName = user.UserName
	return
}
