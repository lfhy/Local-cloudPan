package api

var PicType = []string{"jpeg", "jpg", "png", "svg", "gif", "webp"}

type ApiFileListReq struct{}

type ApiFileListRes struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IsDir         bool   `json:"isDir"`
	Ext           string `json:"ext"`
	Size          int64  `json:"size"`
	Modified      int64  `json:"modified"`
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath,omitempty"`
}
