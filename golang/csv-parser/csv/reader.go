package csv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	errMismatchedEscapeChar = errors.New("mismatched escape char")
	errUnexpectedEscapeChar = errors.New("unexpected escape char")
	errWrongNumFields       = errors.New("wrong number of fields")
)

type CsvReader struct {
	delimiter  byte
	escapeChar byte

	reader      *bufio.Reader
	readerState *readerState
}

// readerState keeps track of the current state of the reader between reads
type readerState struct {
	lineNum             int
	expectedNumOfFields int
	escaping            bool
	escaped             bool
	field               bytes.Buffer
	record              []string
	records             [][]string
}

func NewCsvReader(inputReader io.Reader, readerOptions ...ReaderOption) *CsvReader {
	bufReader := bufio.NewReader(inputReader)
	cr := &CsvReader{
		delimiter:  ',',
		escapeChar: '"',
		reader:     bufReader,
		readerState: &readerState{
			lineNum:  1,
			escaping: false,
			escaped:  false,
			field:    bytes.Buffer{},
			record:   []string{},
			records:  [][]string{},
		},
	}

	for _, op := range readerOptions {
		op(cr)
	}

	return cr
}

func (cr *CsvReader) Read() ([][]string, error) {
	for {
		ch, err := cr.reader.ReadByte()

		// if end of file, append the last line and return
		if err == io.EOF {
			err := cr.appendLine()
			if err != nil {
				return nil, err
			}
			return cr.readerState.records, nil
		}

		if err != nil {
			return nil, err
		}

		switch ch {
		case cr.delimiter:
			err = cr.handleDelimiter()
		case cr.escapeChar:
			err = cr.handleEscapeChar()
		// in windows the newline is \r\n, so we can skip the \r and process the next byte which is the \n
		case '\r':
		case '\n':
			err = cr.handleNewLine()
		default:
			err = cr.handleDefault(ch)
		}

		if err != nil {
			return nil, err
		}
	}
}

func (cr *CsvReader) handleDelimiter() error {
	if cr.readerState.escaping {
		cr.readerState.field.WriteByte(cr.delimiter)
		return nil
	}

	err := cr.appendField()
	if err != nil {
		return err
	}

	return nil
}

func (cr *CsvReader) handleEscapeChar() error {
	if cr.readerState.escaping {
		nextCh, peakErr := cr.reader.Peek(1)
		if peakErr == nil && nextCh[0] == cr.escapeChar {
			cr.readerState.field.WriteByte(cr.escapeChar)
			cr.reader.ReadByte()
		} else {
			cr.readerState.escaping = false
			cr.readerState.escaped = true
		}
	} else if cr.readerState.field.Len() == 0 {
		cr.readerState.escaping = true
		cr.readerState.escaped = false
	} else {
		return fmt.Errorf("%w at line: %d", errUnexpectedEscapeChar, cr.readerState.lineNum)
	}

	return nil
}

func (cr *CsvReader) handleNewLine() error {
	if cr.readerState.escaping {
		cr.readerState.field.WriteByte('\n')
		return nil
	}
	err := cr.appendLine()
	if err != nil {
		return err
	}

	return nil
}

func (cr *CsvReader) handleDefault(ch byte) error {
	if cr.readerState.escaped {
		return fmt.Errorf("%w at line: %d", errMismatchedEscapeChar, cr.readerState.lineNum)
	}

	return cr.readerState.field.WriteByte(ch)
}

func (cr *CsvReader) appendField() error {
	cr.readerState.record = append(cr.readerState.record, cr.readerState.field.String())
	cr.readerState.field.Reset()

	cr.readerState.escaping = false
	cr.readerState.escaped = false

	return nil
}

func (cr *CsvReader) appendLine() error {
	cr.appendField()
	if cr.readerState.lineNum == 1 {
		cr.readerState.expectedNumOfFields = len(cr.readerState.record)
	} else if len(cr.readerState.record) != cr.readerState.expectedNumOfFields {
		return fmt.Errorf("%w at line: %d", errWrongNumFields, cr.readerState.lineNum)
	}

	cr.readerState.records = append(cr.readerState.records, cr.readerState.record)
	cr.readerState.record = []string{}
	cr.readerState.lineNum++

	cr.readerState.escaping = false
	cr.readerState.escaped = false
	return nil
}
