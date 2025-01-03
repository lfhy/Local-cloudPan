package conf

import (
	"os"
	"path/filepath"

	"github.com/lfhy/flag"
)

var (
	Port            string
	ShareFilePath   string
	Bind            string
	ApiPrefix       string
	LocalStaticPath string
	DisableView     bool
	IgnoreDotFiles  bool
	UploadTmpPath   string
	ApiMode         bool
	Gui             bool
)

func init() {
	baseDir, err := os.UserHomeDir()
	if err != nil {
		baseDir = os.TempDir()
	}
	flag.StringConfigVar(&Port, "port", "server", "port", "auto", "服务端口")
	flag.StringConfigVar(&ShareFilePath, "path", "server", "path", baseDir, "分享文件的路径")
	flag.StringConfigVar(&Bind, "bind", "server", "bind", "127.0.0.1", "监听的IP地址")
	flag.StringConfigVar(&ApiPrefix, "api-prefix", "server", "prefix", "/api", "API路由前缀")
	flag.StringConfigVar(&LocalStaticPath, "static", "server", "static", "", "前端编译的静态页面路径")
	flag.BoolConfigVar(&DisableView, "no-view", "server", "noview", false, "不启动静态页面路由")
	flag.BoolConfigVar(&IgnoreDotFiles, "ignore-dot", "server", "ignore-dot", true, "忽略点开头的文件")
	flag.StringConfigVar(&UploadTmpPath, "upload-tmp", "server", "upload-tmp", filepath.Join(baseDir, ".uploads"), "上传临时文件路径")
	flag.BoolVar(&ApiMode, "api-mode", false, "只启动为api模式")
	flag.BoolVar(&Gui, "gui", false, "启动GUI")
	flag.Parse()
}
