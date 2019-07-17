// controller
package controller

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

	"github.com/whitheyxu/snow/context"
	"github.com/whitheyxu/snow/g"
)

const (
	uriParamIndex          = ":"
	paramParseFailedFormat = "Param \"%s\" parse failed. "
)

type ControllerInterface interface {
	Get()
	Post()
	Put()
	Delete()
	Options()
	Patch()
	Head()
	Init(context *context.Context)
}

type Controller struct {
	Ctx *context.Context
}

// Get String param uri / url

func (this *Controller) GetStringForce(key string) (value string) {
	if strings.Index(key, uriParamIndex) == 0 {
		value, _ = this.getStringFromUriParams(key)
		return
	}
	value, _ = this.getStringFromUrlParams(key)
	return
}

func (this *Controller) GetString(key string) (string, error) {
	if strings.Index(key, uriParamIndex) == 0 {
		return this.getStringFromUriParams(key)
	}
	return this.getStringFromUrlParams(key)
}

func (this *Controller) getStringFromUrlParams(key string) (string, error) {
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetEmptyString(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return arr[0], nil
}

func (this *Controller) getStringFromUriParams(key string) (string, error) {
	value, ok := this.Ctx.UriParams[key]
	if !ok {
		return g.GetEmptyString(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return value, nil
}

// Get Int param from uri / url ( get, post, put )

func (this *Controller) GetIntForce(key string) (value int) {
	if strings.Index(key, uriParamIndex) == 0 {
		value, _ = this.getIntFromUriParams(key)
		return
	}
	value, _ = this.getIntFromUrlParams(key)
	return
}

func (this *Controller) GetInt(key string) (int, error) {
	if strings.Index(key, uriParamIndex) == 0 {
		return this.getIntFromUriParams(key)
	}
	return this.getIntFromUrlParams(key)
}

func (this *Controller) getIntFromUrlParams(key string) (int, error) {

	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetZeroInt(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return g.StringToInt(arr[0])
}

func (this *Controller) getIntFromUriParams(key string) (int, error) {
	value, ok := this.Ctx.UriParams[key]
	if !ok {
		return g.GetZeroInt(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return g.StringToInt(value)
}

// Get Int64 param from uri / url ( get ,post, put )

func (this *Controller) GetInt64Force(key string) (value int64) {
	if strings.Index(key, uriParamIndex) == 0 {
		value, _ = this.getInt64FromUriParams(key)
		return
	}
	value, _ = this.getInt64FromUrlParams(key)
	return
}

func (this *Controller) GetInt64(key string) (int64, error) {
	if strings.Index(key, uriParamIndex) == 0 {
		return this.getInt64FromUriParams(key)
	}
	return this.getInt64FromUrlParams(key)
}

func (this *Controller) getInt64FromUrlParams(key string) (int64, error) {

	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetZeroInt64(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return g.StringToInt64(arr[0])
}

func (this *Controller) getInt64FromUriParams(key string) (int64, error) {
	value, ok := this.Ctx.UriParams[key]
	if !ok {
		return g.GetZeroInt64(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return g.StringToInt64(value)
}

// Get String slice param from param [get, put, post]

func (this *Controller) GetStringSlice(key string) ([]string, error) {
	return this.getStringSliceFromUrlParams(key)
}

func (this *Controller) GetStringSliceForce(key string) (value []string) {
	value, _ = this.getStringSliceFromUrlParams(key)
	return value
}

func (this *Controller) getStringSliceFromUrlParams(key string) ([]string, error) {
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetEmptyStringSlice(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return []string(arr), nil
}

// Get Int slice param from param [get, put, post]

func (this *Controller) GetIntSlice(key string) ([]int, error) {
	return this.getIntSliceFromUrlParams(key)
}

func (this *Controller) GetIntSliceForce(key string) (value []int) {
	value, _ = this.getIntSliceFromUrlParams(key)
	return
}

func (this *Controller) getIntSliceFromUrlParams(key string) ([]int, error) {
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetEmptyIntSlice(), fmt.Errorf(paramParseFailedFormat, key)
	}

	return g.ArrStringsToArrInt(arr)
}

// Get Int64 slice param from param [get, put, post]

func (this *Controller) GetInt64Slice(key string) ([]int64, error) {
	return this.getInt64SliceFromUrlParams(key)
}

func (this *Controller) GetInt64SliceForce(key string) (value []int64) {
	value, _ = this.getInt64SliceFromUrlParams(key)
	return
}

func (this *Controller) getInt64SliceFromUrlParams(key string) ([]int64, error) {
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetEmptyInt64Slice(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return g.ArrStringsToArrInt64(arr)
}

// Get Bool param from param [get, put, post]

func (this *Controller) GetBool(key string) (bool, error) {
	return this.getBoolFromUrlParams(key)
}

func (this *Controller) GetBoolForce(key string) (value bool) {
	value, _ = this.getBoolFromUrlParams(key)
	return
}

func (this *Controller) getBoolFromUrlParams(key string) (bool, error) {
	convertErrMsg := "convert to bool type failed"
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetDefaultFalseBool(), fmt.Errorf(paramParseFailedFormat, key)
	}
	if arr[0] == "True" || arr[0] == "true" || arr[0] == "1" {
		return true, nil
	} else if arr[0] == "False" || arr[0] == "false" || arr[0] == "0" {
		return false, nil
	} else {
		return false, fmt.Errorf(convertErrMsg)
	}

}

// Get Float64 param from param [get, put, post]

func (this *Controller) GetFloat64(key string) (float64, error) {
	return this.getFloat64fromUrlParams(key)
}

func (this *Controller) GetFloat64Force(key string) (value float64) {
	value, _ = this.getFloat64fromUrlParams(key)
	return
}

func (this *Controller) getFloat64fromUrlParams(key string) (float64, error) {
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetZeroFloat64(), fmt.Errorf(paramParseFailedFormat, key)
	}

	return g.StringToFloat64(arr[0])
}

// Get Float64 slice from param [get, put, post]

func (this *Controller) GetFloat64Slice(key string) ([]float64, error) {
	return this.getFloat64SlicefromUrlParams(key)
}

func (this *Controller) GetFloat64SliceForce(key string) (value []float64) {
	value, _ = this.getFloat64SlicefromUrlParams(key)
	return
}

func (this *Controller) getFloat64SlicefromUrlParams(key string) ([]float64, error) {
	arr, ok := this.Ctx.Request.Form[key]
	if !ok || len(arr) < 1 {
		return g.GetEmptyFloat64Slice(), fmt.Errorf(paramParseFailedFormat, key)
	}
	return g.ArrStringsToArrFloat64(arr)
}

// Get File from form data
func (this *Controller) GetFile(key string) (file multipart.File, header *multipart.FileHeader, err error) {
	file, header, err = this.Ctx.Request.FormFile(key)
	return
}

func (this *Controller) SaveToFile(file multipart.File, dstFileName string, dstFilePath string) (err error) {
	f, err := os.OpenFile(dstFilePath+dstFileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	_, err = io.Copy(f, file)
	return
}

// handle 404 status

// implement controllerInterface
func (this *Controller) Get() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}
func (this *Controller) Post() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}
func (this *Controller) Put() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}
func (this *Controller) Delete() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}
func (this *Controller) Options() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}

func (this *Controller) Patch() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}

func (this *Controller) Head() {
	this.Ctx.ResponseWriter.WriteHeader(405)
	return
}

func (this *Controller) Response(resp []byte) {
	this.Ctx.ResponseWriter.Write(resp)
	return
}

func (this *Controller) Init(context *context.Context) {
	this.Ctx = context
	return
}
