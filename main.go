// Code generated by hertz generator.

package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzlogrus "github.com/hertz-contrib/logger/logrus"
	"github.com/mcoder2014/simple_short_url/biz/config"
	"github.com/mcoder2014/simple_short_url/biz/handler/simple_short_url"
	"github.com/mcoder2014/simple_short_url/util/mylog"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	port := os.Getenv("HOST_PORT")
	h := server.Default(server.WithHostPorts(fmt.Sprintf("%s", port)))
	InitLogger()
	if err := config.Init("./conf/config.yaml"); err != nil {
		panic(fmt.Errorf("init config failed,err=%v", err))
	}
	if err := simple_short_url.Init(); err != nil {
		panic(fmt.Errorf("init short service failed,err=%v", err))
	}

	register(h)
	h.Spin()
}

func InitLogger() {
	// 可定制的输出目录。
	var logFilePath string
	dir := "./log"
	logFilePath = dir + "/run/"
	if err := os.MkdirAll(logFilePath, 0o777); err != nil {
		log.Println(err.Error())
		return
	}

	// 将文件名设置为日期
	logFileName := time.Now().Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			log.Println(err.Error())
			return
		}
	}

	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logger := hertzlogrus.NewLogger(hertzlogrus.WithLogger(log))
	// 提供压缩和删除
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,   // 一个文件最大可达 100M。
		MaxBackups: 5,     // 最多同时保存 5 个文件。
		MaxAge:     10,    // 一个文件最多可以保存 10 天。
		Compress:   false, // 压缩配置。
	}

	logger.SetOutput(&mylog.MyLogWriter{
		Logger: lumberjackLogger,
		ToStd:  true,
	})
	logger.SetLevel(hlog.LevelDebug)

	hlog.SetLogger(logger)
}
