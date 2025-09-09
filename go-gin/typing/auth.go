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
