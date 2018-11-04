package logging

import (
	"strings"
	"bytes"
	"sync"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

type logger struct{
	logFile string
	buf bytes.Buffer

	output io.Writer
	sync.Mutex
}

//daily,hourly,singly
func New(logFile string, args ...string) *logger {
	if logFile != "" && len(args) > 1 {
		rotate := args[0]
		switch rotate {
		case "hourly":
			logFile += time.Now().Format("2006010215")
		case "daily":
			logFile += time.Now().Format("20060102")
		}
	}

	return &logger{
		logFile: logFile,
	}
}

func (l *logger) Debugf(tpl string, args ...interface{}) {
	l.level(Debug).write(tpl, args)
}

func (l *logger) Infof(tpl string, args ...interface{}) {
	l.level(Info).write(tpl, args)
}

func (l *logger) Warnf(tpl string, args ...interface{}) {
	l.level(Warn).write(tpl, args)
}

func (l *logger) Errorf(tpl string, args ...interface{}) {
	l.level(Error).write(tpl, args)
}

func (l *logger) Panicf(tpl string, args ...interface{}) {
	l.level(Panic).write(tpl, args)
}

func (l *logger) Fatalf(tpl string, args ...interface{}) {
	l.level(Fatal).write(tpl, args)
}

func (l *logger) level(lvl Level) *logger{
	curTime := time.Now().Format("2006-01-02 15:04:05.000")
	l.buf.Write([]byte(curTime+"\t"))
	l.printColored(lvl)
	return l
}

func (l *logger) write(tpl string, args []interface{}) {
	msg := tpl
	if len(args) > 0{
		if msg == "" {
			msg = fmt.Sprint(args...)
		} else {
			msg = fmt.Sprintf(tpl, args...)
		}
	}
	l.buf.Write([]byte(msg))

	if l.logFile != "" && l.logFile != string(os.PathSeparator) {
		dir := path.Dir(l.logFile)
		if dir != "." && dir != string(os.PathSeparator) {
			if _,err:=os.Stat(dir);os.IsNotExist(err) {
				err = os.MkdirAll(dir, 0755)
				if err != nil{
					log.Fatalf("make dir: %s failed, err:%v", dir, err)
					os.Exit(1)
				}
			}
		}
	}

	if l.logFile != "" && path.Base(l.logFile) != "" {
		file, err := os.OpenFile(l.logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		defer file.Close()
		if err!=nil{
			log.Fatalf("open file err: %v", err)
			os.Exit(1)
		} else {
			l.output = file
		}
	} else {
		l.output = os.Stdout
	}

	l.Lock()
	_,err := l.output.Write(l.buf.Bytes())
	if err != nil {
		log.Fatalf("write log err: %v", err)
	}
	l.output.Write([]byte("\n"))
	l.buf.Reset()
	l.Unlock()
}

func (l *logger) printColored(lvl Level) *logger {
	var levelColor int
	switch lvl {
	case Debug:
		levelColor = gray
	case Warn:
		levelColor = yellow
	case Error, Fatal, Panic:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(lvl.String())
	if l.logFile == "" {
		coloredLevel := fmt.Sprintf("\x1b[%dm%s\x1b[0m\t", levelColor, levelText)
		l.buf.Write([]byte(coloredLevel))
	} else {
		l.buf.Write([]byte(levelText+"\t"))
	}

	return l
}