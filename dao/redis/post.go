package redis

import (
	"WebApp/modules"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//GetPostListsByOrder 按某种规则获取部分帖子的id
func GetPostListsByOrder(p *modules.ParamPostList) ([]string, error) {
	var key string
	//1.根据query string里参数，确定查redis时的key
	if p.Order == modules.OrderTime {
		key = GetFullkey(KeyPostTimeZset)
	} else {
		key = GetFullkey(KeyPostScoreZset)
	}
	//2.确定查询的范围,这里key的大到小
	return GetPostIDSByKey(key, p.Page, p.Size)
}

func GetPostIDSByKey(key string, Page, Size int64) ([]string, error) {
	start := (Page - 1) * Size
	stop := start + Size - 1
	return client.ZRevRange(key, start, stop).Result()
}

//GetCommunityPostListsByOrder 按社区查找帖子ids
func GetCommunityPostIDSByOrder(p *modules.ParamPostList) ([]string, error) {
	var orderkey string
	//1.根据query string里参数，确定查redis时的key
	if p.Order == modules.OrderTime {
		orderkey = GetFullkey(KeyPostTimeZset)
	} else {
		orderkey = GetFullkey(KeyPostScoreZset)
	}
	//使用 zinterstore 把分区的帖子与帖子分数的zset生成一个新的zset

	ckey := GetFullkey(keyCommunitSetPF + strconv.Itoa(int(p.CommunityID)))
	//利用缓存key来减少zinterstore执行的次数
	key := orderkey + ckey
	if client.Exists(key).Val() < 1 {
		//如果不存在,需要计算zinterstore
		pipeline := client.TxPipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderkey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}

	//根据key去查询IDS
	return GetPostIDSByKey(key, p.Page, p.Size)
}

//GetPostVoteData 获取每个帖子的得票情况
func GetPostVoteData(ids []string) (data []int64, err error) {
	pipeline := client.TxPipeline()
	for _, id := range ids {
		key := GetFullkey(KeyPostVotedZsetPF + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmdrs, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(ids))
	for _, cmder := range cmdrs {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
