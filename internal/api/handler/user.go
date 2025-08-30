// package handler 存放所有的 API handlers。
// handler 的主要职责是：
// 1. 解析和校验请求参数。
// 2. 调用相应的 service 层方法来处理业务逻辑。
// 3. 根据 service 的返回结果，使用 response 包来生成并返回 HTTP 响应。
package handler

import (
	"github.com/KeLes-Coding/gopress/internal/api/middleware"
	"github.com/KeLes-Coding/gopress/internal/api/response"
	"github.com/KeLes-Coding/gopress/internal/service"
	"github.com/KeLes-Coding/gopress/internal/util"
	"github.com/gin-gonic/gin"
)

// UserHandler 结构体，用于挂载与用户相关的 API 方法。
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler 是 UserHandler 的构造函数。
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// SignUpRequest 定义了用户注册接口的请求体结构。
// 使用 `binding:"required"` tag 来告诉 Gin 框架，这些字段是必需的，
// 如果请求中缺少这些字段，Gin 会自动返回一个错误。
type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignUpHandler 是处理用户注册请求的 Gin Handler。
func (h *UserHandler) SignUpHandler(c *gin.Context) {
	// 1. 绑定和校验请求参数
	var req SignUpRequest
	// c.ShouldBindJSON 会尝试将请求的 JSON body 解析到 req 结构体中。
	// 如果 JSON 格式错误或缺少 `binding:"required"` 的字段，会返回一个错误。
	if err := c.ShouldBindJSON(&req); err != nil {
		// 如果参数绑定失败，说明客户端请求格式有误，返回一个错误响应。
		response.Error("请求参数无效", c)
		return
	}

	// 2. 调用 service 层处理注册逻辑
	err := h.userService.SignUp(req.Username, req.Password)
	if err != nil {
		// 如果 service 层返回错误，将错误信息返回给客户端。
		response.Error(err.Error(), c)
		return
	}

	// 3. 注册成功，返回成功响应
	// 注册成功后，通常不返回具体数据，或者只返回一些基本信息（如用户ID），这里我们返回 nil。
	response.Success(nil, c)
}

// LoginRequest 定义了用户登录接口的请求体结构。
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 定义了用户登录成功后返回的数据结构。
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginHandler 是处理用户登录请求的 Gin Handler。
func (h *UserHandler) LoginHandler(c *gin.Context) {
	// 1. 绑定和校验请求参数
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error("请求参数无效", c)
		return
	}

	// 2. 调用 service 层处理登陆逻辑
	token, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		// 如果 service 返回错误，将其返回给客户端
		response.Error(err.Error(), c)
		return
	}

	// 3. 登录成功，返回 JWT
	response.Success(LoginResponse{Token: token}, c)
}

// GetMyProfileHandler 用于获取当前登录用户的信息。
// 这是一个受保护的接口，需要 JWT 认证。
func (h *UserHandler) GetMyProfileHandler(c *gin.Context) {
	// 从 Gin 上下文中获取由 JWTAuthMiddleware 中间件设置的用户 claims
	_claims, exists := c.Get(middleware.CtxUserClaimsKey)
	if !exists {
		// 如果 claims 不存在，通常意味着中间件逻辑有误，是一个服务端错误
		response.Error("无法获取用户信息", c)
		return
	}

	// 使用类型断言，将 an empty interface 转换为我们需要的 *util.MyClaims 类型
	claims, ok := _claims.(*util.MyClaims)
	if !ok {
		// 如果类型断言失败，也是一个服务端错误
		response.Error("用户信息类型错误", c)
		return
	}

	// 返回 claims 中的用户信息，构造一个简单的 map 作为 data
	response.Success(gin.H{
		"user_id":  claims.UserID,
		"username": claims.Username,
	}, c)
}
