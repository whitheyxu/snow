// filter
package filter

import (
	"github.com/whitheyxu/snow/context"
)

type FilterInterface interface {
	Filter(ctx context.Context)
}
