package csv

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	testCases := []struct {
		name        string
		filePath    string
		stringInput string
		delimiter   byte
		excapeChar  byte
		expected    [][]string
		err         error
	}{
		{
			name:       "base test",
			filePath:   "data/test1.csv",
			delimiter:  ',',
			excapeChar: '"',
			expected:   [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "8", "9"}},
			err:        nil,
		},
		{
			name:       "quotes and multiline test",
			filePath:   "data/test2.csv",
			delimiter:  ',',
			excapeChar: '"',
			expected:   [][]string{{"1", "2", "\"\"3\"\""}, {"4", "5", "\n6"}, {"7", "8", "\",9\""}},
			err:        nil,
		},
		{
			name:       "different delimiter and escape character test",
			filePath:   "data/test3.csv",
			delimiter:  '-',
			excapeChar: '|',
			expected:   [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "|8|", ",9"}},
			err:        nil,
		},
		{
			name:        "string input",
			stringInput: "a,b,c",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{{"a", "b", "c"}},
			err:         nil,
		},
		{
			name:        "standalone quotes",
			stringInput: "a\"a,b,c",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{},
			err:         errMismatchedQuotes,
		},
		{
			name:        "quotes in the middle",
			stringInput: "a\"a\"a,b,c",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{},
			err:         nil,
		},
		{
			name:        "empty field",
			stringInput: ",,",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{{"", "", ""}},
			err:         nil,
		},
	}
	for _, testCase := range testCases {
		currTestCase := testCase
		t.Run(currTestCase.name, func(t *testing.T) {
			t.Parallel()
			csvReader := NewCsvReader(WithDelimiter(currTestCase.delimiter), WithEscapeChar(currTestCase.excapeChar))

			var input io.Reader
			if currTestCase.filePath != "" {
				file, err := os.Open(currTestCase.filePath)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer file.Close()
				input = file
			} else if currTestCase.stringInput != "" {
				input = strings.NewReader(currTestCase.stringInput)
			} else {
				input = nil
			}

			records, err := csvReader.Read(input)
			if err != currTestCase.err {
				t.Errorf("expected error %v, got %v", currTestCase.err, err)
			}
			if !reflect.DeepEqual(records, currTestCase.expected) {
				t.Errorf("expected records %v, got %v", currTestCase.expected, records)
			}
		})
	}
}

func BenchmarkRead(b *testing.B) {
	csvReader := NewCsvReader(WithDelimiter(','), WithEscapeChar('"'))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		file, err := os.Open("data/test1.csv")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		_, err = csvReader.Read(file)
		if err != nil {
			b.Errorf("error %v", err)
		}
	}
}
