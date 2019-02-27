package jl

type JavaException struct {
	FramesOmitted int64            `json:"frames_omitted"`
	Message       string           `json:"message"`
	Module        string           `json:"module"`
	StackTrace    []JavaStacktrace `json:"stack_trace"`
	Type          string           `json:"type"`
}

type JavaStacktrace struct {
	File   string `json:"file"`
	Func   string `json:"func"`
	Line   int64  `json:"line"`
	Module string `json:"module"`
}