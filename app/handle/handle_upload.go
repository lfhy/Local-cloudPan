package handle

import (
	"fmt"
	"io"
	"local-cloud-api/api"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

func init() {
	RegisterHandle("/upload", http.MethodPost, upload)
}

// 上传文件分块
func upload(ctx fiber.Ctx, req *api.ApiUploadReq) (*api.ApiUploadRes, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}
	chunkData := form.File["chunkData"]
	if len(chunkData) == 0 {
		return nil, fmt.Errorf("缺少分块数据")
	}
	file, err := chunkData[0].Open() // 读取文件内容
	if err != nil {
		log.Warnln("读取数据错误:", err)
		return nil, err
	}
	defer file.Close()
	log.Debugf("上传文件分块: %+v\n", req)
	// 保存文件
	workDir := GetFileUploadPath(req.FileId)
	os.MkdirAll(workDir, os.ModeDir)
	chunkPath := filepath.Join(workDir, req.ChunkId)
	savePath, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Warnln("保存文件错误:", err)
		return nil, err
	}
	defer savePath.Close()
	_, err = io.Copy(savePath, file) // 将文件内容写入到指定路径
	if err != nil {
		log.Warnln("复制数据错误:", err)
		return nil, err
	}
	SetResMsg(ctx, "分片上传成功")
	return nil, nil
}
