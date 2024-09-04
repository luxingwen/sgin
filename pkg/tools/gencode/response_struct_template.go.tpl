
type {{.StructName}}InfoResponse struct {
	BaseResponse
	Data {{.ResStructName}} `json:"data"`
}

type {{.StructName}}QueryResponse struct {
	BasePageResponse
	Data []{{.ResStructName}} `json:"data"`
}
