package main

import (
	"GoScheduler/internal/modules/rpc/server"
	"flag"
	"log"
	"os"
	"runtime"
)

func main() {
	var serverAddr string
	var allowRoot bool
	var version bool
	flag.BoolVar(&allowRoot, "allow-root", false, "./scheduler-node -allow-root")
	flag.StringVar(&serverAddr, "s", "0.0.0.0:5921", "./scheduler-node -s ip:port")
	flag.BoolVar(&version, "v", false, "./scheduler-node -v")
	flag.Parse()

	if version {

		return
	}

	if runtime.GOOS != "windows" && os.Getuid() == 0 && !allowRoot {
		log.Fatal("Do not run gocron-node as root user")
		return
	}

	server.Start(serverAddr)
}
