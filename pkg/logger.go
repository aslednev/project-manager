package pkg

import (
	"log"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	FATAL
)

type Logger struct {
	level Level
}

func NewLogger(level Level) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		log.Println("[DEBUG]", v)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.level <= INFO {
		log.Println("[INFO]", v)
	}
}

func (l *Logger) Warning(v ...interface{}) {
	if l.level <= WARNING {
		log.Println("[WARNING]", v)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.level <= FATAL {
		log.Fatalln("[FATAL]", v)
	}
}

func (l *Logger) CheckErr(err error) {
	if err != nil {
		l.Fatal(err)
	}
}
