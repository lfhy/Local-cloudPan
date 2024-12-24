package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

func init() {
	RegisterHandle("/fileList", http.MethodGet, fileList)
}

// 文件列表
func fileList(ctx fiber.Ctx, req *api.ApiFileListReq) (*api.ApiFileListRes, error) {
	log.Debugf("req: %+v\n", req)
	realPath := ChangeToSysPath(req.Path)
	log.Debug("访问路径:", realPath)
	return &api.ApiFileListRes{FileList: ListDir(realPath)}, nil
}
