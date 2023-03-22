package param

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"yycelab.com/go-tools/errors"
)

func AssertParamError(err error, msg string, args ...any) {
	if err != nil {
		panic(&errors.ParamError{Msg: fmt.Sprintf(msg, args...)})
	}
}

func Json(req *http.Request, form any) {
	body := bytes.NewBufferString("")
	buff := bufio.NewWriter(body)
	_, err := io.Copy(buff, req.Body)
	AssertParamError(err, "empty body")
	buff.Flush()
	AssertParamError(json.Unmarshal(body.Bytes(), form), "invalid json")
}
