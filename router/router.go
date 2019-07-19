// router
package router

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/whitheyxu/snow/context"
	"github.com/whitheyxu/snow/controller"
	"github.com/whitheyxu/snow/g/logs"
)

type Router struct {
	Tree *Tree
	Pool sync.Pool
}

var (
	router       *Router
	staticRouter *map[string]string
)

func init() {
	router = new(Router)
	router.Tree = new(Tree)
	router.Pool.New = func() interface{} {
		return context.NewContext()
	}
	sr := make(map[string]string)
	staticRouter = &sr
}

func Route(path string, controller controller.ControllerInterface) (leavesSlice []*Leaves) {
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	leaves := new(Leaves)
	leaves.controllerRunObjects = reflect.Indirect(reflect.ValueOf(controller)).Type()
	leaves.path = path
	leavesSlice = []*Leaves{leaves}

	return
}

func RouteNS(path string, leavesSlices ...[]*Leaves) (resLeavesSlice []*Leaves) {
	if path != "/" && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	for _, leavesSlice := range leavesSlices {
		for _, leaves := range leavesSlice {
			leaves.path = path + leaves.path
			resLeavesSlice = append(resLeavesSlice, leaves)
		}
	}
	return
}

func NewRouter(leavesSlices ...[]*Leaves) *Router {
	leavesFinal := new(Leaves)
	for _, leavesSlice := range leavesSlices {
		for _, leaves := range leavesSlice {
			leavesFinal.InsertLeaves(leaves)
		}
	}
	leavesFinal.split()
	router.Tree.root = leavesFinal
	return router
}

func GetRouter() *Router {
	return router
}

var path = ""
var i = 0

func CheckRouter() {
	router.Tree.root.walkCheck()
	walkCheckStaticRouter()
}

func walkCheckStaticRouter() {
	for path, dir := range *staticRouter {
		logs.Infof("MappingStatic [%s] to dir: [%s]\n", path, dir)
	}
	return
}

func (this *Leaves) walkCheck() {
	if strings.Contains(this.index, ":") && strings.Contains(this.index, "*") {
		logs.Panic(`router subPath '*' conflict with ':'`)
	}
	sub := ""
	path += this.path
	for j := 0; j < i; j++ {
		sub = sub + "-"
	}

	if this.controllerRunObjects != nil {
		controller := strings.Split(fmt.Sprintf("%s", this.controllerRunObjects), ".")[1]
		filters := ""
		for k := 0; k < len(this.filterRunObjects); k++ {
			fn := strings.Split(runtime.FuncForPC(reflect.ValueOf(this.filterRunObjects[k].filterFunc).Pointer()).Name(), ".")[1]
			filters += fn + " "
		}
		//		logs.Info("Mapping [", path, "] to controller: ", this.controllerRunObjects, ", filters: [", filters, "]") // , this.filterRunObjects)
		logs.Infof("Mapping [%s] to controller: [%s] and filters: [%s]\n", path, controller, filters)

	}
	i++
	for _, leaves := range this.children {
		leaves.walkCheck()
	}
	i--
	path = strings.TrimSuffix(path, this.path)
}

func getControllerInterface(c interface{}, ctx *context.Context) (ci controller.ControllerInterface, ok bool) {
	if c == nil {
		return nil, false
	}
	vc := reflect.New(c.(reflect.Type))
	ci, ok = vc.Interface().(controller.ControllerInterface)
	if !ok {
		return nil, false
	}
	ci.Init(ctx)
	return

}

func RouteStatic(path string, dir string) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if path == "/" {
		logs.Panic("illegal statuc router path: ", path)
	}
	if strings.HasSuffix(dir, "/") {
		dir = strings.TrimSuffix(dir, "/")
	}

	(*staticRouter)[path] = dir
	return
}

func serveFile(ctx *context.Context, path string, dir string) {
	filePath := dir + strings.TrimPrefix(ctx.Request.URL.Path, path)
	http.ServeFile(ctx.Response.Writer, ctx.Request, filePath)
	ctx.Response.IsWritten = true
	return

}

func serveFilePre(ctx *context.Context) {
	for path, dir := range *staticRouter {
		if strings.HasPrefix(ctx.Request.URL.Path, path) {
			serveFile(ctx, path, dir)
		}
	}

}

func handleNotFound(ctx *context.Context) {
	ctx.Response.Writer.WriteHeader(404)
	return
}

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx := this.Pool.Get().(*context.Context)
	ctx.Init(w, r)
	defer this.Pool.Put(ctx)

	serveFilePre(ctx)
	if ctx.Response.IsWritten {
		return
	}

	leaves := this.Tree.root.GetLeavesByPath(r.URL.Path, ctx)
	if leaves == nil {
		handleNotFound(ctx)
		return
	}
	if leaves.controllerRunObjects == nil {
		handleNotFound(ctx)
		return
	}

	if c, ok := getControllerInterface(leaves.controllerRunObjects, ctx); ok {
		if leaves.filterRunObjects != nil {
			leaves.execFilter(ctx)
			if ctx.Response.IsWritten {
				return
			}
		}

		switch r.Method {
		case http.MethodGet:
			c.Get()
		case http.MethodPost:
			c.Post()
		case http.MethodDelete:
			c.Delete()
		case http.MethodPut:
			c.Put()
		case http.MethodOptions:
			c.Options()
		}

		return
	}
	handleNotFound(ctx)
	return
}
