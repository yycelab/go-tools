package str

import (
	"strings"
	"testing"
)

func TestFirtCaptial(t *testing.T) {
	s := []string{"age_name", "He_is", "myName", "Test_name", "_name"}
	for i := range s {
		println(s[i], "=>", WordsFirstLetter(s[i], L_Upper|L_Retain|L_After_Undlerline))
	}
}

func TestSplitN(t *testing.T) {
	table := "ELAB_USER"
	sql := "SELECT ELAB_USER   FROM   ELAB_USER"
	sql = strings.ReplaceAll(sql, "\\s+", " ")
	println("sql:", sql)
	strings.Count(sql, "WHERE")
	name := strings.SplitAfterN(sql, table, -1)
	println(len(name))
	if len(name) > 1 {
		println("where:", name[1])
	} else {
		println(name[0])
	}
}
