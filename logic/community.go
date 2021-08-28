package logic

import (
	"WebApp/dao/mysql"
	"WebApp/modules"
)

func GetCommunityList() ([]*modules.Community, error) {
	//查找数据库，返回所有的community
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*modules.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
