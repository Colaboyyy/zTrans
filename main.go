package main

import (
	"github.com/zserge/lorca"
	"os"
	"os/signal"
	"syscall"
	"zTrans/config"
	"zTrans/server"
)

func main() {
	// 启动gin服务
	go server.Run()

	// lorca启动chrome
	var ui lorca.UI
	ui, _ = lorca.New("http://127.0.0.1"+config.GetPort()+"/static/index.html", "", 800, 600, "--disable-sync", "--disable-translate")

	chSignal := listenToInterrupt()
	select {
	case <-ui.Done():
	case <-chSignal:
	}
	ui.Close()
}

// 监听中断信号
func listenToInterrupt() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	return chSignal
}
