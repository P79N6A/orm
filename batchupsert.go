package fiorm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

// batchUpdate 批量更新
func (t *FiDB) batchUpdate(rows interface{}, cols []string) {
	t.upsert(rows, cols)
}

// batchInsert 批量插入
func (t *FiDB) batchInsert(rows interface{}) {
	t.upsert(rows, nil)
}

func (t *FiDB) upsert(rows interface{}, cols []string) {
	defer func() {
		if r := recover(); r != nil {
			t.db.Error = fmt.Errorf("%v", r)
		}
	}()

	kind := reflect.TypeOf(rows).Kind()
	if kind != reflect.Slice {
		t.db.Error = errors.New("只支持 []T 或者 []*T")
		return
	}

	// 获取构造体
	raws := reflect.ValueOf(rows)
	if raws.Len() == 0 {
		fmt.Println("没有发现数据行")
		return
	}

	first := raws.Index(0).Interface()
	var st = new(gorm.ModelStruct)
	st, tableName := getTableStruct(first)

	// 判断更新是是否包含ID为0
	if cols != nil {
		val := reflect.ValueOf(raws.Index(0).Interface())
		value := val.FieldByName("ID").Interface()
		switch v := value.(type) {
		case int:
		case int8:
		case int16:
		case int32:
		case int64:
		case uint:
		case uint8:
		case uint16:
		case uint32:
		case uint64:
			if v == 0 {
				t.Error = errors.New("批量更新时，ID不能为0")
				return
			}
		default:
			msg := fmt.Sprintf("未知的ID类型,%v", reflect.TypeOf(value))
			t.Error = errors.New(msg)
			return
		}

	}

	// 分页处理，每次插入500条
	total := raws.Len()
	start := 0
	end := 500
	stop := false
	var values []interface{}
	var sqlStr string

	for j := 0; j < total; j++ {
		if stop == false && start != total {
			sqlStr, values, stop = generateSQL(st, tableName, rows, cols, start, end)
			d := t.db.Exec(sqlStr, values...)
			if d.Error != nil {
				t.db.Error = d.Error
			}

			start = end
			end = end + 500
		} else {
			break
		}
	}
}

func generateSQL(st *gorm.ModelStruct, tableName string, rows interface{}, cols []string, start int, end int) (string, []interface{}, bool) {
	var (
		column     string
		needColumn = true
		values     = []interface{}{}
		raw        = reflect.ValueOf(rows)
		sql        = "INSERT INTO %s ( %s ) VALUES "
	)

	fieldInfo := st.StructFields
	total := raw.Len()
	stop := false
	if end > total {
		end = total
		stop = true
	}

	for i := start; i < end; i++ {
		val := reflect.ValueOf(raw.Index(i).Interface())
		for val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		line := "("
		tp := reflect.Indirect(val).Type()
		for i := 0; i < val.NumField(); i++ {
			c := getFiled(fieldInfo, tp.Field(i).Name)
			if c == "" {
				continue
			}
			if needColumn {
				column += c + ","
			}
			line += "?,"
			values = append(values, val.Field(i).Interface())
		}
		line = strings.TrimSuffix(line, ",")
		line += "),"
		sql += line
		needColumn = false
	}

	sql = strings.TrimSuffix(sql, ",")
	column = strings.TrimSuffix(column, ",")

	// 构造更新语句
	if cols != nil {
		updateSQL := " on duplicate key update "
		for _, val := range cols {
			updateSQL += val + "=values(" + val + "),"
		}
		updateSQL = strings.TrimSuffix(updateSQL, ",")
		sql += updateSQL
	}

	return fmt.Sprintf(sql, tableName, column), values, stop
}

func getTableStruct(row interface{}) (*gorm.ModelStruct, string) {
	scope := db.NewScope(row)
	return scope.GetModelStruct(), scope.QuotedTableName()
}

func getFiled(fields []*gorm.StructField, name string) string {
	for _, f := range fields {
		if name == f.Name {
			if f.IsIgnored {
				return ""
			}
			return f.DBName
		}
	}
	return ""
}
