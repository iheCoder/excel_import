package util

import (
	"errors"
	"os"
	"testing"
)

func TestRecordImportErrorWithContent(t *testing.T) {
	// init recorder
	recorder := NewDefaultUnexpectedRecorder()

	// record error
	err := errors.New("import error")
	ierr := recorder.RecordImportError(err)
	if ierr != nil {
		t.Fatalf("expect nil, but got %s", ierr.Error())
	}

	// record error with one content
	content1 := "content1"
	ierr = recorder.RecordImportErrorWithContent(err, content1)
	if ierr != nil {
		t.Fatalf("expect nil, but got %s", ierr.Error())
	}

	// record error with three contents
	content2 := "content2"
	content3 := "content3"
	ierr = recorder.RecordImportErrorWithContent(err, content1, content2, content3)
	if ierr != nil {
		t.Fatalf("expect nil, but got %s", ierr.Error())
	}

	// flush
	recorder.Flush()

	// check the csv file
	contents, err := ReadExcelContent(recorder.importFailedPath)
	if err != nil {
		t.Fatalf("expect nil, but got %s", err.Error())
	}
	if len(contents) != 3 {
		t.Fatalf("expect 3, but got %d", len(contents))
	}

	// check the content
	if contents[0][0] != "import error" {
		t.Fatalf("expect import error and content1, but got %s and %s", contents[0][0], contents[0][1])
	}

	if contents[1][0] != "import error" || contents[1][1] != content1 {
		t.Fatalf("expect import error, content1 and content2, but got %s, %s and %s", contents[1][0], contents[1][1], contents[1][2])
	}

	if contents[2][0] != "import error" || contents[2][1] != content1 || contents[2][2] != content2 || contents[2][3] != content3 {
		t.Fatalf("expect import error, content1, content2 and content3, but got %s, %s, %s and %s", contents[2][0], contents[2][1], contents[2][2], contents[2][3])
	}

	// remove the csv file
	err = os.Remove(recorder.importFailedPath)
	if err != nil {
		t.Fatalf("remove csv file expect nil, but got %s", err.Error())
	}

	t.Log("RecordImportErrorWithContent success")
}
