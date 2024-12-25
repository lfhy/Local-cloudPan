package handle

import (
	"encoding/json"
	"io"
	"local-cloud-api/api"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
)

func init() {
	RegisterHandle("/merge", http.MethodPost, merge)
}

// 合并文件
func merge(ctx fiber.Ctx, req *api.ApiMergeReq) (*api.ApiMergeRes, error) {
	workDir := GetFileUploadPath(req.FileId)
	fileInfoPath := filepath.Join(workDir, "fileInfo")
	var saveReq api.ApiVerifyReq
	data, err := os.ReadFile(fileInfoPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &saveReq)
	if err != nil {
		return nil, err
	}
	// 合块
	savePath := ChangeToSysPath(req.Path, saveReq.FileName)
	saveFile, err := os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer saveFile.Close()
	for _, chunk := range saveReq.Chunks {
		chunkF, err := os.Open(filepath.Join(workDir, chunk.ChunkId))
		if err != nil {
			os.Remove(savePath)
			return nil, err
		}
		_, err = io.Copy(saveFile, chunkF)
		if err != nil {
			chunkF.Close()
			os.Remove(savePath)
			return nil, err
		}
		chunkF.Close()
	}
	SetResMsg(ctx, "文件上传成功")
	os.RemoveAll(workDir)
	return nil, nil
}
