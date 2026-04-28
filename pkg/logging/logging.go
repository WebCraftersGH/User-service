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

	l.SetOutput(os.Stdout)

	return logrus.NewEntry(l), nopCloser{}, nil
}

type nopCloser struct{}

func (nopCloser) Close() error {
	return nil
}

func parseLevel(level string) logrus.Level {
	parsed, err := logrus.ParseLevel(strings.ToLower(strings.TrimSpace(level)))
	if err != nil {
		return logrus.InfoLevel
	}
	return parsed
}
