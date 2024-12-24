package handle

import (
	"local-cloud-api/api"

	"github.com/gofiber/fiber/v3"
)

var ApiRouteInfo = make(map[string]api.RouteInfo)

func RegisterHandle[Req any, Res any](route string, method string, fn func(ctx fiber.Ctx, req *Req) (res Res, err error)) {
	ApiRouteInfo[route] = &api.RouteController[Req, Res]{
		Method: method,
		Handle: fn,
		Name:   route,
	}
}
