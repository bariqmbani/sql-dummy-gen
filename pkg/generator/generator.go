package generator

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/bariqmbani/sql-dummy-gen/pkg/parser"
)

type Config struct {
	AuditColumn string
	StartTime   time.Time
	EndTime     time.Time
}

func GenerateRow(schema parser.TableSchema, config Config) map[string]interface{} {
	row := make(map[string]interface{})

	for _, col := range schema.Columns {
		if col.Name == config.AuditColumn {
			// Generate random time
			delta := config.EndTime.Sub(config.StartTime)
			nsec := delta.Nanoseconds()
			if nsec <= 0 {
				// If invalid range, just use start time
				row[col.Name] = config.StartTime
			} else {
				randNsec := rand.Int63n(nsec)
				row[col.Name] = config.StartTime.Add(time.Duration(randNsec))
			}
			continue
		}

		// Basic type mapping
		switch {
		case strings.Contains(col.Type, "tinyint"):
			row[col.Name] = rand.Intn(128) // Safe dummy valid for standard tinyint
		case strings.Contains(col.Type, "int"):
			row[col.Name] = rand.Intn(1000)
		case strings.Contains(col.Type, "char") || strings.Contains(col.Type, "text"):
			// Generate random string
			length := 10
			if col.Length != nil && *col.Length > 0 {
				if *col.Length < 10 {
					length = *col.Length
				}
			}
			row[col.Name] = randomString(length)
		case strings.Contains(col.Type, "bool"):
			row[col.Name] = rand.Intn(2) // 0 or 1
		case strings.Contains(col.Type, "time") || strings.Contains(col.Type, "date"):
			// Just current time for other date columns for now, or random
			row[col.Name] = time.Now()
		case strings.Contains(col.Type, "float") || strings.Contains(col.Type, "double") || strings.Contains(col.Type, "decimal"):
			val := rand.Float64() * 1000.0
			if strings.Contains(col.Type, "decimal") && col.Scale != nil {
				pow := math.Pow(10, float64(*col.Scale))
				val = math.Round(val*pow) / pow
			}
			row[col.Name] = val
		default:
			// Fallback string
			row[col.Name] = fmt.Sprintf("val_%s", col.Name)
		}
	}
	return row
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
