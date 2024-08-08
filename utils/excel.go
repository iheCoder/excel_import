package util

import (
	"encoding/csv"
	"fmt"
	"github.com/tealeg/xlsx"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ExcelOpenMode int32

// WriteExcelColumnContent 写入Excel列内容，支持CSV和XLSX格式
func WriteExcelColumnContent(path string, content map[int][]string) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return writeCSV(path, content)
	case ".xlsx":
		return writeXLSX(path, content)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

func writeCSV(path string, content map[int][]string) error {
	return nil
}

func writeXLSX(path string, content map[int][]string) error {
	f, err := xlsx.OpenFile(path)
	if err != nil {
		return err
	}

	if len(f.Sheets) == 0 {
		_, err = f.AddSheet("Sheet1")
		if err != nil {
			return err
		}
	}

	sheet := f.Sheets[0]

	for i, col := range content {
		for j, cell := range col {
			// skip header
			sheet.Cell(j+1, i).SetString(cell)
		}
	}

	return f.Save(path)
}

func ReadExcelContent(path string) ([][]string, error) {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return readCSV(path)
	case ".xlsx":
		return readXLSX(path)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}
}

// readCSV 读取CSV文件
func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

// readXLSX 读取XLSX文件
func readXLSX(path string) ([][]string, error) {
	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}

	var records [][]string
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			var record []string
			for _, cell := range row.Cells {
				text := cell.String()
				record = append(record, text)
			}
			records = append(records, record)
		}
	}
	return records, nil
}

func FormatCell(cell string) string {
	return strings.TrimSpace(cell)
}

// DivideSheetsIntoTables 将Excel文件中的每个Sheet拆分为单独的文件，并返回拆分后的文件绝对路径
func DivideSheetsIntoTables(path string) ([]string, error) {
	f, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}

	filePaths := make([]string, 0, len(f.Sheets))
	for _, sheet := range f.Sheets {
		table := xlsx.NewFile()
		tableSheet, err := table.AddSheet(sheet.Name)
		if err != nil {
			return nil, err
		}
		for _, row := range sheet.Rows {
			tableRow := tableSheet.AddRow()
			for _, cell := range row.Cells {
				text := cell.String()
				tableCell := tableRow.AddCell()
				tableCell.SetString(text)
			}
		}
		tablePath := strings.TrimSuffix(path, filepath.Ext(path)) + "_" + sheet.Name + ".xlsx"
		if err = table.Save(tablePath); err != nil {
			return nil, err
		}

		filePaths = append(filePaths, tablePath)
	}

	return filePaths, nil
}

// CombineTablesIntoOne 将多个Excel文件合并为一个
func CombineTablesIntoOne(paths ...string) error {
	if len(paths) == 0 {
		return nil
	}

	f, err := xlsx.OpenFile(paths[0])
	if err != nil {
		return err
	}

	for i := 1; i < len(paths); i++ {
		table, err := xlsx.OpenFile(paths[i])
		if err != nil {
			return err
		}

		for _, sheet := range table.Sheets {
			newSheet, ok := f.Sheet[sheet.Name]
			if !ok {
				newSheet, err = f.AddSheet(sheet.Name)
				if err != nil {
					return err
				}
			}

			for _, row := range sheet.Rows {
				newRow := newSheet.AddRow()
				for _, cell := range row.Cells {
					text := cell.String()
					newCell := newRow.AddCell()
					newCell.SetString(text)
				}
			}
		}
	}

	return f.Save(paths[0])
}

// DivideExcelContent 拆分Excel数据到多个文件
func DivideExcelContent(path string, rowLimit int) ([]string, error) {
	records, err := ReadExcelContent(path)
	if err != nil {
		return nil, err
	}

	if len(records) <= rowLimit {
		return []string{path}, nil
	}

	var filePaths []string
	for i := 0; i < len(records); i += rowLimit {
		subRecords := records[i:min(i+rowLimit, len(records))]
		subPath := strings.TrimSuffix(path, filepath.Ext(path)) + fmt.Sprintf("_%d", i) + filepath.Ext(path)

		if err = WriteExcelContent(subPath, subRecords); err != nil {
			return nil, err
		}
		filePaths = append(filePaths, subPath)
	}

	return filePaths, nil
}

func WriteExcelContent(path string, content [][]string) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return writeCSVContent(path, content)
	case ".xlsx":
		return writeXLSXContent(path, content)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

func writeCSVContent(path string, content [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	for _, record := range content {
		if err = writer.Write(record); err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}

func writeXLSXContent(path string, content [][]string) error {
	f := xlsx.NewFile()
	sheet, err := f.AddSheet("Sheet1")
	if err != nil {
		return err
	}

	for _, record := range content {
		row := sheet.AddRow()
		for _, cell := range record {
			row.AddCell().SetString(cell)
		}
	}

	return f.Save(path)
}
