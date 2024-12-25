package api

type ApiCreateDirReq struct {
	Path    string `json:"path"`    // 当前路径
	DirName string `json:"dirName"` // 新建文件夹名称
}

type ApiCreateDirRes struct{}
