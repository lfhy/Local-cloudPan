package handle

import (
	"bytes"
	"local-cloud-api/api"
	"local-cloud-api/conf"
	"net/url"
	"os"
	"path/filepath"
	"sort"

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

func ListDir(rootDir string, sorts ...string) []*api.FileInfo {
	res := make([]*api.FileInfo, 0)
	sortBy := "name"
	if len(sorts) > 0 {
		sortBy = sorts[0]
	}
	// 打开根目录
	dir, err := os.Open(rootDir)
	if err != nil {
		log.Warn("Error opening directory:", err)
		return res
	}
	defer dir.Close()

	// 读取目录项
	entries, err := dir.Readdir(-1)
	if err != nil {
		log.Warn("Error reading directory entries:", err)
		return res
	}

	// 遍历目录项并打印文件名
	for _, entry := range entries {
		name := entry.Name()
		if conf.IgnoreDotFiles && name[0] == '.' {
			continue
		}
		info := api.FilePathToApiFileInfo(filepath.Join(rootDir, name), entry)
		res = append(res, &info)
	}
	// 排序
	switch sortBy {
	default:
		sort.Slice(res, func(i, j int) bool {
			if res[i].IsDir && res[j].IsDir {
				return res[i].Name < res[j].Name
			}
			if res[i].IsDir {
				return true
			}
			if res[j].IsDir {
				return false
			}
			return res[i].Name < res[j].Name
		})
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
