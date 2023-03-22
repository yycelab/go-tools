package dt

import (
	"fmt"
	"runtime"
	"time"
)

const (
	STD_TIME_PATTERN  = "2006-01-02 15:04:05"
	LONG_TIME_PATTERN = "2006-01-02 15:04:05.999"
	STD_DATE          = "2006-01-02"
	DEFAULT_DATE_TIME = "1900-01-01 00:00:00.000"
	TIME_ZONE         = "Asia/Shanghai"
)

type Time time.Time

func (t *Time) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte(DEFAULT_DATE_TIME), nil
	}
	return []byte(time.Time(*t).Format(fmt.Sprintf(`"%s"`, STD_TIME_PATTERN))), nil
}

func (t *Time) UnmarshalJSON(datetime []byte) (err error) {
	defer func(err *error) {
		someErr := recover()
		if someErr != nil {
			switch someErr.(type) {
			case runtime.Error: // 运行时错误
				*err = fmt.Errorf("runtime error:%+v", err)
			default: // 非运行时错误
				*err = fmt.Errorf("error:%+v", err)
			}
		}
	}(&err)
	pattern := STD_TIME_PATTERN
	in := string(datetime[1 : len(datetime)-1])
	inputLen := len(in)
	if inputLen == 0 {
		return
	}
	if inputLen == len(STD_DATE) {
		pattern = STD_DATE
	} else if inputLen == len(LONG_TIME_PATTERN) {
		pattern = LONG_TIME_PATTERN
	}
	// println("json unmarshal pattern", pattern, " ", fmt.Sprintf("[%s]", in), " ", len(in))
	pt, err := time.ParseInLocation(pattern, in, time.FixedZone(TIME_ZONE, 0))
	// println("parse date succ:", fmt.Sprintf("%+v", pt))
	if err == nil {
		tmp := Time(pt)
		*t = tmp
	}
	return err
}

func (t *Time) Scan(src any) error {
	switch typ := src.(type) {
	case time.Time:
		*t = Time(src.(time.Time))
	case *time.Time:
		tmp := src.(*time.Time)
		if tmp != nil {
			v := *tmp
			*t = Time(v)
		}
	default:
		return fmt.Errorf("not support Scan type:%+v", typ)
	}
	return nil
}
