package handle

import (
	"local-cloud-api/api"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/fileMoveOrCopy", http.MethodPost, fileMoveOrCopy)
}

// 文件移动或复制
func fileMoveOrCopy(ctx fiber.Ctx, req *api.ApiFileMoveOrCopyReq) (*api.ApiFileMoveOrCopyRes, error) {
	for _, file := range req.FileList {
		src := ChangeToSysPath(req.Path, file)
		dest := ChangeToSysPath(req.Destination, file)

		err := Copy(src, dest)
		if err != nil {
			return nil, err
		}
		if req.Dtype == "move" {
			os.RemoveAll(src)
		}
	}
	if req.Dtype == "move" {
		SetResMsg(ctx, "移动成功")
	} else {
		SetResMsg(ctx, "复制成功")
	}
	return nil, nil
}
