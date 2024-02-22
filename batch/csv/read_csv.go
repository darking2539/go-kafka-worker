package klcsv

import (
	"encoding/csv"
	"os"
)

func ReadCSV(path string, callback func(record []string) error, separator *rune, skip int) error {

	// os.Open() opens specific file in
	// read-only mode and this return
	// a pointer of type os.File
	file, err := os.Open(path)

	// Checks for the error
	if err != nil {
		return err
	}

	// Closes the file
	defer file.Close()

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameter
	// and this creates a new csv.Reader that reads
	// from the file
	reader := csv.NewReader(file)
	if separator != nil {
		reader.Comma = *separator
	}

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		return err
	}

	// Loop to iterate through
	// and print each of the string slice
	for row, record := range records {
		if row < skip {
			continue
		}
		
		err = callback(record)
		if err != nil {
			return err
		}
	}
	return nil
}