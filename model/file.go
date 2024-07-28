package model

type FileAttachment struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ReqFileDeleteParam struct {
	Filename string `json:"filename" binding:"required"`
}
