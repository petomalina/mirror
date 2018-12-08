package mirror

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

// Logger is an alias to the logrus logger that provides additional
// methods for bundle manipulation and reflection
type Logger struct {
	*logrus.Logger
}

type LogFactory func() *logrus.Entry

// Fields type-aliases the logrus.Fields so the package can be skipped within
// the mirror package
type Fields = logrus.Fields

var (
	// L is a global logger that can be reconfigured by third parties
	// to customize logging
	L = &Logger{logrus.New()}
)

// init instruments third party libraries to work in default
// settings when running mirror code or its bundling extension
func init() {
	L.SetLevel(logrus.InfoLevel)
	L.SetOutput(os.Stdout)
	L.SetFormatter(&logrus.TextFormatter{})
}

func (l *Logger) Method(obj, method string) *logrus.Entry {
	_, callerFile, callerLine, _ := runtime.Caller(0)

	return l.WithFields(logrus.Fields{
		"Object": obj,
		"Method": method,
		"Caller": fmt.Sprintf("%s:%d", callerFile, callerLine),
	})
}
