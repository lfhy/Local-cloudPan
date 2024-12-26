package api

type ApiMergeReq struct {
	FileId string `json:"fileId"`
	Path   string `json:"path"`
}

type ApiMergeRes struct {
	NeedChunkIds []string `json:"needs"`
}
