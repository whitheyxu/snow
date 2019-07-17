// filter
package router

import (
	"github.com/whitheyxu/snow/context"
)

type filterInterface interface {
	getPriority() int
	filter(context.Context)
}

type Filter struct {
	priority   int
	filterFunc func(*context.Context)
}

func (this *Filter) filter(ctx *context.Context) {
	this.filterFunc(ctx)
	return
}

func NewFilter(priority int, filterFunc func(*context.Context)) (filter *Filter) {
	filter = new(Filter)
	filter.priority = priority
	filter.filterFunc = filterFunc
	return
}

func RegisterFilter(path string, priority int, filterFunc func(*context.Context)) {

	filter := NewFilter(priority, filterFunc)
	GetRouter().Tree.root.insertFilters(path, filter)
	return
}
