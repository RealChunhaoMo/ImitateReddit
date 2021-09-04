package mysql

import (
	"WebApp/modules"
	"strings"

	"github.com/jmoiron/sqlx"
)

func CreatePost(p *modules.Post) (err error) {
	sqlStr := `insert into post(
    post_id,title,content,author_id,community_id) values(?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

//GetPostDetailByID 根据帖子id查询单个帖子
func GetPostDetailByID(id int64) (post *modules.Post, err error) {
	sqlStr := `select post_id,author_id,community_id,title,content,
	create_time from post where post_id = ?
	`
	post = new(modules.Post)
	err = db.Get(post, sqlStr, id)
	return
}

//GetPostList 根据query string参数查询帖子列表
func GetPostList(page, size int64) (posts []*modules.Post, err error) {
	sqlStr := `select post_id,author_id,community_id,title,content,
	create_time from post 
	ORDER BY create_time
	DESC 
	    limit ?,?
	`
	posts = make([]*modules.Post, 0, size)
	err = db.Select(&posts, sqlStr, page-1, size)
	return
}

//GetPostListByIDS 根据给定数据查询帖子列表
func GetPostListByIDS(ids []string) (postlist []*modules.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
	from post where post_id in(?)
	order by FIND_IN_SET(post_id,?)`
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postlist, query, args...)
	return
}
