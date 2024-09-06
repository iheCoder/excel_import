package util

import (
	"errors"
	"fmt"
)

const (
	rowErrSep = ", "
)

type ErrBuilder struct {
	errs    []error
	content string
}

func NewErrBuilder() *ErrBuilder {
	return &ErrBuilder{
		errs: make([]error, 0),
	}
}

func (eb *ErrBuilder) Add(err error) {
	if err == nil {
		return
	}

	eb.errs = append(eb.errs, err)
}

func (eb *ErrBuilder) AddContent(content string) {
	eb.content = content
}

func (eb *ErrBuilder) Build() error {
	if len(eb.errs) == 0 {
		return nil
	}

	errStr := eb.content
	if errStr == "" {
		errStr = "数据错误: "
	}

	for _, err := range eb.errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + "; "
	}

	return errors.New(errStr)
}

func CombineErrors(row int, errs ...error) error {
	errStr := fmt.Sprintf("第%d行数据错误: ", row+1)
	for _, err := range errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + "; "
	}
	return errors.New(errStr)
}

func CombineRowsErrors(rows []int, errs ...error) error {
	if len(errs) == 0 {
		return nil
	}

	errStr := "第"
	for _, row := range rows {
		errStr += fmt.Sprintf("%d", row+1)
		errStr += rowErrSep
	}
	errStr = errStr[:len(errStr)-len(rowErrSep)] + "行数据错误: "

	for _, err := range errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + "; "
	}

	return errors.New(errStr)
}
