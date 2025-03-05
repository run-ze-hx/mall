package middleware

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"mall/kitex_gen/auth"
	"mall/kitex_gen/auth/authservice"
	"mall/utils"
	"strings"
)

func JwtAuthMiddleware(authClient authservice.Client) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 1. 提取 Token
		token := extractTokenFromHeader(ctx)
		if token == "" {
			ctx.AbortWithStatusJSON(401, utils.Response{Code: 401, Msg: "未提供 Token"})
			return
		}

		//  调用 Auth 微服务验证 Token,拿到userid
		verifyResp, err := authClient.VerifyTokenByRPC(c, &auth.VerifyTokenReq{Token: token})
		if err != nil || !verifyResp.Res {
			ctx.AbortWithStatusJSON(401, utils.Response{Code: 401, Msg: "无效 Token" + err.Error()})
			return
		}

		//  传递用户上下文
		ctx.Set("userID", verifyResp.UserId)
		ctx.Next(c)
	}
}

// 辅助函数：从 Header 提取 Token
func extractTokenFromHeader(ctx *app.RequestContext) string {
	authHeader := ctx.GetHeader("Authorization")
	if len(authHeader) > 7 && strings.EqualFold(string(authHeader[:7]), "Bearer ") {
		return string(authHeader[7:])
	}
	return ""
}
