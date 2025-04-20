package csv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

type CsvReader struct {
	delimiter  byte
	escapeChar byte
	hasHeader  bool
}

var (
	errMismatchedEscapeChar = errors.New("mismatched escape char")
	errUnexpectedEscapeChar = errors.New("unexpected escape char")
	errUnexpectedChar       = errors.New("unexpected char")
	errWrongNumFields       = errors.New("wrong number of fields")
)

func NewCsvReader(options ...ReaderOption) *CsvReader {
	reader := &CsvReader{delimiter: ',', escapeChar: '"'}
	for _, op := range options {
		op(reader)
	}

	return reader
}

func (c *CsvReader) Read(input io.Reader) ([][]string, error) {
	bufReader := bufio.NewReader(input)
	var records [][]string
	var record []string
	var field bytes.Buffer

	inQuotes := false
	justClosedQuote := false
	lineNum := 1
	colNum := 0
	expectedNumFields := 0

	for {
		ch, err := bufReader.ReadByte()
		colNum++

		if err == io.EOF {
			if inQuotes {
				return nil, fmt.Errorf("%w at line %d, column %d", errMismatchedEscapeChar, lineNum, colNum)
			}
			if justClosedQuote || field.Len() > 0 || len(record) > 0 {
				record = append(record, field.String())
				if lineNum > 1 && len(record) != expectedNumFields {
					return nil, fmt.Errorf("%w at line %d, expected %d, got %d", errWrongNumFields, lineNum, expectedNumFields, len(record))
				}
				records = append(records, record)
			}
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read error at line %d, column %d: %w", lineNum, colNum, err)
		}

		switch ch {
		case c.escapeChar:
			peek, err := bufReader.Peek(1)
			if inQuotes {
				if err == nil && peek[0] == c.escapeChar {
					// Escaped quote
					bufReader.ReadByte()
					colNum++
					field.WriteByte(c.escapeChar)
				} else {
					// Possible closing quote
					inQuotes = false
					justClosedQuote = true
				}
			} else if field.Len() == 0 {
				// Starting quoted field
				inQuotes = true
				justClosedQuote = false
			} else {
				return nil, fmt.Errorf("%w at line %d, column %d", errUnexpectedEscapeChar, lineNum, colNum)
			}

		case c.delimiter:
			if inQuotes {
				field.WriteByte(ch)
			} else if justClosedQuote || !inQuotes {
				record = append(record, field.String())
				field.Reset()
				justClosedQuote = false
			}

		case '\n':
			if inQuotes {
				field.WriteByte(ch)
			} else {
				record = append(record, field.String())
				field.Reset()
				if len(record) > 0 {
					records = append(records, record)
				}
				if lineNum == 1 {
					expectedNumFields = len(record)
				}
				record = []string{}
				lineNum++
				colNum = 0
				justClosedQuote = false
			}

		case '\r':
			// skip, wait for \n
			continue

		default:
			if justClosedQuote {
				return nil, fmt.Errorf("%w at line %d, column %d, char %c", errUnexpectedChar, lineNum, colNum, ch)
			}
			field.WriteByte(ch)
		}
	}

	return records, nil
}
