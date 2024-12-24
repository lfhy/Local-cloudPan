package handle

import (
	"local-cloud-api/api"
	"local-cloud-api/conf"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3/log"
)

func ChangeToSysPath(path string) string {
	return filepath.Join(conf.ShareFilePath, path)
}

func ListDir(rootDir string) []*api.ApiFileListRes {
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
	var res []*api.ApiFileListRes
	// 遍历目录项并打印文件名
	for _, entry := range entries {
		info := api.FilePathToApiFileListRes(entry.Name(), entry)
		res = append(res, &info)
	}
	return res
}
