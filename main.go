package main

import (
	"WebApp/controllers"
	"WebApp/dao/mysql"
	"WebApp/dao/redis"
	"WebApp/logger"
	"WebApp/pkg/snowflake"
	"WebApp/routes"
	"WebApp/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"go.uber.org/zap"
)

// GoWeb开发通用的模板
//1.加载配置
//2.初始化日志
//3.初始化MySql
//4.初始化Redis
//5.路由注册
//6.启动服务（优雅地关机）

func main() {
	//if len(os.Args) < 2 {
	//	fmt.Println("忘记输入命令行参数了!!!!!")
	//	return
	//}
	//1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("Settings Init failed!%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("logger Init failed!%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("Logger init is ok!!")
	//3.初始化MySql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("mysql Init failed!%v\n", err)
		return
	}
	defer mysql.Close()
	//4.初始化Redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("redis Init failed!%v\n", err)
		return
	}
	defer redis.Close()
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("Snowflake Init Failed! error = %v\n", err)
		return
	}

	//设置gin框架内置的参数校验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("controllers.InitTrans Init falied! error = %v\n", err)
		return
	}
	//5.路由注册
	r := routes.Setup(settings.Conf.Mode)
	r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	//6.启动服务（优雅地关机）

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
