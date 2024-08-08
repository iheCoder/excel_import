package util

import (
	"encoding/csv"
	"os"
)

const checkFailedPath = "check_failed.csv"

type UnexpectedRecorder struct {
	checkFailedPath      string
	checkFailedCsvWriter *csv.Writer
}

func NewDefaultUnexpectedRecorder() *UnexpectedRecorder {
	return &UnexpectedRecorder{
		checkFailedPath: checkFailedPath,
	}
}

func (u *UnexpectedRecorder) RecordCheckError(err error) error {
	if err == nil {
		return nil
	}

	// Initialize the csv writer.
	if u.checkFailedCsvWriter == nil {
		uerr := u.initCheckFailedCsvWriter()
		if uerr != nil {
			return uerr
		}
	}

	// Write the error into the csv file.
	return u.checkFailedCsvWriter.Write([]string{err.Error()})
}

func (u *UnexpectedRecorder) initCheckFailedCsvWriter() error {
	file, err := os.Create(u.checkFailedPath)
	if err != nil {
		return err
	}

	u.checkFailedCsvWriter = csv.NewWriter(file)
	return nil
}

func (u *UnexpectedRecorder) Flush() {
	if u.checkFailedCsvWriter != nil {
		u.checkFailedCsvWriter.Flush()
	}
}
