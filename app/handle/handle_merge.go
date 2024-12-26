package handle

import (
	"encoding/hex"
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
	var needs []string
	var isDeletePath bool
	for _, chunk := range saveReq.Chunks {
		if isDeletePath {
			md5, _ := Md5File(GetChunkPath(req.FileId, chunk.Index))
			if md5 != chunk.ChunkId {
				needs = append(needs, chunk.ChunkId)
			}
			continue
		}
		chunkF, err := os.Open(GetChunkPath(req.FileId, chunk.Index))
		if err != nil {
			if !isDeletePath {
				os.RemoveAll(savePath)
				isDeletePath = true
			}
			// return nil, err
			needs = append(needs, chunk.ChunkId)
			continue
		}
		hash, reader := Md5Read(chunkF)
		_, err = io.Copy(saveFile, reader)
		if err != nil {
			chunkF.Close()
			if !isDeletePath {
				os.RemoveAll(savePath)
				isDeletePath = true
			}
			needs = append(needs, chunk.ChunkId)
			continue
		}
		chunkF.Close()
		if hex.EncodeToString(hash.Sum(nil)) != chunk.ChunkId {
			if !isDeletePath {
				os.RemoveAll(savePath)
				isDeletePath = true
			}
			// return nil, fmt.Errorf("分块校验失败文件合并失败")
			needs = append(needs, chunk.ChunkId)
			continue
		}
	}
	if len(needs) > 0 {
		return &api.ApiMergeRes{NeedChunkIds: needs}, api.ErrorCheckFileFailed
	}
	SetResMsg(ctx, "文件上传成功")
	os.RemoveAll(workDir)
	return nil, nil
}
