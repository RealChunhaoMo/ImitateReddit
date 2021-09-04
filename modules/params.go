package modules

//定义请求的参数结构体
const (
	OrderTime   = "time"
	OrderScore  = "score"
	DefaultPage = 1
	DefaultSize = 10
)

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type ParamSignIn struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VoteData struct {
	//UserID 用户id可以直接从发请求的用户获取
	PostID   string `json:"post_id" binding:"required"`              // 帖子的ID
	VoteType int64  `json:"vote_type,string" binding:"oneof=1 0 -1"` // 投票的类型，是顶还是踩
}

type ParamPostList struct {
	Page        int64  `json:"page" form:"page"`
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"`
}
