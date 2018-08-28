package logger

import (
    "io"
    "os"
    "log"
    "path/filepath"
    "github.com/getsentry/raven-go"
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
    Fatal   *log.Logger
    Raven   *raven.Client
)

func Init(logDir, sentryURL string) {
    if _, err := os.Stat(logDir); os.IsNotExist(err) {
        err := os.MkdirAll(logDir, 0777)
        if err != nil {
            panic(err)
        }
    }
    errLogFile, err := os.OpenFile(filepath.Join(logDir, "error.log"), os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        panic(err)
    }
    infoLogFile, err := os.OpenFile(filepath.Join(logDir, "info.log"), os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
    if err != nil {
        panic(err)
    }
    traceHandle := os.Stdout
    infoHandle := io.MultiWriter(os.Stdout, infoLogFile)
    warningHandle := io.MultiWriter(os.Stdout, infoLogFile)
    errorHandle := io.MultiWriter(os.Stderr, errLogFile)
    fatalHandle := io.MultiWriter(os.Stderr, errLogFile)

    Raven, _ = raven.New(sentryURL)
    // Raven, _ = raven.New("")
    // raven.SetDSN()

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
    Fatal = log.New(fatalHandle,
        "FATAL: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}