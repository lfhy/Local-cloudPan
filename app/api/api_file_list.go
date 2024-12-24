package api

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var PicType = []string{"jpeg", "jpg", "png", "svg", "gif", "webp"}

type ApiFileListReq struct {
	Path     string `query:"path"`
	SortMode string `query:"sortMode"`
}

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

func FilePathToApiFileListRes(absPath string, info ...fs.FileInfo) ApiFileListRes {
	var fi fs.FileInfo
	if len(info) > 0 {
		fi = info[0]
	} else {
		fi, _ = os.Stat(absPath)
	}
	res := FileInfoToApiFileListRes(fi)
	res.ID = absPath
	res.FilePath = absPath
	return res
}

func FileInfoToApiFileListRes(fi fs.FileInfo) ApiFileListRes {
	return ApiFileListRes{
		Name:          fi.Name(),
		IsDir:         fi.IsDir(),
		Ext:           strings.TrimPrefix(filepath.Ext(fi.Name()), "."),
		Size:          fi.Size(),
		Modified:      fi.ModTime().Unix(),
		ThumbnailPath: "", // TODO:
	}
}
