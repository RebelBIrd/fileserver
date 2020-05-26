package router

import (
	"fileserver/conf"
	"fileserver/fastDfs"
	"fileserver/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qinyuanmao/go-utils/httputl"
	"github.com/qinyuanmao/go-utils/logutl"
	"github.com/snluu/uuid"
)

type FileFunc struct{}

var fileFunc = FileFunc{}

func (file FileFunc) uploadFile(context *gin.Context) {
	_ = context.Request.ParseMultipartForm(2000000000)
	md5 := httputl.GetParam(context, "md5")
	fileName := httputl.GetParam(context, "fileName")
	if md5 == "" {
		context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_ParamNoFound), "md5参数不能为空！"))
	} else {
		node := model.FileNode{MD5: md5}
		if exist := node.IsFileExist(); exist {
			context.JSON(http.StatusOK, httputl.RespSuccess(node))
		} else {
			file, _ := context.FormFile("file")
			info := fastDfs.GinFileHandler(file, FilePath, fileName)
			if md5 == info.MD5 {
				id := uuid.Rand().Hex()
				node.ID = id
				node.Name = info.Name
				node.Path = info.Path
				node.Url = conf.QsConfig.ApiConf.Url + "/" + info.Path
				node.Suffix = info.Suffix
				node.Type = info.Type
				node.Category = info.Category
				node.MD5 = info.MD5
				if err := node.Save(); err != nil {
					context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), err.Error()))
				} else {
					context.JSON(http.StatusOK, httputl.RespSuccess(node))
				}
			} else {
				_ = fastDfs.DeleteFile(info.Path)
				context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "md5与文件不对应！"))
			}
		}
	}
}

func (file FileFunc) uploadFiles(context *gin.Context) {
	_ = context.Request.ParseMultipartForm(2000000000)
	var nodes []model.FileNode
	form, _ := context.MultipartForm()
	files := form.File["files"]
	for _, file := range files {
		info := fastDfs.GinFileHandler(file, FilePath, "")
		node := model.FileNode{MD5: info.MD5}
		if exist := node.IsFileExist(); exist {
			nodes = append(nodes, node)
			continue
		} else {
			id := uuid.Rand().Hex()
			node = model.FileNode{
				ID:       id,
				Name:     info.Name,
				Path:     info.Path,
				Url:      conf.QsConfig.ApiConf.Url + "/file/" + id,
				Suffix:   info.Suffix,
				Type:     info.Type,
				MD5:      info.MD5,
				Category: info.Category,
			}
		}
		if err := node.Save(); err != nil {
			logutl.Error(err.Error())
		} else {
			nodes = append(nodes, node)
		}
	}
	context.JSON(http.StatusOK, httputl.RespSuccess(nodes))
}

func (file FileFunc) queryByMd5(context *gin.Context) {
	md5 := httputl.GetParam(context, "md5")
	if md5 == "" {
		context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "请求路径异常！"))
	} else {
		fileNode := model.FileNode{MD5: md5}
		_ = fileNode.IsFileExist()
		if fileNode.ID == "" {
			context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "File not Fount!"))
		} else {
			context.JSON(http.StatusOK, httputl.RespSuccess(fileNode))
		}
	}
}

func (file FileFunc) deleteByFileId(context *gin.Context) {
	fileId := httputl.GetParam(context, "id")
	if fileId == "" {
		context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "请求路径异常！"))
	} else {
		node := model.FileNode{ID: fileId}
		err := node.DeleteById()
		if err != nil {
			context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), err.Error()))
		} else {
			context.JSON(http.StatusOK, httputl.RespDefaultSuccess())
		}
	}
}

func (file FileFunc) getFileById(context *gin.Context) {
	fileId := httputl.GetParam(context, "id")
	if fileId == "" {
		context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "请求路径异常！"))
	} else {
		fileNode := model.FileNode{ID: fileId}
		err := fileNode.GetFileById()
		if err != nil {
			context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), err.Error()))
		} else if fileNode.MD5 == "" {
			context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "File not Fount!"))
		} else {
			context.Writer.Header().Add("content-disposition", `attachment; filename=`+fileNode.Name+fileNode.Suffix)
			context.File(fileNode.Path)
		}
	}
}

func (file FileFunc) getFileInfoById(context *gin.Context) {
	fileId := httputl.GetParam(context, "id")
	if fileId == "" {
		context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "请求路径异常！"))
	} else {
		fileNode := model.FileNode{ID: fileId}
		err := fileNode.GetFileById()
		if err != nil {
			context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), err.Error()))
		} else if fileNode.MD5 == "" {
			context.JSON(http.StatusOK, httputl.RespFailed(int(httputl.RPCD_Failed), "File not Fount!"))
		} else {
			context.JSON(http.StatusOK, httputl.RespSuccess(fileNode))
		}
	}
}

func (file FileFunc) FindFileByType(context *gin.Context) {
	fileType := httputl.GetIntParam(context, "type")
	afterTime := httputl.GetInt64Param(context, "afterTime")
	query := httputl.GetParam(context, "query")
	pageIndex := httputl.GetIntParam(context, "pageIndex")
	pageSize := httputl.GetIntParam(context, "pageSize")
	total, _, data := model.FindFindByType(fastDfs.FileType(fileType), afterTime, query, pageSize, pageIndex)
	context.JSON(http.StatusOK, httputl.RespArraySuccess(pageIndex, pageSize, total, data))
}
