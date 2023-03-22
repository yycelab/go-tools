package io

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	ErrInvalidConfigFileType = errors.New("unsupport config file's type")
)

func ReadConfigByViper(file string) (*viper.Viper, error) {
	v := viper.New()
	pos := strings.LastIndex(file, ".")
	if pos < 0 {
		return nil, ErrInvalidConfigFileType
	}
	v.SetConfigType(file[pos+1:])
	f, err := os.OpenFile(file, os.O_RDONLY, 0644)
	if err == nil {
		defer f.Close()
		err = v.ReadConfig(f)
	}
	return v, err
}

func DirEnsure(dest string, create bool) (abspath string, created bool, err error) {
	tmp := filepath.Dir(dest)
	if !filepath.IsAbs(tmp) {
		tmp, err = filepath.Abs(tmp)
		if err != nil {
			tmp, err = filepath.Abs(fmt.Sprintf("./%s", tmp))
		}
		if err != nil {
			return
		}
	}

	if fdir, e := os.Stat(tmp); e != nil {
		err = e
		if create && os.IsNotExist(err) {
			created = true
			err = os.MkdirAll(tmp, 0777)
		}
	} else if !fdir.IsDir() {
		err = fmt.Errorf("%s不是有效的目录", tmp)
	}

	if err == nil {
		if strings.HasSuffix(dest, "/") {
			abspath = fmt.Sprintf("%s/", tmp)
		} else {
			abspath = fmt.Sprintf("%s/%s", tmp, filepath.Base(dest))
		}
	}
	return
}
