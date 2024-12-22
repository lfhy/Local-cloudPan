package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/delete", http.MethodPost, fileMoveOrCopy)
}

// 文件移动或复制
func fileMoveOrCopy(ctx fiber.Ctx, req *api.ApiFileMoveOrCopyReq) (*api.ApiFileMoveOrCopyRes, error) {
	return nil, api.ErrorNoImp
}
