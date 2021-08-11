package settings

import (
	"fmt"
	"net/http"

	"github.com/fsnotify/fsnotify"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//定义一个全局变量，用于保存所有的配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

func Init() (err error) {
	viper.SetConfigFile("config.yaml") // 指定配置文件
	//viper.SetConfigType("yaml")        // 指定配置文件类型
	viper.AddConfigPath(".")   // 指定查找配置文件的路径
	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	//把读取到的配置信息，反序列化到全局变量Conf中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper Unmarshal failed!! error: %v\n", err)
	}
	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("配置文件发生改变了....")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("After config info changed!! viper Unmarshal failed!! error: %v\n", err)
		}
	})
	r := gin.Default()
	// 访问/version的返回值会随配置文件的变化而变化
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("version"))
	})

	if err := r.Run(
		fmt.Sprintf(":%d", Conf.Port)); err != nil {
		panic(err)
	}
	return
}
