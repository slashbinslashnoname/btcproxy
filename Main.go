package main

import (
	"encoding/json"
	"flag"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

func main() {
	// 解析命令行参数
	configFilePath := flag.String("c", "agent_conf.json", "Path of config file")
	logDir := flag.String("l", "", "Log directory")
	flag.Parse()

	if *logDir == "" || *logDir == "stderr" {
		flag.Lookup("logtostderr").Value.Set("true")
	} else {
		flag.Lookup("log_dir").Value.Set(*logDir)
	}

	// 增大文件描述符上限
	IncreaseFDLimit()

	// 读取配置文件
	config := NewConfig()
	err := config.LoadFromFile(*configFilePath)
	if err != nil {
		glog.Fatal("load config failed: ", err)
		return
	}
	config.Init()

	// 打印加载的配置文件（用于调试）
	if glog.V(3) {
		configBytes, _ := json.Marshal(config)
		glog.Info("config: ", string(configBytes))
	}

	// 启动 HTTP 调试服务
	if config.HTTPDebug.Enable {
		glog.Info("HTTP debug enabled: ", config.HTTPDebug.Listen)
		go func() {
			err := http.ListenAndServe(config.HTTPDebug.Listen, nil)
			if err != nil {
				glog.Error("launch http debug service failed: ", err.Error())
			}
		}()
	}

	// 会话管理器
	manager := NewSessionManager(config)

	// 退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		glog.Info("exiting...")
		manager.Stop()
	}()

	// 运行代理
	manager.Run()
}
