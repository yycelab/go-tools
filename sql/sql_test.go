package sql

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type Demo struct {
	Age    int
	Name   string
	Gender string
	Dt     time.Time
}

type BDemo struct {
	*Demo
}

func TestWhereDynamic(t *testing.T) {
	var c any = &Demo{
		Age:    1,
		Name:   "",
		Gender: "men",
		Dt:     time.Now(),
	}
	r := reflect.ValueOf(c)
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	max := r.NumField()
	index := 0
	for index < max {
		field := r.Field(index)
		// sf,ok := field.(*reflect.StructField)
		if !field.IsValid() || field.IsZero() {
			println(index, "invalid value")
		} else {
			println(index, "valid value", fmt.Sprintf("%+v", field.Interface()))
		}
		index++
	}
	rt := reflect.TypeOf(c)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	fnum := rt.NumField()
	j := 0
	for j < fnum {
		f := rt.Field(j)
		println(j, " name:", f.Name)
		j++
	}

}

func TestParseSQL(t *testing.T) {
	sql := "SELECT id as `Id`,age as `Age`,ELAB_USER as `User`,`FROM` as `From`  FROM   ELAB_USER WHERE id=? and age=? and User='123'"
	txt, err := ParseSQL(sql)
	if err == nil {
		println("returned:", txt)
	} else {
		println(err.Error())
	}
}

func TestExtractWhere(t *testing.T) {
	var dt time.Time
	println(fmt.Sprintf("%+v", dt))
	println(dt.Unix())
	dt = time.Now()
	println(dt.Unix())

	sql := "SELECT ELAB_USER,`FROM`   FROM   ELAB_USER WHERE id=? and age=?"
	w, n, ok := CountSQL(sql)
	println("found:", ok, " ,where:", w, " ", n)
}

// func TestParseSQL(t *testing.T) {
// 	// table := "ELAB_USER"
// 	sql := "SELECT ELAB_USER,`FROM`   FROM   ELAB_USER where id=? and name=?"

// 	stmt, err := sqlparser.Parse(sql)

// 	if err == nil {
// 		query := stmt.(*sqlparser.Select)
// 		node := query.Where
// 		where := sqlparser.String(node)
// 		pattern := regexp.MustCompile(`:v\d{1,2}`)
// 		replaced := pattern.ReplaceAll([]byte(where), []byte("?"))
// 		println(string(replaced))
// 		return
// 	}
// 	println(err.Error())
// }
