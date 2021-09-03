package logic

import (
	"WebApp/dao/mysql"
	"WebApp/dao/redis"
	"WebApp/modules"
	"WebApp/pkg/snowflake"
	"fmt"

	"go.uber.org/zap"
)

func CreatePost(p *modules.Post) (err error) {
	//1.生成帖子ID
	p.ID = snowflake.GenID()
	//2.保存到数据库，并把创建帖子的时间存到Redis
	err = mysql.CreatePost(p)
	if err != nil {
		//fmt.Println("MysqlError")
		return err
	}
	err = redis.CreatePost(p.ID)
	if err != nil {
		fmt.Println("RedisError")
		return err
	}
	return err
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

func GetPostList(page, size int64) (data []*modules.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*modules.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID)", zap.Error(err))
			continue
		}
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)", zap.Error(err))
			continue
		}
		PostDetail := &modules.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, PostDetail)
	}
	return
}
