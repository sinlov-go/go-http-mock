package gin_mock_case

import (
	"github.com/sebdah/goldie/v2"
	"github.com/sinlov-go/go-http-mock/gin_mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetPath(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath
	// mock GetPath
	tests := []struct {
		name     string
		path     string
		header   map[string]string
		respCode int
		wantErr  bool
	}{
		{
			name:     "sample 123",
			path:     "/biz/path/123",
			respCode: http.StatusOK,
		},
		{
			name:     "sample 567",
			path:     "/biz/path/567",
			respCode: http.StatusOK,
		},
		{
			name:     "StatusNotFound",
			path:     "/biz/path/",
			respCode: http.StatusNotFound,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			// do GetPath
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodGet).
				Body(nil).
				Header(tc.header).
				NewRecorder()
			assert.False(t, tc.wantErr)
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify GetPath
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}

func TestGetQuery(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath

	type query struct {
		Offset string `form:"offset" json:"offset"`
		Limit  string `form:"limit" json:"limit"`
	}

	// mock GetQuery
	tests := []struct {
		name     string
		path     string
		query    any
		header   map[string]string
		respCode int
		wantErr  bool
	}{
		{
			name: "sample", // testdata/TestGetQuery/sample.golden
			path: "/biz/query",
			query: query{
				Offset: "0",
				Limit:  "10",
			},
			respCode: http.StatusOK,
		},
		{
			name: "fail offset", // testdata/TestGetQuery/sample.golden
			path: "/biz/query",
			query: query{
				Offset: "a",
				Limit:  "10",
			},
			respCode: http.StatusBadRequest,
		},
		{
			name: "fail limit", // testdata/TestGetQuery/sample.golden
			path: "/biz/query",
			query: query{
				Offset: "0",
				Limit:  "abc",
			},
			respCode: http.StatusBadRequest,
		},
		{
			name: "fail not exist url", // testdata/TestGetQuery/sample.golden
			path: "/biz/query/",
			query: query{
				Offset: "0",
				Limit:  "10",
			},
			respCode: http.StatusMovedPermanently,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			// do GetQuery
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodGet).
				BodyForm(tc.query).
				Header(tc.header).
				NewRecorder()
			assert.False(t, tc.wantErr)
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify GetQuery
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}

func TestGetString(t *testing.T) {
	ginMock := gin_mock.NewGinMock(t, basicRouter, basePath, "/biz/string")
	recorder := ginMock.
		Method(http.MethodGet).
		Body(nil).
		NewRecorder()
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "this is biz message", recorder.Body.String())
}
