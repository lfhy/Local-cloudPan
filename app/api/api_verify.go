package api

type ApiVerifyReq struct {
	FileId   string   `json:"fileId"`   // 文件ID
	FileName string   `json:"fileName"` // 文件名
	Chunks   []*Chunk `json:"chunks"`   // 分块信息
}

type Chunk struct {
	ChunkId string `json:"chunkId"` // 分块ID
	Index   int    `json:"index"`   // 分块索引
}

type ApiVerifyRes struct {
	Code   int      `json:"code"` // 状态码
	Chunks []string `json:"shouldUpload"`
}
