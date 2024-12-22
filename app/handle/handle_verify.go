package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/verify", http.MethodPost, verify)
}

// 合并文件
func verify(ctx fiber.Ctx, req *api.ApiVerifyReq) (*api.ApiVerifyRes, error) {
	return nil, api.ErrorNoImp
}
