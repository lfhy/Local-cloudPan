package handle

import (
	"local-cloud-api/api"
	"local-cloud-api/conf"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

var (
	// 正则表达式匹配Markdown中的图片链接
	ImgRe = regexp.MustCompile(`!\[.*?\]\((.*?)\)|<img src="(.*?)"`)
)

func init() {
	RegisterHandle("/view/*", http.MethodGet, view)
}

// 获取文件
func view(ctx fiber.Ctx, req *api.ApiViewReq) (*api.ApiViewRes, error) {
	filePath := ctx.Params("*")
	absPath := ChangeToSysPath(filePath)
	log.Info("filePath:", absPath)
	ext := strings.TrimPrefix(filepath.Ext(absPath), ".")
	ctx.Type(ext).Response().Header.Set("Cache-Control", "max-age=86400")
	if req.Short {
		data, err := GetShortImg(absPath)
		if err == nil {
			ctx.Send(data)
			return nil, api.ErrorNoRes
		}
	}
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	if req.ReplaceImgPath && strings.ToLower(filepath.Ext(absPath)) == ".md" {

		// 存储提取到的图片路径
		var imgPaths []string
		matches := ImgRe.FindAllStringSubmatch(string(data), -1)
		// 提取图片路径
		for _, match := range matches {
			for _, path := range match[1:] {
				if path != "" {
					imgPaths = append(imgPaths, path)
				}
			}
		}

		// 替换路径
		resdata := string(data)
		for _, imgPath := range imgPaths {
			// fmt.Println(imgPath)
			if strings.HasPrefix(imgPath, "http") {
				continue
			}
			resdata = strings.ReplaceAll(resdata, imgPath, conf.ApiPrefix+"/view/"+filepath.Dir(filePath)+"/"+imgPath)
		}
		data = []byte(resdata)
	}
	ctx.Send(data)
	return nil, api.ErrorNoRes
}
