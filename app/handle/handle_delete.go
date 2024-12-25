package handle

import (
	"local-cloud-api/api"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/delete", http.MethodDelete, deleteFile)
}

// 删除文件
func deleteFile(ctx fiber.Ctx, req *api.ApiDeleteReq) (*api.ApiCheckFileListRes, error) {
	for _, file := range req.FilenameList {
		os.Remove(ChangeToSysPath(req.Path, file))
	}

	return nil, nil
}
