package recordtocsv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// RecordToCSVService manages the process of recording data to CSV files.
type RecordToCSVService struct {
	// Dir is the desired folder name, e.g., "files/record".
	Dir string

	// Filename is the desired base filename, e.g., "agoda_booking_record".
	Filename string

	// Column represents the CSV header columns.
	// Example: []string{"id", "request", "response"}
	Column []string

	// RecordType determines the time-based suffix for the filename: "daily", "monthly", "yearly".
	RecordType string
}

// NewRecordToCSV creates and returns a new RecordToCSVService instance.
func NewRecordToCSV(dir, filename string, column []string, recordType string) *RecordToCSVService {
	return &RecordToCSVService{
		Dir:        dir,
		Filename:   filename,
		Column:     column,
		RecordType: recordType,
	}
}

// Record processes the given payload and appends it to a time-suffixed CSV file.
func (r *RecordToCSVService) Record(payload interface{}) error {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Log the error or return a more specific error if needed
		return fmt.Errorf("failed to load time zone 'Asia/Jakarta': %w", err)
	}

	timeNow := time.Now().In(loc)
	var suffix string

	switch r.RecordType {
	case "daily":
		suffix = timeNow.Format("2006_01_02")
	case "monthly":
		suffix = timeNow.Format("2006_01")
	case "yearly":
		suffix = timeNow.Format("2006")
	default:
		return fmt.Errorf("unsupported record type: %q. Must be 'daily', 'monthly', or 'yearly'", r.RecordType)
	}

	// Use filepath.Join for robust path construction across different OS
	filePath := filepath.Join(r.Dir, fmt.Sprintf("%s_%s.csv", r.Filename, suffix))

	// Ensure the directory exists
	if err := os.MkdirAll(r.Dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", r.Dir, err)
	}

	if err := r.Append(filePath, r.Column, payload); err != nil {
		return fmt.Errorf("failed to append record to %q: %w", filePath, err)
	}
	return nil
}

// Append writes a single data record to the specified CSV file.
// It handles creating the file and writing headers if the file doesn't exist.
func (r *RecordToCSVService) Append(filename string, column []string, data interface{}) error {
	// Open the file in append mode. If it doesn't exist, create it.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open/create CSV file %q: %w", filename, err)
	}
	defer file.Close() // Ensure the file is closed

	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush() // Ensure data is flushed to the file

	// Check if the file is empty (newly created or truly empty) to write headers
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %q: %w", filename, err)
	}

	if stat.Size() == 0 { // File is empty, write header
		if err := csvWriter.Write(column); err != nil {
			return fmt.Errorf("failed to write CSV header to %q: %w", filename, err)
		}
	}

	// Convert payload to a map for easy column-based access
	var dataMap map[string]interface{}
	// Using json.Marshal then json.Unmarshal is acceptable for generic interface{}
	// but direct struct field mapping is more efficient if payload type is known.
	// For this generic case, it's a common pattern.
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payload to JSON: %w", err)
	}
	if err := json.Unmarshal(jsonBytes, &dataMap); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to map: %w", err)
	}

	record := make([]string, len(column))
	for i, col := range column {
		if val, ok := dataMap[col]; ok && val != nil {
			record[i] = fmt.Sprintf("%v", val) // Use %v to handle various types
		} else {
			record[i] = "" // Ensure empty string for missing or nil values
		}
	}

	if err := csvWriter.Write(record); err != nil {
		return fmt.Errorf("failed to write CSV record to %q: %w", filename, err)
	}

	// Check for any errors that occurred during writing
	if err := csvWriter.Error(); err != nil && err != io.EOF { // io.EOF can be ignored when flushing
		return fmt.Errorf("CSV writer encountered an error: %w", err)
	}

	return nil
}
