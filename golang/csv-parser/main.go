package main

import (
	"fmt"
	"os"

	"github.com/jeremyseow/csv-parser/csv"
)

func main() {
	csvReader := csv.NewCsvReader(csv.WithDelimiter(','), csv.WithEscapeChar('"'))
	file, err := os.Open("csv/data/test1.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	records, err := csvReader.Read(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	for lineNum, record := range records {
		for _, item := range record {
			fmt.Printf("line %d: %s\n", lineNum+1, item)
		}
	}
}
