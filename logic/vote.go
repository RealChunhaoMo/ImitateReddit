package logic

import (
	"WebApp/dao/redis"
	"WebApp/modules"
	"strconv"

	"go.uber.org/zap"
)

func PostVote(uid int64, data *modules.VoteData) error {
	zap.L().Debug("PostVote",
		zap.Int64("userID", uid),
		zap.String("PostID", data.PostID),
		zap.Int64("VoteType", data.VoteType))
	return redis.VoteForPost(strconv.Itoa(int(uid)), data.PostID, float64(data.VoteType))
}
