// package middleware 存放 Gin 的中间件。
package middleware

import (
	"strings"

	"github.com/KeLes-Coding/gopress/internal/api/response"
	"github.com/KeLes-Coding/gopress/internal/util"
	"github.com/gin-gonic/gin"
)

// CtxUserClaimsKey 是一个常量，用作 Gin Context 中存储用户 Claims 的键。
// 将其导出可以方便地在其他包（如 handler）中安全地引用，避免因手写字符串错误导致 bug。
const CtxUserClaimsKey = "userClaims"

// JWTAuthMiddleware 是一个 Gin 中间件，用于验证 JWT。
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Authorization 请求头中获取 token 字符串
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Unauthorized("请求未携带 token", c)
			return
		}

		// 2. 校验 token 格式
		// 一个标准的 token 格式是 "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized("Token 格式不正确", c)
			c.Abort()
			return
		}

		// 3. 解析并验证 token
		tokenString := parts[1]
		claims, err := util.ParseToken(tokenString)
		if err != nil {
			// 如果 ParseToken 返回错误，则认证失败
			response.Unauthorized("无效的 token", c)
			c.Abort()
			return
		}

		// 4. 将解析出的用户信息（claims）存入 Gin 的 Context
		// 这样，后续的 handler 就可以从 context 中获取到当前登录用户的信息
		c.Set(CtxUserClaimsKey, claims)

		// 5. 调用 c.Next() 将请求传递给下一个处理函数
		c.Next()
	}
}
