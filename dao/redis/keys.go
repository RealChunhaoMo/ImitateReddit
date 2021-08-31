package redis

//redis keys
const (
	KerPreFix          = "WebApp:"
	KeyPostTimeZset    = "post:time"   // 帖子和其发布时间
	KeyPostScoreZset   = "post:score"  // 帖子和其分数
	KeyPostVotedZsetPF = "post:voted:" // 记录用户投票类型，参数是post_id
)
