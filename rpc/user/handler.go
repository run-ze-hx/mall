package main

import (
	"context"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"golang.org/x/crypto/bcrypt"
	"mall/kitex_gen/user"
	"mall/model"
	"mall/rpc/user/userdao/mysql"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
}

// Register implements the UserServiceImpl interface.
// It registers a new user with the provided email and password.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {

	// 验证两次密码是否一致
	if req.Password != req.ConfirmPassword {
		return nil, kerrors.NewBizStatusError(500701, "password mismatch")
	}

	// 检查邮箱是否存在
	if User1, _ := mysql.GetUserByEmail(req.Email); User1 != nil {
		return nil, kerrors.NewBizStatusError(500702, "email already exists")
	}

	// 加密密码
	hashedPwd, _ := mysql.HashPassword(req.Password)

	// 创建用户
	newUser := &model.User{
		Email:        req.Email,
		PasswordHash: hashedPwd,
	}
	if err := mysql.CreateUser(newUser); err != nil {
		return nil, kerrors.NewBizStatusError(500703, err.Error())
	}

	//返回用户Id
	return &user.RegisterResp{UserId: newUser.Id}, nil

}

// Login implements the UserServiceImpl interface.
// user login with email and password
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	// 用邮箱查询用户
	user1, err1 := mysql.GetUserByEmail(req.Email)
	if err1 != nil {
		return nil, kerrors.NewBizStatusError(500704, err1.Error())
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user1.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, kerrors.NewBizStatusError(500705, err.Error())
	}

	return &user.LoginResp{UserId: user1.Id}, nil

}
