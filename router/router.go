// router
package router

import (
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
	router *Router
)

func init() {
	router = new(Router)
	router.Tree = new(Tree)
	router.Pool.New = func() interface{} {
		return context.NewContext()
	}
}

func Route(path string, controller controller.ControllerInterface) (leavesSlice []*Leaves) {
	leaves := new(Leaves)
	leaves.controllerRunObjects = reflect.Indirect(reflect.ValueOf(controller)).Type()
	leaves.path = path
	leavesSlice = []*Leaves{leaves}

	return
}

func RouteNS(path string, leavesSlices ...[]*Leaves) (resLeavesSlice []*Leaves) {
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
			//						logs.Infof("Mappting [%s] to [%s]", leaves.path, leaves.controllerRunObjects)
			leavesFinal.InsertLeaves(leaves)
		}
	}
	leavesFinal.split()
	router.Tree.root = leavesFinal
	//	leavesFinal.walkCheck()
	return router
}

func GetRouter() *Router {
	return router

}

var path = ""
var i = 0

func CheckRouter() {
	router.Tree.root.walkCheck()
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
		filters := ""
		for k := 0; k < len(this.filterRunObjects); k++ {
			fn := strings.Split(runtime.FuncForPC(reflect.ValueOf(this.filterRunObjects[k].filterFunc).Pointer()).Name(), ".")[1]
			filters += fn + " "
		}
		logs.Info("Mapping [", path, "] to controller: ", this.controllerRunObjects, ", filters: [", filters, "]") // , this.filterRunObjects)

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

func (this *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	context := this.Pool.Get().(*context.Context)
	context.Init(w, r)
	defer this.Pool.Put(context)

	leaves := this.Tree.root.GetLeavesByPath(r.URL.Path, context)
	if leaves == nil {
		w.WriteHeader(404)
		return
	}
	if leaves.controllerRunObjects == nil {
		w.WriteHeader(404)
		return
	}

	if c, ok := getControllerInterface(leaves.controllerRunObjects, context); ok {
		if leaves.filterRunObjects != nil {
			leaves.execFilter(context)
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
	w.WriteHeader(404)
	return
}
