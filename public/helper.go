package easycrud

import (
	"github.com/lib/pq"
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

func isIgnore(s structTag) bool {
	_, table := s["table"]
	_, ignore := s["ignore"]

	return table || ignore
}

func (e EasyCRUD[TEntity, TKey]) getValueByFieldName(model TEntity, fieldName string) (interface{}, error) {
	field := reflect.ValueOf(model).FieldByName(fieldName)
	kind := field.Kind()

	switch kind {
	case reflect.Slice:
		return pq.Array(field.Interface()), nil
	case reflect.Array:
		return pq.Array(field.Interface()), nil
	default:
		return field.Interface(), nil
	}
}

func (e EasyCRUD[TEntity, TKey]) getFieldPointers(model *TEntity) ([]interface{}, error) {
	t := reflect.TypeOf(*model)
	v := reflect.ValueOf(model)
	numFields := t.NumField()
	fieldPointers := make([]interface{}, 0)

	for i := 0; i < numFields; i++ {
		fieldValue := v.Elem().Field(i)
		tags := getStructTag(t.Field(i))

		if isIgnore(tags) {
			continue
		}

		if fieldValue.CanAddr() {
			fieldPointers = append(fieldPointers, fieldValue.Addr().Interface())
		} else {
			return nil, ErrFieldCannotBePointer
		}
	}

	return fieldPointers, nil
}
