package csv

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

type CsvReader struct {
	delimiter  byte
	escapeChar byte
	hasHeader  bool
}

var (
	errMismatchedQuotes = errors.New("mismatched quotes")
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

	inEscapeChar := false
	quoteCount := 0

	for {
		ch, err := bufReader.ReadByte()

		// if we reach end of file and the line is well-formed, we should add the line
		// to records as there was no new line.
		if err == io.EOF {
			if inEscapeChar {
				return nil, errMismatchedQuotes
			}
			record = append(record, field.String())

			// skip empty line.
			if len(record) > 0 {
				records = append(records, record)
			}
			break
		}
		if err != nil {
			return nil, err
		}

		switch ch {
		case c.escapeChar:
			quoteCount++
			nextCh, err := bufReader.Peek(1)

			// if there are consecutive quotes, it means that it is being escaped.
			if err == nil && nextCh[0] == c.escapeChar {
				// consume next quote as well but we only write 1 quote to the current line buffer.
				bufReader.ReadByte()
				field.WriteByte(c.escapeChar)
			} else {
				// keeps track of whether we are in between quotes.
				inEscapeChar = !inEscapeChar
			}

		case c.delimiter:
			// if in between quote, then we should consider the delimiter.
			if inEscapeChar {
				field.WriteByte(ch)
			} else {
				record = append(record, field.String())
				field.Reset()
			}

		// in windows, new lines are \r\n instead of just \n. so here we are techincally
		// skipping \r and wait for the next byte which will definitely be \n to
		// append the field to record.
		case '\r':
		case '\n':
			// this handles multiline items.
			if inEscapeChar {
				field.WriteByte(ch)
			} else {
				record = append(record, field.String())
				field.Reset()
				records = append(records, record)
				record = []string{}
			}

		default:
			field.WriteByte(ch)
		}
	}

	return records, nil
}
