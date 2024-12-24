package handle

import (
	"bytes"
	"local-cloud-api/api"
	"local-cloud-api/conf"
	"net/url"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/lfhy/log"
)

func ChangeToSysPath(path string) string {
	absPath := filepath.Join(conf.ShareFilePath, path)
	path, err := url.PathUnescape(absPath)
	if err == nil {
		absPath = path
	}
	return absPath
}

func ListDir(rootDir string) []*api.FileInfo {
	// 打开根目录
	dir, err := os.Open(rootDir)
	if err != nil {
		log.Warn("Error opening directory:", err)
		return nil
	}
	defer dir.Close()

	// 读取目录项
	entries, err := dir.Readdir(-1)
	if err != nil {
		log.Warn("Error reading directory entries:", err)
		return nil
	}
	var res []*api.FileInfo
	// 遍历目录项并打印文件名
	for _, entry := range entries {
		info := api.FilePathToApiFileInfo(filepath.Join(rootDir, entry.Name()), entry)
		res = append(res, &info)
	}
	return res
}

// 获取缩略图
func GetShortImg(imagePath string) ([]byte, error) {
	img, err := imaging.Open(imagePath)
	if err != nil {
		return nil, err
	}
	buffer := new(bytes.Buffer)
	thumb := imaging.Thumbnail(img, 300, 300, imaging.CatmullRom)
	err = imaging.Encode(buffer, thumb, imaging.JPEG)
	return buffer.Bytes(), err
}
