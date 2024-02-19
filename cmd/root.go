package cmd

import (
	"GoScheduler/internal/modules/logger"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var AppVersion = "1.0.0"

var rootCmd = &cobra.Command{
	Use:     "GoScheduler",
	Short:   "GoScheduler定时任务管理系统",
	Long:    `GoScheduler是用Go实现的秒级分布式定时任务执行管理系统`,
	Version: AppVersion,
}

var allCmd = &cobra.Command{
	Use: "all",
	Run: func(cmd *cobra.Command, args []string) {
		allServers := []*cobra.Command{masterCmd, nodeCmd}
		for _, server := range allServers {
			server.Run(server, args)
		}
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}

func Execute() {
	logger.InitLogger()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
