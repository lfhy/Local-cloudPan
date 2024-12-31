package server

import (
	"fmt"
	"local-cloud-api/conf"
	"local-cloud-api/handle"
	"local-cloud-api/static"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fstatic "github.com/gofiber/fiber/v3/middleware/static"
)

func Run() {
	// 初始化应用

	var config fiber.Config
	config.BodyLimit = 1024 * 1024 * 1024 * 1024
	config.AppName = "local-cloud-api"
	app := fiber.New(config)
	app.Use(cors.New())
	prefix := conf.ApiPrefix
	if conf.ApiMode {
		prefix = ""
		conf.DisableView = true
	}
	// 绑定API
	for route, apifunc := range handle.ApiRouteInfo {
		switch apifunc.GetMethod() {
		case http.MethodGet:
			app.Get(prefix+route, apifunc.RouteFunc())
		case http.MethodDelete:
			app.Delete(prefix+route, apifunc.RouteFunc())
		case http.MethodPatch:
			app.Patch(prefix+route, apifunc.RouteFunc())
		case http.MethodPost:
			app.Post(prefix+route, apifunc.RouteFunc())
		}
	}
	// 初始化页面
	if !conf.DisableView {
		if conf.LocalStaticPath != "" {
			app.Get("*", fstatic.New(conf.LocalStaticPath))
		} else {
			app.Get("*", fstatic.New("dist", fstatic.Config{FS: static.GetFs()}))
		}
		app.Get("/", func(c fiber.Ctx) error {
			return c.Redirect().To("/index.html")
		})
	}
	appConfig := fiber.ListenConfig{
		EnablePrefork:     true,
		EnablePrintRoutes: true,
	}
	// 启动服务
	go app.Listen(conf.Bind+":"+fmt.Sprint(conf.Port), appConfig)
}
