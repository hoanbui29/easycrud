package easycrud

import "errors"

var (
	ErrModelMustBeStruct    = errors.New("model must be a struct")
	ErrTableNotDefined      = errors.New("table name not defined")
	ErrPrimaryKeyNotDefined = errors.New("primary key not defined")
)
