package utility

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

func GetJSONValue(objStr string, key string) interface{} {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil
	}
	m := f.(map[string]interface{})
	return m[key]
}

func GetJSONValue2(objStr string, key string) interface{} {
	var m map[string]interface{}
	err := StringToMapWithError(objStr, &m)
	if err != nil {
		return nil
	}
	return m[key]
}

func StringToMap(objStr string) map[string]interface{} {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return nil
	}
	m := f.(map[string]interface{})
	return m
}

func StringToMapWithError(objStr string, m *map[string]interface{}) error {
	b := []byte(objStr)
	var f interface{}
	err := json.Unmarshal(b, &f)
	if err != nil {
		return err
	}
	*m = f.(map[string]interface{})
	return nil
}

func ObjToString(v interface{}) string {
	str := fmt.Sprintf("%v", v)
	return str
}

// func stringToJsonString(objStr) (string,error) {
// 	content, _ := json.Marshal(res)
// 	return string(content), nil
// }

func MapToJsonString(m map[string]interface{}) (string,error){
	content, _ := json.Marshal(m)
	return string(content), nil
}

func ShowJSON(lsKV map[string]interface{}) {
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