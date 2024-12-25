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
	files := ListDir(realPath, req.SortMode)
	return &api.ApiFileListRes{FileList: files}, nil
	// empty := make([]string, 0)
	// return &api.ApiEmptyFileList{FileList: empty}, nil
}
