package models

import (
	"fmt"
	"github.com/albrow/go-data-parser"
)

// Checks multiple key value pairs in the form of a map to see if they are unique.
// The keys of values are taken to mean the column names, and the corresponding values
// are the values at that column name for a given row. If any values not unique, adds
// an error to val with a detailed message.
func mValidateUnique(val *data.Validator, table string, values map[string]string, msgFmt string) error {
	for key, value := range values {
		msg := fmt.Sprintf(msgFmt, key)
		if err := validateUnique(val, table, key, value, msg); err != nil {
			return err
		}
	}
	return nil
}

// Checks if a single value is unique. Uniqueness means that there is no row in table
// with column equal to value. If it is not unique, adds an error to val with a detailed
// message.
func validateUnique(val *data.Validator, table string, column string, value string, msg string) error {
	stmt := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s.%s=$1", table, table, column)
	if count, err := Db.SelectInt(stmt, value); err != nil {
		return err
	} else if count == 1 {
		val.AddError(column, msg)
	}
	return nil
}
