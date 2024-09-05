package util

import (
	"excel_import"
	"reflect"
	"strconv"
	"strings"
)

const (
	excelImportTag = "exi"
	invalidIndex   = -1
	gormTag        = "gorm"
)

func ParseTag(st any) []*excel_import.ExcelImportTagAttr {
	// get struct if st is a pointer
	if reflect.TypeOf(st).Kind() == reflect.Ptr {
		st = reflect.ValueOf(st).Elem().Interface()
	}

	// get reflect type
	t := reflect.TypeOf(st)

	// get the number of fields
	numField := t.NumField()

	// create a slice to store the tag attributes
	tagAttrs := make([]*excel_import.ExcelImportTagAttr, 0, numField)

	var index int
	// iterate over the fields
	for i := 0; i < numField; i++ {
		// get the field
		field := t.Field(i)

		// get the tag
		tag := field.Tag.Get(excelImportTag)

		// parse the tag
		tagAttr := parseTag(tag)

		// handle the case when the column index is not set
		if tagAttr.ColumnIndex == invalidIndex {
			tagAttr.ColumnIndex = index
			index++
		} else {
			index = tagAttr.ColumnIndex + 1
		}

		// append the tag attributes to the slice
		tagAttrs = append(tagAttrs, tagAttr)
	}

	return tagAttrs
}

func parseTag(tag string) *excel_import.ExcelImportTagAttr {
	// create a new tag attribute
	tagAttr := &excel_import.ExcelImportTagAttr{
		ColumnIndex: invalidIndex,
	}

	// handle the case when the tag is empty
	if len(tag) == 0 {
		return tagAttr
	}

	// split the tag by comma
	tagParts := strings.Split(tag, ",")

	// iterate over the tag parts
	for _, part := range tagParts {
		// split the part by colon
		partParts := strings.Split(part, ":")

		// get the key and value
		key := partParts[0]
		value := partParts[1]

		// set the key and value to the tag attribute
		switch key {
		case "index":
			ci, err := strconv.Atoi(value)
			if err == nil {
				tagAttr.ColumnIndex = ci
			}
		case "rewrite":
			rw, err := strconv.ParseBool(value)
			if err == nil {
				tagAttr.Rewrite = rw
			}
		case "chk":
			tagAttr.Check = excel_import.CheckMode(value)
		case "tree":
			tagAttr.Tree = excel_import.TreeFlag(value)
		}
	}

	return tagAttr
}

func ParseGormTag(st any) []*excel_import.GormTag {
	// get struct if st is a pointer
	if reflect.TypeOf(st).Kind() == reflect.Ptr {
		st = reflect.ValueOf(st).Elem().Interface()
	}

	// get reflect type
	t := reflect.TypeOf(st)

	// get the number of fields
	numField := t.NumField()

	// create a slice to store the tag attributes
	tagAttrs := make([]*excel_import.GormTag, 0, numField)

	// iterate over the fields
	for i := 0; i < numField; i++ {
		// get the field
		field := t.Field(i)

		// get the tag
		tag := field.Tag.Get(gormTag)

		// parse the tag
		tagAttr := parseGormTag(tag)

		// append the tag attributes to the slice
		tagAttrs = append(tagAttrs, tagAttr)
	}

	return tagAttrs
}

func parseGormTag(tag string) *excel_import.GormTag {
	// create a new tag attribute
	tagAttr := &excel_import.GormTag{}

	// handle the case when the tag is empty
	if len(tag) == 0 {
		return tagAttr
	}

	// split the tag by comma
	tagParts := strings.Split(tag, ",")

	// iterate over the tag parts
	for _, part := range tagParts {
		// split the part by colon
		partParts := strings.Split(part, ":")

		// get the key and value
		key := partParts[0]
		value := partParts[1]

		// set the key and value to the tag attribute
		switch key {
		case "column":
			tagAttr.Column = value
		case "primary_key":
			tagAttr.PrimaryKey = true
		case "auto_increment":
			tagAttr.AutoIncrement = true
		case "default":
			tagAttr.Default = value
		case "not null":
			tagAttr.NotNull = true
		case "size":
			size, err := strconv.Atoi(value)
			if err == nil {
				tagAttr.Size = size
			}
		case "type":
			tagAttr.Type = value
		}
	}

	return tagAttr
}
