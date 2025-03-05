package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	auth "mall/kitex_gen/auth"
	"time"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	var token string
	token, err = generatetoken(req.GetUserId())
	if err != nil {
		log.Println("token generate error", err)
		return nil, kerrors.NewBizStatusError(500101, err.Error())
	}

	return &auth.DeliveryResp{Token: token}, nil
}

func generatetoken(userid int32) (string, error) {
	// 设置 JWT 令牌的过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建 JWT 声明
	claims := jwt.MapClaims{
		"userId": userid,
		"exp":    expirationTime.Unix(),
	}

	// 生成 JWT 令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {

	res, userid := verifytoken(req.GetToken())
	return &auth.VerifyResp{Res: res, UserId: userid}, nil
}

func verifytoken(token string) (res bool, userid int32) {
	claims := jwt.MapClaims{}

	tokenstr, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})
	if err != nil { //token解析失败
		return false, -1
	}

	if claims, ok := tokenstr.Claims.(jwt.MapClaims); ok && tokenstr.Valid {
		// 验证 exp 是否存在
		exp, ok := claims["exp"].(float64)
		if !ok || exp == 0 {
			return false, -1
		}

		// 验证是否过期
		if time.Now().Unix() > int64(exp) {
			return false, -1
		}
		// 获取 userId 并转换为 int32
		userIdFloat, ok := claims["userId"].(float64)
		if !ok {
			return false, -1
		}
		userid = int32(userIdFloat)
	}
	return true, userid
}
