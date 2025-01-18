
[![go mod version](https://img.shields.io/github/go-mod/go-version/sinlov-go/go-http-mock?label=go.mod)](https://github.com/sinlov-go/go-http-mock)
[![GoDoc](https://godoc.org/github.com/sinlov-go/go-http-mock?status.png)](https://godoc.org/github.com/sinlov-go/go-http-mock)
[![goreportcard](https://goreportcard.com/badge/github.com/sinlov-go/go-http-mock)](https://goreportcard.com/report/github.com/sinlov-go/go-http-mock)

[![GitHub license](https://img.shields.io/github/license/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock/tags)
[![GitHub release)](https://img.shields.io/github/v/release/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock/releases)

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息


[![ci](https://github.com/sinlov-go/go-http-mock/actions/workflows/ci.yml/badge.svg)](https://github.com/sinlov-go/go-http-mock/actions/workflows/ci.yml)

[![go mod version](https://img.shields.io/github/go-mod/go-version/sinlov-go/go-http-mock?label=go.mod)](https://github.com/sinlov-go/go-http-mock)
[![GoDoc](https://godoc.org/github.com/sinlov-go/go-http-mock?status.png)](https://godoc.org/github.com/sinlov-go/go-http-mock)
[![goreportcard](https://goreportcard.com/badge/github.com/sinlov-go/go-http-mock)](https://goreportcard.com/report/github.com/sinlov-go/go-http-mock)

[![GitHub license](https://img.shields.io/github/license/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock)
[![codecov](https://codecov.io/gh/sinlov-go/go-http-mock/branch/main/graph/badge.svg)](https://codecov.io/gh/sinlov-go/go-http-mock)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock/tags)
[![GitHub release)](https://img.shields.io/github/v/release/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock/releases)

## for what

- this project used to http test case mock for golang

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/sinlov-go/go-http-mock)](https://github.com/sinlov-go/go-http-mock/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息

## depends

in go mod project

```bash
# warning use private git host must set
# global set for once
# add private git host like github.com to evn GOPRIVATE
$ go env -w GOPRIVATE='github.com'
# use ssh proxy
# set ssh-key to use ssh as http
$ git config --global url."git@github.com:".insteadOf "https://github.com/"
# or use PRIVATE-TOKEN
# set PRIVATE-TOKEN as gitlab or gitea
$ git config --global http.extraheader "PRIVATE-TOKEN: {PRIVATE-TOKEN}"
# set this rep to download ssh as https use PRIVATE-TOKEN
$ git config --global url."ssh://github.com/".insteadOf "https://github.com/"

# before above global settings
# test version info
$ git ls-remote -q https://github.com/sinlov-go/go-http-mock.git

# test depends see full version
$ go list -mod readonly -v -m -versions github.com/sinlov-go/go-http-mock
# or use last version add go.mod by script
$ echo "go mod edit -require=$(go list -mod=readonly -m -versions github.com/sinlov-go/go-http-mock | awk '{print $1 "@" $NF}')"
$ echo "go mod vendor"
```

## Features

- support [gin](https://github.com/gin-gonic/gin) mock test
- [ ] more perfect test case coverage
- [ ] more perfect benchmark case

## env

- minimum go version: go 1.19
- change `go 1.19`, `^1.19`, `1.19.12-bullseye`, `1.19.12` to new go version

### libs

| lib                                 | version |
|:------------------------------------|:--------|
| https://github.com/stretchr/testify | v1.8.4  |
| https://github.com/sebdah/goldie    | v2.5.3  |

- more libs see [go.mod](https://github.com/sinlov-go/go-http-mock/blob/main/go.mod)

## usage

### gin mock

- for [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin) use package `"github.com/sinlov-go/go-http-mock/gin_mock"`
- you can use `gin_mock.MockEngine` to init mock `gin.Engine`

```go
package gin_mock_test

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
	// bind gin router
	router(basicRouter, basePath)
}

func setupTestRouter() *gin.Engine {
	e := gin_mock.MockEngine(func(engine *gin.Engine) error {
		return nil
	})
	return e
}
```

- test case example

```go
package gin_mock_test_test

import (
	"github.com/sinlov-go/go-http-mock/gin_mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetString(t *testing.T) {
	ginMock := gin_mock.NewGinMock(t, basicRouter, basePath, "/Biz/string")
	recorder := ginMock.
		Method(http.MethodGet).
		Body(nil).
		NewRecorder()
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "this is Biz message", recorder.Body.String())
}
```

- more example see package `github.com/sinlov-go/go-http-mock/gin_mock_test`

# dev

```bash
# It needs to be executed after the first use or update of dependencies.
$ make init dep
```

- test code

```bash
$ make test testBenchmark
```

add main.go file and run

```bash
# run at env dev use cmd/main.go
$ make dev
```

- ci to fast check

```bash
# check style at local
$ make style

# run ci at local
$ make ci
```

## docker

```bash
# then test build as test/Dockerfile
$ make dockerTestRestartLatest
# clean test build
$ make dockerTestPruneLatest

# more info see
$ make helpDocker
```
