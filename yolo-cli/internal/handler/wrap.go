package handler

import (
	"context"
	"net/http"
	"reflect"

	"yolo-team/yolo-cli/internal/common"

	"github.com/gin-gonic/gin"
)

func Wrap[R any, T any](fn func(ctx context.Context, req R) (T, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req R

		if err := c.ShouldBindUri(&req); err != nil {
			common.Error(c, 40001, err.Error())
			return
		}

		switch c.Request.Method {
		case http.MethodGet, http.MethodDelete:
			if err := c.ShouldBindQuery(&req); err != nil {
				common.Error(c, 40001, err.Error())
				return
			}
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := c.ContentType()
			if ct == "application/json" || ct == "" {
				if err := c.ShouldBindJSON(&req); err != nil {
					common.Error(c, 40001, err.Error())
					return
				}
			}
		}

		bindHeaders(c, &req)

		resp, err := fn(c.Request.Context(), req)
		if err != nil {
			if appErr, ok := common.AsAppError(err); ok {
				common.Error(c, appErr.Code, appErr.Message)
			} else {
				common.Error(c, 50001, err.Error())
			}

			return
		}

		common.Success(c, resp)
	}
}

func bindHeaders(c *gin.Context, req interface{}) {
	v := reflect.ValueOf(req)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}

	v = v.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("header")
		if tag == "" {
			continue
		}

		val := c.GetHeader(tag)
		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			fv.SetString(val)
		case reflect.Bool:
			fv.SetBool(val == "true")
		}
	}
}
