package main

import (
	"bookcrawer/config"
	"bookcrawer/models"
	"bookcrawer/router"
	"errors"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

var (
	cfg = pflag.StringP("config", "c", "", "bookcrawer config file path.")
)

func main() {
	flag.Parse()

	// 初始config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// 初始化数据库
	models.DB.Init()
	defer models.DB.Close()

	gin.SetMode(viper.GetString("runmode"))

	g := gin.New()

	//路由
	router.Load(g)

	//是否存活
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()
	// 开始监听
	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// 用一个Goroutine检测是否可用
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}
