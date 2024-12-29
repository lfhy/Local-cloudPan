package api

type ApiRenameReq struct {
	Path    string `json:"path"`
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

type ApiRenameRes struct{}
