// context
package context

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/whitheyxu/snow/g/logs"
)

type Context struct {
	Response Response
	Request  *http.Request

	UriParams   map[string]string
	RequestBody []byte
}

type Response struct {
	Writer    http.ResponseWriter
	IsWritten bool
}

func (this *Context) Init(w http.ResponseWriter, r *http.Request) {

	this.Request = r
	this.Response.Writer = w
	this.Response.IsWritten = false

	if strings.Index(this.Request.Header.Get("Content-Type"), "multipart/form-data") == 0 {
		this.Request.ParseMultipartForm(1 << 32)
	} else {
		this.Request.ParseForm()
	}
	if r.Body != nil {
		this.RequestBody, _ = ioutil.ReadAll(r.Body)
	} else {
		this.RequestBody = nil
	}
	this.UriParams = make(map[string]string)
}

func NewContext() *Context {
	return new(Context)
}

func (this *Response) Response(resp []byte) {
	this.IsWritten = true

	this.Writer.Write(resp)
	return
}

func (this *Response) ResponseJson(v interface{}) {
	this.IsWritten = true

	response, err := json.Marshal(v)
	logs.Info(string(response))
	if err != nil {
		logs.Error(err)
	}
	this.Writer.Write(response)
	return
}

func (this *Response) Abort(code int) {
	this.IsWritten = true
	this.Writer.WriteHeader(code)
	return
}
