package klcsv

import (
	"encoding/csv"
	"os"
	"strings"
)

func WriteCSV(path string, input [][]string, delimeter *rune) error {

	ps := strings.Split(path, "/")
	if len(ps) > 1 {
		directory := strings.Join(ps[:len(ps)-1], "/")
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return nil
		}
		//Create a folder/directory at a full qualified path
		err := os.Mkdir(directory, 0755)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}

	w := csv.NewWriter(file)
	if delimeter != nil {
		w.Comma = *delimeter
	}
	defer w.Flush()

	err = w.WriteAll(input) // calls Flush internally
	if err != nil {
		return err
	}
	return nil
}