package mysql

import (
	"WebApp/modules"
	"WebApp/settings"
	"testing"
)

func init() {
	dbconf := settings.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "mch123",
		DB:           "sql_demo",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbconf)
	if err != nil {
		panic(err)
	}
}
func TestCreatePost(t *testing.T) {

	post := &modules.Post{
		ID:          10,
		AuthorID:    123,
		Title:       "test",
		CommunityID: 1,
		Content:     "just a post test",
	}
	err := CreatePost(post)
	if err != nil {
		t.Fatalf("CreatePost insert into mysql failed!!error = %#v", err)
	}
	t.Logf("CreatePost insert into mysql Success!!")
}
