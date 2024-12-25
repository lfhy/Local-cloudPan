package handle

import (
	"local-cloud-api/api"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

func init() {
	RegisterHandle("/view/*", http.MethodGet, view)
}

// 获取文件
func view(ctx fiber.Ctx, req *api.ApiViewReq) (*api.ApiViewRes, error) {
	filePath := ctx.Params("*")
	absPath := ChangeToSysPath(filePath)
	log.Info("filePath:", absPath)
	ext := strings.TrimPrefix(filepath.Ext(absPath), ".")
	ctx.Type(ext).Response().Header.Set("Cache-Control", "max-age=86400")
	if req.Short {
		data, err := GetShortImg(absPath)
		if err == nil {
			ctx.Send(data)
			return nil, api.ErrorNoRes
		}
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	ctx.Send(data)
	return nil, api.ErrorNoRes
}
