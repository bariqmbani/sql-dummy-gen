package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bariqmbani/sql-dummy-gen/pkg/builder"
	"github.com/bariqmbani/sql-dummy-gen/pkg/generator"
	"github.com/bariqmbani/sql-dummy-gen/pkg/parser"
)

type Config struct {
	DBType      string
	DDLFile     string
	OutputFile  string
	AuditColumn string
	TimeRange   string
	RowCount    int
}

func main() {
	var config Config
	flag.StringVar(&config.DBType, "db", "mysql", "Target database type (mysql)")
	flag.StringVar(&config.DDLFile, "ddl", "", "Path to the DDL file containing CREATE TABLE statements")
	flag.StringVar(&config.OutputFile, "output", "", "Output file for INSERT statements (default: output-{ddl-file}.sql)")
	flag.StringVar(&config.AuditColumn, "created-col", "", "Name of the created date audit column to populate with random dates")
	flag.StringVar(&config.TimeRange, "time-range", "", "Time range for audit column (format: YYYY-MM-DD,YYYY-MM-DD). Defaults to today if empty.")
	flag.IntVar(&config.RowCount, "num", 1, "Number of rows to generate per table")
	flag.Parse()

	if config.DDLFile == "" {
		fmt.Println("Error: -ddl flag is required")
		flag.Usage()
		os.Exit(1)
	}

	if config.OutputFile == "" {
		base := filepath.Base(config.DDLFile)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		config.OutputFile = fmt.Sprintf("output-%s.sql", name)
	}

	// Parsing Time Range
	var startTime, endTime time.Time
	if config.TimeRange != "" {
		parts := strings.Split(config.TimeRange, ",")
		if len(parts) != 2 {
			log.Fatalf("Error: Invalid time-range format. Expected YYYY-MM-DD,YYYY-MM-DD")
		}
		var err error
		startTime, err = time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
		if err != nil {
			log.Fatalf("Error parsing start time: %v", err)
		}
		endTime, err = time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
		if err != nil {
			log.Fatalf("Error parsing end time: %v", err)
		}
		// Make sure endTime is inclusive of the day
		endTime = endTime.Add(24*time.Hour - 1*time.Nanosecond)
	} else {
		// Default to today
		now := time.Now()
		startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endTime = startTime.Add(24*time.Hour - 1*time.Nanosecond)
	}

	fmt.Printf("Parsing DDL file: %s\n", config.DDLFile)
	ddlContent, err := os.ReadFile(config.DDLFile)
	if err != nil {
		log.Fatalf("Error reading DDL file: %v", err)
	}

	schemas, err := parser.ParseDDL(string(ddlContent))
	if err != nil {
		log.Fatalf("Error parsing DDL: %v", err)
	}

	genConfig := generator.Config{
		AuditColumn: config.AuditColumn,
		StartTime:   startTime,
		EndTime:     endTime,
	}

	f, err := os.Create(config.OutputFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer f.Close()

	totalStart := time.Now()
	batchSize := 1000

	for _, schema := range schemas {
		fmt.Printf("Generating data for table: %s\n", schema.TableName)

		rows := make([]map[string]interface{}, 0, batchSize)
		for i := 0; i < config.RowCount; i++ {
			row := generator.GenerateRow(schema, genConfig)
			rows = append(rows, row)

			if len(rows) >= batchSize {
				inserts := builder.BuildInsert(config.DBType, schema.TableName, rows)
				if _, err := f.WriteString(inserts + "\n\n"); err != nil {
					log.Fatalf("Error writing to file: %v", err)
				}
				rows = rows[:0] // Clear and keep capacity
			}

			// Always log progress periodically
			if (i+1)%1000 == 0 || i+1 == config.RowCount {
				fmt.Printf("\rGenerated %d/%d rows", i+1, config.RowCount)
			}
		}

		// Flush remaining
		if len(rows) > 0 {
			inserts := builder.BuildInsert(config.DBType, schema.TableName, rows)
			if _, err := f.WriteString(inserts + "\n\n"); err != nil {
				log.Fatalf("Error writing to file: %v", err)
			}
		}
		fmt.Println()
	}

	elapsed := time.Since(totalStart)
	fmt.Printf("Successfully generated INSERT statements to %s\n", config.OutputFile)
	fmt.Printf("Time taken: %v\n", elapsed)
}
