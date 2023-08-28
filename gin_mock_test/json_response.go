package gin_mock_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	//nolint:golint,unused
	errJsonErrNil = fmt.Errorf("err is nil")
)

//nolint:golint,unused
type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// jsonSuccess
// return
//
//nolint:golint,unused
func jsonSuccess(c *gin.Context, data interface{}) {
	if data != nil {
		c.JSON(http.StatusOK, response{
			Code: 0,
			Msg:  "success",
			Data: data,
		})
	} else {
		c.JSON(http.StatusOK, response{
			Code: 0,
			Msg:  "success",
		})
	}
}

// jsonErr
// return
//
//nolint:golint,unused
func jsonErr(c *gin.Context, err error, errMsg ...string) {
	if err == nil {
		err = errJsonErrNil
	}
	if len(errMsg) == 0 {
		c.JSON(http.StatusBadRequest, response{
			Code: 10,
			Msg:  fmt.Sprintf("%v", err),
		})
		return
	} else {
		message := strings.Join(errMsg, "; ")
		c.JSON(http.StatusBadRequest, response{
			Code: 10,
			Msg:  fmt.Sprintf("msg: %s err: %v", message, err),
		})
	}

}
