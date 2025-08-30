// package handler 存放所有 API handlers。
package handler

import (
	"strconv"

	"github.com/KeLes-Coding/gopress/internal/api/response"
	"github.com/KeLes-Coding/gopress/internal/service"
	"github.com/gin-gonic/gin"
)

// TagHandler 结构体，用于挂载与标签相关的 API 方法。
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler 是 TagHandler 的构造函数。
func NewTagHandler() *TagHandler {
	return &TagHandler{
		tagService: service.NewTagService(),
	}
}

// CreateTagRequest 定义了创建标签接口的请求体。
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// CreateTagHandler 是处理创建标签请求的 Gin Handler。
func (h *TagHandler) CreateTagHandler(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("参数校验失败:"+err.Error(), c)
		return
	}
	tag, err := h.tagService.Create(req.Name)
	if err != nil {
		response.Error(err.Error(), c)
		return
	}
	response.Success(tag, c)
}

// ListTagsHandler 是处理获取标签列表请求的 Gin Handler。
func (h *TagHandler) ListTagsHandler(c *gin.Context) {
	tags, err := h.tagService.List()
	if err != nil {
		response.Error("获取标签列表失败："+err.Error(), c)
		return
	}
	response.Success(tags, c)
}

// UpdateTagRequest 定义了更新标签接口的请求体。
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// UpdateTagHandler 是处理更新标签请求的 Gin Handler。
func (h *TagHandler) UpdateTagHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error("无效的标签 ID", c)
		return
	}
	var req UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("参数校验失败"+err.Error(), c)
		return
	}
	updatedTag, err := h.tagService.Update(uint(id), req.Name)
	if err != nil {
		response.Error(err.Error(), c)
		return
	}
	response.Success(updatedTag, c)
}

// DeleteTagHandler 是处理删除标签请求的 Gin Handler。
func (h *TagHandler) DeleteTagHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error("无效的标签 ID", c)
		return
	}
	if err := h.tagService.Delete(uint(id)); err != nil {
		response.Error(err.Error(), c)
		return
	}
	response.Success(nil, c)
}
