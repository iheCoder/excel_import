package util

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"sync"
)

const (
	checkFailedPath    = "check_failed.csv"
	importFailedPath   = "import_failed.csv"
	unexpectedJsonPath = "unexpected.jsonl"
)

type UnexpectedRecorder struct {
	checkFailedPath       string
	checkFailedCsvWriter  *csv.Writer
	importFailedPath      string
	importFailedCsvWriter *csv.Writer
	mu                    sync.Mutex
	importFailedJsonPath  string
	importFailedJsonFile  *os.File
	jsonDecoder           *json.Decoder
	jsonEncoder           *json.Encoder
}

func NewDefaultUnexpectedRecorder() *UnexpectedRecorder {
	return &UnexpectedRecorder{
		checkFailedPath:      checkFailedPath,
		importFailedPath:     importFailedPath,
		importFailedJsonPath: unexpectedJsonPath,
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

func (u *UnexpectedRecorder) RecordImportErrorWithContent(err error, contents ...string) error {
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

	// Write the error and content into the csv file.
	errWithContent := append([]string{err.Error()}, contents...)
	return u.importFailedCsvWriter.Write(errWithContent)
}

func (u *UnexpectedRecorder) RecordContentJson(content any) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// initialize the json file.
	if u.importFailedJsonFile == nil {
		file, err := os.Create(u.importFailedJsonPath)
		if err != nil {
			return err
		}

		u.importFailedJsonFile = file
		u.jsonEncoder = json.NewEncoder(file)
	}

	// write the content into the json file.
	return u.jsonEncoder.Encode(content)
}

func (u *UnexpectedRecorder) IterateJsonContent(contentObj any) bool {
	// init the json decoder
	if u.jsonDecoder == nil {
		file, err := os.Open(u.importFailedJsonPath)
		if err != nil {
			panic(err)
		}

		u.jsonDecoder = json.NewDecoder(file)
	}

	// decode the content
	if err := u.jsonDecoder.Decode(contentObj); err != nil {
		return false
	}

	return true
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

	if u.importFailedJsonFile != nil {
		u.importFailedJsonFile.Close()
	}
}
