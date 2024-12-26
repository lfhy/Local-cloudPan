package api

type ApiViewReq struct {
	Short          bool `query:"short"`
	ReplaceImgPath bool `query:"replaceImgPath"`
}

type ApiViewRes struct{}
