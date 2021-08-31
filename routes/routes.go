package routes

import (
	"WebApp/controllers"
	"WebApp/controllers/middlewares"
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

	//r.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	c.String(http.StatusOK, "pong")
	//})
	////用户注册
	//r.POST("/signup", controllers.SignUpHandler)
	////用户登录
	//r.POST("/signin", controllers.SignInHandler)
	//r.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "SignUp successfully!!")
	//})

	v1 := r.Group("/api/v1")
	//v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
	//	c.String(http.StatusOK, "pong")
	//})
	v1.Use(middlewares.JWTAuthMiddleware())
	v1.POST("/signup", controllers.SignUpHandler)
	v1.POST("/login", controllers.SignInHandler)
	v1.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "SignUp successfully!!")
	})

	{
		v1.GET("/Community", controllers.CommunityHandler)
		v1.GET("/Community/:id", controllers.CommunityDetailHandler)
		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
	}
	return r
}
