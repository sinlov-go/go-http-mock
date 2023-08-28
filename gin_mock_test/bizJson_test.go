package gin_mock_test_test

import (
	"github.com/sebdah/goldie/v2"
	"github.com/sinlov-go/go-http-mock/gin_mock"
	"github.com/sinlov-go/go-http-mock/gin_mock_test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetJSON(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath
	// mock GetJSON
	tests := []struct {
		name     string
		path     string
		header   map[string]string
		respCode int
		wantErr  bool
	}{
		{
			name:     "sample", // testdata/TestGetJSON/sample.golden
			path:     "/Biz/json",
			respCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			// do GetJSON
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodGet).
				BodyJson(nil).
				Header(tc.header).
				NewRecorder()
			assert.False(t, tc.wantErr)
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify GetJSON
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}

func TestPostJsonModelBiz(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath
	// mock PostJsonModelBiz
	tests := []struct {
		name     string
		path     string
		header   map[string]string
		body     interface{}
		respCode int
		wantErr  bool
	}{
		{
			name: "sample", // testdata/TestPostJsonModelBiz/sample.golden
			path: "/Biz/modelBiz",
			body: gin_mock_test.Biz{
				Info:   "input info here",
				Id:     "foo",
				Offset: 1,
				Limit:  10,
			},
			respCode: http.StatusOK,
		},
		{
			name: "error model", // testdata/TestPostJsonModelBiz/sample.golden
			path: "/Biz/modelBiz",
			body: struct {
				Foo string `json:"foo"`
			}{
				Foo: "foo",
			},
			respCode: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			// do PostJsonModelBiz
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodPost).
				BodyJson(tc.body).
				Header(tc.header).
				NewRecorder()
			assert.False(t, tc.wantErr)
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify PostJsonModelBiz
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}

func TestPostQueryJsonMode(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath
	// mock PostQueryJsonMode
	tests := []struct {
		name     string
		path     string
		header   map[string]string
		query    interface{}
		body     interface{}
		respCode int
		wantErr  bool
	}{
		{
			name: "sample", // testdata/TestPostQueryJsonMode/sample.golden
			path: "/Biz/modelBizQuery",
			query: gin_mock_test.Biz{
				Offset: 1,
				Limit:  10,
			},
			body: gin_mock_test.Biz{
				Info: "input info here",
				Id:   "foo",
			},
			respCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			// do PostQueryJsonMode
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodPost).
				Query(tc.query).
				BodyJson(tc.body).
				Header(tc.header).
				NewRecorder()
			assert.False(t, tc.wantErr)
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify PostQueryJsonMode
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}
