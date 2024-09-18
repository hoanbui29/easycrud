package easycrud

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/hoanbui29/easycrud/internal/validator"
)

type fieldModel struct {
	name       string
	fieldType  reflect.Type
	tags       structTag
	columnName string
}

type EasyCRUDModel[TEntity any, TKey any] interface {
	Create() (TKey, error)
}

type EasyCRUD[TEntity any, TKey any] struct {
	db *sql.DB
}

func (e *EasyCRUD[TEntity, TKey]) validateModel() (string, fieldModel, []fieldModel, error) {
	var model TEntity
	var defaultKey fieldModel
	fields := make([]fieldModel, 0)

	t := reflect.TypeOf(model)
	if !validator.IsStruct(t) {
		return "", defaultKey, fields, ErrModelMustBeStruct
	}

	numFields := t.NumField()

	var tableName string
	var pkeyField fieldModel

	for i := 0; i < numFields; i++ {
		field := t.Field(i)

		tags := getStructTag(field)

		var columnName string

		if isColumnName, cName := getColumnName(tags); isColumnName {
			columnName = cName
		} else {
			columnName = field.Name
		}

		model := fieldModel{
			name:       field.Name,
			fieldType:  field.Type,
			tags:       tags,
			columnName: columnName,
		}

		if isTableName, tName := getTableName(tags); isTableName {
			tableName = tName
			continue
		}

		if _, ok := tags["pkey"]; ok {
			pkeyField = model
			continue
		}

		fields = append(fields, model)
	}

	if tableName == "" {
		return "", defaultKey, []fieldModel{}, ErrTableNotDefined
	}

	if pkeyField.name == "" {
		return "", defaultKey, []fieldModel{}, ErrPrimaryKeyNotDefined
	}

	return tableName, pkeyField, fields, nil
}

func (e *EasyCRUD[TEntity, TKey]) Create(input TEntity) (TKey, error) {
	var defaultKey TKey
	table, pkey, fields, err := e.validateModel()

	if err != nil {
		return defaultKey, err
	}

	var columns []string
	var quotedColumns []string
	var placeholders []string
	var values []interface{}

	for _, field := range fields {
		columns = append(columns, field.columnName)
		quotedColumns = append(quotedColumns, fmt.Sprintf(`"%s"`, field.columnName))
		placeholders = append(placeholders, fmt.Sprintf(`$%d`, len(values)+1))
		values = append(values, reflect.ValueOf(input).FieldByName(field.name).Interface())
	}

	stmt := fmt.Sprintf(`INSERT INTO %s (%s)
        VALUES (%s) RETURNING %s`, table, strings.Join(quotedColumns, ","), strings.Join(placeholders, ","), pkey.columnName)

	row := e.db.QueryRow(stmt, values...)

	if err := row.Scan(&defaultKey); err != nil {
		return defaultKey, err
	}

	return defaultKey, nil
}
