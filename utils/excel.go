package util

import (
	"encoding/csv"
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ExcelOpenMode int32

const (
	defaultStartRow  = 1
	DefaultSuffixKey = ".td"
)

// WriteExcelColumnContent 写入Excel列内容，支持CSV和XLSX格式
func WriteExcelColumnContent(path string, content map[int][]string) error {
	return WriteExcelColumnContentByStartRow(path, content, defaultStartRow)
}

// WriteExcelColumnContentByStartRow 写入Excel列内容，支持CSV和XLSX格式，支持指定起始行
func WriteExcelColumnContentByStartRow(path string, content map[int][]string, startRow int) error {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".csv":
		return writeCSVByStartRow(path, content, startRow)
	case ".xlsx":
		return writeXLSXByStartRow(path, content, startRow)
	default:
		return fmt.Errorf("unsupported file type: %s", ext)
	}
}

func writeCSV(path string, content map[int][]string) error {
	return nil
}

func writeCSVByStartRow(path string, content map[int][]string, startRow int) error {
	return nil
}

func writeXLSX(path string, content map[int][]string) error {
	return writeXLSXByStartRow(path, content, defaultStartRow)
}

func writeXLSXByStartRow(path string, content map[int][]string, startRow int) error {
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
			sheet.Cell(j+startRow, i).SetString(cell)
		}
	}

	return f.Save(path)
}

// ReadExcelValidContentInCommonCase read excel content in common case, support CSV and XLSX format
// the common case is that the end row is the row that all cells are empty or the first cell is empty, and skip the header
func ReadExcelValidContentInCommonCase(path string) ([][]string, error) {
	contents, err := ReadExcelContent(path)
	if err != nil {
		return nil, err
	}

	// skip header
	contents = contents[1:]

	// end row
	for i, content := range contents {
		if DefaultRowEndFunc(content) {
			return contents[:i], nil
		}
	}

	return contents, nil
}

// ReadExcelContent read excel content from file, support CSV and XLSX format
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
	reader.FieldsPerRecord = -1
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

// DivideSheetsIntoTables 将Excel文件中的每个Sheet拆分为单独的文件，并返回拆分后的文件绝对路径
func DivideSheetsIntoTables(path string) ([]string, error) {
	return DivideSheetsIntoTablesBySuffixKey(path, "")
}

func DivideSheetsIntoTablesByDefaultSuffixKey(path string) ([]string, error) {
	return DivideSheetsIntoTablesBySuffixKey(path, DefaultSuffixKey)
}

// DivideSheetsIntoTablesBySuffixKey splits the sheets by suffix key and preserves formatting, styles as possible.
// ATTENTION that the format or style may be lost in some cases.
func DivideSheetsIntoTablesBySuffixKey(path, suffixKey string) ([]string, error) {
	f, err := xlsx.OpenFile(path)
	if err != nil {
		return nil, err
	}

	filePaths := make([]string, 0, len(f.Sheets))
	for _, sheet := range f.Sheets {
		if len(suffixKey) > 0 && !strings.HasSuffix(sheet.Name, suffixKey) {
			continue
		}

		table := xlsx.NewFile()
		tableSheet, err := table.AddSheet(sheet.Name)
		if err != nil {
			return nil, err
		}
		for _, row := range sheet.Rows {
			tableRow := tableSheet.AddRow()
			for _, cell := range row.Cells {
				tableCell := tableRow.AddCell()
				copyCell(cell, tableCell)
			}
		}
		tablePath := strings.TrimSuffix(path, filepath.Ext(path)) + "_" + strings.TrimSuffix(sheet.Name, suffixKey) + ".xlsx"
		if err = table.Save(tablePath); err != nil {
			return nil, err
		}

		filePaths = append(filePaths, tablePath)
	}

	return filePaths, nil
}

// copySheet copies the cell content and style from src to dst
func copyCell(src, dst *xlsx.Cell) {
	dst.SetString(src.String())
	if src.GetStyle() != nil {
		dst.SetStyle(src.GetStyle())
	}
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

			// skip header
			if ok {
				sheet.Rows = sheet.Rows[1:]
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

	header := records[0]

	var filePaths []string
	var tableIndex int
	for i := 1; i < len(records); i += rowLimit {
		tableIndex++
		end := min(i+rowLimit, len(records))
		subRecords := append([][]string{header}, records[i:end]...)
		subPath := strings.TrimSuffix(path, filepath.Ext(path)) + fmt.Sprintf("_%d", tableIndex) + filepath.Ext(path)

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

type treeExcelInfo struct {
	path string
	key  string
	f    *xlsx.File
}

func genTreeKey(ukCols []string) string {
	return strings.Join(ukCols, "_")
}

// DivideMultipleTreesIntoMultipleTables 将多棵树拆分为多个表
// 依据ukColIndex列索引组合唯一键，将多棵树拆分为多个表
func DivideMultipleTreesIntoMultipleTables(path, outputDir string, ukColIndex []int) ([]string, error) {
	records, err := ReadExcelContent(path)
	if err != nil {
		return nil, err
	}

	if len(records) <= 1 {
		return nil, fmt.Errorf("empty excel content")
	}

	header := records[0]

	ukTreeMap := make(map[string]*treeExcelInfo)
	for _, row := range records[1:] {
		// generate unique key
		ukCols := make([]string, 0, len(ukColIndex))
		for _, idx := range ukColIndex {
			ukCols = append(ukCols, FormatCell(row[idx]))
		}
		key := genTreeKey(ukCols)

		// get or create tree
		tree, ok := ukTreeMap[key]
		if !ok {
			f := xlsx.NewFile()
			sheet, err := f.AddSheet("Sheet1")
			if err != nil {
				return nil, err
			}

			for i, cell := range header {
				sheet.Cell(0, i).SetString(cell)
			}

			tree = &treeExcelInfo{
				path: filepath.Join(outputDir, key+".xlsx"),
				key:  key,
				f:    f,
			}
			ukTreeMap[key] = tree
		}

		// write row
		treeSheet := tree.f.Sheets[0]
		exRow := treeSheet.AddRow()
		for _, cell := range row {
			exRow.AddCell().SetString(cell)
		}
	}

	// create output dir if not exists
	if _, err = os.Stat(outputDir); os.IsNotExist(err) {
		if err = os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	var filePaths []string
	for _, tree := range ukTreeMap {
		if err = tree.f.Save(tree.path); err != nil {
			return nil, err
		}
		filePaths = append(filePaths, tree.path)
	}

	return filePaths, nil
}

// SetHyperlinksInColumn 将 Excel 文件中指定列的单元格内容设置为超链接
func SetHyperlinksInColumn(path string, urls []string, index int) error {
	// 打开Excel文件
	f, err := excelize.OpenFile(path)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}

	// 假设我们只处理第一个工作表
	sheetName := f.GetSheetName(0)

	// 将urls写入指定列，并设置为超链接
	for i, url := range urls {
		// skip header
		cell, _ := excelize.CoordinatesToCellName(index+1, i+1+1)
		f.SetCellValue(sheetName, cell, url)
		if err := f.SetCellHyperLink(sheetName, cell, url, "External"); err != nil {
			return fmt.Errorf("failed to set hyperlink: %w", err)
		}
	}

	// 保存文件
	if err := f.Save(); err != nil {
		return fmt.Errorf("failed to save Excel file: %w", err)
	}

	return nil
}
