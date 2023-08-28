package gin_mock_case

import (
	"github.com/gin-gonic/gin"
	"github.com/sinlov-go/go-http-mock/gin_mock"
)

var (
	basePath    = "/api/v1"
	basicRouter *gin.Engine
)

func init() {
	basicRouter = setupTestRouter()
	router(basicRouter, basePath)
}

func setupTestRouter() *gin.Engine {
	e := gin_mock.MockEngine(func(engine *gin.Engine) error {
		return nil
	})
	return e
}
