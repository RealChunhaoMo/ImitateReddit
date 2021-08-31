package modules

//定义请求的参数结构体
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
	PostID   int64 `json:"post_id,string" binding:"required"`              // 帖子的ID
	VoteType int64 `json:"vote_type,string" binding:"required,oneof=1 -1"` // 投票的类型，是顶还是踩
}
