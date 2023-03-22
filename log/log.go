package log

import (
	"io"
	stlog "log"
	"os"
)

const p string = "Downloader: "

type Log struct {
	printlnConsole *stlog.Logger
	printFile      *stlog.Logger
	printMulti     *stlog.Logger
	prefix         string "Downloader: "
}

func New() *Log {
	l := &Log{}
	l.printlnConsole = stlog.New(os.Stdout, p, stlog.LstdFlags|stlog.Lshortfile)
	file, err := os.OpenFile(".\\log\\log.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		stlog.Println(err)
	}
	//defer file.Close()
	l.printFile = stlog.New(file, p, stlog.LstdFlags|stlog.Lshortfile)
	multi := io.MultiWriter(os.Stdout, file)
	l.printMulti = stlog.New(multi, p, stlog.LstdFlags|stlog.Lshortfile)
	return l
}

func (l *Log) PrintlnConsole(v ...any) {
	l.printlnConsole.Println(v...)
}

func (l *Log) PrintFile(v ...any) {
	l.printFile.Println(v...)
}

func (l *Log) PrintMulti(v ...any) {
	l.printMulti.Println(v...)
}
