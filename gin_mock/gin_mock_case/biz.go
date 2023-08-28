package gin_mock_case

//nolint:golint,unused
type biz struct {
	Info   string `json:"info,omitempty" example:"input info here" default:"info"`
	Id     string `json:"id,omitempty"  example:"id123zqqeeadg24qasd" default:"id123zqqeeadg24qasd"`
	Offset int    `json:"offset,omitempty" example:"0" default:"0"`
	Limit  int    `json:"limit,omitempty" example:"10" default:"10"`
}

//nolint:golint,unused
type fileRequest struct {
	FileName   string `json:"file_name" form:"file_name" binding:"required"`
	UploadName string `json:"upload_name" form:"upload_name" binding:"required"`
}
