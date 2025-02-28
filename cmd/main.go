package main

import (
	"apiserver-gin/internal/middleware"
	"apiserver-gin/internal/middleware/trace"
	"apiserver-gin/internal/repo/mysql"
	"apiserver-gin/pkg/config"
	"apiserver-gin/pkg/log"
	"apiserver-gin/pkg/version"
	"apiserver-gin/server"
)

func main() {
	// 解析服务器启动参数
	appOpt := &server.AppOptions{}
	server.ResolveAppOptions(appOpt)
	if appOpt.PrintVersion {
		version.PrintVersion()
	}
	// 加载配置文件
	c := config.Load(appOpt.ConfigFilePath)
	log.InitLogger(&c.LogConfig,
		log.WithOption("appName", c.AppName),
		log.WithOption("requestId", trace.RequestId())) // 日志
	defer log.Sync()

	ds := mysql.NewDefaultMysql(c.DBConfig) // 创建数据库链接，使用默认的实现方式
	// 创建HTTPServer
	srv := server.NewHttpServer(config.GlobalConfig)
	srv.RegisterOnShutdown(func() {
		if ds != nil {
			ds.Close()
		}
	})
	router := initRouter(ds)
	srv.Run(middleware.NewMiddleware(), router)
}
