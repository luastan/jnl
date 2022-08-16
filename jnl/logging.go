package jnl

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "[i] ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(os.Stderr, "[X] ", log.Ldate|log.Ltime)
}
