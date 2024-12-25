package api

type ApiDeleteReq struct {
	Path         string   `query:"path"`
	FilenameList []string `query:"filenameList"`
}

type ApiDeleteRes struct{}
