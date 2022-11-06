package main

import (
	"log"
	"os"
)

var infoLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.LUTC|log.Ltime)
var errorLogger = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile)

type app struct {
	store Store
}

func (a *app) logInfo(msg string) {
	infoLogger.Println(msg)
}

func (a *app) logFatal(msg string) {
	errorLogger.Println(msg)
}
