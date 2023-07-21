package logger

import (
	"log"
	"io"
	"os"
	"time"
	"runtime"

	"github.com/gin-gonic/gin"

	"masmaint/config"
)


const LOG_LEVEL_DEBUG = 1
const LOG_LEVEL_INFO = 2
const LOG_LEVEL_WARNING = 3
const LOG_LEVEL_ERROR = 4
const LOG_LEVEL_FATAL = 5

const LOGFOLDER = "log/"
const FORMAT = "2006-01-02-15-04-05"
const LOGFILE_EX = ".log"

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
		os.Mkdir(LOGFOLDER, 0777)
	}

	t := time.Now()
	fn := LOGFOLDER + t.Format(FORMAT) + LOGFILE_EX

	file, err = os.Create(fn);

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


func SetAccessLogger() {
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
}


func LogDebug(msg string) {
	if logLevel > LOG_LEVEL_DEBUG {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logD.Println(f.Name(), msg)
}


func LogInfo(msg string) {
	if logLevel > LOG_LEVEL_INFO {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logI.Println(f.Name(), msg)
}


func LogWarning(msg string) {
	if logLevel > LOG_LEVEL_WARNING {
		return
	}
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logW.Println(f.Name(), msg)
}


func LogError(msg string) {
	if logLevel > LOG_LEVEL_ERROR {
		return
	}
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logE.Println("\n", "File:", file, "Line:", line, "\n",
		"Func:", f.Name(), "Msg:", msg,
	)
}


func LogFatal(msg string) {
	if logLevel > LOG_LEVEL_FATAL {
		return
	}
	logF.Fatal("Msg:", msg)
}
