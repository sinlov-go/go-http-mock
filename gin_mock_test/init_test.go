package gin_mock_test_test

import (
	"github.com/gin-gonic/gin"
	"github.com/sinlov-go/go-http-mock/gin_mock"
	"github.com/sinlov-go/go-http-mock/gin_mock_test"
)

var (
	basePath    = "/api/v1"
	basicRouter *gin.Engine
)

func init() {
	basicRouter = setupTestRouter()
	gin_mock_test.Router(basicRouter, basePath)
}

func setupTestRouter() *gin.Engine {
	e := gin_mock.MockEngine(func(engine *gin.Engine) error {
		return nil
	})
	return e
}
