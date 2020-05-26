package fastDfs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"path"
	"strings"

	"github.com/qinyuanmao/go-utils/strutl"
	"github.com/tedcy/fdfs_client"
)

type FileType int

const (
	FT_UNDEFINED FileType = iota
	FT_VIDEO
	FT_IMAGE
	FT_DOC
	FT_APP
	FT_ZIP
	FT_MUSIC
	FT_OTHRER
)

type FileInfo struct {
	Name     string
	Suffix   string
	Type     FileType
	Category string
	Path     string
	MD5      string
}

var fastClient *fdfs_client.Client

func GetFastDfsClient() (*fdfs_client.Client, error) {
	var err error
	if fastClient == nil {
		fastClient, err = fdfs_client.NewClientWithConfig("./conf/fdfs.conf")
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return fastClient, nil
	}
	return fastClient, nil
}

func init() {
	GetFastDfsClient()
}

func GinFileHandler(v *multipart.FileHeader, path, fileName string) (info *FileInfo) {
	if fileName == "" {
		fileName = v.Filename
	}
	fileName = strutl.UrlDecode(fileName)
	if fileName == "" {
		return
	} else {
		rename, fileSuffix := DecodeFileName(fileName)
		fileType, category, fileExtName := GetFileType(fileName)
		// PathExistOrCreate(path)
		// PathExistOrCreate(path + "/" + category)
		// if PathExist(info.Path) {
		// 	saveName := rename + strutl.GetRandomString(4)
		// 	info.Path = path + "/" + category + "/" + saveName + fileSuffix
		// }
		// file, _ := v.Open()
		// defer file.Close()
		// out, _ := os.Create(info.Path)
		// defer out.Close()
		// if _, err := io.Copy(out, file); err != nil {
		// 	fmt.Println(err.Error())
		// 	return nil
		// }
		client, err := GetFastDfsClient()
		// defer client.Destory()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fileContent, _ := v.Open()

		byteContainer, err := ioutil.ReadAll(fileContent)
		fileID, err := client.UploadByBuffer(byteContainer, fileExtName)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(fileID)
		info = &FileInfo{
			Name:     rename,
			Path:     fileID,
			Suffix:   fileSuffix,
			Category: category,
			Type:     fileType,
		}
		info.MD5 = GetMd5(v)
		return
	}
}
func DeleteFile(fileID string) error {
	client, err := GetFastDfsClient()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if err = client.DeleteFile(fileID); err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
func GetMd5(v *multipart.FileHeader) string {
	file, err := v.Open()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer file.Close()
	md5h := md5.New()
	_, err = io.Copy(md5h, file)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	} else {
		return hex.EncodeToString(md5h.Sum(nil))
	}
}

func DecodeFileName(filePath string) (fileName, fileType string) {
	filePath = path.Base(filePath)
	fileType = path.Ext(filePath)
	fileName = strings.TrimSuffix(filePath, fileType)
	switch fileType {
	case ".gz":
		if strings.HasSuffix(fileName, ".tar") {
			fileType = ".tar.gz"
			fileName = strings.TrimSuffix(fileName, ".tar")
		}
	case ".exe":
		if strings.HasSuffix(fileName, ".asp") {
			fileType = ".asp.exe"
			fileName = strings.TrimSuffix(fileName, ".exe")
		}
	}
	return
}
func GetFileType(fileName string) (FileType, string, string) {
	_, fileType := DecodeFileName(fileName)
	if fileType == "" {
		return FT_OTHRER, "other", fileType
	}
	fileType = strings.ToLower(fileType)[1:]
	switch fileType {
	case "avi", "mov", "rmvb", "fmv", "mp4", "3gp", "mkv", "f4v":
		return FT_VIDEO, "video", fileType
	case "bmp", "jpg", "png", "ico", "tif", "gif", "pcx", "tga", "exif", "fpx", "svg", "psd", "cdr", "pcd", "dxf", "ufo", "eps", "ai", "raw", "WMF", "webp", "jpeg":
		return FT_IMAGE, "image", fileType
	case "txt", "doc", "docx", "ppt", "pptx", "wpf", "md", "xls", "xlsx", "pdf":
		return FT_DOC, "doc", fileType
	case "apk", "app", "exe", "ipa":
		return FT_APP, "app", fileType
	case "zip", "tar.gz", "dmg", "rar", "tar", "7z", "iso", "bz2":
		return FT_ZIP, "zip", fileType
	case "cd", "aiff", "mp3", "wma", "ogg", "acc", "amr", "mid":
		return FT_MUSIC, "music", fileType
	default:
		return FT_OTHRER, "other", fileType
	}
}
