package gin_mock_test_test

import (
	"github.com/sebdah/goldie/v2"
	"github.com/sinlov-go/go-http-mock/gin_mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSaveFileHandler(t *testing.T) {
	// mock gin at package test init()
	ginEngine := basicRouter
	apiBasePath := basePath
	// mock SaveFileHandler
	tests := []struct {
		name       string
		path       string
		fileName   string
		uploadName string
		header     map[string]string
		body       interface{}
		respCode   int
		wantErr    bool
	}{
		{
			name:       "sample", // testdata/TestSaveFileHandler/sample.golden
			path:       "/Biz/upload",
			fileName:   "test1.txt",
			uploadName: "test1",
			respCode:   http.StatusOK,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t,
				goldie.WithDiffEngine(goldie.ClassicDiff),
			)

			param := make(map[string]interface{})
			param["file_name"] = tc.fileName
			param["upload_name"] = tc.uploadName

			// do SaveFileHandler
			ginMock := gin_mock.NewGinMock(t, ginEngine, apiBasePath, tc.path)
			recorder := ginMock.
				Method(http.MethodPost).
				BodyFileSingleForm(tc.fileName, tc.uploadName, param).
				Header(tc.header).
				NewRecorder()
			if tc.wantErr {
				t.Logf("want err close check case %s", t.Name())
				return
			}
			// verify SaveFileHandler
			assert.Equal(t, tc.respCode, recorder.Code)
			g.Assert(t, t.Name(), recorder.Body.Bytes())
		})
	}
}
