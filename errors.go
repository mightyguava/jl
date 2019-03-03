package jl

// JavaExceptions represents a list of Java exceptions, in the default JSON format.
type JavaExceptions []*JavaException

// JavaException contains a single Java exception in a causal chain.
type JavaException struct {
	FramesOmitted int64           `json:"frames_omitted"`
	Message       string          `json:"message"`
	Module        string          `json:"module"`
	StackTrace    []JavaStackItem `json:"stack_trace"`
	Type          string          `json:"type"`
}

// JavaStackItem is a single line in a stack trace.
type JavaStackItem struct {
	File   string `json:"file"`
	Func   string `json:"func"`
	Line   int64  `json:"line"`
	Module string `json:"module"`
}

// LogrusError encapsulates a logrus style error
type LogrusError struct {
	Error string
	Stack string
}
