package utils

import (
	"encoding/json"
	"reflect"
)

func DeepEqualStruct(vx, vy interface{}) bool {

	var v1, v2 interface{}

	bytes, _ := json.Marshal(vx)
	json.Unmarshal(bytes, &v1)
	bytes, _ = json.Marshal(vy)
	json.Unmarshal(bytes, &v2)

	return deepEqual(v1, v2)
}

// Equal checks equality between 2 struct
func deepEqual(vx, vy interface{}) bool {
	if reflect.TypeOf(vx) != reflect.TypeOf(vy) {
		return false
	}

	switch x := vx.(type) {
	case map[string]interface{}:
		y := vy.(map[string]interface{})

		if len(x) != len(y) {
			return false
		}

		for k, v := range x {
			val2 := y[k]

			if (v == nil) != (val2 == nil) {
				return false
			}

			if !deepEqual(v, val2) {
				return false
			}
		}

		return true
	case []interface{}:
		y := vy.([]interface{})

		if len(x) != len(y) {
			return false
		}

		var matches int
		flagged := make([]bool, len(y))
		for _, v := range x {
			for i, v2 := range y {
				if deepEqual(v, v2) && !flagged[i] {
					matches++
					flagged[i] = true
					break
				}
			}
		}
		return matches == len(x)
	default:
		return vx == vy
	}
}
