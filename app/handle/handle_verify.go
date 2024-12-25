package handle

import (
	"encoding/json"
	"local-cloud-api/api"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/verify", http.MethodPost, verify)
}

// 上传校验
func verify(ctx fiber.Ctx, req *api.ApiVerifyReq) (*api.ApiVerifyRes, error) {
	// 判断分块目录是否存在
	workDir := GetFileUploadPath(req.FileId)
	fileInfoPath := filepath.Join(workDir, "fileInfo")
	var saveReq api.ApiVerifyReq
	_, err := os.Stat(fileInfoPath)
	// 文件存在则读取进行解析
	if err == nil {
		data, _err := os.ReadFile(fileInfoPath)
		if _err == nil {
			err = json.Unmarshal(data, &saveReq)
		} else {
			err = _err
		}
		if err != nil {
			os.Remove(fileInfoPath)
		}
	}
	// 读取失败或者解析失败则说明文件不存在
	if err != nil {
		// 上传目录不存在则说明全部需要上传
		os.MkdirAll(workDir, os.ModeDir)
		// 保存上传文件信息
		data, _ := json.Marshal(req)
		os.WriteFile(fileInfoPath, data, os.ModePerm)
		Chunks := make([]string, 0)
		for _, chunk := range req.Chunks {
			Chunks = append(Chunks, chunk.ChunkId)
		}
		ctx.JSON(&api.ApiVerifyRes{
			Code:   200,
			Chunks: Chunks,
		})
		return nil, api.ErrorNoRes
	}
	// 文件存在则判断是否需要上传
	Chunks := make([]string, 0)
	for _, chunk := range saveReq.Chunks {
		_, err := os.Stat(filepath.Join(workDir, chunk.ChunkId))
		if err != nil {
			Chunks = append(Chunks, chunk.ChunkId)
		}
	}
	ctx.JSON(&api.ApiVerifyRes{
		Code:   200,
		Chunks: Chunks,
	})
	return nil, api.ErrorNoRes
}
