package server

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"zTrans/config"
	"zTrans/server/controller"
	"zTrans/server/ws"
)

//由于go embed不支持.. 需要先执行go generate将前端生成文件拷贝到当前目录下
//go:generate cp -r ../frontend/dist ./frontend_assets
//go:embed frontend_assets/dist/*
var Fs embed.FS

func Run() {
	r := gin.Default()
	staticFiles, _ := fs.Sub(Fs, "frontend_assets/dist")
	r.StaticFS("/static", http.FS(staticFiles))
	r.POST("/api/v1/files", controller.FileController)
	r.GET("/api/v1/qrcodes", controller.QrcodeController)
	r.GET("/uploads/:path", controller.UploadsController)
	r.GET("/api/v1/addresses", controller.AddressesController)
	r.POST("/api/v1/texts", controller.TextController)
	hub := ws.NewHub()
	go hub.Run()
	r.GET("/ws", func(c *gin.Context) {
		ws.WsController(c, hub)
	})
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	r.Run(config.GetPort())
}
