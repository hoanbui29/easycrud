package easycrud

import "errors"

var (
	ErrModelMustBeStruct    = errors.New("model must be a struct")
	ErrTableNotDefined      = errors.New("table name not defined")
	ErrPrimaryKeyNotDefined = errors.New("primary key not defined")
	ErrFieldCannotInterface = errors.New("field cannot be interface")
	ErrFieldCannotBePointer = errors.New("field cannot be pointer")
)
