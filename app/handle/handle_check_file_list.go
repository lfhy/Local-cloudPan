package handle

import (
	"fmt"
	"local-cloud-api/api"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

func init() {
	RegisterHandle("/checkFileList", http.MethodGet, checkFileList)
}

// 校验文件是否存在
func checkFileList(ctx fiber.Ctx, req *api.ApiCheckFileListReq) (*api.ApiCheckFileListRes, error) {
	log.Printf("校验文件列表: %+v\n", req)
	for _, file := range req.FilenameLists {
		_, err := os.Stat(ChangeToSysPath(file))
		if err != nil {
			SetResMsg(ctx, fmt.Sprintf("%v 文件不存在！本次操作无效！", file))
			return nil, api.ErrorCheckFileFailed
		}
	}
	return nil, nil
}
