package sql

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/xwb1989/sqlparser"
)

var (
	restore_param = regexp.MustCompile(`:v\d{1,2}`)
)

// func DynamicWhere(like any) (where string, params []any) {

// }

func ParseSQL(exp string) (pretty string, err error) {
	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(exp)
	if err == nil {
		// w := new(bytes.Buffer)

		buff := sqlparser.NewTrackedBuffer(func(buf *sqlparser.TrackedBuffer, node sqlparser.SQLNode) {
			rewrite := ""
			switch t := node.(type) {
			case *sqlparser.ColName:
				rewrite = sqlparser.String(t)
			case sqlparser.TableName:
				rewrite = sqlparser.String(t)
			case *sqlparser.SQLVal:
				if strings.HasPrefix(sqlparser.String(t), ":") {
					t.Val = []byte{'?'}
				}
				// case *sqlparser.AndExpr:
				// 	right := sqlparser.String(t.Right)
				// 	println("right:", right)
			}
			if len(rewrite) > 0 {
				if !strings.Contains(rewrite, "`") {
					rewrite = fmt.Sprintf("`%s`", rewrite)
				}
				buf.WriteString(rewrite)
			} else {
				node.Format(buf)
			}
		})
		stmt.Format(buff)
		pretty = buff.String()
	}
	return
}

func CountSQL(selectSQL string) (count string, pn int, err error) {
	var stmt sqlparser.Statement
	stmt, err = sqlparser.Parse(selectSQL)
	if err == nil {
		query := stmt.(*sqlparser.Select)
		from := sqlparser.String(query.From)
		if query.Where != nil {
			where := restore_param.ReplaceAllString(sqlparser.String(query.Where), "?")
			pn = strings.Count(where, "?")
			count = fmt.Sprintf("SELECT COUNT(1) FROM %s %s", from, where)
		} else {
			count = fmt.Sprintf("SELECT COUNT(1) FROM %s", from)
		}
	}
	return
}
