package main

import (
	"fileserver/conf"
	"fileserver/router"

	"github.com/gin-gonic/gin"

	"github.com/qinyuanmao/go-utils/httputl"
)

func main() {
	httputl.StartServer(router.Groups, nil, conf.QsConfig.Port, func(engine *gin.Engine) {
		engine.Static("/files", "./files")
	})
}
