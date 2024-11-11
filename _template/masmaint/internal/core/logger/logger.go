package logger

import (
	"log"
	"io"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"

	"masmaint/config"
)


const LOG_LEVEL_DEBUG = 1
const LOG_LEVEL_INFO = 2
const LOG_LEVEL_WARNING = 3
const LOG_LEVEL_ERROR = 4
const LOG_LEVEL_FATAL = 5

const LOGFOLDER = "log"
const LOGFILE = "app.log"

var file *os.File

var logD *log.Logger
var logI *log.Logger
var logW *log.Logger
var logE *log.Logger
var logF *log.Logger

var logLevel int


func init() {
	_, err := os.Stat(LOGFOLDER)
	if err != nil {
		os.Mkdir(LOGFOLDER, 0755)
	}

	fp := LOGFOLDER + "/" + LOGFILE
	_, err = os.Stat(fp)
	if err != nil {
		file, err = os.Create(fp)
	} else {
		file, err = os.OpenFile(fp, os.O_APPEND|os.O_WRONLY, 0644)
	}

	if err != nil {
		log.Panic(err)
	}

	logLevel = getLogLevel()

	logD = log.New(
		os.Stdout,
		"[DEBUG]",
		log.LstdFlags,
	)

	logI = log.New(
		io.MultiWriter(os.Stdout, file),
		"[INFO]",
		log.LstdFlags,
	)

	logW = log.New(
		io.MultiWriter(os.Stdout, file),
		"[WARNING]",
		log.LstdFlags,
	)

	logE = log.New(
		io.MultiWriter(os.Stdout, file),
		"[ERROR]",
		log.LstdFlags,
	)

	logF = log.New(
		io.MultiWriter(os.Stdout, file),
		"[FATAL]",
		log.LstdFlags,
	)

	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
}


func getLogLevel() int {
	level := config.GetConfig().LogLevel
	switch level {
	case "DEBUG", "debug":
		return LOG_LEVEL_DEBUG
	case "INFO", "info":
		return LOG_LEVEL_INFO
	case "WARNGING", "warning":
		return LOG_LEVEL_WARNING
	case "ERROR", "error":
		return LOG_LEVEL_ERROR
	case "FATAL", "fatal":
		return LOG_LEVEL_FATAL
	default:
		return LOG_LEVEL_INFO
	}
}


func Debug(msg string) {
	if logLevel > LOG_LEVEL_DEBUG {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logD.Println("\n", "File:", file, "Line:", line, "\n",
		"Func:", f.Name(), "Msg:", msg,
	)
}


func Info(msg string) {
	if logLevel > LOG_LEVEL_INFO {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logI.Println(f.Name(), msg)
}


func Warning(msg string) {
	if logLevel > LOG_LEVEL_WARNING {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logW.Println("\n", "File:", file, "Line:", line, "\n",
		"Func:", f.Name(), "Msg:", msg,
	)
}


func Error(msg string) {
	if logLevel > LOG_LEVEL_ERROR {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logE.Println("\n", "File:", file, "Line:", line, "\n",
		"Func:", f.Name(), "Msg:", msg,
	)
}


func Fatal(msg string) {
	if logLevel > LOG_LEVEL_FATAL {
		return
	}
	logF.Fatal("Msg:", msg)
}
