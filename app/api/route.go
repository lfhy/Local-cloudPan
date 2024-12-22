package api

import "github.com/gofiber/fiber/v3"

type RouteController[Req any, Res any] struct {
	Method string
	Handle func(ctx fiber.Ctx, req Req) (res Res, err error)
}

type RouteInfo interface {
	RouteFunc() func(fiber.Ctx) error
	GetMethod() string
}

func (r *RouteController[Req, Res]) RouteFunc() func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		// TODO:获取请求参数
		var req Req
		c.JSON(&req)
		res, err := r.Handle(c, req)
		if err != nil {
			return ResError(c, err)
		}
		return ResOK(c, "", res)
	}
}

func (r *RouteController[Req, Res]) GetMethod() string {
	return r.Method
}
