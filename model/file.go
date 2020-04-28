package model

import (
	"time"

	"github.com/LyricTian/gin-admin/pkg/errors"
	"github.com/RebelBIrd/fileserver/fastDfs"

	"github.com/snluu/uuid"

	"github.com/qinyuanmao/go-utils/logutl"
	"github.com/qinyuanmao/go-utils/ormutl"
	"github.com/qinyuanmao/go-utils/pageutl"
	"github.com/qinyuanmao/go-utils/strutl"
	"github.com/qinyuanmao/go-utils/timeutl"
)

// type FileType int
// type FileInfo struct {
// 	Name     string
// 	Suffix   string
// 	Type     FileType
// 	Category string
// 	Path     string
// 	MD5      string
// }

type FileNode struct {
	ID       string           `xorm:"NOT NULL PK" json:"fileId"` //文件ID
	Name     string           `json:"fileName"`                  //文件名
	Url      string           `xorm:"NOT NULL" json:"fileUrl"`   //文件外部访问地址
	Path     string           `xorm:"NOT NULL" json:"-"`         //文件存储路径
	CreateAt timeutl.Time     `xorm:"created" json:"createAt"`   //文件创建时间
	UpdateAt timeutl.Time     `xorm:"updated" json:"updateAt"`   //文件更新时间
	MD5      string           `xorm:"NOT NULL" json:"md5"`       //文件MD5
	Suffix   string           `json:"suffix,omitempty"`          //文件后缀
	Type     fastDfs.FileType `json:"type"`                      //文件类型
	Category string           `json:"category"`                  //文件分类目录
	Other    string           `json:"other,omitempty"`           //文件备注
}

func (node *FileNode) Save() (err error) {
	_, err = ormutl.GetEngine().Insert(node)
	if err != nil {
		logutl.Error(err.Error())
	}
	return
}

func (node *FileNode) Update() (err error) {
	if node.ID == "" {
		node.ID = uuid.Rand().Hex()
		return node.Save()
	} else {
		_, err = ormutl.GetEngine().ID(node.ID).Update(node)
		if err != nil {
			logutl.Error(err.Error())
		}
		return
	}
}

func (node *FileNode) GetFileById() (err error) {
	if node.ID == "" {
		return errors.New("File not exist!")
	}
	_, err = ormutl.GetEngine().ID(node.ID).Get(node)
	if err != nil {
		logutl.Error(err.Error())
	}
	return
}

func (node *FileNode) IsFileExist() (exist bool) {
	if node.MD5 == "" {
		return false
	}
	exist, _ = ormutl.GetEngine().Where("MD5 = ?", node.MD5).Get(node)
	return
}

func (node *FileNode) DeleteById() (err error) {
	if node.ID == "" {
		return errors.New("File not exist!")
	}
	err = node.GetFileById()
	if err != nil {
		logutl.Error(err.Error())
	} else {
		err = fastDfs.DeleteFile(node.Path)
		if err != nil {
			logutl.Error(err.Error())
		}
	}
	_, err = ormutl.GetEngine().Exec("Delete From FileNode Where ID = ?", node.ID)
	if err != nil {
		logutl.Error(err.Error())
	}
	return
}

func (node *FileNode) DeleteByMd5() (err error) {
	if node.MD5 == "" {
		return errors.New("File not exist!")
	}
	exist := node.IsFileExist()
	if exist {
		_, err = ormutl.GetEngine().Exec("Delete From FileNode Where MD5 = ?", node.MD5)
		if err != nil {
			logutl.Error(err.Error())
		}
		err = fastDfs.DeleteFile(node.Path)
		if err != nil {
			logutl.Error(err.Error())
		}
	} else {
		return errors.New("File not exist!")
	}
	return
}

func (node *FileNode) DeleteByFileUrl() (err error) {
	if node.Url == "" {
		return errors.New("File not exist!")
	}
	err = node.GetByFileUrl()
	if err != nil {
		logutl.Error(err.Error())
	} else {
		err = fastDfs.DeleteFile(node.Path)
		if err != nil {
			logutl.Error(err.Error())
		}
	}
	if err = node.DeleteById(); err != nil {
		logutl.Error(err.Error())
	}
	return
}

func (node *FileNode) GetByFileUrl() (err error) {
	if node.Url == "" {
		return errors.New("File not exist!")
	}
	_, err = ormutl.GetEngine().Where("Url = ?", node.Url).Get(&node)
	if err != nil {
		logutl.Error(err.Error())
	}
	return
}

func FindFindByType(fileType fastDfs.FileType, afterTime int64, query string, pageSize, pageIndex int) (total int, count int, nodes []FileNode) {
	if pageSize == 0 {
		pageSize = 20
	}
	if pageIndex == 0 {
		pageIndex = 1
	}
	_, _ = ormutl.GetEngine().Select("Count(*)").Table(FileNode{}).Where("Type = ? And UpdateAt > ? And Name Like ?",
		fileType, timeutl.Time(time.Unix(afterTime, 0)), strutl.ConnString("%", query, "%")).Get(&total)
	_ = ormutl.GetEngine().SQL("Select * From FileNode Where Type = ? And UpdateAt > ? And Name Like ? Order By UpdateAt Desc Limit ?,?",
		fileType, timeutl.Time(time.Unix(afterTime, 0)), strutl.ConnString("%", query, "%"),
		pageSize*(pageIndex-1), pageSize).Find(&nodes)
	count = pageutl.GetPageCount(pageSize, total)
	return
}
