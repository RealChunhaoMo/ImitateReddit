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
