package logger

import (
	"log"
	"io"
	"os"
	"time"
	"runtime"

	"github.com/gin-gonic/gin"
)


const LOGFOLDER = "log/"
const FORMAT = "2006-01-02-15-04-05"
const LOGFILE_EX = ".log"

var file *os.File

var logF *log.Logger
var logE *log.Logger
var logW *log.Logger
var logI *log.Logger
var logD *log.Logger


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

	logF = log.New(
		io.MultiWriter(os.Stdout, file),
		"[FATAL]",
		log.LstdFlags,
	)

	logE = log.New(
		io.MultiWriter(os.Stdout, file),
		"[ERROR]",
		log.LstdFlags,
	)

	logW = log.New(
		io.MultiWriter(os.Stdout, file),
		"[WARN]",
		log.LstdFlags,
	)

	logI = log.New(
		io.MultiWriter(os.Stdout, file),
		"[INFO]",
		log.LstdFlags,
	)

	logD = log.New(
		os.Stdout,
		"[DEBUG]",
		log.LstdFlags,
	)

}

func SetAccessLogger() {
	gin.DefaultWriter = io.MultiWriter(os.Stdout, file)
}

func LogFatal(msg string) {
	logF.Fatal("Msg:", msg)
}


func LogError(msg string) {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logE.Println("\n", "File:", file, "Line:", line, "\n",
		"Func:", f.Name(), "Msg:", msg,
	)
}


func LogWarn(msg string) {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logW.Println(f.Name(), msg)
}


func LogInfo(msg string) {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logI.Println(f.Name(), msg)
}


func LogDebug(msg string) {
	pc, _, _, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	logD.Println(f.Name(), msg)
}