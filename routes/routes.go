package routes

import (
	"WebApp/controllers"
	"WebApp/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//用户注册
	r.POST("/signup", controllers.SignUpHandler)

	//用户登录
	r.POST("/signin", controllers.SignInHandler)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "SignUp successfully!!")
	})
	return r
}
