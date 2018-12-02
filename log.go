package fiorm

import "github.com/sirupsen/logrus"

// GoLogWriter 自定义LOG
type GoLogWriter struct {
}

// Println 写日志
func (t *GoLogWriter) Println(v ...interface{}) {
	logrus.Debug(v...)
}
