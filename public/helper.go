package easycrud

import (
	"reflect"
	"strings"
)

type structTag map[string]string

func getStructTag(t reflect.StructField) structTag {
	result := make(map[string]string)
	rawTagValue := t.Tag.Get("easycrud")
	if rawTagValue == "" {
		return result
	}

	tagValues := strings.Split(rawTagValue, ",")

	for _, tagValue := range tagValues {
		tag := strings.Split(tagValue, "=")
		if len(tag) == 1 {
			result[tag[0]] = ""
		} else {
			result[tag[0]] = tag[1]
		}
	}

	return result
}

func getTagValue(s structTag, key string) (bool, string) {
	value, ok := s[key]

	return ok, value
}

func getColumnName(s structTag) (bool, string) {
	return getTagValue(s, "column")
}

func getTableName(s structTag) (bool, string) {
	return getTagValue(s, "table")
}
