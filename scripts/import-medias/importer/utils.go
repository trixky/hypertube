package importer

import (
	"database/sql"
	"fmt"
	"os"
)

func makeNullInt32(value *int32) (null_int32 sql.NullInt32) {
	if value == nil || *value == 0 {
		return
	}
	null_int32.Int32 = *value
	null_int32.Valid = true
	return
}

func makeNullFloat64(value *float64) (null_float64 sql.NullFloat64) {
	if value == nil || *value == 0 {
		return
	}
	null_float64.Float64 = *value
	null_float64.Valid = true
	return
}

func makeNullString(value *string) (null_string sql.NullString) {
	if value == nil || *value == "\\N" {
		return
	}
	null_string.String = *value
	null_string.Valid = true
	return
}

func openTsvFile(path string) (file *os.File, err error) {
	file, err = os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
