package helper

import (
	"reflect"
	"strings"
)

func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func GetBSONTagMap(s interface{}, v map[string]any) map[string]any {
	m := make(map[string]any)
	t := reflect.TypeOf(s).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		bsonTag := getStructTag(f, "bson")
		if bsonTag != "" {
			bson := strings.Split(bsonTag, ",")
			if bson[0] != "" {
				m[bson[0]] = 1
			}
		}
	}

	for k, v := range v {
		m[k] = v
	}

	// remove value 0
	for k, v := range m {
		if v == 0 {
			delete(m, k)
		}
	}

	return m
}
