package handle

import (
	"local-cloud-api/api"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/upload", http.MethodPost, upload)
}

// 合并文件
func upload(ctx fiber.Ctx, req *api.ApiUploadReq) (*api.ApiUploadRes, error) {
	// form, err := c.MultipartForm()
	// if err != nil {
	// 	return api.ResError(c, err)
	// }

	// 处理文件内容
	// return api.ResOK(c, "分片上传成功", nil)
	return nil, api.ErrorNoImp
}
