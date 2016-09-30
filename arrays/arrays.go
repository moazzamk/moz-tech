package arrays

import "reflect"


func InArray(haystack interface{}, needle interface{}) bool {
	switch reflect.TypeOf(haystack).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(haystack)

			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(needle, s.Index(i).Interface()) == true {
					return true
				}
			}
	}

	return false
}
