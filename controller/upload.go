package controller

import (
	"net/http"
	"sgin/pkg/app"
)

type UploadController struct {
}

// 文件上传
// @Summary 文件上传
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "token"
// @Param file formData file true "文件"
// @Success 200 {string} app.Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /upload [post]
func (u *UploadController) UploadFile(ctx *app.Context) {
	// Multipart form
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Logger.Error("上传文件失败:", err)
		ctx.JSONError(http.StatusBadRequest, "上传文件失败")
		return
	}

	for _, files := range form.File {
		for _, file := range files {
			file.Filename = ctx.Config.Upload.Dir + "/" + file.Filename
			if err := ctx.SaveUploadedFile(file, file.Filename); err != nil {
				ctx.Logger.Error("上传文件失败:", err)
				ctx.JSONError(http.StatusBadRequest, "上传文件失败")
				return
			}
		}
	}

	ctx.JSONSuccess("上传文件成功")
}
