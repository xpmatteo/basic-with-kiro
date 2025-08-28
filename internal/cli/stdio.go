package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// StdInputReader implements InputReader using standard input
type StdInputReader struct {
	scanner *bufio.Scanner
}

// NewStdInputReader creates a new standard input reader
func NewStdInputReader() *StdInputReader {
	return &StdInputReader{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// ReadLine reads a line from standard input
func (r *StdInputReader) ReadLine() (string, error) {
	if r.scanner.Scan() {
		return r.scanner.Text(), nil
	}
	
	if err := r.scanner.Err(); err != nil {
		return "", err
	}
	
	return "", io.EOF
}

// StdOutputWriter implements OutputWriter using standard output
type StdOutputWriter struct{}

// NewStdOutputWriter creates a new standard output writer
func NewStdOutputWriter() *StdOutputWriter {
	return &StdOutputWriter{}
}

// WriteLine writes a line to standard output
func (w *StdOutputWriter) WriteLine(line string) error {
	_, err := fmt.Println(line)
	return err
}

// Write writes data to standard output
func (w *StdOutputWriter) Write(data []byte) (int, error) {
	return os.Stdout.Write(data)
}