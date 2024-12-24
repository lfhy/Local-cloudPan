package api

import (
	"fmt"
	"reflect"
	"strconv"

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
		c.JSON(&req)
		ParseQuery(&req, c)
		res, err := r.Handle(c, &req)
		if err != nil {
			return ResError(c, err)
		}
		return ResOK(c, "", res)
	}
}

func (r *RouteController[Req, Res]) GetMethod() string {
	return r.Method
}

func ParseQuery(data any, ctx fiber.Ctx) {
	defer recover()
	// query := ctx.Queries()
	// 解析结构体
	t := reflect.ValueOf(data)
	if t.Kind() == reflect.Ptr {
		if t.IsNil() {
			t = reflect.New(t.Type())
		}
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			dtype := t.Type()
			value := dtype.Field(i)
			if !value.IsExported() {
				continue
			}
			v, ok := value.Tag.Lookup("query")
			if ok {
				// 获取到定义的值
				val := reflect.ValueOf(data).Elem().Field(i)
				defaultValue, _ := value.Tag.Lookup("default")
				setData := ctx.Query(v, defaultValue)
				if setData != "" {
					// jtool.Debug(setData)
					ParseStringToAny(setData, val)
				}
			}
		}
	default:
		log.Warn("不支持的类型:", t.Kind())
	}
}

func ParseStringToAny(value string, val any) error {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if rv, ok := val.(reflect.Value); ok {
		reflectValue = rv
	} else {
		reflectValue = reflect.ValueOf(val)
	}
	// 取出真实类型

	for {
		reflectKind = reflectValue.Kind()
		switch reflectKind {
		case reflect.Ptr:
			if !reflectValue.IsValid() || reflectValue.IsNil() {
				// 为空就创一个默认值出来
				reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			} else {
				reflectValue = reflectValue.Elem()
			}
		case reflect.Int:
			data, err := strconv.Atoi(value)
			if err != nil {
				return err
			} else {
				reflectValue.SetInt(int64(data))
				return nil
			}
		case reflect.String:
			reflectValue.SetString(value)
			return nil
		case reflect.Bool:
			data, err := strconv.ParseBool(value)
			if err != nil {
				return err
			} else {
				reflectValue.SetBool(data)
				return nil
			}
		case reflect.Float64:
			data, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return err
			} else {
				reflectValue.SetFloat(data)
				return nil
			}
		case reflect.Uint:
			data, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return err
			} else {
				reflectValue.SetUint(data)
				return nil
			}
		default:
			// 不支持的类型
			return fmt.Errorf("不支持的类型: %v", reflectKind)
		}
	}
}
