package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func init() {
	RegisterHandle("/thumbnail/*", http.MethodGet, thumbnail)
}

// 图片文件压缩
func thumbnail(ctx fiber.Ctx, req *api.ApiFileListReq) (*api.ApiFileListRes, error) {
	filePath := ctx.Params("*")
	log.Info("filePath:", filePath)
	return nil, api.ErrorNoImp
}
