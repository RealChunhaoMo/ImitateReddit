package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const AWeekInSeconds = 24 * 3600 * 7

var (
	ErrVoteTimeHasPassed = errors.New("投票时间已过")
	ScorePerVote         = 432 //每一票的分值
)

/*
基于用户投票的分数统计：
看了阮一峰的博客 https://www.ruanyifeng.com/blog/2012/02/ranking_algorithm_hacker_news.html

这里就用比较简单的算法吧，他博客上面的比较麻烦吧。
投一票获得432分，如果帖子获得200票就可以让该帖子续上一天，一天是24*360 = 86400s = 432*200

投票的分类
1.VoteType = 1(赞成票)
1.1之前没投过票  	1 - 0 = 1
1.2之前投过反对票	1 - (-1) = 2

2.VoteType = 0(取消投票)
2.1之前投赞成票 	0 - 1 = -1
2.2之前投反对票 	0 - (-1) = 1

3.VoteType = -1(反对票)
3.1之前没投票 	-1 - 0 = -1
3.2之前投赞成票 	-1 - 1 = -2

投票的限制
为了减轻数据库的压力，只允许对最近几天发布的帖子投票，
如果超过了时间限制，那么投票的结果就国定了，存到MySql里就好了
*/

func CreatePost(PostID, CommunityID int64) error {
	pipeline := client.TxPipeline()
	//添加帖子创建的时间
	pipeline.ZAdd(GetFullkey(KeyPostTimeZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: PostID,
	})

	//初始化帖子的分数
	pipeline.ZAdd(GetFullkey(KeyPostScoreZset), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: PostID,
	})

	//把帖子id放入社区的set中
	ckey := GetFullkey(keyCommunitSetPF + strconv.Itoa(int(CommunityID)))
	pipeline.SAdd(ckey, PostID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, PostID string, CurrentVote float64) error {
	//1.检查投票的限制
	PostTime := client.ZScore(GetFullkey(KeyPostTimeZset), PostID).Val()
	if float64(time.Now().Unix())-PostTime > AWeekInSeconds {
		return ErrVoteTimeHasPassed
	}
	//2.更新帖子的分数
	//2.1去查当前用户的投票情况
	LastVote := client.ZScore(GetFullkey(KeyPostVotedZsetPF+PostID), userID).Val()
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(GetFullkey(KeyPostScoreZset), (CurrentVote-LastVote)*float64(ScorePerVote), PostID).Result()
	//3.处理用户的投票
	if CurrentVote == 0 {
		pipeline.ZRem(GetFullkey(KeyPostVotedZsetPF+PostID), userID)
	} else {
		pipeline.ZAdd(GetFullkey(KeyPostVotedZsetPF+PostID), redis.Z{
			Score:  CurrentVote,
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
