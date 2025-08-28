package cli

// InputReader interface for reading user input
type InputReader interface {
	ReadLine() (string, error)
}

// OutputWriter interface for writing output
type OutputWriter interface {
	WriteLine(string) error
	Write([]byte) (int, error)
}