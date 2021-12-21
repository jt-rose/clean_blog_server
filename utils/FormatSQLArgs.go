package utils

import (
	"strconv"
	"strings"
)

// takes a slice of []int and converts it to a formatted string
// that can be used as an argument in a SQL = ANY() query
func FormatSliceForSQLParams(ids []int) string {
	// convert ids to []string
var stringArgs []string
for _, id := range ids {
	stringArgs = append(stringArgs, strconv.Itoa(id))
}

// format param of SQL query
queryParam := "{"+ strings.Join(stringArgs, ",") + "}"

return queryParam
}
