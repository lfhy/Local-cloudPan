package api

type ApiCheckFileListReq struct {
	FilenameLists []string `query:"filenameList"`
}

type ApiCheckFileListRes struct{}
