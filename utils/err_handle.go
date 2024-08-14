package util

import (
	"errors"
	"fmt"
)

const (
	rowErrSep = ", "
)

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
