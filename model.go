package potato

import (
    "fmt"
    "strings"
    "database/sql"
    _"github.com/go-sql-driver/mysql"
)

type Model struct {
    Table string
    Columns []string
}


func (m *Model) SqlColumnsPart(Columns []string) string {
    return "`" + strings.Join(Columns, "`,`") + "`"
}

func (m *Model) SqlWherePart(query map[string]interface{}) (string, []interface{}) {
    var where string
    length := len(query)
    values := make([]interface{}, 0, length)

    if length > 0 {
        conditions := make([]string, 0, length)
        for k, v := range query {
            i := strings.Split(k, " ")
            c := i[0]
            o := "="
            if len(i) > 1 {
                o = i[1]
            }

            conditions = append(conditions, fmt.Sprintf("`%s` %s ?", c, o))
            values = append(values, v)
        }

        where = "WHERE " + strings.Join(conditions, " AND ")
    }
    
    return where, values
}

func (m *Model) SqlLimitPart(limit ...int64) string {
    length := len(limit)

    if length == 0 {
        return ""
    }
    
    if length == 1 {
        return fmt.Sprintf("LIMIT %d", limit[0])
    }

    return fmt.Sprintf("LIMIT %d, %d", limit[0], limit[1])
}

func (m *Model) CreateFindStmt(query map[string]interface{}, order string, limit ...int64) (string, []interface{}) {
    where, values := m.SqlWherePart(query)
    stmt := fmt.Sprintf("SELECT %s FROM `%s` %s ORDER BY %s %s",
            m.SqlColumnsPart(m.Columns), m.Table, where, order, m.SqlLimitPart(limit...))

    return stmt, values
}

func (m *Model) Find(query map[string]interface{}, order string, limit ...int64) (*sql.Rows, error) {
    stmt, values := m.CreateFindStmt(query, order, limit...)
    rows, e := D.Query(stmt, values...)
    return rows, e
}

func (m *Model) FindOne(query map[string]interface{}, order string) *sql.Row {
    stmt, values := m.CreateFindStmt(query, order, 1)
    return D.QueryRow(stmt, values...)
}

func (m *Model) Insert(data map[string]interface{}) int64 {
    length  := len(data)
    columns := make([]string, 0, length)
    holders := make([]string, 0, length)
    values  := make([]interface{}, 0, length)
    for k, v := range data {
        columns = append(columns, k)
        holders = append(holders, "?")
        values  = append(values, v)
    }

    stmt := fmt.Sprintf("INSERT INTO `%s` (%s)VALUES(%s)",
            m.Table, m.SqlColumnsPart(columns), strings.Join(holders, ","))

    result, e := D.Exec(stmt, values...)
    if e != nil {
        L.Println(e)
        return 0
    }

    id, e := result.LastInsertId()
    if e!= nil {
        L.Println(e)
        return 0
    }

    return id
}

func (m *Model) Update(data map[string]interface{}, query map[string]interface{}, limit ...int64) int64 {
    length := len(data)
    sets   := make([]string, 0, length)
    values := make([]interface{}, 0, length + len(query))
    for k, v := range data {
        sets   = append(sets, fmt.Sprintf("`%s`=?", k))
        values = append(values, v)
    }

    where, v := m.SqlWherePart(query)
    values = append(values, v...)
    stmt := fmt.Sprintf("UPDATE `%s` SET %s %s %s",
            m.Table, strings.Join(sets, ","), where, m.SqlLimitPart(limit...))

    result, e := D.Exec(stmt, values...)
    if e != nil {
        L.Println(e)
        return 0
    }

    n, e := result.RowsAffected()
    if e!= nil {
        L.Println(e)
        return 0
    }

    return n
}

func (m *Model) Delete(query map[string]interface{}, limit ...int64) int64 {
    where, values := m.SqlWherePart(query)
    stmt := fmt.Sprintf("DELETE `%s` %s %s", m.Table, where, m.SqlLimitPart(limit...))
    result, e := D.Exec(stmt, values...)
    if e != nil {
        L.Println(e)
        return 0
    }

    n, e := result.RowsAffected()
    if e!= nil {
        L.Println(e)
        return 0
    }

    return n
}