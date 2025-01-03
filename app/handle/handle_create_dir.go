package handle

import (
	"local-cloud-api/api"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/createDir", http.MethodPost, createDir)
}

// 文件夹创建
func createDir(ctx fiber.Ctx, req *api.ApiCreateDirReq) (*api.ApiCreateDirRes, error) {
	return nil, os.Mkdir(ChangeToSysPath(filepath.Join(req.Path, req.DirName)), os.ModeDir)
}
