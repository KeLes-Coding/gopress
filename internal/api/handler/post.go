package handler

import (
	"strconv"

	"github.com/KeLes-Coding/gopress/internal/api/middleware"
	"github.com/KeLes-Coding/gopress/internal/api/response"
	"github.com/KeLes-Coding/gopress/internal/service"
	"github.com/KeLes-Coding/gopress/internal/util"
	"github.com/gin-gonic/gin"
)

// PostHandler 结构体...
type PostHandler struct {
	postService *service.PostService
}

// NewPostHandler 是 PostHandler 的构造函数。
func NewPostHandler() *PostHandler {
	return &PostHandler{
		postService: service.NewPostService(),
	}
}

// CreatePostRequest 定义了创建文章接口的请求体。
type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,min=2,max=255"`
	Content    string `json:"content" binding:"required,min=10"`
	Summary    string `json:"summary"`
	Status     *int   `json:"status" binding:"required,oneof=0 1"` // 使用指针以区分 0 和未提供
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
}

// CreatePostHandler ...
func (h *PostHandler) CreatePostHandler(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("参数校验失败: "+err.Error(), c)
		return
	}

	// 从 JWT claims 中获取当前登录用户的 ID
	_claims, _ := c.Get(middleware.CtxUserClaimsKey)
	claims := _claims.(*util.MyClaims)

	dto := &service.CreatePostDTO{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Status:     *req.Status,
		UserID:     claims.UserID,
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
	}

	post, err := h.postService.Create(dto)
	if err != nil {
		response.Error(err.Error(), c)
		return
	}

	response.Success(post, c)
}

// ListPostsHandler ...
func (h *PostHandler) ListPostsHandler(c *gin.Context) {
	// 从查询参数获取分页信息，并提供默认值
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	dto := &service.ListPostsDTO{
		Page:     page,
		PageSize: pageSize,
	}

	result, err := h.postService.List(dto)
	if err != nil {
		response.Error("获取文章列表失败: "+err.Error(), c)
		return
	}
	response.Success(result, c)
}

// GetPostHandler ...
func (h *PostHandler) GetPostHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error("无效的文章 ID", c)
		return
	}
	post, err := h.postService.GetByID(uint(id))
	if err != nil {
		response.Error(err.Error(), c)
		return
	}
	response.Success(post, c)
}

// UpdatePostRequest 定义了更新文章接口的请求体。
type UpdatePostRequest struct {
	Title      string `json:"title" binding:"required,min=2,max=255"`
	Content    string `json:"content" binding:"required,min=10"`
	Summary    string `json:"summary"`
	Status     *int   `json:"status" binding:"required,oneof=0 1"`
	CategoryID uint   `json:"category_id" binding:"required"`
	TagIDs     []uint `json:"tag_ids"`
}

// UpdatePostHandler ...
func (h *PostHandler) UpdatePostHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error("无效的文章 ID", c)
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("参数校验失败: "+err.Error(), c)
		return
	}

	dto := &service.UpdatePostDTO{
		ID:         uint(id),
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Status:     *req.Status,
		CategoryID: req.CategoryID,
		TagIDs:     req.TagIDs,
	}

	post, err := h.postService.Update(dto)
	if err != nil {
		response.Error(err.Error(), c)
		return
	}

	response.Success(post, c)
}

// DeletePostHandler ...
func (h *PostHandler) DeletePostHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.Error("无效的文章 ID", c)
		return
	}

	if err := h.postService.Delete(uint(id)); err != nil {
		response.Error(err.Error(), c)
		return
	}

	response.Success(nil, c)
}
