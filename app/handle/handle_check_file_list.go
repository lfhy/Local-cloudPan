package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/checkFileList", http.MethodGet, checkFileList)
}

// 校验文件是否存在
func checkFileList(ctx fiber.Ctx, req *api.ApiCheckFileListReq) (*api.ApiCheckFileListRes, error) {
	return nil, api.ErrorNoImp
}
