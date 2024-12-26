package handle

import (
	"bytes"
	"fmt"
	"io"
	"local-cloud-api/api"
	"local-cloud-api/conf"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/lfhy/log"
)

func ChangeToSysPath(path ...string) string {
	absPath := filepath.Join(conf.ShareFilePath, filepath.Join(path...))
	upath, err := url.PathUnescape(absPath)
	if err == nil {
		absPath = upath
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
	case "size":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Size > res[j].Size
		})
	case "modified":
		sort.Slice(res, func(i, j int) bool {
			return res[i].Modified > res[j].Modified
		})
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

// 获取文件上传路径
func GetFileUploadPath(fileId string) string {
	return ChangeToSysPath(".uploads", fileId)
}

func StatUntilFileNameOK(dest string) string {
	name := dest
	ext := filepath.Ext(dest)
	noext := strings.TrimSuffix(name, ext)
	count := 0
	for {
		_, err := os.Stat(name)
		if err == nil {
			name = fmt.Sprintf("%v-%v%v", noext, count, ext)
			count++
		} else {
			break
		}
	}
	return name
}

func Copy(src, dest string) error {
	// 打开源文件
	sourceFile, err := os.Open(src)
	if err != nil {
		log.Errorln("打开源失败:", err)
		return err
	}
	defer sourceFile.Close()
	dest = StatUntilFileNameOK(dest)
	// 创建目标文件
	destinationFile, err := os.Create(dest)
	if err != nil {
		log.Errorln("创建目标文件失败:", err)
		return err
	}
	defer destinationFile.Close()

	// 复制文件内容
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		log.Errorln("复制文件失败:", err)
		return err
	}
	return nil
}
