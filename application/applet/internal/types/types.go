// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	Mobile           string `json:"mobile"`
	VerificationCode string `json:"verification_code"`
}

type LoginResponse struct {
	UserId int64 `json:"userId"`
	Token  Token `json:"token"`
}

type RegisterRequest struct {
	Name             string `json:"name"`
	Mobile           string `json:"mobile"`
	Password         string `json:"password"`
	VerificationCode string `json:"verification_code"`
}

type RegisterResponse struct {
	UserId int64 `json:"user_id"`
	Token  Token `json:"token"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	AccessExpire int64  `json:"access_expire"`
}

type UserInfoResponse struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type VerificationRequest struct {
	Mobile string `json:"mobile"`
}

type VerificationResponse struct {
}
