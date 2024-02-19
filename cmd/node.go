package cmd

import (
	"GoScheduler/internal/modules/rpc/server"
	"fmt"
	"github.com/spf13/cobra"
)

var nodeCmd = &cobra.Command{
	Use: "node",
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		port, _ := cmd.Flags().GetInt("port")
		serverAddr := fmt.Sprintf("%s:%d", host, port)
		server.Start(serverAddr)
	},
}

func init() {
	nodeCmd.Flags().StringP("host", "H", "0.0.0.0", "server address")
	nodeCmd.Flags().IntP("port", "p", 5921, "server port")
	rootCmd.AddCommand(nodeCmd)
}
