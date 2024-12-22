package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/fileList", http.MethodGet, fileList)
}

// 文件列表
func fileList(ctx fiber.Ctx, req *api.ApiFileListReq) (*api.ApiFileListRes, error) {
	return nil, api.ErrorNoImp
}
