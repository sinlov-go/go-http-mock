package gin_mock

import "fmt"

var (
	ErrMethodNotSupported = fmt.Errorf("method is not supported")
	ErrMIMENotSupported   = fmt.Errorf("mime is not supported")
	ErrNotSetTesting      = fmt.Errorf("please set testing.T")
)
