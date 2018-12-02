package db2struct

import (
	"fmt"
	"io"
	"os"
	"strings"

	"git.code.oa.com/fip-team/fiorm/internal"
	"github.com/iancoleman/strcase"
)

// db2struct --host 10.123.17.81 -d go_testdb -t user --package myGoPackage --struct User -p itcloud@123 --user root --guregu --gorm
// os.Getenv("MYSQL_HOST")
var mariadbHost = ""

// goopt.String([]string{"-H", "--host"}, "", "Host to check mariadb status of")
// var mariadbHostPassed = ""

// goopt.Int([]string{"--mysql_port"}, 3306, "Specify a port to connect to")
var mariadbPort = 3306

// goopt.String([]string{"-t", "--table"}, "", "Table to build struct from")
var mariadbTable = ""

// goopt.String([]string{"-d", "--database"}, "nil", "Database to for connection")
var mariadbDatabase = ""

var mariadbPassword = ""

// goopt.String([]string{"-u", "--user"}, "user", "user to connect to database")
var mariadbUser = ""

// goopt.Flag([]string{"-v", "--verbose"}, []string{}, "Enable verbose output", "")
// var verbose = ""

// goopt.String([]string{"--package"}, "", "name to set for package")
var packageName = ""

// goopt.String([]string{"--struct"}, "", "name to set for struct")
var structName = ""

// goopt.Flag([]string{"--json"}, []string{"--no-json"}, "Add json annotations (default)", "Disable json annotations")
var jsonAnnotation = true

// goopt.Flag([]string{"--gorm"}, []string{}, "Add gorm annotations (tags)", "")
var gormAnnotation = true

// goopt.Flag([]string{"--guregu"}, []string{}, "Add guregu null types", "")
var gureguTypes = true

// BuildFullParam 将表结构转换为实体
func BuildFullParam(table string, packageName string, jsonFlag bool, gureguFlag bool, filepath string) {
	s := &internal.DbSetting
	mariadbHost = s.Host
	mariadbPort = s.Port
	mariadbTable = table
	mariadbDatabase = s.DbName
	pwd, err1 := internal.Decrypt(s.Password)
	if err1 != nil {
		panic("密码解密错误")
	}
	mariadbPassword = pwd
	mariadbUser = s.User
	jsonAnnotation = jsonFlag
	gureguTypes = gureguFlag

	// Username is required
	if mariadbUser == "" {
		fmt.Println("Username is required! Add it with --user=name")
		return
	}

	// If a mariadb host is passed use it
	if mariadbHost == "" {
		fmt.Println("host is required! Add it with --host=host")
		return
	}

	if mariadbPassword == "" {
		fmt.Println("Error reading password:")
		return
	}

	if mariadbDatabase == "" {
		fmt.Println("Database can not be null")
		return
	}

	if mariadbTable == "" {
		fmt.Println("Table can not be null")
		return
	}

	// If structName is not set we need to default it
	if structName == "" {
		structName = strcase.ToCamel(mariadbTable)
	}
	// If packageName is not set we need to default it
	if packageName == "" {
		packageName = "newpackage"
	}

	columnDataTypes, err := GetColumnsFromMysqlTable(mariadbUser, mariadbPassword, mariadbHost, mariadbPort, mariadbDatabase, mariadbTable)

	if err != nil {
		fmt.Println("Error in selecting column data information from mysql information schema")
		return
	}

	// Generate struct string based on columnDataTypes
	struc, err := Generate(*columnDataTypes, mariadbTable, structName, packageName, jsonAnnotation, gormAnnotation, gureguTypes)

	if err != nil {
		fmt.Println("Error in creating struct from json: " + err.Error())
		return
	}

	output := string(struc)
	fmt.Println(output)

	// 生成文件
	path := filepath
	ferr := writeStringToFile(path, table, output)
	if ferr != nil {
		fmt.Println(ferr)
	}
}

// Build 将表结构转换为实体
func Build(table string, packageName string) {
	BuildFullParam(table, packageName, false, true, "c:/FiormModel")
}

func writeStringToFile(filepath, table, s string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		os.Mkdir(filepath, 0700)
	}
	path := filepath + "/" + table + ".go"
	fo, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}
	fmt.Println("实体类生成路径：", filepath)
	return nil
}
