package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"log"
	"mall/kitex_gen/auth"
	"mall/kitex_gen/user"
	"mall/kitex_gen/user/userservice"
	"mall/utils"
	"net/http"
)

var UserCli userservice.Client

// UserCliInit 创建 user 服务的 Client，并连接到 etcd
func UserCliInit(r discovery.Resolver) {

	uc, err1 := userservice.NewClient("user",
		client.WithResolver(r), // 使用 etcd 注册中心
	)
	if err1 != nil {
		log.Fatal(err1)
	}
	UserCli = uc
}

func UserRegister(c context.Context, ctx *app.RequestContext) {
	user1 := &user.RegisterReq{
		Email:           ctx.PostForm("email"),
		Password:        ctx.PostForm("password"),
		ConfirmPassword: ctx.PostForm("confirm_password"),
	}

	resp, err := UserCli.Register(c, user1)
	if err != nil {
		// 处理 RPC 调用错误

		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// 将 RPC 响应返回给 HTTP 客户端
	ctx.JSON(http.StatusOK, utils.Response{
		Code: 200,
		Data: map[string]interface{}{
			"user_id": resp.UserId,
		},
		Msg: "User registered successfully",
	})
}

func UserLogin(c context.Context, ctx *app.RequestContext) {
	user1 := &user.LoginReq{
		Email:    ctx.PostForm("email"),
		Password: ctx.PostForm("password"),
	}
	if user1.Email == "" || user1.Password == "" {
		ctx.JSON(http.StatusOK, map[string]string{"error": "email or password is empty"})
		return
	}

	resp, err := UserCli.Login(c, user1)
	if err != nil {
		// 处理 RPC 调用错误
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	//login 成功 ，发token令牌
	authreq := &auth.DeliverTokenReq{UserId: resp.UserId}
	authresp, autherr := AuthCli.DeliverTokenByRPC(c, authreq)
	if autherr != nil { //deliver Token error
		ctx.JSON(http.StatusInternalServerError, map[string]string{"error": autherr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, utils.Response{
		Code: 200,
		Data: map[string]interface{}{
			"Token": authresp.Token,
		},
		Msg: "User login successfully",
	})

}
