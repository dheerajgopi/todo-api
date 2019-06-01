package common

import (
	"strconv"
	"strings"
)

// Page stores pagination data
type Page struct {
	Limit  int
	Offset int64
	Sort   []*Sort
	Cursor []*Sort
}

// TODO: sort object -> cursor based pagination based on sort
type Sort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
	LastVal   string `json:"lastVal"`
}

func (sort *Sort) SetAscending() *Sort {
	sort.Direction = "asc"

	return sort
}

func (sort *Sort) SetDescending() *Sort {
	sort.Direction = "desc"

	return sort
}

func (sort *Sort) ValidateLastVal(fieldType FieldType) bool {
	val := sort.LastVal

	if val == "" {
		return false
	}

	isValid := true

	switch fieldType {
	case BOOLEAN:
		valLowerCase := strings.ToLower(val)
		isValid = valLowerCase == "true" || valLowerCase == "false"
	case INT64:
		_, err := strconv.ParseInt(val, 10, 64)

		if err != nil {
			isValid = false
		}
	case UNIXTIME:
		_, err := strconv.ParseUint(val, 10, 32)

		if err != nil {
			isValid = false
		}
	}

	return isValid
}

type FieldType int

const (
	BOOLEAN FieldType = 1 + iota
	INT64
	UNIXTIME
	STRING
)

var fieldTypes = [...]string{
	"BOOLEAN",
	"INT64",
	"UNIXTIME",
	"STRING",
}

func (fieldType FieldType) String() string {
	return fieldTypes[fieldType-1]
}
