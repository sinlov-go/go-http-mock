package gin_mock

import (
	"github.com/gin-gonic/gin"
	"os"
)

// MockEngine
// for unit test
func MockEngine(init func(engine *gin.Engine) error) *gin.Engine {
	gin.SetMode(fetchGinRunMode())
	e := gin.Default()

	if init != nil {
		err := init(e)
		if err != nil {
			panic(err)
		}
	}

	return e
}

func fetchGinRunMode() string {
	ginMode := os.Getenv(gin.EnvGinMode)
	if ginMode == "" {
		ginMode = gin.TestMode
	}
	return ginMode
}
