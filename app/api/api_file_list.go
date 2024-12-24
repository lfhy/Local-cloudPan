package api

import (
	"io/fs"
	"local-cloud-api/conf"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var PicType = []string{"jpeg", "jpg", "png", "svg", "gif", "webp"}

type ApiFileListReq struct {
	Path     string `query:"path"`
	SortMode string `query:"sortMode"`
}

type ApiFileListRes struct {
	FileList []*FileInfo `json:"fileList"`
}

type FileInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	IsDir         bool   `json:"isDir"`
	Ext           string `json:"ext"`
	Size          int64  `json:"size"`
	Modified      int64  `json:"modified"`
	FilePath      string `json:"filePath"`
	ThumbnailPath string `json:"thumbnailPath,omitempty"`
}

func FilePathToApiFileInfo(absPath string, info ...fs.FileInfo) FileInfo {
	var fi fs.FileInfo
	if len(info) > 0 {
		fi = info[0]
	} else {
		fi, _ = os.Stat(absPath)
	}
	absPath = strings.TrimPrefix(absPath, conf.ShareFilePath)
	res := FileInfoToApiFilInfo(fi)
	absPath = strings.ReplaceAll(url.PathEscape(absPath), "%2F", "/")
	res.ID = absPath
	res.FilePath = conf.ApiPrefix + "/view" + absPath
	if slices.Index(PicType, res.Ext) >= 0 {
		res.ThumbnailPath = res.FilePath + "?short=true"
	}
	return res
}

func FileInfoToApiFilInfo(fi fs.FileInfo) FileInfo {
	ext := strings.TrimPrefix(filepath.Ext(fi.Name()), ".")

	return FileInfo{
		Name:     fi.Name(),
		IsDir:    fi.IsDir(),
		Ext:      ext,
		Size:     fi.Size(),
		Modified: fi.ModTime().UnixMilli(),
		// ThumbnailPath: ThumbnailPath,
	}
}
