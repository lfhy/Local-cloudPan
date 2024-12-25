package api

type ApiUploadReq struct {
	FileId  string `form:"fileId"`
	Index   int    `form:"index"`
	ChunkId string `form:"chunkId"`
	// ChunkData []byte `form:"chunkData"`
}

type ApiUploadRes struct{}
