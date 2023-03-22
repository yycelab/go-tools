package dt

import "time"

const (
	Standard_Long   = "2006-01-02 15:04:05.999"
	Standard_YMDHMS = "2006-01-02 15:04:05"
	Standard_YMD    = "2006-01-02"
)

func Format(t time.Time, pattern string) string {
	return t.Format(pattern)
}

func Parser(tstring string, pattern string) (time.Time, error) {
	return time.Parse(pattern, tstring)
}
