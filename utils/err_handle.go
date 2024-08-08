package util

import (
	"errors"
	"fmt"
)

func CombineErrors(row int, errs ...error) error {
	errStr := fmt.Sprintf("第%d行数据格式错误: ", row+1)
	for _, err := range errs {
		if err == nil {
			continue
		}
		errStr += err.Error() + "; "
	}
	return errors.New(errStr)
}
