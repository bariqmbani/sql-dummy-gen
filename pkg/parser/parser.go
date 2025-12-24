package parser

import (
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

type Column struct {
	Name   string
	Type   string
	Length *int
	Scale  *int
}

type TableSchema struct {
	TableName string
	Columns   []Column
}

func ParseDDL(ddlContent string) ([]TableSchema, error) {
	// split ddlContent by semicolon to handle multiple statements?
	// sqlparser.Parse parses a single statement.
	// But `sqlparser` might not support multiple statements in one string directly if it expects one.
	// However, sqlparser.SplitStatementToPieces exists or we can just iterate.
	// Let's rely on Tokenizer if needed, or simple split for now, assuming well-formed input.
	// Actually, let's use sqlparser.Tokenize to be safer? Or just try to parse repeatedly.
	// A simple approach for this MVP: split by ;

	var schemas []TableSchema
	// Remove comments? sqlparser handles some, but let's just feed raw first.

	// Naive split by semicolon for multiple statements
	statements := strings.Split(ddlContent, ";")

	for _, stmtStr := range statements {
		stmtStr = strings.TrimSpace(stmtStr)
		if stmtStr == "" {
			continue
		}

		stmt, err := sqlparser.Parse(stmtStr)
		if err != nil {
			fmt.Printf("Warning: skipping statement due to parse error or non-SQL: %v\nSQL: %s\n", err, stmtStr)
			continue
		}

		ddl, ok := stmt.(*sqlparser.DDL)
		if !ok {
			continue
		}

		if ddl.Action != sqlparser.CreateStr {
			continue
		}

		// Debug print
		// fmt.Printf("DDL Action: %v\n", ddl.Action)

		tableName := ddl.NewName.Name.String()
		var columns []Column

		for _, col := range ddl.TableSpec.Columns {
			colName := col.Name.String()
			colType := col.Type.Type // e.g., "int", "varchar"

			// Normalize type
			colType = strings.ToLower(colType)

			var length *int
			var scale *int

			if col.Type.Length != nil {
				valStr := string(col.Type.Length.Val)
				var l int
				if _, err := fmt.Sscanf(valStr, "%d", &l); err == nil {
					length = &l
				}
			}

			if col.Type.Scale != nil {
				valStr := string(col.Type.Scale.Val)
				var s int
				if _, err := fmt.Sscanf(valStr, "%d", &s); err == nil {
					scale = &s
				}
			}

			columns = append(columns, Column{
				Name:   colName,
				Type:   colType,
				Length: length,
				Scale:  scale,
			})
		}

		schemas = append(schemas, TableSchema{
			TableName: tableName,
			Columns:   columns,
		})
	}

	return schemas, nil
}
