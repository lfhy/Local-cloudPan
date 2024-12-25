package api

type ApiCheckFileListReq struct {
	Path          string   `query:"path"`
	FilenameLists []string `query:"filenameList"`
}

type ApiCheckFileListRes struct{}
