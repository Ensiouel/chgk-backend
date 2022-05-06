package gocket

import (
	"log"
	"os"
)

type Blogger struct {
	Err  *log.Logger
	Info *log.Logger
}

func NewBlogger() *Blogger {
	infoFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(infoFile, "", log.Ldate|log.Ltime)

	errFile, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	errLog := log.New(errFile, "", log.Ldate|log.Ltime)

	return &Blogger{
		Info: infoLog,
		Err:  errLog,
	}
}
