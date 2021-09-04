package logic

import (
	"WebApp/dao/mysql"
	"WebApp/dao/redis"
	"WebApp/modules"
	"WebApp/pkg/snowflake"
	"fmt"

	"go.uber.org/zap"
)

//CreatePost 创建帖子
func CreatePost(p *modules.Post) (err error) {
	//1.生成帖子ID
	p.ID = snowflake.GenID()
	//2.保存到数据库，并把创建帖子的时间存到Redis
	err = mysql.CreatePost(p)
	if err != nil {
		//fmt.Println("MysqlError")
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	if err != nil {
		fmt.Println("RedisError")
		return err
	}
	return err
}

//GetPostDetail 根据帖子id获取帖子详情
func GetPostDetail(id int64) (data *modules.ApiPostDetail, err error) {
	post, err := mysql.GetPostDetailByID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailByID(id) failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
		return
	}
	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(post.AuthorID) failed", zap.Error(err))
		return
	}
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(post.CommunityID)", zap.Error(err))
		return
	}
	data = &modules.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

//GetPostList 根据page和size获取帖子
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

//GetPostList2 获取帖子列表，帖子列表可根据创建时间和分数来排列
func GetPostList2(p *modules.ParamPostList) (data []*modules.ApiPostDetail, err error) {
	//2.去redis里获取帖子id
	//3.根据帖子id查询帖子详情
	ids, err := redis.GetPostListsByOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListsByOrder(p) return 0 data")
		return
	}
	postlist, err := mysql.GetPostListByIDS(ids)
	if err != nil {
		return
	}

	//获取每个帖子的得票情况
	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData(ids) failed", zap.Error(err))
		return
	}
	for index, post := range postlist {
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
			VoteNum:         votedata[index],
			CommunityDetail: community,
		}
		data = append(data, PostDetail)
	}
	return
}

//GetCommunitPostList 按社区分类获取同一社区的帖子
func GetCommunitPostList(p *modules.ParamPostList) (data []*modules.ApiPostDetail, err error) {
	ids, err := redis.GetCommunityPostIDSByOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostListsByOrder(p) return 0 data")
		return
	}
	postlist, err := mysql.GetPostListByIDS(ids)
	if err != nil {
		return
	}

	votedata, err := redis.GetPostVoteData(ids)
	if err != nil {
		zap.L().Error("redis.GetPostVoteData(ids) failed", zap.Error(err))
		return
	}
	for index, post := range postlist {
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
			VoteNum:         votedata[index],
			CommunityDetail: community,
		}
		data = append(data, PostDetail)
	}
	return
}

//GetPostListNew 将上面的GetPostList2和GetCommunitPostList合二为一
func GetPostListUnion(p *modules.ParamPostList) (data []*modules.ApiPostDetail, err error) {
	//根据请求参数的不同，使用不同的处理逻辑
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunitPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
