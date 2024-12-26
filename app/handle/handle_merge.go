package handle

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
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
		chunkF, err := os.Open(GetChunkPath(req.FileId, chunk.Index))
		if err != nil {
			os.Remove(savePath)
			return nil, err
		}
		hash, reader := Md5Read(chunkF)
		_, err = io.Copy(saveFile, reader)
		if err != nil {
			chunkF.Close()
			os.Remove(savePath)
			return nil, err
		}
		chunkF.Close()
		if hex.EncodeToString(hash.Sum(nil)) != chunk.ChunkId {
			os.Remove(savePath)
			return nil, fmt.Errorf("分块校验失败文件合并失败")
		}

	}
	SetResMsg(ctx, "文件上传成功")
	os.RemoveAll(workDir)
	return nil, nil
}
