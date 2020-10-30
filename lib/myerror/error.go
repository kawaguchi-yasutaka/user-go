package myerror

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
)

type ErrorType string

type CustomError struct {
	Message    string
	ErrorType  ErrorType
	StatusCode int
}

func (e CustomError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

func NewCustomError(msg string, errorType ErrorType, StatusCode int) CustomError {
	return CustomError{
		Message:    msg,
		ErrorType:  errorType,
		StatusCode: StatusCode,
	}
}

const (
	ErrorDBDuplicateEntry = "db_duplicate_entry"
	ErrorDBError          = "db_error"
)

var duplicateEntryRegExp = regexp.MustCompile("Error 1062")

func DBError(err error) CustomError {
	switch {
	case duplicateEntryRegExp.MatchString(err.Error()):
		return NewCustomError(err.Error(), ErrorDBDuplicateEntry, http.StatusInternalServerError)
	default:
		return NewCustomError(err.Error(), ErrorDBError, http.StatusInternalServerError)
	}
}

func EqualErrorType(x error, y error) bool {
	if reflect.DeepEqual(x, y) {
		return true
	}
	xCErr, ok := x.(CustomError)
	if !ok {
		return false
	}
	yCErr, ok := y.(CustomError)
	if !ok {
		return false
	}
	if xCErr.ErrorType == yCErr.ErrorType {
		return true
	}
	return false
}
