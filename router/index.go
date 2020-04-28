package router

import (
	"net/http"
	"regexp"

	"github.com/RebelBIrd/fileserver/utils"
	"github.com/gin-gonic/gin"
	"github.com/qinyuanmao/go-utils/httputl"
	"github.com/qinyuanmao/go-utils/logutl"
)

const FilePath = "files"

var Groups []httputl.BaseGroup

func init() {
	// fileutl.PathExistOrCreate(FilePath)
	Groups = append(Groups, httputl.BaseGroup{
		Path: "file",
		Routers: []httputl.BaseRouter{{
			Type:    httputl.POST,
			Path:    "uploadFile",
			Handler: fileFunc.uploadFile,
		}, {
			Type:    httputl.POST,
			Path:    "uploadFiles",
			Handler: fileFunc.uploadFiles,
		}, {
			Type:    httputl.POST,
			Path:    "queryByMd5",
			Handler: fileFunc.queryByMd5,
		}, {
			Type:    httputl.GET,
			Path:    ":id",
			Handler: fileFunc.getFileById,
		}, {
			Type:    httputl.DELETE,
			Path:    ":id",
			Handler: fileFunc.deleteByFileId,
		}, {
			Type:    httputl.POST,
			Path:    "info",
			Handler: fileFunc.getFileInfoById,
		}, {
			Type:    httputl.POST,
			Path:    "findFile",
			Handler: fileFunc.FindFileByType,
		}},
		Middleware: func(context *gin.Context) {
			match, err := regexp.Match("(\\w{8}(-\\w{4}){3}-\\w{12}?)", []byte(context.Request.URL.Path))
			if err != nil {
				logutl.Error(err.Error())
			}
			if match {
				context.Next()
			} else {
				timestamp := context.Request.Header.Get("timestamp")
				sign := context.Request.Header.Get("sign")
				if timestamp == "" || sign == "" {
					context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "签名认证失败！"))
					context.Abort()
				} else {
					if utils.CheckSign(timestamp, sign) {
						context.Next()
					} else {
						context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "签名认证失败！"))
						context.Abort()
					}
				}
			}
		},
	})
}
