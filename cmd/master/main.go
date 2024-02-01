package main

import (
	"GoScheduler/internal/models"
	"GoScheduler/internal/modules/app"
	"GoScheduler/internal/modules/global"
	setting2 "GoScheduler/internal/modules/setting"
	"GoScheduler/internal/routers"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// web服务器默认端口
var AppVersion = "1.0"

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = "GoScheduler"
	cliApp.Usage = "A Timed Task Management System"
	cliApp.Version = AppVersion
	cliApp.Commands = getCommands()
	cliApp.Flags = append(cliApp.Flags, []cli.Flag{}...)
	err := cliApp.Run(os.Args)
	if err != nil {
		zap.S().Fatal(err)
	}
}

func getCommands() []cli.Command {
	command := cli.Command{
		Name:   "web",
		Usage:  "run web server",
		Action: runWeb,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "host",
				Value: "0.0.0.0",
				Usage: "bind host",
			},
			cli.IntFlag{
				Name:  "port,p",
				Value: global.DefaultPort,
				Usage: "bind port",
			},
		},
	}

	return []cli.Command{command}
}

func runWeb(ctx *cli.Context) {
	// 初始化应用
	app.InitEnv()
	// 初始化模块 DB、定时任务等
	initModule(ctx)
	// 捕捉信号,配置热更新等
	go catchSignal()
}

func initModule(ctx *cli.Context) {
	// 初始化路由
	routers.InitRouter(ctx)

	if !global.Installed {
		return
	}

	setting2.InitConfig(global.AppConfig)

	models.InitDB()

	// 初始化定时任务
	service.ServiceTask.Initialize()
}

// 捕捉信号
func catchSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-c
		zap.S().Info("收到信号 -- ", s)
		switch s {
		case syscall.SIGHUP:
			zap.S().Info("收到终端断开信号, 忽略")
		case syscall.SIGINT, syscall.SIGTERM:
			shutdown()
		}
	}
}

// 应用退出
func shutdown() {
	defer func() {
		zap.S().Info("已退出")
		os.Exit(0)
	}()

	if !global.Installed {
		return
	}
	zap.S().Info("应用准备退出")
	// 停止所有任务调度
	zap.S().Info("停止定时任务调度")
	service.ServiceTask.WaitAndExit()
}
