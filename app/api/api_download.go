package api

type ApiDownloadReq struct {
	FilenameLists string `query:"filenameList"`
}

type ApiDownloadRes struct{}
