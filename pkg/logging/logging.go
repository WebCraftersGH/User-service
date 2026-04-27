package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	LogsPath     = "logs"
	LogsFilename = "log.log"
)

type Logger = *logrus.Entry

func New(level string) (Logger, io.Closer, error) {
	l := logrus.New()
	l.SetReportCaller(true)
	l.SetLevel(parseLevel(level))
	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			funcName := filepath.Base(f.Function)
			fileName := filepath.Base(f.File)
			return funcName + "()", fmt.Sprintf("%s:%d", fileName, f.Line)
		},
	})

	if err := os.MkdirAll(LogsPath, 0o755); err != nil {
		return nil, nil, fmt.Errorf("create logs dir: %w", err)
	}

	logPath := filepath.Join(LogsPath, LogsFilename)
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o640)
	if err != nil {
		return nil, nil, fmt.Errorf("open log file: %w", err)
	}

	l.SetOutput(io.MultiWriter(os.Stdout, logFile))

	return logrus.NewEntry(l), logFile, nil
}

func parseLevel(level string) logrus.Level {
	parsed, err := logrus.ParseLevel(strings.ToLower(strings.TrimSpace(level)))
	if err != nil {
		return logrus.InfoLevel
	}
	return parsed
}
