package main

type Level int64

const (
	Error Level = iota
)

type Issue struct {
	level    Level
	message  string
	row, col int
}

func NewError(message string) *Issue {
	return &Issue{level: Error, message: message}
}
