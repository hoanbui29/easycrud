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
	Create(input TEntity) (TKey, error)
	Detail(key TKey) (TEntity, error)
	Update(input TEntity) (bool, error)
	Delete(key TKey) (bool, error)
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
		value, err := e.getValueByFieldName(input, field.name)

		if err != nil {
			return defaultKey, err
		}
		values = append(values, value)
	}

	stmt := fmt.Sprintf(`INSERT INTO %s (%s)
        VALUES (%s) RETURNING %s`, table, strings.Join(quotedColumns, ","), strings.Join(placeholders, ","), pkey.columnName)

	row := e.db.QueryRow(stmt, values...)

	if err := row.Scan(&defaultKey); err != nil {
		return defaultKey, err
	}

	return defaultKey, nil
}

func (e *EasyCRUD[TEntity, TKey]) Detail(key TKey) (TEntity, error) {
	var defaultEntity TEntity
	table, pkey, _, err := e.validateModel()

	if err != nil {
		return defaultEntity, err
	}

	row := e.db.QueryRow(fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1`, table, pkey.columnName), key)

	var entity TEntity

	scanArgs, err := e.getFieldPointers(&entity, false)

	if err != nil {
		return defaultEntity, err
	}

	err = row.Scan(scanArgs...)

	if err != nil {
		return defaultEntity, err
	}

	return entity, nil
}

func (e *EasyCRUD[TEntity, TKey]) Update(input TEntity) (bool, error) {
	table, pkey, fields, err := e.validateModel()

	if err != nil {
		return false, err
	}

	var columns []string
	var updateColumns []string
	var values []interface{}

	pkeyValue, err := e.getValueByFieldName(input, pkey.name)

	if err != nil {
		return false, err
	}

	for _, field := range fields {
		columns = append(columns, field.columnName)
		updateColumns = append(updateColumns, fmt.Sprintf(`"%s" = $%d`, field.columnName, len(values)+1))
		value, err := e.getValueByFieldName(input, field.name)

		if err != nil {
			return false, err
		}
		values = append(values, value)
	}

	stmt := fmt.Sprintf(`UPDATE %s SET %s WHERE %s = $%d`, table, strings.Join(updateColumns, ","), pkey.columnName, len(values)+1)
	result, err := e.db.Exec(stmt, append(values, pkeyValue)...)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	return rowsAffected > 0, err
}

func (e *EasyCRUD[TEntity, TKey]) Delete(key TKey) (bool, error) {
	table, pkey, _, err := e.validateModel()

	if err != nil {
		return false, err
	}

	result, err := e.db.Exec(fmt.Sprintf(`DELETE FROM %s WHERE %s = $1`, table, pkey.columnName), key)

	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()

	return rowsAffected > 0, err

}
