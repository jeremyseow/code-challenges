package csv

import (
	ocsv "encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
			expected:   [][]string{{"1", "2", "\"3\""}, {"4", "5", "\n6"}, {"7", "8", "\",9\""}},
			err:        nil,
		},
		{
			name:       "different delimiter test",
			filePath:   "data/test3.csv",
			delimiter:  '-',
			excapeChar: '"',
			expected:   [][]string{{"1", "2", "3"}, {"4", "5", "6"}, {"7", "\"8\"", ",9"}},
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
			expected:    nil,
			err:         errUnexpectedEscapeChar,
		},
		{
			name:        "quotes in the middle",
			stringInput: "a,b,c\na,b,c\n\"aa\"a,b,c",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    nil,
			err:         errMismatchedEscapeChar,
		},
		{
			name:        "empty field",
			stringInput: ",,",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{{"", "", ""}},
			err:         nil,
		},
		{
			name:        "missing column",
			stringInput: "a,b,c\nd,e",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    nil,
			err:         errWrongNumFields,
		},
		{
			name:        "missing column in the middle",
			stringInput: "a,b\nd,e\nf,g,h",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    nil,
			err:         errWrongNumFields,
		},
		{
			name:        "mix",
			stringInput: "a,b,c\nd,e,\"\"\"f\"\"\"\ng,h,\"i\n\"",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{{"a", "b", "c"}, {"d", "e", "\"f\""}, {"g", "h", "i\n"}},
			err:         nil,
		},
		{
			name:        "single field",
			stringInput: "a",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{{"a"}},
			err:         nil,
		},
		{
			name:        "whitespace",
			stringInput: " ",
			delimiter:   ',',
			excapeChar:  '"',
			expected:    [][]string{{" "}},
			err:         nil,
		},
	}

	for _, testCase := range testCases {
		currTestCase := testCase
		t.Run(currTestCase.name, func(t *testing.T) {
			t.Parallel()

			var testInput io.Reader
			var validateInput io.Reader
			if currTestCase.filePath != "" {
				file, err := os.Open(currTestCase.filePath)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer file.Close()
				testInput = file

				file2, err := os.Open(currTestCase.filePath)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer file2.Close()
				validateInput = file2
			} else if currTestCase.stringInput != "" {
				testInput = strings.NewReader(currTestCase.stringInput)
				validateInput = strings.NewReader(currTestCase.stringInput)
			} else {
				t.Errorf("test case %s has no input", currTestCase.name)
			}

			ocsvReader := ocsv.NewReader(validateInput)
			ocsvReader.Comma = rune(currTestCase.delimiter)
			recs, errValidate := ocsvReader.ReadAll()
			fmt.Println(currTestCase.name, recs)
			if errValidate != nil {
				fmt.Println(currTestCase.name, errValidate)
			}

			csvReader := NewCsvReader(testInput, WithDelimiter(currTestCase.delimiter), WithEscapeChar(currTestCase.excapeChar))
			records, errTest := csvReader.Read()
			assert.True(t, errors.Is(errTest, currTestCase.err))
			assert.Equal(t, currTestCase.expected, records)

			if errTest != nil {
				fmt.Println(currTestCase.name, errTest)
			}

			assert.Equal(t, currTestCase.expected, recs)
		})
	}
}

func BenchmarkRead(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		file, err := os.Open("data/test1.csv")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		csvReader := NewCsvReader(file, WithDelimiter(','), WithEscapeChar('"'))
		_, err = csvReader.Read()
		if err != nil {
			b.Errorf("error %v", err)
		}
	}
}
