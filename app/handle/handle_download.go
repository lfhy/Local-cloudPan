package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/download", http.MethodGet, download)
}

// 下载文件
func download(ctx fiber.Ctx, req *api.ApiDownloadReq) (*api.ApiDownloadRes, error) {
	return nil, api.ErrorNoImp
}
