package handler

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/neobarter/server/internal/pkg/response"
	"github.com/neobarter/server/internal/pkg/storage"
	"github.com/neobarter/server/internal/service"
)

type UploadHandler struct {
	uploadSvc *service.UploadService
}

func NewUploadHandler(uploadSvc *service.UploadService) *UploadHandler {
	return &UploadHandler{uploadSvc: uploadSvc}
}

// UploadImage 上传单张图片
// POST /v1/upload/image  (multipart/form-data, field: file)
func (h *UploadHandler) UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	// 大小预检查（避免读取超大文件进内存）
	if fileHeader.Size > storage.MaxImageSize {
		response.BadRequest(c, storage.ErrFileTooLarge.Error())
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.ServerError(c, "读取文件失败")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		response.ServerError(c, "读取文件失败")
		return
	}

	url, err := h.uploadSvc.UploadImage(data, fileHeader.Filename)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"url": url})
}
