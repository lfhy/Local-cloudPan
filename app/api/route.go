package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

type RouteController[Req any, Res any] struct {
	Method string
	Name   string
	Handle func(ctx fiber.Ctx, req *Req) (res Res, err error)
}

type RouteInfo interface {
	RouteFunc() func(fiber.Ctx) error
	GetMethod() string
}

func (r *RouteController[Req, Res]) RouteFunc() func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		// TODO:获取请求参数
		var req Req
		log.Debug("请求:", r.Name)
		bind := c.Bind()
		bind.JSON(&req)
		bind.Query(&req)
		bind.Header(&req)
		bind.Form(&req)
		bind.MultipartForm(&req)
		res, err := r.Handle(c, &req)
		if err != nil {
			return ResError(c, err)
		}
		msg := c.UserContext().Value(ResApiMsg)
		if msg != nil {
			return ResOK(c, fmt.Sprint(msg), res)
		}
		return ResOK(c, "", res)
	}
}

type ApiMsg string

var ResApiMsg ApiMsg

func (r *RouteController[Req, Res]) GetMethod() string {
	return r.Method
}
