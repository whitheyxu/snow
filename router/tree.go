// tree
package router

import (
	"strings"

	"github.com/whitheyxu/snow/context"
	"github.com/whitheyxu/snow/g/logs"
)

type sortableInterface interface {
	getPriority() int
}

type Leaves struct {
	children             []*Leaves
	controllerRunObjects interface{}
	filterRunObjects     []*Filter
	path                 string
	index                string
}

type Tree struct {
	root *Leaves
}

func (this *Leaves) split() {
	sIndex := strings.Index(this.path, "/")
	if sIndex >= 0 && sIndex < len(this.path)-1 {
		leaves := newLeaves(this.path[sIndex+1:], this.controllerRunObjects)
		leaves.children = this.children
		leaves.index = this.index

		this.index = string(this.path[sIndex+1])
		this.path = this.path[:sIndex+1]
		this.controllerRunObjects = nil
		this.children = nil
		this.children = []*Leaves{leaves}
		leaves.split()
	}
	for _, leaves := range this.children {
		leaves.split()
	}
	return
}

func (this *Leaves) sortAppendFilterRunObject(filter *Filter) {
	for i := 0; i < len(this.filterRunObjects); i++ {
		if filter.priority == this.filterRunObjects[i].priority {
			logs.Panic(`filter priority conflict`)
			return
		} else if filter.priority > this.filterRunObjects[i].priority {
			continue
		} else {
			this.filterRunObjects = append(this.filterRunObjects[:i+1], this.filterRunObjects[i:]...)
			this.filterRunObjects[i] = filter

			return
		}
	}
	this.filterRunObjects = append(this.filterRunObjects, filter)
	return
}

func newLeaves(path string, c interface{}) (leaves *Leaves) {
	leaves = new(Leaves)
	leaves.path = path
	leaves.controllerRunObjects = c
	return
}

func getMaxPrefix(c1 string, c2 string) (prefix string) {
	if len(c1) > len(c2) {
		c1, c2 = c2, c1
	}
	for i := 0; i <= len(c1); i++ {
		if c1[:i] == c2[:i] {
			prefix = c1[:i]
		}
	}
	return
}

func (this *Leaves) isEmptyLeaves() (isEmpty bool) {
	if this.path != "" {
		isEmpty = false
		return
	}
	if this.index != "" {
		isEmpty = false
		return
	}
	if this.controllerRunObjects != nil {
		isEmpty = false
		return
	}
	isEmpty = true

	return
}

func pathMatch(path string, subPath string) (isPrefixMatch bool, isMatch bool) {

	isPrefixMatch = false
	isMatch = false
	if subPath[0] == '*' || subPath[0] == ':' {
		isPrefixMatch = true
		if strings.Index(path, "/") < 0 || strings.Index(path, "/") == len(path)-1 {
			isMatch = true
		}
		return
	}
	if path == subPath || path == subPath+"/" || subPath == path+"/" {
		isMatch = true
	}
	if strings.HasPrefix(path, subPath) || strings.HasPrefix(subPath, path) {
		isPrefixMatch = true
	}
	return
}

func getWildcardKey(path string) (key string) {
	if strings.HasSuffix(path, "/") {
		key = path[:len(path)-1]
	} else {
		key = path
	}
	return
}
func getWildcardValue(path string) (value string) {
	index := strings.Index(path, "/")
	if index < 0 {
		value = path
	} else {
		value = path[:index]
	}
	return
}

func (this *Leaves) execFilter(ctx *context.Context) {
	for i := 0; i < len(this.filterRunObjects); i++ {
		this.filterRunObjects[i].filter(ctx)
	}
	return

}

func (this *Leaves) passonFiltersToChildren(filter *Filter) {
	for _, leaves := range this.children {
		if leaves.controllerRunObjects != nil {
			leaves.sortAppendFilterRunObject(filter)
		}
		if leaves.children != nil {
			leaves.passonFiltersToChildren(filter)
		}
	}
	return
}

func (this *Router) InsertFilters(path string, filter *Filter) {
	this.Tree.root.insertFilters(path, filter)
	return
}

func (this *Leaves) insertFilters(path string, filter *Filter) {

	var isWildcard bool
	if strings.Count(path, "*") == 0 {
		isWildcard = false
	} else if strings.Count(path, "*") == 1 && strings.HasSuffix(path, "*") {
		isWildcard = true
		path = path[:len(path)-1]
	} else {
		logs.Panic("Do not support the path :" + path)
	}
	leaves := this.GetLeavesByPath(path, nil)
	if leaves.controllerRunObjects != nil {
		leaves.sortAppendFilterRunObject(filter)
		if isWildcard == true {
			leaves.passonFiltersToChildren(filter)
		}

	} else if isWildcard == true {
		leaves.passonFiltersToChildren(filter)
	}
	return

}

func (this *Leaves) InsertLeaves(leaves *Leaves) {

	// do not support insert a leaf with children
	if leaves.children != nil {
		logs.Panic("Do not support insert a parent leaf node")
		return
	}

	if this.isEmptyLeaves() {
		this.path = leaves.path
		this.controllerRunObjects = leaves.controllerRunObjects
		return
	}

	prefix := getMaxPrefix(this.path, leaves.path)

	l1 := newLeaves(string(this.path[len(prefix):]), this.controllerRunObjects)
	l2 := newLeaves(string(leaves.path[len(prefix):]), leaves.controllerRunObjects)

	if l1.path == "" && l2.path == "" {
		this.controllerRunObjects = l2.controllerRunObjects
		return
	}

	this.path = prefix

	if l1.path == "" {
		isMatchChildPrefix := false
		for _, childLeaves := range this.children {
			if prefix := getMaxPrefix(childLeaves.path, l2.path); prefix != "" {
				isMatchChildPrefix = true
				childLeaves.InsertLeaves(l2)
			}
		}
		if !isMatchChildPrefix {
			this.index = this.index + string(l2.path[0])
			this.children = append(this.children, l2)
		}
		return
	}
	if l2.path == "" {
		this.controllerRunObjects = leaves.controllerRunObjects
		l1.children = this.children
		l1.index = this.index
		this.children = []*Leaves{l1}
		this.index = string(l1.path[0])
		return
	}
	if l1.path != "" && l2.path != "" {
		this.controllerRunObjects = nil
		l1.children = this.children
		l1.index = this.index
		this.children = nil
		this.children = append(this.children, l1, l2)
		this.index = string(l1.path[0]) + string(l2.path[0])
		return
	}

	return
}

func (this *Leaves) GetChildLeavesByIndex(index byte) (leaves *Leaves) {
	var leavesWild *Leaves
	for i := 0; i < len(this.index); i++ {
		if '*' == this.index[i] || ':' == this.index[i] {
			if leavesWild != nil {
				logs.Panic(`router conflict .  ':' conflict widh '*' `)
			}
			leavesWild = this.children[i]
		}
		if this.index[i] == index {
			leaves = this.children[i]
			return
		}
	}
	if leavesWild != nil {
		leaves = leavesWild
		return
	}
	return

}

func (this *Leaves) GetLeavesByPath(path string, ctx *context.Context) (leaves *Leaves) {
	isPrefixMatch, isMatch := pathMatch(path, this.path)
	if isMatch {
		return this
	} else {
		if isPrefixMatch {
			if this.path[0] == ':' {
				key := getWildcardKey(this.path)
				value := getWildcardValue(path)
				ctx.UriParams[key] = value
			}

			var subIndex byte
			var subPath string
			if this.path[0] == ':' || this.path[0] == '*' {
				subIndex = path[strings.Index(path, "/")+1]
				subPath = path[strings.Index(path, "/")+1:]
			} else {
				subIndex = path[len(this.path)]
				subPath = path[len(this.path):]
			}

			leaves := this.GetChildLeavesByIndex(subIndex)
			if leaves != nil {
				subLeaves := leaves.GetLeavesByPath(subPath, ctx)
				if subLeaves == nil {
					subLeaves = this.GetChildLeavesByIndex('*')
					return subLeaves
				}
				return subLeaves
			}
			return nil
		} else {
			return nil
		}

	}
	return
}
