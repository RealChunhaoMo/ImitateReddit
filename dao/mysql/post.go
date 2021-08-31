package mysql

import (
	"WebApp/modules"
)

func CreatePost(p *modules.Post) (err error) {
	sqlStr := `insert into post(
    post_id,title,content,author_id,community_id) values(?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostDetailByID(id int64) (post *modules.Post, err error) {
	sqlStr := `select post_id,author_id,community_id,title,content,
	create_time from post where post_id = ?
	`
	post = new(modules.Post)
	err = db.Get(post, sqlStr, id)
	return
}

func GetPostList(page, size int64) (posts []*modules.Post, err error) {
	sqlStr := `select post_id,author_id,community_id,title,content,
	create_time from post limit ?,?
	`
	posts = make([]*modules.Post, 0, 2)
	err = db.Select(&posts, sqlStr, page-1, size)
	return
}
