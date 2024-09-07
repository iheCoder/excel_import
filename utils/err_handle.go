package util

import (
	"errors"
	"fmt"
)

const (
	rowErrSep  = ", "
	errLineSep = "; \n"
)

type ErrBuilder struct {
	errs   []error
	header string
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

func (eb *ErrBuilder) AddWithContent(s string, err error) {
	if err == nil {
		return
	}

	eb.errs = append(eb.errs, errors.New(fmt.Sprintf("\t内容 %s 错误: %s", s, err.Error())))
}

func (eb *ErrBuilder) AddHeader(header string) {
	eb.header = header
}

func (eb *ErrBuilder) Build() error {
	if len(eb.errs) == 0 {
		return nil
	}

	errStr := eb.header

	for _, err := range eb.errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + errLineSep
	}

	// remove last "; \n"
	errStr = errStr[:len(errStr)-len(errLineSep)]

	return errors.New(errStr)
}

func CombineErrors(row int, errs ...error) error {
	errStr := fmt.Sprintf("第%d行数据错误 \n", row+1)
	for _, err := range errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + "\n"
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
	errStr = errStr[:len(errStr)-len(rowErrSep)] + "行数据错误 \n"

	for _, err := range errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + "\n"
	}

	return errors.New(errStr)
}
