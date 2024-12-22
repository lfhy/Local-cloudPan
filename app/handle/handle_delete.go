package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/delete", http.MethodDelete, deleteFile)
}

// 删除文件
func deleteFile(ctx fiber.Ctx, req *api.ApiDeleteReq) (*api.ApiCheckFileListRes, error) {
	return nil, api.ErrorNoImp
}
