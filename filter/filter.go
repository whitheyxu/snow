// filter
package filter

import (
	"snow/context"
)

type FilterInterface interface {
	Filter(ctx context.Context)
}
