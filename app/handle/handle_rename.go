package handle

import (
	"local-cloud-api/api"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/rename", http.MethodPost, rename)
}

// 重命名
func rename(ctx fiber.Ctx, req *api.ApiRenameReq) (*api.ApiRenameRes, error) {
	// 校验文件是否不存在
	fromPath := ChangeToSysPath(req.Path, req.OldName)
	_, err := os.Stat(fromPath)
	if err != nil {
		return nil, api.ErrorCheckFileFailed
	}
	// 校验文件是否已经存在
	toPath := ChangeToSysPath(req.Path, req.NewName)
	_, err = os.Stat(toPath)
	if err == nil {
		return nil, api.ErrorRenameFailed
	}
	return nil, os.Rename(fromPath, toPath)
}
