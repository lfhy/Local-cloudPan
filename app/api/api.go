package api

import "github.com/gofiber/fiber/v3"

type ApiRes struct {
	c        fiber.Ctx
	httpCode int
	Code     int    `json:"code"`
	Msg      string `json:"msg,omitempty"`
	Data     any    `json:"data,omitempty"`
}

func (res *ApiRes) Send() error {
	res.c.Status(res.httpCode)
	return res.c.JSON(res)
}

func ResOK(ctx fiber.Ctx, msg string, data any) error {
	res := ApiRes{c: ctx, Code: 200, httpCode: 200, Msg: msg, Data: data}
	return res.Send()
}

func ResError(ctx fiber.Ctx, err error) error {
	res := ApiRes{c: ctx, Msg: err.Error()}
	apiError, ok := err.(Error)
	if ok {
		if apiError.code == 0 {
			return nil
		}
		res.httpCode = 200
		res.Code = apiError.code
	} else {
		res.httpCode = 500
		res.Code = 500
	}
	return res.Send()
}
