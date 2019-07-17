// context
package context

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	UriParams   map[string]string
	RequestBody []byte
}

func (this *Context) Init(w http.ResponseWriter, r *http.Request) {

	this.Request = r
	this.ResponseWriter = w

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
