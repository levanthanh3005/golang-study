package functions

import (
	"fmt"
	"encoding/json"
	)

// func studyJson () {
// 	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
// 	var f interface{}
// 	err := json.Unmarshal(b, &f)
// 	fmt.Println(err)
	
// 	m := f.(map[string]interface{})

// 	for k, v := range m {
// 	    switch vv := v.(type) {
// 	    case string:
// 	        fmt.Println(k, "is string", vv)
// 	    case float64:
// 	        fmt.Println(k, "is float64", vv)
// 	    case []interface{}:
// 	        fmt.Println(k, "is an array:")
// 	        for i, u := range vv {
// 	            fmt.Println(i, u)
// 	        }
// 	    default:
// 	        fmt.Println(k, "is of a type I don't know how to handle")
// 	    }
// 	}
// 	fmt.Println(m["Name"])
// 	// fmt.Println(f.(map["Age"]int{}))
// }

func getJSONValue(objStr string, key string) interface{} {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil
	}
	m := f.(map[string]interface{})
	return m[key]
}

func getJSONValue2(objStr string, key string) interface{} {
	var m map[string]interface{}
	err := stringToMapWithError(objStr, &m)
	if err != nil {
		return nil
	}
	return m[key]
}

func stringToMap(objStr string) map[string]interface{} {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil
	}
	m := f.(map[string]interface{})
	return m
}

func stringToMapWithError(objStr string, m *map[string]interface{}) error {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	*m = f.(map[string]interface{})
	return nil
}

func objToString(v interface{}) string {
	str := fmt.Sprintf("%v", v)
	return str
}

func showJSON(lsKV map[string]interface{}) {
	for k, v := range lsKV {
	    switch vv := v.(type) {
	    case string:
	        fmt.Println(k, "is string", vv)
	    case float64:
	        fmt.Println(k, "is float64", vv)
	    case []interface{}:
	        fmt.Println(k, "is an array:")
	        for i, u := range vv {
	            fmt.Println(i, u)
	        }
	    default:
	        fmt.Println(k, "is of a type I don't know how to handle")
	    }
	}
}