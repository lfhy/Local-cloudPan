package handle

import (
	"archive/zip"
	"io"
	"local-cloud-api/api"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

func init() {
	RegisterHandle("/download", http.MethodGet, download)
}

// 下载文件
func download(ctx fiber.Ctx, req *api.ApiDownloadReq) (*api.ApiDownloadRes, error) {
	log.Debug("下载文件列表:", req.FilenameLists)
	files := strings.Split(strings.TrimSuffix(strings.TrimPrefix(req.FilenameLists, "["), "]"), ",")
	ctx.Response().Header.Set("Content-Type", "application/octet-stream;charset=utf-8-sig")
	if len(files) == 1 {
		ctx.Response().Header.Set("Content-Disposition", "attachment;filename="+filepath.Base(files[0]))
		realPath := ChangeToSysPath(files[0])
		log.Debugln("读取文件路径:", realPath)
		ctx.SendFile(realPath, fiber.SendFile{Compress: true})
		return nil, api.ErrorNoRes
	}
	// 多个文件打包返回
	zipWriter := zip.NewWriter(ctx.Response().BodyWriter())
	for _, file := range files {
		f, err := os.Open(ChangeToSysPath(file))
		if err != nil {
			continue
		}
		fi, _ := f.Stat()
		zf, _ := zipWriter.CreateHeader(&zip.FileHeader{
			Name:     file,
			Modified: fi.ModTime(),
			Method:   zip.Deflate,
			Flags:    0x800,
		})
		io.Copy(zf, f)
		f.Close()
	}
	zipWriter.Close()
	return nil, api.ErrorNoRes
}
