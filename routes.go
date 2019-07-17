// routes
package snow

import (
	"snow/context"
	"snow/controller"
	"snow/router"
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

func NewRouter(leavesSlices ...[]*router.Leaves) *router.Router {
	return router.NewRouter(leavesSlices...)
}
