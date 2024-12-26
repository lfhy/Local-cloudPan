package api

import (
	"fmt"
	"local-cloud-api/conf"
	"path/filepath"
)

type ApiUploadReq struct {
	FileId  string `form:"fileId"`
	Index   int    `form:"index"`
	ChunkId string `form:"chunkId"`
	// ChunkData []byte `form:"chunkData"`
}

func (req *ApiUploadReq) GetChunkPath() string {
	return filepath.Join(conf.UploadTmpPath, req.FileId, fmt.Sprintf("%05d.part", req.Index))
}

type ApiUploadRes struct{}
