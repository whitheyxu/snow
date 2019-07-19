// routes
package snow

import (
	"github.com/whitheyxu/snow/context"
	"github.com/whitheyxu/snow/controller"
	"github.com/whitheyxu/snow/router"
)

func Route(path string, controller controller.ControllerInterface) (leavesSlice []*router.Leaves) {
	return router.Route(path, controller)
}

func RouteNS(path string, leavesSlices ...[]*router.Leaves) (resLeavesSlice []*router.Leaves) {
	return router.RouteNS(path, leavesSlices...)
}

func RegisterFilter(path string, priority int, filterFunc func(*context.Context)) {

	filter := router.NewFilter(priority, filterFunc)
	router.GetRouter().InsertFilters(path, filter)
	return
}

func RouteStatic(path string, dir string) {
	router.RouteStatic(path, dir)
	return
}

func NewRouter(leavesSlices ...[]*router.Leaves) *router.Router {
	return router.NewRouter(leavesSlices...)
}
