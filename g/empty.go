// empty
package g

const (
	emptyString      string  = ""
	zeroInt          int     = 0
	zeroInt64        int64   = 0
	zeroFloat64      float64 = 0.0
	defaultFalseBool bool    = false
)

var emptyStruct = struct{}{}
var emptyIntSlice = []int{}
var emptyInt64Slice = []int64{}
var emptyStringSlice = []string{}
var emptyFloat64Slice = []float64{}

func GetZeroInt() int {
	return zeroInt
}

func GetZeroInt64() int64 {
	return zeroInt64
}

func GetDefaultFalseBool() bool {
	return defaultFalseBool
}

func GetZeroFloat64() float64 {
	return zeroFloat64
}

func GetEmptyString() string {
	return emptyString
}

func GetEmptyStruct() struct{} {
	return emptyStruct
}

func GetEmptyIntSlice() []int {
	return emptyIntSlice
}

func GetEmptyInt64Slice() []int64 {
	return emptyInt64Slice
}

func GetEmptyFloat64Slice() []float64 {
	return emptyFloat64Slice
}

func GetEmptyStringSlice() []string {
	return emptyStringSlice
}
