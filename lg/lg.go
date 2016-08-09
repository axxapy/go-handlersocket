package lg

import (
	"log"
	"fmt"
)

const (
	DEBUG   = 1
	ERROR   = 2
	INFO    = 4
	VERBOSE = 8
	WARNING = 16
	WTF     = 32
	ALL     = DEBUG ^ ERROR ^ INFO ^ VERBOSE ^ WARNING ^ WTF
)

var (
	log_level = WARNING ^ ERROR ^ INFO
)

func SetLogLevel(level int) {
	log_level = level
}

func _print(tag string, level string, msg ...interface{}) {
	log.Print("[" + level + "] [" + tag + "] " + fmt.Sprintln(msg...))
}

func V(tag string, msg ...interface{}) {
	if log_level & VERBOSE != 0 {
		_print(tag, "V", msg)
	}
}

func E(tag string, msg ...interface{}) {
	if log_level & ERROR != 0 {
		_print(tag, "E", msg)
	}
}

func W(tag string, msg ...interface{}) {
	if log_level & WARNING != 0 {
		_print(tag, "W", msg)
	}
}

func I(tag string, msg ...interface{}) {
	if log_level & INFO != 0 {
		_print(tag, "I", msg)
	}
}