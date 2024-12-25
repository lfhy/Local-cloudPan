package api

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/lfhy/log"
)

type ApiRes struct {
	c        fiber.Ctx
	httpCode int
	Code     int    `json:"code"`
	Msg      string `json:"msg,omitempty"`
	Data     any    `json:"data,omitempty"`
}

func (res *ApiRes) Send() error {
	res.c.Status(res.httpCode)
	err := res.c.JSON(res)
	if err != nil {
		log.Warnln("发送失败:", err)
	}
	return err
}

func ResOK(ctx fiber.Ctx, msg string, data any) error {
	res := ApiRes{c: ctx, Code: 200, httpCode: 200, Msg: msg, Data: data}
	return res.Send()
}

func ResError(ctx fiber.Ctx, err error) error {
	res := ApiRes{c: ctx, Msg: err.Error()}
	userRes := ctx.UserContext().Value(ResApiMsg)
	if userRes != nil {
		res.Msg = fmt.Sprint(res)
	}
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
	log.PError("返回错误:", err)
	return res.Send()
}
