package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"sgin/model"
	"sgin/pkg/app"
	"strings"

	"github.com/google/uuid"
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
// @Router /api/v1/upload [post]
func (u *UploadController) UploadFile(ctx *app.Context) {
	// Multipart form

	dstpath := getFilePath(ctx)

	const maxUploadSize = 100 * 1024 * 1024 // 100 MB
	err := ctx.Request.ParseMultipartForm(maxUploadSize)
	if err != nil {
		ctx.Logger.Error("上传文件失败:", err)
		ctx.JSONError(http.StatusBadRequest, "超过最大上传限制")
		return
	}

	ctx.Logger.Info("上传文件:", dstpath)

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.Logger.Error("上传文件失败:", err)
		ctx.JSONError(http.StatusBadRequest, "上传文件失败")
		return
	}

	fileattachments := make([]model.FileAttachment, 0)

	for _, files := range form.File {
		for _, file := range files {
			attatchment := model.FileAttachment{
				Name: file.Filename,
			}
			filename := uuid.New().String() + filepath.Ext(file.Filename)
			filename = filepath.Join(dstpath, filename)
			attatchment.Url = filename

			if err := ctx.SaveUploadedFile(file, filepath.Join(ctx.Config.Upload.Dir, filename)); err != nil {
				ctx.Logger.Error("上传文件失败:", err)
				ctx.JSONError(http.StatusBadRequest, "上传文件失败")
				return
			}
			fileattachments = append(fileattachments, attatchment)
		}
	}

	ctx.JSONSuccess(fileattachments)
}

func getFilePath(c *app.Context) (r string) {
	dstpath := c.Param("path")
	if !strings.HasPrefix(dstpath, "/") {
		dstpath = "/" + dstpath
	}
	dstpath = filepath.Clean(dstpath)
	return dstpath
}

// 删除文件
// @Summary 删除文件
// @Tags 上传
// @Accept json
// @Produce json
// @Param Authorization header string true "token"
// @Param param body model.ReqUuidParam true "文件路径"
// @Success 200 {string} app.Response "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/upload/delete [post]
func (u *UploadController) DeleteFile(ctx *app.Context) {
	var param model.ReqFileDeleteParam
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSONError(http.StatusBadRequest, err.Error())
		return
	}

	err := os.Remove(ctx.Config.Upload.Dir + param.Filename)
	if err != nil {
		ctx.Logger.Error("删除文件失败:", err)

		if os.IsNotExist(err) {
			ctx.Logger.Error("文件不存在:", param.Filename)
			ctx.JSONSuccess("删除文件成功")
			return
		}

		ctx.JSONError(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSONSuccess("删除文件成功")
}
