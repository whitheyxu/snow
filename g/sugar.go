// sugar
package g

import "reflect"

func If(condition bool, trueVal, falseVal interface{}) interface{} {

	if condition {
		return trueVal
	}

	return falseVal

}

func IsEmpty(a interface{}) bool {

	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	return v.Interface() == reflect.Zero(v.Type()).Interface()

}
