package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/merge", http.MethodPost, merge)
}

// 合并文件
func merge(ctx fiber.Ctx, req *api.ApiMergeReq) (*api.ApiMergeRes, error) {
	return nil, api.ErrorNoImp
}
