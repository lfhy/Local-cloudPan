package handle

import (
	"encoding/hex"
	"fmt"
	"io"
	"local-cloud-api/api"
	"net/http"
	"os"

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
	os.MkdirAll(workDir, 0755)
	chunkPath := req.GetChunkPath()
	savePath, err := os.OpenFile(chunkPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Warnln("保存文件错误:", err)
		return nil, err
	}
	defer savePath.Close()
	hash, md5Write := Md5Write(savePath)
	_, err = io.Copy(md5Write, file) // 将文件内容写入到指定路径
	if err != nil {
		log.Warnln("复制数据错误:", err)
		return nil, err
	}
	if hex.EncodeToString(hash.Sum(nil)) != req.ChunkId {
		os.RemoveAll(chunkPath)
		return nil, fmt.Errorf("分块数据校验失败")
	}
	log.Debugln(req.ChunkId, "上传结束")
	SetResMsg(ctx, "分片上传成功")
	return nil, nil
	// return nil, api.ErrorNoRes
}
