package example_test

import (
	"os"
	"testing"

	"git.code.oa.com/fip-team/fiorm"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	pwd, err := fiorm.Encrypt("aaa")
	if err != nil {
		panic(err)
	}

	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	logrus.SetFormatter(&logrus.JSONFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	logrus.SetOutput(os.Stdout)
	//设置最低loglevel
	logrus.SetLevel(logrus.DebugLevel)
	conf := &fiorm.DbSettings{}
	conf.Dialect = "mysql"
	conf.DbName = "go_testdb"
	conf.Host = ""
	conf.User = "root"
	conf.Password = pwd
	conf.Port = 3306
	fiorm.InitDB(conf)
	m.Run()
}
