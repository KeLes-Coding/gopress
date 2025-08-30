// package handler 存放所有 API handlers。
package handler

import (
	"strconv"

	"github.com/KeLes-Coding/gopress/internal/api/response"
	"github.com/KeLes-Coding/gopress/internal/service"
	"github.com/gin-gonic/gin"
)

// CategoryHandler 结构体，用于挂载与分类相关的 API 方法。
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler 是 CategoryHandler 的构造函数。
func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		categoryService: service.NewCategoryService(),
	}
}

// CreateCategoryRequest 定义了创建分类接口的请求体。
// 使用 binding tag 来进行参数校验，确保 name 字段存在且长度在 2 到 100 之间。
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// CreateCategoryHandler 是处理创建分类请求的 Gin Handler。
func (h *CategoryHandler) CreateCategoryHandler(c *gin.Context) {
	var req CreateCategoryRequest
	// 绑定并校验 JSON 请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("参数校验失败:"+err.Error(), c)
		return
	}

	// 调用 service 层来处理业务逻辑
	category, err := h.categoryService.Create(req.Name)
	if err != nil {
		response.Error(err.Error(), c)
		return
	}

	// 成功后，返回新创建的分类信息
	response.Success(category, c)
}

// ListCategoriesHandler 是处理获取分类列表请求的 Gin Handler。
func (h *CategoryHandler) ListCategoriesHandler(c *gin.Context) {
	categories, err := h.categoryService.List()
	if err != nil {
		response.Error("获取分类列表失败:"+err.Error(), c)
		return
	}

	response.Success(categories, c)
}

// UpdateCategoryRequest 定义了更新分类接口的请求体。
type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

// UpdateCategoryHandler 是处理更新分类请求的 Gin Handler。
func (h *CategoryHandler) UpdateCategoryHandler(c *gin.Context) {
	// 1. 从 URL 路径参数中分类 ID
	// c.Param("id") 返回的是字符串，需要转化为 uint 类型
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error("无效的分类 ID", c)
		return
	}

	// 2. 绑定并校验请求体
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("参数校验失败"+err.Error(), c)
		return
	}

	// 3. 调用 service 层处理更新逻辑
	updatedCategory, err := h.categoryService.Update(uint(id), req.Name)
	if err != nil {
		response.Error(err.Error(), c)
		return
	}

	// 4. 返回更新后的分类信息
	response.Success(updatedCategory, c)
}

// DeleteCategoryHandler 是处理删除分类请求的 Gin Handler。
func (h *CategoryHandler) DeleteCategoryHandler(c *gin.Context) {
	// 1. 从 URL 路径参数中获取分类 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error("无效的分类 ID", c)
		return
	}

	// 2. 调用 service 层处理删除逻辑
	if err := h.categoryService.Delete(uint(id)); err != nil {
		response.Error(err.Error(), c)
		return
	}

	// 3. 删除成功，返回成功的响应，通常 data 为 nil
	response.Success(nil, c)
}
