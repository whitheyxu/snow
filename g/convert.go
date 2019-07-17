// convet.go
package g

import (
	"errors"
	"fmt"
	"strconv"
)

func ArrIntToString(arr []int) (result string, err error) {

	result = ""

	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}

	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
		return
	}

	return

}

func ArrIntToStringForce(arr []int) (result string) {

	result, _ = ArrIntToString(arr)

	return

}

func ArrInt64ToString(arr []int64) (result string, err error) {

	result = ""

	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("%v", a)
		} else {
			result = fmt.Sprintf("%v,%v", result, a)
		}
	}

	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}

	return

}

func ArrInt64ToStringForce(arr []int64) (result string) {

	result, _ = ArrInt64ToString(arr)

	return

}

func ArrStringsToString(arr []string) (result string, err error) {

	result = ""

	for indx, a := range arr {
		if indx == 0 {
			result = fmt.Sprintf("\"%v\"", a)
		} else {
			result = fmt.Sprintf("%v,\"%v\"", result, a)
		}
	}

	if result == "" {
		err = errors.New(fmt.Sprintf("array is empty, err: %v", arr))
	}

	return

}

func ArrStringsToStringForce(arr []string) (result string) {

	result, _ = ArrStringsToString(arr)

	return

}

func ArrStringsToArrIntForce(arrStr []string) (arrInt []int) {

	arrInt = []int{}

	for _, str := range arrStr {
		i, _ := strconv.Atoi(str)
		arrInt = append(arrInt, i)
	}

	return arrInt

}

func ArrStringsToArrInt(arrStr []string) (arrInt []int, err error) {

	arrInt = []int{}

	for _, str := range arrStr {
		i, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		arrInt = append(arrInt, i)
	}

	return

}

func ArrStringsToArrFloat64(arrStr []string) (arrFloat64 []float64, err error) {
	arrFloat64 = []float64{}
	for _, str := range arrStr {
		i, err := StringToFloat64(str)
		if err != nil {
			return nil, err
		}
		arrFloat64 = append(arrFloat64, i)
	}
	return
}

func ArrStringsToArrFloat64Force(arrStr []string) (arrFloat64 []float64) {
	arrFloat64 = []float64{}
	for _, str := range arrStr {
		i, _ := StringToFloat64(str)
		arrFloat64 = append(arrFloat64, i)
	}
	return
}

func ArrStringsToArrInt64Force(arrStr []string) (arrInt64 []int64, err error) {

	arrInt64 = []int64{}

	for _, str := range arrStr {
		i, _ := strconv.ParseInt(str, 10, 64)
		arrInt64 = append(arrInt64, i)
	}

	return

}

func ArrStringsToArrInt64(arrStr []string) (arrInt64 []int64, err error) {

	arrInt64 = []int64{}

	for _, str := range arrStr {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		arrInt64 = append(arrInt64, i)
	}

	return

}

func ArrIntToArrString(arrInt []int) (arrStr []string) {

	arrStr = []string{}

	for _, i := range arrInt {
		str := IntToString(i)
		arrStr = append(arrStr, str)
	}

	return arrStr

}

func StringToFloat64(str string) (result float64, err error) {

	result, err = strconv.ParseFloat(str, 64)

	return
}

func StringToFloat64Force(str string) (result float64) {
	result, _ = strconv.ParseFloat(str, 64)
	return
}

func StringToInt(str string) (result int, err error) {

	result, err = strconv.Atoi(str)

	return

}

func StringToIntForce(str string) (result int) {

	result, _ = strconv.Atoi(str)

	return

}

func Int64ToString(intNum int64) (result string) {

	result = strconv.FormatInt(intNum, 10)

	return

}

func StringToInt64(str string) (result int64, err error) {

	result, err = strconv.ParseInt(str, 10, 64)

	return

}

func StringToInt64Force(str string) (result int64) {
	result, _ = strconv.ParseInt(str, 10, 64)

	return
}

func IntToString(intNum int) (result string) {

	result = strconv.Itoa(intNum)

	return

}

func Float64ToString(f float64) (result string) {

	result = strconv.FormatFloat(f, 'f', -1, 64)

	return

}
