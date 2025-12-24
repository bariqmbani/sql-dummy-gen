package builder

import (
	"fmt"
	"strings"
	"time"
)

func BuildInsert(dbType string, tableName string, rows []map[string]interface{}) string {
	if len(rows) == 0 {
		return ""
	}

	// Assuming all rows have same keys
	var columns []string
	firstRow := rows[0]
	for k := range firstRow {
		columns = append(columns, k)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(columns, ", ")))

	var values []string
	for _, row := range rows {
		var rowValues []string
		for _, col := range columns {
			val := row[col]
			rowValues = append(rowValues, formatValue(dbType, val))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(rowValues, ", ")))
	}

	sb.WriteString(strings.Join(values, ", "))
	sb.WriteString(";")

	return sb.String()
}

func formatValue(dbType string, val interface{}) string {
	switch v := val.(type) {
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case time.Time:
		// Format based on DB type if needed, assume standard SQL 'YYYY-MM-DD HH:MM:SS'
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case int, int64, float64:
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprintf("'%v'", v)
	}
}
