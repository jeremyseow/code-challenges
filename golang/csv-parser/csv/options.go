package csv

type ReaderOption func(*CsvReader)

func WithDelimiter(delimiter byte) ReaderOption {
	return func(reader *CsvReader) {
		reader.delimiter = delimiter
	}
}

func WithEscapeChar(escapeChar byte) ReaderOption {
	return func(reader *CsvReader) {
		reader.escapeChar = escapeChar
	}
}

func WithHeader(hasHeader bool) ReaderOption {
	return func(reader *CsvReader) {
		reader.hasHeader = hasHeader
	}
}
