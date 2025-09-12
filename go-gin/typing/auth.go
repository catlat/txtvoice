package typing

type LoginSimpleReq struct {
	Identity string `form:"identity" binding:"required" label:"账号标识(邮箱/手机号)"`
}

type LoginUser struct {
	Identity    string `json:"identity"`
	DisplayName string `json:"display_name"`
	Status      int    `json:"status"`
}

type LoginSimpleReply struct {
	Token string    `json:"token"`
	User  LoginUser `json:"user"`
}

// PhoneLoginReq 手机号+密码登录
type PhoneLoginReq struct {
	Phone    string `form:"phone" binding:"required" label:"手机号"`
	Password string `form:"password" binding:"required,min=6" label:"密码"`
}

// ChangePasswordReq 修改密码
type ChangePasswordReq struct {
	NewPassword string `form:"new_password" binding:"required,min=6" label:"新密码"`
}
