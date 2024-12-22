package server

import (
	"fmt"
	"local-cloud-api/conf"
	"local-cloud-api/handle"
	"local-cloud-api/static"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	fstatic "github.com/gofiber/fiber/v3/middleware/static"
)

func Run() {
	// 初始化应用
	app := fiber.New()
	// 绑定API
	for route, apifunc := range handle.ApiRouteInfo {
		switch apifunc.GetMethod() {
		case http.MethodGet:
			app.Get(conf.ApiPrefix+route, apifunc.RouteFunc())
		case http.MethodDelete:
			app.Delete(conf.ApiPrefix+route, apifunc.RouteFunc())
		case http.MethodPatch:
			app.Patch(conf.ApiPrefix+route, apifunc.RouteFunc())
		case http.MethodPost:
			app.Post(conf.ApiPrefix+route, apifunc.RouteFunc())
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

	// 启动服务
	err := app.Listen(conf.Bind + ":" + fmt.Sprint(conf.Port))
	if err != nil {
		log.Fatal(err)
	}
}
