package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/rename", http.MethodPost, rename)
}

// 重命名
func rename(ctx fiber.Ctx, req *api.ApiRenameReq) (*api.ApiRenameRes, error) {
	return nil, api.ErrorNoImp
}
