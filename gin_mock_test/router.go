package gin_mock_test

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//nolint:golint,unused
func Router(g *gin.Engine, basePath string) {
	bizRouteGroup := g.Group(basePath + "/Biz")
	{
		bizRouteGroup.GET("/string", getString)
		bizRouteGroup.GET("/path/:some_id", getPath)
		bizRouteGroup.GET("/query", getQuery)
		bizRouteGroup.GET("/json", getJSON)

		// post
		bizRouteGroup.POST("/form", postForm)
		bizRouteGroup.POST("/modelBiz", postJsonModelBiz)
		bizRouteGroup.POST("/modelBizQuery", postQueryJsonMode)
		bizRouteGroup.POST("/upload", saveFileHandler)
	}
}

//nolint:golint,unused
func getString(c *gin.Context) {
	message := "this is Biz message"
	c.String(http.StatusOK, message)
}

//nolint:golint,unused
func getPath(c *gin.Context) {
	id := c.Param("some_id")
	if id == "" {
		jsonErr(c, nil, "id not found")
		return
	}
	resp := Biz{
		Id: id,
	}
	jsonSuccess(c, resp)
}

//nolint:golint,unused
func getQuery(c *gin.Context) {
	offset, limit, err := parseQueryCommonOffsetAndLimit(c)
	if err != nil {
		jsonErr(c, err)
		return
	}
	resp := Biz{
		Offset: offset,
		Limit:  limit,
	}
	jsonSuccess(c, resp)
}

//nolint:golint,unused
func getJSON(c *gin.Context) {
	resp := Biz{
		Info: "message",
	}
	jsonSuccess(c, struct {
		NewInfo string `json:"new_info"`
	}{NewInfo: resp.Info})
}

//nolint:golint,unused
func postForm(c *gin.Context) {
	c.GetHeader("")
	r := c.Request
	err := r.ParseForm()
	if err != nil {
		jsonErr(c, err, "Form parse error")
		return
	}
	formContent := make(map[string]string)
	for k, v := range r.PostForm {
		formContent[k] = strings.Join(v, "")
	}
	jsonSuccess(c, struct {
		PostFormContent map[string]string `json:"post_form_content,omitempty"`
	}{
		PostFormContent: formContent,
	})
}

//nolint:golint,unused
func postJsonModelBiz(c *gin.Context) {
	var req Biz
	if err := c.BindJSON(&req); err != nil {
		jsonErr(c, err)
		return
	}
	if req.Id == "" {
		jsonErr(c, nil, "id", "not found, set id and retry")
		return
	}
	c.JSON(http.StatusOK, req)
}

//nolint:golint,unused
func postQueryJsonMode(c *gin.Context) {
	offset, limit, err := parseQueryCommonOffsetAndLimit(c)
	if err != nil {
		jsonErr(c, err)
		return
	}
	var req Biz
	if errBind := c.BindJSON(&req); errBind != nil {
		jsonErr(c, err)
		return
	}
	req.Offset = offset
	req.Limit = limit
	c.JSON(http.StatusOK, req)
}
