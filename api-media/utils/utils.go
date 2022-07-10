package utils

import "database/sql"

func MakeNullInt32(value *int32) (null_int32 sql.NullInt32) {
	if value == nil {
		return
	}
	null_int32.Int32 = *value
	null_int32.Valid = true
	return
}

func MakeNullFloat64(value *float64) (null_float64 sql.NullFloat64) {
	if value == nil {
		return
	}
	null_float64.Float64 = *value
	null_float64.Valid = true
	return
}

func MakeNullString(value *string) (null_string sql.NullString) {
	if value == nil {
		return
	}
	null_string.String = *value
	null_string.Valid = true
	return
}
