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
			err:         errUnexpectedChar,
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
		{
			name:       "empty",
			filePath:   "data/test4.csv",
			delimiter:  ',',
			excapeChar: '"',
			expected:   nil,
			err:        nil,
		},
	}
	for _, testCase := range testCases {
		currTestCase := testCase
		t.Run(currTestCase.name, func(t *testing.T) {
			t.Parallel()

			var input io.Reader
			var input2 io.Reader
			if currTestCase.filePath != "" {
				file, err := os.Open(currTestCase.filePath)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer file.Close()
				input = file

				file2, err := os.Open(currTestCase.filePath)
				if err != nil {
					fmt.Println("Error opening file:", err)
					return
				}
				defer file2.Close()
				input2 = file2
			} else if currTestCase.stringInput != "" {
				input = strings.NewReader(currTestCase.stringInput)
				input2 = strings.NewReader(currTestCase.stringInput)
			} else {
				input = nil
			}

			ocsvReader := ocsv.NewReader(input2)
			ocsvReader.Comma = rune(currTestCase.delimiter)
			recs, err2 := ocsvReader.ReadAll()
			fmt.Println(currTestCase.name, recs)
			if err2 != nil {
				fmt.Println(currTestCase.name, err2)
			}

			csvReader := NewCsvReader(WithDelimiter(currTestCase.delimiter), WithEscapeChar(currTestCase.excapeChar))
			records, err := csvReader.Read(input)
			assert.True(t, errors.Is(err, currTestCase.err))
			assert.Equal(t, currTestCase.expected, records)

			if err != nil {
				fmt.Println(currTestCase.name, err)
			}

			assert.Equal(t, currTestCase.expected, recs)
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
