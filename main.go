package main

import (
	"github.com/gin-gonic/gin"

	"github.com/qinyuanmao/fileserver/conf"
	_ "github.com/qinyuanmao/fileserver/conf"
	"github.com/qinyuanmao/fileserver/router"
	_ "github.com/qinyuanmao/fileserver/router"
	"github.com/qinyuanmao/go-utils/httputl"
)

func main() {
	httputl.StartServer(router.Groups, nil, conf.QsConfig.Port, func(engine *gin.Engine) {
		engine.Static("/files", "./files")
	})
}
