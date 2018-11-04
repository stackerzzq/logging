package logging

import (
	"fmt"
)

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	gray    = 37
)

type Level int8

const (
	Debug Level = iota - 1
	Info
	Warn
	Error
	Panic
	Fatal

	minLevel = Debug
	maxLevel = Fatal
)

func (l Level) String() string {
	switch l{
	case Debug:
		return "debug"
	case Info:
		return "info"
	case Warn:
		return "warn"
	case Error:
		return "error"
	case Panic:
		return "panic"
	case Fatal:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}