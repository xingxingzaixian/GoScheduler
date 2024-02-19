package cmd

import (
	"GoScheduler/internal/modules/global"
	"GoScheduler/internal/routers"
	"GoScheduler/internal/task"
	"context"
	"github.com/spf13/cobra"
)

var masterCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"serve"},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		startSever(cmd.Context(), host, port)
	},
}

func init() {
	masterCmd.Flags().StringP("host", "H", "0.0.0.0", "server address")
	masterCmd.Flags().IntP("port", "p", 5320, "server port")
	rootCmd.AddCommand(masterCmd)
}

func startSever(ctx context.Context, host string, port int) {
	// 1. 初始化路由
	routers.InitRouter(host, port)

	// 2. 如果没安装，就执行安装
	if !global.Installed {
		return
	}

	// 3. 初始化任务调度
	task.Initialize()
}
