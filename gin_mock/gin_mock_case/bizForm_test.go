package gin_mock_case

import (
	"github.com/sebdah/goldie/v2"
	"github.com/sinlov-go/go-http-mock/gin_mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestPostForm(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath
	// mock PostForm
	tests := []struct {
		name     string
		path     string
		header   map[string]string
		body     interface{}
		respCode int
		wantErr  bool
	}{
		{
			name: "sample", // testdata/TestPostForm/sample.golden
			path: "/biz/form",
			body: biz{
				Info:   "input info here",
				Id:     "id123zqqeeadg24qasd",
				Offset: 0,
				Limit:  10,
			},
			respCode: http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			// do PostForm
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodPost).
				BodyForm(tc.body).
				Header(tc.header).
				NewRecorder()
			assert.False(t, tc.wantErr)
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify PostForm
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}
