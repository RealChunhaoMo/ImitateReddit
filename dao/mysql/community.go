package mysql

import (
	"WebApp/modules"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (communitylist []*modules.Community, err error) {
	sqlStr := "select community_id,community_name from community"
	if err := db.Select(&communitylist, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("No community in database")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (community *modules.CommunityDetail, err error) {
	sqlStr := `select community_id,community_name,introduction,
	create_time from community where community_id = ?`
	community = new(modules.CommunityDetail)
	if err := db.Get(community, sqlStr, id); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
		}
	}
	return community, err
}
