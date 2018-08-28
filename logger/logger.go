package logger

import (
    "io"
    "log"
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

func Init(
    sentryURL string,
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer,
    fatalHandle io.Writer) {

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