package util

import (
	"encoding/csv"
	"os"
	"sync"
)

const checkFailedPath = "check_failed.csv"
const importFailedPath = "import_failed.csv"

type UnexpectedRecorder struct {
	checkFailedPath       string
	checkFailedCsvWriter  *csv.Writer
	importFailedPath      string
	importFailedCsvWriter *csv.Writer
	mu                    sync.Mutex
}

func NewDefaultUnexpectedRecorder() *UnexpectedRecorder {
	return &UnexpectedRecorder{
		checkFailedPath:  checkFailedPath,
		importFailedPath: importFailedPath,
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

// RecordImportError records the import error into the csv file.
// thread safe.
func (u *UnexpectedRecorder) RecordImportError(err error) error {
	if err == nil {
		return nil
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	// Initialize the csv writer.
	if u.importFailedCsvWriter == nil {
		uerr := u.initImportFailedCsvWriter()
		if uerr != nil {
			return uerr
		}
	}

	// Write the error into the csv file.
	return u.importFailedCsvWriter.Write([]string{err.Error()})
}

func (u *UnexpectedRecorder) initImportFailedCsvWriter() error {
	file, err := os.Create(u.importFailedPath)
	if err != nil {
		return err
	}

	u.importFailedCsvWriter = csv.NewWriter(file)
	return nil
}

func (u *UnexpectedRecorder) Flush() {
	if u.checkFailedCsvWriter != nil {
		u.checkFailedCsvWriter.Flush()
	}

	if u.importFailedCsvWriter != nil {
		u.importFailedCsvWriter.Flush()
	}
}
