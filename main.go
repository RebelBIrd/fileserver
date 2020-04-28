package main

import (
	"github.com/gin-gonic/gin"

	"github.com/RebelBIrd/fileserver/conf"
	_ "github.com/RebelBIrd/fileserver/conf"
	"github.com/RebelBIrd/fileserver/router"
	_ "github.com/RebelBIrd/fileserver/router"
	"github.com/qinyuanmao/go-utils/httputl"
)

func main() {
	httputl.StartServer(router.Groups, nil, conf.QsConfig.Port, func(engine *gin.Engine) {
		engine.Static("/files", "./files")
	})
}
