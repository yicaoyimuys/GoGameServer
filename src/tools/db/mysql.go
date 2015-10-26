package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strconv"
	"strings"
)

import (
	. "tools"
	"tools/cfg"
)

type Model struct {
	db        *sql.DB
	tablename string
	param     []string
	columnstr string
	where     string
	pk        string
	orderby   string
	limit     string
	join      string
}

var DBOrm *Model

//数据库初始化
func init() {
	hostname := cfg.GetValue("db_hostname")
	port := cfg.GetValue("db_port")
	username := cfg.GetValue("db_username")
	password := cfg.GetValue("db_password")
	database := cfg.GetValue("db_database")
	charset := cfg.GetValue("db_charset")

	sqlDb, err := sql.Open("mysql", username+":"+password+"@tcp("+hostname+":"+port+")/"+database+"?charset="+charset)
	if err != nil {
		ERR("Sql Open Error", err)
	}

	sqlDb.SetMaxOpenConns(5000)
	sqlDb.SetMaxIdleConns(1000)
	err = sqlDb.Ping()
	if err != nil {
		ERR("Sql Ping Error", err)
	}

	DBOrm = new(Model)
	DBOrm.db = sqlDb
}

func (m *Model) FindAll() map[int]map[string]string {
	if m.db == nil {
		ERR("[FindAll] mysql not connect")
		return nil
	}

	if len(m.param) == 0 {
		m.columnstr = "*"
	} else {
		if len(m.param) == 1 {
			m.columnstr = m.param[0]
		} else {
			m.columnstr = strings.Join(m.param, ",")
		}
	}

	sql := fmt.Sprintf("Select %v from %v %v %v %v %v", m.columnstr, m.tablename, m.join, m.where, m.orderby, m.limit)
	rows, err := m.db.Query(sql)
	if err != nil {
		ERR(sql, err)
		return nil
	}
	return QueryResult(rows)
}

func (m *Model) FindOne() map[int]map[string]string {
	if m.db == nil {
		ERR("[FindOne] mysql not connect")
		return nil
	}

	return m.Limit(1).FindAll()
}

func (m *Model) Insert(param map[string]interface{}) (num int, err error) {
	if m.db == nil {
		ERR("[Insert] mysql not connect")
		return 0, errors.New("mysql not connect")
	}
	var keys []string
	var values []string
	if len(m.pk) != 0 {
		delete(param, m.pk)
	}
	for key, value := range param {
		keys = append(keys, key)
		switch value.(type) {
		case int, int64, int32:
			values = append(values, strconv.Itoa(value.(int)))
		case string:
			values = append(values, value.(string))
		case float32, float64:
			values = append(values, strconv.FormatFloat(value.(float64), 'f', -1, 64))
		}
	}
	fileValue := "'" + strings.Join(values, "','") + "'"
	fileds := "`" + strings.Join(keys, "`,`") + "`"
	sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", m.tablename, fileds, fileValue)
	result, err := m.db.Exec(sql)
	if err != nil {
		ERR(sql, err)
		return 0, err
	}
	i, err := result.LastInsertId()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if err != nil {
		err = errors.New("insert failure")
	}
	return s, err

}

func (m *Model) Fileds(param ...string) *Model {
	m.param = param
	return m
}

func (m *Model) Update(param map[string]interface{}) (num int, err error) {
	if m.db == nil {
		ERR("[Update] mysql not connect")
		return 0, errors.New("mysql not connect")
	}
	var setValue []string
	for key, value := range param {
		switch value.(type) {
		case int, int64, int32:
			set := fmt.Sprintf("%v = %v", key, value.(int))
			setValue = append(setValue, set)
		case string:
			set := fmt.Sprintf("%v = '%v'", key, value.(string))
			setValue = append(setValue, set)
		case float32, float64:
			set := fmt.Sprintf("%v = '%v'", key, strconv.FormatFloat(value.(float64), 'f', -1, 64))
			setValue = append(setValue, set)
		}

	}
	setData := strings.Join(setValue, ",")
	sql := fmt.Sprintf("UPDATE %v SET %v %v", m.tablename, setData, m.where)
	result, err := m.db.Exec(sql)
	if err != nil {
		ERR(sql, err)
		return 0, err
	}
	i, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))

	return s, err
}

func (m *Model) Delete(param string) (num int, err error) {
	if m.db == nil {
		ERR("[Delete] mysql not connect")
		return 0, errors.New("mysql not connect")
	}
	h := m.Where(param).FindOne()
	if len(h) == 0 {
		return 0, errors.New("no Value")
	}
	sql := fmt.Sprintf("DELETE FROM %v WHERE %v", m.tablename, param)
	result, err := m.db.Exec(sql)
	if err != nil {
		ERR(sql, err)
		return 0, err
	}
	i, err := result.RowsAffected()
	s, _ := strconv.Atoi(strconv.FormatInt(i, 10))
	if i == 0 {
		err = errors.New("delete failure")
	}

	return s, err
}

func (m *Model) Query(sql string) interface{} {
	if m.db == nil {
		ERR("[Query] mysql not connect")
		return errors.New("mysql not connect")
	}
	var query = strings.TrimSpace(sql)
	s, err := regexp.MatchString(`(?i)^select`, query)
	if err == nil && s == true {
		result, _ := m.db.Query(sql)
		c := QueryResult(result)
		return c
	}
	exec, err := regexp.MatchString(`(?i)^(update|delete)`, query)
	if err == nil && exec == true {
		m_exec, err := m.db.Exec(query)
		if err != nil {
			return err
		}
		num, _ := m_exec.RowsAffected()
		id := strconv.FormatInt(num, 10)
		return id
	}

	insert, err := regexp.MatchString(`(?i)^insert`, query)
	if err == nil && insert == true {
		m_exec, err := m.db.Exec(query)
		if err != nil {
			return err
		}
		num, _ := m_exec.LastInsertId()
		id := strconv.FormatInt(num, 10)
		return id
	}
	result, _ := m.db.Exec(query)

	return result

}

func QueryResult(rows *sql.Rows) map[int]map[string]string {
	defer rows.Close()

	var result = make(map[int]map[string]string)
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanargs := make([]interface{}, len(values))
	for i := range values {
		scanargs[i] = &values[i]
	}

	var n = 1
	for rows.Next() {
		result[n] = make(map[string]string)
		err := rows.Scan(scanargs...)

		if err != nil {
			fmt.Println(err)
		}

		for i, v := range values {
			result[n][columns[i]] = string(v)
		}
		n++
	}
	return result
}

func (m *Model) SetTable(tablename string) *Model {
	m.tablename = tablename
	return m
}

func (m *Model) Where(param string) *Model {
	m.where = fmt.Sprintf(" where %v", param)
	return m
}

func (m *Model) SetPk(pk string) *Model {
	m.pk = pk
	return m
}

func (m *Model) OrderBy(param string) *Model {
	m.orderby = fmt.Sprintf("ORDER BY %v", param)
	return m
}

func (m *Model) Limit(size ...int) *Model {
	var end int
	start := size[0]
	if len(size) > 1 {
		end = size[1]
		m.limit = fmt.Sprintf("Limit %d,%d", start, end)
		return m
	}
	m.limit = fmt.Sprintf("Limit %d", start)
	return m
}

func (m *Model) LeftJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("LEFT JOIN %v ON %v", table, condition)
	return m
}

func (m *Model) RightJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("RIGHT JOIN %v ON %v", table, condition)
	return m
}

func (m *Model) Join(table, condition string) *Model {
	m.join = fmt.Sprintf("INNER JOIN %v ON %v", table, condition)
	return m
}

func (m *Model) FullJoin(table, condition string) *Model {
	m.join = fmt.Sprintf("FULL JOIN %v ON %v", table, condition)
	return m
}

//the function will use friendly way to print the data
func Print(slice map[int]map[string]string) {
	for _, v := range slice {
		for key, value := range v {
			fmt.Println(key, value)
		}
		fmt.Println("---------------")
	}
}

func (m *Model) DbClose() {
	m.db.Close()
}
