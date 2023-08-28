package gin_mock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

const (
	typeJson = "json"
	typeForm = "form"
)

// NewGinMock
//   - router: gin.Engine
//   - baseUrl: is the base url of the server, will trim the last "/"
//   - url: is the url of the request
func NewGinMock(t *testing.T, router *gin.Engine, baseUrl string, url string) GinMock {

	baseUrl = strings.TrimSuffix(baseUrl, "/")

	fullUrl := baseUrl + url

	return &ginMock{
		t:       t,
		router:  router,
		baseUrl: baseUrl,
		url:     url,
		fullUrl: fullUrl,
	}
}

type ginMock struct {
	t       *testing.T
	router  *gin.Engine
	baseUrl string
	url     string
	fullUrl string

	method string

	request *http.Request
}

func (g *ginMock) NewRecorder() *httptest.ResponseRecorder {
	if g.request == nil {
		if g.t != nil {
			g.t.Fatalf("please call Method Body() first")
		} else {
			panic(ErrNotSetTesting)
		}
	}
	recorder := httptest.NewRecorder()
	g.router.ServeHTTP(recorder, g.request)
	return recorder
}

func (g *ginMock) FullUrl() string {
	return g.fullUrl
}

// Method
//   - if method is not supported, and not set testing will panic
//
// please use this before Body()
func (g *ginMock) Method(method string) GinMock {
	switch method {
	default:
		if g.t != nil {
			g.t.Fatalf("mock request name %s method [ %s ] url %v error %v", g.t.Name(), method, g.fullUrl, ErrMethodNotSupported)
		} else {
			panic(ErrNotSetTesting)
		}
	case http.MethodGet:
		fallthrough
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodPatch:
		fallthrough
	case http.MethodHead:
		fallthrough
	case http.MethodOptions:
		fallthrough
	case http.MethodConnect:
		fallthrough
	case http.MethodTrace:
		g.method = method
	}
	return g
}

// Query
//   - if query is nil, will not add query string
//
// please use this before Body()
func (g *ginMock) Query(query interface{}) GinMock {
	if query != nil {
		g.fullUrl = fmt.Sprintf("%s?%s", g.fullUrl, mockQueryStrFrom(query))
	}
	return g
}

// Body
//   - if body is nil, will not set empty body
func (g *ginMock) Body(body io.Reader) GinMock {
	newRequest, err := http.NewRequest(g.method, g.fullUrl, body)
	if err != nil {
		if g.t != nil {
			g.t.Fatalf("mock request name %s method [ %s ] url %v error %v", g.t.Name(), g.method, g.fullUrl, err)
		} else {
			panic(ErrNotSetTesting)
		}
	}
	g.request = newRequest
	return g
}

// BodyForm
//   - if param is nil, will not set empty body
//   - if param is not struct or map, will not set empty body
//   - will auto add "Content-Type" header "application/x-www-form-urlencoded;charset=utf-8"
func (g *ginMock) BodyForm(param interface{}) GinMock {
	request, api, err := makeRequest(g.method, typeForm, g.fullUrl, param)
	if err != nil {
		if g.t != nil {
			g.t.Fatalf("mock makeRequest name %s method [ %s ] url %v error %v", g.t.Name(), g.method, g.fullUrl, err)
		} else {
			panic(ErrNotSetTesting)
		}
	}
	g.request = request
	g.fullUrl = api
	return g
}

// BodyJson
//   - if body is nil, will not set empty body
//   - if body is not struct or map, will not set empty body
//   - will auto add "Content-Type" header "application/json;charset=utf-8"
func (g *ginMock) BodyJson(body interface{}) GinMock {
	request, api, err := makeRequest(g.method, typeJson, g.fullUrl, body)
	if err != nil {
		if g.t != nil {
			g.t.Fatalf("mock makeRequest name %s method [ %s ] url %v error %v", g.t.Name(), g.method, g.fullUrl, err)
		} else {
			panic(ErrNotSetTesting)
		}
	}
	g.request = request
	g.fullUrl = api
	return g
}

func (g *ginMock) BodyFileForm(fileName string, fieldName string, param interface{}) GinMock {
	request, api, err := makeFileRequest(g.method, g.fullUrl, fileName, fieldName, param)
	if err != nil {
		if g.t != nil {
			g.t.Fatalf("mock makeFileRequest name %s method [ %s ] url %v error %v", g.t.Name(), g.method, g.fullUrl, err)
		} else {
			panic(ErrNotSetTesting)
		}
	}
	g.request = request
	g.fullUrl = api
	return g
}

// Header
//   - if request is nil, will call Body() first or BodyJson() first
func (g *ginMock) Header(header map[string]string) GinMock {
	if g.request == nil {
		if g.t != nil {
			g.t.Fatalf("please call Method Body() first")
		} else {
			panic(ErrNotSetTesting)
		}
	}
	if len(header) > 0 {
		for k, v := range header {
			g.request.Header.Add(k, v)
		}
	}
	return g
}

type GinMock interface {
	FullUrl() string

	Method(method string) GinMock

	Query(query interface{}) GinMock

	Body(body io.Reader) GinMock

	BodyForm(param interface{}) GinMock

	BodyJson(body interface{}) GinMock

	BodyFileForm(fileName string, fieldName string, param interface{}) GinMock

	Header(header map[string]string) GinMock

	NewRecorder() *httptest.ResponseRecorder
}

// make request
func makeRequest(method, mime, api string, param interface{}) (request *http.Request, finalApi string, err error) {
	method = strings.ToUpper(method)
	mime = strings.ToLower(mime)

	finalApi = api

	switch mime {
	case typeJson:
		var (
			contentBuffer *bytes.Buffer
			jsonBytes     []byte
		)
		jsonBytes, err = json.Marshal(param)
		if err != nil {
			return
		}
		contentBuffer = bytes.NewBuffer(jsonBytes)
		request, err = http.NewRequest(string(method), finalApi, contentBuffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/json;charset=utf-8")
	case typeForm:
		queryStr := mockQueryStrFrom(param)
		var buffer io.Reader

		if (method == http.MethodDelete || method == http.MethodGet) && queryStr != "" {
			finalApi += "?" + queryStr
		} else {
			buffer = bytes.NewReader([]byte(queryStr))
		}

		request, err = http.NewRequest(string(method), finalApi, buffer)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	default:
		err = ErrMIMENotSupported
		return
	}
	return
}

// makeFileRequest
// make request which contains uploading file
func makeFileRequest(method, api, fileName, fieldName string, param interface{}) (request *http.Request, finalApi string, err error) {
	method = strings.ToUpper(method)

	if method != http.MethodPost && method != http.MethodPut {
		err = ErrMethodNotSupported
		return
	}

	finalApi = api

	// create form file
	buf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(buf)
	fileWriter, err := bodyWriter.CreateFormFile(fieldName, fileName)
	if err != nil {
		return
	}

	// read the file
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return
	}

	// read the file to the fileWriter
	length, err := fileWriter.Write(fileBytes)
	if err != nil {
		return
	}

	errFileFormError := bodyWriter.Close()
	if errFileFormError != nil {
		return nil, finalApi, errFileFormError
	}

	// make request
	queryStr := mockQueryStrFrom(param)
	if queryStr != "" {
		finalApi += "?" + queryStr
	}
	request, err = http.NewRequest(string(method), finalApi, buf)
	if err != nil {
		return
	}

	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	err = request.ParseMultipartForm(int64(length))
	return
}

// mockQueryStrFrom
//
//	make query string from params
func mockQueryStrFrom(params interface{}) (result string) {
	if params == nil {
		return
	}
	value := reflect.ValueOf(params)

	switch value.Kind() {
	case reflect.Struct:
		var formName string
		for i := 0; i < value.NumField(); i++ {
			if formName = value.Type().Field(i).Tag.Get("form"); formName == "" {
				// don't tag the form name, use camel name
				formName = getCamelNameFrom(value.Type().Field(i).Name)
			}
			result += "&" + formName + "=" + fmt.Sprintf("%v", value.Field(i).Interface())
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			result += "&" + fmt.Sprintf("%v", key.Interface()) + "=" + fmt.Sprintf("%v", value.MapIndex(key).Interface())
		}
	default:
		return
	}

	if result != "" {
		result = result[1:]
	}
	return
}

// getCamelNameFrom
//
//	get the Camel name of the original name
func getCamelNameFrom(name string) string {
	result := ""
	i := 0
	j := 0
	r := []rune(name)
	for m, v := range r {
		// if the char is the capital
		if v >= 'A' && v < 'a' {
			// if the prior is the lower-case || if the prior is the capital and the latter is the lower-case
			if (m != 0 && r[m-1] >= 'a') || ((m != 0 && r[m-1] >= 'A' && r[m-1] < 'a') && (m != len(r)-1 && r[m+1] >= 'a')) {
				i = j
				j = m
				result += name[i:j] + "_"
			}
		}
	}

	result += name[j:]
	return strings.ToLower(result)
}
