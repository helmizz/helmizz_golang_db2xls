package helmizz_db2xls

import (
	"database/sql"
	"encoding/json"
	"bytes"
	"fmt"
	"log"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
)

// Fungsi untuk membersihkan karakter yang tidak diinginkan
func sanitizeJSON(input []byte) []byte {
    sanitized := bytes.ReplaceAll(input, []byte("\\\""), []byte("\""))
    sanitized = bytes.ReplaceAll(sanitized, []byte("\\"), []byte("\\\\"))
	sanitized = bytes.ReplaceAll(sanitized, []byte("\r"), []byte(" "))
	sanitized = bytes.ReplaceAll(sanitized, []byte("\n"), []byte(" "))
    return sanitized
}

type QueryConfig struct {
	SQL       string `json:"sql"`
	SheetName string `json:"sheet_name"`
}

type Config struct {
	TypeDB     string        `json:"type_db"`
	Host       string        `json:"host"`
	Port       string        `json:"port"`
	Database   string        `json:"database"`
	Username   string        `json:"username"`
	Password   string        `json:"password"`
	Query      []QueryConfig `json:"query"`
	OutputFile string        `json:"output_file"`
}

func ExportToXLSX(configJSON string, queryJSON string) error {
	var config Config

	if configJSON != "" {
		configFile, err := os.ReadFile(configJSON)
		if err != nil {
			return fmt.Errorf("failed to read config file: %v", err)
		}

		configFile = sanitizeJSON(configFile)

		err = json.Unmarshal(configFile, &config)
		if err != nil {
			return fmt.Errorf("failed to parse config JSON: %v", err)
		}
	} else if queryJSON != "" {
		queryJSON := string(sanitizeJSON([]byte(queryJSON)))
		err := json.Unmarshal([]byte(queryJSON), &config.Query)
		if err != nil {
			return fmt.Errorf("failed to parse inline query JSON: %v", err)
		}
		config.OutputFile = "result/output.xlsx"
	} else {
		return fmt.Errorf("either config_json or query must be provided")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.Username, config.Password, config.Host, config.Port, config.Database)
	dbConn, err := sql.Open(config.TypeDB, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	file := xlsx.NewFile()

	headerStyle := xlsx.NewStyle()
	headerStyle.Font.Bold = true
	headerStyle.Font.Size = 10
	headerStyle.Font.Color = "FFFFFFFF"
	headerStyle.Fill = *xlsx.NewFill("solid", "FF4F81BD", "FF4F81BD")
	headerStyle.Alignment.Horizontal = "center"
	headerStyle.Alignment.Vertical = "center"
	headerStyle.ApplyFont = true
	headerStyle.ApplyFill = true
	headerStyle.ApplyAlignment = true

	dataStyle := xlsx.NewStyle()
	dataStyle.Font.Size = 8
	dataStyle.Alignment.Horizontal = "left"
	dataStyle.Alignment.Vertical = "center"
	dataStyle.ApplyFont = true
	dataStyle.ApplyAlignment = true

	for _, query := range config.Query {
		rows, err := dbConn.Query(query.SQL)
		if err != nil {
			return fmt.Errorf("failed to execute query (%s): %v", query.SQL, err)
		}
		defer rows.Close()

		sheet, err := file.AddSheet(query.SheetName)
		if err != nil {
			return fmt.Errorf("failed to create sheet (%s): %v", query.SheetName, err)
		}

		row := sheet.AddRow()
		columns, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns for query (%s): %v", query.SQL, err)
		}
		for _, colName := range columns {
			cell := row.AddCell()
			cell.Value = colName
			cell.SetStyle(headerStyle)
		}

		for rows.Next() {
			values := make([]interface{}, len(columns))
			valuePointers := make([]interface{}, len(columns))
			for i := range values {
				valuePointers[i] = &values[i]
			}

			err := rows.Scan(valuePointers...)
			if err != nil {
				return fmt.Errorf("failed to scan row for query (%s): %v", query.SQL, err)
			}

			row := sheet.AddRow()
			for _, value := range values {
				cell := row.AddCell()
				cell.Value = fmt.Sprintf("%v", value)
				cell.SetStyle(dataStyle)
			}
		}

		if err := rows.Err(); err != nil {
			return fmt.Errorf("row iteration error for query (%s): %v", query.SQL, err)
		}
	}

	err = file.Save(config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to save XLSX file (%s): %v", config.OutputFile, err)
	}

	return nil
}
