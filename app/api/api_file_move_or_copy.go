package api

type ApiFileMoveOrCopyReq struct {
	Path        string   `json:"path"`
	FileList    []string `json:"fileList"`
	Destination string   `json:"destination"`
	Dtype       string   `json:"dtype"` //move or copy
}

type ApiFileMoveOrCopyRes struct{}
