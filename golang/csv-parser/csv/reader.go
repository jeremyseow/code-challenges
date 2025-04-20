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
	// bufio reads or writes data in chunks rather than one byte at a time, which is more efficient.
	bufReader := bufio.NewReader(input)
	var records [][]string
	var record []string
	var field bytes.Buffer

	// keep track of if the field is being escaped.
	inEscapeChar := false

	// keep track of if we have escaped the entire field.
	closedEscapeChar := false
	lineNum := 1
	colNum := 0

	// based on the 1st row, throw error if any row has a different number of fields.
	expectedNumFields := 0

	for {
		ch, err := bufReader.ReadByte()
		colNum++

		if err == io.EOF {
			// if we are still escaping when we have reached the end, then we have a mismatched escape char.
			if inEscapeChar {
				return nil, fmt.Errorf("%w at line %d, column %d", errMismatchedEscapeChar, lineNum, colNum)
			}
			if field.Len() > 0 || len(record) > 0 {
				record = append(record, field.String())

				// if this is the only row, then we can skip the check. else we throw error if
				// any row has a different number of fields.
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
			// we want to check the next char without consuming it.
			peek, err := bufReader.Peek(1)
			if inEscapeChar {
				// if we are escaping and we have a consecutive escape char, it means we are escaping it.
				if err == nil && peek[0] == c.escapeChar {
					// consume the next escape char but only write once.
					bufReader.ReadByte()
					colNum++
					field.WriteByte(c.escapeChar)
				} else {
					// we have possible finished escape the field, but need to check if there are more
					// chars after this escape char.
					inEscapeChar = false
					closedEscapeChar = true
				}
			} else if field.Len() == 0 {
				// we are escaping the current field.
				inEscapeChar = true
				closedEscapeChar = false
			} else {
				// throw error for unexpected escape char.
				return nil, fmt.Errorf("%w at line %d, column %d", errUnexpectedEscapeChar, lineNum, colNum)
			}

		case c.delimiter:
			if inEscapeChar {
				field.WriteByte(ch)
			} else if closedEscapeChar || !inEscapeChar {
				record = append(record, field.String())
				field.Reset()
				closedEscapeChar = false
			}

		// in windows, the new line is \r\n. we can just do nothing for \r and wait for the \n.
		case '\r':
		case '\n':
			// this is for multi-line fields.
			if inEscapeChar {
				field.WriteByte(ch)
			} else {
				record = append(record, field.String())
				field.Reset()

				// skip empty lines.
				if len(record) > 0 {
					records = append(records, record)
				}

				// set the expected number of fields based on the 1st row.
				if lineNum == 1 {
					expectedNumFields = len(record)
				}

				record = []string{}
				lineNum++
				colNum = 0
				closedEscapeChar = false
			}

		default:
			// throw error if we have more chars after we are done escaping for the field.
			if closedEscapeChar {
				return nil, fmt.Errorf("%w at line %d, column %d, char %c", errUnexpectedChar, lineNum, colNum, ch)
			}
			field.WriteByte(ch)
		}
	}

	return records, nil
}
